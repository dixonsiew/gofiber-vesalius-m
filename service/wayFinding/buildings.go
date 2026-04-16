package wayFinding

import (
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

func (s *WayFindingService) SaveBuilding(o *model.WayFindingBuildings) error {
    query := `
        INSERT INTO BUILDINGS
         (BUILDING_CODE, BUILDING_NAME, BUILDING_DISPLAY_ORDER)
         VALUES
        (:buildingCode, :buildingName, :buildingDisplayOrder)
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("buildingCode", o.BuildingCode.String),
        sql.Named("buildingName", o.BuildingName.String),
        sql.Named("buildingDisplayOrder", o.BuildingDisplayOrder.Int32),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) UpdateBuilding(o *model.WayFindingBuildings) error {
    query := `
        UPDATE WAY_FINDING_BUILDINGS SET
          BUILDING_CODE = :buildingCode,
          BUILDING_NAME = :buildingName,
          BUILDING_DISPLAY_ORDER = :buildingDisplayOrder
        WHERE BUILDING_CODE = :buildingCode
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("buildingCode", o.BuildingCode.String),
        sql.Named("buildingName", o.BuildingName.String),
        sql.Named("buildingDisplayOrder", o.BuildingDisplayOrder.Int32),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) ListBuildings(page string, limit string) (*model.PagedList, error) {
    total, err := s.BuildingsCount(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllBuildings(pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) BuildingsCount(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(*) AS COUNT FROM WAY_FINDING_BUILDINGS`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, err
}

func (s *WayFindingService) FindAllBuildings(offset int, limit int, conn *sqlx.DB) ([]model.WayFindingBuildings, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM WAY_FINDING_BUILDINGS ORDER BY BUILDING_DISPLAY_ORDER OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingBuildings{}, ""), 1)
    list := make([]model.WayFindingBuildings, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, err
}

func (s *WayFindingService) DeleteBuildingsByBuildingCode(buildingCode string) error {
    query := `DELETE FROM WAY_FINDING_BUILDINGS WHERE BUILDING_CODE = :buildingCode`
    _, err := s.db.ExecContext(s.ctx, query, buildingCode)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) FindBuildingsByBuildingCode(buildingCode string) (*model.WayFindingBuildings, error) {
    query := `SELECT * FROM WAY_FINDING_BUILDINGS WHERE BUILDING_CODE = :buildingCode`
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingBuildings{}, ""), 1)
    var o model.WayFindingBuildings
    err := s.db.GetContext(s.ctx, &o, query, buildingCode)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, err
}
