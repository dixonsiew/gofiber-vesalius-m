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
    api.Post("/appointment/returning-patient", c.GetGuestReturningPatient)
    api.Post("/appointment/new-patient", c.MakeGuestNewPatient)
    api.Get("/appointment/get-doctor-appointments/:doctorId/:month/:year/:needAppt", c.GetDoctorAppointments)
    api.Post("/appointment/check-make-appointment/:branchId/:prn", c.CheckGuestPatientAppointment)
    api.Post("/appointment/make-appointment/:branchId/:prn", c.GetMakeGuestAppointment)
    api.Get("/notification/all/:playerId", c.GetAllGuestNotificationLists)
    api.Get("/notification/unseen/count/:playerId", c.GetGuestUnseenNotificationCount)
    api.Get("/notification/seen/:notificationId/:playerId", c.SetGuestNotificationSeen)
    api.Get("/clubs/littlekids/about-us", c.GetLittleKidsAboutUs)
    api.Post("/clubs/littlekids/membership", c.CreateGuestLittleKidsMembership)
    api.Get("/clubs/littlekids/membership/:identificationNumber", c.GetAllAppLittleKidsMemberships)
    api.Get("/clubs/littlekids/activity/all/mobile/:isHome", c.GetAllAppLittleKidsActivities)
    api.Post("/clubs/littlekids/activity/participate", c.ParticipateLittleKidsActivity)
    api.Get("/clubs/goldenpearl/about-us", c.GetGoldenPearlAboutUs)
    api.Post("/clubs/goldenpearl/membership", c.CreateGuestGoldenPearlMembership)
    api.Get("/clubs/goldenpearl/membership/:identificationNumber", c.GetAllAppGoldenPearlMemberships)
    api.Get("/clubs/goldenpearl/activity/all/mobile/:isHome", c.GetAllAppGoldenPearlActivities)
    api.Post("/clubs/goldenpearl/activity/participate", c.ParticipateGoldenPearlActivity)
    api.Get("/package/all/mobile/:isHome", c.GetAllAppPackages)
    api.Get("/package/packageStatus/:packageId", c.GetPackageStatusById)
    api.Get("/package/:packageId", c.GetPackageById)
    api.Post("/purchase/:paymentMethod", c.CreateGuestPurchaseDetails)
    api.Post("/package/check/expiry-maxpurchase", c.CheckPackageExpiryMaxpurchase)
}
