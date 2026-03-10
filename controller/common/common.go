package common

import (
	novaCountryService "vesaliusm/service/country"
    "vesaliusm/database"

	"github.com/gofiber/fiber/v3"
)

var novaCountrySvc *novaCountryService.CountryService = 
    novaCountryService.NewCountryService(database.GetDb(), database.GetCtx())

// GetCountriesTelCode
//
// @Tags Common
// @Produce json
// @Success 200 {array} model.CountryTelCode
// @Router /common/telcode/list [get]
func GetCountriesTelCode(c fiber.Ctx) error {
    lx, err := novaCountrySvc.FindAllCountryTelCode()
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
func GetCountries(c fiber.Ctx) error {
    lx, err := novaCountrySvc.FindAllCountries()
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
func GetNationalities(c fiber.Ctx) error {
    lx, err := novaCountrySvc.FindAllNationalities()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}
