package novaDoctorPatientAppt

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/model/vesaliusGeo"
    sqx "vesaliusm/sql"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var NovaDoctorPatientApptSvc *NovaDoctorPatientApptService = NewNovaDoctorPatientApptService(database.GetDb(), database.GetCtx())

type NovaDoctorPatientApptService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaDoctorPatientApptService(db *sqlx.DB, ctx context.Context) *NovaDoctorPatientApptService {
    return &NovaDoctorPatientApptService{db: db, ctx: ctx}
}

func (s *NovaDoctorPatientApptService) List(page string, limit string) (*model.PagedList, error) {
    total, err := s.Count(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAll(pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *NovaDoctorPatientApptService) Count(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    var count int
    query := `SELECT COUNT(*) AS COUNT FROM NOVA_DOCTOR_PATIENT_APPT`
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *NovaDoctorPatientApptService) FindAll(offset int, limit int, conn *sqlx.DB) ([]model.DoctorPatientAppointment, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT * FROM NOVA_DOCTOR_PATIENT_APPT ndpa
        ORDER BY DATE_APPT DESC 
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.DoctorPatientAppointment{}, "ndpa."), 1)
    list := make([]model.DoctorPatientAppointment, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaDoctorPatientApptService) FindAllByDoctorId(doctorId int64, month int, year int, needAppt string) ([]model.NovaDoctorAppointment, []model.NovaDoctorAppointmentLists, error) {
    monthArray := []string{"EMPTY", "JAN-", "FEB-", "MAR-", "APR-", "MAY-", "JUN-", "JUL-", "AUG-", "SEP-", "OCT-", "NOV-", "DEC-"}
    monthYear := fmt.Sprintf("%s%d", monthArray[month], year)
    query := sqx.GET_ALL_DOCTOR_APPOINTMENTS
    list := make([]model.NovaDoctorAppointment, 0)
    lx := make([]model.NovaDoctorAppointmentLists, 0)
    err := s.db.SelectContext(s.ctx, &lx, query, sql.Named("doctorId", doctorId), sql.Named("monthYear", monthYear))
    if err != nil {
        utils.LogError(err)
        return nil, nil, err
    }

    if needAppt == "1" {
        query := `SELECT DAY_OF_WEEK, SLOT_TYPE, SESSION_TYPE, START_TIME, END_TIME FROM NOVA_DOCTOR_APPT_SLOT WHERE DOCTOR_ID = :doctorId`
        err := s.db.SelectContext(s.ctx, &list, query, doctorId)
        if err != nil {
            utils.LogError(err)
            return nil, nil, err
        }
    }

    return list, lx, nil
}

func (s *NovaDoctorPatientApptService) ListByKeyword(keyword string, keyword2 string, keyword3 string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountByKeyword(keyword, keyword2, keyword3, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindByKeyword(keyword, keyword2, keyword3, pager.GetLowerBound(), pager.PageSize, s.db)
    if err != nil {
        return nil, err
    }
    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *NovaDoctorPatientApptService) CountByKeyword(keyword string, keyword2 string, keyword3 string, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    conds, args := buildKeywordConditions(keyword, keyword2, keyword3)
    base := `SELECT COUNT(*) AS COUNT FROM NOVA_DOCTOR_PATIENT_APPT ndpa`
    query := base + whereClause(conds)

    var count int
    err := db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *NovaDoctorPatientApptService) FindByKeyword(keyword string, keyword2 string, keyword3 string, offset int, limit int, conn *sqlx.DB) ([]model.DoctorPatientAppointment, error) {
    db := database.GetFromCon(conn, s.db)
    conditions, args := buildKeywordConditions(keyword, keyword2, keyword3)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    base := `SELECT * FROM NOVA_DOCTOR_PATIENT_APPT ndpa`
    query := base + whereClause(conditions) +
        ` ORDER BY DATE_APPT DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.DoctorPatientAppointment{}, "ndpa."), 1)

    list := make([]model.DoctorPatientAppointment, 0)
    err := db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaDoctorPatientApptService) ExistsByPrnDateTime(prn string, dateTime string, conn *sqlx.DB) (bool, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS
        (
          SELECT 1 FROM NOVA_DOCTOR_PATIENT_APPT
          WHERE APPT_STATUS != 'CANCELLED'
          AND PATIENT_PRN = :prn
          AND DATE_APPT = TO_DATE(:dt, 'DD-MON-YYYY hh24:mi')
        )
    `
    var count int
    err := db.GetContext(s.ctx, &count, query, prn, dateTime)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, nil
}

func (s *NovaDoctorPatientApptService) ExistsByPrnDateSessionType(prn string, date string, sessionType string) (bool, error) {
    query := `
        SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS
        (
          SELECT 1 FROM NOVA_DOCTOR_PATIENT_APPT 
          WHERE APPT_SESSIONTYPE = :sessionType 
          AND APPT_STATUS != 'CANCELLED' 
          AND PATIENT_PRN = :prn
          AND DATE_APPT LIKE TO_DATE(:dt, 'YYYY-MM-DD')
        )
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, sessionType, prn, date)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, nil
}

func (s *NovaDoctorPatientApptService) FindApptSlotsByDoctorId(doctorId int64) ([]model.NovaDoctorApptSlot, error) {
    query := `SELECT * FROM NOVA_DOCTOR_APPT_SLOT WHERE DOCTOR_ID = :doctorId`
    list := make([]model.NovaDoctorApptSlot, 0)
    err := s.db.SelectContext(s.ctx, &list, query, doctorId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaDoctorPatientApptService) Save(prn string, apptSessionType string, doctorId int64, rd vesaliusGeo.AppointmentBookingConfirmation, remark string, tx *sqlx.Tx) error {
    patientName, err := s.getPatientName(s.db, prn)
    if err != nil {
        return err
    }

    rem := &remark
    if remark == "" {
        rem = nil
    }

    query := `
        INSERT INTO NOVA_DOCTOR_PATIENT_APPT (
         DOCTOR_ID, DOCTOR_NAME, DOCTOR_SPECIALTY, PATIENT_PRN, PATIENT_NAME, 
         APPT_STATUS, APPT_NO, APPT_DAY, APPT_SESSIONTYPE, APPT_CLINIC, 
         APPT_ROOM, APPT_CASETYPE, DATE_APPT, PACKAGE_PURCHASE_NO
         ) VALUES (
         :doctor_id, :doctor_name, :doctor_specialty, :patient_prn, :patient_name, 
         :appointment_status, :appointment_number, :appointment_day, :appointment_sessionType, :appointment_clinic, 
         :appointment_room, :appointment_caseType, TO_DATE(:appointment_dateTime, 'DD-MON-YYYY hh24:mi'), :purchase_no
        )
    `
    args := []any{
        sql.Named("doctor_id", doctorId),
        sql.Named("doctor_name", rd.DoctorName),
        sql.Named("doctor_specialty", rd.Specialty),
        sql.Named("patient_prn", prn),
        sql.Named("patient_name", patientName),
        sql.Named("appointment_status", rd.AppointmentStatus),
        sql.Named("appointment_number", rd.AppointmentNumber),
        sql.Named("appointment_day", rd.Day),
        sql.Named("appointment_sessionType", apptSessionType),
        sql.Named("appointment_clinic", rd.Clinic),
        sql.Named("appointment_room", rd.Room),
        sql.Named("appointment_caseType", rd.CaseType),
        sql.Named("appointment_dateTime", fmt.Sprintf("%s %s", rd.Date, rd.StartTime)),
        sql.Named("purchase_no", rem),
    }
    if tx == nil {
        _, err = s.db.ExecContext(s.ctx, query, args...)
    } else {
        _, err = tx.ExecContext(s.ctx, query, args...)
    }
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *NovaDoctorPatientApptService) UpdateToReschedule(prn string, apptSessionType string, o vesaliusGeo.AppointmentChangeConfirmation, tx *sqlx.Tx) error {
    query := `
        UPDATE NOVA_DOCTOR_PATIENT_APPT SET
          DOCTOR_NAME = :doctor_name,
          DOCTOR_SPECIALTY = :doctor_specialty,
          APPT_STATUS = :appointment_status,
          APPT_NO = :appointment_number,
          APPT_DAY = :appointment_day,
          APPT_SESSIONTYPE = :appointment_sessionType,
          APPT_CLINIC = :appointment_clinic,
          APPT_ROOM = :appointment_room,
          APPT_CASETYPE = :appointment_caseType,
          DATE_APPT = TO_DATE(:appointment_dateTime, 'DD-MON-YYYY hh24:mi')
        WHERE PATIENT_PRN = :patient_prn AND APPT_NO = :appointment_number
    `
    args := []any{
        sql.Named("doctor_name", o.DoctorName),
        sql.Named("doctor_specialty", o.Specialty),
        sql.Named("appointment_status", o.AppointmentStatus),
        sql.Named("appointment_number", o.AppointmentNumber),
        sql.Named("appointment_day", o.Day),
        sql.Named("appointment_sessionType", apptSessionType),
        sql.Named("appointment_clinic", o.Clinic),
        sql.Named("appointment_room", o.Room),
        sql.Named("appointment_caseType", o.CaseType),
        sql.Named("appointment_dateTime", fmt.Sprintf("%s %s", o.Date, o.StartTime)),
        sql.Named("patient_prn", prn),
        sql.Named("appointment_number", o.AppointmentNumber),
    }
    var err error
    if tx == nil {
        _, err = s.db.ExecContext(s.ctx, query, args...)
    } else {
        _, err = tx.ExecContext(s.ctx, query, args...)
    }
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *NovaDoctorPatientApptService) UpdateToCancel(prn string, o vesaliusGeo.AppointmentCancelConfirmation, tx *sqlx.Tx) error {
    query := `
        UPDATE NOVA_DOCTOR_PATIENT_APPT SET
          APPT_STATUS = :appointment_status,
          APPT_NO = :appointment_number
        WHERE PATIENT_PRN = :patient_prn AND APPT_NO = :appointment_number
    `
    args := []any{
        sql.Named("appointment_status", o.AppointmentStatus),
        sql.Named("appointment_number", o.AppointmentNumber),
        sql.Named("patient_prn", prn),
    }
    var err error
    if tx == nil {
        _, err = s.db.ExecContext(s.ctx, query, args...)
    } else {
        _, err = tx.ExecContext(s.ctx, query, args...)
    }
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *NovaDoctorPatientApptService) getPatientName(conn *sqlx.DB, prn string) (string, error) {
    query := `
        SELECT COALESCE(
          (SELECT FIRST_NAME FROM APPLICATION_USER WHERE MASTER_PRN = :prn and INACTIVE_FLAG = 'N' FETCH FIRST 1 ROW ONLY),
          (SELECT FULLNAME FROM APPLICATION_USER_FAMILY WHERE NOK_PRN = :prn FETCH FIRST 1 ROW ONLY),
          (SELECT GUEST_NAME FROM APPLICATION_GUEST WHERE GUEST_PRN = :prn FETCH FIRST 1 ROW ONLY)
        ) AS PATIENT_NAME FROM DUAL
    `
    var patientName string
    err := conn.GetContext(s.ctx, &patientName, query, sql.Named("prn", prn))
    if err != nil {
        utils.LogError(err)
        return "", err
    }
    return patientName, nil
}

func buildKeywordConditions(keyword string, keyword2 string, keyword3 string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `(LOWER(ndpa.PATIENT_NAME) LIKE :keyword OR ndpa.PATIENT_PRN LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        conds = append(conds, `LOWER(ndpa.DOCTOR_NAME) LIKE :keyword2`)
        args = append(args, sql.Named("keyword2", strings.ToLower(keyword2)))
    }
    if keyword3 != "" {
        conds = append(conds, `TRUNC(ndpa.DATE_APPT) = TO_DATE(:keyword3, 'dd/mm/yyyy')`)
        args = append(args, sql.Named("keyword3", keyword3))
    }

    return conds, args
}

func whereClause(conds []string) string {
    if len(conds) == 0 {
        return ""
    }
    return " WHERE " + strings.Join(conds, " AND ")
}
