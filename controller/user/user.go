package user

import (
	"strconv"
	"vesaliusm/dto"
	"vesaliusm/middleware"
	"vesaliusm/service/applicationUser"
	"vesaliusm/utils"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	applicationUserService *applicationUser.ApplicationUserService
}

func NewUserController() *UserController {
	return &UserController{
		applicationUserService: applicationUser.ApplicationUserSvc,
	}
}

// GetAllUsers
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
// @Success 200 {array} model.ApplicationUser
// @Router /user/all [get]
func (cr *UserController) GetAllUsers(c fiber.Ctx) error {
	page := c.Query("_page", "1")
	limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
	m, err := cr.applicationUserService.List(page, limit)
	if err != nil {
		return err
	}

	c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
	c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
	return c.JSON(m.List)
}

// GetAllActiveUsers
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
// @Success 200 {array} model.ApplicationUser
// @Router /user/all/active [get]
func (cr *UserController) GetAllActiveUsers(c fiber.Ctx) error {
	page := c.Query("_page", "1")
	limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
	m, err := cr.applicationUserService.ListActive(page, limit)
	if err != nil {
		return err
	}

	c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
	c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
	return c.JSON(m.List)
}

// UpdatePlayerId
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        playerId              path      string  true  "PlayerId"
// @Success 200
// @Router /user/update-playerid/{playerId} [post]
func (cr *UserController) UpdatePlayerId(c fiber.Ctx) error {
	id, _, err := middleware.ValidateToken(c)
	if err != nil {
		return err
	}

	playerId := c.Params("playerId")
	err = cr.applicationUserService.UpdatePlayerId(playerId, id, nil)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"successMessage": "PlayerID succuessfully updated",
	})
}

// AddMachineId
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param request body dto.PostMachineInfo true "AddMachineId Request"
// @Success 200
// @Router /user/add-machine-id [post]
func (cr *UserController) AddMachineId(c fiber.Ctx) error {
	id, _, err := middleware.ValidateToken(c)
	if err != nil {
		return err
	}

	data := new(dto.PostMachineInfo)
	if err := utils.BindNValidate(c, data); err != nil {
		return err
	}

	err = cr.applicationUserService.UpdateMachineId(data.MachineId, id, nil)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"successMessage": "MachineID successfully updated",
	})
}
