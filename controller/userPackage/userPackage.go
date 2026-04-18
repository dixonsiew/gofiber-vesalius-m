package userPackage

import (
    "fmt"
    "strconv"
    "strings"
    "time"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    upck "vesaliusm/model/userPackage"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/exportExcel"
    "vesaliusm/service/patientPurchaseDetails"
    "vesaliusm/service/payment"
    "vesaliusm/service/vesalius"
    "vesaliusm/service/wallex"
    "vesaliusm/utils"
    "vesaliusm/utils/constants"

    "github.com/gofiber/fiber/v3"
)

type UserPackageController struct {
    applicationUserService        *applicationUser.ApplicationUserService
    patientPurchaseDetailsService *patientPurchaseDetails.PatientPurchaseDetailsService
    paymentService                *payment.PaymentService
    vesaliusService               *vesalius.VesaliusService
    wallexService                 *wallex.WallexService
    exportExcelService            *exportExcel.ExportExcelService
}

func NewUserPackageController() *UserPackageController {
    return &UserPackageController{
        applicationUserService:        applicationUser.ApplicationUserSvc,
        patientPurchaseDetailsService: patientPurchaseDetails.PatientPurchaseDetailsSvc,
        paymentService:                payment.PaymentSvc,
        vesaliusService:               vesalius.VesaliusSvc,
        wallexService:                 wallex.WallexSvc,
        exportExcelService:            exportExcel.ExportExcelSvc,
    }
}

// CheckPackageExpiryMaxpurchase
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param    request    body    dto.CheckPackageExpiryMaxpurchaseDto    true    "CheckPackageExpiryMaxpurchaseDto"
// @Success 200
// @Router /user-package/check/expiry-maxpurchase [post]
func (cr *UserPackageController) CheckPackageExpiryMaxpurchase(c fiber.Ctx) error {
    data := new(dto.CheckPackageExpiryMaxpurchaseDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    cartIsValid := true
    cartResult := make([]upck.PackageCheckResult, 0)

    for _, pkg := range data.Package {
        r, err := cr.patientPurchaseDetailsService.CheckPackageExpiryMaxPurchase(pkg.PackageId, pkg.QuantityPurchased)
        if err != nil {
            return err
        }

        if r.Expired == 1 || r.Soldout == 1 || r.ExceedPurchase == 1 {
            cartIsValid = false
        }

        cartResult = append(cartResult, *r)
    }

    return c.JSON(fiber.Map{
        "cartIsValid": cartIsValid,
        "cartResult":  cartResult,
    })
}

// CreateUserPurchaseDetails
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        paymentMethod     path      int                      true   "paymentMethod"
// @Param        request           body      dto.CreateUserPackageDto true   "CreateUserPackageDto"
// @Success 200
// @Router /user-package/purchase/{paymentMethod} [post]
func (cr *UserPackageController) CreateUserPurchaseDetails(c fiber.Ctx) error {
    paymentMethod := c.Params("paymentMethod")
    ipaymentMethod, _ := strconv.ParseInt(paymentMethod, 10, 32)
    data := new(dto.CreateUserPackageDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    fullContact := fmt.Sprintf("%s%s", data.UserPackagePayment.BillingContactCode, data.UserPackagePayment.BillingContactNo)
    if len(fullContact) > 20 {
        return fiber.NewError(fiber.StatusBadRequest, "Billing Contact Number is too long")
    }

    userPackage := make([]upck.UserPackage, 0)
    products := make([]utils.Map, 0)
    totalPrice := 0.0

    paymentRefNo, err := cr.paymentService.GeneratePaymentRefNo()
    if err != nil {
        return err
    }

    for _, pkg := range data.UserPackage {
        p := utils.Map{
            "name":        pkg.PackageName,
            "description": pkg.PackageName,
            "price":       pkg.PackagePrice,
            "quantity":    pkg.QuantityPurchased,
            "createdAt":   strconv.FormatInt(time.Now().UnixMilli(), 10),
            "updatedAt":   strconv.FormatInt(time.Now().UnixMilli(), 10),
            "deletedAt":   "",
        }
        products = append(products, p)

        totalPrice = totalPrice + pkg.PackagePrice*float64(pkg.QuantityPurchased)

        o := upck.UserPackage{
            PatientPrn:        utils.NewNullString("-"),
            PatientName:       utils.NewNullString(data.UserPackagePayment.BillingFullname),
            PackageId:         utils.NewInt64(pkg.PackageId),
            QuantityPurchased: pkg.QuantityPurchased,
            PackageStatus:     utils.NewNullString(string(constants.PackageStatusOrdered)),
        }
        userPackage = append(userPackage, o)
    }

    lsaddr := make([]string, 0)
    pyt := data.UserPackagePayment
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
            "paymentPurpose":          "SCVE",
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
            "products": products,
        }
        wallexRes, err := cr.wallexService.SubmitPaymentRequest(wallexPrm)
        if err != nil {
            return err
        }

        if totalPrice != wallexRes.Amount {
            return fiber.NewError(fiber.StatusBadRequest, "Incorrect Total Price")
        }

        userPackagePayment := upck.PackagePaymentDetails{
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
        }
        err = cr.paymentService.SaveGuestWallex(userPackagePayment, userPackage)
        if err != nil {
            return err
        }

        return c.JSON(fiber.Map{
            "message": "User Purchase Details created",
            "wallexDetails": fiber.Map{
                "expiredAt":  wallexRes.ExpiredAt,
                "paymentUrl": wallexRes.PaymentUrl,
            },
        })
    } else {
        userPackagePayment := upck.PackagePaymentDetails{
            PaymentGateway:         utils.NewInt32(int32(ipaymentMethod)),
            PaymentRequestNo:       utils.NewNullString(paymentRefNo),
            PaymentRequestCurrency: utils.NewNullString("MYR"),
            PaymentAmount:          utils.NewFloat(totalPrice),
            PaymentPurpose:         utils.NewNullString("Hospital Pkg Purchase"),
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
        }
        err = cr.paymentService.SaveGuestIPay(userPackagePayment, userPackage)
        if err != nil {
            return err
        }

        return c.JSON(paymentRefNo)
    }
}

// GetAllUserPurchaseHistory
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        _page              query      int  false  "_page"  default:"1"
// @Param        _limit             query      int  false  "_limit" default:"10"
// @Success 200 {array} userPackage.UserPackage
// @Router /user-package/all/mobile [get]
func (cr *UserPackageController) GetAllUserPurchaseHistory(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.patientPurchaseDetailsService.ListByUserId(userId, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllPurchaseHistory
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        _page              query      int  false  "_page"  default:"1"
// @Param        _limit             query      int  false  "_limit" default:"10"
// @Success 200 {array} userPackage.UserPackage
// @Router /user-package/all [get]
func (cr *UserPackageController) GetAllPurchaseHistory(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.patientPurchaseDetailsService.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllPurchaseHistory
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        _page              query      int  false  "_page"  default:"1"
// @Param        _limit             query      int  false  "_limit" default:"10"
// @Param request body dto.SearchPurchaseHistoryDto false "Keyword"
// @Success 200 {array} userPackage.UserPackage
// @Router /user-package/all [post]
func (cr *UserPackageController) SearchAllPurchaseHistory(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key2 := data.GetString("keyword2")
    key3 := data.GetString("keyword3")
    key4 := data.GetString("keyword4")

    if key != "" {
        key = "%" + key + "%"
    }
    if key2 != "" {
        key2 = "%" + key2 + "%"
    }
    if key3 != "" {
        key3 = "%" + key3 + "%"
    }
    if key4 != "" && key4 != "All" {
        key4 = "%" + key4 + "%"
    }

    x := dto.SearchKeyword4Dto{
        Keyword:  key,
        Keyword2: key2,
        Keyword3: key3,
        Keyword4: key4,
    }
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.patientPurchaseDetailsService.ListByKeyword(x, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetUserPackageById
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        purchaseId         path      int  true  "purchaseId"
// @Success 200 {object} userPackage.UserPackage
// @Router /user-package/{purchaseId} [get]
func (cr *UserPackageController) GetUserPackageById(c fiber.Ctx) error {
    purchaseId := c.Params("purchaseId")
    ipurchaseId, _ := strconv.ParseInt(purchaseId, 10, 64)
    o, err := cr.patientPurchaseDetailsService.FindByPurchaseId(ipurchaseId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// UpdateUserPackageStatus
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        purchaseId         path      string  true  "purchaseId"
// @Success 200
// @Router /user-package/status/{purchaseId} [post]
func (cr *UserPackageController) UpdateUserPackageStatus(c fiber.Ctx) error {
    data := new(dto.UserPackageStatusDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    purchaseId := c.Params("purchaseId")
    ipurchaseId, _ := strconv.ParseInt(purchaseId, 10, 64)

    if data.Status != constants.PackageStatusPurchased &&
        data.Status != constants.PackageStatusBooked &&
        data.Status != constants.PackageStatusRedeemed &&
        data.Status != constants.PackageStatusCancelled {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid Package Status",
        })
    } else {
        if data.Status == constants.PackageStatusCancelled {
            apptRes, err := cr.patientPurchaseDetailsService.GetAppointmentDetailsByPurchaseId(ipurchaseId)
            if err != nil {
                return err
            }

            if apptRes != nil {
                dx := &dto.PostCancelAppointmentDto{
                    AppointmentNumber: apptRes.ApptNo.String,
                    Remark:            apptRes.PackagePurchaseNo.String,
                }
                _, err = cr.vesaliusService.VesaliusGetCancelAppointment(apptRes.PatientPrn.String, dx)
                if err != nil {
                    return err
                }
            } else {
                err := cr.patientPurchaseDetailsService.UpdatePackageStatusByPurchaseId(ipurchaseId, constants.PackageStatusPurchased)
                if err != nil {
                    return err
                }
            }

        } else {
            err := cr.patientPurchaseDetailsService.UpdatePackageStatusByPurchaseId(ipurchaseId, data.Status)
            if err != nil {
                return err
            }
        }
    }

    return c.JSON(fiber.Map{
        "message": "User Package Status updated",
    })
}

// GetAllExportPurchaseHistory
//
// @Tags User Package
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} upck.UserPackage
// @Router /user-package/export/all [get]
func (cr *UserPackageController) GetAllExportPurchaseHistory(c fiber.Ctx) error {
    lx, err := cr.exportExcelService.ExportAllPurchaseHistory()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetSearchExportPurchaseHistory
//
// @Tags User Package
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param        request           body        dto.SearchKeyword4Dto    false  "Keyword"
// @Success 200 {array} upck.UserPackage
// @Router /user-package/export/search [post]
func (cr *UserPackageController) GetSearchExportPurchaseHistory(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key2 := data.GetString("keyword2")
    key3 := data.GetString("keyword3")
    key4 := data.GetString("keyword4")

    if key != "" {
        key = "%" + key + "%"
    }
    if key2 != "" {
        key2 = "%" + key2 + "%"
    }
    if key3 != "" {
        key3 = "%" + key3 + "%"
    }
    if key4 != "" && key4 != "All" {
        key4 = "%" + key4 + "%"
    }

    x := dto.SearchKeyword4Dto{
        Keyword:  key,
        Keyword2: key2,
        Keyword3: key3,
        Keyword4: key4,
    }
    lx, err := cr.exportExcelService.ExportPurchaseHistoryByKeyword(x)
    if err != nil {
        return err
    }

    return c.JSON(lx)
}
