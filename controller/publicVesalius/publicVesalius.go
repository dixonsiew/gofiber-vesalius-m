package publicVesalius

import (
    "strconv"
    "vesaliusm/dto"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/novaDoctor"
    "vesaliusm/service/vesalius"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

type PublicVesaliusController struct {
    applicationUserService *applicationUser.ApplicationUserService
    novaDoctorService      *novaDoctor.NovaDoctorService
    vesaliusService        *vesalius.VesaliusService
}

func NewPublicVesaliusController() *PublicVesaliusController {
    return &PublicVesaliusController{
        applicationUserService: applicationUser.ApplicationUserSvc,
        novaDoctorService:      novaDoctor.NovaDoctorSvc,
        vesaliusService:        vesalius.VesaliusSvc,
    }
}

// GetPatientData
//
// @Tags Public Vesalius
// @Produce json
// @Param        prn             path       string  true  "prn"
// @Success 200 {array} vesaliusGeo.Patient
// @Router /patient-data/:branchId/:prn [get]
func (cr *PublicVesaliusController) GetPatientData(c fiber.Ctx) error {
    prn := c.Params("prn")
    o, _, err := cr.vesaliusService.VesaliusGetPatientDataByPrn(prn)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// GetNextSessionAvailableSlots
//
// @Tags Public Vesalius
// @Produce json
// @Param        prn          path       string                        true  "prn"
// @Param        request      body       dto.PostNextAvailableSlotsDto true  "PostNextAvailableSlotsDto"
// @Success 200 {array} vesaliusGeo.Slot
// @Router /get-next-session-available-slots/:branchId/:prn [post]
func (cr *PublicVesaliusController) GetNextSessionAvailableSlots(c fiber.Ctx) error {
    data := new(dto.PostNextAvailableSlotsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    prn := c.Params("prn")
    lx, err := cr.vesaliusService.VesaliusGetNextSessionAvailableSlots(prn, data)
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetNextAvailableSlots
//
// @Tags Public Vesalius
// @Produce json
// @Param        prn          path       string                        true  "prn"
// @Param        request      body       dto.PostNextAvailableSlotsDto true  "PostNextAvailableSlotsDto"
// @Success 200 {array} vesaliusGeo.Slot
// @Router /get-next-available-slots/:branchId/:prn [post]
func (cr *PublicVesaliusController) GetNextAvailableSlots(c fiber.Ctx) error {
    data := new(dto.PostNextAvailableSlotsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    prn := c.Params("prn")
    lx, err := cr.vesaliusService.VesaliusGetNextAvailableSlots(prn, data)
    if err != nil {
        return err
    }
    
    return c.JSON(lx)
}

// GetMakeAppointment
//
// @Tags Public Vesalius
// @Produce json
// @Param        prn          path       string                        true  "prn"
// @Param        request      body       dto.PostMakeAppointmentDto    true  "PostMakeAppointmentDto"
// @Success 200 {object} dto.PostMakeAppointmentDto
// @Router /make-appointment/:branchId/:prn [post]
func (cr *PublicVesaliusController) GetMakeAppointment(c fiber.Ctx) error {
    data := new(dto.PostMakeAppointmentDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }
    
    prn := c.Params("prn")
    o, err := cr.vesaliusService.VesaliusGetMakeAppointment(prn, data)
    if err != nil {
        return err
    }
    
    return c.JSON(o)
}

// GetAllDoctorInformation
//
// @Tags Public Vesalius
// @Produce json
// @Param        branchId     path       string                        true  "branchId"
// @Param        webadmin     path       string                        true  "webadmin"
// @Success 200 {array} model.NovaDoctor
// @Router /getAllDoctorInformation/:branchId/:webadmin [get]
func (cr *PublicVesaliusController) GetAllDoctorInformation(c fiber.Ctx) error {
    webadmin := c.Params("webadmin")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.novaDoctorService.List(page, limit, webadmin == "1")
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllHSDoctorInformation
//
// @Tags Public Vesalius
// @Produce json
// @Param        branchId     path       string                        true  "branchId"
// @Success 200 {array} model.NovaDoctor
// @Router /getAllHSDoctorInformation/:branchId [get]
func (cr *PublicVesaliusController) GetAllHSDoctorInformation(c fiber.Ctx) error {
    lx, err := cr.novaDoctorService.FindAllHSMcrAndName()
    if err != nil {
        return err
    }
    
    return c.JSON(lx)
}

// SearchAllDoctorInformation
//
// @Tags Public Vesalius
// @Produce json
// @Param        branchId          path        string                        true  "branchId"
// @Param        _page             query       string                        false  "_page"  default:"1"
// @Param        _limit            query       string                        false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto          false  "Search"
// @Success 200 {array} model.NovaDoctor
// @Router /getAllDoctorInformation/:branchId [post]
func (cr *PublicVesaliusController) SearchAllDoctorInformation(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    if key != "" {
        key = "%" + key + "%"
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.novaDoctorService.ListByKeyword(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetDoctorInformationByMCR
//
// @Tags Public Vesalius
// @Produce json
// @Param        mcr     path       string                        true  "mcr"
// @Success 200 {array} model.NovaDoctor
// @Router /getDoctorInformationByMCR/:branchId/:mcr [get]
func (cr *PublicVesaliusController) GetDoctorInformationByMCR(c fiber.Ctx) error {
    mcr := c.Params("mcr")
    lx, err := cr.novaDoctorService.FindAllByMcr(mcr)
    if err != nil {
        return err
    }
    
    return c.JSON(lx)
}

// GetSpecialtyLookup
//
// @Tags Public Vesalius
// @Produce json
// @Success 200 {array} model.NovaSpecialty
// @Router /lookup/specialty [get]
func (cr *PublicVesaliusController) GetSpecialtyLookup(c fiber.Ctx) error {
    lx, err := cr.novaDoctorService.FindAllNovaSpecialties()
    if err != nil {
        return err
    }
    
    return c.JSON(lx)
}

// CheckFileSize
//
// @Tags Public Vesalius
// @Produce json
// @Param        size     path       string                        true  "size"
// @Success 200
// @Router /check-file-size/:size [get]
func (cr *PublicVesaliusController) CheckFileSize(c fiber.Ctx) error {
    size := c.Params("size")
    isize, _ := strconv.ParseInt(size, 10, 64)
    if isize > 5242880 {
        return fiber.NewError(fiber.StatusBadRequest, "File size limit exceeded 5 MB")
    }

    return c.JSON(fiber.Map{
        "success": 1,
    })
}
