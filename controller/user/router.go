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
    api.Post("/set-new-password", c.SetNewPasswordUser)
    api.Post("/verify", c.VerifyUser)
    api.Post("/verify-email", c.VerifyUserEmail)
    api.Post("/verify-smstac", c.VerifyUserSmsTac)
    
    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Get("/all", c.GetAllUsers)
    api.Get("/notification/unseen/count", c.GetOldAppUnseenNotificationCount)
    api.Get("/notification/list", c.GetOldAppNotificationLists)
    api.Get("/all/active", c.GetAllActiveUsers)
    api.Post("/active-user/:userId", c.SetActive)
    api.Post("/inactive-user/:userId", c.SetInactive)
    api.Post("/all", c.SearchAllUsers)
    api.Get("/userId/:userId", c.GetUserById)
    api.Get("", c.GetUser)
    api.Get("/branches", c.GetUserBranches)
    api.Post("/change-password", c.ChangePassword)
    api.Post("/delete-account", c.UserDeleteAccount)
    api.Post("/disable-firsttime-bio", c.DisableFirstTimeBiometricUser)
    api.Post("/update-playerid/:playerId", c.UpdatePlayerId)
    api.Post("/add-machine-id", c.AddMachineId)
}
