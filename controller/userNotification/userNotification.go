package userNotification

import (
    "strconv"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    "vesaliusm/model"
    "vesaliusm/service/applicationUserNotification"
    "vesaliusm/service/generalNotificationMaster"
    "vesaliusm/service/notification"
    "vesaliusm/utils"
    "vesaliusm/utils/constants"

    "github.com/OneSignal/onesignal-go-api"
    "github.com/gofiber/fiber/v3"
    "github.com/nleeper/goment"
)

type UserNotificationController struct {
    applicationUserNotificationService *applicationUserNotification.ApplicationUserNotificationService
    generalNotificationMasterService   *generalNotificationMaster.GeneralNotificationMasterService
    notificationService                *notification.NotificationService
}

func NewUserNotificationController() *UserNotificationController {
    return &UserNotificationController{
        applicationUserNotificationService: applicationUserNotification.ApplicationUserNotificationSvc,
        generalNotificationMasterService:   generalNotificationMaster.GeneralNotificationMasterSvc,
        notificationService:                notification.NewNotificationService(),
    }
}

// SendNotification
//
// @Tags Notification
// @Produce json
// @Success 200
// @Router /notification/send-notification/{playerId} [post]
func (cr *UserNotificationController) SendNotification(c fiber.Ctx) error {
    // msgType := c.Query("msgType")
    playerId := c.Params("playerId")
    playerIdList := []string{playerId}

    notification := model.OnesignalNotification{
        NotificationTitle: utils.NewNullString("Testing Push Notification"),
        ShortMessage:      utils.NewNullString("Kindly ignore"),
        FullMessage:       utils.NewNullString("From OneSignal"),
        UserId:            utils.NewInt64(365),
        VisitType:         utils.NewNullString("Testing Visit Type"),
        AccountNo:         utils.NewNullString("A123"),
    }

    if len(playerIdList) > 0 {
        n := *onesignal.NewNotification(cr.notificationService.GetAppId())
        n.SetIncludePlayerIds(playerIdList)
        n.SetHeadings(onesignal.StringMap{En: &notification.NotificationTitle.String})
        n.SetContents(onesignal.StringMap{En: &notification.ShortMessage.String})
        n.SetIosBadgeType("Increase")
        n.SetIosBadgeCount(1)
        n.SetData(map[string]any{
            "count": 0,
        })

        res, err := cr.notificationService.CreateNotification(&n)
        if err != nil {
            return err
        }
        return c.JSON(fiber.Map{
            "success": res.Id,
        })
    }

    return c.JSON(fiber.Map{
        "success": -1,
    })
}

// SetNotificationSeen
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Param    notificationId    path    int    true    "notificationId"
// @Success 200
// @Router /notification/seen/{notificationId} [post]
func (cr *UserNotificationController) SetNotificationSeen(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    notificationId := c.Params("notificationId")
    inotificationId, _ := strconv.ParseInt(notificationId, 10, 64)
    err = cr.applicationUserNotificationService.UpdateSeenByUserId(userId, inotificationId)
    if err != nil {
        return err
    }

    n, err := cr.applicationUserNotificationService.CountUnseenByUserId(userId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage":  "Notification has been set to seen",
        "userUnseenCount": n,
    })
}

// GetUnseenNotificationCount
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Success 200 {integer} int
// @Router /notification/unseen/count [get]
func (cr *UserNotificationController) GetUnseenNotificationCount(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    count, err := cr.applicationUserNotificationService.CountUnseenByUserId(userId)
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
// @Param        _page              query      int  false  "_page"  default:"1"
// @Param        _limit             query      int  false  "_limit" default:"10"
// @Success 200 {array} model.OnesignalNotification
// @Router /notification/all [get]
func (cr *UserNotificationController) GetNotificationList(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.applicationUserNotificationService.ListByUserId(userId, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// CreateGeneralNotification
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Param        request               body       dto.GeneralNotificationDto  true  "GeneralNotificationDto"
// @Success 200
// @Router /notification/general/master [post]
func (cr *UserNotificationController) CreateGeneralNotification(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    data := new(dto.GeneralNotificationDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    if data.StartDate != "" {
        if _, err := goment.New(data.StartDate, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong Start Date format")
        }
    }

    if data.EndDate != "" {
        if _, err := goment.New(data.EndDate, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong End Date format")
        }
    }

    o := &model.GeneralNotification{
        NotificationTitle: utils.NewNullString(data.NotificationTitle),
        ShortMessage:      utils.NewNullString(data.ShortMessage),
        FullMessage:       utils.NewNullString(data.FullMessage),
        StartDate:         utils.NewNullString(data.StartDate),
        EndDate:           utils.NewNullString(data.EndDate),
        TargetAgeFrom:     utils.NewInt64(int64(data.TargetAgeFrom)),
        TargetAgeTo:       utils.NewInt64(int64(data.TargetAgeTo)),
        TargetGender:      utils.NewNullString(data.TargetGender),
        TargetNationality: utils.NewNullString(data.TargetNationality),
        TargetCity:        utils.NewNullString(data.TargetCity),
        TargetState:       utils.NewNullString(data.TargetState),
    }
    err = cr.generalNotificationMasterService.Save(o, userId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "General Info Notification created",
    })
}

// UpdateGeneralNotification
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Param        notificationMasterId     path      int                            true  "notificationMasterId"
// @Param        request                  body      dto.GeneralNotificationDto     true  "GeneralNotificationDto"
// @Success 200
// @Router /notification/general/master/{notificationMasterId} [put]
func (cr *UserNotificationController) UpdateGeneralNotification(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    data := new(dto.GeneralNotificationDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    notificationMasterId := c.Params("notificationMasterId")
    inotificationMasterId, _ := strconv.ParseInt(notificationMasterId, 10, 64)

    o := &model.GeneralNotification{
        NotificationMasterId: utils.NewInt64(inotificationMasterId),
        NotificationTitle:    utils.NewNullString(data.NotificationTitle),
        ShortMessage:         utils.NewNullString(data.ShortMessage),
        FullMessage:          utils.NewNullString(data.FullMessage),
        StartDate:            utils.NewNullString(data.StartDate),
        EndDate:              utils.NewNullString(data.EndDate),
        TargetAgeFrom:        utils.NewInt64(int64(data.TargetAgeFrom)),
        TargetAgeTo:          utils.NewInt64(int64(data.TargetAgeTo)),
        TargetGender:         utils.NewNullString(data.TargetGender),
        TargetNationality:    utils.NewNullString(data.TargetNationality),
        TargetCity:           utils.NewNullString(data.TargetCity),
        TargetState:          utils.NewNullString(data.TargetState),
    }
    err = cr.generalNotificationMasterService.Update(o, userId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "General Info Notification updated",
    })
}

// GetGeneralNotificationList
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Param        _page              query      int  false  "_page"  default:"1"
// @Param        _limit             query      int  false  "_limit" default:"10"
// @Success 200 {array} model.GeneralNotification
// @Router /notification/general/master/all [get]
func (cr *UserNotificationController) GetGeneralNotificationList(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.generalNotificationMasterService.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetByNotificationMasterId
//
// @Tags Notification
// @Produce json
// @Security BearerAuth
// @Param        notificationMasterId              path      int  true  "notificationMasterId"
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
// @Param        notificationId              path      int  true  "notificationId"
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
