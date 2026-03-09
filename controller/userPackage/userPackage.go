package userPackage

import (
	"fmt"
	"strings"
	"vesaliusm/database"
	"vesaliusm/middleware"
	_ "vesaliusm/model/userPackage"
	patientPurchaseDetailsService "vesaliusm/service/patientPurchaseDetails"
	"vesaliusm/utils"

	"github.com/gofiber/fiber/v3"
)

var patientPurchaseDetailsSvc *patientPurchaseDetailsService.PatientPurchaseDetailsService = 
    patientPurchaseDetailsService.NewPatientPurchaseDetailsService(database.GetDb(), database.GetCtx())

// SearchAllPurchaseHistory
//
// @Tags User Package
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
// @Param request body map[string]string false "Keyword"
// @Success 200 {array} userPackage.UserPackage
// @Router /user-package/all [post]
func SearchAllPurchaseHistory(c fiber.Ctx) error {
    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    var data fiber.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    var (
        key string
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