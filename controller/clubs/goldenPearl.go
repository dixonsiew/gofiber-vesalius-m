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
    "github.com/nleeper/goment"
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
        strings.EqualFold(data.GoldenDocNumber, strings.TrimSpace(data.NokDocNumber)) ||
        strings.EqualFold(data.GoldenDocType, utils.ClubsDocTypePassport) &&
        strings.EqualFold(data.NokDocType, utils.ClubsDocTypePassport) &&
        strings.EqualFold(data.GoldenDocNumber, strings.TrimSpace(data.NokDocNumber)) ||
        strings.EqualFold(data.GoldenDocType, utils.ClubsDocTypeBirthCert) &&
        strings.EqualFold(data.NokDocType, utils.ClubsDocTypeBirthCert) &&
        strings.EqualFold(data.GoldenDocNumber, strings.TrimSpace(data.NokDocNumber)) {
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
        "email":    "",
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
