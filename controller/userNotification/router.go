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
    api.Post("/send-notification/:playerId", c.SendNotification)
    
    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Post("/seen/:notificationId", c.SetNotificationSeen)
    api.Get("/unseen/count", c.GetUnseenNotificationCount)
    api.Get("/all", c.GetNotificationList)
    api.Post("/general/master", c.CreateGeneralNotification)
    api.Put("/general/master/:notificationMasterId", c.UpdateGeneralNotification)
    api.Get("/general/master/all", c.GetGeneralNotificationList)
    api.Get("/general/master/:notificationMasterId", c.GetByNotificationMasterId)
    api.Get("/:notificationId", c.GetNotificationById)
}
