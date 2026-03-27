package userNotification

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    userNotificationController := NewUserNotificationController()
    userNotificationController.registerRoutes(router)
}

func (c *UserNotificationController) registerRoutes(router fiber.Router) {
    api := router.Group("/notification")
    
    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Get("/unseen/count", c.GetUnseenNotificationCount)
    api.Get("/all", c.GetNotificationList)
    api.Get("/:notificationId", c.GetNotificationById)
    api.Get("/general/master/all", c.GetGeneralNotificationList)
    api.Get("/general/master/:notificationMasterId", c.GetByNotificationMasterId)
}
