package user

import (
    "vesaliusm/database"
    "vesaliusm/middleware"
    applicationuserService "vesaliusm/service/applicationUser"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    var applicationUserSvc *applicationuserService.ApplicationUserService = 
        applicationuserService.NewApplicationUserService(database.GetDb(), database.GetCtx())

    userController := NewUserController(applicationUserSvc)
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
