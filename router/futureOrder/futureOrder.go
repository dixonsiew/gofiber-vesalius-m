package futureOrder

import (
    "vesaliusm/controller/futureOrder"
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/future-order")
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Get("/all/:prn", futureOrder.GetAllFutureOrder)
}
