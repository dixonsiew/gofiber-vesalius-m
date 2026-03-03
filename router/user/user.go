package user

import (
    "vesaliusm/controller/user"
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/user")
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Get("/all", user.GetAllUsers)
    api.Get("/all/active", user.GetAllActiveUsers)
    api.Post("/update-playerid/:playerId", user.UpdatePlayerId)
    api.Post("/add-machine-id", user.AddMachineId)
}
