package publicBranch

import (
    branchService "vesaliusm/service/branch"

    "github.com/gofiber/fiber/v3"
)

type PublicBranchController struct {
    branchSvc *branchService.BranchService
}

func NewPublicBranchController(branchSvc *branchService.BranchService) *PublicBranchController {
    return &PublicBranchController{
        branchSvc: branchSvc,
    }
}

// GetList
//
// @Tags Public Branch
// @Produce json
// @Success 200 {array} model.Branch
// @Router /public/branch/list [get]
func (cr *PublicBranchController) GetList(c fiber.Ctx) error {
    lx, err := cr.branchSvc.FindAll()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}
