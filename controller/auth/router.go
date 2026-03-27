package auth

import (
    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    authController := NewAuthController()
    authController.registerRoutes(router)
}

func (c *AuthController) registerRoutes(router fiber.Router) {
    router.Post("/login", c.Login)
    router.Post("/login/v2", c.NewLogin)
}
