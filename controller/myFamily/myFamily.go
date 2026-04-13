package myFamily

import (
    "strconv"
    "vesaliusm/middleware"
    "vesaliusm/service/applicationUserFamily"
    "vesaliusm/utils/constants"

    "github.com/gofiber/fiber/v3"
)

type MyFamilyController struct {
    applicationUserFamilyService *applicationUserFamily.ApplicationUserFamilyService
}
    
func NewMyFamilyController() *MyFamilyController {
    return &MyFamilyController{
        applicationUserFamilyService: applicationUserFamily.ApplicationUserFamilySvc,
    }
}

// GetAllActiveUserFamilies
//
// @Tags MyFamily
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"       default:"1"
// @Param        _limit             query      string  false  "_limit"      default:"10"
// @Param        _self              query      string  false  "_self"       default:"0"
// @Param        _isForAppt         query      string  false  "_isForAppt"  default:"0"
// @Success 200 {array} model.ApplicationUserFamily
// @Router /my-family [get]
func (cr *MyFamilyController) GetAllActiveUserFamilies(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    self := c.Query("_self", "0")
    isForAppt := c.Query("_isForAppt", "0")
    m, err := cr.applicationUserFamilyService.ListActiveByUserId(userId, self != "0", isForAppt != "0", self, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllUserFamilies
//
// @Tags MyFamily
// @Produce json
// @Security BearerAuth
// @Param        userId             path       string  true   "UserId"
// @Param        _page              query      string  false  "_page"       default:"1"
// @Param        _limit             query      string  false  "_limit"      default:"10"
// @Param        _self              query      string  false  "_self"       default:"0"
// @Param        _isForAppt         query      string  false  "_isForAppt"  default:"0"
// @Success 200 {array} model.ApplicationUserFamily
// @Router /my-family/all/{userId} [get]
func (cr *MyFamilyController) GetAllUserFamilies(c fiber.Ctx) error {
    userId := c.Params("userId")
    iuserId, _ := strconv.ParseInt(userId, 10, 64)
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    self := c.Query("_self", "0")
    isForAppt := c.Query("_isForAppt", "0")
    m, err := cr.applicationUserFamilyService.ListByUserId(iuserId, self != "0", isForAppt != "0", self, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetFamilyById
//
// @Tags MyFamily
// @Produce json
// @Security BearerAuth
// @Param        familyId             path       string  true   "FamilyId"
// @Success 200 {object} model.ApplicationUserFamily
// @Router /my-family/familyId/{familyId} [get]
func (cr *MyFamilyController) GetFamilyById(c fiber.Ctx) error {
    familyId := c.Params("familyId")
    ifamilyId, _ := strconv.ParseInt(familyId, 10, 64)
    o, err := cr.applicationUserFamilyService.FindByFamilyId(ifamilyId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}
