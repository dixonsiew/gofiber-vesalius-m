package userNotification

import (
    "vesaliusm/controller/userNotification"
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/notification")
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Get("/unseen/count", userNotification.GetUnseenNotificationCount)
    api.Get("/all", userNotification.GetNotificationList)
    api.Get("/general/master/all", userNotification.GetGeneralNotificationList)
}
