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
