package userPackage

import (
    "strconv"
    "strings"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    "vesaliusm/model/userPackage"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/patientPurchaseDetails"
    "vesaliusm/service/vesalius"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

type UserPackageController struct {
    applicationUserService        *applicationUser.ApplicationUserService
    patientPurchaseDetailsService *patientPurchaseDetails.PatientPurchaseDetailsService
    vesaliusService               *vesalius.VesaliusService
}

func NewUserPackageController() *UserPackageController {
    return &UserPackageController{
        applicationUserService:        applicationUser.ApplicationUserSvc,
        patientPurchaseDetailsService: patientPurchaseDetails.PatientPurchaseDetailsSvc,
        vesaliusService:               vesalius.VesaliusSvc,
    }
}

// CheckPackageExpiryMaxpurchase
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param request body dto.CheckPackageExpiryMaxpurchaseDto true "CheckPackageExpiryMaxpurchaseDto"
// @Success 200
// @Router /user-package/check/expiry-maxpurchase [post]
func (cr *UserPackageController) CheckPackageExpiryMaxpurchase(c fiber.Ctx) error {
    data := new(dto.CheckPackageExpiryMaxpurchaseDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    cartIsValid := true
    cartResult := make([]userPackage.PackageCheckResult, 0)

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

// GetAllUserPurchaseHistory
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
// @Success 200 {array} userPackage.UserPackage
// @Router /user-package/all/mobile [get]
func (cr *UserPackageController) GetAllUserPurchaseHistory(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.patientPurchaseDetailsService.ListByUserId(userId, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllPurchaseHistory
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
// @Success 200 {array} userPackage.UserPackage
// @Router /user-package/all [get]
func (cr *UserPackageController) GetAllPurchaseHistory(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.patientPurchaseDetailsService.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllPurchaseHistory
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
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
        key = "%" + strings.ToLower(key) + "%"
    }
    if key2 != "" {
        key2 = "%" + strings.ToLower(key2) + "%"
    }
    if key3 != "" {
        key3 = "%" + strings.ToLower(key3) + "%"
    }
    if key4 != "" && key4 != "All" {
        key4 = "%" + strings.ToLower(key4) + "%"
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.patientPurchaseDetailsService.ListByKeyword(key, key2, key3, key4, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetUserPackageById
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        purchaseId         path      string  true  "purchaseId"
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

    if data.Status != utils.PackageStatusPurchased &&
        data.Status != utils.PackageStatusBooked &&
        data.Status != utils.PackageStatusRedeemed &&
        data.Status != utils.PackageStatusCancelled {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid Package Status",
        })
    } else {
        if data.Status == utils.PackageStatusCancelled {
            apptRes, err := cr.patientPurchaseDetailsService.GetAppointmentDetailsByPurchaseId(ipurchaseId)
            if err != nil {
                return err
            }

            if apptRes != nil {
                dx := &dto.PostCancelAppointmentDto{
                    AppointmentNumber: apptRes.ApptNo.String,
                    Remark: apptRes.PackagePurchaseNo.String,
                }
                _, err = cr.vesaliusService.VesaliusGetCancelAppointment(apptRes.PatientPrn.String, dx)
                if err != nil {
                    return err
                }
            } else {
                err := cr.patientPurchaseDetailsService.UpdatePackageStatusByPurchaseId(ipurchaseId, utils.PackageStatusPurchased)
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
