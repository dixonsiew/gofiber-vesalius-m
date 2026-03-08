package publicBranch

import (
    branchService "vesaliusm/service/branch"
    "vesaliusm/database"

    "github.com/gofiber/fiber/v3"
)

var branchSvc *branchService.BranchService = branchService.NewBranchService(database.GetDb(), database.GetCtx())

// GetList
//
// @Tags Public Branch
// @Produce json
// @Success 200 {array} model.Branch
// @Router /public/branch/list [get]
func GetList(c fiber.Ctx) error {
    lx, err := branchSvc.FindAll()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}
