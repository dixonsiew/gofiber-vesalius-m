package clubs

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    clubController := NewClubsController()
    clubController.registerRoutes(router)
}

func (c *ClubsController) registerRoutes(router fiber.Router) {
    api := router.Group("/clubs")
    api.Post("/littlekids/activity/participate", c.ParticipateLittleKidsActivity)
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Post("/littlekids/membership", c.CreateLittleKidsMembership)
    api.Post("/littlekids/membership/webadmin", c.CreateLittleKidsMembershipViaWebportal)
    api.Get("/littlekids/membership/:membershipId", c.GetLittleKidsMembershipById)
    api.Get("/littlekids/membership/all/mobile", c.GetAllAppLittleKidsMemberships)
    api.Get("/littlekids/membership/export/all", c.GetAllExportLittleKidsMembership)
    api.Post("/littlekids/membership/export/search", c.GetSearchExportLittleKidsMembership)
    api.Get("/littlekids/membership/all", c.GetAllLittleKidsMemberships)
    api.Post("/littlekids/membership/all", c.SearchAllLittleKidsMembership)
    api.Get("/littlekids/my-activity/all", c.GetAllUserLittleKidsActivities)
    api.Get("/goldenpearl/about-us", c.GetGoldenPearlAboutUs)
}
