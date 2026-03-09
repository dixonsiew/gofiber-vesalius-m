package generalNotificationMaster

import (
    "context"
    "database/sql"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/guregu/null/v6"
    "github.com/jmoiron/sqlx"
)

type GeneralNotificationMasterService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewGeneralNotificationMasterService(db *sqlx.DB, ctx context.Context) *GeneralNotificationMasterService {
    return &GeneralNotificationMasterService{db: db, ctx: ctx}
}

// save inserts a new general notification.
func (s *GeneralNotificationMasterService) Save(o *model.GeneralNotification, adminId int64) error {
    // Convert string dates (DD/MM/YYYY) to values for binding.
    // If the string is "-", treat as NULL.
    startDateVal := nullStringIfDash(o.StartDate)
    endDateVal := nullStringIfDash(o.EndDate)

    query := `
        INSERT INTO GENERAL_NOTIFICATION_MASTER
        (
            NOTIFICATION_TITLE, SHORT_MESSAGE, FULL_MESSAGE, START_DATE_TIME,
            END_DATE_TIME, TARGET_AGE_FROM, TARGET_AGE_TO, TARGET_GENDER,
            TARGET_NATIONALITY, TARGET_CITY, TARGET_STATE, USER_CREATE
        )
        VALUES
        (
            :notificationTitle, :shortMessage, :fullMessage, 
            TO_DATE(:startDateTime, 'DD/MM/YYYY hh24:mi'),
            TO_DATE(:endDateTime, 'DD/MM/YYYY hh24:mi'),
            :targetAgeFrom, :targetAgeTo, :targetGender,
            :targetNationality, :targetCity, :targetState, :adminId
        )
    `
    params := map[string]interface{}{
        "notificationTitle": o.NotificationTitle.String,
        "shortMessage":      o.ShortMessage.String,
        "fullMessage":       o.FullMessage.String,
        "startDateTime":     startDateVal,
        "endDateTime":       endDateVal,
        "targetAgeFrom":     o.TargetAgeFrom.Int64,
        "targetAgeTo":       o.TargetAgeTo.Int64,
        "targetGender":      o.TargetGender.String,
        "targetNationality": o.TargetNationality.String,
        "targetCity":        o.TargetCity.String,
        "targetState":       o.TargetState.String,
        "adminId":           adminId,
    }
    _, err := s.db.NamedExecContext(s.ctx, query, params)
    return err
}

func (s *GeneralNotificationMasterService) Update(o *model.GeneralNotification, adminId int64) error {
    startDateVal := nullStringIfDash(o.StartDate)
    endDateVal := nullStringIfDash(o.EndDate)

    query := `
        UPDATE GENERAL_NOTIFICATION_MASTER SET
            NOTIFICATION_TITLE = :notificationTitle,
            SHORT_MESSAGE = :shortMessage,
            FULL_MESSAGE = :fullMessage,
            START_DATE_TIME = TO_DATE(:startDateTime, 'DD/MM/YYYY hh24:mi'),
            END_DATE_TIME = TO_DATE(:endDateTime, 'DD/MM/YYYY hh24:mi'),
            TARGET_AGE_FROM = :targetAgeFrom,
            TARGET_AGE_TO = :targetAgeTo,
            TARGET_GENDER = :targetGender,
            TARGET_NATIONALITY = :targetNationality,
            TARGET_CITY = :targetCity,
            TARGET_STATE = :targetState,
            USER_UPDATE = :adminId,
            DATE_UPDATE = CURRENT_TIMESTAMP
        WHERE NOTIFICATION_MASTER_ID = :notification_master_id
    `
    params := map[string]interface{}{
        "notification_master_id": o.NotificationMasterID.Int64,
        "notificationTitle":      o.NotificationTitle.String,
        "shortMessage":           o.ShortMessage.String,
        "fullMessage":            o.FullMessage.String,
        "startDateTime":          startDateVal,
        "endDateTime":            endDateVal,
        "targetAgeFrom":          o.TargetAgeFrom.Int64,
        "targetAgeTo":            o.TargetAgeTo.Int64,
        "targetGender":           o.TargetGender.String,
        "targetNationality":      o.TargetNationality.String,
        "targetCity":             o.TargetCity.String,
        "targetState":            o.TargetState.String,
        "adminId":                adminId,
    }
    _, err := s.db.NamedExecContext(s.ctx, query, params)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *GeneralNotificationMasterService) FindByNotificationMasterId(notificationMasterId int64) (*model.GeneralNotification, error) {
    query := `SELECT ` + getGeneralNotificationMasterCols() + ` FROM GENERAL_NOTIFICATION_MASTER WHERE NOTIFICATION_MASTER_ID = :1`
    var o model.GeneralNotification
    err := s.db.GetContext(s.ctx, &o, query, notificationMasterId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    return &o, err
}

func (s *GeneralNotificationMasterService) List(page string, limit string) (*model.PagedList, error) {
    total, err := s.Count(s.db)
    if err != nil {
        utils.LogError(err)
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

func (s *GeneralNotificationMasterService) Count(conn *sqlx.DB) (int, error) {
    db := s.db
    if conn != nil {
        db = conn
    }
    var count int
    err := db.GetContext(s.ctx, &count, `SELECT COUNT(NOTIFICATION_MASTER_ID) FROM GENERAL_NOTIFICATION_MASTER`)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *GeneralNotificationMasterService) FindAll(offset int, limit int, conn *sqlx.DB) ([]model.GeneralNotification, error) {
    db := s.db
    if conn != nil {
        db = conn
    }
    query := `
        SELECT ` + getGeneralNotificationMasterCols() + ` FROM GENERAL_NOTIFICATION_MASTER
        ORDER BY DATE_CREATE DESC OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY
    `
    var lx []model.GeneralNotification
    err := db.SelectContext(s.ctx, &lx, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for _, o := range lx {
        o.Set()
    }
    return lx, nil
}

func nullStringIfDash(s null.String) interface{} {
    if !s.Valid {
        return nil
    }
    if s.String == "-" {
        return nil
    }
    return s.String
}

func getGeneralNotificationMasterCols() string {
    return `
        NOTIFICATION_MASTER_ID,
        NOTIFICATION_TITLE,
        SHORT_MESSAGE,
        FULL_MESSAGE,
        START_DATE_TIME,
        END_DATE_TIME,
        TARGET_AGE_FROM,
        TARGET_AGE_TO,
        TARGET_GENDER,
        TARGET_NATIONALITY,
        TARGET_CITY,
        TARGET_STATE
    `
}
