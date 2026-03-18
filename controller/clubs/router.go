package clubs

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/clubs")
    clubsLittleKidsController := NewClubsLittleKidsController()
    clubsLittleKidsController.registerRoutes(api)

    clubsGoldenPearlController := NewClubsGoldenPearlController()
    clubsGoldenPearlController.registerRoutes(api)
}

func (c *ClubsLittleKidsController) registerRoutes(router fiber.Router) {
    api := router.Group("/littlekids")
    api.Post("/activity/participate", c.ParticipateLittleKidsActivity)
    api.Get("/activity/all", c.GetAllLittleKidsActivities)
    api.Get("/activity/all/mobile/:isHome", c.GetAllAppLittleKidsActivities)
    api.Get("/activity/:activityId", c.GetLittleKidsActivityById)
    api.Get("/about-us", c.GetLittleKidsAboutUs)

    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Post("/membership", c.CreateLittleKidsMembership)
    api.Post("/membership/webadmin", c.CreateLittleKidsMembershipViaWebportal)
    api.Put("/membership/webadmin/:membershipId", c.UpdateLittleKidsMembership)
    api.Get("/membership/:membershipId", c.GetLittleKidsMembershipById)
    api.Get("/membership/all/mobile", c.GetAllAppLittleKidsMemberships)
    api.Get("/membership/export/all", c.GetAllExportLittleKidsMembership)
    api.Post("/membership/export/search", c.GetSearchExportLittleKidsMembership)
    api.Get("/membership/all", c.GetAllLittleKidsMemberships)
    api.Post("/membership/all", c.SearchAllLittleKidsMembership)
    api.Get("/my-activity/all", c.GetAllUserLittleKidsActivities)
    api.Post("/activity", c.CreateLittleKidsActivity)
    api.Get("/activity/export/all", c.GetAllExportLittleKidsActivity)
    api.Post("/activity/export/search", c.GetSearchExportLittleKidsActivity)
    api.Post("/activity/all", c.SearchAllLittleKidsActivities)
    api.Get("/activity/attendees/:activityId/export/all", c.GetAllExportLittleKidsAttendees)
    api.Post("/activity/attendees/:activityId/export/search", c.GetSearchExportLittleKidsAttendees)
    api.Get("/activity/attendees/:activityId", c.GetLittleKidsActivityAttendeesById)
    api.Post("/activity/attendees/:activityId", c.SearchAllLittleKidsAttendees)
    api.Get("/activity/name/:activityId", c.GetLittleKidsActivityNameById)
    api.Put("/activity/:activityId", c.UpdateLittleKidsActivity)
    api.Post("/about-us", c.CreateLittleKidsAboutUs)
    api.Put("/about-us/:kidsClubId", c.UpdateLittleKidsAboutUs)
}

func (c *ClubsGoldenPearlController) registerRoutes(router fiber.Router) {
    api := router.Group("/goldenpearl")
    api.Post("/activity/participate", c.ParticipateGoldenPearlActivity)
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
}
