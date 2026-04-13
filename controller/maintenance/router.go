package maintenance

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    maintenanceController := NewMaintenanceController()
    maintenanceController.registerRoutes(router)
}

func (c *MaintenanceController) registerRoutes(router fiber.Router) {
    api := router.Group("/maintenance")

    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Get("/hospital-profile/all", c.GetAllHospitalProfiles)
    api.Post("/hospital-profile/search", c.SearchAllHospitalProfiles)
    api.Post("/hospital-profile/update", c.UpdateHospitalProfile)
    api.Get("/param-setting/all", c.GetAllParamSettings)
    api.Post("/param-setting/search", c.SearchAllParamSettings)
    api.Post("/param-setting/update", c.UpdateParamSetting)
    api.Get("/notification-setting/all", c.GetAllNotificationSettings)
    api.Post("/notification-setting/search", c.SearchAllNotificationSettings)
    api.Post("/notification-setting/update", c.UpdateNotificationSetting)
    api.Get("/cronjob-history/all", c.GetAllCronjobHistories)
    api.Post("/cronjob-history/search/:cronjobName", c.SearchAllCronjobHistories)
    api.Get("/dynamic-email-master/all", c.GetAllDynamicEmailSettings)
    api.Post("/dynamic-email-master/search", c.SearchAllDynamicEmailSettings)
    api.Get("/dynamic-email-master/export/all", c.GetAllExportDynamicEmailSettings)
    api.Post("/dynamic-email-master/export/search", c.GetSearchDynamicEmailSettings)
    api.Get("/dynamic-email-master/view/:emailFunctionName", c.GetDynamicEmailSettingByFunctionName)
    api.Post("/dynamic-email-master/update", c.UpdateDynamicEmailSetting)
    api.Get("/statistic-appointment/all", c.GetAllStatisticAppointments)
    api.Get("/statistic-mobile/registration/all", c.GetAllStatisticMobileRegistrations)
    api.Get("/statistic-mobile/feedback/all", c.GetAllStatisticMobileFeedbacks)
    api.Get("/statistic-mobile/package/all", c.GetAllStatisticMobilePackages)
    api.Get("/statistic-mobile/clubs/kids", c.GetAllStatisticMobileClubsKids)
    api.Get("/statistic-mobile/clubs/goldenpearl", c.GetAllStatisticMobileClubsGoldenPearl)
}
