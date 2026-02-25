package user

import (
    "vesaliusm/controller/user"
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/user")
    api.Use(middleware.JWTProtected)
    api.Get("/all", user.GetAllUsers)
}
