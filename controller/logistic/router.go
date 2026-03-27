package logistic

import (
	"vesaliusm/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
	logisticController := NewLogisticController()
	logisticController.registerRoutes(router)
}

func (c *LogisticController) registerRoutes(router fiber.Router) {
	api := router.Group("/logistic")

	authenticated := api.Group("/")
	authenticated.Use(middleware.JWTProtected, middleware.ValidateUser)
	authenticated.Get("/setup", c.GetLogisticSetup)

	app := api.Group("/")
	app.Use(middleware.JWTProtected, middleware.ValidateAppUser)
	app.Post("/slot/all/mobile", c.GetAllAppLogisticSlots)
	app.Post("/request", c.CreateLogisticRequest)
	app.Get("/request/all/mobile", c.GetAllAppLogisticRequests)
	app.Post("/request/status", c.UpdateAppLogisticRequestStatus)

	admin := api.Group("/")
	admin.Use(middleware.JWTProtected, middleware.ValidateAdminUser)
	admin.Post("/setup", c.CreateLogisticSetup)
	admin.Put("/setup/:logisticSetupId", c.UpdateLogisticSetup)
	admin.Post("/slot", c.CreateLogisticSlot)
	admin.Get("/slot/all", c.GetAllLogisticSlots)
	admin.Get("/request/all", c.GetAllLogisticRequests)
	admin.Post("/request/all", c.SearchAllLogisticRequests)
	admin.Get("/request/export/all", c.ExportAllLogisticRequests)
	admin.Post("/request/export/search", c.ExportSearchLogisticRequests)
	admin.Get("/request/:requestId", c.GetLogisticRequestByID)
	admin.Post("/request/status/webadmin", c.UpdateLogisticRequestStatus)
}
