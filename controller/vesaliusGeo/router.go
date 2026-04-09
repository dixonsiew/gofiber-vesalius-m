package vesaliusGeo

import (
    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    vesaliusGeoController := NewVesaliusGeoController()
    vesaliusGeoController.registerRoutes(router)
}

func (c *VesaliusGeoController) registerRoutes(router fiber.Router) {
    api := router.Group("/vesalius-geo")
    api.Get("/authenticate", c.Login)
    api.Get("/logout/:token", c.Logout)
}
