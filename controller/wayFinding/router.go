package wayfinding

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    wayFindingController := NewWayFindingController()
    wayFindingController.registerRoutes(router)
}

func (c *WayFindingController) registerRoutes(router fiber.Router) {
    api := router.Group("/way-finding")
    api.Get("/buildings", c.GetAllBuildings)
    
    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Post("/dropdown", c.GetDropdowns)
    api.Post("/buildings/create", c.CreateBuilding)
    api.Put("/buildings/update/:buildingCode", c.UpdateBuilding)
    api.Delete("/buildings/delete/:buildingCode", c.DeleteBuildingsByBuildingCode)
    api.Get("/buildings/:buildingCode", c.GetBuildingsByBuildingCode)
    api.Post("/floors/create", c.CreateFloor)
}
