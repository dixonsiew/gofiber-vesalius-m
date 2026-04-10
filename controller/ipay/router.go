package ipay

import (
    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    ipayController := NewIpayController()
    ipayController.registerRoutes(router)
}

func (c *IpayController) registerRoutes(router fiber.Router) {
    api := router.Group("/payment/ipay88")
    api.Get("/requery", c.Requery)
    api.Get("/submit", c.Submit)
    api.Post("/backend/response", c.BackendResponse)
    api.Post("/backend/backend_response", c.BackendPostResponse)
}
