package guest

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    "golang.org/x/crypto/bcrypt"
)

var GuestSvc *GuestService = NewGuestService(database.GetDb(), database.GetCtx())

type GuestService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewGuestService(db *sqlx.DB, ctx context.Context) *GuestService {
    return &GuestService{
        db:  db,
        ctx: ctx,
    }
}

func (s *GuestService) FindGuestAccount() (*model.ApplicationUser, error) {
    query := `SELECT * FROM APPLICATION_USER WHERE USERNAME = 'nova_default'`
    query = strings.Replace(query, "*", utils.GetDbCols(model.ApplicationUser{}, ""), 1)
    var o model.ApplicationUser
    err := s.db.GetContext(s.ctx, &o, query)
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

func (s *GuestService) GetTempGuestInfo(prn string) (*model.TempGuest, error) {
    query := `SELECT * FROM APPLICATION_GUEST WHERE GUEST_PRN = :prn`
    query = strings.Replace(query, "*", utils.GetDbCols(model.TempGuest{}, ""), 1)
    var o model.TempGuest
    err := s.db.GetContext(s.ctx, &o, query, prn)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, err
}

func (s *GuestService) DeleteTempGuestInfo(prn string) error {
    qeury := `DELETE FROM APPLICATION_GUEST WHERE GUEST_PRN = :prn`
    _, err := s.db.ExecContext(s.ctx, qeury, prn)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *GuestService) UpdateSeenNotificationByPlayerId(playerId string, notificationId int64) error {
    query := `
        UPDATE APPLICATION_USER_NOTIFICATION SET 
          IS_SEEN = 'Y',
          DATE_SEEN = CURRENT_TIMESTAMP
        WHERE NOTIFICATION_ID = :notificationId AND GUEST_PLAYER_ID = :playerId
    `
    _, err := s.db.ExecContext(s.ctx, query, notificationId, playerId)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *GuestService) CountUnseenNotificationByGuestPlayerId(playerId string) (int, error) {
    var count int
    query := `
        SELECT COUNT(GUEST_PLAYER_ID) AS COUNT FROM APPLICATION_USER_NOTIFICATION 
        WHERE GUEST_PLAYER_ID = :playerId AND MSG_TYPE = 'GENERAL_INFO' AND DATE_SENT IS NOT NULL AND IS_SEEN = 'N'
    `
    err := s.db.GetContext(s.ctx, &count, query, playerId)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *GuestService) ListGuestNotification(playerId string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountGuestNotification(playerId, s.db)
	if err != nil {
		return nil, err
	}
	pager := model.GetPager(total, page, limit)
	list, err := s.FindAllGuestNotification(playerId, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *GuestService) CountGuestNotification(playerId string, conn *sqlx.DB) (int, error) {
    var count int
    query := `
        SELECT COUNT(UNIQUE MASTER_ID) AS COUNT FROM APPLICATION_USER_NOTIFICATION 
        WHERE MSG_TYPE = 'GENERAL_INFO'
        AND DATE_SENT IS NOT NULL
        AND GUEST_PLAYER_ID = :playerId
    `
    err := conn.GetContext(s.ctx, &count, query, playerId)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *GuestService) FindAllGuestNotification(playerId string, offset int, limit int, conn *sqlx.DB) ([]model.OnesignalNotification, error) {
    query := `
        SELECT * FROM APPLICATION_USER_NOTIFICATION
        WHERE MSG_TYPE = 'GENERAL_INFO'
        AND GUEST_PLAYER_ID = :playerId
        AND NOTIFICATION_ID IN (
          SELECT NOTIFICATION_ID FROM (
            SELECT MASTER_ID, MAX(NOTIFICATION_ID) AS NOTIFICATION_ID
            FROM APPLICATION_USER_NOTIFICATION
            WHERE MSG_TYPE = 'GENERAL_INFO'
            AND GUEST_PLAYER_ID = :playerId
            GROUP BY MASTER_ID
          )   
         )  
        AND DATE_SENT IS NOT NULL
        ORDER BY DATE_CREATE DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.OnesignalNotification{}, ""), 1)
    list := make([]model.OnesignalNotification, 0)
    err := s.db.SelectContext(s.ctx, &list, query,
        sql.Named("playerId", playerId),
        sql.Named("offset", offset),
        sql.Named("limit", limit),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *GuestService) ValidateGuestCredentials(user model.ApplicationUser, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
    return err == nil
}
