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
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Get("/goldenpearl/about-us", c.GetGoldenPearlAboutUs)
}
