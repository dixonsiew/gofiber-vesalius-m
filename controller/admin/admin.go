package admin

import (
    "fmt"
    "vesaliusm/database"
    "vesaliusm/middleware"
    adminUserService "vesaliusm/service/adminUser"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

var adminUserSvc *adminUserService.AdminUserService = adminUserService.NewAdminUserService(database.GetDb(), database.GetCtx())

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

    admin, err := adminUserSvc.FindWithAssignBranchByAdminId(user.AdminID.Int64)
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
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := adminUserSvc.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllAdmin
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Param        keyword           body        string  false  "Keyword"
// @Success 200 {array} model.AdminUser
// @Router /admin/all [post]
func SearchAllAdmin(c fiber.Ctx) error {
    var data fiber.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    var key string
    if keyword, ok := data["keyword"]; ok {
        key = keyword.(string)
        if key != "" {
            key = "%" + key + "%"
        }
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := adminUserSvc.ListByKeyword(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// GetAllAuditMobileUser
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.MobileUserAuditLog
// @Router /admin/adminportal/mobile-user/log/all [get]
func GetAllAuditMobileUser(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := adminUserSvc.ListMobileUserAuditLog(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllAuditMobileUser
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Param        keyword           body        string  false  "Keyword"
// @Param        keyword2          body        string  false  "Keyword2"
// @Success 200 {array} model.MobileUserAuditLog
// @Router /admin/adminportal/mobile-user/log/all [post]
func SearchAllAuditMobileUser(c fiber.Ctx) error {
    var data fiber.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    var (
        key  string
        key2 string
    )
    if keyword, ok := data["keyword"]; ok {
        key = keyword.(string)
        if key != "" {
            key = "%" + key + "%"
        }
    }
    if keyword, ok := data["keyword2"]; ok {
        key2 = keyword.(string)
        if key2 != "" {
            key2 = "%" + key2 + "%"
        }
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := adminUserSvc.ListMobileUserAuditLogByKeyword(key, key2, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}
