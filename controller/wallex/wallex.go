package wallex

import (
	"strings"
	"vesaliusm/dto"
	"vesaliusm/service/billPaymentDetails"
	"vesaliusm/service/packagePaymentDetails"
	"vesaliusm/service/wallex"
	"vesaliusm/utils"

	"github.com/gofiber/fiber/v3"
)

type WallexController struct {
    billPaymentDetailsService    *billPaymentDetails.BillPaymentDetailsService
    packagePaymentDetailsService *packagePaymentDetails.PackagePaymentDetailsService
    wallexService                *wallex.WallexService
}

func NewWallexController() *WallexController {
    return &WallexController{
        billPaymentDetailsService:    billPaymentDetails.BillPaymentDetailsSvc,
        packagePaymentDetailsService: packagePaymentDetails.PackagePaymentDetailsSvc,
        wallexService:                wallex.WallexSvc,
    }
}

// Login
//
// @Tags Wallex
// @Produce json
// @Success 200
// @Router /wallex/authenticate [get]
func (cr *WallexController) Login(c fiber.Ctx) error {
    token, err := cr.wallexService.Authenticate()
    if err != nil {
        return err
    }
    return c.JSON(fiber.Map{
        "token": token,
    })
}

// WallexWebhook
//
// @Tags Wallex
// @Produce json
// @Param        request       body        dto.WallexWebhookDto  true  "WallexWebhookDto"
// @Success 200
// @Router /wallex/backend/backend_response [post]
func (cr *WallexController) WallexWebhook(c fiber.Ctx) error {
    data := new(dto.WallexWebhookDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    if data.Resource == "collection_request" && data.Status == "full-paid" {
        if strings.Contains(data.Remarks, "MB") {
            err := cr.billPaymentDetailsService.UpdateWallexPaymentStatus(data.ResourceId)
            if err != nil {
                return err
            }
        } else if strings.Contains(data.Remarks, "MP") {
            err := cr.packagePaymentDetailsService.UpdateWallexPaymentStatus(data.ResourceId)
            if err != nil {
                return err
            }
        } else {
            return fiber.NewError(fiber.StatusBadRequest, "Invalid Request Number")
        }

        return c.JSON(fiber.Map{
            "message": "Wallex Payment Details Updated",
        })
    } else {
        return c.JSON(fiber.Map{
            "message": "No Payment Wallex Details to update",
        })
    }
}
