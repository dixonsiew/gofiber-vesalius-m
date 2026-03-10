package userPackage

import (
    "vesaliusm/controller/userPackage"
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/user-package")
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Post("/check/expiry-maxpurchase", userPackage.CheckPackageExpiryMaxpurchase)
    api.Get("/all/mobile", userPackage.GetAllUserPurchaseHistory)
    api.Get("/all", userPackage.GetAllPurchaseHistory)
    api.Post("/all", userPackage.SearchAllPurchaseHistory)
    api.Get("/:purchaseId", userPackage.GetUserPackageById)
}
