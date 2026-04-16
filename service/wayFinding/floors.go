package wayFinding

import (
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

func (s *WayFindingService) SaveFloor(o *model.WayFindingFloors) error {
    query := `
        INSERT INTO WAY_FINDING_FLOORS
         (FLOOR_CODE, FLOOR_NAME, FLOOR_DISPLAY_ORDER, FLOOR_IMAGE_RAW)
         VALUES
        (:floorCode, :floorName, :floorDisplayOrder, :floorImageRaw)
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("floorCode", o.FloorCode.String),
        sql.Named("floorName", o.FloorName.String),
        sql.Named("floorDisplayOrder", o.FloorId.Int64),
        sql.Named("floorImageRaw", o.FloorImageRaw.String),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) UpdateFloor(o *model.WayFindingFloors) error {
    query := `
        UPDATE WAY_FINDING_FLOORS SET
          FLOOR_CODE = :floorCode,
          FLOOR_NAME = :floorName,
          FLOOR_DISPLAY_ORDER = :floorDisplayOrder,
          FLOOR_IMAGE_RAW = :floorImageRaw
        WHERE FLOOR_CODE = :floorCode
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("floorCode", o.FloorCode.String),
        sql.Named("floorName", o.FloorName.String),
        sql.Named("floorDisplayOrder", o.FloorId.Int64),
        sql.Named("floorImageRaw", o.FloorImageRaw.String),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) ListFloors() ([]model.WayFindingFloors, error) {
    return s.FindAllFloors(s.db)
}

func (s *WayFindingService) ListFloorsWebAdmin(page string, limit string) (*model.PagedList, error) {
    total, err := s.FloorsCount(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllFloorsWebAdmin(pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) FloorsCount(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(*) AS COUNT FROM WAY_FINDING_FLOORS`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *WayFindingService) FindAllFloors(conn *sqlx.DB) ([]model.WayFindingFloors, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM WAY_FINDING_FLOORS ORDER BY FLOOR_DISPLAY_ORDER`
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingFloors{}, ""), 1)
    list := make([]model.WayFindingFloors, 0)
    err := db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *WayFindingService) FindAllFloorsWebAdmin(offset int, limit int, conn *sqlx.DB) ([]model.WayFindingFloors, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM WAY_FINDING_FLOORS ORDER BY FLOOR_DISPLAY_ORDER OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingFloors{}, ""), 1)
    list := make([]model.WayFindingFloors, 0)
    err := db.SelectContext(s.ctx, &list, query, limit, offset)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *WayFindingService) ListFloorsByKeyword(keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.FloorsCountByKeyword(keyword, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllFloorsByKeyword(keyword, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) ListFloorsByKeywordWebAdmin(keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.FloorsCountByKeyword(keyword, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllFloorsByKeywordWebAdmin(keyword, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) FloorsCountByKeyword(keyword string, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT COUNT(*) AS COUNT 
        FROM WAY_FINDING_FLOORS
        WHERE LOWER(FLOOR_CODE) LIKE :keyword
        OR LOWER(FLOOR_NAME) LIKE :keyword
    `
    var count int
    err := db.GetContext(s.ctx, &count, query, sql.Named("keyword", strings.ToLower(keyword)))
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *WayFindingService) FindAllFloorsByKeyword(keyword string, offset int, limit int, conn *sqlx.DB) ([]model.WayFindingFloors, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT * FROM WAY_FINDING_FLOORS
        WHERE LOWER(FLOOR_CODE) LIKE :keyword
        OR LOWER(FLOOR_NAME) LIKE :keyword
        ORDER BY FLOOR_DISPLAY_ORDER
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingFloors{}, ""), 1)
    list := make([]model.WayFindingFloors, 0)
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

func (s *WayFindingService) FindAllFloorsByKeywordWebAdmin(keyword string, offset int, limit int, conn *sqlx.DB) ([]model.WayFindingFloors, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT * FROM WAY_FINDING_FLOORS
        WHERE LOWER(FLOOR_CODE) LIKE :keyword
        OR LOWER(FLOOR_NAME) LIKE :keyword
        ORDER BY FLOOR_DISPLAY_ORDER
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingFloors{}, ""), 1)
    list := make([]model.WayFindingFloors, 0)
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

func (s *WayFindingService) DeleteFloorsByFloorCode(floorCode string) error {
    query := `DELETE FROM WAY_FINDING_FLOORS WHERE FLOOR_CODE = :floorCode`
    _, err := s.db.ExecContext(s.ctx, query, floorCode)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) FindFloorsByFloorCode(floorCode string) (*model.WayFindingFloors, error) {
    query := `SELECT * FROM WAY_FINDING_FLOORS WHERE FLOOR_CODE = :floorCode`
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingFloors{}, ""), 1)
    var o model.WayFindingFloors
    err := s.db.GetContext(s.ctx, &o, query, floorCode)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}
