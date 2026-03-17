package userNotification

import (
	"strconv"
	"vesaliusm/middleware"
	"vesaliusm/service/applicationUserNotification"
	"vesaliusm/service/generalNotificationMaster"
	"vesaliusm/utils"

	"github.com/gofiber/fiber/v3"
)

type UserNotificationController struct {
	applicationUserNotificationService *applicationUserNotification.ApplicationUserNotificationService
	generalNotificationMasterService   *generalNotificationMaster.GeneralNotificationMasterService
}

func NewUserNotificationController() *UserNotificationController {
	return &UserNotificationController{
		applicationUserNotificationService: applicationUserNotification.ApplicationUserNotificationSvc,
		generalNotificationMasterService:   generalNotificationMaster.GeneralNotificationMasterSvc,
	}
}

// GetUnseenNotificationCount
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Success 200 {integer} int
// @Router /notification/unseen/count [get]
func (cr *UserNotificationController) GetUnseenNotificationCount(c fiber.Ctx) error {
	_, user, err := middleware.ValidateToken(c)
	if err != nil {
		return err
	}

	if user == nil {
		return middleware.Unauthorized(c)
	}

	count, err := cr.applicationUserNotificationService.CountUnseenByUserId(user.UserId.Int64)
	if err != nil {
		return err
	}

	return c.Send([]byte(strconv.Itoa(count)))
}

// GetNotificationList
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
// @Success 200 {array} model.OnesignalNotification
// @Router /notification/all [get]
func (cr *UserNotificationController) GetNotificationList(c fiber.Ctx) error {
	_, user, err := middleware.ValidateToken(c)
	if err != nil {
		return err
	}

	if user == nil {
		return middleware.Unauthorized(c)
	}

	page := c.Query("_page", "1")
	limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
	m, err := cr.applicationUserNotificationService.ListByUserId(user.UserId.Int64, page, limit)
	if err != nil {
		return err
	}

	c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
	c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
	return c.JSON(m.List)
}

// GetGeneralNotificationList
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
// @Success 200 {array} model.GeneralNotification
// @Router /notification/general/master/all [get]
func (cr *UserNotificationController) GetGeneralNotificationList(c fiber.Ctx) error {
	page := c.Query("_page", "1")
	limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
	m, err := cr.generalNotificationMasterService.List(page, limit)
	if err != nil {
		return err
	}

	c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
	c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
	return c.JSON(m.List)
}

// GetByNotificationMasterId
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Param        notificationMasterId              path      string  true  "Notification MasterId"
// @Success 200 {object} model.GeneralNotification
// @Router /notification/general/master/{notificationMasterId} [get]
func (cr *UserNotificationController) GetByNotificationMasterId(c fiber.Ctx) error {
	notificationMasterId := c.Params("notificationMasterId")
	id, _ := strconv.ParseInt(notificationMasterId, 10, 64)
	o, err := cr.generalNotificationMasterService.FindByNotificationMasterId(id)
	if err != nil {
		return err
	}

	return c.JSON(o)
}

// GetNotificationById
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Param        notificationId              path      string  true  "Notification Id"
// @Success 200 {object} model.OnesignalNotification
// @Router /notification/{notificationId} [get]
func (cr *UserNotificationController) GetNotificationById(c fiber.Ctx) error {
	notificationId := c.Params("notificationId")
	id, _ := strconv.ParseInt(notificationId, 10, 64)
	o, err := cr.applicationUserNotificationService.FindByNotificationId(id)
	if err != nil {
		return err
	}

	return c.JSON(o)
}
