package goldenPearl

import (
    "strconv"
    "strings"
    "vesaliusm/controller/clubs/shared"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    model "vesaliusm/model/clubs"
    "vesaliusm/service/clubs"
    "vesaliusm/service/mail"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/nleeper/goment"
)

type ClubsGoldenPearlController struct {
    clubService *clubs.ClubService
    mailService *mail.MailService
}

func NewClubsGoldenPearlController() *ClubsGoldenPearlController {
    return &ClubsGoldenPearlController{
        clubService: clubs.ClubSvc,
        mailService: mail.MailSvc,
    }
}

// CreateGoldenPearlMembership
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param request body dto.GoldenPearlMembershipDto true "GoldenPearlMembershipDto"
// @Success 200
// @Router /clubs/goldenpearl/membership [post]
func (cr *ClubsGoldenPearlController) CreateGoldenPearlMembership(c fiber.Ctx) error {
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
        "message": "Golden Pearl Membership created",
    })
}

// CreateGoldenPearlMembershipViaWebportal
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param request body dto.GoldenPearlMembershipDto true "GoldenPearlMembershipDto"
// @Success 200
// @Router /clubs/goldenpearl/membership/webadmin [post]
func (cr *ClubsGoldenPearlController) CreateGoldenPearlMembershipViaWebportal(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

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

    err = cr.clubService.SaveGoldenPearlMembershipViaWebportal(o, user.AdminId.Int64)
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
        "message": "Golden Pearl Membership created",
    })
}

// UpdateGoldenPearlMembership
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param        membershipId               path        string  true  "membershipId"
// @Param        request                    body        dto.GoldenPearlMembershipDto  true  "GoldenPearlMembershipDto"
// @Success 200
// @Router /clubs/goldenpearl/membership/webadmin/{membershipId} [put]
func (cr *ClubsGoldenPearlController) UpdateGoldenPearlMembership(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    membershipId := c.Params("membershipId")
    imembershipId, _ := strconv.ParseInt(membershipId, 10, 64)

    data := new(dto.GoldenPearlMembershipDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
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

    var o model.GoldenPearlMembership
    o.GoldenMembershipId = utils.NewInt64(imembershipId)
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
    o.NokPostCode = utils.NewNullString(data.NokPostCode)
    o.NokState = utils.NewNullString(data.NokState)
    o.NokCountryCode = utils.NewNullString(data.NokCountryCode)
    o.Relationship = utils.NewNullString(data.Relationship)
    o.PreferredLanguage = utils.NewNullString(data.PreferredLanguage)

    err = cr.clubService.UpdateGoldenPearlMembershipViaWebportal(o, user.AdminId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Golden Pearl Membership updated",
    })
}

// GetGoldenPearlMembershipById
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.GoldenPearlMembership
// @Router /clubs/goldenpearl/membership/{membershipId} [get]
func (cr *ClubsGoldenPearlController) GetGoldenPearlMembershipById(c fiber.Ctx) error {
    membershipId := c.Params("membershipId")
    imembershipId, _ := strconv.ParseInt(membershipId, 10, 64)

    o, err := cr.clubService.FindGoldenPearlMembershipByMembershipId(imembershipId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// GetAllAppGoldenPearlMemberships
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []clubs.GoldenPearlMembership
// @Router /clubs/goldenpearl/membership/all/mobile [get]
func (cr *ClubsGoldenPearlController) GetAllAppGoldenPearlMemberships(c fiber.Ctx) error {
    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    lx, err := cr.clubService.FindAllAppGoldenPearlMemberships(user.UserId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetAllExportGoldenPearlMembership
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.GoldenPearlMembership
// @Router /clubs/goldenpearl/membership/export/all [get]
func (cr *ClubsGoldenPearlController) GetAllExportGoldenPearlMembership(c fiber.Ctx) error {
    return nil
}

// GetSearchExportGoldenPearlMembership
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         keyword      body        dto.SearchKeyword2Dto false  "Search"
// @Success 200 {array} model.GoldenPearlMembership
// @Router /clubs/goldenpearl/membership/export/search [post]
func (cr *ClubsGoldenPearlController) GetSearchExportGoldenPearlMembership(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key2 := data.GetString("keyword2")
    if key != "" {
        key = "%" + key + "%"
    }
    if key2 != "" {
        key2 = "%" + key2 + "%"
    }

    return nil
}

// GetAllGoldenPearlMemberships
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.GoldenPearlMembership
// @Router /clubs/goldenpearl/membership/all [get]
func (cr *ClubsGoldenPearlController) GetAllGoldenPearlMemberships(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.clubService.ListGoldenPearlMemberships(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllGoldenPearlMembership
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         _page        query       string                false  "_page"  default:"1"
// @Param         _limit       query       string                false  "_limit" default:"10"
// @Param         keyword      body        dto.SearchKeyword2Dto false  "Search"
// @Success 200 {array} model.GoldenPearlMembership
// @Router /clubs/goldenpearl/membership/all [post]
func (cr *ClubsGoldenPearlController) SearchAllGoldenPearlMembership(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key2 := data.GetString("keyword2")
    if key != "" {
        key = "%" + key + "%"
    }
    if key2 != "" {
        key2 = "%" + key2 + "%"
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.clubService.ListGoldenPearlMembershipByKeyword(key, key2, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllUserGoldenPearlActivities
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.GoldenPearlMyActivity
// @Router /clubs/goldenpearl/my-activity/all [get]
func (cr *ClubsGoldenPearlController) GetAllUserGoldenPearlActivities(c fiber.Ctx) error {
    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    lx, err := cr.clubService.FindAllUserGoldenPearlActivities(user.UserId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// ParticipateGoldenPearlActivity
//
// @Tags Clubs
// @Produce json
// @Param data body dto.GoldenPearlActvParticipationDto true "GoldenPearlActvParticipationDto"
// @Success 200
// @Router /clubs/goldenpearl/activity/participate [post]
func (cr *ClubsGoldenPearlController) ParticipateGoldenPearlActivity(c fiber.Ctx) error {
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

// CreateGoldenPearlActivity
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param request body dto.GoldenPearlActivityDto true "GoldenPearlActivityDto"
// @Success 200
// @Router /clubs/goldenpearl/activity [post]
func (cr ClubsGoldenPearlController) CreateGoldenPearlActivity(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    data := new(dto.GoldenPearlActivityDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    maxParticipant, _ := strconv.ParseInt(data.ActivityMaxParticipant, 10, 32)

    var o model.GoldenPearlActivity
    o.GoldenActivityCode = utils.NewNullString(data.GoldenActivityCode)
    o.GoldenActivityName = utils.NewNullString(data.GoldenActivityName)
    o.GoldenActivityDesc = utils.NewNullString(data.GoldenActivityDesc)
    o.GoldenActivityImage = utils.NewNullString(data.GoldenActivityImage)
    o.ActivityStartDateTime = utils.NewNullString(data.ActivityStartDateTime)
    o.ActivityEndDateTime = utils.NewNullString(data.ActivityEndDateTime)
    o.ActivityMaxParticipant = utils.NewInt32(int32(maxParticipant))
    o.ActivityTnc = utils.NewNullString(data.ActivityTnc)
    o.ActivityDisplayOrder = utils.NewNullString(data.ActivityDisplayOrder)

    err = cr.clubService.SaveGoldenPearlActivity(o, user.AdminId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Golden Pearl Activity created",
    })
}

// UpdateGoldenPearlActivity
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param        activityId      path      string                             true  "activityId"
// @Param        request         body      dto.GoldenPearlActivityDto true  "GoldenPearlActivityDto"
// @Success 200
// @Router /clubs/goldenpearl/activity/{activityId} [put]
func (cr *ClubsGoldenPearlController) UpdateGoldenPearlActivity(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    activityId := c.Params("activityId")
    iactivityId, _ := strconv.ParseInt(activityId, 10, 64)
    data := new(dto.GoldenPearlActivityDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    maxParticipant, _ := strconv.ParseInt(data.ActivityMaxParticipant, 10, 32)

    var o model.GoldenPearlActivity
    o.GoldenActivityId = utils.NewInt64(iactivityId)
    o.GoldenActivityCode = utils.NewNullString(data.GoldenActivityCode)
    o.GoldenActivityName = utils.NewNullString(data.GoldenActivityName)
    o.GoldenActivityDesc = utils.NewNullString(data.GoldenActivityDesc)
    o.GoldenActivityImage = utils.NewNullString(data.GoldenActivityImage)
    o.ActivityStartDateTime = utils.NewNullString(data.ActivityStartDateTime)
    o.ActivityEndDateTime = utils.NewNullString(data.ActivityEndDateTime)
    o.ActivityMaxParticipant = utils.NewInt32(int32(maxParticipant))
    o.ActivityTnc = utils.NewNullString(data.ActivityTnc)
    o.ActivityDisplayOrder = utils.NewNullString(data.ActivityDisplayOrder)

    err = cr.clubService.UpdateGoldenPearlActivity(o, user.AdminId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Golden Pearl Activity updated",
    })
}

// GetAllExportGoldenPearlActivity
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.GoldenPearlActivity
// @Router /clubs/goldenpearl/activity/export/all [get]
func (cr *ClubsGoldenPearlController) GetAllExportGoldenPearlActivity(c fiber.Ctx) error {
    return nil
}

// GetSearchExportGoldenPearlActivity
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param keyword body dto.SearchKeyword3Dto false "Search"
// @Success 200 {array} model.GoldenPearlActivity
// @Router /clubs/goldenpearl/activity/export/search [post]
func (cr *ClubsGoldenPearlController) GetSearchExportGoldenPearlActivity(c fiber.Ctx) error {
    var data fiber.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }
    
    return nil
}

// GetAllGoldenPearlActivities
//
// @Tags Clubs
// @Produce json
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.GoldenPearlActivity
// @Router /clubs/goldenpearl/activity/all [get]
func (cr *ClubsGoldenPearlController) GetAllGoldenPearlActivities(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.clubService.ListGoldenPearlActivities(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllAppGoldenPearlActivities
//
// @Tags Clubs
// @Produce json
// @Param         _page        query       string                false  "_page"  default:"1"
// @Param         _limit       query       string                false  "_limit" default:"10"
// @Param         isHome       path        string                true   "isHome"
// @Success 200 {array} model.GoldenPearlActivity
// @Router /clubs/goldenpearl/activity/all/mobile/{isHome} [get]
func (cr *ClubsGoldenPearlController) GetAllAppGoldenPearlActivities(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    isHome := c.Params("isHome")
    m, err := cr.clubService.ListAppGoldenPearlActivities(isHome == "1", page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllGoldenPearlActivities
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         _page        query       string                false  "_page"  default:"1"
// @Param         _limit       query       string                false  "_limit" default:"10"
// @Param         keyword      body        dto.SearchKeyword3Dto false  "Search"
// @Success 200 {array} model.GoldenPearlActivity
// @Router /clubs/goldenpearl/activity/all [post]
func (cr *ClubsGoldenPearlController) SearchAllGoldenPearlActivities(c fiber.Ctx) error {
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
        if _, err := goment.New(key2, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong start date format")
        }
    }
    if key3 != "" {
        if _, err := goment.New(key3, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong end date format")
        }
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.clubService.ListGoldenPearlActivitiesByKeyword(key, key2, key3, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllExportGoldenPearlAttendees
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param        activityId      path      string                             true  "activityId"
// @Success 200 {array} model.GoldenPearlMembership
// @Router /clubs/goldenpearl/activity/attendees/{activityId}/export/all [get]
func (cr *ClubsGoldenPearlController) GetAllExportGoldenPearlAttendees(c fiber.Ctx) error {
    return nil
}

// GetSearchExportGoldenPearlAttendees
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param        activityId      path      string                             true  "activityId"
// @Success 200 {array} model.GoldenPearlMembership
// @Router /clubs/goldenpearl/activity/attendees/{activityId}/export/search [post]
func (cr *ClubsGoldenPearlController) GetSearchExportGoldenPearlAttendees(c fiber.Ctx) error {
    return nil
}

// GetGoldenPearlActivityAttendeesById
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         activityId      path     string                true   "activityId"
// @Param         _page           query    string                false  "_page"  default:"1"
// @Param         _limit          query    string                false  "_limit" default:"10"
// @Success 200 {array} model.GoldenPearlMembership
// @Router /clubs/goldenpearl/activity/attendees/{activityId} [get]
func (cr *ClubsGoldenPearlController) GetGoldenPearlActivityAttendeesById(c fiber.Ctx) error {
    activityId := c.Params("activityId")
    iactivityId, _ := strconv.ParseInt(activityId, 10, 64)
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.clubService.ListGoldenPearlActivityAttendees(iactivityId, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllGoldenPearlAttendees
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         activityId      path     string                true   "activityId"
// @Param         _page           query    string                false  "_page"  default:"1"
// @Param         _limit          query    string                false  "_limit" default:"10"
// @Success 200 {array} model.GoldenPearlMembership
// @Router /clubs/goldenpearl/activity/attendees/{activityId} [post]
func (cr *ClubsGoldenPearlController) SearchAllGoldenPearlAttendees(c fiber.Ctx) error {
    activityId := c.Params("activityId")
    iactivityId, _ := strconv.ParseInt(activityId, 10, 64)
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))

    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key2 := data.GetString("keyword2")
    if key != "" {
        key = "%" + key + "%"
    }
    if key2 != "" {
        key2 = "%" + key2 + "%"
    }

    m, err := cr.clubService.ListGoldenPearlAttendeesByKeyword(iactivityId, key, key2, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetGoldenPearlActivityNameById
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         activityId   path        string                true   "activityId"
// @Success 200 {object} model.GoldenPearlActivity
// @Router /clubs/goldenpearl/activity/name/{activityId} [get]
func (cr *ClubsGoldenPearlController) GetGoldenPearlActivityNameById(c fiber.Ctx) error {
    activityId := c.Params("activityId")
    iactivityId, _ := strconv.ParseInt(activityId, 10, 64)
    o, err := cr.clubService.FindGoldenPearlActivityNameById(iactivityId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// GetGoldenPearlActivityById
//
// @Tags Clubs
// @Produce json
// @Param         activityId   path        string                true   "activityId"
// @Success 200 {object} model.GoldenPearlActivity
// @Router /clubs/goldenpearl/activity/{activityId} [get]
func (cr *ClubsGoldenPearlController) GetGoldenPearlActivityById(c fiber.Ctx) error {
    activityId := c.Params("activityId")
    iactivityId, _ := strconv.ParseInt(activityId, 10, 64)
    o, err := cr.clubService.FindGoldenPearlActivitiesByActivityId(iactivityId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// CreateGoldenPearlAboutUs
//
// @Tags Clubs
// @Produce json
// @Param request body dto.GoldenPearlAboutUsDto true "GoldenPearlAboutUsDto"
// @Success 200
// @Router /clubs/goldenpearl/about-us [post]
func (cr *ClubsGoldenPearlController) CreateGoldenPearlAboutUs(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    data := new(dto.GoldenPearlAboutUsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.clubService.ExistsGoldenPearlAboutUs()
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, "Already setup Golden Pearl - About Us previously")
    }

    var o model.GoldenPearlAboutUs
    o.GoldenClubTitle = utils.NewNullString(data.GoldenClubTitle)
    o.GoldenClubDesc = utils.NewNullString(data.GoldenClubDesc)
    o.GoldenClubImage = utils.NewNullString(data.GoldenClubImage)
    o.GoldenClubTnc = utils.NewNullString(data.GoldenClubTnc)
    o.GoldenClubExtLink = utils.NewNullString(data.GoldenClubExtLink)

    goldenClubId, err :=cr.clubService.SaveGoldenPearlAboutUs(o, user.AdminId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message":      "Golden Pearl About Us created",
        "golden_club_id": goldenClubId,
    })
}

// UpdateGoldenPearlAboutUs
//
// @Tags Clubs
// @Produce json
// @Param goldenPearlId path int true "goldenPearlId"
// @Param request body dto.GoldenPearlAboutUsDto true "GoldenPearlAboutUsDto"
// @Success 200
// @Router /clubs/goldenpearl/about-us/{goldenPearlId} [put]
func (cr *ClubsGoldenPearlController) UpdateGoldenPearlAboutUs(c fiber.Ctx) error {
    goldenPearlId := c.Params("goldenPearlId")
    igoldenPearlId, _ := strconv.ParseInt(goldenPearlId, 10, 64)

    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    data := new(dto.GoldenPearlAboutUsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.clubService.ExistsGoldenPearlAboutUs()
    if err != nil {
        return err
    }

    if !b {
        return fiber.NewError(fiber.StatusBadRequest, "Golden Pearl - About us does not exist")
    }

    var o model.GoldenPearlAboutUs
    o.GoldenClubId = utils.NewInt64(igoldenPearlId)
    o.GoldenClubTitle = utils.NewNullString(data.GoldenClubTitle)
    o.GoldenClubDesc = utils.NewNullString(data.GoldenClubDesc)
    o.GoldenClubImage = utils.NewNullString(data.GoldenClubImage)
    o.GoldenClubTnc = utils.NewNullString(data.GoldenClubTnc)
    o.GoldenClubExtLink = utils.NewNullString(data.GoldenClubExtLink)

    err = cr.clubService.UpdateGoldenPearlAboutUs(o, user.AdminId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Golden Pearl About Us updated",
    })
}

// GetGoldenPearlAboutUs
//
// @Tags Clubs
// @Produce json
// @Success 200 {object} model.GoldenPearlAboutUs
// @Router /clubs/goldenpearl/about-us [get]
func (cr *ClubsGoldenPearlController) GetGoldenPearlAboutUs(c fiber.Ctx) error {
    o, err := cr.clubService.FindGoldenPearlAboutUs()
    if err != nil {
        return err
    }

    return c.JSON(o)
}
