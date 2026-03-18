package clubs

import (
    "vesaliusm/controller/clubs/goldenpearl"
	"vesaliusm/controller/clubs/littlekids"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/clubs")
    littleKids.SetupRoutes(api)
    goldenPearl.SetupRoutes(api)
}
