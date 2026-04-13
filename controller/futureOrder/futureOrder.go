package futureOrder

import (
	"strconv"
	"vesaliusm/service/futureOrder"
	"vesaliusm/utils/constants"

	"github.com/gofiber/fiber/v3"
)

type FutureOrderController struct {
	futureOrderService *futureOrder.FutureOrderService
}

func NewFutureOrderController() *FutureOrderController {
	return &FutureOrderController{
		futureOrderService: futureOrder.FutureOrderSvc,
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
// @Success 200 {array} futureOrder.FutureOrder
// @Router /future-order/all/{prn} [get]
func (cr *FutureOrderController) GetAllFutureOrder(c fiber.Ctx) error {
	prn := c.Query("prn", "")
	page := c.Query("_page", "1")
	limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
	m, err := cr.futureOrderService.List(prn, page, limit)
	if err != nil {
		return err
	}

	c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
	c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
	return c.JSON(m.List)
}
