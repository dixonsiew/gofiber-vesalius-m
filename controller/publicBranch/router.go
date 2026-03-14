package publicBranch

import (
    "vesaliusm/database"
    branchService "vesaliusm/service/branch"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    var branchSvc *branchService.BranchService = 
        branchService.NewBranchService(database.GetDb(), database.GetCtx())

    publicBranchController := NewPublicBranchController(branchSvc)
    publicBranchController.registerRoutes(router)
}

func (c *PublicBranchController) registerRoutes(router fiber.Router) {
    api := router.Group("/public")
    api.Get("/branch/list", c.GetList)
}
