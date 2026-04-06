package vesalius

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "image"
    "image/png"
    "slices"
    "strconv"
    "strings"
    "vesaliusm/dto"
    "vesaliusm/model"
    gm "vesaliusm/model/vesaliusGeo"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/novaDoctor"
    "vesaliusm/service/novaDoctorPatientAppt"
    "vesaliusm/service/vesalius"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/nfnt/resize"
    "github.com/nleeper/goment"
)

type VesaliusController struct {
    applicationUserService       applicationUser.ApplicationUserService
    novaDoctorService            *novaDoctor.NovaDoctorService
    novaDoctorPatientApptService *novaDoctorPatientAppt.NovaDoctorPatientApptService
    vesaliusService              *vesalius.VesaliusService
}

func NewVesaliusController() *VesaliusController {
    return &VesaliusController{
        applicationUserService:       *applicationUser.ApplicationUserSvc,
        novaDoctorService:            novaDoctor.NovaDoctorSvc,
        novaDoctorPatientApptService: novaDoctorPatientAppt.NovaDoctorPatientApptSvc,
        vesaliusService:              vesalius.VesaliusSvc,
    }
}

// ProcessResizeImage
//
// @Tags Vesalius
// @Produce text/plain
// @Success 200
// @Router /vesalius/process-resize-image [post]
func (cr *VesaliusController) ProcessResizeImage(c fiber.Ctx) error {
    lx, err := cr.novaDoctorService.FindAll(0, 500, true)
    if err != nil {
        return err
    }

    for i := range lx {
        if lx[i].Image.Valid {
            resized, err := resizeBase64Image(lx[i].Image.String)
            if err != nil {
                return err
            }

            cr.novaDoctorService.ResizeAllDoctorImage(resized, lx[i].DoctorId.Int64)
        }
    }

    return c.SendString("done")
}

// GetAllAppointments
//
// @Tags Vesalius
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.DoctorPatientAppointment
// @Router /vesalius/doctor/appointment/all [get]
func (cr *VesaliusController) GetAllAppointments(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.novaDoctorPatientApptService.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllAppointments
//
// @Tags Vesalius
// @Produce json
// @Security BearerAuth
// @Param         _page        query       string                false  "_page"  default:"1"
// @Param         _limit       query       string                false  "_limit" default:"10"
// @Param         keyword      body        dto.SearchKeyword3Dto false  "Search"
// @Success 200 {array} model.DoctorPatientAppointment
// @Router /vesalius/doctor/appointment/all [post]
func (cr *VesaliusController) SearchAllAppointments(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key2 := data.GetString("keyword2")
    key3 := data.GetString("keyword3")
    if key != "" {
        key = "%" + key + "%"
    }
    if key2 != "" {
        key2 = "%" + key2 + "%"
    }
    if key3 != "" {
        if _, err := goment.New(key3, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong date format")
        }
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.novaDoctorPatientApptService.ListByKeyword(key, key2, key3, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// CreateDoctor
//
// @Tags Vesalius
// @Produce json
// @Security BearerAuth
// @Param      request     body     dto.NovaDoctorDto     true  "NovaDoctorDto"
// @Success 200
// @Router /vesalius/doctor [post]
func (cr *VesaliusController) CreateDoctor(c fiber.Ctx) error {
    data := new(dto.NovaDoctorDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.novaDoctorService.ExistsByMcr(data.Mcr)
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, "A doctor with that mcr already exists")
    }

    resizeValidation, err := resizeBase64Image(data.Image)
    if err != nil {
        return err
    }

    o := &model.NovaDoctor{
        Gender:           utils.NewNullString(data.Gender),
        MCR:              utils.NewNullString(data.Mcr),
        Name:             utils.NewNullString(data.Name),
        Nationality:      utils.NewNullString(data.Nationality),
        Qualifications:   utils.NewNullString(data.Qualifications),
        RegistrationNum:  utils.NewNullString(data.RegistrationNum),
        DisplaySequence:  utils.NewInt32(int32(data.DisplaySequence)),
        AllowAppointment: utils.NewNullString(data.AllowAppointment),
        ConsultantType:   utils.NewNullString(data.ConsultantType),
        IsForPackage:     utils.NewNullString(data.IsForPackage),
        Image:            utils.NewNullString(data.Image),
        ResizeImage:      resizeValidation,
    }

    o.DoctorClinicHours = make([]model.NovaDoctorClinicHours, len(data.DoctorClinicHours))
    for i, v := range data.DoctorClinicHours {
        o.DoctorClinicHours[i] = v.ToDbModel()
        o.DoctorClinicHours[i].Set()
    }

    o.DoctorClinicLocation = make([]model.NovaDoctorClinicLocation, len(data.DoctorClinicLocation))
    for i, v := range data.DoctorClinicLocation {
        o.DoctorClinicLocation[i] = v.ToDbModel()
    }

    o.DoctorAppointment = make([]model.NovaDoctorAppointment, len(data.DoctorAppointment))
    for i, v := range data.DoctorAppointment {
        o.DoctorAppointment[i] = v.ToDbModel()
    }

    o.DoctorContact = make([]model.NovaDoctorContact, len(data.DoctorContact))
    for i, v := range data.DoctorContact {
        o.DoctorContact[i] = v.ToDbModel()
    }

    o.DoctorQualifications = make([]model.NovaDoctorQualifications, len(data.DoctorQualifications))
    for i, v := range data.DoctorQualifications {
        o.DoctorQualifications[i] = v.ToDbModel()
    }

    o.DoctorSpecialities = make([]model.NovaDoctorSpecialities, len(data.DoctorSpecialities))
    for i, v := range data.DoctorSpecialities {
        o.DoctorSpecialities[i] = v.ToDbModel()
    }

    o.DoctorSpecialty = make([]model.NovaDoctorSpecialty, len(data.DoctorSpecialty))
    for i, v := range data.DoctorSpecialty {
        o.DoctorSpecialty[i] = v.ToDbModel()
        o.DoctorSpecialty[i].Set()
    }

    if o.AllowAppointment.String == "Y" {
        if len(o.DoctorSpecialty) > 0 {
            hasMainSpecialty := slices.ContainsFunc(o.DoctorSpecialty, func(x model.NovaDoctorSpecialty) bool {
                return x.PrimarySpecialty == true
            })
            if !hasMainSpecialty {
                return fiber.NewError(fiber.StatusBadRequest, "Doctor must have at least one Main Specialty")
            }
        } else {
            return fiber.NewError(fiber.StatusBadRequest, "Doctor must have at least one Specialty")
        }
    }

    o.DoctorSpokenLanguage = make([]model.NovaDoctorSpokenLanguage, len(data.DoctorSpokenLanguage))
    for i, v := range data.DoctorSpokenLanguage {
        o.DoctorSpokenLanguage[i] = v.ToDbModel()
    }

    err = cr.novaDoctorService.Save(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "success": 1,
    })
}

// GetDoctorInformationByDoctorId
//
// @Tags Vesalius
// @Produce json
// @Security BearerAuth
// @Param       doctorId       path      string  true  "doctorId"
// @Success 200
// @Router /vesalius/doctor/{doctorId} [get]
func (cr *VesaliusController) GetDoctorInformationByDoctorId(c fiber.Ctx) error {
    doctorId := c.Params("doctorId")
    idoctorId, _ := strconv.ParseInt(doctorId, 10, 64)
    o, err := cr.novaDoctorService.FindAllByDoctorId(idoctorId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// UpdateDoctor
//
// @Tags Vesalius
// @Produce json
// @Security BearerAuth
// @Param       doctorId       path      string             true  "doctorId"
// @Param       request        body      dto.NovaDoctorDto  true  "NovaDoctorDto"
// @Success 200
// @Router /vesalius/doctor/{doctorId} [put]
func (cr *VesaliusController) UpdateDoctor(c fiber.Ctx) error {
    doctorId := c.Params("doctorId")
    idoctorId, _ := strconv.ParseInt(doctorId, 10, 64)
    data := new(dto.NovaDoctorDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.novaDoctorService.ExistsByOtherMcr(data.Mcr, int(idoctorId))
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, "A doctor with that mcr already exists")
    }

    resizeValidation, err := resizeBase64Image(data.Image)
    if err != nil {
        return err
    }

    o := &model.NovaDoctor{
        DoctorId:         utils.NewInt64(idoctorId),
        Gender:           utils.NewNullString(data.Gender),
        MCR:              utils.NewNullString(data.Mcr),
        Name:             utils.NewNullString(data.Name),
        Nationality:      utils.NewNullString(data.Nationality),
        Qualifications:   utils.NewNullString(data.Qualifications),
        RegistrationNum:  utils.NewNullString(data.RegistrationNum),
        DisplaySequence:  utils.NewInt32(int32(data.DisplaySequence)),
        AllowAppointment: utils.NewNullString(data.AllowAppointment),
        ConsultantType:   utils.NewNullString(data.ConsultantType),
        IsForPackage:     utils.NewNullString(data.IsForPackage),
        Image:            utils.NewNullString(data.Image),
        ResizeImage:      resizeValidation,
    }

    o.DoctorClinicHours = make([]model.NovaDoctorClinicHours, len(data.DoctorClinicHours))
    for i, v := range data.DoctorClinicHours {
        o.DoctorClinicHours[i] = v.ToDbModel()
        o.DoctorClinicHours[i].Set()
    }

    o.DoctorClinicLocation = make([]model.NovaDoctorClinicLocation, len(data.DoctorClinicLocation))
    for i, v := range data.DoctorClinicLocation {
        o.DoctorClinicLocation[i] = v.ToDbModel()
    }

    o.DoctorAppointment = make([]model.NovaDoctorAppointment, len(data.DoctorAppointment))
    for i, v := range data.DoctorAppointment {
        o.DoctorAppointment[i] = v.ToDbModel()
    }

    o.DoctorContact = make([]model.NovaDoctorContact, len(data.DoctorContact))
    for i, v := range data.DoctorContact {
        o.DoctorContact[i] = v.ToDbModel()
    }

    o.DoctorQualifications = make([]model.NovaDoctorQualifications, len(data.DoctorQualifications))
    for i, v := range data.DoctorQualifications {
        o.DoctorQualifications[i] = v.ToDbModel()
    }

    o.DoctorSpecialities = make([]model.NovaDoctorSpecialities, len(data.DoctorSpecialities))
    for i, v := range data.DoctorSpecialities {
        o.DoctorSpecialities[i] = v.ToDbModel()
    }

    o.DoctorSpecialty = make([]model.NovaDoctorSpecialty, len(data.DoctorSpecialty))
    for i, v := range data.DoctorSpecialty {
        o.DoctorSpecialty[i] = v.ToDbModel()
        o.DoctorSpecialty[i].Set()
    }

    if o.AllowAppointment.String == "Y" {
        if len(o.DoctorSpecialty) > 0 {
            hasMainSpecialty := slices.ContainsFunc(o.DoctorSpecialty, func(x model.NovaDoctorSpecialty) bool {
                return x.PrimarySpecialty == true
            })
            if !hasMainSpecialty {
                return fiber.NewError(fiber.StatusBadRequest, "Doctor must have at least one Main Specialty")
            }
        } else {
            return fiber.NewError(fiber.StatusBadRequest, "Doctor must have at least one Specialty")
        }
    }

    o.DoctorSpokenLanguage = make([]model.NovaDoctorSpokenLanguage, len(data.DoctorSpokenLanguage))
    for i, v := range data.DoctorSpokenLanguage {
        o.DoctorSpokenLanguage[i] = v.ToDbModel()
    }

    err = cr.novaDoctorService.Update(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "success": 1,
    })
}

// RemoveDoctorImage
//
// @Tags Vesalius
// @Produce json
// @Security BearerAuth
// @Param       doctorId       path      string       true  "doctorId"
// @Success 200
// @Router /vesalius/doctor-image-delete/{doctorId} [put]
func (cr *VesaliusController) RemoveDoctorImage(c fiber.Ctx) error {
    doctorId := c.Params("doctorId")
    idoctorId, _ := strconv.ParseInt(doctorId, 10, 64)
    err := cr.novaDoctorService.DeleteImageById(idoctorId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "success": 1,
    })
}

// Remove
//
// @Tags Vesalius
// @Produce json
// @Security BearerAuth
// @Param       doctorId       path      string       true  "doctorId"
// @Success 200
// @Router /vesalius/doctor/{doctorId} [delete]
func (cr *VesaliusController) Remove(c fiber.Ctx) error {
    doctorId := c.Params("doctorId")
    idoctorId, _ := strconv.ParseInt(doctorId, 10, 64)
    err := cr.novaDoctorService.DeleteByDoctorId(idoctorId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "success": 1,
    })
}

// GetDoctorAppointments
//
// @Tags Vesalius
// @Produce json
// @Security BearerAuth
// @Param       doctorId       path      string       true  "doctorId"
// @Param       month          path      string       true  "month"
// @Param       year           path      string       true  "year"
// @Param       needAppt       path      string       true  "needAppt"
// @Success 200
// @Router /vesalius/get-doctor-appointments/{doctorId}/{month}/{year}/{needAppt} [get]
func (cr *VesaliusController) GetDoctorAppointments(c fiber.Ctx) error {
    doctorId := c.Params("doctorId")
    month := c.Params("month")
    year := c.Params("year")
    needAppt := c.Params("needAppt")
    idoctorId, _ := strconv.ParseInt(doctorId, 10, 64)
    imonth, _ := strconv.ParseInt(month, 10, 32)
    iyear, _ := strconv.ParseInt(year, 10, 32)
    la, lb, err := cr.novaDoctorPatientApptService.FindAllByDoctorId(idoctorId, int(imonth), int(iyear), needAppt)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "calendarDailyStatus": lb,
        "doctorAppointment": la,
    })
}

// GetPatientOutstandingBillData
//
// @Tags Vesalius
// @Produce json
// @Security BearerAuth
// @Param       branchId         path      string       true  "branchId"
// @Param       prn              path      string       true  "prn"
// @Param       billNumber       path      string       true  "billNumber"
// @Success 200 {array} byte
// @Router /vesalius/outstanding-bill/{branchId}/{prn}/{billNumber} [get]
func (cr *VesaliusController) GetPatientOutstandingBillData(c fiber.Ctx) error {
    prn := c.Params("prn")
    billNumber := c.Params("billNumber")
    data, err := cr.vesaliusService.VesaliusGetPatientOutstandingBillDetails(prn, billNumber)
    if err != nil {
        return err
    }

    return c.Send(data)
}

// GetPatientData
//
// @Tags Vesalius
// @Produce json
// @Security BearerAuth
// @Param       branchId         path      string       true  "branchId"
// @Param       prn              path      string       true  "prn"
// @Success 200 {object} gm.Patient
// @Router /vesalius/patient-data/{branchId}/{prn} [get]
func (cr *VesaliusController) GetPatientData(c fiber.Ctx) error {
    prn := c.Params("prn")
    o, err := cr.vesaliusService.VesaliusGetPatientDataByPrn(prn)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

func (cr *VesaliusController) GetPatientFutureAppointments(c fiber.Ctx) error {
    prn := c.Params("prn")
    isHome := c.Params("isHome")
    cr.vesaliusService.VesaliusGetPatientFutureAppointments(prn, isHome == "1")
}

func (cr *VesaliusController) 

func (cr *VesaliusController) 

func (cr *VesaliusController) 

func (cr *VesaliusController) 

func resizeBase64Image(base64s string) (string, error) {
    i := strings.Index(base64s, "base64,")

    if i < 0 {
        base64Data := base64s
        imgBytes, err := base64.StdEncoding.DecodeString(base64Data)
        if err != nil {
            return "", err
        }

        srcImg, _, err := image.Decode(bytes.NewReader(imgBytes))
        if err != nil {
            return "", err
        }

        if len(imgBytes) > 20000 {
            resizedImg := resize.Resize(125, 0, srcImg, resize.Lanczos3)
            var buf bytes.Buffer
            if err := png.Encode(&buf, resizedImg); err != nil {
                return "", err
            }
            return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
        }
    } else {
        base64Data := base64s[i+7:]
        m := base64s[:i+7]
        imgBytes, err := base64.StdEncoding.DecodeString(base64Data)
        if err != nil {
            return "", err
        }

        srcImg, _, err := image.Decode(bytes.NewReader(imgBytes))
        if err != nil {
            return "", err
        }

        resizedImg := resize.Resize(125, 0, srcImg, resize.Lanczos3)
        var buf bytes.Buffer
        if err := png.Encode(&buf, resizedImg); err != nil {
            return "", err
        }
        s := fmt.Sprintf("%s%s", m, base64.StdEncoding.EncodeToString(buf.Bytes()))
        return s, nil
    }

    return base64s, nil
}
