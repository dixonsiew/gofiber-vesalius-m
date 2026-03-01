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
}
