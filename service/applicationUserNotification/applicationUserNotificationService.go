package applicationUserNotification

import (
    "context"
    "database/sql"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

type ApplicationUserNotificationService struct {
    db  *sqlx.DB
    ctx context.Context
}

// NewApplicationUserNotificationService creates a new instance.
func NewApplicationUserNotificationService(db *sqlx.DB, ctx context.Context) *ApplicationUserNotificationService {
    return &ApplicationUserNotificationService{db: db, ctx: ctx}
}

func (s *ApplicationUserNotificationService) Save(o *model.OnesignalNotification) error {
    query := `
        INSERT INTO APPLICATION_USER_NOTIFICATION 
        (USER_ID, NOTIFICATION_TITLE, MSG_TYPE, SHORT_MESSAGE, FULL_MESSAGE, VISIT_TYPE, ACCOUNT_NO)
        VALUES 
        (:user_id, :notification_title, :msg_type, :short_message, :full_message, :visit_type, :account_no)
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("user_id", o.UserID.Int64),
        sql.Named("notification_title", o.NotificationTitle.String),
        sql.Named("msg_type", o.MsgType.String),
        sql.Named("short_message", o.ShortMessage.String),
        sql.Named("full_message", o.FullMessage.String),
        sql.Named("visit_type", o.VisitType.String),
        sql.Named("account_no", o.AccountNo.String),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserNotificationService) CountUnseenByUserId(userId int64) (int, error) {
    query := `SELECT COUNT(USER_ID) FROM APPLICATION_USER_NOTIFICATION WHERE USER_ID = :userId AND DATE_SENT IS NOT NULL AND IS_SEEN = 'N'`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, userId)
    if err != nil {
        utils.LogError(err)
    }
    return count, nil
}

func (s *ApplicationUserNotificationService) ListByUserId(userId int64, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountByUserId(userId, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllByUserId(userId, pager.GetLowerBound(), pager.PageSize, nil)
    if err != nil {
        return nil, err
    }
    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ApplicationUserNotificationService) CountByUserId(userId int64, conn *sqlx.DB) (int, error) {
    db := s.db
    if conn != nil {
        db = conn
    }
    query := `SELECT COUNT(USER_ID) FROM APPLICATION_USER_NOTIFICATION WHERE USER_ID = :userId AND DATE_SENT IS NOT NULL`
    var count int
    err := db.GetContext(s.ctx, &count, query, userId)
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (s *ApplicationUserNotificationService) FindAllByUserId(userId int64, offset int, limit int, conn *sqlx.DB) ([]model.OnesignalNotification, error) {
    db := s.db
    if conn != nil {
        db = conn
    }
    query := `
        SELECT ` + getOnesignalNotificationCols() + ` FROM APPLICATION_USER_NOTIFICATION 
        WHERE USER_ID = :userId AND DATE_SENT IS NOT NULL
        ORDER BY DATE_SENT DESC 
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    notifications := make([]model.OnesignalNotification, 0)
    err := db.SelectContext(s.ctx, &notifications, query, userId, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for _, o := range notifications {
        o.Set()
    }
    return notifications, nil
}

func (s *ApplicationUserNotificationService) FindByNotificationId(notificationId int64) (*model.OnesignalNotification, error) {
    query := `SELECT ` + getOnesignalNotificationCols() + ` FROM APPLICATION_USER_NOTIFICATION WHERE NOTIFICATION_ID = :notificationId`
    var o model.OnesignalNotification
    err := s.db.GetContext(s.ctx, &o, query, notificationId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    return &o, nil
}

func (s *ApplicationUserNotificationService) UpdateSeenByUserId(userId int64, notificationId int64) error {
    query := `
        UPDATE APPLICATION_USER_NOTIFICATION SET 
            IS_SEEN = 'Y',
            DATE_SEEN = CURRENT_TIMESTAMP
        WHERE NOTIFICATION_ID = :notificationId AND USER_ID = :userId
    `
    _, err := s.db.ExecContext(s.ctx, query, notificationId, userId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func getOnesignalNotificationCols() string {
    return `
        NOTIFICATION_ID,
        USER_ID,
        VISIT_TYPE,
        ACCOUNT_NO,
        NOTIFICATION_TITLE,
        MSG_TYPE,
        SHORT_MESSAGE,
        FULL_MESSAGE,
        MASTER_ID,
        IS_SEEN,
        DATE_SENT
    `
}
