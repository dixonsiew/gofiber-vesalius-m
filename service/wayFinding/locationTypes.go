package wayfinding

import (
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

func (s *WayFindingService) SaveLocationType(o *model.WayFindingLocationTypes) error {
    query := `
        INSERT INTO WAY_FINDING_LOCATION_TYPE
        (LOCATION_TYPE_CODE, LOCATION_TYPE_NAME, LOCATION_TYPE_DISPLAY_ORDER)
        VALUES
        (:locationTypeCode, :locationTypeName, :locationTypeDisplayOrder)
    `
    _, err := s.db.ExecContext(s.ctx, query, 
        sql.Named("locationTypeCode", o.LocationTypeCode.String),
        sql.Named("locationTypeName", o.LocationTypeName.String),
        sql.Named("locationTypeDisplayOrder", o.LocationTypeId.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) UpdateLocationType(o *model.WayFindingLocationTypes, locationTypeCode string) error {
    query := `
        UPDATE WAY_FINDING_LOCATION_TYPE SET
          LOCATION_TYPE_NAME = :locationTypeName,
          LOCATION_TYPE_DISPLAY_ORDER = :locationTypeDisplayOrder
        WHERE LOCATION_TYPE_CODE = :locationTypeCode
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("locationTypeCode", o.LocationTypeCode.String),
        sql.Named("locationTypeName", o.LocationTypeName.String),
        sql.Named("locationTypeDisplayOrder", o.LocationTypeDisplayOrder.Int32),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) ListLocationTypes(page string, limit string) (*model.PagedList, error) {
    total, err := s.LocationTypesCount(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllWayFindingLocationTypes(pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) LocationTypesCount(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(*) AS COUNT FROM WAY_FINDING_LOCATION_TYPE`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *WayFindingService) FindAllWayFindingLocationTypes(offset int, limit int, conn *sqlx.DB) ([]model.WayFindingLocationTypes, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM WAY_FINDING_LOCATION_TYPE ORDER BY LOCATION_TYPE_DISPLAY_ORDER OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingLocationTypes{}, ""), 1)
    list := make([]model.WayFindingLocationTypes, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *WayFindingService) ListLocationTypesByKeyword(keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.LocationTypesCountByKeyword(keyword, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllWayFindingLocationTypesByKeyword(keyword, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) LocationTypesCountByKeyword(keyword string, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT COUNT(*) AS COUNT 
        FROM WAY_FINDING_LOCATION_TYPE
        WHERE LOWER(LOCATION_TYPE_CODE) LIKE :keyword
        OR LOWER(LOCATION_TYPE_NAME) LIKE :keyword
    `
    var count int
    err := db.GetContext(s.ctx, &count, query, sql.Named("keyword", strings.ToLower(keyword)))
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *WayFindingService) FindAllWayFindingLocationTypesByKeyword(keyword string, offset int, limit int, conn *sqlx.DB) ([]model.WayFindingLocationTypes, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT * FROM WAY_FINDING_LOCATION_TYPE
        WHERE LOWER(LOCATION_TYPE_CODE) LIKE :keyword
        OR LOWER(LOCATION_TYPE_NAME) LIKE :keyword
        ORDER BY LOCATION_TYPE_DISPLAY_ORDER
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingLocationTypes{}, ""), 1)
    list := make([]model.WayFindingLocationTypes, 0)
    err := db.SelectContext(s.ctx, &list, query, 
        sql.Named("keyword", strings.ToLower(keyword)), 
        sql.Named("offset", offset), 
        sql.Named("limit", limit),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *WayFindingService) DeleteLocationTypesByLocationTypeCode(locationTypeCode string) error {
    query := `DELETE FROM WAY_FINDING_LOCATION_TYPE WHERE LOCATION_TYPE_CODE = :locationTypeCode`
    _, err := s.db.ExecContext(s.ctx, query, locationTypeCode)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) FindLocationTypesByLocationTypeCode(locationTypeCode string) (*model.WayFindingLocationTypes, error) {
    query := `SELECT * FROM WAY_FINDING_LOCATION_TYPE WHERE LOCATION_TYPE_CODE = :locationTypeCode`
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingLocationTypes{}, ""), 1)
    var o model.WayFindingLocationTypes
    err := s.db.GetContext(s.ctx, &o, query, locationTypeCode)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}
