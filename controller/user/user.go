package user

import (
	"fmt"
	applicationuserService "vesaliusm/service/applicationUser"
	"vesaliusm/utils"

	"github.com/gofiber/fiber/v3"
)

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
    m, err := applicationuserService.List(page, limit)
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
    m, err := applicationuserService.ListActive(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}
