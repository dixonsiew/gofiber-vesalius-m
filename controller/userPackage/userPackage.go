package userPackage

import (
    "fmt"
    "strconv"
    "strings"
    "vesaliusm/database"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    patientPurchaseDetailsService "vesaliusm/service/patientPurchaseDetails"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

var patientPurchaseDetailsSvc *patientPurchaseDetailsService.PatientPurchaseDetailsService = patientPurchaseDetailsService.NewPatientPurchaseDetailsService(database.GetDb(), database.GetCtx())

// CheckPackageExpiryMaxpurchase
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param request body dto.CheckPackageExpiryMaxpurchaseDto true "CheckPackageExpiryMaxpurchaseDto"
// @Success 200
// @Router /user-package/check/expiry-maxpurchase [post]
func CheckPackageExpiryMaxpurchase(c fiber.Ctx) error {
    data := new(dto.CheckPackageExpiryMaxpurchaseDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    cartIsValid := true
    cartResult := make([]interface{}, 0)

    for _, pkg := range data.Package {
        r, err := patientPurchaseDetailsSvc.CheckPackageExpiryMaxPurchase(pkg.PackageID, pkg.QuantityPurchased)
        if err != nil {
            return err
        }

        if r.Expired == 1 || r.Soldout == 1 || r.ExceedPurchase == 1 {
            cartIsValid = false
        }

        cartResult = append(cartResult, r)
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
func GetAllUserPurchaseHistory(c fiber.Ctx) error {
    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := patientPurchaseDetailsSvc.ListByPrn(user.UserID.Int64, page, limit)
    if err != nil {
        return err
    }
    
    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
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
func GetAllPurchaseHistory(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := patientPurchaseDetailsSvc.List(page, limit)
    if err != nil {
        return err
    }
    
    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
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
func SearchAllPurchaseHistory(c fiber.Ctx) error {
    var data fiber.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    var (
        key  string
        key2 string
        key3 string
        key4 string
    )
    if keyword, ok := data["keyword"]; ok {
        key = keyword.(string)
        if key != "" {
            key = "%" + strings.ToLower(key) + "%"
        }
    }
    if keyword, ok := data["keyword2"]; ok {
        key2 = keyword.(string)
        if key2 != "" {
            key2 = "%" + strings.ToLower(key2) + "%"
        }
    }
    if keyword, ok := data["keyword3"]; ok {
        key3 = keyword.(string)
        if key3 != "" {
            key3 = "%" + strings.ToLower(key3) + "%"
        }
    }
    if keyword, ok := data["keyword4"]; ok {
        v := keyword.(string)
        if v != "All" && v != "" {
            key4 = "%" + strings.ToLower(v) + "%"
        } else {
            key4 = v
        }
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := patientPurchaseDetailsSvc.ListByKeyword(key, key2, key3, key4, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// GetUserPackageById
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        purchaseId         path      string  true  "PurchaseId"
// @Success 200 {object} userPackage.UserPackage
// @Router /user-package/{purchaseId} [get]
func GetUserPackageById(c fiber.Ctx) error {
    purchaseId := c.Params("purchaseId")
    ipurchaseId, err := strconv.ParseInt(purchaseId, 10, 64)
    if err != nil {
        return err
    }

    o, err := patientPurchaseDetailsSvc.FindByPurchaseId(ipurchaseId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}
