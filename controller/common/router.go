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
    api.Get("/app/hospital-profile", c.GetAppHospitalProfile)
    api.Get("/app/version", c.GetAppVersion)
    api.Post("/app/version", c.UpdateAppVersion)
    api.Get("/app/release/version", c.GetReleaseVersion)
    api.Post("/app/release/version", c.UpdateReleaseVersion)
    api.Post("/request-delete-account", c.RequestDeleteAccount)
    api.Get("/service/guest/list", c.GetGuestModeServices)
    api.Get("/service/auth/list", c.GetAuthModeServices)
    api.Get("/telcode/list", c.GetCountriesTelCode)
    api.Get("/country/list", c.GetCountries)
    api.Get("/nationality/list", c.GetNationalities)
    api.Post("/downloaded-app/:playerId", c.UserDownloadedApp)
    api.Post("/downloaded-app/v2", c.UserDownloadedAppV2)
}
