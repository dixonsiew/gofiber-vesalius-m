package guest

import (
    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    guestController := NewGuestController()
    guestController.registerRoutes(router)
}

func (c *GuestController) registerRoutes(router fiber.Router) {
    api := router.Group("/guest")
    api.Get("/vesalius/getAllDoctorInformation/:branchId", c.GetAllDoctorInformation)
    api.Post("/vesalius/getAllDoctorInformation/:branchId", c.SearchAllDoctorInformation)
    api.Get("/notification/all/:playerId", c.getAllGuestNotificationLists)
    api.Get("/notification/seen/:notificationId/:playerId", c.GetGuestUnseenNotificationCount)
    api.Get("/clubs/littlekids/about-us", c.GetLittleKidsAboutUs)
}
