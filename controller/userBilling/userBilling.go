package userBilling

import (
    "fmt"
    "strconv"
    "strings"
    "time"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    model "vesaliusm/model/userBilling"
    "vesaliusm/service/billPaymentDetails"
    "vesaliusm/service/patientOutstandingBill"
    "vesaliusm/service/payment"
    "vesaliusm/service/wallex"
    "vesaliusm/utils"
    "vesaliusm/utils/constants"

    "github.com/gofiber/fiber/v3"
)

type UserBillingController struct {
    billPaymentDetailsService     *billPaymentDetails.BillPaymentDetailsService
    paymentService                *payment.PaymentService
    patientOutstandingBillService *patientOutstandingBill.PatientOutstandingBillService
    wallexService                 *wallex.WallexService
}

func NewUserBillingController() *UserBillingController {
    return &UserBillingController{
        billPaymentDetailsService:     billPaymentDetails.BillPaymentDetailsSvc,
        paymentService:                payment.PaymentSvc,
        patientOutstandingBillService: patientOutstandingBill.PatientOutstandingBillSvc,
        wallexService:                 wallex.WallexSvc,
    }
}

// CreateUserBillingDetails
// @Tags User Billing
// @Produce json
// @Security BearerAuth
// @Param    paymentMethod    path    int                         true    "paymentMethod"
// @Param    request          body    dto.CreateUserBillingDto    true    "CreateUserBillingDto"
// @Success 200
// @Router /user-billing/pay/{paymentMethod} [post]
func (cr *UserBillingController) CreateUserBillingDetails(c fiber.Ctx) error {
    paymentMethod := c.Params("paymentMethod")
    ipaymentMethod, _ := strconv.ParseInt(paymentMethod, 10, 32)
    data := new(dto.CreateUserBillingDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    bills := make([]utils.Map, 0)

    paymentRefNo, err := cr.paymentService.GeneratePaymentRefNo()
    if err != nil {
        return err
    }

    bill := data.UserBilling
    userBilling := model.UserBilling{
        PatientPrn:           bill.PatientPrn,
        PatientName:          bill.PatientName,
        PatientDocNumber:     bill.PatientDocNumber,
        VesRegDateTime:       bill.VesRegDateTime,
        VesBillNumber:        bill.VesBillNumber,
        VesInvoiceNumber:     bill.VesInvoiceNumber,
        VesInvoiceDateTime:   bill.VesInvoiceDateTime,
        VesInvoiceAmount:     bill.VesInvoiceAmount,
        VesOutstandingAmount: bill.VesOutstandingAmount,
        PaymentStatus:        constants.PaymentStatusSubmitted,
    }

    floatVesOutstandingAmount, _ := strconv.ParseFloat(strings.ReplaceAll(bill.VesOutstandingAmount, ",", ""), 64)

    m := utils.Map{
        "name":        bill.VesInvoiceNumber,
        "description": fmt.Sprintf("Outstanding Bill Payment for: %s", bill.VesInvoiceNumber),
        "price":       floatVesOutstandingAmount,
        "quantity":    1,
        "createdAt":   strconv.FormatInt(time.Now().UnixMilli(), 10),
        "updatedAt":   strconv.FormatInt(time.Now().UnixMilli(), 10),
        "deletedAt":   "",
    }
    bills = append(bills, m)

    lsaddr := make([]string, 0)
    pyt := data.UserBillingPayment
    if pyt.BillingAddress1 != "" {
        lsaddr = append(lsaddr, pyt.BillingAddress1)
    }

    if pyt.BillingAddress2 != "" {
        lsaddr = append(lsaddr, pyt.BillingAddress2)
    }

    if pyt.BillingAddress3 != "" {
        lsaddr = append(lsaddr, pyt.BillingAddress3)
    }

    if pyt.BillingTowncity != "" {
        lsaddr = append(lsaddr, pyt.BillingTowncity)
    }

    if pyt.BillingState != "" {
        lsaddr = append(lsaddr, pyt.BillingState)
    }

    addr := strings.Join(lsaddr, ", ")

    if ipaymentMethod == constants.PaymentMethodWallex {
        wallexPrm := utils.Map{
            "collectionRequestNumber": paymentRefNo,
            "currency":                "MYR",
            "paymentPurpose":          "WX14",
            "paymentCurrency":         "IDR",
            "paymentPartial":          false,
            "remarks":                 paymentRefNo,
            "customerInfo": utils.Map{
                "name":             pyt.BillingFullname,
                "ituTelephoneCode": pyt.BillingContactCode,
                "mobileNumber":     pyt.BillingContactNo,
                "email":            pyt.BillingEmail,
                "address":          addr,
            },
            "products": bills,
        }
        wallexRes, err := cr.wallexService.SubmitPaymentRequest(wallexPrm)
        if err != nil {
            return err
        }

        if floatVesOutstandingAmount != wallexRes.Amount {
            return fiber.NewError(fiber.StatusBadRequest, "Incorrect Total Price")
        }

        userBillPayment := model.BillingPayment{
            PaymentGateway:         utils.NewInt32(int32(ipaymentMethod)),
            PaymentRequestId:       utils.NewNullString(wallexRes.ID),
            PaymentRequestNo:       utils.NewNullString(wallexRes.CollectionRequestNumber),
            PaymentRefId:           utils.NewNullString(wallexRes.ReferenceId),
            PaymentRequestCurrency: utils.NewNullString(wallexRes.Currency),
            PaymentAmount:          utils.NewFloat(wallexRes.Amount),
            PaymentPurpose:         utils.NewNullString(wallexRes.PaymentPurpose),
            PaymentCurrency:        utils.NewNullString(wallexRes.PaymentCurrency),
            PaymentAmountCollected: utils.NewFloat(wallexRes.PaymentAmountCollected),
            PaymentRemarks:         utils.NewNullString(wallexRes.Remarks),
            PaymentStatus:          utils.NewNullString(wallexRes.Status),
            PaymentAuthCode:        utils.NewNullString(""),
            PaymentErrorDesc:       utils.NewNullString(""),
            PaymentTransDate:       utils.NewNullString(""),
            BillingFullname:        utils.NewNullString(pyt.BillingFullname),
            BillingAddress1:        utils.NewNullString(pyt.BillingAddress1),
            BillingAddress2:        utils.NewNullString(pyt.BillingAddress2),
            BillingAddress3:        utils.NewNullString(pyt.BillingAddress3),
            BillingTowncity:        utils.NewNullString(pyt.BillingTowncity),
            BillingState:           utils.NewNullString(pyt.BillingState),
            BillingPostcode:        utils.NewNullString(pyt.BillingPostcode),
            BillingCountryCode:     utils.NewNullString(pyt.BillingCountryCode),
            BillingContactNo:       utils.NewNullString(pyt.BillingContactNo),
            BillingContactCode:     utils.NewNullString(pyt.BillingContactCode),
            BillingEmail:           utils.NewNullString(pyt.BillingEmail),
            PaymentUrl:             wallexRes.PaymentUrl,
            VesPaymentReceiptNo:    "",
        }
        err = cr.paymentService.SaveBillingWallex(userBillPayment, userBilling)
        if err != nil {
            return err
        }

        return c.JSON(fiber.Map{
            "message": "User Billing Payment Details created",
            "wallexDetails": fiber.Map{
                "expiredAt":  wallexRes.ExpiredAt,
                "paymentUrl": wallexRes.PaymentUrl,
            },
        })
    } else {
        userBillPayment := model.BillingPayment{
            PaymentGateway:         utils.NewInt32(int32(ipaymentMethod)),
            PaymentRequestNo:       utils.NewNullString(paymentRefNo),
            PaymentRequestCurrency: utils.NewNullString("MYR"),
            PaymentAmount:          utils.NewFloat(floatVesOutstandingAmount),
            PaymentPurpose:         utils.NewNullString(fmt.Sprintf("Hospital Outstanding Bill Payment for: %s", bill.VesInvoiceNumber)),
            PaymentCurrency:        utils.NewNullString("MYR"),
            PaymentRemarks:         utils.NewNullString(paymentRefNo),
            PaymentStatus:          utils.NewNullString("unpaid"),
            PaymentAuthCode:        utils.NewNullString(""),
            PaymentErrorDesc:       utils.NewNullString(""),
            PaymentTransDate:       utils.NewNullString(""),
            BillingFullname:        utils.NewNullString(pyt.BillingFullname),
            BillingAddress1:        utils.NewNullString(pyt.BillingAddress1),
            BillingAddress2:        utils.NewNullString(pyt.BillingAddress2),
            BillingAddress3:        utils.NewNullString(pyt.BillingAddress3),
            BillingTowncity:        utils.NewNullString(pyt.BillingTowncity),
            BillingState:           utils.NewNullString(pyt.BillingState),
            BillingPostcode:        utils.NewNullString(pyt.BillingPostcode),
            BillingCountryCode:     utils.NewNullString(pyt.BillingCountryCode),
            BillingContactNo:       utils.NewNullString(pyt.BillingContactNo),
            BillingContactCode:     utils.NewNullString(pyt.BillingContactCode),
            BillingEmail:           utils.NewNullString(pyt.BillingEmail),
            VesPaymentReceiptNo:    "",
        }
        err := cr.paymentService.SaveBillingIPay(userBillPayment, userBilling)
        if err != nil {
            return err
        }

        return c.JSON(paymentRefNo)
    }
}

// GetAllUserPaidBillingHistory
//
// @Tags User Billing
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.UserBillingPayment
// @Router /user-billing/paid/all/mobile [get]
func (cr *UserBillingController) GetAllUserPaidBillingHistory(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.billPaymentDetailsService.ListPaidByPrn(userId, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// UpdateUserBillingPaymentStatus
//
// @Tags User Billing
// @Produce json
// @Security BearerAuth
// @Param    billPaymentId    path    int                                true    "billPaymentId"
// @Param    request          body    dto.UserBillingPaymentStatusDto    true    "UserBillingPaymentStatusDto"
// @Success 200
// @Router /user-billing/status/{billPaymentId} [post]
func (cr *UserBillingController) UpdateUserBillingPaymentStatus(c fiber.Ctx) error {
    billPaymentId := c.Params("billPaymentId")
    ibillPaymentId, _ := strconv.ParseInt(billPaymentId, 10, 64)
    data := new(dto.UserBillingPaymentStatusDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    if data.Status != constants.PaymentStatusPaid {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Billing Payment Status")
    }

    err := cr.patientOutstandingBillService.UpdatePaymentStatusByBillPaymentId(ibillPaymentId, data.Status, nil)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "User Package Status updated",
    })
}
