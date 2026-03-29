package guest

import (
    "strconv"
    "strings"
    "vesaliusm/controller/clubs/shared"
    "vesaliusm/dto"
    "vesaliusm/model/clubs"
    "vesaliusm/model/userPackage"
    clubsSvc "vesaliusm/service/clubs"
    "vesaliusm/service/guest"
    "vesaliusm/service/mail"
    "vesaliusm/service/novaDoctor"
    "vesaliusm/service/patientPurchaseDetails"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

type GuestController struct {
    clubService                   *clubsSvc.ClubService
    guestService                  *guest.GuestService
    novaDoctorService             *novaDoctor.NovaDoctorService
    patientPurchaseDetailsService *patientPurchaseDetails.PatientPurchaseDetailsService
    mailService                   *mail.MailService
}

func NewGuestController() *GuestController {
    return &GuestController{
        clubService:                   clubsSvc.ClubSvc,
        guestService:                  guest.GuestSvc,
        novaDoctorService:             novaDoctor.NovaDoctorSvc,
        patientPurchaseDetailsService: patientPurchaseDetails.PatientPurchaseDetailsSvc,
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
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.novaDoctorService.List(page, limit, false)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
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
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.novaDoctorService.ListByKeywordGuest(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
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
func (cr *GuestController) getAllGuestNotificationLists(c fiber.Ctx) error {
    playerId := c.Params("playerId")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.guestService.ListGuestNotification(playerId, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
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

    if !strings.EqualFold(data.KidsDocType, utils.ClubsDocTypeNRIC) &&
        !strings.EqualFold(data.KidsDocType, utils.ClubsDocTypeBirthCert) &&
        !strings.EqualFold(data.KidsDocType, utils.ClubsDocTypePassport) {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Kids Document Type")
    }

    if !strings.EqualFold(data.GuardianDocType, utils.ClubsDocTypeNRIC) &&
        !strings.EqualFold(data.GuardianDocType, utils.ClubsDocTypeBirthCert) &&
        !strings.EqualFold(data.GuardianDocType, utils.ClubsDocTypePassport) {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Guardian Document Type")
    }

    if strings.EqualFold(data.KidsDocType, utils.ClubsDocTypeNRIC) &&
        strings.EqualFold(data.GuardianDocType, utils.ClubsDocTypeNRIC) &&
        strings.TrimSpace(data.KidsDocNumber) == strings.TrimSpace(data.GuardianDocNumber) ||
        strings.EqualFold(data.KidsDocType, utils.ClubsDocTypePassport) &&
            strings.EqualFold(data.GuardianDocType, utils.ClubsDocTypePassport) &&
            strings.TrimSpace(data.KidsDocNumber) == strings.TrimSpace(data.GuardianDocNumber) ||
        strings.EqualFold(data.KidsDocType, utils.ClubsDocTypeBirthCert) &&
            strings.EqualFold(data.GuardianDocType, utils.ClubsDocTypeBirthCert) &&
            strings.TrimSpace(data.KidsDocNumber) == strings.TrimSpace(data.GuardianDocNumber) {
        return fiber.NewError(fiber.StatusBadRequest, "Kids Identification Number and Guardian Identification Number cannot be same")
    }

    if strings.EqualFold(data.KidsDocType, utils.ClubsDocTypeNRIC) {
        data.KidsDocNumber = strings.ReplaceAll(data.KidsDocNumber, "-", "")
    }

    if strings.EqualFold(data.GuardianDocType, utils.ClubsDocTypeNRIC) {
        data.GuardianDocNumber = strings.ReplaceAll(data.GuardianDocNumber, "-", "")
    }

    var o clubs.LittleExplorersKidsMembership
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

    emailPrm := utils.Map{
        "kidsName": o.KidsName.String,
        "email":    "",
    }
    if o.KidsEmail.Valid {
        emailPrm["email"] = o.KidsEmail.String
        go func() {
            cr.mailService.SendLittleKids(emailPrm)
        }()
        emailPrm["email"] = ""
    }
    if o.GuardianEmail.Valid {
        emailPrm["email"] = o.GuardianEmail.String
        go func() {
            cr.mailService.SendLittleKids(emailPrm)
        }()
        emailPrm["email"] = ""
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
// @Success 200 (object) clubs.LittleExplorersKidsMembership
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
// @Param        isHome    path        string  true  "isHome"
// @Success 200 (object) model.PagedList
// @Router /guest/clubs/littlekids/activity/all/mobile/{isHome} [get]
func (cr *GuestController) GetAllAppLittleKidsActivities(c fiber.Ctx) error {
    isHome := c.Params("isHome")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.clubService.ListAppLittleKidsActivities(isHome == "1", page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
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

    actvParticipation := make([]clubs.LittleExplorersKidsActvParticipation, 0)
    for i := range data.KidsActvParticipation {
        participation := data.KidsActvParticipation[i]
        var o clubs.LittleExplorersKidsActvParticipation
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

        emailPrm := utils.Map{
            "activityName": activity.KidsActivityName.String,
            "memberName":   kidsMember.KidsName.String,
            "email":        "",
        }
        if kidsMember.KidsEmail.Valid {
            emailPrm["email"] = kidsMember.KidsEmail.String
            go func() {
                cr.mailService.SendClubsEventRegistrationToMember(emailPrm)
            }()
            emailPrm["email"] = ""
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
// @Success 200 (object) clubs.GoldenPearlAboutUs
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

    if !strings.EqualFold(data.GoldenDocType, utils.ClubsDocTypeNRIC) &&
        !strings.EqualFold(data.GoldenDocType, utils.ClubsDocTypeBirthCert) &&
        !strings.EqualFold(data.GoldenDocType, utils.ClubsDocTypePassport) {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Golden Pearl Document Type")
    }

    if !strings.EqualFold(data.NokDocType, utils.ClubsDocTypeNRIC) &&
        !strings.EqualFold(data.NokDocType, utils.ClubsDocTypeBirthCert) &&
        !strings.EqualFold(data.NokDocType, utils.ClubsDocTypePassport) {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid NOK Document Type")
    }

    if strings.EqualFold(data.GoldenDocType, utils.ClubsDocTypeNRIC) &&
        strings.EqualFold(data.NokDocType, utils.ClubsDocTypeNRIC) &&
        strings.EqualFold(strings.TrimSpace(data.GoldenDocNumber), strings.TrimSpace(data.NokDocNumber)) ||
        strings.EqualFold(data.GoldenDocType, utils.ClubsDocTypePassport) &&
            strings.EqualFold(data.NokDocType, utils.ClubsDocTypePassport) &&
            strings.EqualFold(strings.TrimSpace(data.GoldenDocNumber), strings.TrimSpace(data.NokDocNumber)) ||
        strings.EqualFold(data.GoldenDocType, utils.ClubsDocTypeBirthCert) &&
            strings.EqualFold(data.NokDocType, utils.ClubsDocTypeBirthCert) &&
            strings.EqualFold(strings.TrimSpace(data.GoldenDocNumber), strings.TrimSpace(data.NokDocNumber)) {
        return fiber.NewError(fiber.StatusBadRequest, "Golden Pearl Identification Number and NOK Identification Number cannot be same")
    }

    if strings.EqualFold(data.GoldenDocType, utils.ClubsDocTypeNRIC) {
        data.GoldenDocNumber = strings.ReplaceAll(data.GoldenDocNumber, "-", "")
    }

    if strings.EqualFold(data.NokDocType, utils.ClubsDocTypeNRIC) {
        data.NokDocNumber = strings.ReplaceAll(data.NokDocNumber, "-", "")
    }

    var o clubs.GoldenPearlMembership
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

    emailPrm := utils.Map{
        "goldenName": o.GoldenName.String,
        "email":      "",
    }
    if o.GoldenEmail.Valid {
        emailPrm["email"] = o.GoldenEmail.String
        go func() {
            cr.mailService.SendGoldenPearl(emailPrm)
        }()
        emailPrm["email"] = ""
    }
    if o.NokEmail.Valid {
        emailPrm["email"] = o.NokEmail.String
        go func() {
            cr.mailService.SendGoldenPearl(emailPrm)
        }()
        emailPrm["email"] = ""
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
// @Success 200 (object) []clubs.GuestLittleKidsMembership
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
// @Param isHome path string true "isHome"
// @Success 200 (object) []clubs.GoldenPearlActivity
// @Router /guest/clubs/goldenpearl/activity/all/mobile/{isHome} [get]
func (cr *GuestController) GetAllAppGoldenPearlActivities(c fiber.Ctx) error {
    isHome := c.Params("isHome")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.clubService.ListAppGoldenPearlActivities(isHome == "1", page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
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

    actvParticipation := make([]clubs.GoldenPearlActvParticipation, 0)
    for i := range data.GoldenActvParticipation {
        participation := data.GoldenActvParticipation[i]

        var o clubs.GoldenPearlActvParticipation
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

        emailPrm := utils.Map{
            "activityName": activity.GoldenActivityName,
            "memberName":   goldenMember.GoldenName,
            "email":        "",
        }
        if goldenMember.GoldenEmail.Valid {
            emailPrm["email"] = goldenMember.GoldenEmail.String
            go func() {
                cr.mailService.SendClubsEventRegistrationToMember(emailPrm)
            }()
            emailPrm["email"] = ""
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

func (cr *GuestController) GetAllAppPackages(c fiber.Ctx) error {
    // isHome := c.Params("isHome")
    // page := c.Query("_page", "1")
    // limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    return nil
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
    cartResult := make([]*userPackage.PackageCheckResult, 0)
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
