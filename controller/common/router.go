package common

import (
    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    commonController := NewCommonController()
    commonController.registerRoutes(router)
}

func (c *CommonController) registerRoutes(router fiber.Router) {
    api := router.Group("/common")
    api.Get("/telcode/list", c.GetCountriesTelCode)
    api.Get("/country/list", c.GetCountries)
    api.Get("/nationality/list", c.GetNationalities)
}
