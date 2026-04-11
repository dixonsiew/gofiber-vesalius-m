package qms

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    qmsController := NewQmsController()
    qmsController.registerRoutes(router)
}

func (c *QmsController) registerRoutes(router fiber.Router) {
    api := router.Group("/qms")
    api.Post("/backend/qms_response", c.QmsServerWebhook)
    
    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Post("/backend/qms_request", c.QmsClientWebhook)
}
