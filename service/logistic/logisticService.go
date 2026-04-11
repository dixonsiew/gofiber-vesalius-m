package logistic

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"vesaliusm/database"
	"vesaliusm/dto"
	"vesaliusm/model"
	logisticModel "vesaliusm/model/logistic"
	"vesaliusm/service/applicationUser"
	"vesaliusm/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/guregu/null/v6"
	"github.com/jmoiron/sqlx"
)

var LogisticSvc *LogisticService = NewLogisticService(database.GetDb(), database.GetCtx())

type LogisticService struct {
	db                     *sqlx.DB
	ctx                    context.Context
	applicationUserService *applicationUser.ApplicationUserService
}

func NewLogisticService(db *sqlx.DB, ctx context.Context) *LogisticService {
	return &LogisticService{
		db:                     db,
		ctx:                    ctx,
		applicationUserService: applicationUser.ApplicationUserSvc,
	}
}

func (s *LogisticService) ExistsSetup() (bool, error) {
	query := `SELECT COUNT(*) FROM DUAL WHERE EXISTS (SELECT 1 FROM LOGISTIC_ARRANGEMENT_SETUP WHERE LOGISTIC_SETUP_CODE = 'TERMS_AND_CONDITIONS')`
	var count int
	err := s.db.GetContext(s.ctx, &count, query)
	if err != nil {
		utils.LogError(err)
		return false, err
	}
	return count == 1, nil
}

func (s *LogisticService) CreateSetup(data *dto.LogisticSetupDto, adminID int64) (*logisticModel.LogisticSetup, error) {
	tx, err := s.db.BeginTxx(s.ctx, nil)
	if err != nil {
		utils.LogError(err)
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(s.ctx, `
		INSERT INTO LOGISTIC_ARRANGEMENT_SETUP (
			LOGISTIC_SETUP_CODE, LOGISTIC_SETUP_VALUE, USER_CREATE
		) VALUES (
			'TERMS_AND_CONDITIONS', :setupValue, :adminID
		)`,
		sql.Named("setupValue", data.LogisticSetupValue),
		sql.Named("adminID", adminID),
	)
	if err != nil {
		utils.LogError(err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		utils.LogError(err)
		return nil, err
	}

	return s.FindSetup()
}

func (s *LogisticService) UpdateSetup(setupID int64, data *dto.LogisticSetupDto, adminID int64) error {
	_, err := s.db.ExecContext(s.ctx, `
		UPDATE LOGISTIC_ARRANGEMENT_SETUP SET
			LOGISTIC_SETUP_VALUE = :setupValue,
			USER_UPDATE = :adminID,
			DATE_UPDATE = CURRENT_TIMESTAMP
		WHERE LOGISTIC_SETUP_ID = :setupID
		AND LOGISTIC_SETUP_CODE = 'TERMS_AND_CONDITIONS'`,
		sql.Named("setupValue", data.LogisticSetupValue),
		sql.Named("adminID", adminID),
		sql.Named("setupID", setupID),
	)
	if err != nil {
		utils.LogError(err)
	}

	return err
}

func (s *LogisticService) FindSetup() (*logisticModel.LogisticSetup, error) {
	query := `
		SELECT
			LOGISTIC_SETUP_ID,
			LOGISTIC_SETUP_CODE,
			LOGISTIC_SETUP_VALUE
		FROM LOGISTIC_ARRANGEMENT_SETUP
		WHERE LOGISTIC_SETUP_CODE = 'TERMS_AND_CONDITIONS'`

	var setup logisticModel.LogisticSetup
	err := s.db.GetContext(s.ctx, &setup, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		utils.LogError(err)
		return nil, err
	}
	return &setup, nil
}

func (s *LogisticService) ReplaceSlots(data *dto.LogisticSlotsDto, adminID int64) error {
	tx, err := s.db.BeginTxx(s.ctx, nil)
	if err != nil {
		utils.LogError(err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(s.ctx, `DELETE FROM LOGISTIC_ARRANGEMENT_SLOT`)
	if err != nil {
		utils.LogError(err)
		return err
	}

	query := `
		INSERT INTO LOGISTIC_ARRANGEMENT_SLOT (
			DAY_OF_WEEK, PICKUP_TIME, MAX_SLOTS, DISPLAY_SEQUENCE, USER_CREATE
		) VALUES (
			:dayOfWeek, :pickUpTime, :maxSlots, :displaySequence, :adminID
		)`

	for _, slot := range data.LogisticSlots {
		_, err = tx.ExecContext(s.ctx, query,
			sql.Named("dayOfWeek", slot.DayOfWeek),
			sql.Named("pickUpTime", slot.PickUpTime),
			sql.Named("maxSlots", slot.MaxSlots),
			sql.Named("displaySequence", slot.DisplaySequence),
			sql.Named("adminID", adminID),
		)
		if err != nil {
			utils.LogError(err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		utils.LogError(err)
	}

	return err
}

func (s *LogisticService) FindAppSlots(data *dto.LogisticSlotMobileDto) ([]logisticModel.LogisticSlot, error) {
	query := `
		SELECT
			DAY_OF_WEEK,
			PICK_UP_DATE,
			PICK_UP_TIME,
			AVAILABLE_SLOTS
		FROM (
			SELECT
				lars.DAY_OF_WEEK AS DAY_OF_WEEK,
				TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY'), 'DD/MM/YYYY') AS PICK_UP_DATE,
				lars.PICKUP_TIME AS PICK_UP_TIME,
				lars.MAX_SLOTS - COALESCE(SUM(
					CASE WHEN lar.VISIT_WITH_COMPANION = 'Y' THEN 2 ELSE 1 END * lar.REQUESTED_SLOTS
				), 0) AS AVAILABLE_SLOTS
			FROM LOGISTIC_ARRANGEMENT_SLOT lars
			LEFT JOIN (
				SELECT
					REQUESTED_PICKUP_DAY,
					REQUESTED_PICKUP_TIME,
					VISIT_WITH_COMPANION,
					COUNT(*) AS REQUESTED_SLOTS
				FROM LOGISTIC_ARRANGEMENT_REQUESTER
				WHERE REQUESTED_PICKUP_DATE >= TO_DATE(:flightArrivalDate, 'DD/MM/YYYY')
				AND REQUESTED_PICKUP_DATE <= TO_DATE(:flightArrivalDate, 'DD/MM/YYYY')
				AND LOGISTIC_REQUEST_STATUS IN ('Requested', 'Confirmed')
				GROUP BY REQUESTED_PICKUP_DAY, REQUESTED_PICKUP_TIME, VISIT_WITH_COMPANION
			) lar ON TRIM(lars.DAY_OF_WEEK) = TRIM(lar.REQUESTED_PICKUP_DAY)
				AND TRIM(lars.PICKUP_TIME) = TRIM(lar.REQUESTED_PICKUP_TIME)
			WHERE TRIM(lars.DAY_OF_WEEK) = TRIM(TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY'), 'Day'))
			AND lars.PICKUP_TIME > :flightArrivalTime
			GROUP BY lars.DAY_OF_WEEK, TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY'), 'DD/MM/YYYY'), lars.PICKUP_TIME, lars.MAX_SLOTS
			HAVING lars.MAX_SLOTS - COALESCE(SUM(
				CASE WHEN lar.VISIT_WITH_COMPANION = 'Y' THEN 2 ELSE 1 END * lar.REQUESTED_SLOTS
			), 0) > 0

			UNION

			SELECT
				lars.DAY_OF_WEEK AS DAY_OF_WEEK,
				TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY') + INTERVAL '1' DAY, 'DD/MM/YYYY') AS PICK_UP_DATE,
				lars.PICKUP_TIME AS PICK_UP_TIME,
				lars.MAX_SLOTS - COALESCE(SUM(
					CASE WHEN lar.VISIT_WITH_COMPANION = 'Y' THEN 2 ELSE 1 END * lar.REQUESTED_SLOTS
				), 0) AS AVAILABLE_SLOTS
			FROM LOGISTIC_ARRANGEMENT_SLOT lars
			LEFT JOIN (
				SELECT
					REQUESTED_PICKUP_DAY,
					REQUESTED_PICKUP_TIME,
					VISIT_WITH_COMPANION,
					COUNT(*) AS REQUESTED_SLOTS
				FROM LOGISTIC_ARRANGEMENT_REQUESTER
				WHERE REQUESTED_PICKUP_DATE >= TO_DATE(:flightArrivalDate, 'DD/MM/YYYY') + INTERVAL '1' DAY
				AND REQUESTED_PICKUP_DATE <= TO_DATE(:flightArrivalDate, 'DD/MM/YYYY') + INTERVAL '1' DAY
				AND LOGISTIC_REQUEST_STATUS IN ('Requested', 'Confirmed')
				GROUP BY REQUESTED_PICKUP_DAY, REQUESTED_PICKUP_TIME, VISIT_WITH_COMPANION
			) lar ON TRIM(lars.DAY_OF_WEEK) = TRIM(lar.REQUESTED_PICKUP_DAY)
				AND TRIM(lars.PICKUP_TIME) = TRIM(lar.REQUESTED_PICKUP_TIME)
			WHERE TRIM(lars.DAY_OF_WEEK) = TRIM(TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY') + INTERVAL '1' DAY, 'Day'))
			GROUP BY lars.DAY_OF_WEEK, TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY') + INTERVAL '1' DAY, 'DD/MM/YYYY'), lars.PICKUP_TIME, lars.MAX_SLOTS
			HAVING lars.MAX_SLOTS - COALESCE(SUM(
				CASE WHEN lar.VISIT_WITH_COMPANION = 'Y' THEN 2 ELSE 1 END * lar.REQUESTED_SLOTS
			), 0) > 0
		)
		ORDER BY PICK_UP_DATE, PICK_UP_TIME`

	slots := make([]logisticModel.LogisticSlot, 0)
	err := s.db.SelectContext(s.ctx, &slots, query,
        sql.Named("flightArrivalDate", data.FlightArrivalDate),
        sql.Named("flightArrivalTime", data.FlightArrivalTime),
	)
	if err != nil {
		utils.LogError(err)
		return slots, err
	}

	if data.WithCompanion {
		filtered := make([]logisticModel.LogisticSlot, 0, len(slots))
		for _, slot := range slots {
			if slot.AvailableSlots.Int32 <= 1 {
				continue
			}
			filtered = append(filtered, slot)
		}
		return filtered, nil
	}

	return slots, nil
}

func (s *LogisticService) FindAllSlots() ([]logisticModel.LogisticSlot, error) {
	query := `
		SELECT
			LOGISTIC_SLOT_ID,
			DAY_OF_WEEK,
			PICKUP_TIME AS PICK_UP_TIME,
			MAX_SLOTS,
			DISPLAY_SEQUENCE
		FROM LOGISTIC_ARRANGEMENT_SLOT
		ORDER BY
			CASE
				WHEN DAY_OF_WEEK = 'Monday' THEN 1
				WHEN DAY_OF_WEEK = 'Tuesday' THEN 2
				WHEN DAY_OF_WEEK = 'Wednesday' THEN 3
				WHEN DAY_OF_WEEK = 'Thursday' THEN 4
				WHEN DAY_OF_WEEK = 'Friday' THEN 5
				WHEN DAY_OF_WEEK = 'Saturday' THEN 6
				WHEN DAY_OF_WEEK = 'Sunday' THEN 7
				ELSE 8
			END,
			DISPLAY_SEQUENCE,
			PICKUP_TIME`

	slots := make([]logisticModel.LogisticSlot, 0)
	err := s.db.SelectContext(s.ctx, &slots, query)
	if err != nil {
		utils.LogError(err)
	}
	return slots, err
}

func (s *LogisticService) CreateRequest(data *dto.LogisticRequestDto) (*logisticModel.LogisticRequest, error) {
	requestDate, err := time.Parse("02/01/2006", data.RequestedPickupDate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Wrong date format")
	}

	requestNumber, err := s.GenerateRequestNumber()
	if err != nil {
		return nil, err
	}

	request := &logisticModel.LogisticRequest{
		LogisticRequestNumber:   null.NewString(requestNumber, true),
		RequesterPrn:            utils.NewNullString(data.RequesterPrn),
		RequesterName:           utils.NewNullString(data.RequesterName),
		RequesterDob:            utils.NewNullString(data.RequesterDob),
		RequesterDocType:        utils.NewNullString(data.RequesterDocType),
		RequesterDocNumber:      utils.NewNullString(data.RequesterDocNumber),
		RequesterNationality:    utils.NewNullString(data.RequesterNationality),
		RequesterEmail:          utils.NewNullString(data.RequesterEmail),
		PrimaryDoctor:           utils.NewNullString(data.PrimaryDoctor),
		VisitWithCompanion:      utils.NewNullString(data.VisitWithCompanion),
		CompanionName:           utils.NewNullString(data.CompanionName),
		CompanionDob:            utils.NewNullString(data.CompanionDob),
		CompanionDocType:        utils.NewNullString(data.CompanionDocType),
		CompanionDocNumber:      utils.NewNullString(data.CompanionDocNumber),
		RelationshipToRequester: utils.NewNullString(data.RelationshipToRequester),
		FlightAirlineName:       utils.NewNullString(data.FlightAirlineName),
		FlightNumber:            utils.NewNullString(data.FlightNumber),
		FlightArrivalDate:       utils.NewNullString(data.FlightArrivalDate),
		FlightArrivalTime:       utils.NewNullString(data.FlightArrivalTime),
		RequestedPickupDate:     utils.NewNullString(data.RequestedPickupDate),
		RequestedPickupTime:     utils.NewNullString(data.RequestedPickupTime),
		RequestedPickupDay:      utils.NewNullString(requestDate.Weekday().String()),
	}

	_, err = s.db.ExecContext(s.ctx, `
		INSERT INTO LOGISTIC_ARRANGEMENT_REQUESTER (
			LOGISTIC_REQUEST_NUMBER,
			REQUESTER_PRN,
			REQUESTER_NAME,
			REQUESTER_DOB,
			REQUESTER_DOC_TYPE,
			REQUESTER_DOC_NUMBER,
			REQUESTER_NATIONALITY,
			REQUESTER_EMAIL,
			PRIMARY_DOCTOR,
			VISIT_WITH_COMPANION,
			COMPANION_NAME,
			COMPANION_DOB,
			COMPANION_DOC_TYPE,
			COMPANION_DOC_NUMBER,
			RELATIONSHIP_TO_REQUESTER,
			FLIGHT_AIRLINE_NAME,
			FLIGHT_NUMBER,
			FLIGHT_ARRIVAL_DATE,
			FLIGHT_ARRIVAL_TIME,
			REQUESTED_PICKUP_DATE,
			REQUESTED_PICKUP_TIME,
			REQUESTED_PICKUP_DAY
		) VALUES (
			:requestNumber,
			:requesterPrn,
			:requesterName,
			TO_DATE(:requesterDob, 'DD/MM/YYYY'),
			:requesterDocType,
			:requesterDocNumber,
			:requesterNationality,
			:requesterEmail,
			:primaryDoctor,
			:visitWithCompanion,
			:companionName,
			TO_DATE(:companionDob, 'DD/MM/YYYY'),
			:companionDocType,
			:companionDocNumber,
			:relationshipToRequester,
			:flightAirlineName,
			:flightNumber,
			TO_DATE(:flightArrivalDate, 'DD/MM/YYYY'),
			:flightArrivalTime,
			TO_DATE(:requestedPickupDate, 'DD/MM/YYYY'),
			:requestedPickupTime,
			:requestedPickupDay
		)`,
		sql.Named("requestNumber", request.LogisticRequestNumber.String),
		sql.Named("requesterPrn", request.RequesterPrn.String),
		sql.Named("requesterName", request.RequesterName.String),
		sql.Named("requesterDob", request.RequesterDob.String),
		sql.Named("requesterDocType", request.RequesterDocType.String),
		sql.Named("requesterDocNumber", request.RequesterDocNumber.String),
		sql.Named("requesterNationality", request.RequesterNationality.String),
		sql.Named("requesterEmail", request.RequesterEmail.String),
		sql.Named("primaryDoctor", request.PrimaryDoctor.String),
		sql.Named("visitWithCompanion", request.VisitWithCompanion.String),
		sql.Named("companionName", request.CompanionName.String),
		sql.Named("companionDob", request.CompanionDob.String),
		sql.Named("companionDocType", request.CompanionDocType.String),
		sql.Named("companionDocNumber", request.CompanionDocNumber.String),
		sql.Named("relationshipToRequester", request.RelationshipToRequester.String),
		sql.Named("flightAirlineName", request.FlightAirlineName.String),
		sql.Named("flightNumber", request.FlightNumber.String),
		sql.Named("flightArrivalDate", request.FlightArrivalDate.String),
		sql.Named("flightArrivalTime", request.FlightArrivalTime.String),
		sql.Named("requestedPickupDate", request.RequestedPickupDate.String),
		sql.Named("requestedPickupTime", request.RequestedPickupTime.String),
		sql.Named("requestedPickupDay", request.RequestedPickupDay.String),
	)
	if err != nil {
		utils.LogError(err)
		return nil, err
	}
	return request, nil
}

func (s *LogisticService) ListAppRequests(userID int64, page string, limit string) (*model.PagedList, error) {
	user, err := s.applicationUserService.FindByUserId(userID, s.db)
	if err != nil {
		return nil, err
	}
	total, err := s.CountAppRequests(user.MasterPrn.String)
	if err != nil {
		return nil, err
	}
	pager := model.GetPager(total, page, limit)
	list, err := s.FindAllAppRequests(user.MasterPrn.String, pager.GetLowerBound(), pager.PageSize)
	if err != nil {
		return nil, err
	}
	return &model.PagedList{
		List: list, 
		Total: total, 
		TotalPages: pager.GetTotalPages(),
	}, nil
}

func (s *LogisticService) CountAppRequests(prn string) (int, error) {
	query := `
		SELECT COUNT(LOGISTIC_REQUEST_ID)
		FROM LOGISTIC_ARRANGEMENT_REQUESTER
		WHERE REQUESTER_PRN = :1
		OR REQUESTER_PRN IN (
			SELECT NOK_PRN FROM APPLICATION_USER_FAMILY WHERE PATIENT_PRN = :2
		)`
	var count int
	err := s.db.GetContext(s.ctx, &count, query, prn, prn)
	if err != nil {
		utils.LogError(err)
	}
	return count, err
}

func (s *LogisticService) FindAllAppRequests(prn string, offset int, limit int) ([]logisticModel.LogisticRequest, error) {
	query := `
		SELECT
			LOGISTIC_REQUEST_ID,
			LOGISTIC_REQUEST_NUMBER,
			LOGISTIC_REQUEST_STATUS,
			REQUESTER_NAME,
			NVL(PRIMARY_DOCTOR, '-') AS PRIMARY_DOCTOR_NAME,
			VISIT_WITH_COMPANION,
			NVL(COMPANION_NAME, '-') AS COMPANION_NAME,
			FLIGHT_AIRLINE_NAME,
			FLIGHT_NUMBER,
			TO_CHAR(FLIGHT_ARRIVAL_DATE, 'DD-MON-YYYY') AS FLIGHT_ARRIVAL_DATE,
			FLIGHT_ARRIVAL_TIME,
			TO_CHAR(REQUESTED_PICKUP_DATE, 'DD-MON-YYYY') AS REQUESTED_PICKUP_DATE,
			REQUESTED_PICKUP_TIME
		FROM LOGISTIC_ARRANGEMENT_REQUESTER
		WHERE REQUESTER_PRN = :1
		OR REQUESTER_PRN IN (
			SELECT NOK_PRN FROM APPLICATION_USER_FAMILY WHERE PATIENT_PRN = :2
		)
		ORDER BY FLIGHT_ARRIVAL_DATE DESC, LOGISTIC_REQUEST_ID DESC
		OFFSET :3 ROWS FETCH NEXT :4 ROWS ONLY`

	list := make([]logisticModel.LogisticRequest, 0)
	err := s.db.SelectContext(s.ctx, &list, query, prn, prn, offset, limit)
	if err != nil {
		utils.LogError(err)
		return list, err
	}
	for i := range list {
		list[i].SetApp()
	}
	return list, nil
}

func (s *LogisticService) List(page string, limit string) (*model.PagedList, error) {
	total, err := s.Count()
	if err != nil {
		return nil, err
	}
	pager := model.GetPager(total, page, limit)
	list, err := s.FindAll(pager.GetLowerBound(), pager.PageSize)
	if err != nil {
		return nil, err
	}
	return &model.PagedList{
		List: list, 
		Total: total, 
		TotalPages: pager.GetTotalPages(),
	}, nil
}

func (s *LogisticService) Count() (int, error) {
	var count int
	err := s.db.GetContext(s.ctx, &count, `SELECT COUNT(LOGISTIC_REQUEST_ID) FROM LOGISTIC_ARRANGEMENT_REQUESTER`)
	if err != nil {
		utils.LogError(err)
	}
	return count, err
}

func (s *LogisticService) FindAll(offset int, limit int) ([]logisticModel.LogisticRequest, error) {
	query := s.getWebAdminSelect() + `
		ORDER BY lar.FLIGHT_ARRIVAL_DATE
		OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY`

	list := make([]logisticModel.LogisticRequest, 0)
	err := s.db.SelectContext(s.ctx, &list, query, offset, limit)
	if err != nil {
		utils.LogError(err)
		return list, err
	}
	for i := range list {
		list[i].SetWebAdmin()
	}
	return list, nil
}

func (s *LogisticService) ListByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string, page string, limit string) (*model.PagedList, error) {
	total, err := s.CountByKeyword(keyword, keyword2, keyword3, keyword4)
	if err != nil {
		return nil, err
	}
	pager := model.GetPager(total, page, limit)
	list, err := s.FindByKeyword(keyword, keyword2, keyword3, keyword4, pager.GetLowerBound(), pager.PageSize)
	if err != nil {
		return nil, err
	}
	return &model.PagedList{
		List: list, 
		Total: total, 
		TotalPages: pager.GetTotalPages(),
	}, nil
}

func (s *LogisticService) CountByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string) (int, error) {
	whereClause, args := buildLogisticSearchConditions(keyword, keyword2, keyword3, keyword4)
	query := `SELECT COUNT(lar.LOGISTIC_REQUEST_ID) FROM LOGISTIC_ARRANGEMENT_REQUESTER lar` + whereClause
	var count int
	err := s.db.GetContext(s.ctx, &count, query, args...)
	if err != nil {
		utils.LogError(err)
	}

	return count, err
}

func (s *LogisticService) FindByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string, offset int, limit int) ([]logisticModel.LogisticRequest, error) {
	whereClause, args := buildLogisticSearchConditions(keyword, keyword2, keyword3, keyword4)
	query := s.getWebAdminSelect() + whereClause + `
		ORDER BY lar.FLIGHT_ARRIVAL_DATE
		OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
	args = append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	list := make([]logisticModel.LogisticRequest, 0)
	err := s.db.SelectContext(s.ctx, &list, query, args...)
	if err != nil {
		utils.LogError(err)
		return list, err
	}
	for i := range list {
		list[i].SetWebAdmin()
	}
	return list, nil
}

func (s *LogisticService) FindByRequestID(requestID int64) (*logisticModel.LogisticRequest, error) {
	query := s.getWebAdminSelect() + ` WHERE lar.LOGISTIC_REQUEST_ID = :1`
	var request logisticModel.LogisticRequest
	err := s.db.GetContext(s.ctx, &request, query, requestID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		utils.LogError(err)
		return nil, err
	}
	request.SetWebAdmin()
	return &request, nil
}

func (s *LogisticService) UpdateRequestStatusByNumber(requestNumber string, status string, adminID int64, userID int64) error {
	query := `
		UPDATE LOGISTIC_ARRANGEMENT_REQUESTER SET
			LOGISTIC_REQUEST_STATUS = :status,
			%s = :updatedBy,
			%s = CURRENT_TIMESTAMP
		WHERE LOGISTIC_REQUEST_NUMBER = :requestNumber`

	field := "USER_UPDATE"
	fieldDate := "USER_DATE_UPDATE"
	updatedBy := userID
	if adminID > 0 {
		field = "ADMIN_UPDATE"
		fieldDate = "ADMIN_DATE_UPDATE"
		updatedBy = adminID
	}
	_, err := s.db.ExecContext(s.ctx, fmt.Sprintf(query, field, fieldDate),
		sql.Named("status", status),
		sql.Named("updatedBy", updatedBy),
		sql.Named("requestNumber", requestNumber),
	)
	if err != nil {
		utils.LogError(err)
	}
	return err
}

func (s *LogisticService) FindByRequestNumberForMail(requestNumber string) (*logisticModel.LogisticRequest, error) {
	query := `
		SELECT
			LOGISTIC_REQUEST_NUMBER,
			NVL(REQUESTER_PRN, '-') AS REQUESTER_PRN,
			REQUESTER_NAME,
			VISIT_WITH_COMPANION,
			NVL(COMPANION_NAME, '-') AS COMPANION_NAME,
			TO_CHAR(REQUESTED_PICKUP_DATE, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS REQUESTED_PICKUP_DATE,
			REQUESTED_PICKUP_TIME
		FROM LOGISTIC_ARRANGEMENT_REQUESTER
		WHERE LOGISTIC_REQUEST_NUMBER = :requestNumber`
	var request logisticModel.LogisticRequest
	err := s.db.GetContext(s.ctx, &request, query, requestNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		utils.LogError(err)
		return nil, err
	}
	return &request, nil
}

func (s *LogisticService) ExportAll() ([]logisticModel.LogisticRequest, error) {
	query := s.getWebAdminSelect() + ` ORDER BY lar.FLIGHT_ARRIVAL_DATE`
	list := make([]logisticModel.LogisticRequest, 0)
	err := s.db.SelectContext(s.ctx, &list, query)
	if err != nil {
		utils.LogError(err)
		return list, err
	}
	for i := range list {
		list[i].SetWebAdmin()
	}
	return list, nil
}

func (s *LogisticService) ExportByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string) ([]logisticModel.LogisticRequest, error) {
	whereClause, args := buildLogisticSearchConditions(keyword, keyword2, keyword3, keyword4)
	query := s.getWebAdminSelect() + whereClause + ` ORDER BY lar.FLIGHT_ARRIVAL_DATE`
	list := make([]logisticModel.LogisticRequest, 0)
	err := s.db.SelectContext(s.ctx, &list, query, args...)
	if err != nil {
		utils.LogError(err)
		return list, err
	}
	for i := range list {
		list[i].SetWebAdmin()
	}
	return list, nil
}

func (s *LogisticService) GenerateRequestNumber() (string, error) {
	var requestNumber string
	err := s.db.GetContext(s.ctx, &requestNumber, `SELECT GEN_LOGISTIC_REQUEST_NUMBER() FROM DUAL`)
	if err != nil {
		utils.LogError(err)
	}
	return requestNumber, err
}

func (s *LogisticService) getWebAdminSelect() string {
	return `
		SELECT
			lar.LOGISTIC_REQUEST_ID,
			lar.LOGISTIC_REQUEST_NUMBER,
			lar.LOGISTIC_REQUEST_STATUS,
			NVL(lar.REQUESTER_PRN, '-') AS REQUESTER_PRN,
			lar.REQUESTER_NAME,
			TO_CHAR(lar.REQUESTER_DOB, 'DD/MM/YYYY') AS REQUESTER_DOB,
			lar.REQUESTER_DOC_TYPE,
			lar.REQUESTER_DOC_NUMBER,
			lar.REQUESTER_NATIONALITY,
			lar.REQUESTER_EMAIL,
			NVL(lar.PRIMARY_DOCTOR, '-') AS PRIMARY_DOCTOR,
			lar.VISIT_WITH_COMPANION,
			NVL(lar.COMPANION_NAME, '-') AS COMPANION_NAME,
			TO_CHAR(lar.COMPANION_DOB, 'DD/MM/YYYY') AS COMPANION_DOB,
			lar.COMPANION_DOC_TYPE,
			NVL(lar.COMPANION_DOC_NUMBER, '-') AS COMPANION_DOC_NUMBER,
			lar.RELATIONSHIP_TO_REQUESTER,
			lar.FLIGHT_AIRLINE_NAME,
			lar.FLIGHT_NUMBER,
			TO_CHAR(lar.FLIGHT_ARRIVAL_DATE, 'DD/MM/YYYY') AS FLIGHT_ARRIVAL_DATE,
			lar.FLIGHT_ARRIVAL_TIME,
			TO_CHAR(lar.REQUESTED_PICKUP_DATE, 'DD/MM/YYYY') AS REQUESTED_PICKUP_DATE,
			lar.REQUESTED_PICKUP_TIME,
			lar.REQUESTED_PICKUP_DAY,
			TO_CHAR(lar.DATE_CREATE, 'DD/MM/YYYY HH24:MI') AS DATE_CREATE,
			NVL(TO_CHAR(lar.USER_UPDATE), '-') AS USER_UPDATE,
			NVL(TO_CHAR(lar.USER_DATE_UPDATE, 'DD/MM/YYYY HH24:MI'), '-') AS USER_DATE_UPDATE,
			NVL(TO_CHAR(lar.ADMIN_UPDATE), '-') AS ADMIN_UPDATE,
			NVL(TO_CHAR(lar.ADMIN_DATE_UPDATE, 'DD/MM/YYYY HH24:MI'), '-') AS ADMIN_DATE_UPDATE
		FROM LOGISTIC_ARRANGEMENT_REQUESTER lar`
}

func buildLogisticSearchConditions(keyword string, keyword2 string, keyword3 string, keyword4 string) (string, []any) {
	conds := make([]string, 0)
	args := make([]any, 0)

	if keyword != "" {
		conds = append(conds, `(LOWER(lar.REQUESTER_PRN) LIKE :keyword OR LOWER(lar.REQUESTER_NAME) LIKE :keyword)`)
		args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
	}
	if keyword2 != "" {
		conds = append(conds, `LOWER(lar.PRIMARY_DOCTOR) LIKE :keyword2`)
		args = append(args, sql.Named("keyword2", strings.ToLower(keyword2)))
	}
	if keyword3 != "" {
		conds = append(conds, `TRUNC(lar.REQUESTED_PICKUP_DATE) = TO_DATE(:keyword3, 'dd/mm/yyyy')`)
		args = append(args, sql.Named("keyword3", keyword3))
	}
	if keyword4 != "" {
		conds = append(conds, `LOWER(lar.LOGISTIC_REQUEST_STATUS) LIKE :keyword4`)
		args = append(args, sql.Named("keyword4", strings.ToLower(keyword4)))
	}

	if len(conds) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conds, " AND "), args
}
