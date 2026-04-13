package common

import (
    "strconv"
    "vesaliusm/dto"
    "vesaliusm/model"
    "vesaliusm/service/app"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/country"
    "vesaliusm/service/requestDeleteAccount"
    "vesaliusm/utils"
    "vesaliusm/utils/constants"

    "github.com/gofiber/fiber/v3"
)

type CommonController struct {
    appService                  *app.AppService
    applicationUserService      *applicationUser.ApplicationUserService
    novaCountryService          *country.CountryService
    requestDeleteAccountService *requestDeleteAccount.RequestDeleteAccountService
}

func NewCommonController() *CommonController {
    return &CommonController{
        appService:                  app.AppSvc,
        applicationUserService:      applicationUser.ApplicationUserSvc,
        novaCountryService:          country.CountrySvc,
        requestDeleteAccountService: requestDeleteAccount.RequestDeleteAccountSvc,
    }
}

// GetAppHospitalProfile
//
// @Tags Common
// @Produce json
// @Success 200 {array} model.HospitalProfile
// @Router /common/app/hospital-profile [get]
func (cr *CommonController) GetAppHospitalProfile(c fiber.Ctx) error {
    lx, err := cr.appService.FindAllAppHospitalProfile()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetAppVersion
//
// @Tags Common
// @Produce json
// @Success 200 {array} model.AppVersion
// @Router /common/app/version [get]
func (cr *CommonController) GetAppVersion(c fiber.Ctx) error {
    lx, err := cr.appService.FindAllAppVersion()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// UpdateAppVersion
//
// @Tags Common
// @Accept json
// @Produce json
// @Param request body dto.AppVersionDto true "AppVersionDto"
// @Success 200
// @Router /common/app/version [post]
func (cr *CommonController) UpdateAppVersion(c fiber.Ctx) error {
    data := new(dto.AppVersionDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    if data.OsPlatform != "Android" && data.OsPlatform != "IOS" {
        return fiber.NewError(fiber.StatusBadRequest, "Incorrect OS Platform")
    }

    err := cr.appService.UpdateAppVersion(data.LatestVersion, data.OsPlatform, data.Status)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "App Version Updated",
    })
}

// GetReleaseVersion
//
// @Tags Common
// @Produce json
// @Success 200 {array} model.ReleaseVersion
// @Router /common/app/release/version [get]
func (cr *CommonController) GetReleaseVersion(c fiber.Ctx) error {
    lx, err := cr.appService.FindAllReleaseVersion()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// UpdateReleaseVersion
//
// @Tags Common
// @Accept json
// @Produce json
// @Param request body dto.ReleaseVersionDto true "ReleaseVersionDto"
// @Success 200
// @Router /common/app/release/version [post]
func (cr *CommonController) UpdateReleaseVersion(c fiber.Ctx) error {
    data := new(dto.ReleaseVersionDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    if data.StackPlatform != "NestJS Backend" && data.StackPlatform != "NestJS Cronjob" && data.StackPlatform != "Angular WebAdmin" {
        return fiber.NewError(fiber.StatusBadRequest, "Incorrect Stack Platform")
    }

    currentVersion, err := cr.appService.GetCurrentReleaseVersion(data.StackPlatform)
    if err != nil {
        return err
    }

    if data.LatestVersion == currentVersion {
        return fiber.NewError(fiber.StatusBadRequest, "Latest version cannot be same as Current version")
    }

    err = cr.appService.UpdateReleaseVersion(data.LatestVersion, data.StackPlatform)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Release Version Updated",
    })
}

// RequestDeleteAccount
//
// @Tags Common
// @Accept json
// @Produce json
// @Param request body dto.DeleteAccountDto true "DeleteAccountDto"
// @Success 200
// @Router /common/request-delete-account [post]
func (cr *CommonController) RequestDeleteAccount(c fiber.Ctx) error {
    data := new(dto.DeleteAccountDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o := model.RequestDeleteAccount{
        PRN: data.PRN,
        Fullname: data.FullName,
        DocumentNumber: data.DocumentNumber,
        DOB: data.DOB,
        ContactNumber: data.ContactNumber,
        Email: data.Email,
    }
    err := cr.requestDeleteAccountService.SaveRequest(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Account Delete Requested",
    })
}

// GetGuestModeServices
//
// @Tags Common
// @Produce json
// @Success 200 {array} model.AppServices
// @Router /common/service/guest/list [get]
func (cr *CommonController) GetGuestModeServices(c fiber.Ctx) error {
    m, err := cr.appService.ListByGuestMode()
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    return c.JSON(m.List)
}

// GetAuthModeServices
//
// @Tags Common
// @Produce json
// @Success 200 {array} model.AppServices
// @Router /common/service/auth/list [get]
func (cr *CommonController) GetAuthModeServices(c fiber.Ctx) error {
    m, err := cr.appService.ListByAuthMode()
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    return c.JSON(m.List)
}

// GetCountriesTelCode
//
// @Tags Common
// @Produce json
// @Success 200 {array} model.CountryTelCode
// @Router /common/telcode/list [get]
func (cr *CommonController) GetCountriesTelCode(c fiber.Ctx) error {
    lx, err := cr.novaCountryService.FindAllCountryTelCode()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetCountries
//
// @Tags Common
// @Produce json
// @Success 200 {array} model.Country
// @Router /common/country/list [get]
func (cr *CommonController) GetCountries(c fiber.Ctx) error {
    lx, err := cr.novaCountryService.FindAllCountries()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetNationalities
//
// @Tags Common
// @Produce json
// @Success 200 {array} model.Nationality
// @Router /common/nationality/list [get]
func (cr *CommonController) GetNationalities(c fiber.Ctx) error {
    lx, err := cr.novaCountryService.FindAllNationalities()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// UserDownloadedApp
//
// @Tags Common
// @Produce json
// @Param playerId path string true "playerId"
// @Success 200
// @Router /common/downloaded-app/{playerId} [post]
func (cr *CommonController) UserDownloadedApp(c fiber.Ctx) error {
    playerId := c.Params("playerId")
    err := cr.applicationUserService.InsertDownloadApp(playerId, nil)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "PlayerID successfully captured",
    })
}

// UserDownloadedAppV2
//
// @Tags Common
// @Produce json
// @Param request body dto.AppDownloadedUserDto true "AppDownloadedUserDto"
// @Success 200
// @Router /common/downloaded-app-v2 [post]
func (cr *CommonController) UserDownloadedAppV2(c fiber.Ctx) error {
    data := new(dto.AppDownloadedUserDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    err := cr.applicationUserService.InsertDownloadAppV2(data.MachineId, data.PlayerId, nil)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "PlayerID successfully captured",
    })
} 
