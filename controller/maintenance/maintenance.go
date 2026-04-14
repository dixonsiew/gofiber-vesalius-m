package maintenance

import (
    "strconv"
    // "vesaliusm/model"
    "vesaliusm/dto"
    "vesaliusm/service/exportExcel"
    "vesaliusm/service/maintenance"
    "vesaliusm/utils"
    "vesaliusm/utils/constants"

    "github.com/gofiber/fiber/v3"
    "github.com/nleeper/goment"
)

type MaintenanceController struct {
    maintenanceService *maintenance.MaintenanceService
    exportExcelService *exportExcel.ExportExcelService
}

func NewMaintenanceController() *MaintenanceController {
    return &MaintenanceController{
        maintenanceService: maintenance.MaintenanceSvc,
        exportExcelService: exportExcel.ExportExcelSvc,
    }
}

// GetAllHospitalProfiles
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int  false  "_page"  default:"1"
// @Param        _limit            query       int  false  "_limit" default:"10"
// @Success 200 {array} model.HospitalProfile
// @Router /maintenance/hospital-profile/all [get]
func (cr *MaintenanceController) GetAllHospitalProfiles(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.maintenanceService.ListAllHospitalProfiles(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllHospitalProfiles
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int                   false  "_page"  default:"1"
// @Param        _limit            query       int                   false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.HospitalProfile
// @Router /maintenance/hospital-profile/search [post]
func (cr *MaintenanceController) SearchAllHospitalProfiles(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    if key != "" {
        key = "%" + key + "%"
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.maintenanceService.SearchAllHospitalProfilesByKeyword(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// UpdateHospitalProfile
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        request       body        dto.HospitalProfileDto  true  "HospitalProfileDto"
// @Success 200
// @Router /maintenance/hospital-profile/update [post]
func (cr *MaintenanceController) UpdateHospitalProfile(c fiber.Ctx) error {
    data := new(dto.HospitalProfileDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    err := cr.maintenanceService.UpdateHospitalProfileByDescName(data)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Hospital Profile Updated",
    })
}

// GetAllParamSettings
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int  false  "_page"  default:"1"
// @Param        _limit            query       int  false  "_limit" default:"10"
// @Success 200 {array} model.ParamSetting
// @Router /maintenance/param-setting/all [get]
func (cr *MaintenanceController) GetAllParamSettings(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.maintenanceService.ListAllParamSettings(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllParamSettings
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int                   false  "_page"  default:"1"
// @Param        _limit            query       int                   false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.ParamSetting
// @Router /maintenance/param-setting/search [post]
func (cr *MaintenanceController) SearchAllParamSettings(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    if key != "" {
        key = "%" + key + "%"
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.maintenanceService.SearchAllParamSettingsByKeyword(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// UpdateParamSetting
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        request       body        dto.ParamSettingDto  true  "ParamSettingDto"
// @Success 200
// @Router /maintenance/param-setting/update [post]
func (cr *MaintenanceController) UpdateParamSetting(c fiber.Ctx) error {
    data := new(dto.ParamSettingDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    err := cr.maintenanceService.UpdateParamSettingByParamCode(data)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Param Setting Updated",
    })
}

// GetAllNotificationSettings
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int  false  "_page"  default:"1"
// @Param        _limit            query       int  false  "_limit" default:"10"
// @Success 200 {array} model.NotificationSetting
// @Router /maintenance/notification-setting/all [get]
func (cr *MaintenanceController) GetAllNotificationSettings(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.maintenanceService.ListAllNotificationSettings(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllNotificationSettings
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int                   false  "_page"  default:"1"
// @Param        _limit            query       int                   false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.NotificationSetting
// @Router /maintenance/notification-setting/search [post]
func (cr *MaintenanceController) SearchAllNotificationSettings(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    if key != "" {
        key = "%" + key + "%"
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.maintenanceService.SearchAllNotificationSettingsByKeyword(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// UpdateNotificationSetting
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        request        body        dto.NotificationSettingDto  true  "NotificationSettingDto"
// @Success 200
// @Router /maintenance/notification-setting/update [post]
func (cr *MaintenanceController) UpdateNotificationSetting(c fiber.Ctx) error {
    data := new(dto.NotificationSettingDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    err := cr.maintenanceService.UpdateNotificationSettingByNotificationCode(data)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Notification Setting Updated",
    })
}

// GetAllCronjobHistories
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int  false  "_page"  default:"1"
// @Param        _limit            query       int  false  "_limit" default:"10"
// @Success 200 {array} model.CronjobHistory
// @Router /maintenance/cronjob-history/all [get]
func (cr *MaintenanceController) GetAllCronjobHistories(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.maintenanceService.ListAllCronjobHistories(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllCronjobHistories
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        cronjobName       path        string            true   "cronjobName"
// @Param        _page             query       int               false  "_page"  default:"1"
// @Param        _limit            query       int               false  "_limit" default:"10"
// @Param         keyword      body        dto.SearchKeyword2Dto false  "Search"
// @Success 200 {array} model.CronjobHistory
// @Router /maintenance/cronjob-history/search/{cronjobName} [post]
func (cr *MaintenanceController) SearchAllCronjobHistories(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key2 := data.GetString("keyword2")
    if key != "" {
        if _, err := goment.New(key, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong start date format")
        }
    }
    if key2 != "" {
        if _, err := goment.New(key2, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong end date format")
        }
    }

    cronjobName := c.Params("cronjobName")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.maintenanceService.SearchAllCronjobHistoriesByKeyword(cronjobName, key, key2, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllDynamicEmailSettings
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int  false  "_page"  default:"1"
// @Param        _limit            query       int  false  "_limit" default:"10"
// @Success 200 {array} model.DynamicEmailMaster
// @Router /maintenance/dynamic-email-master/all [get]
func (cr *MaintenanceController) GetAllDynamicEmailSettings(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.maintenanceService.ListAllDynamicEmailSettings(page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllDynamicEmailSettings
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        _page             query       int                   false  "_page"  default:"1"
// @Param        _limit            query       int                   false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.DynamicEmailMaster
// @Router /maintenance/dynamic-email-master/search [post]
func (cr *MaintenanceController) SearchAllDynamicEmailSettings(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    if key != "" {
        key = "%" + key + "%"
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.maintenanceService.SearchAllDynamicEmailSettingsByKeyword(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllExportDynamicEmailSettings
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.DynamicEmailMaster
// @Router /maintenance/dynamic-email-master/export/all [get]
func (cr *MaintenanceController) GetAllExportDynamicEmailSettings(c fiber.Ctx) error {
    lx, err := cr.exportExcelService.ExportAllDynamicEmailSettings()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetSearchDynamicEmailSettings
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.DynamicEmailMaster
// @Router /maintenance/dynamic-email-master/export/search [post]
func (cr *MaintenanceController) GetSearchDynamicEmailSettings(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    if key != "" {
        key = "%" + key + "%"
    }

    lx, err := cr.exportExcelService.ExportDynamicEmailSettingsByKeyword(key)
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetDynamicEmailSettingByFunctionName
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        emailFunctionName path string true "emailFunctionName"
// @Success 200 {object} model.DynamicEmailMaster
// @Router /maintenance/dynamic-email-master/view/{emailFunctionName} [get]
func (cr *MaintenanceController) GetDynamicEmailSettingByFunctionName(c fiber.Ctx) error {
    emailFunctionName := c.Params("emailFunctionName")
    o, err := cr.maintenanceService.ViewDynamicEmailSettingByFunctionName(emailFunctionName)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// UpdateDynamicEmailSetting
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Param        request        body      dto.DynamicEmailMasterDto true "DynamicEmailMasterDto"
// @Success 200
// @Router /maintenance/dynamic-email-master/update [post]
func (cr *MaintenanceController) UpdateDynamicEmailSetting(c fiber.Ctx) error {
    data := new(dto.DynamicEmailMasterDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    err := cr.maintenanceService.UpdateDynamicEmailSettingByFunctionName(data)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Email Setting Updated",
    })
}

// GetAllStatisticAppointments
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.StatisticAppointment
// @Router /maintenance/statistic/appointment/all [get]
func (cr *MaintenanceController) GetAllStatisticAppointments(c fiber.Ctx) error {
    lx, err := cr.maintenanceService.GetAllStatisticAppointments()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetAllStatisticMobileRegistrations
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.StatisticMobileRegistration
// @Router /maintenance/statistic/mobile/registration/all [get]
func (cr *MaintenanceController) GetAllStatisticMobileRegistrations(c fiber.Ctx) error {
    lx, err := cr.maintenanceService.GetAllStatisticMobileRegistrations()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetAllStatisticMobileFeedbacks
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.StatisticMobileFeedback
// @Router /maintenance/statistic/mobile/feedback/all [get]
func (cr *MaintenanceController) GetAllStatisticMobileFeedbacks(c fiber.Ctx) error {
    lx, err := cr.maintenanceService.GetAllStatisticMobileFeedbacks()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetAllStatisticMobilePackages
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.StatisticMobilePackage
// @Router /maintenance/statistic/mobile/package/all [get]
func (cr *MaintenanceController) GetAllStatisticMobilePackages(c fiber.Ctx) error {
    lx, err := cr.maintenanceService.GetAllStatisticMobilePackages()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetAllStatisticMobileClubsKids
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.StatisticMobileClubs
// @Router /maintenance/statistic/mobile/clubs/kids [get]
func (cr *MaintenanceController) GetAllStatisticMobileClubsKids(c fiber.Ctx) error {
    lx, err := cr.maintenanceService.GetAllStatisticMobileClubsKids()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetAllStatisticMobileClubsGoldenPearl
//
// @Tags Maintenance
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.StatisticMobileClubs
// @Router /maintenance/statistic/mobile/clubs/goldenpearl [get]
func (cr *MaintenanceController) GetAllStatisticMobileClubsGoldenPearl(c fiber.Ctx) error {
    lx, err := cr.maintenanceService.GetAllStatisticMobileClubsGoldenPearl()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}
