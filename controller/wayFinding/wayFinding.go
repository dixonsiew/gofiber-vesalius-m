package wayfinding

import (
    "fmt"
    "strconv"
    "vesaliusm/dto"

    //"vesaliusm/middleware"
    "vesaliusm/model"
    "vesaliusm/service/wayfinding"
    "vesaliusm/utils"
    "vesaliusm/utils/constants"

    "github.com/gofiber/fiber/v3"
)

type WayFindingController struct {
    wayFindingService *wayfinding.WayFindingService
}

func NewWayFindingController() *WayFindingController {
    return &WayFindingController{
        wayFindingService: wayfinding.WayFindingSvc,
    }
}

// GetDropdowns
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param    request    body    dto.DropdownsDto    true    "DropdownsDto"
// @Success 200
// @Router /way-finding/dropdown [post]
func (cr *WayFindingController) GetDropdowns(c fiber.Ctx) error {
    data := new(dto.DropdownsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    m, err := cr.wayFindingService.ListDropdowns(data.BuildingCode)
    if err != nil {
        return err
    }

    return c.JSON(m)
}

// CreateBuilding
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param    request    body    dto.BuildingsDto    true    "BuildingsDto"
// @Success 200
// @Router /way-finding/buildings/create [post]
func (cr *WayFindingController) CreateBuilding(c fiber.Ctx) error {
    data := new(dto.BuildingsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o := &model.WayFindingBuildings{
        BuildingCode:         utils.NewNullString(data.BuildingCode),
        BuildingName:         utils.NewNullString(data.BuildingName),
        BuildingDisplayOrder: utils.NewInt32(int32(data.BuildingDisplayOrder)),
    }
    err := cr.wayFindingService.SaveBuilding(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Building created",
    })
}

// UpdateBuilding
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param    request    body    dto.BuildingsDto    true    "BuildingsDto"
// @Success 200
// @Router /way-finding/buildings/update/{buildingCode} [put]
func (cr *WayFindingController) UpdateBuilding(c fiber.Ctx) error {
    buildingCode := c.Params("buildingCode")
    data := new(dto.BuildingsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o := &model.WayFindingBuildings{
        BuildingCode:         utils.NewNullString(buildingCode),
        BuildingName:         utils.NewNullString(data.BuildingName),
        BuildingDisplayOrder: utils.NewInt32(int32(data.BuildingDisplayOrder)),
    }
    err := cr.wayFindingService.UpdateBuilding(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Building updated",
    })
}

// GetAllBuildings
//
// @Tags Way Finding
// @Produce json
// @Param        _page              query      int  false  "_page"  default:"1"
// @Param        _limit             query      int  false  "_limit" default:"10"
// @Success 200 {array} model.WayFindingBuildings
// @Router /way-finding/buildings [get]
func (cr *WayFindingController) GetAllBuildings(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.wayFindingService.ListBuildings(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// DeleteBuildingsByBuildingCode
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param    buildingCode    path    string    true    "buildingCode"
// @Success 200
// @Router /way-finding/buildings/delete/{buildingCode} [delete]
func (cr *WayFindingController) DeleteBuildingsByBuildingCode(c fiber.Ctx) error {
    buildingCode := c.Params("buildingCode")
    b, err := cr.wayFindingService.ExistsByBuildingCode(buildingCode)
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("The Building Code: %s is linked to existing Location records. To delete, you must first delete or update the related Location records", buildingCode))
    }

    err = cr.wayFindingService.DeleteBuildingsByBuildingCode(buildingCode)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": 1,
    })
}

// GetBuildingsByBuildingCode
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param    buildingCode    path    string    true    "buildingCode"
// @Success 200
// @Router /way-finding/buildings/{buildingCode} [get]
func (cr *WayFindingController) GetBuildingsByBuildingCode(c fiber.Ctx) error {
    buildingCode := c.Params("buildingCode")
    o, err := cr.wayFindingService.FindBuildingsByBuildingCode(buildingCode)
    if err != nil {
        return err
    }
    return c.JSON(o)
}

// CreateFloor
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param    request    body    dto.FloorsDto    true    "FloorsDto"
// @Success 200
// @Router /way-finding/floors/create [post]
func (cr *WayFindingController) CreateFloor(c fiber.Ctx) error {
    data := new(dto.FloorsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o := &model.WayFindingFloors{
        FloorCode:     utils.NewNullString(data.FloorCode),
        FloorName:     utils.NewNullString(data.FloorName),
        FloorId:       utils.NewInt64(int64(data.FloorDisplayOrder)),
        FloorImageRaw: utils.NewNullString(data.FloorImageRaw),
    }
    err := cr.wayFindingService.SaveFloor(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Floor created",
    })
}

// UpdateFloor
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param    floorCode    path    string           true    "floorCode"
// @Param    request      body    dto.FloorsDto    true    "FloorsDto"
// @Success 200
// @Router /way-finding/floors/update/{floorCode} [put]
func (cr *WayFindingController) UpdateFloor(c fiber.Ctx) error {
    floorCode := c.Params("floorCode")
    data := new(dto.FloorsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o := &model.WayFindingFloors{
        FloorCode:     utils.NewNullString(floorCode),
        FloorName:     utils.NewNullString(data.FloorName),
        FloorId:       utils.NewInt64(int64(data.FloorDisplayOrder)),
        FloorImageRaw: utils.NewNullString(data.FloorImageRaw),
    }
    err := cr.wayFindingService.UpdateFloor(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Floor updated",
    })
}

// GetAllFloors
//
// @Tags Way Finding
// @Produce json
// @Success 200 {array} model.WayFindingFloors
// @Router /way-finding/floors [get]
func (cr *WayFindingController) GetAllFloors(c fiber.Ctx) error {
    lx, err := cr.wayFindingService.ListFloors()
    if err != nil {
        return err
    }
    return c.JSON(lx)
}

// GetAllFloorsWebAdmin
//
// @Tags Way Finding
// @Produce json
// @Param        _page              query      int  false  "_page"  default:"1"
// @Param        _limit             query      int  false  "_limit" default:"10"
// @Success 200 {array} model.WayFindingFloors
// @Router /way-finding/floors/webadmin [get]
func (cr *WayFindingController) GetAllFloorsWebAdmin(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.wayFindingService.ListFloorsWebAdmin(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllFloors
//
// @Tags Way Finding
// @Produce json
// @Param        _page             query       int                   false  "_page"  default:"1"
// @Param        _limit            query       int                   false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.WayFindingFloors
// @Router /way-finding/floors [post]
func (cr *WayFindingController) SearchAllFloors(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    if key != "" {
        key = "%" + key + "%"
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.wayFindingService.ListFloorsByKeyword(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllFloorsWebAdmin
//
// @Tags Way Finding
// @Produce json
// @Param        _page             query       int                   false  "_page"  default:"1"
// @Param        _limit            query       int                   false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.WayFindingFloors
// @Router /way-finding/floors/webadmin [post]
func (cr *WayFindingController) SearchAllFloorsWebAdmin(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    if key != "" {
        key = "%" + key + "%"
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.wayFindingService.ListFloorsByKeywordWebAdmin(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// DeleteFloorsByFloorCode
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param        floorCode             path       string                   true  "floorCode"
// @Success 200
// @Router /way-finding/floors/delete/{floorCode} [delete]
func (cr *WayFindingController) DeleteFloorsByFloorCode(c fiber.Ctx) error {
    floorCode := c.Params("floorCode")
    b, err := cr.wayFindingService.ExistsByFloorCode(floorCode)
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("The Floor Code: %s is linked to existing Location records. To delete, you must first delete or update the related Location records", floorCode))
    }

    err = cr.wayFindingService.DeleteFloorsByFloorCode(floorCode)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": 1,
    })
}

// GetFloorsByFloorCode
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param        floorCode             path       string                   true  "floorCode"
// @Success 200
// @Router /way-finding/floors/{floorCode} [get]
func (cr *WayFindingController) GetFloorsByFloorCode(c fiber.Ctx) error {
    floorCode := c.Params("floorCode")
    o, err := cr.wayFindingService.FindFloorsByFloorCode(floorCode)
    if err != nil {
        return err
    }
    return c.JSON(o)
}

// CreateLocation
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param    request      body    dto.LocationsDto    true    "LocationsDto"
// @Success 200
// @Router /way-finding/locations/create [post]
func (cr *WayFindingController) CreateLocation(c fiber.Ctx) error {
    data := new(dto.LocationsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o := &model.WayFindingLocations{
        LocationBuildingCode: utils.NewNullString(data.LocationBuildingCode),
        LocationFloorCode:    utils.NewNullString(data.LocationFloorCode),
        LocationTypeCode:     utils.NewNullString(data.LocationTypeCode),
        LocationCode:         utils.NewNullString(data.LocationCode),
        LocationName:         utils.NewNullString(data.LocationName),
    }
    err := cr.wayFindingService.SaveLocation(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Location created",
    })
}

// UpdateLocation
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param    locationId    path    int                 true    "locationId"
// @Param    request       body    dto.LocationsDto    true    "LocationsDto"
// @Success 200
// @Router /way-finding/locations/update/{locationId} [put]
func (cr *WayFindingController) UpdateLocation(c fiber.Ctx) error {
    locationId := c.Params("locationId")
    ilocationId, _ := strconv.ParseInt(locationId, 10, 64)
    data := new(dto.LocationsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o := &model.WayFindingLocations{
        LocationId:           utils.NewInt64(ilocationId),
        LocationBuildingCode: utils.NewNullString(data.LocationBuildingCode),
        LocationFloorCode:    utils.NewNullString(data.LocationFloorCode),
        LocationTypeCode:     utils.NewNullString(data.LocationTypeCode),
        LocationCode:         utils.NewNullString(data.LocationCode),
        LocationName:         utils.NewNullString(data.LocationName),
    }
    err := cr.wayFindingService.UpdateLocation(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Location updated",
    })
}

// GetAllLocations
//
// @Tags Way Finding
// @Produce json
// @Param        _page             query       int                   false  "_page"  default:"1"
// @Param        _limit            query       int                   false  "_limit" default:"10"
// @Success 200 {array} model.WayFindingLocations
// @Router /way-finding/locations [get]
func (cr *WayFindingController) GetAllLocations(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.wayFindingService.ListLocations(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllLocations
//
// @Tags Way Finding
// @Produce json
// @Param        _page             query       int                      false  "_page"  default:"1"
// @Param        _limit            query       int                      false  "_limit" default:"10"
// @Param        request           body        dto.SearchKeyword4Dto    false  "Keyword"
// @Success 200 {array} model.WayFindingLocations
// @Router /way-finding/locations [post]
func (cr *WayFindingController) SearchAllLocations(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key2 := data.GetString("keyword2")
    key3 := data.GetString("keyword3")
    key4 := data.GetString("keyword4")

    if key != "" {
        key = "%" + key + "%"
    }
    if key2 != "" {
        key2 = "%" + key2 + "%"
    }
    if key3 != "" {
        key3 = "%" + key3 + "%"
    }
    if key4 != "" {
        key4 = "%" + key4 + "%"
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.wayFindingService.ListLocationsByKeyword(key, key2, key3, key4, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// DeleteLocationsById
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param        locationId        path        int                      true  "locationId"
// @Success 200
// @Router /way-finding/locations/delete/{locationId} [delete]
func (cr *WayFindingController) DeleteLocationsById(c fiber.Ctx) error {
    locationId := c.Params("locationId")
    ilocationId, _ := strconv.ParseInt(locationId, 10, 64)
    err := cr.wayFindingService.DeleteLocationsById(ilocationId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": 1,
    })
}

// GetLocationsById
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param        locationId        path        int                      true  "locationId"
// @Success 200
// @Router /way-finding/locations/{locationId} [get]
func (cr *WayFindingController) GetLocationsById(c fiber.Ctx) error {
    locationId := c.Params("locationId")
    ilocationId, _ := strconv.ParseInt(locationId, 10, 64)
    o, err := cr.wayFindingService.FindLocationsById(ilocationId)
    if err != nil {
        return err
    }
    return c.JSON(o)
}

// CreateLocationType
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param    request       body    dto.LocationTypesDto    true    "LocationTypesDto"
// @Success 200
// @Router /way-finding/location-types/create [post]
func (cr *WayFindingController) CreateLocationType(c fiber.Ctx) error {
    data := new(dto.LocationTypesDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o := model.WayFindingLocationTypes{
        LocationTypeCode: utils.NewNullString(data.LocationTypeCode),
        LocationTypeName: utils.NewNullString(data.LocationTypeName),
        LocationTypeId:   utils.NewInt64(int64(data.LocationTypeDisplayOrder)),
    }
    b, err := cr.wayFindingService.ExistsByLocationTypeCode(data.LocationTypeCode)
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("The Location Type Code: %s already existed in our records. Please try again with different Location Type Code", data.LocationTypeCode))
    }

    err = cr.wayFindingService.SaveLocationType(&o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Location Type created",
    })
}

// UpdateLocationType
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param    locationTypeCode        path        string                  true    "locationTypeCode"
// @Param    request                 body        dto.LocationTypesDto    true    "LocationTypesDto"
// @Success 200
// @Router /way-finding/location-types/update/{locationTypeCode} [put]
func (cr *WayFindingController) UpdateLocationType(c fiber.Ctx) error {
    locationTypeCode := c.Params("locationTypeCode")
    data := new(dto.LocationTypesDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o := &model.WayFindingLocationTypes{
        LocationTypeName: utils.NewNullString(data.LocationTypeName),
        LocationTypeDisplayOrder: utils.NewInt32(int32(data.LocationTypeDisplayOrder)),
    }
    err := cr.wayFindingService.UpdateLocationType(o, locationTypeCode)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Location Type updated",
    })
}

// GetAllLocationTypes
//
// @Tags Way Finding
// @Produce json
// @Param        _page             query       int  false  "_page"  default:"1"
// @Param        _limit            query       int  false  "_limit" default:"10"
// @Success 200 {array} model.WayFindingLocationTypes
// @Router /way-finding/location-types [get]
func (cr *WayFindingController) GetAllLocationTypes(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.wayFindingService.ListLocationTypes(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllLocationTypes
//
// @Tags Way Finding
// @Produce json
// @Param        _page             query       int                   false  "_page"  default:"1"
// @Param        _limit            query       int                   false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.WayFindingLocationTypes
// @Router /way-finding/location-types [post]
func (cr *WayFindingController) SearchAllLocationTypes(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    if key != "" {
        key = "%" + key + "%"
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.wayFindingService.ListLocationTypesByKeyword(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// DeleteLocationTypesByLocationTypeCode
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param        locationTypeCode  path        string  true  "locationTypeCode"
// @Success 200
// @Router /way-finding/location-types/delete/{locationTypeCode} [delete]
func (cr *WayFindingController) DeleteLocationTypesByLocationTypeCode(c fiber.Ctx) error {
    locationTypeCode := c.Params("locationTypeCode")
    b, err := cr.wayFindingService.ExistsByLocationTypeCode(locationTypeCode)
    if err != nil {
        return err
    }
    
    if b {
        return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("The Location Type Code: %s is linked to existing Location records. To delete, you must first delete or update the related Location records", locationTypeCode))
    }
    
    err = cr.wayFindingService.DeleteLocationTypesByLocationTypeCode(locationTypeCode)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": 1,
    })
}

// GetLocationTypesByLocationTypeCode
//
// @Tags Way Finding
// @Produce json
// @Security BearerAuth
// @Param        locationTypeCode  path        string  true  "locationTypeCode"
// @Success 200
// @Router /way-finding/location-types/{locationTypeCode} [get]
func (cr *WayFindingController) GetLocationTypesByLocationTypeCode(c fiber.Ctx) error {
    locationTypeCode := c.Params("locationTypeCode")
    o, err := cr.wayFindingService.FindLocationTypesByLocationTypeCode(locationTypeCode)
    if err != nil {
        return err
    }
    return c.JSON(o)
}

func (cr *WayFindingController) xxx(c fiber.Ctx) error {

}

func (cr *WayFindingController) xxx(c fiber.Ctx) error {

}

func (cr *WayFindingController) xxx(c fiber.Ctx) error {

}
