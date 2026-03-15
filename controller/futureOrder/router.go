package futureOrder

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    futureOrderController := NewFutureOrderController()
    futureOrderController.registerRoutes(router)
}

func (c *FutureOrderController) registerRoutes(router fiber.Router) {
    api := router.Group("/future-order")
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Get("/all/:prn", c.GetAllFutureOrder)
}
