package wayfinding

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

var WayFindingSvc *WayFindingService = NewWayFindingService(database.GetDb(), database.GetCtx())

type WayFindingService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewWayFindingService(db *sqlx.DB, ctx context.Context) *WayFindingService {
    return &WayFindingService{
        db:  db,
        ctx: ctx,
    }
}

func (s *WayFindingService) ExistsByBuildingCode(buildingCode string) (bool, error) {
    query := `
        SELECT COUNT(*) AS COUNT FROM DUAL 
        WHERE EXISTS (SELECT 1 FROM WAY_FINDING_LOCATIONS WHERE BUILDING_CODE = :buildingCode)
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, buildingCode)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, err
}

func (s *WayFindingService) ExistsByFloorCode(floorCode string) (bool, error) {
    query := `
        SELECT COUNT(*) AS COUNT FROM DUAL 
        WHERE EXISTS (SELECT 1 FROM WAY_FINDING_LOCATIONS WHERE LOCATION_FLOOR_CODE = :floorCode)
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, floorCode)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, err
}

func (s *WayFindingService) ExistsByLocationTypeCode(locationTypeCode string) (bool, error) {
    query := `
        SELECT COUNT(*) AS COUNT FROM DUAL 
        WHERE EXISTS (SELECT 1 FROM WAY_FINDING_LOCATIONS WHERE LOCATION_TYPE_CODE = :locationTypeCode)
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, floorCode)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, err
}

func (s *WayFindingService) ExistsByRouteFromToLocationId(fromLocId int64, toLocId int64) (bool, error) {
    query := `
        SELECT COUNT(*) AS COUNT FROM DUAL 
        WHERE EXISTS (SELECT 1 FROM WAY_FINDING_ROUTES WHERE ROUTE_FROM_LOC_ID = :fromLocId AND ROUTE_TO_LOC_ID = :toLocId)
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, fromLocId, toLocId)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, err
}

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

// func (s *WayFindingService) 

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

func (s *WayFindingService) FindAllWayFindingBuildings(offset int, limit int, conn *sqlx.DB) ([]model.WayFindingBuildings, error) {
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
        utils.LogError(err)
        return nil, err
    }
    return &o, err
}

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

func (s *WayFindingService) FindAllWayFindingFloors(conn *sqlx.DB) ([]model.WayFindingFloors, error) {
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

func (s *WayFindingService) FindAllWayFindingFloorsWebAdmin(offset int, limit int, conn *sqlx.DB) ([]model.WayFindingFloors, error) {
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

func (s *WayFindingService) FindAllWayFindingFloorsByKeyword(keyword string, offset int, limit int, conn *sqlx.DB) ([]model.WayFindingFloors, error) {
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

func (s *WayFindingService) FindAllWayFindingFloorsByKeywordWebAdmin(keyword string, offset int, limit int, conn *sqlx.DB) ([]model.WayFindingFloors, error) {
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
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *WayFindingService) SaveLocation(o *model.WayFindingLocations) {
    query := `
        SELECT COUNT(*) AS COUNT FROM WAY_FINDING_LOCATIONS 
        WHERE :locationCode IN (LOCATION_CODE)
    `
    var count int
    err := db.GetContext(s.ctx, &count, query, o.locationCode.String)
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
    _, err := s.db.ExecContext(s.ctx, query,
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
    err := db.GetContext(s.ctx, &count, query, locationId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    if count > 0 {
        return fiber.NewError(fiber.StatusBadRequest, "One of the Routes are using this Location")
    }

    query = `DELETE FROM WAY_FINDING_LOCATIONS WHERE LOCATION_ID = :locationId`
    _, err := s.db.ExecContext(s.ctx, query, locationId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) FindLocationsById(locationId int64) (*model.WayFindingLocations, error) {
    query := `SELECT * FROM WAY_FINDING_LOCATIONS WHERE LOCATION_ID = :locationId`
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingLocations{}, ""), 1)
    var o model.WayFindingLocations
    err := db.GetContext(s.ctx, &o, query, locationId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

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

func (s *WayFindingService) 

func (s *WayFindingService) 

func (s *WayFindingService) 

func (s *WayFindingService) 

func buildLocationsConditions(keyword string, keyword2 string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `(LOWER(wffloor.FLOOR_CODE) LIKE :keyword OR LOWER(wffloor.FLOOR_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        conds = append(conds, `(LOWER(wfloctype.LOCATION_TYPE_CODE) LIKE :keyword2 OR LOWER(wfloctype.LOCATION_TYPE_NAME) LIKE :keyword2)`)
        args = append(args, sql.Named("keyword2", strings.ToLower(keyword2)))
    }
    if keyword3 != "" {
        conds = append(conds, `(LOWER(wfloc.LOCATION_CODE) LIKE :keyword3 OR LOWER(wfloc.LOCATION_NAME) LIKE :keyword3)`)
        args = append(args, sql.Named("keyword3", strings.ToLower(keyword3)))
    }
    if keyword4 != "" {
        conds = append(conds, `(LOWER(wfbuilding.BUILDING_CODE) LIKE :keyword4 OR LOWER(wfbuilding.BUILDING_NAME) LIKE :keyword4)`)
        args = append(args, sql.Named("keyword4", strings.ToLower(keyword4)))
    }

    return conds, args
}

func whereClause(conds []string) string {
    if len(conds) == 0 {
        return ""
    }
    return " WHERE " + strings.Join(conds, " AND ")
}
