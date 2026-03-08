package userNotification

import (
    "fmt"
    "strconv"
    "vesaliusm/database"
    "vesaliusm/middleware"
    applicationUserNotificationService "vesaliusm/service/applicationUserNotification"
    generalNotificationMasterService "vesaliusm/service/generalNotificationMaster"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

var applicationUserNotificationSvc *applicationUserNotificationService.ApplicationUserNotificationService = applicationUserNotificationService.NewApplicationUserNotificationService(database.GetDb(), database.GetCtx())
var generalNotificationMasterSvc *generalNotificationMasterService.GeneralNotificationMasterService = generalNotificationMasterService.NewGeneralNotificationMasterService(database.GetDb(), database.GetCtx())

// GetUnseenNotificationCount
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Success 200 {integer} int
// @Router /notification/unseen/count [get]
func GetUnseenNotificationCount(c fiber.Ctx) error {
    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    count, err := applicationUserNotificationSvc.CountUnseenByUserId(user.UserID.Int64)
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
// @Success 200 {array} model.OnesignalNotification
// @Router /notification/all [get]
func GetNotificationList(c fiber.Ctx) error {
    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := applicationUserNotificationSvc.ListByUserId(user.UserID.Int64, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// GetGeneralNotificationList
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.GeneralNotification
// @Router /notification/general/master/all [get]
func GetGeneralNotificationList(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := generalNotificationMasterSvc.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
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
func GetByNotificationMasterId(c fiber.Ctx) error {
    notificationMasterId := c.Params("notificationMasterId")
    id, _ := strconv.ParseInt(notificationMasterId, 10, 64)
    o, err := generalNotificationMasterSvc.FindByNotificationMasterId(id)
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
func GetNotificationById(c fiber.Ctx) error {
    notificationId := c.Params("notificationId")
    id, _ := strconv.ParseInt(notificationId, 10, 64)
    o, err := applicationUserNotificationSvc.FindByNotificationId(id)
    if err != nil {
        return err
    }

    return c.JSON(o)
}
