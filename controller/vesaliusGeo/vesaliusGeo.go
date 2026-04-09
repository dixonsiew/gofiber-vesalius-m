package vesaliusGeo

import (
    "vesaliusm/service/vesaliusGeo"

    "github.com/gofiber/fiber/v3"
)

type VesaliusGeoController struct {
    vesaliusGeoService *vesaliusGeo.VesaliusGeoService
}

func NewVesaliusGeoController() *VesaliusGeoController {
    return &VesaliusGeoController{
        vesaliusGeoService: vesaliusGeo.VesaliusGeoSvc,
    }
}

// Login
//
// @Tags Vesalius Geo
// @Produce json
// @Success 200
// @Router /vesalius-geo/authenticate [get]
func (cr *VesaliusGeoController) Login(c fiber.Ctx) error {
    token, _, ex := cr.vesaliusGeoService.Login()
    if ex != nil {
        return fiber.NewError(fiber.StatusBadRequest, ex.ToString())
    }
    return c.JSON(fiber.Map{
        "token": token,
    })
}

// Logout
//
// @Tags Vesalius Geo
// @Produce json
// @Param token path string true "token"
// @Success 200
// @Router /vesalius-geo/logout/{token} [get]
func (cr *VesaliusGeoController) Logout(c fiber.Ctx) error {
    token := c.Params("token")
    res, ex := cr.vesaliusGeoService.Logout(token)
    if ex != nil {
        return fiber.NewError(fiber.StatusBadRequest, ex.ToString())
    }
    return c.JSON(res)
}
