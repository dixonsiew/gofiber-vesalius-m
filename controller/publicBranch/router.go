package publicBranch

import (
    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    publicBranchController := NewPublicBranchController()
    publicBranchController.registerRoutes(router)
}

func (c *PublicBranchController) registerRoutes(router fiber.Router) {
    api := router.Group("/public")
    api.Get("/branch/list", c.GetList)
}
