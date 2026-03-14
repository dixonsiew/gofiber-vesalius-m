package userPackage

import (
    "vesaliusm/database"
    "vesaliusm/middleware"
    patientPurchaseDetailsService "vesaliusm/service/patientPurchaseDetails"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    var patientPurchaseDetailsSvc *patientPurchaseDetailsService.PatientPurchaseDetailsService = 
        patientPurchaseDetailsService.NewPatientPurchaseDetailsService(database.GetDb(), database.GetCtx())

    userPackageController := NewUserPackageController(patientPurchaseDetailsSvc)
    userPackageController.registerRoutes(router)
}

func (c *UserPackageController) registerRoutes(router fiber.Router) {
    api := router.Group("/user-package")
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Post("/check/expiry-maxpurchase", c.CheckPackageExpiryMaxpurchase)
    api.Get("/all/mobile", c.GetAllUserPurchaseHistory)
    api.Get("/all", c.GetAllPurchaseHistory)
    api.Post("/all", c.SearchAllPurchaseHistory)
    api.Get("/:purchaseId", c.GetUserPackageById)
}
