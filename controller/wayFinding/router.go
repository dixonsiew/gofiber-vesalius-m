package wayFinding

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
    api.Get("/floors", c.GetAllFloors)
    api.Get("/floors/webadmin", c.GetAllFloorsWebAdmin)
    api.Post("/floors", c.SearchAllFloors)
    api.Post("/floors/webadmin", c.SearchAllFloorsWebAdmin)
    api.Get("/locations", c.GetAllLocations)
    api.Post("/locations", c.SearchAllLocations)
    api.Get("/location-types", c.GetAllLocationTypes)
    api.Get("/routes", c.GetAllRoutes)
    api.Post("/routes", c.SearchAllRoutes)
    api.Post("/location/:code", c.SearchAllLocationsByCode)
    api.Get("/location/:code", c.GetAllLocationsByTypeCode)
    api.Get("/route/:fromId/:toId", c.GetRoute)
    api.Get("/location-qr/:locationId/:locationTypeId", c.GetLocationsByLocationIdAndLocationTypeId)
    
    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Post("/dropdown", c.GetDropdowns)
    api.Post("/buildings/create", c.CreateBuilding)
    api.Put("/buildings/update/:buildingCode", c.UpdateBuilding)
    api.Delete("/buildings/delete/:buildingCode", c.DeleteBuildingsByBuildingCode)
    api.Get("/buildings/:buildingCode", c.GetBuildingsByBuildingCode)
    api.Post("/floors/create", c.CreateFloor)
    api.Put("/floors/update/:floorCode", c.UpdateFloor)
    api.Delete("/floors/delete/:floorCode", c.DeleteFloorsByFloorCode)
    api.Get("/floors/:floorCode", c.GetFloorsByFloorCode)
    api.Post("/locations/create", c.CreateLocation)
    api.Put("/locations/update/:locationId", c.UpdateLocation)
    api.Delete("/locations/delete/:locationId", c.DeleteLocationsById)
    api.Get("/locations/:locationId", c.GetLocationsById)
    api.Post("/location-types/create", c.CreateLocationType)
    api.Put("/location-types/update/:locationTypeCode", c.UpdateLocationType)
    api.Delete("/location-types/delete/:locationTypeCode", c.DeleteLocationTypesByLocationTypeCode)
    api.Get("/location-types/:locationTypeCode", c.GetLocationTypesByLocationTypeCode)
    api.Post("/routes/create", c.CreateRoute)
    api.Put("/routes/update/:routeId", c.UpdateRoute)
    api.Delete("/routes/delete/:routeId", c.DeleteRoutesByRouteId)
    api.Get("/routes/:routeId", c.GetRoutesByRouteId)
}
