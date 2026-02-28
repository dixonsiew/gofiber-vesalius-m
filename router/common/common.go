package common

import (
    "vesaliusm/controller/common"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/common")
    api.Get("/telcode/list", common.GetCountriesTelCode)
    api.Get("/country/list", common.GetCountries)
    api.Get("/nationality/list", common.GetNationalities)
}
