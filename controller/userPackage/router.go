package userPackage

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    userPackageController := NewUserPackageController()
    userPackageController.registerRoutes(router)
}

func (c *UserPackageController) registerRoutes(router fiber.Router) {
    api := router.Group("/user-package")
    
    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Post("/check/expiry-maxpurchase", c.CheckPackageExpiryMaxpurchase)
    api.Get("/all/mobile", c.GetAllUserPurchaseHistory)
    api.Get("/all", c.GetAllPurchaseHistory)
    api.Post("/all", c.SearchAllPurchaseHistory)
    api.Get("/:purchaseId", c.GetUserPackageById)
    api.Post("/status/:purchaseId", c.UpdateUserPackageStatus)
}
