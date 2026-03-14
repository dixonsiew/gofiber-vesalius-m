package clubs

import (
    "vesaliusm/database"
    "vesaliusm/middleware"
    clubService "vesaliusm/service/clubs"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    var clubSvc *clubService.ClubService = clubService.NewClubService(database.GetDb(), database.GetCtx())

    clubController := NewClubsController(clubSvc)
    clubController.registerRoutes(router)
}

func (c *ClubsController) registerRoutes(router fiber.Router) {
    api := router.Group("/clubs")
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Get("/goldenpearl/about-us", c.GetGoldenPearlAboutUs)
}
