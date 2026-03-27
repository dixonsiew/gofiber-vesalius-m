package maintenance

import (
	"database/sql"
	"strings"
	"vesaliusm/dto"
	"vesaliusm/model"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

func (s *MaintenanceService) ListAllParamSettings(page string, limit string) (*model.PagedList, error) {
    total, err := s.CountAllParamSettings(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllParamSettings(pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *MaintenanceService) SearchAllParamSettingsByKeyword(keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountParamSettingsByKeyword(keyword, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindParamSettingsByKeyword(keyword, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *MaintenanceService) CountAllParamSettings(conn *sqlx.DB) (int, error){
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT COUNT(*) AS COUNT FROM PARAM_SETTINGS`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *MaintenanceService) CountParamSettingsByKeyword(keyword string, conn *sqlx.DB) (int, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    conds, args := buildParamSettingsConditions(keyword)
    base := `SELECT COUNT(*) AS COUNT FROM PARAM_SETTINGS ps`
    query := base + whereClause(conds)

    var count int
	err := s.db.GetContext(s.ctx, &count, query, args...)
	if err != nil {
		utils.LogError(err)
		return 0, err
	}
	return count, nil
}

func (s *MaintenanceService) FindAllParamSettings(offset int, limit int, conn *sqlx.DB) ([]model.ParamSetting, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM PARAM_SETTINGS ORDER BY PARAM_CODE OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.ParamSetting{}, ""), 1)
    list := make([]model.ParamSetting, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) FindParamSettingsByKeyword(keyword string, offset int, limit int, conn *sqlx.DB) ([]model.ParamSetting, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    conditions, args := buildParamSettingsConditions(keyword)
    args = append(args, sql.Named("offset", offset))
	args = append(args, sql.Named("limit", limit))

    base := `SELECT ps.* FROM PARAM_SETTINGS ps`
    query := base + whereClause(conditions) + ` OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "ps.*", utils.GetDbCols(model.ParamSetting{}, "ps."), 1)
    list := make([]model.ParamSetting, 0)
    err := db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) UpdateParamSettingByParamCode(data *dto.ParamSettingDto) error {
    query := `UPDATE PARAM_SETTINGS SET PARAM_VALUE = :pVal WHERE PARAM_CODE = :pCode`
    _, err := s.db.ExecContext(s.ctx, query, 
        sql.Named("pVal", data.ParamValue),
        sql.Named("pCode", data.ParamCode),
    )
    if err != nil {
		utils.LogError(err)
		return err
    }
    return err
}

func buildParamSettingsConditions(keyword string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `(LOWER(ps.PARAM_CODE) LIKE :keyword OR LOWER(ps.PARAM_DESC) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }

    return conds, args
}
