package clubs

import (
    clubService "vesaliusm/service/club"

    "github.com/gofiber/fiber/v3"
)

// GetGoldenPearlAboutUs
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.GoldenPearlAboutUs
// @Router /clubs/goldenpearl/about-us [get]
func GetGoldenPearlAboutUs(c fiber.Ctx) error {
    o, err := clubService.FindGoldenPearlAboutUs()
    if err != nil {
        return err
    }

    return c.JSON(o)
}
