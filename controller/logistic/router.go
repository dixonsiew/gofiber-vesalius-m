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
	
	api.Use(middleware.JWTProtected, middleware.ValidateUser)
	api.Post("/setup", c.CreateLogisticSetup)
	api.Put("/setup/:logisticSetupId", c.UpdateLogisticSetup)
	api.Get("/setup", c.GetLogisticSetup)
	api.Post("/slot", c.CreateLogisticSlot)
	api.Post("/slot/all/mobile", c.GetAllAppLogisticSlots)
	api.Get("/slot/all", c.GetAllLogisticSlots)
	api.Post("/request", c.CreateLogisticRequest)
	api.Get("/request/all/mobile", c.GetAllAppLogisticRequests)
	api.Get("/request/all", c.GetAllLogisticRequests)
	api.Post("/request/all", c.SearchAllLogisticRequests)
	api.Get("/request/export/all", c.GetAllExportLogisticRequest)
	api.Post("/request/export/search", c.GetSearchExportLogisticRequest)
	api.Get("/request/:requestId", c.GetLogisticRequestById)
	api.Post("/request/status", c.UpdateAppLogisticRequestStatus)
	api.Post("/request/status/webadmin", c.UpdateLogisticRequestStatus)
}
