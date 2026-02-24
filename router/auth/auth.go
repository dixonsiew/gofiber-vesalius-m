package auth

import (
    "vesaliusm/controller/auth"

    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
    router.Post("/login", auth.Login)
}
