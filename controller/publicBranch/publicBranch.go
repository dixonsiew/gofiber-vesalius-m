package publicBranch

import (
    "vesaliusm/service/branch"

    "github.com/gofiber/fiber/v3"
)

type PublicBranchController struct {
    branchService *branch.BranchService
}

func NewPublicBranchController() *PublicBranchController {
    return &PublicBranchController{
        branchService: branch.BranchSvc,
    }
}

// GetList
//
// @Tags Public Branch
// @Produce json
// @Success 200 {array} model.Branch
// @Router /public/branch/list [get]
func (cr *PublicBranchController) GetList(c fiber.Ctx) error {
    lx, err := cr.branchService.FindAll()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}
