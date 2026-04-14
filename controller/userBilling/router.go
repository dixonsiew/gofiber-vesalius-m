package userBilling

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    userBillingController := NewUserBillingController()
    userBillingController.registerRoutes(router)
}

func (c *UserBillingController) registerRoutes(router fiber.Router) {
    api := router.Group("/user-billing")
    api.Post("/status/:billPaymentId", c.UpdateUserBillingPaymentStatus)

    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Post("/pay/:paymentMethod", c.CreateUserBillingDetails)
    api.Get("/paid/all/mobile", c.GetAllUserPaidBillingHistory)
}
