package clubs

import (
    "vesaliusm/controller/clubs"
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/clubs")
    api.Use(middleware.JWTProtected, middleware.ValidateAppUser)
    api.Get("/goldenpearl/about-us", clubs.GetGoldenPearlAboutUs)
}
