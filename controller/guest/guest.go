package guest

import (
    "strconv"
    "vesaliusm/service/clubs"
    "vesaliusm/service/guest"
    "vesaliusm/service/novaDoctor"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

type GuestController struct {
    clubsService      *clubs.ClubService
    guestService      *guest.GuestService
    novaDoctorService *novaDoctor.NovaDoctorService
}

func NewGuestController() *GuestController {
    return &GuestController{
        clubsService:      clubs.ClubSvc,
        guestService:      guest.GuestSvc,
        novaDoctorService: novaDoctor.NovaDoctorSvc,
    }
}

// GetAllDoctorInformation
//
// @Tags Guest
// @Produce json
// @Param        branchId          path        string  true  "branchId"
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.NovaDoctor
// @Router /guest/vesalius/getAllDoctorInformation/{branchId} [get]
func (cr *GuestController) GetAllDoctorInformation(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.novaDoctorService.List(page, limit, false)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllDoctorInformation
//
// @Tags Guest
// @Produce json
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.NovaDoctor
// @Router /guest/vesalius/getAllDoctorInformation/{branchId} [post]
func (cr *GuestController) SearchAllDoctorInformation(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key = "%" + key + "%"

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.novaDoctorService.ListByKeywordGuest(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllGuestNotificationLists
//
// @Tags Guest
// @Produce json
// @Param        playerId          path        string  true  "playerId"
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.OnesignalNotification
// @Router /guest/notification/all/{playerId} [get]
func (cr *GuestController) getAllGuestNotificationLists(c fiber.Ctx) error {
    playerId := c.Params("playerId")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.guestService.ListGuestNotification(playerId, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetGuestUnseenNotificationCount
//
// @Tags Guest
// @Produce json
// @Param        notificationId    path        string  true  "notificationId"
// @Param        playerId          path        string  true  "playerId"
// @Success 200
// @Router /guest/notification/seen/{notificationId}/{playerId} [get]
func (cr *GuestController) GetGuestUnseenNotificationCount(c fiber.Ctx) error {
    notificationId := c.Params("notificationId")
    playerId := c.Params("playerId")
    inotificationId, _ := strconv.ParseInt(notificationId, 10, 64)
    err := cr.guestService.UpdateSeenNotificationByPlayerId(playerId, inotificationId)
    if err != nil {
        return err
    }

    n, err := cr.guestService.CountUnseenNotificationByGuestPlayerId(playerId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage":  "Guest Notification has been set to seen",
        "userUnseenCount": n,
    })
}

// GetLittleKidsAboutUs
//
// @Tags Guest
// @Produce json
// @Success 200 {object} clubs.LittleExplorersKidsAboutUs
// @Router /guest/clubs/littlekids/about-us [get]
func (cr *GuestController) GetLittleKidsAboutUs(c fiber.Ctx) error {
    o, err := cr.clubsService.FindLittleKidsAboutUs()
    if err != nil {
        return err
    }

    return c.JSON(o)
}
