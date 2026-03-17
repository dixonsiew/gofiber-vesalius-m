package clubs

import (
    "fmt"
    "strconv"
    "strings"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    "vesaliusm/model/clubs"
    clubsSvc "vesaliusm/service/clubs"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

type ClubsController struct {
    clubService *clubsSvc.ClubService
}

func NewClubsController() *ClubsController {
    return &ClubsController{
        clubService: clubsSvc.ClubSvc,
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
func (cr *ClubsController) CreateLittleKidsMembership(c fiber.Ctx) error {
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

    eligibleAge := littleKidsEligibleAge(data.KidsDob)
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

    emailPrm := map[string]interface{}{
        "kidsName": o.KidsName.String,
        "email":    "",
    }
    if o.KidsEmail.Valid {
        emailPrm["email"] = o.KidsEmail.String

        emailPrm["email"] = ""
    }
    if o.GuardianEmail.Valid {
        emailPrm["email"] = o.GuardianEmail.String

        emailPrm["email"] = ""
    }

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
func (cr *ClubsController) CreateLittleKidsMembershipViaWebportal(c fiber.Ctx) error {
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
    
    eligibleAge := littleKidsEligibleAge(data.KidsDob)
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
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Kids Document Type")
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

    emailPrm := map[string]interface{}{
        "kidsName": o.KidsName.String,
        "email":    "",
    }
    if o.KidsEmail.Valid {
        emailPrm["email"] = o.KidsEmail.String

        emailPrm["email"] = ""
    }
    if o.GuardianEmail.Valid {
        emailPrm["email"] = o.GuardianEmail.String

        emailPrm["email"] = ""
    }

    return c.JSON(fiber.Map{
        "message": "Little Explorers Kids Membership created",
    })
}

func (cr *ClubsController) UpdateLittleKidsMembership(c fiber.Ctx) error {
    return nil
}

// GetLittleKidsMembershipById
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/membership/{membershipId} [get]
func (cr *ClubsController) GetLittleKidsMembershipById(c fiber.Ctx) error {
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
func (cr *ClubsController) GetAllAppLittleKidsMemberships(c fiber.Ctx) error {
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
func (cr *ClubsController) GetAllExportLittleKidsMembership(c fiber.Ctx) error {
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
func (cr *ClubsController) GetSearchExportLittleKidsMembership(c fiber.Ctx) error {
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

// GetAllLittleKidsMemberships
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/membership/all [get]
func (cr *ClubsController) GetAllLittleKidsMemberships(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := cr.clubService.ListLittleKidsMemberships(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllLittleKidsMembership
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Param keyword body dto.SearchKeyword2Dto false "Search"
// @Success 200 {array} clubs.LittleExplorersKidsMembership
// @Router /clubs/littlekids/membership/all [post]
func (cr *ClubsController) SearchAllLittleKidsMembership(c fiber.Ctx) error {
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
    limit := c.Query("_limit", "10")
    m, err := cr.clubService.ListLittleKidsMembershipByKeyword(key, key2, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// GetAllUserLittleKidsActivities
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {array} clubs.LittleExplorersKidsMyActivity
// @Router /clubs/littlekids/my-activity/all [get]
func (cr *ClubsController) GetAllUserLittleKidsActivities(c fiber.Ctx) error {
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
func (cr *ClubsController) ParticipateLittleKidsActivity(c fiber.Ctx) error {
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

        emailPrm := map[string]string{
            "activityName": activity.KidsActivityName.String,
            "memberName":   kidsMember.KidsName.String,
            "email": "",
        }
        if kidsMember.KidsEmail.Valid {
            emailPrm["email"] = kidsMember.KidsEmail.String
            emailPrm["email"] = ""
        }
    }

    err := cr.clubService.ParticipateLittleKidsActivity(actvParticipation)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Activity Participate successfully",
    })
}

func (cr *ClubsController) CreateLittleKidsActivity(c fiber.Ctx) error {
    
}

// GetGoldenPearlAboutUs
//
// @Tags Clubs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} clubs.GoldenPearlAboutUs
// @Router /clubs/goldenpearl/about-us [get]
func (cr *ClubsController) GetGoldenPearlAboutUs(c fiber.Ctx) error {
    o, err := cr.clubService.FindGoldenPearlAboutUs()
    if err != nil {
        return err
    }

    return c.JSON(o)
}
