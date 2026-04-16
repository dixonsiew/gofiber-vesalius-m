package wayFinding

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/dto"
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

func (s *WayFindingService) ListDropdowns(buildingCode string) (fiber.Map, error) {
    var err error
    locations := []model.WayFindingLocations{}
    buildings, err := s.FindAllBuildings(0, 100, s.db)
    if err != nil {
        return nil, err
    }

    floors, err := s.FindAllFloorsWebAdmin(0, 100, s.db)
    if err != nil {
        return nil, err
    }

    types, err := s.FindAllLocationTypes(0, 100, s.db)
    if err != nil {
        return nil, err
    }

    if buildingCode != "" {
        x := dto.SearchKeyword4Dto{
            Keyword4: buildingCode,
        }
        locations, err = s.FindAllLocationsByKeyword(x, 0, 100, s.db)
    } else {
        locations, err = s.FindAllLocations(0, 100, s.db)
    }

    if err != nil {
        return nil, err
    }

    return fiber.Map{
        "buildings": buildings,
        "floors": floors,
        "types": types,
        "locations": locations,
    }, nil
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
    err := s.db.GetContext(s.ctx, &count, query, locationTypeCode)
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

func (s *WayFindingService) FindAllLocationsByLocationTypeCode(code string, offset int, limit int, conn *sqlx.DB) ([]model.WayFindingLocations, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM WAY_FINDING_LOCATIONS WHERE LOCATION_TYPE_CODE = :code OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingLocations{}, ""), 1)
    list := make([]model.WayFindingLocations, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *WayFindingService) FindAllLocationsByTypeCodeByKeyword(code string, keyword string, offset int, limit int, conn *sqlx.DB) ([]model.WayFindingLocations, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT * FROM WAY_FINDING_LOCATIONS
        WHERE LOCATION_TYPE_CODE = :code
        AND LOWER(LOCATION_NAME) LIKE :keyword
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingLocations{}, ""), 1)
    list := make([]model.WayFindingLocations, 0)
    err := db.SelectContext(s.ctx, &list, query, code, strings.ToLower(keyword), offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *WayFindingService) FindRoutes(fromId int64, toId int64) (*model.WayFindingRoutes, error) {
    query := `SELECT * FROM WAY_FINDING_ROUTES WHERE ROUTE_FROM_LOC_ID = :fromId AND ROUTE_TO_LOC_ID = :toId`
    query = strings.Replace(query, "*", utils.GetDbCols(model.WayFindingRoutes{}, ""), 1)
    var o model.WayFindingRoutes
    err := s.db.GetContext(s.ctx, &o, query, fromId, toId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *WayFindingService) FindLocationByLocationIdAndLocationTypeId(locationId int64, locationTypeId int64) (*model.WayFindingLocations, error) {
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

func (s *WayFindingService) ListLocationByTypeCode(code string, page string, limit string) (*model.PagedList, error) {
    total, err := s.LocationByTypeCodeCount(code, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllLocationsByLocationTypeCode(code, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) LocationByTypeCodeCount(code string, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(*) AS COUNT FROM WAY_FINDING_LOCATIONS WHERE LOCATION_TYPE_CODE = :code`
    var count int
    err := db.GetContext(s.ctx, &count, query, code)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *WayFindingService) ListLocationByTypeCodeByKeyword(code string, keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.LocationByTypeCodeCountByKeyword(code, keyword, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllLocationsByTypeCodeByKeyword(code, keyword, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) LocationByTypeCodeCountByKeyword(code string, keyword string, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT COUNT(*) AS COUNT 
        FROM WAY_FINDING_LOCATIONS 
        WHERE LOCATION_TYPE_CODE = :code
        AND LOWER(LOCATION_NAME) LIKE :keyword
    `
    var count int
    err := db.GetContext(s.ctx, &count, query, code, strings.ToLower(keyword))
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func buildLocationsConditions(x dto.SearchKeyword4Dto) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if x.Keyword != "" {
        conds = append(conds, `(LOWER(wffloor.FLOOR_CODE) LIKE :keyword OR LOWER(wffloor.FLOOR_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(x.Keyword)))
    }
    if x.Keyword2 != "" {
        conds = append(conds, `(LOWER(wfloctype.LOCATION_TYPE_CODE) LIKE :keyword2 OR LOWER(wfloctype.LOCATION_TYPE_NAME) LIKE :keyword2)`)
        args = append(args, sql.Named("keyword2", strings.ToLower(x.Keyword2)))
    }
    if x.Keyword3 != "" {
        conds = append(conds, `(LOWER(wfloc.LOCATION_CODE) LIKE :keyword3 OR LOWER(wfloc.LOCATION_NAME) LIKE :keyword3)`)
        args = append(args, sql.Named("keyword3", strings.ToLower(x.Keyword3)))
    }
    if x.Keyword4 != "" {
        conds = append(conds, `(LOWER(wfbuilding.BUILDING_CODE) LIKE :keyword4 OR LOWER(wfbuilding.BUILDING_NAME) LIKE :keyword4)`)
        args = append(args, sql.Named("keyword4", strings.ToLower(x.Keyword4)))
    }

    return conds, args
}

func buildRoutesConditions(x dto.SearchKeyword4Dto) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if x.Keyword != "" {
        conds = append(conds, `LOWER(wflocfrom.LOCATION_BUILDING_CODE) LIKE :keyword`)
        args = append(args, sql.Named("keyword", strings.ToLower(x.Keyword)))
    }
    if x.Keyword2 != "" {
        conds = append(conds, `LOWER(wflocfrom.LOCATION_NAME) LIKE :keyword2`)
        args = append(args, sql.Named("keyword2", strings.ToLower(x.Keyword2)))
    }
    if x.Keyword3 != "" {
        conds = append(conds, `LOWER(wflocto.LOCATION_BUILDING_CODE) LIKE :keyword3`)
        args = append(args, sql.Named("keyword3", strings.ToLower(x.Keyword3)))
    }
    if x.Keyword4 != "" {
        conds = append(conds, `LOWER(wflocto.LOCATION_NAME) LIKE :keyword4`)
        args = append(args, sql.Named("keyword4", strings.ToLower(x.Keyword4)))
    }

    return conds, args
}

func whereClause(conds []string) string {
    if len(conds) == 0 {
        return ""
    }
    return " WHERE " + strings.Join(conds, " AND ")
}
