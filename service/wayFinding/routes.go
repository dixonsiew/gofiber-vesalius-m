package wayfinding

import (
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

func (s *WayFindingService) SaveRoute(o *model.WayFindingRoutes) error {
    query := `
        INSERT INTO WAY_FINDING_ROUTES
        (ROUTE_FROM_LOC_ID, ROUTE_TO_LOC_ID, ROUTE_FROM_IMAGE_RAW, ROUTE_TO_IMAGE_RAW)
        VALUES
        (:routeFromLocationId, :routeToLocationId, :routeFromImageRaw, :routeToImageRaw)
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("routeFromLocationId", o.RouteFromLocationId.Int64),
        sql.Named("routeToLocationId", o.RouteToLocationId.Int64),
        sql.Named("routeFromImageRaw", o.RouteFromImageRaw.String),
        sql.Named("routeToImageRaw", o.RouteToImageRaw.String),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) UpdateRoute(o *model.WayFindingRoutes) error {
    query := `
        UPDATE WAY_FINDING_ROUTES SET
          ROUTE_FROM_LOC_ID = :routeFromLocationId,
          ROUTE_TO_LOC_ID = :routeToLocationId,
          ROUTE_FROM_IMAGE_RAW = :routeFromImageRaw,
          ROUTE_TO_IMAGE_RAW = :routeToImageRaw
        WHERE ROUTE_ID = :routeId
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("routeFromLocationId", o.RouteFromLocationId.Int64),
        sql.Named("routeToLocationId", o.RouteToLocationId.Int64),
        sql.Named("routeFromImageRaw", o.RouteFromImageRaw.String),
        sql.Named("routeToImageRaw", o.RouteToImageRaw.String),
        sql.Named("routeId", o.RouteId.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) ListRoutes(page string, limit string) (*model.PagedList, error) {
    total, err := s.RoutesCount(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllWayFindingRoutes(pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) RoutesCount(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(*) AS COUNT FROM WAY_FINDING_ROUTES`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *WayFindingService) FindAllWayFindingRoutes(offset int, limit int, conn *sqlx.DB) ([]model.WayFindingRoutes, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT
          wfroute.ROUTE_ID,
          wfroute.ROUTE_FROM_LOC_ID,
          wfroute.ROUTE_FROM_IMAGE_RAW,
          wflocfrom.LOCATION_BUILDING_CODE AS FROM_BUILDING_CODE,
          wfbldfrom.BUILDING_NAME AS FROM_BUILDING_NAME,
          wflocfrom.LOCATION_FLOOR_CODE AS FROM_FLOOR_CODE,
          wflocfrom.LOCATION_TYPE_CODE AS FROM_TYPE_CODE,
          wflocfrom.LOCATION_CODE AS FROM_CODE,
          wflocfrom.LOCATION_NAME AS FROM_NAME,
          wfroute.ROUTE_TO_LOC_ID,
          wfroute.ROUTE_TO_IMAGE_RAW,
          wflocto.LOCATION_BUILDING_CODE AS TO_BUILDING_CODE,
          wfbldto.BUILDING_NAME AS TO_BUILDING_NAME,
          wflocto.LOCATION_FLOOR_CODE AS TO_FLOOR_CODE,
          wflocto.LOCATION_TYPE_CODE AS TO_TYPE_CODE,
          wflocto.LOCATION_CODE AS TO_CODE,
          wflocto.LOCATION_NAME AS TO_NAME
        FROM WAY_FINDING_ROUTES wfroute
         JOIN WAY_FINDING_LOCATIONS wflocfrom ON wfroute.ROUTE_FROM_LOC_ID = wflocfrom.LOCATION_ID
         JOIN WAY_FINDING_LOCATIONS wflocto ON wfroute.ROUTE_TO_LOC_ID = wflocto.LOCATION_ID
         JOIN WAY_FINDING_BUILDINGS wfbldfrom ON wflocfrom.LOCATION_BUILDING_CODE = wfbldfrom.BUILDING_CODE
         JOIN WAY_FINDING_BUILDINGS wfbldto ON wflocto.LOCATION_BUILDING_CODE = wfbldto.BUILDING_CODE
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    list := make([]model.WayFindingRoutes, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *WayFindingService) ListRoutesByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string, page string, limit string) (*model.PagedList, error) {
    total, err := s.RoutesCountByKeyword(keyword, keyword2, keyword3, keyword4, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllWayFindingRoutesByKeyword(keyword, keyword2, keyword3, keyword4, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *WayFindingService) RoutesCountByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    conds, args := buildRoutesConditions(keyword, keyword2, keyword3, keyword4)
    base := `
        SELECT COUNT(*) AS COUNT
        FROM WAY_FINDING_ROUTES wfroute
        JOIN WAY_FINDING_LOCATIONS wflocfrom ON wfroute.ROUTE_FROM_LOC_ID = wflocfrom.LOCATION_ID
        JOIN WAY_FINDING_LOCATIONS wflocto ON wfroute.ROUTE_TO_LOC_ID = wflocto.LOCATION_ID
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

func (s *WayFindingService) FindAllWayFindingRoutesByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string, offset int, limit int, conn *sqlx.DB) ([]model.WayFindingRoutes, error) {
    db := database.GetFromCon(conn, s.db)
    conds, args := buildRoutesConditions(keyword, keyword2, keyword3, keyword4)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    base := `
        SELECT
          wfroute.ROUTE_ID,
          wfroute.ROUTE_FROM_LOC_ID,
          wflocfrom.LOCATION_BUILDING_CODE AS FROM_BUILDING_CODE,
          wfbldfrom.BUILDING_NAME AS FROM_BUILDING_NAME,
          wflocfrom.LOCATION_FLOOR_CODE AS FROM_FLOOR_CODE,
          wflocfrom.LOCATION_TYPE_CODE AS FROM_TYPE_CODE,
          wflocfrom.LOCATION_CODE AS FROM_CODE,
          wflocfrom.LOCATION_NAME AS FROM_NAME,
          wfroute.ROUTE_TO_LOC_ID,
          wflocto.LOCATION_BUILDING_CODE AS TO_BUILDING_CODE,
          wfbldto.BUILDING_NAME AS TO_BUILDING_NAME,
          wflocto.LOCATION_FLOOR_CODE AS TO_FLOOR_CODE,
          wflocto.LOCATION_TYPE_CODE AS TO_TYPE_CODE,
          wflocto.LOCATION_CODE AS TO_CODE,
          wflocto.LOCATION_NAME AS TO_NAME
        FROM WAY_FINDING_ROUTES wfroute
         JOIN WAY_FINDING_LOCATIONS wflocfrom ON wfroute.ROUTE_FROM_LOC_ID = wflocfrom.LOCATION_ID
         JOIN WAY_FINDING_LOCATIONS wflocto ON wfroute.ROUTE_TO_LOC_ID = wflocto.LOCATION_ID
         JOIN WAY_FINDING_BUILDINGS wfbldfrom ON wflocfrom.LOCATION_BUILDING_CODE = wfbldfrom.BUILDING_CODE
         JOIN WAY_FINDING_BUILDINGS wfbldto ON wflocto.LOCATION_BUILDING_CODE = wfbldto.BUILDING_CODE
    `
    query := base + whereClause(conds) +
        ` OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    list := make([]model.WayFindingRoutes, 0)
    err := db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *WayFindingService) DeleteRoutesByRouteId(routeId int64) error {
    query := `DELETE FROM WAY_FINDING_ROUTES WHERE ROUTE_ID = :routeId`
    _, err := s.db.ExecContext(s.ctx, query, routeId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *WayFindingService) FindRoutesByRouteId(routeId int64) (*model.WayFindingRoutes, error) {
    m := map[string]string{
        "r.FROM_BUILDING_CODE": "",
        "r.FROM_BUILDING_NAME": "",
        "r.TO_BUILDING_CODE": "",
        "r.TO_BUILDING_NAME": "",
    }
    query := `
        SELECT 
          r.*,
          fromLoc.LOCATION_BUILDING_CODE AS FROM_BUILDING_CODE,
          fromBld.BUILDING_NAME AS FROM_BUILDING_NAME,
          toLoc.LOCATION_BUILDING_CODE AS TO_BUILDING_CODE,
          toBld.BUILDING_NAME AS TO_BUILDING_NAME
        FROM WAY_FINDING_ROUTES r
         LEFT JOIN WAY_FINDING_LOCATIONS fromLoc ON r.ROUTE_FROM_LOC_ID = fromLoc.LOCATION_ID
         LEFT JOIN WAY_FINDING_BUILDINGS fromBld ON fromLoc.LOCATION_BUILDING_CODE = fromBld.BUILDING_CODE
         LEFT JOIN WAY_FINDING_LOCATIONS toLoc ON r.ROUTE_TO_LOC_ID = toLoc.LOCATION_ID
         LEFT JOIN WAY_FINDING_BUILDINGS toBld ON toLoc.LOCATION_BUILDING_CODE = toBld.BUILDING_CODE
        WHERE r.ROUTE_ID = :routeId
    `
    query = strings.Replace(query, "r.*", utils.GetDbColsWithReplace(model.WayFindingRoutes{}, "r.", m), 1)
    var o model.WayFindingRoutes
    err := s.db.GetContext(s.ctx, &o, query, routeId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}
