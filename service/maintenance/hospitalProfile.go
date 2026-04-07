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

func (s *MaintenanceService) ListAllHospitalProfiles(page string, limit string) (*model.PagedList, error) {
    total, err := s.CountAllHospitalProfiles(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllHospitalProfiles(pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *MaintenanceService) SearchAllHospitalProfilesByKeyword(keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountHospitalProfilesByKeyword(keyword, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindHospitalProfilesByKeyword(keyword, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *MaintenanceService) CountAllHospitalProfiles(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(*) AS COUNT FROM HOSPITAL_PROFILE`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *MaintenanceService) CountHospitalProfilesByKeyword(keyword string, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    conds, args := buildKeywordConditions(keyword)
    base := `SELECT COUNT(*) AS COUNT FROM HOSPITAL_PROFILE hp`
    query := base + whereClause(conds)

    var count int
    err := db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *MaintenanceService) FindAllHospitalProfiles(offset int, limit int, conn *sqlx.DB) ([]model.HospitalProfile, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM HOSPITAL_PROFILE OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.HospitalProfile{}, ""), 1)
    list := make([]model.HospitalProfile, 0)
    err := db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) FindHospitalProfilesByKeyword(keyword string, offset int, limit int, conn *sqlx.DB) ([]model.HospitalProfile, error) {
    db := database.GetFromCon(conn, s.db)
    conditions, args := buildKeywordConditions(keyword)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    base := `SELECT hp.* FROM HOSPITAL_PROFILE hp`
    base = strings.Replace(base, "hp.*", utils.GetDbCols(model.HospitalProfile{}, "hp."), 1)
    query := base + whereClause(conditions) + ` OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    list := make([]model.HospitalProfile, 0)
    err := db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) UpdateHospitalProfileByDescName(data *dto.HospitalProfileDto) error {
    query := `UPDATE HOSPITAL_PROFILE SET PROFILE_VALUE = :profileVal WHERE PROFILE_DESC = :profileDesc`
    _, err := s.db.ExecContext(s.ctx, query, 
        sql.Named("profileVal", data.ProfileValue),
        sql.Named("profileDesc", data.ProfileDesc),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return err
}

func buildKeywordConditions(keyword string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `(LOWER(hp.PROFILE_DESC) LIKE :keyword OR LOWER(hp.PROFILE_VALUE) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }

    return conds, args
}
