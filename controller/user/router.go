package user

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    userController := NewUserController()
    userController.registerRoutes(router)
}

func (c *UserController) registerRoutes(router fiber.Router) {
    api := router.Group("/user")
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Get("/all", c.GetAllUsers)
    api.Get("/all/active", c.GetAllActiveUsers)
    api.Post("/update-playerid/:playerId", c.UpdatePlayerId)
    api.Post("/add-machine-id", c.AddMachineId)
}
