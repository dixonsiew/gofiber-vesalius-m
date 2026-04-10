package wallex

import (
    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    wallexController := NewWallexController()
    wallexController.registerRoutes(router)
}

func (c *WallexController) registerRoutes(router fiber.Router) {
    api := router.Group("/wallex")
    api.Get("/authenticate", c.Login)
    api.Post("/backend/backend_response", c.WallexWebhook)
}
