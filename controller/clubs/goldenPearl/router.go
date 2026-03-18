package goldenPearl

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    clubsGoldenPearlController := NewClubsGoldenPearlController()
    clubsGoldenPearlController.registerRoutes(router)
}

func (c *ClubsGoldenPearlController) registerRoutes(router fiber.Router) {
    api := router.Group("/goldenpearl")
    api.Post("/activity/participate", c.ParticipateGoldenPearlActivity)
    api.Get("/activity/all", c.GetAllGoldenPearlActivities)
    api.Get("/activity/all/mobile/:isHome", c.GetAllAppGoldenPearlActivities)
    api.Get("/activity/:activityId", c.GetGoldenPearlActivityById)
    api.Get("/about-us", c.GetGoldenPearlAboutUs)

    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Post("/membership", c.CreateGoldenPearlMembership)
    api.Post("/membership/webadmin", c.CreateGoldenPearlMembershipViaWebportal)
    api.Put("/membership/webadmin/:membershipId", c.UpdateGoldenPearlMembership)
    api.Get("/membership/:membershipId", c.GetGoldenPearlMembershipById)
    api.Get("/membership/all/mobile", c.GetAllAppGoldenPearlMemberships)
    api.Get("/membership/export/all", c.GetAllExportGoldenPearlMembership)
    api.Post("/membership/export/search", c.GetSearchExportGoldenPearlMembership)
    api.Get("/membership/all", c.GetAllGoldenPearlMemberships)
    api.Post("/membership/all", c.SearchAllGoldenPearlMembership)
    api.Get("/my-activity/all", c.GetAllUserGoldenPearlActivities)
    api.Post("/activity", c.CreateGoldenPearlActivity)
    api.Put("/activity/:activityId", c.UpdateGoldenPearlActivity)
    api.Get("/activity/export/all", c.GetAllExportGoldenPearlActivity)
    api.Post("/activity/all", c.SearchAllGoldenPearlActivities)
    api.Get("/activity/attendees/:activityId/export/all", c.GetAllExportGoldenPearlAttendees)
    api.Post("/activity/attendees/:activityId/export/search", c.GetSearchExportGoldenPearlAttendees)
    api.Get("/activity/attendees/:activityId", c.GetGoldenPearlActivityAttendeesById)
    api.Post("/activity/attendees/:activityId", c.SearchAllGoldenPearlAttendees)
    api.Get("/activity/name/:activityId", c.GetGoldenPearlActivityNameById)
    api.Post("/about-us", c.CreateGoldenPearlAboutUs)
    api.Put("/about-us/:goldenPearlId", c.UpdateGoldenPearlAboutUs)
}
