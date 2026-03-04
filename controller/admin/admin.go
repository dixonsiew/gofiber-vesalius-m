package admin

import (
    "fmt"
    "vesaliusm/middleware"
    adminUserService "vesaliusm/service/adminUser"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

// GetAdmin
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.AdminUser
// @Router /admin [get]
func GetAdmin(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    admin, err := adminUserService.FindWithAssignBranchByAdminId(user.AdminID.Int64)
    if err != nil {
        return err
    }

    return c.JSON(admin)
}

// GetAllAdmin
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.AdminUser
// @Router /admin/all [get]
func GetAllAdmin(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return fiber.NewError(fiber.StatusUnauthorized, "User not found")
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := adminUserService.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}
