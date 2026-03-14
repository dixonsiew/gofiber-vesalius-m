package futureOrder

import (
    "vesaliusm/database"
    "vesaliusm/middleware"
    futureOrderService "vesaliusm/service/futureOrder"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    var futureOrderSvc *futureOrderService.FutureOrderService = futureOrderService.NewFutureOrderService(database.GetDb(), database.GetCtx())

    futureOrderController := NewFutureOrderController(futureOrderSvc)
    futureOrderController.registerRoutes(router)
}

func (c *FutureOrderController) registerRoutes(router fiber.Router) {
    api := router.Group("/future-order")
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Get("/all/:prn", c.GetAllFutureOrder)
}
