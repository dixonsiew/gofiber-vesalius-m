package logistic

import (
    "strconv"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    lg "vesaliusm/model/logistic"
    "vesaliusm/service/exportExcel"
    "vesaliusm/service/logistic"
    "vesaliusm/service/mail"
    "vesaliusm/utils"
    "vesaliusm/utils/constants"

    "github.com/gofiber/fiber/v3"
    "github.com/nleeper/goment"
)

type LogisticController struct {
    logisticService    *logistic.LogisticService
    exportExcelService *exportExcel.ExportExcelService
    mailService        *mail.MailService
}

func NewLogisticController() *LogisticController {
    return &LogisticController{
        logisticService:    logistic.LogisticSvc,
        exportExcelService: exportExcel.ExportExcelSvc,
        mailService:        mail.MailSvc,
    }
}

// CreateLogisticSetup
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param    request    body    dto.LogisticSetupDto    true    "LogisticSetupDto"
// @Success 200
// @Router /logistic/setup [post]
func (cr *LogisticController) CreateLogisticSetup(c fiber.Ctx) error {
    data := new(dto.LogisticSetupDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    adminId, _, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    b, err := cr.logisticService.ExistsLogisticSetup()
    if err != nil {
        return err
    }
    if b {
        return fiber.NewError(fiber.StatusBadRequest, "Already setup Logistic Arrangement previously")
    }

    o := lg.LogisticSetup{
        LogisticSetupValue: utils.NewNullString(data.LogisticSetupValue),
    }

    logisticSetupId, err := cr.logisticService.SaveLogisticSetup(o, adminId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message":      "Logistic Arrangement Setup created",
        "kids_club_id": logisticSetupId,
    })
}

// UpdateLogisticSetup
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param     logisticSetupId    path    int                      true     "logisticSetupId"
// @Param     request            body    dto.LogisticSetupDto     true     "LogisticSetupDto"
// @Success 200
// @Router /logistic/setup/{logisticSetupId} [put]
func (cr *LogisticController) UpdateLogisticSetup(c fiber.Ctx) error {
    logisticSetupId := c.Params("logisticSetupId")
    ilogisticSetupId, _ := strconv.ParseInt(logisticSetupId, 10, 64)

    data := new(dto.LogisticSetupDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    adminId, _, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    b, err := cr.logisticService.ExistsLogisticSetup()
    if err != nil {
        return err
    }
    if !b {
        return fiber.NewError(fiber.StatusBadRequest, "Logistic Arrangement Setup does not exist")
    }

    o := lg.LogisticSetup{
        LogisticSetupId:    utils.NewInt64(ilogisticSetupId),
        LogisticSetupValue: utils.NewNullString(data.LogisticSetupValue),
    }

    err = cr.logisticService.UpdateLogisticSetup(o, adminId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Logistic Arrangement Setup updated",
    })
}

// GetLogisticSetup
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Success 200 {object} lg.LogisticSetup
// @Router /logistic/setup [get]
func (cr *LogisticController) GetLogisticSetup(c fiber.Ctx) error {
    o, err := cr.logisticService.FindLogisticSetup()
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// CreateLogisticSlot
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param    request    body    dto.LogisticSlotsDto    true    "LogisticSlotsDto"
// @Success 200
// @Router /logistic/slot [post]
func (cr *LogisticController) CreateLogisticSlot(c fiber.Ctx) error {
    data := new(dto.LogisticSlotsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    adminId, _, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    o := lg.LogisticSlots{
        LogisticSlots: []lg.LogisticSlot{},
    }
    if len(data.LogisticSlots) > 0 {
        for i := range data.LogisticSlots {
            x := data.LogisticSlots[i]
            o.LogisticSlots = append(o.LogisticSlots, x.ToDbModel())
        }
    }

    err = cr.logisticService.SaveLogisticSlot(o, adminId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Logistic Slot created",
    })
}

// GetAllAppLogisticSlots
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param    request    body    dto.LogisticSlotMobileDto    true    "LogisticSlotMobileDto"
// @Success 200 {array} lg.LogisticSlot
// @Router /logistic/slot/all/mobile [post]
func (cr *LogisticController) GetAllAppLogisticSlots(c fiber.Ctx) error {
    data := new(dto.LogisticSlotMobileDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    lx, err := cr.logisticService.FindAllAppLogisticSlots(data)
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetAllLogisticSlots
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Success 200 {array} lg.LogisticSlot
// @Router /logistic/slot/all [get]
func (cr *LogisticController) GetAllLogisticSlots(c fiber.Ctx) error {
    lx, err := cr.logisticService.FindAllLogisticSlots()
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "logisticArrangementSlots": lx,
    })
}

// CreateLogisticRequest
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param    request    body    dto.LogisticRequestDto    true    "LogisticRequestDto"
// @Success 200
// @Router /logistic/request [post]
func (cr *LogisticController) CreateLogisticRequest(c fiber.Ctx) error {
    data := new(dto.LogisticRequestDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    logisticRequestNumber, err := cr.logisticService.GenerateLogisticRequestNumber()
    if err != nil {
        return err
    }

    g, _ := goment.New(data.RequestedPickupDate, "DD/MM/YYYY")
    requestedPickupDay := g.Format("dddd")

    o := lg.LogisticRequest{
        LogisticRequestNumber:   utils.NewNullString(logisticRequestNumber),
        RequesterPrn:            utils.NewNullString(data.RequesterPrn),
        RequesterName:           utils.NewNullString(data.RequesterName),
        RequesterDob:            utils.NewNullString(data.RequesterDob),
        RequesterDocType:        utils.NewNullString(data.RequesterDocType),
        RequesterDocNumber:      utils.NewNullString(data.RequesterDocNumber),
        RequesterNationality:    utils.NewNullString(data.RequesterNationality),
        RequesterEmail:          utils.NewNullString(data.RequesterEmail),
        PrimaryDoctor:           utils.NewNullString(data.PrimaryDoctor),
        VisitWithCompanion:      utils.NewNullString(data.VisitWithCompanion),
        CompanionName:           utils.NewNullString(data.CompanionName),
        CompanionDob:            utils.NewNullString(data.CompanionDob),
        CompanionDocType:        utils.NewNullString(data.CompanionDocType),
        CompanionDocNumber:      utils.NewNullString(data.CompanionDocNumber),
        RelationshipToRequester: utils.NewNullString(data.RelationshipToRequester),
        FlightAirlineName:       utils.NewNullString(data.FlightAirlineName),
        FlightNumber:            utils.NewNullString(data.FlightNumber),
        FlightArrivalDate:       utils.NewNullString(data.FlightArrivalDate),
        FlightArrivalTime:       utils.NewNullString(data.FlightArrivalTime),
        RequestedPickupDate:     utils.NewNullString(data.RequestedPickupDate),
        RequestedPickupTime:     utils.NewNullString(data.RequestedPickupTime),
        RequestedPickupDay:      utils.NewNullString(requestedPickupDay),
    }
    err = cr.logisticService.SaveLogisticRequest(o)
    if err != nil {
        return err
    }

    go func() {
        cr.mailService.SendLogisticConfirmation(&o)
    }()

    return c.JSON(fiber.Map{
        "message": "Logistic Request created",
    })
}

// GetAllAppLogisticRequests
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int  false  "_page"  default:"1"
// @Param        _limit            query       int  false  "_limit" default:"10"
// @Success 200 {array} lg.LogisticRequest
// @Router /logistic/request/all/mobile [get]
func (cr *LogisticController) GetAllAppLogisticRequests(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.logisticService.ListAppLogisticRequests(userId, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllLogisticRequests
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int  false  "_page"  default:"1"
// @Param        _limit            query       int  false  "_limit" default:"10"
// @Success 200 {array} lg.LogisticRequest
// @Router /logistic/request/all [get]
func (cr *LogisticController) GetAllLogisticRequests(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.logisticService.ListLogisticRequests(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllLogisticRequests
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int                      false  "_page"  default:"1"
// @Param        _limit            query       int                      false  "_limit" default:"10"
// @Param        request           body        dto.SearchKeyword4Dto    false  "Keyword"
// @Success 200 {array} lg.LogisticRequest
// @Router /logistic/request/all [post]
func (cr *LogisticController) SearchAllLogisticRequests(c fiber.Ctx) error {
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
        if _, err := goment.New(key3, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong date format")
        }
    }
    if key4 != "" {
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
    m, err := cr.logisticService.ListLogisticRequestsByKeyword(x, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllExportLogisticRequest
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Success 200 {array} lg.LogisticRequest
// @Router /logistic/request/export/all [get]
func (cr *LogisticController) GetAllExportLogisticRequest(c fiber.Ctx) error {
    lx, err := cr.exportExcelService.ExportAllLogisticRequest()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetSearchExportLogisticRequest
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param        keyword           body        dto.SearchKeyword4Dto  false  "Search"
// @Success 200 {array} lg.LogisticRequest
// @Router /logistic/request/export/search [post]
func (cr *LogisticController) GetSearchExportLogisticRequest(c fiber.Ctx) error {
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
        if _, err := goment.New(key3, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong date format")
        }
    }
    if key4 != "" {
        key4 = "%" + key4 + "%"
    }

    x := dto.SearchKeyword4Dto{
        Keyword:  key,
        Keyword2: key2,
        Keyword3: key3,
        Keyword4: key4,
    }
    lx, err := cr.exportExcelService.ExportLogisticRequestByKeyword(x)
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetLogisticRequestById
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param    requestId    path    int    true    "requestId"
// @Success 200 {object} lg.LogisticRequest
// @Router /logistic/request/{requestId} [get]
func (cr *LogisticController) GetLogisticRequestById(c fiber.Ctx) error {
    requestId := c.Params("requestId")
    irequestId, _ := strconv.ParseInt(requestId, 10, 64)
    o, err := cr.logisticService.FindLogisticRequestByRequestId(irequestId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// UpdateAppLogisticRequestStatus
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param    request    body    dto.LogisticRequestStatusDto    true    "LogisticRequestStatusDto"
// @Success 200
// @Router /logistic/request/status [post]
func (cr *LogisticController) UpdateAppLogisticRequestStatus(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    data := new(dto.LogisticRequestStatusDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    if data.Status != constants.LogisticRequestStatusConfirmed &&
        data.Status != constants.LogisticRequestStatusCancelled &&
        data.Status != constants.LogisticRequestStatusRejected {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Logistic Request Status")
    } else {
        err := cr.logisticService.UpdateLogisticRequestStatusByRequestNumber(data.RequestNumber, data.Status, 0, userId)
        if err != nil {
            return err
        }

        if data.Status == constants.LogisticRequestStatusCancelled {
            o, err := cr.logisticService.FindLogisticRequestByRequestNumber(data.RequestNumber)
            if err != nil {
                return err
            }

            go func() {
                cr.mailService.SendLogisticCancellation(o)
            }()
        }
    }

    return c.JSON(fiber.Map{
        "message": "Logistic Request Status updated",
    })
}

// UpdateLogisticRequestStatus
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param    request    body    dto.LogisticRequestStatusDto    true    "LogisticRequestStatusDto"
// @Success 200
// @Router /logistic/request/status/webadmin [post]
func (cr *LogisticController) UpdateLogisticRequestStatus(c fiber.Ctx) error {
    adminId, _, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    data := new(dto.LogisticRequestStatusDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    if data.Status != constants.LogisticRequestStatusConfirmed &&
        data.Status != constants.LogisticRequestStatusCancelled &&
        data.Status != constants.LogisticRequestStatusRejected {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Logistic Request Status")
    } else {
        err := cr.logisticService.UpdateLogisticRequestStatusByRequestNumber(data.RequestNumber, data.Status, adminId, 0)
        if err != nil {
            return err
        }
    }

    return c.JSON(fiber.Map{
        "message": "Logistic Request Status updated",
    })
}
