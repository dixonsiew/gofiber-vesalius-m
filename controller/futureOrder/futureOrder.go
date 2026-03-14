package futureOrder

import (
    "fmt"
    futureOrderService "vesaliusm/service/futureOrder"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

type FutureOrderController struct {
    futureOrderSvc *futureOrderService.FutureOrderService
}

func NewFutureOrderController(futureOrderSvc *futureOrderService.FutureOrderService) *FutureOrderController {
    return &FutureOrderController{
        futureOrderSvc: futureOrderSvc,
    }
}

// GetAllFutureOrder
//
// @Tags Future Order
// @Produce json
// @Security BearerAuth
// @Param        prn               path        string  true  "prn"
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.FutureOrder
// @Router /future-order/all/{prn} [get]
func (cr *FutureOrderController) GetAllFutureOrder(c fiber.Ctx) error {
    prn := c.Query("prn", "")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := cr.futureOrderSvc.List(prn, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}
