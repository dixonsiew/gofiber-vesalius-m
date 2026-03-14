package userNotification

import (
    "vesaliusm/database"
    "vesaliusm/middleware"
    applicationUserNotificationService "vesaliusm/service/applicationUserNotification"
    generalNotificationMasterService "vesaliusm/service/generalNotificationMaster"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    var (
        applicationUserNotificationSvc *applicationUserNotificationService.ApplicationUserNotificationService = 
            applicationUserNotificationService.NewApplicationUserNotificationService(database.GetDb(), database.GetCtx())
        generalNotificationMasterSvc *generalNotificationMasterService.GeneralNotificationMasterService = 
            generalNotificationMasterService.NewGeneralNotificationMasterService(database.GetDb(), database.GetCtx())
    )

    userNotificationController := NewUserNotificationController(applicationUserNotificationSvc, generalNotificationMasterSvc)
    userNotificationController.registerRoutes(router)
}

func (c *UserNotificationController) registerRoutes(router fiber.Router) {
    api := router.Group("/notification")
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Get("/unseen/count", c.GetUnseenNotificationCount)
    api.Get("/all", c.GetNotificationList)
    api.Get("/:notificationId", c.GetNotificationById)
    api.Get("/general/master/all", c.GetGeneralNotificationList)
    api.Get("/general/master/:notificationMasterId", c.GetByNotificationMasterId)
}
