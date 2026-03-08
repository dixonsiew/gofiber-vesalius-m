package publicBranch

import (
    "vesaliusm/controller/publicBranch"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/public")
    api.Get("/branch/list", publicBranch.GetList)
}
