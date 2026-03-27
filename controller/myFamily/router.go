package myFamily

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    myFamilyController := NewMyFamilyController()
    myFamilyController.registerRoutes(router)
}

func (c *MyFamilyController) registerRoutes(router fiber.Router) {
    api := router.Group("/my-family")

    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Get("", c.GetAllActiveUserFamilies)
    api.Get("/all/:userId", c.GetAllUserFamilies)
    api.Get("/familyId/:familyId", c.GetFamilyById)
}
