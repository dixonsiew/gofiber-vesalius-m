package littleKids

import (
    "strconv"
    "strings"
    "vesaliusm/controller/clubs/shared"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    "vesaliusm/model/clubs"
    clubsSvc "vesaliusm/service/clubs"
    "vesaliusm/service/mail"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/nleeper/goment"
)

type ClubsLittleKidsController struct {
    clubService *clubsSvc.ClubService
    mailService *mail.MailService
}

func NewClubsLittleKidsController() *ClubsLittleKidsController {
    return &ClubsLittleKidsController{
        clubService: clubsSvc.ClubSvc,
        mailService: mail.MailSvc,
    }
}

// CreateLittleKidsMembership
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param request body dto.LittleExplorersKidsMembershipDto true "LittleExplorersKidsMembershipDto"
// @Success 200
// @Router /clubs/littlekids/membership [post]
func (cr *ClubsLittleKidsController) CreateLittleKidsMembership(c fiber.Ctx) error {
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
        // remove the dash from the NRIC
        data.KidsDocNumber = strings.ReplaceAll(data.KidsDocNumber, "-", "")
    }

    if strings.EqualFold(data.GuardianDocType, utils.ClubsDocTypeNRIC) {
        // remove the dash from the NRIC
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
        "message": "Little Explorers Kids Membership created",
    })
}

// CreateLittleKidsMembershipViaWebportal
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param request body dto.LittleExplorersKidsMembershipDto true "LittleExplorersKidsMembershipDto"
// @Success 200
// @Router /clubs/littlekids/membership/webadmin [post]
func (cr *ClubsLittleKidsController) CreateLittleKidsMembershipViaWebportal(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

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
        // remove the dash from the NRIC
        data.KidsDocNumber = strings.ReplaceAll(data.KidsDocNumber, "-", "")
    }

    if strings.EqualFold(data.GuardianDocType, utils.ClubsDocTypeNRIC) {
        // remove the dash from the NRIC
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

    err = cr.clubService.SaveLittleKidsMembershipViaWebportal(o, user.AdminId.Int64)
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
        "message": "Little Explorers Kids Membership created",
    })
}

// UpdateLittleKidsMembership
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param        membershipId               path        string  true  "membershipId"
// @Param        data                       body        dto.LittleExplorersKidsMembershipDto  true  "LittleExplorersKidsMembershipDto"
// @Success 200
// @Router /clubs/littlekids/membership/webadmin/{membershipId} [put]
func (cr *ClubsLittleKidsController) UpdateLittleKidsMembership(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    membershipId := c.Params("membershipId")
    imembershipId, _ := strconv.ParseInt(membershipId, 10, 64)

    data := new(dto.LittleExplorersKidsMembershipDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
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
        // remove the dash from the NRIC
        data.KidsDocNumber = strings.ReplaceAll(data.KidsDocNumber, "-", "")
    }

    if strings.EqualFold(data.GuardianDocType, utils.ClubsDocTypeNRIC) {
        // remove the dash from the NRIC
        data.GuardianDocNumber = strings.ReplaceAll(data.GuardianDocNumber, "-", "")
    }

    var o clubs.LittleExplorersKidsMembership
    o.KidsMembershipId = utils.NewInt64(imembershipId)
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
    o.GuardianPostCode = utils.NewNullString(data.GuardianPostCode)
    o.GuardianState = utils.NewNullString(data.GuardianState)
    o.GuardianCountryCode = utils.NewNullString(data.GuardianCountryCode)
    o.Relationship = utils.NewNullString(data.Relationship)
    o.PreferredLanguage = utils.NewNullString(data.PreferredLanguage)

    err = cr.clubService.UpdateLittleKidsMembershipViaWebportal(o, user.AdminId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Little Explorers Kids Membership updated",
    })
}

// GetLittleKidsMembershipById
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param        membershipId               path        string  true  "membershipId"
// @Success 200 {object} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/membership/{membershipId} [get]
func (cr *ClubsLittleKidsController) GetLittleKidsMembershipById(c fiber.Ctx) error {
    membershipId := c.Params("membershipId")
    imembershipId, _ := strconv.ParseInt(membershipId, 10, 64)
    o, err := cr.clubService.FindLittleKidsMembershipByMembershipId(imembershipId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// GetAllAppLittleKidsMemberships
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/membership/all/mobile [get]
func (cr *ClubsLittleKidsController) GetAllAppLittleKidsMemberships(c fiber.Ctx) error {
    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    lx, err := cr.clubService.FindAllAppLittleKidsMemberships(user.UserId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetAllExportLittleKidsMembership
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/membership/export/all [get]
func (cr *ClubsLittleKidsController) GetAllExportLittleKidsMembership(c fiber.Ctx) error {
    return nil
}

// GetSearchExportLittleKidsMembership
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param keyword body dto.SearchKeyword2Dto false "Search"
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/membership/export/search [post]
func (cr *ClubsLittleKidsController) GetSearchExportLittleKidsMembership(c fiber.Ctx) error {
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

// GetAllLittleKidsMemberships
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/membership/all [get]
func (cr *ClubsLittleKidsController) GetAllLittleKidsMemberships(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.clubService.ListLittleKidsMemberships(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllLittleKidsMembership
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         _page        query       string                false  "_page"  default:"1"
// @Param         _limit       query       string                false  "_limit" default:"10"
// @Param         keyword      body        dto.SearchKeyword2Dto false  "Search"
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/membership/all [post]
func (cr *ClubsLittleKidsController) SearchAllLittleKidsMembership(c fiber.Ctx) error {
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
    m, err := cr.clubService.ListLittleKidsMembershipByKeyword(key, key2, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllUserLittleKidsActivities
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} clubs.LittleExplorersKidsMyActivity
// @Router /clubs/littlekids/my-activity/all [get]
func (cr *ClubsLittleKidsController) GetAllUserLittleKidsActivities(c fiber.Ctx) error {
    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    lx, err := cr.clubService.FindAllUserLittleKidsActivities(user.UserId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// ParticipateLittleKidsActivity
//
// @Tags Clubs
// @Produce json
// @Param request body dto.LittleExplorersKidsActvParticipationDto true "LittleExplorersKidsActvParticipationDto"
// @Success 200
// @Router /clubs/littlekids/activity/participate [post]
func (cr *ClubsLittleKidsController) ParticipateLittleKidsActivity(c fiber.Ctx) error {
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

// ParticipateLittleKidsActivity
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param request body dto.LittleExplorersKidsActivityDto true "LittleExplorersKidsActivityDto"
// @Success 200
// @Router /clubs/littlekids/activity [post]
func (cr *ClubsLittleKidsController) CreateLittleKidsActivity(c fiber.Ctx) error {
    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    data := new(dto.LittleExplorersKidsActivityDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    maxParticipant, _ := strconv.ParseInt(data.ActivityMaxParticipant, 10, 32)

    var o clubs.LittleExplorersKidsActivity
    o.KidsActivityCode = utils.NewNullString(data.KidsActivityCode)
    o.KidsActivityName = utils.NewNullString(data.KidsActivityName)
    o.KidsActivityDesc = utils.NewNullString(data.KidsActivityDesc)
    o.KidsActivityImage = utils.NewNullString(data.KidsActivityImage)
    o.ActivityStartDateTime = utils.NewNullString(data.ActivityStartDateTime)
    o.ActivityEndDateTime = utils.NewNullString(data.ActivityEndDateTime)
    o.ActivityMaxParticipant = utils.NewInt32(int32(maxParticipant))
    o.ActivityTnc = utils.NewNullString(data.ActivityTnc)
    o.ActivityDisplayOrder = utils.NewNullString(data.ActivityDisplayOrder)

    err = cr.clubService.SaveLittleKidsActivity(o, user.UserId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Little Kids Activity created",
    })
}

// UpdateLittleKidsActivity
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param        activityId      path      string                             true  "activityId"
// @Param        request         body      dto.LittleExplorersKidsActivityDto true  "LittleExplorersKidsActivityDto"
// @Success 200
// @Router /clubs/littlekids/activity/{activityId} [put]
func (cr *ClubsLittleKidsController) UpdateLittleKidsActivity(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    activityId := c.Params("activityId")
    iactivityId, _ := strconv.ParseInt(activityId, 10, 64)
    data := new(dto.LittleExplorersKidsActivityDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    maxParticipant, _ := strconv.ParseInt(data.ActivityMaxParticipant, 10, 32)

    var o clubs.LittleExplorersKidsActivity
    o.KidsActivityId = utils.NewInt64(iactivityId)
    o.KidsActivityCode = utils.NewNullString(data.KidsActivityCode)
    o.KidsActivityName = utils.NewNullString(data.KidsActivityName)
    o.KidsActivityDesc = utils.NewNullString(data.KidsActivityDesc)
    o.KidsActivityImage = utils.NewNullString(data.KidsActivityImage)
    o.ActivityStartDateTime = utils.NewNullString(data.ActivityStartDateTime)
    o.ActivityEndDateTime = utils.NewNullString(data.ActivityEndDateTime)
    o.ActivityMaxParticipant = utils.NewInt32(int32(maxParticipant))
    o.ActivityTnc = utils.NewNullString(data.ActivityTnc)
    o.ActivityDisplayOrder = utils.NewNullString(data.ActivityDisplayOrder)

    err = cr.clubService.UpdateLittleKidsActivity(o, user.AdminId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Little Kids Activity updated",
    })
}

// GetAllExportLittleKidsActivity
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} clubs.LittleExplorersKidsActivity
// @Router /clubs/littlekids/activity/export/all [get]
func (cr *ClubsLittleKidsController) GetAllExportLittleKidsActivity(c fiber.Ctx) error {
    return nil
}

// GetSearchExportLittleKidsActivity
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         keyword      body        dto.SearchKeyword3Dto false  "Search"
// @Success 200 {array} clubs.LittleExplorersKidsActivity
// @Router /clubs/littlekids/activity/export/search [post]
func (cr *ClubsLittleKidsController) GetSearchExportLittleKidsActivity(c fiber.Ctx) error {
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

    return nil
}

// GetAllLittleKidsActivities
//
// @Tags Clubs
// @Produce json
// @Param         _page        query       string                false  "_page"  default:"1"
// @Param         _limit       query       string                false  "_limit" default:"10"
// @Success 200 {array} clubs.LittleExplorersKidsActivity
// @Router /clubs/littlekids/activity/all [get]
func (cr *ClubsLittleKidsController) GetAllLittleKidsActivities(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.clubService.ListLittleKidsActivities(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllAppLittleKidsActivities
//
// @Tags Clubs
// @Produce json
// @Param         _page        query       string                false  "_page"  default:"1"
// @Param         _limit       query       string                false  "_limit" default:"10"
// @Param         isHome       path        string                true   "isHome"
// @Success 200 {array} clubs.LittleExplorersKidsActivity
// @Router /clubs/littlekids/activity/all/mobile/{isHome} [get]
func (cr *ClubsLittleKidsController) GetAllAppLittleKidsActivities(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    isHome := c.Params("isHome")
    m, err := cr.clubService.ListAppLittleKidsActivities(isHome == "1", page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllLittleKidsActivities
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         _page        query       string                false  "_page"  default:"1"
// @Param         _limit       query       string                false  "_limit" default:"10"
// @Param         keyword      body        dto.SearchKeyword3Dto false  "Search"
// @Success 200 {array} clubs.LittleExplorersKidsActivity
// @Router /clubs/littlekids/activity/all [post]
func (cr *ClubsLittleKidsController) SearchAllLittleKidsActivities(c fiber.Ctx) error {
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
    m, err := cr.clubService.ListLittleKidsActivitiesByKeyword(key, key2, key3, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllExportLittleKidsAttendees
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param        activityId      path      string                             true  "activityId"
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/activity/attendees/{activityId}/export/all [get]
func (cr *ClubsLittleKidsController) GetAllExportLittleKidsAttendees(c fiber.Ctx) error {
    return nil
}

// GetSearchExportLittleKidsAttendees
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param        activityId      path      string                             true  "activityId"
// @Param        keyword         body      dto.SearchKeyword2Dto              false  "Search"
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/activity/attendees/:activityId/export/search [post]
func (cr *ClubsLittleKidsController) GetSearchExportLittleKidsAttendees(c fiber.Ctx) error {
    return nil
}

// GetLittleKidsActivityAttendeesById
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         activityId      path     string                true   "activityId"
// @Param         _page           query    string                false  "_page"  default:"1"
// @Param         _limit          query    string                false  "_limit" default:"10"
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/activity/attendees/{activityId} [get]
func (cr *ClubsLittleKidsController) GetLittleKidsActivityAttendeesById(c fiber.Ctx) error {
    activityId := c.Params("activityId")
    iactivityId, _ := strconv.ParseInt(activityId, 10, 64)
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.clubService.ListLittleKidsActivityAttendees(iactivityId, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllLittleKidsAttendees
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         activityId   path        string                true   "activityId"
// @Param         _page        query       string                false  "_page"  default:"1"
// @Param         _limit       query       string                false  "_limit" default:"10"
// @Param         keyword      body        dto.SearchKeyword2Dto false  "Search"
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/activity/attendees/{activityId} [post]
func (cr *ClubsLittleKidsController) SearchAllLittleKidsAttendees(c fiber.Ctx) error {
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

    m, err := cr.clubService.ListLittleKidsAttendeesByKeyword(iactivityId, key, key2, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetLittleKidsActivityNameById
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         activityId   path        string                true   "activityId"
// @Success 200 {object} clubs.LittleExplorersKidsActivity
// @Router /clubs/littlekids/activity/name/{activityId} [get]
func (cr *ClubsLittleKidsController) GetLittleKidsActivityNameById(c fiber.Ctx) error {
    activityId := c.Params("activityId")
    iactivityId, _ := strconv.ParseInt(activityId, 10, 64)
    o, err := cr.clubService.FindLittleKidsActivityNameById(iactivityId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// GetLittleKidsActivityById
//
// @Tags Clubs
// @Produce json
// @Param         activityId   path        string                true   "activityId"
// @Success 200
// @Router /clubs/littlekids/activity/{activityId} [get]
func (cr *ClubsLittleKidsController) GetLittleKidsActivityById(c fiber.Ctx) error {
    activityId := c.Params("activityId")
    iactivityId, _ := strconv.ParseInt(activityId, 10, 64)
    o, err := cr.clubService.FindLittleKidsActivitiesByActivityId(iactivityId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// CreateLittleKidsAboutUs
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param request body dto.LittleExplorersKidsAboutUsDto true "LittleExplorersKidsAboutUsDto"
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/about-us [post]
func (cr *ClubsLittleKidsController) CreateLittleKidsAboutUs(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    data := new(dto.LittleExplorersKidsAboutUsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.clubService.ExistsLittleKidsAboutUs()
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, "Already setup Little Explorers Kids - About Us previously")
    }

    var o clubs.LittleExplorersKidsAboutUs
    o.KidsClubTitle = utils.NewNullString(data.KidsClubTitle)
    o.KidsClubDesc = utils.NewNullString(data.KidsClubDesc)
    o.KidsClubImage = utils.NewNullString(data.KidsClubImage)
    o.KidsClubTnc = utils.NewNullString(data.KidsClubTnc)
    o.KidsClubPartnerLink = utils.NewNullString(data.KidsClubPartnerLink)

    kidsClubId, err := cr.clubService.SaveLittleKidsAboutUs(o, user.AdminId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message":      "Little Kids About Us created",
        "kids_club_id": kidsClubId,
    })
}

// UpdateLittleKidsAboutUs
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param         kidsClubId      path        string                            true    "kidsClubId"
// @Param         request         body        dto.LittleExplorersKidsAboutUsDto true    "LittleExplorersKidsAboutUsDto"
// @Success 200
// @Router /clubs/littlekids/about-us/{kidsClubId} [put]
func (cr *ClubsLittleKidsController) UpdateLittleKidsAboutUs(c fiber.Ctx) error {
    kidsClubId := c.Params("kidsClubId")
    ikidsClubId, _ := strconv.ParseInt(kidsClubId, 10, 64)

    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    data := new(dto.LittleExplorersKidsAboutUsDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.clubService.ExistsLittleKidsAboutUs()
    if err != nil {
        return err
    }

    if !b {
        return fiber.NewError(fiber.StatusBadRequest, "Little Explorers Kids - About us does not exist")
    }

    var o clubs.LittleExplorersKidsAboutUs
    o.KidsClubId = utils.NewInt64(ikidsClubId)
    o.KidsClubTitle = utils.NewNullString(data.KidsClubTitle)
    o.KidsClubDesc = utils.NewNullString(data.KidsClubDesc)
    o.KidsClubImage = utils.NewNullString(data.KidsClubImage)
    o.KidsClubTnc = utils.NewNullString(data.KidsClubTnc)
    o.KidsClubPartnerLink = utils.NewNullString(data.KidsClubPartnerLink)

    err = cr.clubService.UpdateLittleKidsAboutUs(o, user.AdminId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Little Kids About Us updated",
    })
}

// GetLittleKidsAboutUs
//
// @Tags Clubs
// @Produce json
// @Success 200 {object} clubs.LittleExplorersKidsAboutUs
// @Router /clubs/littlekids/about-us [get]
func (cr *ClubsLittleKidsController) GetLittleKidsAboutUs(c fiber.Ctx) error {
    o, err := cr.clubService.FindLittleKidsAboutUs()
    if err != nil {
        return err
    }

    return c.JSON(o)
}
