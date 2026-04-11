package logistic

import (
    "bytes"
    "fmt"
    "strconv"
    "strings"
    "time"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    logisticModel "vesaliusm/model/logistic"
    logisticService "vesaliusm/service/logistic"
    "vesaliusm/service/mail"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/xuri/excelize/v2"
)

type LogisticController struct {
    logisticService *logisticService.LogisticService
    mailService     *mail.MailService
}

func NewLogisticController() *LogisticController {
    return &LogisticController{
        logisticService: logisticService.LogisticSvc,
        mailService:     mail.MailSvc,
    }
}

// CreateLogisticSetup
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param request body dto.LogisticSetupDto true "LogisticSetupDto"
// @Success 200 {object} map[string]interface{}
// @Router /logistic/setup [post]
func (cr *LogisticController) CreateLogisticSetup(c fiber.Ctx) error {
    data := new(dto.LogisticSetupDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    _, admin, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    exists, err := cr.logisticService.ExistsSetup()
    if err != nil {
        return err
    }
    if exists {
        return fiber.NewError(fiber.StatusBadRequest, "Already setup Logistic Arrangement previously")
    }

    setup, err := cr.logisticService.CreateSetup(data, admin.AdminId.Int64)
    if err != nil {
        return err
    }

    result := fiber.Map{"message": "Logistic Arrangement Setup created"}
    if setup != nil {
        result["logisticSetupId"] = setup.LogisticSetupId.Int64
    }

    return c.JSON(result)
}

// UpdateLogisticSetup
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param logisticSetupId path string true "LogisticSetupId"
// @Param request body dto.LogisticSetupDto true "LogisticSetupDto"
// @Success 200 {object} map[string]interface{}
// @Router /logistic/setup/{logisticSetupId} [put]
func (cr *LogisticController) UpdateLogisticSetup(c fiber.Ctx) error {
    setupID, err := strconv.ParseInt(c.Params("logisticSetupId"), 10, 64)
    if err != nil {
        return err
    }

    data := new(dto.LogisticSetupDto)
    if err = utils.BindNValidate(c, data); err != nil {
        return err
    }

    _, admin, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    exists, err := cr.logisticService.ExistsSetup()
    if err != nil {
        return err
    }
    if !exists {
        return fiber.NewError(fiber.StatusBadRequest, "Logistic Arrangement Setup does not exist")
    }

    err = cr.logisticService.UpdateSetup(setupID, data, admin.AdminId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{"message": "Logistic Arrangement Setup updated"})
}

// GetLogisticSetup
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Success 200 {object} logistic.LogisticSetup
// @Router /logistic/setup [get]
func (cr *LogisticController) GetLogisticSetup(c fiber.Ctx) error {
    setup, err := cr.logisticService.FindSetup()
    if err != nil {
        return err
    }
    return c.JSON(setup)
}

// CreateLogisticSlot
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param request body dto.LogisticSlotsDto true "LogisticSlotsDto"
// @Success 200 {object} map[string]interface{}
// @Router /logistic/slot [post]
func (cr *LogisticController) CreateLogisticSlot(c fiber.Ctx) error {
    data := new(dto.LogisticSlotsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    _, admin, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    err = cr.logisticService.ReplaceSlots(data, admin.AdminId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{"message": "Logistic Request created"})
}

// GetAllAppLogisticSlots
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param request body dto.LogisticSlotMobileDto true "LogisticSlotMobileDto"
// @Success 200 {array} logistic.LogisticSlot
// @Router /logistic/slot/all/mobile [post]
func (cr *LogisticController) GetAllAppLogisticSlots(c fiber.Ctx) error {
    data := new(dto.LogisticSlotMobileDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }
    if err := validateDateDDMMYYYY(data.FlightArrivalDate); err != nil {
        return err
    }

    lx, err := cr.logisticService.FindAppSlots(data)
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
// @Success 200 {object} map[string]interface{}
// @Router /logistic/slot/all [get]
func (cr *LogisticController) GetAllLogisticSlots(c fiber.Ctx) error {
    lx, err := cr.logisticService.FindAllSlots()
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{"logisticArrangementSlots": lx})
}

// CreateLogisticRequest
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param request body dto.LogisticRequestDto true "LogisticRequestDto"
// @Success 200 {object} map[string]interface{}
// @Router /logistic/request [post]
func (cr *LogisticController) CreateLogisticRequest(c fiber.Ctx) error {
    data := new(dto.LogisticRequestDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }
    if err := validateDateDDMMYYYY(data.FlightArrivalDate); err != nil {
        return err
    }
    if err := validateDateDDMMYYYY(data.RequestedPickupDate); err != nil {
        return err
    }

    request, err := cr.logisticService.CreateRequest(data)
    if err != nil {
        return err
    }

    go func() {
        cr.mailService.SendLogisticConfirmation(request)
    }()

    return c.JSON(fiber.Map{"message": "Logistic Request created"})
}

// GetAllAppLogisticRequests
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param _page query string false "_page" default:"1"
// @Param _limit query string false "_limit" default:"10"
// @Success 200 {array} logistic.LogisticRequest
// @Router /logistic/request/all/mobile [get]
func (cr *LogisticController) GetAllAppLogisticRequests(c fiber.Ctx) error {
    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.logisticService.ListAppRequests(user.UserId.Int64, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllLogisticRequests
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param _page query string false "_page" default:"1"
// @Param _limit query string false "_limit" default:"10"
// @Success 200 {array} logistic.LogisticRequest
// @Router /logistic/request/all [get]
func (cr *LogisticController) GetAllLogisticRequests(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.logisticService.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllLogisticRequests
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param _page query string false "_page" default:"1"
// @Param _limit query string false "_limit" default:"10"
// @Param request body dto.SearchKeyword2Dto false "Keyword"
// @Success 200 {array} logistic.LogisticRequest
// @Router /logistic/request/all [post]
func (cr *LogisticController) SearchAllLogisticRequests(c fiber.Ctx) error {
    keyword, keyword2, keyword3, keyword4, err := getSearchKeywords(c)
    if err != nil {
        return err
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.logisticService.ListByKeyword(keyword, keyword2, keyword3, keyword4, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// ExportAllLogisticRequests
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Success 200 {array} logistic.LogisticRequest
// @Router /logistic/request/export/all [get]
func (cr *LogisticController) ExportAllLogisticRequests(c fiber.Ctx) error {
    lx, err := cr.logisticService.ExportAll()
    if err != nil {
        return err
    }

    return cr.sendLogisticExportFile(c, lx, "logistic_request_all")
}

// ExportSearchLogisticRequests
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param request body dto.SearchKeyword2Dto false "Keyword"
// @Success 200 {array} logistic.LogisticRequest
// @Router /logistic/request/export/search [post]
func (cr *LogisticController) ExportSearchLogisticRequests(c fiber.Ctx) error {
    keyword, keyword2, keyword3, keyword4, err := getSearchKeywords(c)
    if err != nil {
        return err
    }

    lx, err := cr.logisticService.ExportByKeyword(keyword, keyword2, keyword3, keyword4)
    if err != nil {
        return err
    }

    return cr.sendLogisticExportFile(c, lx, "logistic_request_search")
}

// GetLogisticRequestByID
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param requestId path string true "RequestId"
// @Success 200 {object} logistic.LogisticRequest
// @Router /logistic/request/{requestId} [get]
func (cr *LogisticController) GetLogisticRequestByID(c fiber.Ctx) error {
    requestID, err := strconv.ParseInt(c.Params("requestId"), 10, 64)
    if err != nil {
        return err
    }

    o, err := cr.logisticService.FindByRequestID(requestID)
    if err != nil {
        return err
    }
    if o == nil {
        return fiber.NewError(fiber.StatusNotFound, "Logistic Request not found")
    }

    return c.JSON(o)
}

// UpdateAppLogisticRequestStatus
//
// @Tags Logistic
// @Produce json
// @Security BearerAuth
// @Param request body dto.LogisticRequestStatusDto true "LogisticRequestStatusDto"
// @Success 200 {object} map[string]interface{}
// @Router /logistic/request/status [post]
func (cr *LogisticController) UpdateAppLogisticRequestStatus(c fiber.Ctx) error {
    data := new(dto.LogisticRequestStatusDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }
    if err := validateLogisticStatus(data.Status); err != nil {
        return err
    }

    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    err = cr.logisticService.UpdateRequestStatusByNumber(data.RequestNumber, data.Status, 0, user.UserId.Int64)
    if err != nil {
        return err
    }

    if data.Status == utils.LogisticRequestStatusCancelled {
        request, err := cr.logisticService.FindByRequestNumberForMail(data.RequestNumber)
        if err != nil {
            return err
        } else if request != nil {
            go func() {
                cr.mailService.SendLogisticCancellation(request)
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
// @Param request body dto.LogisticRequestStatusDto true "LogisticRequestStatusDto"
// @Success 200 {object} map[string]interface{}
// @Router /logistic/request/status/webadmin [post]
func (cr *LogisticController) UpdateLogisticRequestStatus(c fiber.Ctx) error {
    data := new(dto.LogisticRequestStatusDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }
    if err := validateLogisticStatus(data.Status); err != nil {
        return err
    }

    _, admin, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    err = cr.logisticService.UpdateRequestStatusByNumber(data.RequestNumber, data.Status, admin.AdminId.Int64, 0)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Logistic Request Status updated",
    })
}

func getSearchKeywords(c fiber.Ctx) (string, string, string, string, error) {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return "", "", "", "", err
    }

    keyword := data.GetString("keyword")
    keyword2 := data.GetString("keyword2")
    keyword3 := data.GetString("keyword3")
    keyword4 := data.GetString("keyword4")

    if keyword3 != "" {
        if err := validateDateDDMMYYYY(keyword3); err != nil {
            return "", "", "", "", err
        }
    }
    if keyword != "" {
        keyword = "%" + strings.ToLower(keyword) + "%"
    }
    if keyword2 != "" {
        keyword2 = "%" + strings.ToLower(keyword2) + "%"
    }
    if keyword4 != "" && !strings.EqualFold(keyword4, "All") {
        keyword4 = "%" + strings.ToLower(keyword4) + "%"
    } else if strings.EqualFold(keyword4, "All") {
        keyword4 = ""
    }

    return keyword, keyword2, keyword3, keyword4, nil
}

func validateDateDDMMYYYY(value string) error {
    if _, err := time.Parse("02/01/2006", value); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Wrong date format")
    }
    return nil
}

func validateLogisticStatus(status string) error {
    if status != utils.LogisticRequestStatusConfirmed && status != utils.LogisticRequestStatusCancelled && status != utils.LogisticRequestStatusRejected {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Logistic Request Status")
    }
    return nil
}

func (cr *LogisticController) sendLogisticExportFile(c fiber.Ctx, list []logisticModel.LogisticRequest, filePrefix string) error {
    f := excelize.NewFile()
    defer f.Close()

    sheet := "LogisticRequests"
    defaultSheet := f.GetSheetName(0)
    f.SetSheetName(defaultSheet, sheet)

    headers := []string{
        "Request Number",
        "Request Status",
        "Requester PRN",
        "Requester Name",
        "Requester DOB",
        "Requester Doc Type",
        "Requester Doc Number",
        "Requester Nationality",
        "Requester Email",
        "Primary Doctor",
        "Visit With Companion",
        "Companion Name",
        "Companion DOB",
        "Companion Doc Type",
        "Companion Doc Number",
        "Relationship To Requester",
        "Flight Airline Name",
        "Flight Number",
        "Flight Arrival Date",
        "Flight Arrival Time",
        "Requested Pickup Date",
        "Requested Pickup Time",
        "Requested Pickup Day",
        "Date Create",
        "User Update",
        "User Date Update",
        "Admin Update",
        "Admin Date Update",
    }

    for i, h := range headers {
        cell, _ := excelize.CoordinatesToCellName(i+1, 1)
        if err := f.SetCellValue(sheet, cell, h); err != nil {
            return err
        }
    }

    for i, row := range list {
        r := i + 2
        values := []any{
            row.LogisticRequestNumber.String,
            row.LogisticRequestStatus.String,
            row.RequesterPrn.String,
            row.RequesterName.String,
            row.RequesterDob.String,
            row.RequesterDocType.String,
            row.RequesterDocNumber.String,
            row.RequesterNationality.String,
            row.RequesterEmail.String,
            row.PrimaryDoctor.String,
            row.VisitWithCompanion.String,
            row.CompanionName.String,
            row.CompanionDob.String,
            row.CompanionDocType.String,
            row.CompanionDocNumber.String,
            row.RelationshipToRequester.String,
            row.FlightAirlineName.String,
            row.FlightNumber.String,
            row.FlightArrivalDate.String,
            row.FlightArrivalTime.String,
            row.RequestedPickupDate.String,
            row.RequestedPickupTime.String,
            row.RequestedPickupDay.String,
            row.DateCreate.String,
            row.UserUpdate.String,
            row.UserDateUpdate.String,
            row.AdminUpdate.String,
            row.AdminDateUpdate.String,
        }

        for j, v := range values {
            cell, _ := excelize.CoordinatesToCellName(j+1, r)
            if err := f.SetCellValue(sheet, cell, v); err != nil {
                return err
            }
        }
    }

    lastCol, _ := excelize.ColumnNumberToName(len(headers))
    tableRange := fmt.Sprintf("A1:%s1", lastCol)
    style, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
    if err == nil {
        _ = f.SetCellStyle(sheet, "A1", tableRange, style)
    }

    for i := range headers {
        col, _ := excelize.ColumnNumberToName(i + 1)
        _ = f.SetColWidth(sheet, col, col, 20)
    }

    buf := bytes.NewBuffer(nil)
    if err := f.Write(buf); err != nil {
        return err
    }

    filename := fmt.Sprintf("%s_%s.xlsx", filePrefix, time.Now().Format("20060102_150405"))
    c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
    return c.Send(buf.Bytes())
}
