package hpackage

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    packageController := NewPackageController()
    packageController.registerRoutes(router)
}

func (c *PackageController) registerRoutes(router fiber.Router) {
    api := router.Group("/package")
    api.Post("/process-resize-image", c.ProcessResizeImage)
    api.Get("/all", c.GetAllPackages)
    api.Get("/all/mobile/:isHome", c.GetAllAppPackages)

    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Post("/all", c.SearchAllPackages)
    api.Get("/status/:packageId", c.GetPackageStatusById)
    api.Get("/:packageId", c.GetPackageById)
    api.Post("/", c.CreatePackage)
    api.Put("/:packageId", c.UpdatePackage)
}
