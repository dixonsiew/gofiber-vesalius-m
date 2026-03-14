package clubs

import (
    clubService "vesaliusm/service/clubs"

    "github.com/gofiber/fiber/v3"
)

type ClubsController struct {
    clubSvc *clubService.ClubService
}

func NewClubsController(clubSvc *clubService.ClubService) *ClubsController {
    return &ClubsController{
        clubSvc: clubSvc,
    }
}

// GetGoldenPearlAboutUs
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} clubs.GoldenPearlAboutUs
// @Router /clubs/goldenpearl/about-us [get]
func (cr *ClubsController) GetGoldenPearlAboutUs(c fiber.Ctx) error {
    o, err := cr.clubSvc.FindGoldenPearlAboutUs()
    if err != nil {
        return err
    }

    return c.JSON(o)
}
