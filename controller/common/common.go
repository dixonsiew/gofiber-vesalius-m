package common

import (
    "vesaliusm/service/country"

    "github.com/gofiber/fiber/v3"
)

type CommonController struct {
    novaCountryService *country.CountryService
}

func NewCommonController() *CommonController {
    return &CommonController{
        novaCountryService: country.CountrySvc,
    }
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
