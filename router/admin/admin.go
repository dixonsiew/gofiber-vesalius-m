package admin

import (
    "vesaliusm/controller/admin"
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/admin")
    api.Use(middleware.JWTProtected, middleware.ValidateAdminUser)
    api.Get("/", admin.GetAdmin)
    api.Get("/all", admin.GetAllAdmin)
    api.Post("/all", admin.SearchAllAdmin)
    api.Get("/adminportal/mobile-user/log/all", admin.GetAllAuditMobileUser)
    api.Post("/adminportal/mobile-user/log/all", admin.SearchAllAuditMobileUser)
    api.Get("/adminportal/log/all", admin.GetAllAuditLog)
    api.Post("/adminportal/log/all", admin.SearchAllAuditLog)
    api.Get("/adminId/:adminId", admin.GetUserById)
    api.Post("/reset-admin-password/:email", admin.ResetAdminPassword)
    api.Post("/reset-user-password/:email", admin.ResetUserPassword)
    api.Post("/delete-user/:userId", admin.DeleteUser)
    api.Post("/link-user-prn", admin.LinkUserPrn)
    api.Post("/change-password", admin.ChangePassword)
}
