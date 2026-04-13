package guest

import (
	"strconv"
	"strings"
	"time"
	"vesaliusm/config"
	"vesaliusm/controller/clubs/shared"
	"vesaliusm/dto"
	model "vesaliusm/model/clubs"
    mm "vesaliusm/model/mail"
	upck "vesaliusm/model/userPackage"
	"vesaliusm/service/clubs"
	"vesaliusm/service/country"
	"vesaliusm/service/guest"
	"vesaliusm/service/hpackage"
	"vesaliusm/service/mail"
	"vesaliusm/service/novaDoctor"
	"vesaliusm/service/novaDoctorPatientAppt"
	"vesaliusm/service/patientPurchaseDetails"
	"vesaliusm/service/payment"
	"vesaliusm/service/vesalius"
	"vesaliusm/service/wallex"
	"vesaliusm/utils"
    "vesaliusm/utils/constants"

	"github.com/gofiber/fiber/v3"
	"github.com/nleeper/goment"
)

type GuestController struct {
    clubService                   *clubs.ClubService
    guestService                  *guest.GuestService
    countryService                *country.CountryService
    novaDoctorService             *novaDoctor.NovaDoctorService
    novaDoctorPatientApptService  *novaDoctorPatientAppt.NovaDoctorPatientApptService
    packageService                *hpackage.PackageService
    patientPurchaseDetailsService *patientPurchaseDetails.PatientPurchaseDetailsService
    paymentService                *payment.PaymentService
    vesaliusService               *vesalius.VesaliusService
    wallexService                 *wallex.WallexService
    mailService                   *mail.MailService
}

func NewGuestController() *GuestController {
    return &GuestController{
        clubService:                   clubs.ClubSvc,
        guestService:                  guest.GuestSvc,
        countryService:                country.CountrySvc,
        novaDoctorService:             novaDoctor.NovaDoctorSvc,
        novaDoctorPatientApptService:  novaDoctorPatientAppt.NovaDoctorPatientApptSvc,
        packageService:                hpackage.PackageSvc,
        patientPurchaseDetailsService: patientPurchaseDetails.PatientPurchaseDetailsSvc,
        paymentService:                payment.PaymentSvc,
        vesaliusService:               vesalius.VesaliusSvc,
        wallexService:                 wallex.WallexSvc,
        mailService:                   mail.MailSvc,
    }
}

// GetAllDoctorInformation
//
// @Tags Guest
// @Produce json
// @Param        branchId          path        string  true  "branchId"
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.NovaDoctor
// @Router /guest/vesalius/getAllDoctorInformation/{branchId} [get]
func (cr *GuestController) GetAllDoctorInformation(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.novaDoctorService.List(page, limit, false)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllDoctorInformation
//
// @Tags Guest
// @Produce json
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.NovaDoctor
// @Router /guest/vesalius/getAllDoctorInformation/{branchId} [post]
func (cr *GuestController) SearchAllDoctorInformation(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key = "%" + key + "%"

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.novaDoctorService.ListByKeywordGuest(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetGuestReturningPatient
//
// @Tags Guest
// @Produce json
// @Param        request  body        dto.GuestGetReturningPatientDto  false  "GuestGetReturningPatientDto"
// @Success 200
// @Router /guest/appointment/returning-patient [post]
func (cr *GuestController) GetGuestReturningPatient(c fiber.Ctx) error {
    data := new(dto.GuestGetReturningPatientDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    patient, _, err := cr.vesaliusService.VesaliusGetPatientDataByNric(data.IdentificationNumber)
    if patient == nil {
        return fiber.NewError(fiber.StatusBadRequest, "The Identification Number provided does not exist in our hospital records. Please retry.")
    }

    if err != nil {
        return err
    }

    lname := []string{strings.TrimSpace(patient.Name.FirstName)}
    if strings.TrimSpace(patient.Name.MiddleName) != "" {
        lname = append(lname, strings.TrimSpace(patient.Name.MiddleName))
    }

    if strings.TrimSpace(patient.Name.LastName) != "" {
        lname = append(lname, strings.TrimSpace(patient.Name.LastName))
    }

    patientFullName := strings.Join(lname, " ")
    err = cr.guestService.SavePatientTempGuestInfo(patient, patientFullName)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "prn":  patient.Prn,
        "name": patientFullName,
    })
}

// MakeGuestNewPatient
//
// @Tags Guest
// @Produce json
// @Param        request  body        dto.GuestMakeNewPatientDto  false  "GuestMakeNewPatientDto"
// @Success 200
// @Router /guest/appointment/new-patient [post]
func (cr *GuestController) MakeGuestNewPatient(c fiber.Ctx) error {
    data := new(dto.GuestMakeNewPatientDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    vesCountryCode, err := cr.countryService.FindCountryCodeByNationality(data.Nationality)
    if err != nil {
        return err
    }

    if vesCountryCode == "" {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Nationality")
    } else {
        data.Nationality = vesCountryCode
    }

    patient, ex := cr.vesaliusService.VesaliusCheckpatientDataByNric(data.IdentificationNumber)
    if ex.Code == "WS-00014" {
        newPerson, err := cr.vesaliusService.VesaliusProcessPersonBiodata(data)
        if err != nil {
            return err
        }

        lname := []string{strings.TrimSpace(newPerson.Name.FirstName)}
        if strings.TrimSpace(newPerson.Name.MiddleName) != "" {
            lname = append(lname, strings.TrimSpace(newPerson.Name.MiddleName))
        }

        if strings.TrimSpace(newPerson.Name.LastName) != "" {
            lname = append(lname, strings.TrimSpace(newPerson.Name.LastName))
        }

        personFullName := strings.Join(lname, " ")
        err = cr.guestService.SavePersonTempGuestInfo(newPerson, personFullName)
        if err != nil {
            return err
        }

        prn := "Invalid Person Number"
        if newPerson.Document[0].Code == "X2" {
            prn = newPerson.Document[0].Value
        }
        return c.JSON(fiber.Map{
            "prn":  prn,
            "name": personFullName,
        })
    } else {
        lname := []string{strings.TrimSpace(patient.Name.FirstName)}
        if strings.TrimSpace(patient.Name.MiddleName) != "" {
            lname = append(lname, strings.TrimSpace(patient.Name.MiddleName))
        }

        if strings.TrimSpace(patient.Name.LastName) != "" {
            lname = append(lname, strings.TrimSpace(patient.Name.LastName))
        }

        patientFullName := strings.Join(lname, " ")
        err = cr.guestService.SavePatientTempGuestInfo(patient, patientFullName)
        if err != nil {
            return err
        }

        return c.JSON(fiber.Map{
            "prn":  patient.Prn,
            "name": patientFullName,
        })
    }
}

// GetDoctorAppointments
//
// @Tags Guest
// @Produce json
// @Param       doctorId       path      string       true  "doctorId"
// @Param       month          path      string       true  "month"
// @Param       year           path      string       true  "year"
// @Param       needAppt       path      string       true  "needAppt"
// @Success 200
// @Router /guest/appointment/get-doctor-appointments/{doctorId}/{month}/{year}/{needAppt} [get]
func (cr *GuestController) GetDoctorAppointments(c fiber.Ctx) error {
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
        "doctorAppointment":   la,
    })
}

// CheckGuestPatientAppointment
//
// @Tags Guest
// @Produce json
// @Param       branchId     path      string                         true  "branchId"
// @Param       prn          path      string                         true  "prn"
// @Param       request      body      dto.PostCheckAppointmentDto    true  "PostCheckAppointmentDto"
// @Success 200 {boolean} boolean
// @Router /guest/appointment/check-make-appointment/{branchId}/{prn} [post]
func (cr *GuestController) CheckGuestPatientAppointment(c fiber.Ctx) error {
    data := new(dto.PostCheckAppointmentDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    prn := c.Params("prn")
    g, _ := goment.New(data.ApptDate, "YYYY-MM-DD")
    convertedDate := g.Format("YYYY-MM-DD")
    b, err := cr.novaDoctorPatientApptService.ExistsByPrnDateSessionType(prn, convertedDate, data.ApptSessionType)
    if err != nil {
        return err
    }

    return c.JSON(b)
}

// GetMakeGuestAppointment
//
// @Tags Guest
// @Produce json
// @Param       branchId     path      string                         true  "branchId"
// @Param       prn          path      string                         true  "prn"
// @Param       request      body      dto.PostMakeAppointmentDto     true  "PostMakeAppointmentDto"
// @Success 200
// @Router /guest/appointment/make-appointment/{branchId}/{prn} [post]
func (cr *GuestController) GetMakeGuestAppointment(c fiber.Ctx) error {
    data := new(dto.PostMakeAppointmentDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    prn := c.Params("prn")
    o, err := cr.vesaliusService.VesaliusGetMakeAppointment(prn, data)
    if err != nil {
        return err
    }

    if o != nil {
        guest, err := cr.guestService.GetTempGuestInfo(prn)
        if err != nil {
            return err
        }

        email := ""
        if guest.GuestEmail.Valid {
            email = guest.GuestEmail.String
        }

        emailPrm := mm.MailGuestAppointmentConfirmation{
            GuestName:       guest.GuestName.String,
            DoctorName:      o.DoctorName,
            AppointmentDate: o.Date,
            AppointmentTime: o.StartTime,
            ClinicLocation:  o.Clinic,
            Email:           email,
        }
        if email != "" {
            go func() {
                err := cr.mailService.SendGuestAppointmentConfirmationToPatient(emailPrm)
                if err == nil {
                    _ = cr.guestService.DeleteTempGuestInfo(guest.GuestPRN.String)
                }
            }()
            emailPrm.Email = ""
        }

        go func() {
            cr.mailService.SendGuestAppointmentConfirmationToIH(emailPrm)
        }()
    }

    return c.JSON(fiber.Map{
        "success": 1,
    })
}

// GetAllGuestNotificationLists
//
// @Tags Guest
// @Produce json
// @Param        playerId          path        string  true  "playerId"
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.OnesignalNotification
// @Router /guest/notification/all/{playerId} [get]
func (cr *GuestController) GetAllGuestNotificationLists(c fiber.Ctx) error {
    playerId := c.Params("playerId")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.guestService.ListGuestNotification(playerId, page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetGuestUnseenNotificationCount
//
// @Tags Guest
// @Produce json
// @Param        notificationId    path        string  true  "notificationId"
// @Param        playerId          path        string  true  "playerId"
// @Success 200
// @Router /guest/notification/seen/{notificationId}/{playerId} [get]
func (cr *GuestController) GetGuestUnseenNotificationCount(c fiber.Ctx) error {
    notificationId := c.Params("notificationId")
    playerId := c.Params("playerId")
    inotificationId, _ := strconv.ParseInt(notificationId, 10, 64)
    err := cr.guestService.UpdateSeenNotificationByPlayerId(playerId, inotificationId)
    if err != nil {
        return err
    }

    n, err := cr.guestService.CountUnseenNotificationByGuestPlayerId(playerId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage":  "Guest Notification has been set to seen",
        "userUnseenCount": n,
    })
}

// GetLittleKidsAboutUs
//
// @Tags Guest
// @Produce json
// @Success 200 {object} clubs.LittleExplorersKidsAboutUs
// @Router /guest/clubs/littlekids/about-us [get]
func (cr *GuestController) GetLittleKidsAboutUs(c fiber.Ctx) error {
    o, err := cr.clubService.FindLittleKidsAboutUs()
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// CreateGuestLittleKidsMembership
//
// @Tags Guest
// @Produce json
// @Param request body dto.LittleExplorersKidsMembershipDto true "LittleExplorersKidsMembershipDto"
// @Success 200
// @Router /guest/clubs/littlekids/membership [post]
func (cr *GuestController) CreateGuestLittleKidsMembership(c fiber.Ctx) error {
    data := new(dto.LittleExplorersKidsMembershipDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.clubService.ExistsLittleKidsByDocTypeDocNo(data.KidsDocType, strings.TrimSpace(data.KidsDocNumber))
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, "There is an existing member registered with the identification type / number. Please double check")
    }

    b2, err := cr.clubService.ExistsLittleKidsByPrn(strings.TrimSpace(data.KidsDocNumber))
    if err != nil {
        return err
    }

    if b2 {
        return fiber.NewError(fiber.StatusBadRequest, "You have already registered previously. Please reach out to our customer service at +604-238 3388 for further action")
    }

    eligibleAge := shared.LittleKidsEligibleAge(data.KidsDob)
    if !eligibleAge {
        return fiber.NewError(fiber.StatusBadRequest, "Only 12 years old and below")
    }

    if !strings.EqualFold(data.KidsDocType, constants.ClubsDocTypeNRIC) &&
        !strings.EqualFold(data.KidsDocType, constants.ClubsDocTypeBirthCert) &&
        !strings.EqualFold(data.KidsDocType, constants.ClubsDocTypePassport) {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Kids Document Type")
    }

    if !strings.EqualFold(data.GuardianDocType, constants.ClubsDocTypeNRIC) &&
        !strings.EqualFold(data.GuardianDocType, constants.ClubsDocTypeBirthCert) &&
        !strings.EqualFold(data.GuardianDocType, constants.ClubsDocTypePassport) {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Guardian Document Type")
    }

    if strings.EqualFold(data.KidsDocType, constants.ClubsDocTypeNRIC) &&
        strings.EqualFold(data.GuardianDocType, constants.ClubsDocTypeNRIC) &&
        strings.TrimSpace(data.KidsDocNumber) == strings.TrimSpace(data.GuardianDocNumber) ||
        strings.EqualFold(data.KidsDocType, constants.ClubsDocTypePassport) &&
            strings.EqualFold(data.GuardianDocType, constants.ClubsDocTypePassport) &&
            strings.TrimSpace(data.KidsDocNumber) == strings.TrimSpace(data.GuardianDocNumber) ||
        strings.EqualFold(data.KidsDocType, constants.ClubsDocTypeBirthCert) &&
            strings.EqualFold(data.GuardianDocType, constants.ClubsDocTypeBirthCert) &&
            strings.TrimSpace(data.KidsDocNumber) == strings.TrimSpace(data.GuardianDocNumber) {
        return fiber.NewError(fiber.StatusBadRequest, "Kids Identification Number and Guardian Identification Number cannot be same")
    }

    if strings.EqualFold(data.KidsDocType, constants.ClubsDocTypeNRIC) {
        data.KidsDocNumber = strings.ReplaceAll(data.KidsDocNumber, "-", "")
    }

    if strings.EqualFold(data.GuardianDocType, constants.ClubsDocTypeNRIC) {
        data.GuardianDocNumber = strings.ReplaceAll(data.GuardianDocNumber, "-", "")
    }

    var o model.LittleExplorersKidsMembership
    kidsMembershipNo, err := cr.clubService.GenerateKidsMembershipNo()
    if err != nil {
        return err
    }

    o.KidsMembershipNumber = utils.NewNullString(kidsMembershipNo)
    o.KidsName = utils.NewNullString(data.KidsName)
    o.KidsDob = utils.NewNullString(data.KidsDob)
    o.KidsDocType = utils.NewNullString(data.KidsDocType)
    o.KidsDocNumber = utils.NewNullString(data.KidsDocNumber)
    o.KidsGender = utils.NewNullString(data.KidsGender)
    o.KidsNationality = utils.NewNullString(data.KidsNationality)
    o.KidsEmail = utils.NewNullString(data.KidsEmail)
    o.GuardianName = utils.NewNullString(data.GuardianName)
    o.GuardianDob = utils.NewNullString(data.GuardianDob)
    o.GuardianDocType = utils.NewNullString(data.GuardianDocType)
    o.GuardianDocNumber = utils.NewNullString(data.GuardianDocNumber)
    o.GuardianGender = utils.NewNullString(data.GuardianGender)
    o.GuardianNationality = utils.NewNullString(data.GuardianNationality)
    o.GuardianEmail = utils.NewNullString(data.GuardianEmail)
    o.GuardianHomeContact = utils.NewNullString(data.GuardianHomeContact)
    o.GuardianMobileContact = utils.NewNullString(data.GuardianMobileContact)
    o.GuardianAddress1 = utils.NewNullString(data.GuardianAddress1)
    o.GuardianAddress2 = utils.NewNullString(data.GuardianAddress2)
    o.GuardianAddress3 = utils.NewNullString(data.GuardianAddress3)
    o.GuardianPostCode = utils.NewNullString(data.GuardianPostCode)
    o.GuardianState = utils.NewNullString(data.GuardianState)
    o.GuardianCountryCode = utils.NewNullString(data.GuardianCountryCode)
    o.Relationship = utils.NewNullString(data.Relationship)
    o.PreferredLanguage = utils.NewNullString(data.PreferredLanguage)

    err = cr.clubService.SaveLittleKidsMembership(o)
    if err != nil {
        return err
    }

    emailPrm := mm.MailLittleKids{
        KidsName: o.KidsName.String,
        Email:    "",
    }
    if o.KidsEmail.Valid {
        emailPrm.Email = o.KidsEmail.String
        go func() {
            cr.mailService.SendLittleKids(emailPrm)
        }()
        emailPrm.Email = ""
    }
    if o.GuardianEmail.Valid {
        emailPrm.Email = o.GuardianEmail.String
        go func() {
            cr.mailService.SendLittleKids(emailPrm)
        }()
        emailPrm.Email = ""
    }

    go func() {
        cr.mailService.SendLittleKids(emailPrm)
    }()

    return c.JSON(fiber.Map{
        "message": "Guest Little Explorers Kids Membership created",
    })
}

// GetAllAppLittleKidsMemberships
//
// @Tags Guest
// @Produce json
// @Param        identificationNumber    path        string  true  "identificationNumber"
// @Success 200 {object} model.LittleExplorersKidsMembership
// @Router /guest/clubs/littlekids/membership/{identificationNumber} [get]
func (cr *GuestController) GetAllAppLittleKidsMemberships(c fiber.Ctx) error {
    identificationNumber := c.Params("identificationNumber")
    o, err := cr.clubService.FindGuestLittleKidsMembershipByIc(identificationNumber)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// GetAllAppLittleKidsActivities
//
// @Tags Guest
// @Produce json
// @Param        isHome    path        string  true   "isHome"
// @Param        _page     query       string  false  "_page"  default:"1"
// @Param        _limit    query       string  false  "_limit" default:"10"
// @Success 200 {object} model.PagedList
// @Router /guest/clubs/littlekids/activity/all/mobile/{isHome} [get]
func (cr *GuestController) GetAllAppLittleKidsActivities(c fiber.Ctx) error {
    isHome := c.Params("isHome")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.clubService.ListAppLittleKidsActivities(isHome == "1", page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// ParticipateLittleKidsActivity
//
// @Tags Guest
// @Produce json
// @Param request body dto.LittleExplorersKidsActvParticipationDto true "LittleExplorersKidsActvParticipationDto"
// @Success 200
// @Router /guest/clubs/littlekids/activity/participate [post]
func (cr *GuestController) ParticipateLittleKidsActivity(c fiber.Ctx) error {
    data := new(dto.LittleExplorersKidsActvParticipationDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    actvParticipation := make([]model.LittleExplorersKidsActvParticipation, 0)
    for i := range data.KidsActvParticipation {
        participation := data.KidsActvParticipation[i]
        var o model.LittleExplorersKidsActvParticipation
        o.KidsActivityId = participation.KidsActivityId
        o.KidsMembershipId = participation.KidsMembershipId
        o.ActivityDateTime = participation.ActivityDateTime

        actvParticipation = append(actvParticipation, o)

        activity, err := cr.clubService.FindLittleKidsActivityNameById(int64(o.KidsActivityId))
        if err != nil {
            return err
        }
        kidsMember, err := cr.clubService.FindLittleKidsMembershipById(int64(o.KidsMembershipId))
        if err != nil {
            return err
        }

        emailPrm := mm.MailClubsEventRegistrationToMember{
            ActivityName: activity.KidsActivityName.String,
            MemberName:   kidsMember.KidsName.String,
            Email:        "",
        }
        if kidsMember.KidsEmail.Valid {
            emailPrm.Email = kidsMember.KidsEmail.String
            go func() {
                cr.mailService.SendClubsEventRegistrationToMember(emailPrm)
            }()
            emailPrm.Email = ""
        }

        go func() {
            cr.mailService.SendClubsEventRegistrationToIH(emailPrm)
        }()
    }

    err := cr.clubService.ParticipateLittleKidsActivity(actvParticipation)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Activity Participate successfully",
    })
}

// GetGoldenPearlAboutUs
//
// @Tags Guest
// @Produce json
// @Success 200 {object} model.GoldenPearlAboutUs
// @Router /guest/clubs/goldenpearl/about-us [get]
func (cr *GuestController) GetGoldenPearlAboutUs(c fiber.Ctx) error {
    o, err := cr.clubService.FindGoldenPearlAboutUs()
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// CreateGuestGoldenPearlMembership
//
// @Tags Guest
// @Produce json
// @Param request body dto.GoldenPearlMembershipDto true "GoldenPearlMembershipDto"
// @Success 200
// @Router /guest/clubs/goldenpearl/membership [post]
func (cr *GuestController) CreateGuestGoldenPearlMembership(c fiber.Ctx) error {
    data := new(dto.GoldenPearlMembershipDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.clubService.ExistsGoldenPearlByDocTypeDocNo(data.GoldenDocType, strings.TrimSpace(data.GoldenDocNumber))
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, "There is an existing member registered with the identification type / number. Please double check")
    }

    b2, err := cr.clubService.ExistsGoldenPearlByPrn(strings.TrimSpace(data.GoldenDocNumber))
    if err != nil {
        return err
    }

    if b2 {
        return fiber.NewError(fiber.StatusBadRequest, "You have already registered previously. Please reach out to our customer service at +604-238 3388 for further action")
    }

    eligibleAge := shared.GoldenPearlEligibleAge(data.GoldenDob)
    if !eligibleAge {
        return fiber.NewError(fiber.StatusBadRequest, "Only 60 years old and above")
    }

    if !strings.EqualFold(data.GoldenDocType, constants.ClubsDocTypeNRIC) &&
        !strings.EqualFold(data.GoldenDocType, constants.ClubsDocTypeBirthCert) &&
        !strings.EqualFold(data.GoldenDocType, constants.ClubsDocTypePassport) {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Golden Pearl Document Type")
    }

    if !strings.EqualFold(data.NokDocType, constants.ClubsDocTypeNRIC) &&
        !strings.EqualFold(data.NokDocType, constants.ClubsDocTypeBirthCert) &&
        !strings.EqualFold(data.NokDocType, constants.ClubsDocTypePassport) {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid NOK Document Type")
    }

    if strings.EqualFold(data.GoldenDocType, constants.ClubsDocTypeNRIC) &&
        strings.EqualFold(data.NokDocType, constants.ClubsDocTypeNRIC) &&
        strings.EqualFold(strings.TrimSpace(data.GoldenDocNumber), strings.TrimSpace(data.NokDocNumber)) ||
        strings.EqualFold(data.GoldenDocType, constants.ClubsDocTypePassport) &&
            strings.EqualFold(data.NokDocType, constants.ClubsDocTypePassport) &&
            strings.EqualFold(strings.TrimSpace(data.GoldenDocNumber), strings.TrimSpace(data.NokDocNumber)) ||
        strings.EqualFold(data.GoldenDocType, constants.ClubsDocTypeBirthCert) &&
            strings.EqualFold(data.NokDocType, constants.ClubsDocTypeBirthCert) &&
            strings.EqualFold(strings.TrimSpace(data.GoldenDocNumber), strings.TrimSpace(data.NokDocNumber)) {
        return fiber.NewError(fiber.StatusBadRequest, "Golden Pearl Identification Number and NOK Identification Number cannot be same")
    }

    if strings.EqualFold(data.GoldenDocType, constants.ClubsDocTypeNRIC) {
        data.GoldenDocNumber = strings.ReplaceAll(data.GoldenDocNumber, "-", "")
    }

    if strings.EqualFold(data.NokDocType, constants.ClubsDocTypeNRIC) {
        data.NokDocNumber = strings.ReplaceAll(data.NokDocNumber, "-", "")
    }

    var o model.GoldenPearlMembership
    goldenMembershipNo, err := cr.clubService.GenerateGoldenMembershipNo()
    if err != nil {
        return err
    }

    o.GoldenMembershipNumber = utils.NewNullString(goldenMembershipNo)
    o.GoldenName = utils.NewNullString(data.GoldenName)
    o.GoldenDob = utils.NewNullString(data.GoldenDob)
    o.GoldenDocType = utils.NewNullString(data.GoldenDocType)
    o.GoldenDocNumber = utils.NewNullString(data.GoldenDocNumber)
    o.GoldenGender = utils.NewNullString(data.GoldenGender)
    o.GoldenNationality = utils.NewNullString(data.GoldenNationality)
    o.GoldenEmail = utils.NewNullString(data.GoldenEmail)
    o.NokName = utils.NewNullString(data.NokName)
    o.NokDob = utils.NewNullString(data.NokDob)
    o.NokDocType = utils.NewNullString(data.NokDocType)
    o.NokDocNumber = utils.NewNullString(data.NokDocNumber)
    o.NokGender = utils.NewNullString(data.NokGender)
    o.NokNationality = utils.NewNullString(data.NokNationality)
    o.NokEmail = utils.NewNullString(data.NokEmail)
    o.NokHomeContact = utils.NewNullString(data.NokHomeContact)
    o.NokMobileContact = utils.NewNullString(data.NokMobileContact)
    o.NokAddress1 = utils.NewNullString(data.NokAddress1)
    o.NokAddress2 = utils.NewNullString(data.NokAddress2)
    o.NokAddress3 = utils.NewNullString(data.NokAddress3)
    o.NokPostCode = utils.NewNullString(data.NokPostCode)
    o.NokState = utils.NewNullString(data.NokState)
    o.NokCountryCode = utils.NewNullString(data.NokCountryCode)
    o.Relationship = utils.NewNullString(data.Relationship)
    o.PreferredLanguage = utils.NewNullString(data.PreferredLanguage)

    err = cr.clubService.SaveGoldenPearlMembership(o)
    if err != nil {
        return err
    }

    emailPrm := mm.MailGoldenPearl{
        GoldenName: o.GoldenName.String,
        Email:      "",
    }
    if o.GoldenEmail.Valid {
        emailPrm.Email = o.GoldenEmail.String
        go func() {
            cr.mailService.SendGoldenPearl(emailPrm)
        }()
        emailPrm.Email = ""
    }
    if o.NokEmail.Valid {
        emailPrm.Email = o.NokEmail.String
        go func() {
            cr.mailService.SendGoldenPearl(emailPrm)
        }()
        emailPrm.Email = ""
    }

    go func() {
        cr.mailService.SendGoldenPearl(emailPrm)
    }()

    return c.JSON(fiber.Map{
        "message": "Guest Golden Pearl Membership created",
    })
}

// GetAllAppGoldenPearlMemberships
//
// @Tags Guest
// @Produce json
// @Param identificationNumber path string true "identificationNumber"
// @Success 200 {array} model.GoldenPearlMembership
// @Router /guest/clubs/goldenpearl/membership/{identificationNumber} [get]
func (cr *GuestController) GetAllAppGoldenPearlMemberships(c fiber.Ctx) error {
    identificationNumber := c.Params("identificationNumber")
    o, err := cr.clubService.FindGuestGoldenPearlMembershipByIc(identificationNumber)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// GetAllAppGoldenPearlActivities
//
// @Tags Guest
// @Produce json
// @Param        isHome            path        string  true   "isHome"
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.GoldenPearlActivity
// @Router /guest/clubs/goldenpearl/activity/all/mobile/{isHome} [get]
func (cr *GuestController) GetAllAppGoldenPearlActivities(c fiber.Ctx) error {
    isHome := c.Params("isHome")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.clubService.ListAppGoldenPearlActivities(isHome == "1", page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// ParticipateGoldenPearlActivity
//
// @Tags Guest
// @Produce json
// @Param request body dto.GoldenPearlActvParticipationDto true "GoldenPearlActvParticipationDto"
// @Success 200
// @Router /guest/clubs/goldenpearl/activity/participate [post]
func (cr *GuestController) ParticipateGoldenPearlActivity(c fiber.Ctx) error {
    data := new(dto.GoldenPearlActvParticipationDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    actvParticipation := make([]model.GoldenPearlActvParticipation, 0)
    for i := range data.GoldenActvParticipation {
        participation := data.GoldenActvParticipation[i]

        var o model.GoldenPearlActvParticipation
        o.GoldenActivityId = participation.GoldenActivityId
        o.GoldenMembershipId = participation.GoldenMembershipId
        o.ActivityDateTime = participation.ActivityDateTime

        actvParticipation = append(actvParticipation, o)

        activity, err := cr.clubService.FindGoldenPearlActivityNameById(int64(o.GoldenActivityId))
        if err != nil {
            return err
        }

        goldenMember, err := cr.clubService.FindGoldenPearlMembershipById(int64(o.GoldenMembershipId))
        if err != nil {
            return err
        }

        emailPrm := mm.MailClubsEventRegistrationToMember{
            ActivityName: activity.GoldenActivityName.String,
            MemberName:   goldenMember.GoldenName.String,
            Email:        "",
        }
        if goldenMember.GoldenEmail.Valid {
            emailPrm.Email = goldenMember.GoldenEmail.String
            go func() {
                cr.mailService.SendClubsEventRegistrationToMember(emailPrm)
            }()
            emailPrm.Email = ""
        }

        go func() {
            cr.mailService.SendClubsEventRegistrationToIH(emailPrm)
        }()
    }

    err := cr.clubService.ParticipateGoldenPearlActivity(actvParticipation)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Activity Participate successfully",
    })
}

// GetAllAppPackages
// @Tags Guest
// @Produce json
// @Param        isHome            path        string  true   "isHome"
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200
// @Router /guest/package/all/mobile/{isHome} [get]
func (cr *GuestController) GetAllAppPackages(c fiber.Ctx) error {
    isHome := c.Params("isHome")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(constants.PAGE_SIZE))
    m, err := cr.packageService.ListApp(isHome == "1", page, limit)
    if err != nil {
        return err
    }

    c.Set(constants.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(constants.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetPackageStatusById
// @Tags Guest
// @Produce json
// @Param        packageId            path        string  true   "packageId"
// @Success 200
// @Router /guest/package/packageStatus/{packageId} [get]
func (cr *GuestController) GetPackageStatusById(c fiber.Ctx) error {
    packageId := c.Params("packageId")
    ipackageId, _ := strconv.ParseInt(packageId, 10, 64)
    o, err := cr.packageService.FindPackageStatusByPackageId(ipackageId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// GetPackageById
// @Tags Guest
// @Produce json
// @Param        packageId            path        string  true   "packageId"
// @Success 200
// @Router /guest/package/{packageId} [get]
func (cr *GuestController) GetPackageById(c fiber.Ctx) error {
    packageId := c.Params("packageId")
    ipackageId, _ := strconv.ParseInt(packageId, 10, 64)
    o, err := cr.packageService.FindByPackageId(ipackageId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// CreateGuestPurchaseDetails
// @Tags Guest
// @Produce json
// @Param        paymentMethod      path      string                     true   "paymentMethod"
// @Param        request            body      dto.CreateGuestPackageDto  true   "request"
// @Success 200
// @Router /guest/purchase/{paymentMethod} [post]
func (cr *GuestController) CreateGuestPurchaseDetails(c fiber.Ctx) error {
    paymentMethod := c.Params("paymentMethod")
    ipaymentMethod, _ := strconv.ParseInt(paymentMethod, 10, 32)
    data := new(dto.CreateGuestPackageDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    guestPackage := make([]upck.UserPackage, 0)
    products := make([]utils.Map, 0)
    totalPrice := 0.0

    paymentRefNo, err := cr.paymentService.GeneratePaymentRefNo()
    if err != nil {
        return err
    }

    for _, pkg := range data.GuestPackage {
        p := utils.Map{
            "name":        pkg.PackageName,
            "description": pkg.PackageName,
            "price":       pkg.PackagePrice,
            "quantity":    pkg.QuantityPurchased,
            "createdAt":   strconv.FormatInt(time.Now().UnixMilli(), 10),
            "updatedAt":   strconv.FormatInt(time.Now().UnixMilli(), 10),
            "deletedAt":   "",
        }
        products = append(products, p)

        totalPrice = totalPrice + pkg.PackagePrice*float64(pkg.QuantityPurchased)

        o := upck.UserPackage{
            PatientPrn:        utils.NewNullString("-"),
            PatientName:       utils.NewNullString(data.GuestPackagePayment.BillingFullname),
            PackageId:         utils.NewInt64(pkg.PackageId),
            QuantityPurchased: pkg.QuantityPurchased,
            PackageStatus:     utils.NewNullString(string(constants.PackageStatusOrdered)),
        }
        guestPackage = append(guestPackage, o)
    }

    lsaddr := make([]string, 0)
    pyt := data.GuestPackagePayment
    if pyt.BillingAddress1 != "" {
        lsaddr = append(lsaddr, pyt.BillingAddress1)
    }

    if pyt.BillingAddress2 != "" {
        lsaddr = append(lsaddr, pyt.BillingAddress2)
    }

    if pyt.BillingAddress3 != "" {
        lsaddr = append(lsaddr, pyt.BillingAddress3)
    }

    if pyt.BillingTowncity != "" {
        lsaddr = append(lsaddr, pyt.BillingTowncity)
    }

    if pyt.BillingState != "" {
        lsaddr = append(lsaddr, pyt.BillingState)
    }

    addr := strings.Join(lsaddr, ", ")

    if ipaymentMethod == constants.PaymentMethodWallex {
        wallexPrm := utils.Map{
            "collectionRequestNumber": paymentRefNo,
            "currency":                "MYR",
            "paymentPurpose":          "SCVE",
            "paymentCurrency":         "IDR",
            "paymentPartial":          false,
            "remarks":                 paymentRefNo,
            "customerInfo": utils.Map{
                "name":             pyt.BillingFullname,
                "ituTelephoneCode": pyt.BillingContactCode,
                "mobileNumber":     pyt.BillingContactNo,
                "email":            pyt.BillingEmail,
                "address":          addr,
            },
            "products": products,
        }
        wallexRes, err := cr.wallexService.SubmitPaymentRequest(wallexPrm)
        if err != nil {
            return err
        }

        if totalPrice != wallexRes.Amount {
            return fiber.NewError(fiber.StatusBadRequest, "Incorrect Total Price")
        }

        guestPackagePayment := upck.PackagePaymentDetails{
            PaymentGateway:         utils.NewInt32(int32(ipaymentMethod)),
            PaymentRequestId:       utils.NewNullString(wallexRes.ID),
            PaymentRequestNo:       utils.NewNullString(wallexRes.CollectionRequestNumber),
            PaymentRefId:           utils.NewNullString(wallexRes.ReferenceId),
            PaymentRequestCurrency: utils.NewNullString(wallexRes.Currency),
            PaymentAmount:          utils.NewFloat(wallexRes.Amount),
            PaymentPurpose:         utils.NewNullString(wallexRes.PaymentPurpose),
            PaymentCurrency:        utils.NewNullString(wallexRes.PaymentCurrency),
            PaymentAmountCollected: utils.NewFloat(wallexRes.PaymentAmountCollected),
            PaymentRemarks:         utils.NewNullString(wallexRes.Remarks),
            PaymentStatus:          utils.NewNullString(wallexRes.Status),
            PaymentAuthCode:        utils.NewNullString(""),
            PaymentErrorDesc:       utils.NewNullString(""),
            PaymentTransDate:       utils.NewNullString(""),
            BillingFullname:        utils.NewNullString(pyt.BillingFullname),
            BillingAddress1:        utils.NewNullString(pyt.BillingAddress1),
            BillingAddress2:        utils.NewNullString(pyt.BillingAddress2),
            BillingAddress3:        utils.NewNullString(pyt.BillingAddress3),
            BillingTowncity:        utils.NewNullString(pyt.BillingTowncity),
            BillingState:           utils.NewNullString(pyt.BillingState),
            BillingPostcode:        utils.NewNullString(pyt.BillingPostcode),
            BillingCountryCode:     utils.NewNullString(pyt.BillingCountryCode),
            BillingContactNo:       utils.NewNullString(pyt.BillingContactNo),
            BillingContactCode:     utils.NewNullString(pyt.BillingContactCode),
            BillingEmail:           utils.NewNullString(pyt.BillingEmail),
            PaymentUrl:             wallexRes.PaymentUrl,
        }
        err = cr.paymentService.SaveGuestWallex(guestPackagePayment, guestPackage)
        if err != nil {
            return err
        }

        return  c.JSON(fiber.Map{
            "message": "Guest Purchase Details created",
            "wallexDetails": fiber.Map{
                "expiredAt": wallexRes.ExpiredAt,
                "paymentUrl": wallexRes.PaymentUrl,
            },
        })
    } else {
        paymentAmt := totalPrice
        if config.GetIpayTestEnv() == "Y" {
            paymentAmt = 1
        }
        guestPackagePayment := &upck.PackagePaymentDetails{
            PaymentGateway:         utils.NewInt32(int32(ipaymentMethod)),
            PaymentRequestNo:       utils.NewNullString(paymentRefNo),
            PaymentRequestCurrency: utils.NewNullString("MYR"),
            PaymentAmount:          utils.NewFloat(paymentAmt),
            PaymentPurpose:         utils.NewNullString("Hospital Pkg Purchase"),
            PaymentCurrency:        utils.NewNullString("MYR"),
            PaymentRemarks:         utils.NewNullString(paymentRefNo),
            PaymentStatus:          utils.NewNullString("unpaid"),
            PaymentAuthCode:        utils.NewNullString(""),
            PaymentErrorDesc:       utils.NewNullString(""),
            PaymentTransDate:       utils.NewNullString(""),
            BillingFullname:        utils.NewNullString(pyt.BillingFullname),
            BillingAddress1:        utils.NewNullString(pyt.BillingAddress1),
            BillingAddress2:        utils.NewNullString(pyt.BillingAddress2),
            BillingAddress3:        utils.NewNullString(pyt.BillingAddress3),
            BillingTowncity:        utils.NewNullString(pyt.BillingTowncity),
            BillingState:           utils.NewNullString(pyt.BillingState),
            BillingPostcode:        utils.NewNullString(pyt.BillingPostcode),
            BillingCountryCode:     utils.NewNullString(pyt.BillingCountryCode),
            BillingContactNo:       utils.NewNullString(pyt.BillingContactNo),
            BillingContactCode:     utils.NewNullString(pyt.BillingContactCode),
            BillingEmail:           utils.NewNullString(pyt.BillingEmail),
        }
        err = cr.paymentService.SaveGuestIPay(*guestPackagePayment, guestPackage)
        if err != nil {
            return err
        }

        return c.JSON(paymentRefNo)
    }
}

// CheckPackageExpiryMaxpurchase
// @Tags Guest
// @Produce json
// @Param request body dto.CheckPackageExpiryMaxpurchaseDto true "CheckPackageExpiryMaxpurchaseDto"
// @Success 200
// @Router /guest/package/check/expiry-maxpurchase [post]
func (cr *GuestController) CheckPackageExpiryMaxpurchase(c fiber.Ctx) error {
    data := new(dto.CheckPackageExpiryMaxpurchaseDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    cartIsValid := true
    cartResult := make([]*upck.PackageCheckResult, 0)
    for _, pkg := range data.Package {
        r, err := cr.patientPurchaseDetailsService.CheckPackageExpiryMaxPurchase(pkg.PackageId, pkg.QuantityPurchased)
        if err != nil {
            return err
        }

        if r.Expired == 1 || r.Soldout == 1 || r.ExceedPurchase == 1 {
            cartResult = append(cartResult, r)
            cartIsValid = false
        }
    }

    return c.JSON(fiber.Map{
        "cartIsValid": cartIsValid,
        "cartResult":  cartResult,
    })
}
