package maintenance

import (
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/dto"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

func (s *MaintenanceService) ListAllDynamicEmailSettings(page string, limit string) (*model.PagedList, error) {
    total, err := s.CountAllDynamicEmailSettings(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllDynamicEmailSettings(pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *MaintenanceService) CountAllDynamicEmailSettings(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(*) AS COUNT FROM EMAIL_MASTER`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *MaintenanceService) FindAllDynamicEmailSettings(offset int, limit int, conn *sqlx.DB) ([]model.DynamicEmailMaster, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM EMAIL_MASTER OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.DynamicEmailMaster{}, ""), 1)
    list := make([]model.DynamicEmailMaster, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) SearchAllDynamicEmailSettingsByKeyword(keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountDynamicEmailSettingsByKeyword(keyword, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindDynamicEmailSettingsByKeyword(keyword, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *MaintenanceService) CountDynamicEmailSettingsByKeyword(keyword string, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    conds, args := buildDynamicEmailSettingsConditions(keyword)
    base := `SELECT COUNT(*) AS COUNT FROM EMAIL_MASTER em`
    query := base + whereClause(conds)

    var count int
    err := db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *MaintenanceService) FindDynamicEmailSettingsByKeyword(keyword string, offset int, limit int, conn *sqlx.DB) ([]model.DynamicEmailMaster, error) {
    db := database.GetFromCon(conn, s.db)
    conditions, args := buildDynamicEmailSettingsConditions(keyword)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    base := `SELECT * FROM EMAIL_MASTER em`
    query := base + whereClause(conditions) + ` OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.DynamicEmailMaster{}, ""), 1)
    list := make([]model.DynamicEmailMaster, 0)
    err := db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) ViewDynamicEmailSettingByFunctionName(functionName string) (*model.DynamicEmailMaster, error) {
    query := `SELECT * FROM EMAIL_MASTER WHERE EMAIL_FUNCTION_NAME = :fname`
    query = strings.Replace(query, "*", utils.GetDbCols(model.DynamicEmailMaster{}, ""), 1)
    var o model.DynamicEmailMaster
    err := s.db.GetContext(s.ctx, &o, query, functionName)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, err
}

func (s *MaintenanceService) UpdateDynamicEmailSettingByFunctionName(data *dto.DynamicEmailMasterDto) error {
    query := `
        UPDATE EMAIL_MASTER SET 
          EMAIL_MODULE = :emmod,
          EMAIL_FOR = :emfor,
          EMAIL_SUBJECT = :emsub,
          EMAIL_RECIPIENT = :emrec,
          EMAIL_SENDER = :emsender,
          EMAIL_SENDER_NAME = :emsendname,
          EMAIL_TEMPLATE = :emtpl 
        WHERE EMAIL_FUNCTION_NAME = :emfname
    `
    _, err := s.db.ExecContext(s.ctx, query, 
        data.EmailModule, 
        data.EmailFor, 
        data.EmailSubject,
        data.EmailRecipient,
        data.EmailSender,
        data.EmailSenderName,
        data.EmailTemplate,
        data.EmailFunctionName,
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func buildDynamicEmailSettingsConditions(keyword string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `(LOWER(em.EMAIL_SUBJECT) LIKE :keyword OR LOWER(em.EMAIL_RECIPIENT) LIKE :keyword OR LOWER(em.EMAIL_SENDER) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }

    return conds, args
}
