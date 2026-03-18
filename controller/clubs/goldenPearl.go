package clubs

import (
    // "fmt"
    "strconv"
    "strings"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    "vesaliusm/model/clubs"
    clubsSvc "vesaliusm/service/clubs"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    // "github.com/nleeper/goment"
)

type ClubsGoldenPearlController struct {
    clubService *clubsSvc.ClubService
}

func NewClubsGoldenPearlController() *ClubsGoldenPearlController {
    return &ClubsGoldenPearlController{
        clubService: clubsSvc.ClubSvc,
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

    eligibleAge := goldenPearlEligibleAge(data.GoldenDob)
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

    emailPrm := map[string]interface{}{
        "goldenName": o.GoldenName.String,
        "email":      "",
    }
    if o.GoldenEmail.Valid {
        emailPrm["email"] = o.GoldenEmail.String

        emailPrm["email"] = ""
    }
    if o.NokEmail.Valid {
        emailPrm["email"] = o.NokEmail.String

        emailPrm["email"] = ""
    }

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

    eligibleAge := goldenPearlEligibleAge(data.GoldenDob)
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

    err = cr.clubService.SaveGoldenPearlMembershipViaWebportal(o, user.AdminId.Int64)
    if err != nil {
        return err
    }

    emailPrm := map[string]interface{}{
        "goldenName": o.GoldenName.String,
        "email":      "",
    }
    if o.GoldenEmail.Valid {
        emailPrm["email"] = o.GoldenEmail.String

        emailPrm["email"] = ""
    }
    if o.NokEmail.Valid {
        emailPrm["email"] = o.NokEmail.String

        emailPrm["email"] = ""
    }

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

    eligibleAge := goldenPearlEligibleAge(data.GoldenDob)
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
// @Success 200 {object} clubs.GoldenPearlMembership
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
// @Success 200 {array} clubs.GoldenPearlMembership
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
// @Success 200 {array} clubs.GoldenPearlMembership
// @Router /clubs/goldenpearl/membership/export/search [post]
func (cr *ClubsGoldenPearlController) GetSearchExportGoldenPearlMembership(c fiber.Ctx) error {
    var data fiber.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    var (
        key  string
        key2 string
    )
    if keyword, ok := data["keyword"]; ok {
        key = keyword.(string)
        if key != "" {
            key = "%" + key + "%"
        }
    }
    if keyword, ok := data["keyword2"]; ok {
        key2 = keyword.(string)
        if key2 != "" {
            key2 = "%" + key2 + "%"
        }
    }

    return nil
}

// GetAllGoldenPearlMemberships
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} clubs.GoldenPearlMembership
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
// @Success 200 {array} clubs.GoldenPearlMembership
// @Router /clubs/goldenpearl/membership/all [post]
func (cr *ClubsGoldenPearlController) SearchAllGoldenPearlMembership(c fiber.Ctx) error {
    var data fiber.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    var (
        key  string
        key2 string
    )
    if keyword, ok := data["keyword"]; ok {
        key = keyword.(string)
        if key != "" {
            key = "%" + key + "%"
        }
    }
    if keyword, ok := data["keyword2"]; ok {
        key2 = keyword.(string)
        if key2 != "" {
            key2 = "%" + key2 + "%"
        }
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
// @Success 200 {array} clubs.GoldenPearlMyActivity
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

        emailPrm := map[string]interface{}{
            "activityName": activity.GoldenActivityName,
            "memberName":   goldenMember.GoldenName,
            "email":        "",
        }

        if goldenMember.GoldenEmail.Valid {
            emailPrm["email"] = goldenMember.GoldenEmail.String

            emailPrm["email"] = ""
        }
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

    var o clubs.GoldenPearlActivity
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

// GetGoldenPearlAboutUs
//
// @Tags Clubs
// @Produce json
// @Success 200 {object} clubs.GoldenPearlAboutUs
// @Router /clubs/goldenpearl/about-us [get]
func (cr *ClubsGoldenPearlController) GetGoldenPearlAboutUs(c fiber.Ctx) error {
    o, err := cr.clubService.FindGoldenPearlAboutUs()
    if err != nil {
        return err
    }

    return c.JSON(o)
}
