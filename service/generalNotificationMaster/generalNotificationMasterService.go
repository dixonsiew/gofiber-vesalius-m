package generalNotificationMaster

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/guregu/null/v6"
    "github.com/jmoiron/sqlx"
)

var GeneralNotificationMasterSvc *GeneralNotificationMasterService = NewGeneralNotificationMasterService(database.GetDb(), database.GetCtx())

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
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("notificationTitle", o.NotificationTitle.String),
        sql.Named("shortMessage", o.ShortMessage.String),
        sql.Named("fullMessage", o.FullMessage.String),
        sql.Named("startDateTime", startDateVal),
        sql.Named("endDateTime", endDateVal),
        sql.Named("targetAgeFrom", o.TargetAgeFrom.Int64),
        sql.Named("targetAgeTo", o.TargetAgeTo.Int64),
        sql.Named("targetGender", o.TargetGender.String),
        sql.Named("targetNationality", o.TargetNationality.String),
        sql.Named("targetCity", o.TargetCity.String),
        sql.Named("targetState", o.TargetState.String),
        sql.Named("adminId", adminId),
    )
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
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("notificationTitle", o.NotificationTitle.String),
        sql.Named("shortMessage", o.ShortMessage.String),
        sql.Named("fullMessage", o.FullMessage.String),
        sql.Named("startDateTime", startDateVal),
        sql.Named("endDateTime", endDateVal),
        sql.Named("targetAgeFrom", o.TargetAgeFrom.Int64),
        sql.Named("targetAgeTo", o.TargetAgeTo.Int64),
        sql.Named("targetGender", o.TargetGender.String),
        sql.Named("targetNationality", o.TargetNationality.String),
        sql.Named("targetCity", o.TargetCity.String),
        sql.Named("targetState", o.TargetState.String),
        sql.Named("adminId", adminId),
        sql.Named("notification_master_id", o.NotificationMasterId.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *GeneralNotificationMasterService) FindByNotificationMasterId(notificationMasterId int64) (*model.GeneralNotification, error) {
    query := `SELECT * FROM GENERAL_NOTIFICATION_MASTER WHERE NOTIFICATION_MASTER_ID = :notification_master_id`
    query = strings.Replace(query, "*", utils.GetDbCols(model.GeneralNotification{}, ""), 1)
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
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAll(pager.GetLowerBound(), pager.PageSize, s.db)
    if err != nil {
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
        SELECT * FROM GENERAL_NOTIFICATION_MASTER
        ORDER BY DATE_CREATE DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.GeneralNotification{}, ""), 1)
    list := make([]model.GeneralNotification, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
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

/* func getGeneralNotificationMasterCols() string {
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
} */
