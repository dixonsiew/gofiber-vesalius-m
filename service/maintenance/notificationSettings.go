package maintenance

import (
	"database/sql"
	"strings"
	"vesaliusm/dto"
	"vesaliusm/model"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

func (s *MaintenanceService) ListAllNotificationSettings(page string, limit string) (*model.PagedList, error) {
    total, err := s.CountAllNotificationSettings(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllNotificationSettings(pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *MaintenanceService) SearchAllNotificationSettingsByKeyword(keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountNotificationSettingsByKeyword(keyword, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindNotificationSettingsByKeyword(keyword, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *MaintenanceService) CountAllNotificationSettings(conn *sqlx.DB) (int, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT COUNT(*) AS COUNT FROM NOTIFICATION_SETTINGS`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *MaintenanceService) CountNotificationSettingsByKeyword(keyword string, conn *sqlx.DB) (int, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    conds, args := buildNotificationSettingsConditions(keyword)
    base := `SELECT COUNT(*) AS COUNT FROM NOTIFICATION_SETTINGS ns`
    query := base + whereClause(conds)

    var count int
    err := db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *MaintenanceService) FindAllNotificationSettings(offset int, limit int, conn *sqlx.DB) ([]model.NotificationSetting, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM NOTIFICATION_SETTINGS ORDER BY NOTIFICATION_CODE OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NotificationSetting{}, ""), 1)
    list := make([]model.NotificationSetting, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) FindNotificationSettingsByKeyword(keyword string, offset int, limit int, conn *sqlx.DB) ([]model.NotificationSetting, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    conditions, args := buildNotificationSettingsConditions(keyword)
    args = append(args, sql.Named("offset", offset))
	args = append(args, sql.Named("limit", limit))

    base := `SELECT * FROM NOTIFICATION_SETTINGS ns`
    query := base + whereClause(conditions) + ` OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NotificationSetting{}, ""), 1)
    list := make([]model.NotificationSetting, 0)
    err := db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) UpdateNotificationSettingByNotificationCode(data *dto.NotificationSettingDto) error {
    query := `UPDATE NOTIFICATION_SETTINGS SET NOTIFICATION_PARAM_1 = :p1, NOTIFICATION_PARAM_2 = :p2 WHERE NOTIFICATION_CODE = :p3`
    _, err := s.db.ExecContext(s.ctx, query, data.NotificationParam1, data.NotificationParam2, data.NotificationCode)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return err
}

func buildNotificationSettingsConditions(keyword string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `LOWER(ns.NOTIFICATION_DESC) LIKE :keyword`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }

    return conds, args
}
