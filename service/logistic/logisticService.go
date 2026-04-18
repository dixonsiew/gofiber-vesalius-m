package logistic

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/dto"
    "vesaliusm/model"
    lg "vesaliusm/model/logistic"
    "vesaliusm/service/applicationUser"
    sqx "vesaliusm/sql"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    "github.com/nleeper/goment"
    go_ora "github.com/sijms/go-ora/v2"
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

func (s *LogisticService) SaveLogisticSetup(o lg.LogisticSetup, adminId int64) (int64, error) {
    query := `
        INSERT INTO LOGISTIC_ARRANGEMENT_SETUP
         (LOGISTIC_SETUP_CODE, LOGISTIC_SETUP_VALUE, USER_CREATE)
         VALUES 
         ('TERMS_AND_CONDITIONS', :logisticSetupValue, :adminId) 
        RETURNING LOGISTIC_SETUP_ID INTO :logistic_setup_id
    `
    var logisticSetupId go_ora.Number
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("logisticSetupValue", o.LogisticSetupValue.String),
        sql.Named("adminId", adminId),
        go_ora.Out{Dest: &logisticSetupId},
    )
    if err != nil {
        utils.LogError(err)
    }
    ilogisticSetupId, _ := logisticSetupId.Int64()
    return ilogisticSetupId, err
}

func (s *LogisticService) UpdateLogisticSetup(o lg.LogisticSetup, adminId int64) error {
    query := `
        UPDATE LOGISTIC_ARRANGEMENT_SETUP SET
          LOGISTIC_SETUP_VALUE = :logisticSetupValue,
          USER_UPDATE = :adminId,
          DATE_UPDATE = CURRENT_TIMESTAMP
        WHERE LOGISTIC_SETUP_ID = :logistic_setup_id
        AND LOGISTIC_SETUP_CODE = 'TERMS_AND_CONDITIONS'
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("logisticSetupValue", o.LogisticSetupValue.String),
        sql.Named("adminId", adminId),
        sql.Named("logistic_setup_id", o.LogisticSetupId.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *LogisticService) ExistsLogisticSetup() (bool, error) {
    query := `SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM LOGISTIC_ARRANGEMENT_SETUP WHERE LOGISTIC_SETUP_CODE = 'TERMS_AND_CONDITIONS')`
    var count int
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, err
}

func (s *LogisticService) FindLogisticSetup() (*lg.LogisticSetup, error) {
    query := `
        SELECT * FROM LOGISTIC_ARRANGEMENT_SETUP
        WHERE LOGISTIC_SETUP_CODE = 'TERMS_AND_CONDITIONS'
    `
    query = strings.Replace(query, "*", utils.GetDbCols(lg.LogisticSetup{}, ""), 1)
    var o lg.LogisticSetup
    err := s.db.GetContext(s.ctx, &o, query)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, err
}

func (s *LogisticService) SaveLogisticSlot(o lg.LogisticSlots, adminId int64) error {
    tx, err := s.db.BeginTxx(s.ctx, nil)
    if err != nil {
        utils.LogError(err)
        return err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()

    _, err = tx.ExecContext(s.ctx, `DELETE FROM LOGISTIC_ARRANGEMENT_SLOT`)
    if err != nil {
        utils.LogError(err)
        return err
    }

    query := `
        INSERT INTO LOGISTIC_ARRANGEMENT_SLOT
        (DAY_OF_WEEK, PICKUP_TIME, MAX_SLOTS, DISPLAY_SEQUENCE, USER_CREATE)
        VALUES
        (:dayOfWeek, :pickUpTime, :maxSlots, :displaySequence, :adminId)
    `
    for _, r := range o.LogisticSlots {
        _, err := tx.ExecContext(s.ctx, query,
            sql.Named("dayOfWeek", r.DayOfWeek.String),
            sql.Named("pickUpTime", r.PickUpTime.String),
            sql.Named("maxSlots", r.MaxSlots.Int32),
            sql.Named("displaySequence", r.DisplaySequence.Int32),
            sql.Named("adminId", adminId),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }

    return tx.Commit()
}

func (s *LogisticService) FindAllAppLogisticSlots(data *dto.LogisticSlotMobileDto) ([]lg.LogisticSlot, error) {
    list := make([]lg.LogisticSlot, 0)
    err := s.db.SelectContext(s.ctx, &list, sqx.GET_MOBILEAPP_SLOTS, 
        sql.Named("flightArrivalDate", data.FlightArrivalDate),
        sql.Named("flightArrivalTime", data.FlightArrivalTime),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        if list[i].AvailableSlots.Int32 <= 1 && data.WithCompanion == true {
            continue
        }

        if list[i].PickUpDate.Valid {
            g, _ := goment.New(list[i].PickUpDate.String, "YYYY-MM-DD")
            list[i].PickUpDate = utils.NewNullString(g.Format("DD/MM/YYYY"))
        }
    }
    return list, nil
}

func (s *LogisticService) FindAllLogisticSlots() ([]lg.LogisticSlot, error) {
    query := `
        SELECT * FROM LOGISTIC_ARRANGEMENT_SLOT
         ORDER BY 
          CASE 
            WHEN DAY_OF_WEEK = 'Monday' THEN 1
            WHEN DAY_OF_WEEK = 'Tuesday' THEN 2
            WHEN DAY_OF_WEEK = 'Wednesday' THEN 3
            WHEN DAY_OF_WEEK = 'Thursday' THEN 4
            WHEN DAY_OF_WEEK = 'Friday' THEN 5
            WHEN DAY_OF_WEEK = 'Saturday' THEN 6
            WHEN DAY_OF_WEEK = 'Sunday' THEN 7
        END
    `
    list := make([]lg.LogisticSlot, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *LogisticService) FindLogisticSlotBySlotId(slotId int64) (*lg.LogisticSlot, error) {
    query := `SELECT * FROM LOGISTIC_ARRANGEMENT_SLOT WHERE LOGISTIC_SLOT_ID = :slotId`
    query = strings.Replace(query, "*", utils.GetDbCols(lg.LogisticSlot{}, ""), 1)
    var o lg.LogisticSlot
    err := s.db.GetContext(s.ctx, &o, query, slotId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *LogisticService) SaveLogisticRequest(o lg.LogisticRequest) error {
    query := `
        INSERT INTO LOGISTIC_ARRANGEMENT_REQUESTER
        (
          LOGISTIC_REQUEST_NUMBER, REQUESTER_PRN, REQUESTER_NAME, REQUESTER_DOB,
          REQUESTER_DOC_TYPE, REQUESTER_DOC_NUMBER, REQUESTER_NATIONALITY, REQUESTER_EMAIL,
          PRIMARY_DOCTOR, VISIT_WITH_COMPANION, COMPANION_NAME, COMPANION_DOB,
          COMPANION_DOC_TYPE, COMPANION_DOC_NUMBER, RELATIONSHIP_TO_REQUESTER, FLIGHT_AIRLINE_NAME,
          FLIGHT_NUMBER, FLIGHT_ARRIVAL_DATE, FLIGHT_ARRIVAL_TIME, REQUESTED_PICKUP_DATE,
          REQUESTED_PICKUP_TIME, REQUESTED_PICKUP_DAY
        ) VALUES 
        (
          :logisticRequestNumber, :requesterPrn, :requesterName, TO_DATE(:requesterDob, 'DD/MM/YYYY'),
          :requesterDocType, :requesterDocNumber, :requesterNationality, :requesterEmail,
          :primaryDoctor, :visitWithCompanion, :companionName, TO_DATE(:companionDob, 'DD/MM/YYYY'),
          :companionDocType, :companionDocNumber, :relationshipToRequester, :flightAirlineName,
          :flightNumber, TO_DATE(:flightArrivalDate, 'DD/MM/YYYY'), :flightArrivalTime, TO_DATE(:requestedPickupDate, 'DD/MM/YYYY'),
          :requestedPickupTime, :requestedPickupDay
        )
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("logisticRequestNumber", o.LogisticRequestNumber.String),
        sql.Named("requesterPrn", o.RequesterPrn.String),
        sql.Named("requesterName", o.RequesterName.String),
        sql.Named("requesterDob", o.RequesterDob.String),
        sql.Named("requesterDocType", o.RequesterDocType.String),
        sql.Named("requesterDocNumber", o.RequesterDocNumber.String),
        sql.Named("requesterNationality", o.RequesterNationality.String),
        sql.Named("requesterEmail", o.RequesterEmail.String),
        sql.Named("primaryDoctor", o.PrimaryDoctor.String),
        sql.Named("visitWithCompanion", o.VisitWithCompanion.String),
        sql.Named("companionName", o.CompanionName.String),
        sql.Named("companionDob", o.CompanionDob.String),
        sql.Named("companionDocType", o.CompanionDocType.String),
        sql.Named("companionDocNumber", o.CompanionDocNumber.String),
        sql.Named("relationshipToRequester", o.RelationshipToRequester.String),
        sql.Named("flightAirlineName", o.FlightAirlineName.String),
        sql.Named("flightNumber", o.FlightNumber.String),
        sql.Named("flightArrivalDate", o.FlightArrivalDate.String),
        sql.Named("flightArrivalTime", o.FlightArrivalTime.String),
        sql.Named("requestedPickupDate", o.RequestedPickupDate.String),
        sql.Named("requestedPickupTime", o.RequestedPickupTime.String),
        sql.Named("requestedPickupDay", o.RequestedPickupDay.String),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *LogisticService) ListLogisticRequests(page string, limit string) (*model.PagedList, error) {
    total, err := s.CountLogisticRequests(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllLogisticRequests(pager.GetLowerBound(), pager.PageSize, s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *LogisticService) CountLogisticRequests(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(LOGISTIC_REQUEST_ID) AS COUNT FROM LOGISTIC_ARRANGEMENT_REQUESTER`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *LogisticService) ListAppLogisticRequests(userId int64, page string, limit string) (*model.PagedList, error) {
    user, err := s.applicationUserService.FindByUserId(userId, s.db)
    if err != nil {
        return nil, err
    }
    total, err := s.CountAppLogisticRequests(user.MasterPrn.String, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllAppLogisticRequests(user.MasterPrn.String, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *LogisticService) CountAppLogisticRequests(prn string, conn *sqlx.DB) (int, error) {
    query := `
        SELECT COUNT(LOGISTIC_REQUEST_ID) AS COUNT FROM LOGISTIC_ARRANGEMENT_REQUESTER
        WHERE REQUESTER_PRN = :prn OR
        REQUESTER_PRN IN (SELECT NOK_PRN FROM APPLICATION_USER_FAMILY WHERE PATIENT_PRN = :prn)
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, sql.Named("prn", prn))
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *LogisticService) FindAllAppLogisticRequests(prn string, offset int, limit int) ([]lg.LogisticRequest, error) {
    query := `
        SELECT lar.*
        FROM LOGISTIC_ARRANGEMENT_REQUESTER lar
        WHERE lar.REQUESTER_PRN = :prn OR 
        lar.REQUESTER_PRN IN (SELECT NOK_PRN FROM APPLICATION_USER_FAMILY WHERE PATIENT_PRN = :prn)
        ORDER BY lar.FLIGHT_ARRIVAL_DATE DESC, lar.LOGISTIC_REQUEST_ID DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "lar.*", utils.GetDbCols(lg.LogisticRequest{}, "lar."), 1)
    var list []lg.LogisticRequest
    err := s.db.SelectContext(s.ctx, &list, query,
        sql.Named("prn", prn),
        sql.Named("offset", offset),
        sql.Named("limit", limit),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
}

func (s *LogisticService) FindAllLogisticRequests(offset int, limit int, conn *sqlx.DB) ([]lg.LogisticRequest, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT lar.*
        FROM LOGISTIC_ARRANGEMENT_REQUESTER lar
        ORDER BY lar.FLIGHT_ARRIVAL_DATE
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "lar.*", utils.GetDbCols(lg.LogisticRequest{}, "lar."), 1)
    var list []lg.LogisticRequest
    err := db.SelectContext(s.ctx, &list, query,
        sql.Named("offset", offset),
        sql.Named("limit", limit),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].SetWebAdmin()
    }
    return list, nil
}

func (s *LogisticService) ListLogisticRequestsByKeyword(x dto.SearchKeyword4Dto, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountLogisticRequestsByKeyword(x, s.db)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindLogisticRequestsByKeyword(x, pager.GetLowerBound(), pager.PageSize, s.db)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *LogisticService) CountLogisticRequestsByKeyword(x dto.SearchKeyword4Dto, conn *sqlx.DB) (int, error) {
    conds, args := buildKeywordConditions(x)
    base := `SELECT COUNT(lar.LOGISTIC_REQUEST_ID) AS COUNT FROM LOGISTIC_ARRANGEMENT_REQUESTER lar`
    query := base + whereClause(conds)

    var count int
    err := s.db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *LogisticService) FindLogisticRequestsByKeyword(x dto.SearchKeyword4Dto, offset int, limit int, conn *sqlx.DB) ([]lg.LogisticRequest, error) {
    conditions, args := buildKeywordConditions(x)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    base := `
        SELECT lar.* 
        FROM LOGISTIC_ARRANGEMENT_REQUESTER lar
    `
    base = strings.Replace(base, "lar.*", utils.GetDbCols(lg.LogisticRequest{}, "lar."), 1)
    query := base + whereClause(conditions) +
        ` ORDER BY lar.FLIGHT_ARRIVAL_DATE OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    list := make([]lg.LogisticRequest, 0)
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].SetWebAdmin()
    }
    return list, nil
}

func (s *LogisticService) FindLogisticRequestByRequestId(requestId int64) (*lg.LogisticRequest, error) {
    query := `SELECT * FROM LOGISTIC_ARRANGEMENT_REQUESTER WHERE LOGISTIC_REQUEST_ID = :requestId`
    query = strings.Replace(query, "*", utils.GetDbCols(lg.LogisticRequest{}, ""), 1)
    var o lg.LogisticRequest
    err := s.db.GetContext(s.ctx, &o, query, requestId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    o.SetWebAdmin()
    return &o, nil
}

func (s *LogisticService) FindLogisticRequestByRequestNumber(requestNumber string) (*lg.LogisticRequest, error) {
    query := `SELECT * FROM LOGISTIC_ARRANGEMENT_REQUESTER WHERE LOGISTIC_REQUEST_NUMBER = :requestNumber`
    query = strings.Replace(query, "*", utils.GetDbCols(lg.LogisticRequest{}, ""), 1)
    var o lg.LogisticRequest
    err := s.db.GetContext(s.ctx, &o, query, requestNumber)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    o.SetWebAdmin()
    return &o, nil
}

func (s *LogisticService) UpdateLogisticRequestStatusByRequestNumber(requestNumber string, status string, adminId int64, userId int64) error {
    var err error
    if adminId > 0 {
        query := `
            UPDATE LOGISTIC_ARRANGEMENT_REQUESTER SET 
            LOGISTIC_REQUEST_STATUS = :status,
            ADMIN_UPDATE = :adminId,
            ADMIN_DATE_UPDATE = CURRENT_TIMESTAMP
            WHERE LOGISTIC_REQUEST_NUMBER = :requestNumber
        `
        _, err = s.db.ExecContext(s.ctx, query, 
            sql.Named("status", status), 
            sql.Named("adminId", adminId), 
            sql.Named("requestNumber", requestNumber),
        )
        if err != nil {
            utils.LogError(err)
        }
    } else if userId > 0 {
        query := `
            UPDATE LOGISTIC_ARRANGEMENT_REQUESTER SET 
            LOGISTIC_REQUEST_STATUS = :status,
            USER_UPDATE = :userId,
            USER_DATE_UPDATE = CURRENT_TIMESTAMP
            WHERE LOGISTIC_REQUEST_NUMBER = :requestNumber
        `
        _, err = s.db.ExecContext(s.ctx, query, 
            sql.Named("status", status), 
            sql.Named("userId", userId), 
            sql.Named("requestNumber", requestNumber),
        )
        if err != nil {
            utils.LogError(err)
        }
    }
    return err
}

func (s *LogisticService) FindAllRequestsForExcel() ([]lg.LogisticRequest, error) {
    query := `SELECT * FROM LOGISTIC_ARRANGEMENT_REQUESTER ORDER BY FLIGHT_ARRIVAL_DATE`
    query = strings.Replace(query, "*", utils.GetDbCols(lg.LogisticRequest{}, ""), 1)
    list := make([]lg.LogisticRequest, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].SetWebAdmin()
    }
    return list, nil
}

func (s *LogisticService) GenerateLogisticRequestNumber() (string, error) {
    query := `SELECT GEN_LOGISTIC_REQUEST_NUMBER() AS LOGISTIC_REQUEST_NUMBER FROM DUAL`
    var requestNumber string
    err := s.db.GetContext(s.ctx, &requestNumber, query)
    if err != nil {
        utils.LogError(err)
        return "", err
    }
    return requestNumber, nil
}

func buildKeywordConditions(x dto.SearchKeyword4Dto) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if x.Keyword != "" {
        conds = append(conds, `(LOWER(lar.REQUESTER_PRN) LIKE :keyword OR LOWER(lar.REQUESTER_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(x.Keyword)))
    }
    if x.Keyword2 != "" {
        conds = append(conds, `LOWER(lar.PRIMARY_DOCTOR) LIKE :keyword2`)
        args = append(args, sql.Named("keyword2", strings.ToLower(x.Keyword2)))
    }
    if x.Keyword3 != "" {
        conds = append(conds, `TRUNC(lar.REQUESTED_PICKUP_DATE) = TO_DATE(:keyword3, 'dd/mm/yyyy')`)
        args = append(args, sql.Named("keyword3", x.Keyword3))
    }
    if x.Keyword4 != "" {
        conds = append(conds, `LOWER(lar.LOGISTIC_REQUEST_STATUS) LIKE :keyword4`)
        args = append(args, sql.Named("keyword4", strings.ToLower(x.Keyword4)))
    }

    return conds, args
}

func whereClause(conds []string) string {
    if len(conds) == 0 {
        return ""
    }
    return " WHERE " + strings.Join(conds, " AND ")
}
