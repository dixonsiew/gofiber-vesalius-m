package admin

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    adminController := NewAdminController()
    adminController.registerRoutes(router)
}

func (c *AdminController) registerRoutes(router fiber.Router) {
	api := router.Group("/admin")
    api.Post("/reset-signup-email/user/:email", c.ResetSignUpUserByEmail)
    api.Post("/self-reset-password/:branchId/:email", c.SelfResetPassword)
    api.Use(middleware.JWTProtected, middleware.ValidateAdminUser)
    api.Get("/", c.GetAdmin)
    api.Get("/all", c.GetAllAdmin)
    api.Post("/all", c.SearchAllAdmin)
    api.Get("/adminportal/mobile-user/log/all", c.GetAllAuditMobileUser)
    api.Post("/adminportal/mobile-user/log/all", c.SearchAllAuditMobileUser)
    api.Get("/adminportal/log/all", c.GetAllAuditLog)
    api.Post("/adminportal/log/all", c.SearchAllAuditLog)
    api.Get("/adminId/:adminId", c.GetUserById)
    api.Post("/reset-admin-password/:email", c.ResetAdminPassword)
    api.Post("/reset-user-password/:email", c.ResetUserPassword)
    api.Post("/delete-user/:userId", c.DeleteUser)
    api.Post("/link-user-prn", c.LinkUserPrn)
    api.Post("/change-password", c.ChangePassword)
    api.Post("/delete-admin/:email", c.DeleteAdmin)
    api.Post("/adminportal/save-log", c.SaveAdminPortalLog)
    api.Post("/change-user-password", c.ChangeUserPassword)
}
