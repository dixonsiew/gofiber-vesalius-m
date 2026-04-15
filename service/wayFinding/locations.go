package wayfinding

import (
    "database/sql"
    "fmt"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/jmoiron/sqlx"
)

func (s *WayFindingService) SaveLocation(o *model.WayFindingLocations) error {
    query := `
        SELECT COUNT(*) AS COUNT FROM WAY_FINDING_LOCATIONS 
        WHERE :locationCode IN (LOCATION_CODE)
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, o.LocationCode.String)
    if err != nil {
        utils.LogError(err)
        return err
    }

    if count > 0 {
        return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("The Location Code: %s already exist in our system. Please try again.", o.LocationCode.String))
    }
    
    query = `
        INSERT INTO WAY_FINDING_LOCATIONS
        (LOCATION_BUILDING_CODE, LOCATION_FLOOR_CODE, LOCATION_TYPE_CODE, LOCATION_CODE, LOCATION_NAME)
        VALUES
        (:locationBuildingCode, :locationFloorCode, :locationTypeCode, :locationCode, :locationName)
    `
    _, err = s.db.ExecContext(s.ctx, query,
        sql.Named("locationBuildingCode", o.LocationBuildingCode.String),
        sql.Named("locationFloorCode", o.LocationFloorCode.String),
        sql.Named("locationTypeCode", o.LocationTypeCode.String),
        sql.Named("locationCode", o.LocationCode.String),
        sql.Named("locationName", o.LocationName.String),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) UpdateLocation(o *model.WayFindingLocations) error {
    query := `
        UPDATE WAY_FINDING_LOCATIONS SET
          LOCATION_BUILDING_CODE = :locationBuildingCode,
          LOCATION_FLOOR_CODE = :locationFloorCode,
          LOCATION_TYPE_CODE = :locationTypeCode,
          LOCATION_CODE = :locationCode,
          LOCATION_NAME = :locationName
        WHERE LOCATION_ID = :locationId
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("locationBuildingCode", o.LocationBuildingCode.String),
        sql.Named("locationFloorCode", o.LocationFloorCode.String),
        sql.Named("locationTypeCode", o.LocationTypeCode.String),
        sql.Named("locationCode", o.LocationCode.String),
        sql.Named("locationName", o.LocationName.String),
        sql.Named("locationId", o.LocationId.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) ListLocations(page string, limit string) (*model.PagedList, error) {
    total, err := s.LocationsCount(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllWayFindingLocations(pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) LocationsCount(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    var count int
    query := `SELECT COUNT(*) AS COUNT FROM WAY_FINDING_LOCATIONS`
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *WayFindingService) FindAllWayFindingLocations(offset int, limit int, conn *sqlx.DB) ([]model.WayFindingLocations, error) {
    query := `
        SELECT
          wfloc.LOCATION_ID,
          wfbuilding.BUILDING_CODE AS LOCATION_BUILDING_CODE,
          wfbuilding.BUILDING_NAME AS LOCATION_BUILDING_NAME,
          wffloor.FLOOR_CODE AS LOCATION_FLOOR_CODE, 
          wffloor.FLOOR_NAME AS LOCATION_FLOOR_NAME,
          wfloctype.LOCATION_TYPE_CODE AS LOCATION_TYPE_CODE, 
          wfloctype.LOCATION_TYPE_NAME AS LOCATION_TYPE_NAME,
          wfloc.LOCATION_CODE, 
          wfloc.LOCATION_NAME
        FROM WAY_FINDING_LOCATIONS wfloc
         JOIN WAY_FINDING_BUILDINGS wfbuilding ON wfloc.LOCATION_BUILDING_CODE = wfbuilding.BUILDING_CODE
         JOIN WAY_FINDING_FLOORS wffloor ON wfloc.LOCATION_FLOOR_CODE = wffloor.FLOOR_CODE
         JOIN WAY_FINDING_LOCATION_TYPE wfloctype ON wfloc.LOCATION_TYPE_CODE = wfloctype.LOCATION_TYPE_CODE 
        ORDER BY LOCATION_NAME OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    list := make([]model.WayFindingLocations, 0)
    err := conn.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *WayFindingService) ListLocationsByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string, page string, limit string) (*model.PagedList, error) {
    total, err := s.LocationsCountByKeyword(keyword, keyword2, keyword3, keyword4, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllWayFindingLocationsByKeyword(keyword, keyword2, keyword3, keyword4, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) LocationsCountByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    conds, args := buildLocationsConditions(keyword, keyword2, keyword3, keyword4)
    base := `
        SELECT COUNT(*) AS COUNT
        FROM WAY_FINDING_LOCATIONS wfloc
        JOIN WAY_FINDING_BUILDINGS wfbuilding ON wfloc.LOCATION_BUILDING_CODE = wfbuilding.BUILDING_CODE
        JOIN WAY_FINDING_FLOORS wffloor ON wfloc.LOCATION_FLOOR_CODE = wffloor.FLOOR_CODE
        JOIN WAY_FINDING_LOCATION_TYPE wfloctype ON wfloc.LOCATION_TYPE_CODE = wfloctype.LOCATION_TYPE_CODE
    `
    query := base + whereClause(conds)

    var count int
    err := db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *WayFindingService) FindAllWayFindingLocationsByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string, offset int, limit int, conn *sqlx.DB) ([]model.WayFindingLocations, error) {
    db := database.GetFromCon(conn, s.db)
    conds, args := buildLocationsConditions(keyword, keyword2, keyword3, keyword4)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    base := `
        SELECT
          wfloc.LOCATION_ID,
          wfbuilding.BUILDING_CODE AS LOCATION_BUILDING_CODE,
          wfbuilding.BUILDING_NAME AS LOCATION_BUILDING_NAME,
          wffloor.FLOOR_CODE AS LOCATION_FLOOR_CODE, 
          wffloor.FLOOR_NAME AS LOCATION_FLOOR_NAME,
          wfloctype.LOCATION_TYPE_CODE AS LOCATION_TYPE_CODE, 
          wfloctype.LOCATION_TYPE_NAME AS LOCATION_TYPE_NAME,
          wfloc.LOCATION_CODE, 
          wfloc.LOCATION_NAME
        FROM WAY_FINDING_LOCATIONS wfloc
         JOIN WAY_FINDING_BUILDINGS wfbuilding ON wfloc.LOCATION_BUILDING_CODE = wfbuilding.BUILDING_CODE
         JOIN WAY_FINDING_FLOORS wffloor ON wfloc.LOCATION_FLOOR_CODE = wffloor.FLOOR_CODE
         JOIN WAY_FINDING_LOCATION_TYPE wfloctype ON wfloc.LOCATION_TYPE_CODE = wfloctype.LOCATION_TYPE_CODE
    `
    query := base + whereClause(conds) +
        ` OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    list := make([]model.WayFindingLocations, 0)
    err := db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *WayFindingService) DeleteLocationsById(locationId int64) error {
    query := `
        SELECT COUNT(*) AS COUNT FROM WAY_FINDING_ROUTES 
        WHERE :locationId IN (ROUTE_FROM_LOC_ID, ROUTE_TO_LOC_ID)
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, locationId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    if count > 0 {
        return fiber.NewError(fiber.StatusBadRequest, "One of the Routes are using this Location")
    }

    query = `DELETE FROM WAY_FINDING_LOCATIONS WHERE LOCATION_ID = :locationId`
    _, err = s.db.ExecContext(s.ctx, query, locationId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) FindLocationsById(locationId int64) (*model.WayFindingLocations, error) {
    query := `SELECT * FROM WAY_FINDING_LOCATIONS WHERE LOCATION_ID = :locationId`
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingLocations{}, ""), 1)
    var o model.WayFindingLocations
    err := s.db.GetContext(s.ctx, &o, query, locationId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}
