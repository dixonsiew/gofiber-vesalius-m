package ipay

import (
    "fmt"
    "strings"
    "vesaliusm/dto"
    "vesaliusm/service/billPaymentDetails"
    "vesaliusm/service/ipay"
    "vesaliusm/service/packagePaymentDetails"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

type IpayController struct {
    billPaymentDetailsService    *billPaymentDetails.BillPaymentDetailsService
    packagePaymentDetailsService *packagePaymentDetails.PackagePaymentDetailsService
    ipayService                  *ipay.IpayService
}

func NewIpayController() *IpayController {
    return &IpayController{
        billPaymentDetailsService:    billPaymentDetails.BillPaymentDetailsSvc,
        packagePaymentDetailsService: packagePaymentDetails.PackagePaymentDetailsSvc,
        ipayService:                  ipay.IpaySvc,
    }
}

// Requery
//
// @Tags IPAY
// @Produce json
// @Param        refno             query       string  false  "refno"  default:""
// @Param        amount            query       string  false  "amount" default:""
// @Success 200
// @Router /payment/ipay88/requery [get]
func (cr *IpayController) Requery(c fiber.Ctx) error {
    refno := c.Query("refno")
    amount := c.Query("amount")
    s, err := cr.ipayService.Requery(refno, amount)
    if err != nil {
        return err
    }

    return c.SendString(s)
}

// Submit
//
// @Tags IPAY
// @Produce json
// @Param       requestNo       query       string  false  "requestNo"  default:""
// @Success 200
// @Router /payment/ipay88/submit [get]
func (cr *IpayController) Submit(c fiber.Ctx) error {
    requestNo := c.Query("requestNo")
    amount := ""
    m := fiber.Map{
        "MerchantCode":  cr.ipayService.MerchantCode,
        "PaymentId":     "",
        "Lang":          "UTF-8",
        "SignatureType": "HMACSHA512",
    }
    if requestNo != "" {
        if strings.Contains(requestNo, "MB") {
            paymentDetails, err := cr.billPaymentDetailsService.FindByRequestNo(requestNo)
            if err != nil {
                return err
            }

            if cr.ipayService.TestEnv == "Y" {
                amount = "1"
            } else {
                amount = utils.GetAmount(paymentDetails.PaymentAmount.Float64)
            }

            m["RefNo"] = paymentDetails.PaymentRequestNo
            m["Amount"] = amount
            m["Currency"] = paymentDetails.PaymentCurrency
            m["ProdDesc"] = paymentDetails.PaymentPurpose
            m["UserName"] = paymentDetails.BillingFullname
            m["UserEmail"] = paymentDetails.BillingEmail
            m["UserContact"] = fmt.Sprintf("%s%s", paymentDetails.BillingContactCode, paymentDetails.BillingContactNo)
            m["Remark"] = paymentDetails.PaymentRequestNo
            m["Signature"] = cr.ipayService.BuildSignature(requestNo, amount, "MYR")

        } else if strings.Contains(requestNo, "MP") {
            paymentDetails, err := cr.packagePaymentDetailsService.FindByRequestNo(requestNo)
            if err != nil {
                return err
            }

            if cr.ipayService.TestEnv == "Y" {
                amount = "1"
            } else {
                amount = utils.GetAmount(paymentDetails.PaymentAmount.Float64)
            }

            m["RefNo"] = paymentDetails.PaymentRequestNo
            m["Amount"] = amount
            m["Currency"] = paymentDetails.PaymentCurrency
            m["ProdDesc"] = paymentDetails.PaymentPurpose
            m["UserName"] = paymentDetails.BillingFullname
            m["UserEmail"] = paymentDetails.BillingEmail
            m["UserContact"] = fmt.Sprintf("%s%s", paymentDetails.BillingContactCode, paymentDetails.BillingContactNo)
            m["Remark"] = paymentDetails.PaymentRequestNo
            m["Signature"] = cr.ipayService.BuildSignature(requestNo, amount, "MYR")

        } else {
            return fiber.NewError(fiber.StatusBadRequest, "Invalid Request Number")
        }
    }

    return c.Render("submit", m)
}

// BackendResponse
//
// @Tags IPAY
// @Produce json
// @Param        request       body        dto.IPayPaymentResponseDto  true  "IPayPaymentResponseDto"
// @Success 200
// @Router /payment/ipay88/backend/response [post]
func (cr *IpayController) BackendResponse(c fiber.Ctx) error {
    data := new(dto.IPayPaymentResponseDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    c.Set(fiber.HeaderContentType, "text/plain")
    qs := cr.ipayService.BuildResponseSignature(data.PaymentId, data.RefNo, data.Amount, data.Currency, data.Status)
    r := "Payment fail."
    if data.Status == "1" && qs == data.Signature {
        if strings.Contains(data.RefNo, "MB") {
            _ = cr.billPaymentDetailsService.UpdateIPayPaymentStatus(data.RefNo)
        } else if strings.Contains(data.RefNo, "MP") {
            _ = cr.packagePaymentDetailsService.UpdateIPayPaymentStatus(data.RefNo)
        } else {
            return fiber.NewError(fiber.StatusBadRequest, "Invalid Request Number")
        }

        r = "Thank you for payment."
    }

    return c.SendString(r)
}

// BackendPostResponse
//
// @Tags IPAY
// @Produce json
// @Param        request       body        dto.IPayPaymentResponseDto  true  "IPayPaymentResponseDto"
// @Success 200
// @Router /payment/ipay88/backend/backend_response [post]
func (cr *IpayController) BackendPostResponse(c fiber.Ctx) error {
    data := new(dto.IPayPaymentResponseDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    c.Set(fiber.HeaderContentType, "text/plain")
    qs := cr.ipayService.BuildResponseSignature(data.PaymentId, data.RefNo, data.Amount, data.Currency, data.Status)
    r := "Fail"
    if data.Status == "1" && qs == data.Signature {
        if strings.Contains(data.RefNo, "MB") {
            _ = cr.billPaymentDetailsService.UpdateIPayPaymentStatus(data.RefNo)
        } else if strings.Contains(data.RefNo, "MP") {
            _ = cr.packagePaymentDetailsService.UpdateIPayPaymentStatus(data.RefNo)
        } else {
            return fiber.NewError(fiber.StatusBadRequest, "Invalid Request Number")
        }

        r = "RECEIVEOK"
    }
    
    return c.SendString(r)
}
