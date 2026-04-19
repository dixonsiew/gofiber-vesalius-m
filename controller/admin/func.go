package admin

import (
    "errors"
    "fmt"
    "strings"
    "vesaliusm/config"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    "vesaliusm/model"
    "vesaliusm/utils"
    "vesaliusm/utils/constants"

    "github.com/gofiber/fiber/v3"
    "github.com/nleeper/goment"
)

func (cr *AdminController) signUpPatient(c fiber.Ctx, data *dto.NewSignupUserDto, appPatient *model.ApplicationUser) error {
    isExistsByPrn, err := cr.applicationUserService.ExistsByPRN(data.UserPrn)
    if err != nil {
        return err
    }
    if isExistsByPrn {
        return fiber.NewError(fiber.StatusBadRequest, "Sorry, an account with the provided PRN already exists. Please sign in to your existing account or contact our Customer Service for assistance at info@islandhospital.com")
    }
    switch appPatient.InactiveFlag.String {
    case "N":
        switch appPatient.SignInType.Int32 {
        case 1:
            return fiber.NewError(fiber.StatusBadRequest, "Sorry, an account with the provided mobile number already exists. Please sign in or use a different mobile number to register. Contact our Customer Service for assistance at info@islandhospital.com")
        case 2:
            return fiber.NewError(fiber.StatusBadRequest, "Sorry, an account with the provided email address already exists. Please sign in or use a different mobile number to register. Contact our Customer Service for assistance at info@islandhospital.com")
        }
    case "Y":
        vesPatient, ex, err := cr.vesaliusService.VesaliusGetPatientDataByPrn(data.UserPrn)
        if vesPatient == nil {
            return fiber.NewError(fiber.StatusBadRequest, "Incorrect PRN: The Patient PRN Number provided does not exist. Please retry")
        }

        if ex != nil {
            if ex.Code == "99" {
                return fiber.NewError(fiber.StatusBadRequest, "Duplicate Patient Profile found. Please contact customer service for assistance.")
            } else {
                return fiber.NewError(fiber.StatusBadRequest, ex.ToString())
            }
        }

        if err != nil {
            var e *fiber.Error
            if errors.As(err, &e) {
                if e.Code == fiber.StatusNoContent {
                    return fiber.NewError(fiber.StatusBadRequest, "Information provided does not match with hospital patient profile. Please retry.")
                }
            }
            return err
        }

        isPatientWithPrn := vesPatient.Prn == data.UserPrn
        if !isPatientWithPrn {
            return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Your patient profile (PRN: %s) is inactive. Please reach out to our customer service at +604-238 3388 for further action", data.UserPrn))
        }

        if len(vesPatient.Documents) > 0 && vesPatient.DOB != "" {
            checkDocument := true
            checkDob := true
            b := false

            for i := range vesPatient.Documents {
                doc := vesPatient.Documents[i]
                if strings.EqualFold(doc.Code, config.GetPatientDocumentCode()) {
                    if strings.TrimSpace(doc.Value) == data.UserPersonNumber {
                        b = true
                    }
                }
            }
            if !b {
                checkDocument = false
            }

            vesPatientDOB, _ := goment.New(strings.TrimSpace(vesPatient.DOB), "DD-MMM-YYYY")
            inputPatientDOB, _ := goment.New(strings.TrimSpace(data.UserDOB), "DD/MM/YYYY")
            if !inputPatientDOB.IsSame(vesPatientDOB) {
                checkDob = false
            }

            if !checkDocument && !checkDob {
                return fiber.NewError(fiber.StatusBadRequest, "Please verify your PRN, NRIC / Passport / Birth Cert and Date of Birth as they do not match our hospital records. For assistance, contact our Customer Service at the Front Desk or at info@islandhospital.com")
            }
        }

        patientDocIDValue := ""
        if len(vesPatient.Documents) > 0 {
            b := false
            for i := range vesPatient.Documents {
                doc := vesPatient.Documents[i]
                if strings.EqualFold(doc.Code, config.GetPatientDocumentCode()) {
                    patientDocIDValue = strings.TrimSpace(doc.Value)
                    if config.GetWSVesaliusConfig().NricWithDash == "N" {
                        data.UserPersonNumber = strings.ReplaceAll(data.UserPersonNumber, "-", "")
                    }
                    if patientDocIDValue == data.UserPersonNumber {
                        b = true
                        if doc.ExpireDate != "" {
                            patientDocExpiry, _ := goment.New(strings.TrimSpace(doc.ExpireDate), "DD-MMM-YYYY")
                            currentDate, _ := goment.New()
                            if patientDocExpiry.IsBefore(currentDate) {
                                return fiber.NewError(fiber.StatusBadRequest, "The NRIC / Passport / Birth Cert provided is not valid. Please confirm your details at the Front Desk or contact Customer Service at info@islandhospital.com")
                            }
                        }
                    }
                }
            }
            if !b {
                return fiber.NewError(fiber.StatusBadRequest, "The NRIC / Passport / Birth Cert provided does not match our hospital records. If you have updated your passport, kindly update it at the Front Desk or contact Customer Service at info@islandhospital.com")
            }
        } else {
            return fiber.NewError(fiber.StatusBadRequest, "The NRIC / Passport / Birth Cert provided does not match our hospital records. If you have updated your passport, kindly update it at the Front Desk or contact Customer Service at info@islandhospital.com")
        }

        vesPatientDOB, _ := goment.New(strings.TrimSpace(vesPatient.DOB), "DD-MMM-YYYY")
        inputPatientDOB, _ := goment.New(strings.TrimSpace(data.UserDOB), "DD/MM/YYYY")
        if !inputPatientDOB.IsSame(vesPatientDOB) {
            return fiber.NewError(fiber.StatusBadRequest, "Incorrect DOB: The Date of Birth provided does not match our hospital records. Please retry")
        }

        switch data.SignInType {
        case 1:
            if data.UserMobileNo != "" {
                isSameMobileNo := strings.EqualFold(strings.TrimSpace(data.UserMobileNo), strings.TrimSpace(appPatient.Username.String))
                if !isSameMobileNo {
                    isExistsByMobileNo, err := cr.applicationUserService.ExistsByMobileNo(data.UserMobileNo)
                    if err != nil {
                        return err
                    }
                    if isExistsByMobileNo {
                        return fiber.NewError(fiber.StatusBadRequest, "Sorry, an account with the provided mobile number already exists. Please sign in or use a different mobile number to register. Contact our Customer Service for assistance at info@islandhospital.com")
                    }
                }
            } else {
                return fiber.NewError(fiber.StatusBadRequest, "Mobile Number is required")
            }
        case 2:
            if data.UserEmail != "" {
                isSameEmail := strings.EqualFold(strings.TrimSpace(data.UserEmail), strings.TrimSpace(appPatient.Username.String))
                if !isSameEmail {
                    isExistsByEmail, err := cr.applicationUserService.ExistsByEmail(data.UserEmail)
                    if err != nil {
                        return err
                    }
                    if isExistsByEmail {
                        return fiber.NewError(fiber.StatusBadRequest, "Sorry, an account with the provided email already exists. Please sign in or use a different email to register. Contact our Customer Service for assistance at info@islandhospital.com")
                    }
                }
            } else {
                return fiber.NewError(fiber.StatusBadRequest, "Email Address is required")
            }
        default:
            return fiber.NewError(fiber.StatusBadRequest, "Invalid Sign In Method")
        }

        username := data.UserEmail
        pw := data.UserPassword
        if data.SignInType == 1 {
            pw = ""
            username = data.UserMobileNo
        }

        h := vesPatient.HomeAddress
        fullAddress := fmt.Sprintf("%s, %s, %s, %s, %s, %s", h.Address1, h.Address2, h.Address3, h.PostalCode, h.CityState, h.Country)
        fullAddress = strings.TrimSpace(fullAddress)
        o := &model.ApplicationUser{
            Address:         utils.NewNullString(fullAddress),
            Address1:        utils.NewNullString(h.Address1),
            Address2:        utils.NewNullString(h.Address2),
            Address3:        utils.NewNullString(h.Address3),
            CityState:       utils.NewNullString(h.CityState),
            Postcode:        utils.NewNullString(h.PostalCode),
            Country:         utils.NewNullString(h.Country),
            Nationality:     utils.NewNullString(utils.ToTitleCase(vesPatient.Nationality.Description)),
            Race:            utils.NewNullString("-"),
            Sex:             utils.NewNullString(vesPatient.Sex.Description),
            Title:           utils.NewNullString(vesPatient.Name.Title),
            ContactNumber:   utils.NewNullString(vesPatient.ContactNumber.Home),
            Dob:             utils.NewNullString(data.UserDOB),
            Email:           utils.NewNullString(vesPatient.ContactNumber.Email),
            MasterPrn:       utils.NewNullString(vesPatient.Prn),
            FirstName:       utils.NewNullString(vesPatient.Name.FirstName),
            MiddleName:      utils.NewNullString(vesPatient.Name.MiddleName),
            LastName:        utils.NewNullString(vesPatient.Name.LastName),
            FullName:        utils.NewNullString(data.UserFullName),
            Password:        utils.NewNullString(pw),
            Resident:        utils.NewNullString(vesPatient.Resident),
            Role:            utils.NewNullString(constants.ROLE_USER),
            Username:        utils.NewNullString(username),
            FirstTimeLogin:  true,
            FirstTimeLoginV: utils.NewInt32(1),
            PlayerId:        utils.NewNullString(data.PlayerId),
            SignInType:      utils.NewInt32(int32(data.SignInType)), // 1 = Mobile No, 2 = Email Address
            DocNoSignup:     utils.NewNullString(data.UserPersonNumber),
            FullnameSignup:  utils.NewNullString(data.UserFullName),
        }
        middleware.TrimStructFieldsRecursive(o)
        err = cr.applicationUserService.UpdateInactiveSignup(o)
        if err != nil {
            return err
        }

        err = cr.applicationUserFamilyService.SignupSync(appPatient.MasterPrn.String, appPatient.UserId.Int64)
        if err != nil {
            return err
        }

        switch data.SignInType {
        case 1:
            return c.JSON(fiber.Map{
                "successMessage": "Sign up successful",
            })
        case 2:
            go func() {
                cr.mailService.SendSignUp(o, "")
            }()
            return c.JSON(fiber.Map{
                "successMessage": "Thanks for signing up! We have sent you an account activation email, please check your email and follow the steps given.",
            })
        }
    }
    return fiber.NewError(fiber.StatusBadRequest)
}

func (cr *AdminController) signUpPrn(c fiber.Ctx, data *dto.NewSignupUserDto) error {
    isExistsByPrn, err := cr.applicationUserService.ExistsByPRN(data.UserPrn)
    if err != nil {
        return err
    }
    if isExistsByPrn {
        return fiber.NewError(fiber.StatusBadRequest, "Sorry, an account with the provided PRN already exists. Please sign in to your existing account or contact our Customer Service for assistance at info@islandhospital.com")
    }
    vesPatient, ex, err := cr.vesaliusService.VesaliusGetPatientDataByPrn(data.UserPrn)
    if vesPatient == nil {
        return fiber.NewError(fiber.StatusBadRequest, "Incorrect PRN: The Patient PRN Number provided does not exist. Please retry")
    }

    if ex != nil {
        if ex.Code == "99" {
            return fiber.NewError(fiber.StatusBadRequest, "Duplicate Patient Profile found. Please contact customer service for assistance.")
        } else {
            return fiber.NewError(fiber.StatusBadRequest, ex.ToString())
        }
    }

    if err != nil {
        var e *fiber.Error
        if errors.As(err, &e) {
            if e.Code == fiber.StatusNoContent {
                return fiber.NewError(fiber.StatusBadRequest, "Information provided does not match with hospital patient profile. Please retry.")
            }
        }
        return err
    }

    isPatientWithPrn := vesPatient.Prn == data.UserPrn
    if !isPatientWithPrn {
        return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Your patient profile (PRN: %s) is inactive. Please reach out to our customer service at +604-238 3388 for further action", data.UserPrn))
    }

    if len(vesPatient.Documents) > 0 && vesPatient.DOB != "" {
        checkDocument := true
        checkDob := true
        b := false

        for i := range vesPatient.Documents {
            doc := vesPatient.Documents[i]
            if strings.EqualFold(doc.Code, config.GetPatientDocumentCode()) {
                if strings.TrimSpace(doc.Value) == data.UserPersonNumber {
                    b = true
                }
            }
        }
        if !b {
            checkDocument = false
        }

        vesPatientDOB, _ := goment.New(strings.TrimSpace(vesPatient.DOB), "DD-MMM-YYYY")
        inputPatientDOB, _ := goment.New(strings.TrimSpace(data.UserDOB), "DD/MM/YYYY")
        if !inputPatientDOB.IsSame(vesPatientDOB) {
            checkDob = false
        }

        if !checkDocument && !checkDob {
            return fiber.NewError(fiber.StatusBadRequest, "Please verify your PRN, NRIC / Passport / Birth Cert and Date of Birth as they do not match our hospital records. For assistance, contact our Customer Service at the Front Desk or at info@islandhospital.com")
        }
    }

    patientDocIDValue := ""
    if len(vesPatient.Documents) > 0 {
        b := false
        for i := range vesPatient.Documents {
            doc := vesPatient.Documents[i]
            if strings.EqualFold(doc.Code, config.GetPatientDocumentCode()) {
                patientDocIDValue = strings.TrimSpace(doc.Value)
                if config.GetWSVesaliusConfig().NricWithDash == "N" {
                    data.UserPersonNumber = strings.ReplaceAll(data.UserPersonNumber, "-", "")
                }
                if patientDocIDValue == data.UserPersonNumber {
                    b = true
                    if doc.ExpireDate != "" {
                        patientDocExpiry, _ := goment.New(strings.TrimSpace(doc.ExpireDate), "DD-MMM-YYYY")
                        currentDate, _ := goment.New()
                        if patientDocExpiry.IsBefore(currentDate) {
                            return fiber.NewError(fiber.StatusBadRequest, "The NRIC / Passport / Birth Cert provided is not valid. Please confirm your details at the Front Desk or contact Customer Service at info@islandhospital.com")
                        }
                    }
                }
            }
        }
        if !b {
            return fiber.NewError(fiber.StatusBadRequest, "The NRIC / Passport / Birth Cert provided does not match our hospital records. If you have updated your passport, kindly update it at the Front Desk or contact Customer Service at info@islandhospital.com")
        }
    } else {
        return fiber.NewError(fiber.StatusBadRequest, "The NRIC / Passport / Birth Cert provided does not match our hospital records. If you have updated your passport, kindly update it at the Front Desk or contact Customer Service at info@islandhospital.com")
    }

    vesPatientDOB, _ := goment.New(strings.TrimSpace(vesPatient.DOB), "DD-MMM-YYYY")
    inputPatientDOB, _ := goment.New(strings.TrimSpace(data.UserDOB), "DD/MM/YYYY")
    if !inputPatientDOB.IsSame(vesPatientDOB) {
        return fiber.NewError(fiber.StatusBadRequest, "Incorrect DOB: The Date of Birth provided does not match our hospital records. Please retry")
    }

    switch data.SignInType {
    case 1:
        if data.UserMobileNo != "" {
            isExistsByMobileNo, err := cr.applicationUserService.ExistsByMobileNo(data.UserMobileNo)
            if err != nil {
                return err
            }
            if isExistsByMobileNo {
                return fiber.NewError(fiber.StatusBadRequest, "Sorry, an account with the provided mobile number already exists. Please sign in or use a different mobile number to register. Contact our Customer Service for assistance at info@islandhospital.com")
            }
        } else {
            return fiber.NewError(fiber.StatusBadRequest, "Mobile Number is required")
        }
    case 2:
        if data.UserEmail != "" {
            isExistsByEmail, err := cr.applicationUserService.ExistsByEmail(data.UserEmail)
            if err != nil {
                return err
            }
            if isExistsByEmail {
                return fiber.NewError(fiber.StatusBadRequest, "Sorry, an account with the provided email address already exists. Please sign in or use a different email address to register. Contact our Customer Service for assistance at info@islandhospital.com")
            }
        } else {
            return fiber.NewError(fiber.StatusBadRequest, "Email Address is required")
        }
    default:
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Sign In Method")
    }

    username := data.UserEmail
    pw := data.UserPassword
    if data.SignInType == 1 {
        pw = ""
        username = data.UserMobileNo
    }

    h := vesPatient.HomeAddress
    fullAddress := fmt.Sprintf("%s, %s, %s, %s, %s, %s", h.Address1, h.Address2, h.Address3, h.PostalCode, h.CityState, h.Country)
    fullAddress = strings.TrimSpace(fullAddress)
    o := &model.ApplicationUser{
        Address:         utils.NewNullString(fullAddress),
        Address1:        utils.NewNullString(h.Address1),
        Address2:        utils.NewNullString(h.Address2),
        Address3:        utils.NewNullString(h.Address3),
        CityState:       utils.NewNullString(h.CityState),
        Postcode:        utils.NewNullString(h.PostalCode),
        Country:         utils.NewNullString(h.Country),
        Nationality:     utils.NewNullString(utils.ToTitleCase(vesPatient.Nationality.Description)),
        Race:            utils.NewNullString("-"),
        Sex:             utils.NewNullString(vesPatient.Sex.Description),
        Title:           utils.NewNullString(vesPatient.Name.Title),
        ContactNumber:   utils.NewNullString(vesPatient.ContactNumber.Home),
        Dob:             utils.NewNullString(data.UserDOB),
        Email:           utils.NewNullString(vesPatient.ContactNumber.Email),
        MasterPrn:       utils.NewNullString(vesPatient.Prn),
        FirstName:       utils.NewNullString(vesPatient.Name.FirstName),
        MiddleName:      utils.NewNullString(vesPatient.Name.MiddleName),
        LastName:        utils.NewNullString(vesPatient.Name.LastName),
        FullName:        utils.NewNullString(data.UserFullName),
        Password:        utils.NewNullString(pw),
        Resident:        utils.NewNullString(vesPatient.Resident),
        Role:            utils.NewNullString(constants.ROLE_USER),
        Username:        utils.NewNullString(username),
        FirstTimeLogin:  true,
        FirstTimeLoginV: utils.NewInt32(1),
        PlayerId:        utils.NewNullString(data.PlayerId),
        SignInType:      utils.NewInt32(int32(data.SignInType)), // 1 = Mobile No, 2 = Email Address
        DocNoSignup:     utils.NewNullString(data.UserPersonNumber),
        FullnameSignup:  utils.NewNullString(data.UserFullName),
    }
    userId, err := cr.applicationUserService.SaveNewSignup(int64(data.BranchId), o)
    if err != nil {
        return err
    }

    if userId < 0 {
        return fiber.NewError(fiber.StatusBadRequest, "Patient failed to register")
    }

    appPatient, err := cr.applicationUserService.FindByUserId(userId, nil)
    if err != nil {
        return err
    }

    if appPatient == nil {
        return fiber.NewError(fiber.StatusBadRequest, "Patient failed to register")
    }

    err = cr.applicationUserFamilyService.SignupSync(appPatient.MasterPrn.String, appPatient.UserId.Int64)
    if err != nil {
        return err
    }

    switch data.SignInType {
    case 1:
        return c.JSON(fiber.Map{
            "successMessage": "Sign up successful",
        })
    case 2:
        go func() {
            cr.mailService.SendSignUp(o, "")
        }()
        return c.JSON(fiber.Map{
            "successMessage": "Thanks for signing up! We have sent you an account activation email, please check your email and follow the steps given.",
        })
    }
    return fiber.NewError(fiber.StatusBadRequest)
}
