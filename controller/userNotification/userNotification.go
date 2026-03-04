package userNotification

import (
    "fmt"
    "strconv"
    "vesaliusm/middleware"
    applicationUserNotificationService "vesaliusm/service/applicationUserNotification"
    generalNotificationMasterService "vesaliusm/service/generalNotificationMaster"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

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
    
    count, err := applicationUserNotificationService.CountUnseenByUserId(user.UserID.Int64)
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
// @Success 200 {array} model.OneSignalNotification
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
    m, err := applicationUserNotificationService.ListByUserId(user.UserID.Int64, page, limit)
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
    m, err := generalNotificationMasterService.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}
