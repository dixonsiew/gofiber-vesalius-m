package user

import (
    "fmt"
    "vesaliusm/database"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    applicationuserService "vesaliusm/service/applicationUser"
    "vesaliusm/utils"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v3"
)

var applicationUserSvc *applicationuserService.ApplicationUserService = 
    applicationuserService.NewApplicationUserService(database.GetDb(), database.GetCtx())

// GetAllUsers
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
// @Success 200 {array} model.ApplicationUser
// @Router /user/all [get]
func GetAllUsers(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := applicationUserSvc.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
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
func GetAllActiveUsers(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := applicationUserSvc.ListActive(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
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
func UpdatePlayerId(c fiber.Ctx) error {
    id, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    playerId := c.Params("playerId")
    err = applicationUserSvc.UpdatePlayerId(playerId, id, nil)
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
func AddMachineId(c fiber.Ctx) error {
    id, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    data := new(dto.PostMachineInfo)
    if err := c.Bind().Body(data); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            errs := utils.GetValidationErrors(validationErrors)
            if errs != nil {
                return errs
            }
        }

        return err
    }

    err = applicationUserSvc.UpdateMachineId(data.MachineId, id, nil)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "MachineID successfully updated",
    })
}
