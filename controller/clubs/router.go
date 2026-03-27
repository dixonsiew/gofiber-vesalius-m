package clubs

import (
    "vesaliusm/controller/clubs/goldenPearl"
	"vesaliusm/controller/clubs/littleKids"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/clubs")
    littleKids.SetupRoutes(api)
    goldenPearl.SetupRoutes(api)
}
