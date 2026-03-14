package common

import (
    novaCountryService "vesaliusm/service/country"

    "github.com/gofiber/fiber/v3"
)

type CommonController struct {
    novaCountrySvc *novaCountryService.CountryService
}

func NewCommonController(novaCountrySvc *novaCountryService.CountryService) *CommonController {
    return &CommonController{
        novaCountrySvc: novaCountrySvc,
    }
}

// GetCountriesTelCode
//
// @Tags Common
// @Produce json
// @Success 200 {array} model.CountryTelCode
// @Router /common/telcode/list [get]
func (cr *CommonController) GetCountriesTelCode(c fiber.Ctx) error {
    lx, err := cr.novaCountrySvc.FindAllCountryTelCode()
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
    lx, err := cr.novaCountrySvc.FindAllCountries()
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
    lx, err := cr.novaCountrySvc.FindAllNationalities()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}
