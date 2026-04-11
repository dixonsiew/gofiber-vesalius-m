package user

import (
    "fmt"
    "strconv"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    mm "vesaliusm/model/mail"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/applicationUserNotification"
    "vesaliusm/service/assignBranch"
    "vesaliusm/service/mail"
    "vesaliusm/service/token"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

type UserController struct {
    applicationUserService             *applicationUser.ApplicationUserService
    applicationUserNotificationService *applicationUserNotification.ApplicationUserNotificationService
    assignBranchService                *assignBranch.AssignBranchService
    tokenService                       *token.TokenService
    mailService                        *mail.MailService
}

func NewUserController() *UserController {
    return &UserController{
        applicationUserService:             applicationUser.ApplicationUserSvc,
        applicationUserNotificationService: applicationUserNotification.ApplicationUserNotificationSvc,
        assignBranchService:                assignBranch.AssignBranchSvc,
        tokenService:                       token.TokenSvc,
        mailService:                        mail.MailSvc,
    }
}

// GetAllUsers
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
// @Success 200 {array} model.ApplicationUser
// @Router /user/all [get]
func (cr *UserController) GetAllUsers(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.applicationUserService.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetOldAppUnseenNotificationCount
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {integer} integer
// @Router /user/notification/unseen/count [get]
func (cr *UserController) GetOldAppUnseenNotificationCount(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    n, err := cr.applicationUserNotificationService.CountUnseenByUserId(userId)
    if err != nil {
        return err
    }

    byteData := []byte(strconv.Itoa(n))
    return c.Send(byteData)
}

// GetOldAppNotificationLists
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
// @Success 200 {array} model.OnesignalNotification
// @Router /user/notification/list [get]
func (cr *UserController) GetOldAppNotificationLists(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.applicationUserNotificationService.ListByUserId(userId, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllActiveUsers
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        _page              query      string  false  "_page"  default:"1"
// @Param        _limit             query      string  false  "_limit" default:"10"
// @Success 200 {array} model.ApplicationUser
// @Router /user/all/active [get]
func (cr *UserController) GetAllActiveUsers(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.applicationUserService.ListActive(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SetActive
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        userId              path      string  true  "UserId"
// @Success 200
// @Router /user/active-user/{userId} [post]
func (cr *UserController) SetActive(c fiber.Ctx) error {
    userId := c.Params("userId")
    iuserId, _ := strconv.ParseInt(userId, 10, 64)
    err := cr.applicationUserService.SetActive(iuserId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": fmt.Sprintf("User %s has been set to Active", userId),
    })
}

// SetInactive
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        userId              path      string  true  "UserId"
// @Success 200
// @Router /user/inactive-user/{userId} [post]
func (cr *UserController) SetInactive(c fiber.Ctx) error {
    userId := c.Params("userId")
    iuserId, _ := strconv.ParseInt(userId, 10, 64)
    _, admin, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    if admin == nil {
        return middleware.Unauthorized(c)
    }

    user, err := cr.applicationUserService.FindByUserId(iuserId, nil)
    if err != nil {
        return err
    }

    if user == nil {
        return fiber.NewError(fiber.StatusBadRequest, "User Not Found")
    }

    err = cr.applicationUserService.DeleteUserAccount(user, admin)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": fmt.Sprintf("User %s has been set to Inactive", userId),
    })
}

// SearchAllUsers
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.ApplicationUser
// @Router /user/all [post]
func (cr *UserController) SearchAllUsers(c fiber.Ctx) error {
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
    m, err := cr.applicationUserService.ListByKeyword(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetUserById
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        userId              path      string  true  "UserId"
// @Success 200 {object} model.ApplicationUser
// @Router /user/userId/{userId} [get]
func (cr *UserController) GetUserById(c fiber.Ctx) error {
    userId := c.Params("userId")
    iuserId, _ := strconv.ParseInt(userId, 10, 64)
    user, err := cr.applicationUserService.FindWithAssignBranchByUserId(iuserId)
    if err != nil {
        return err
    }

    return c.JSON(user)
}

// GetUser
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.ApplicationUser
// @Router /user [get]
func (cr *UserController) GetUser(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    user, err := cr.applicationUserService.FindWithAssignBranchByUserId(userId)
    if err != nil {
        return err
    }

    return c.JSON(user)
}

// GetUserBranches
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.AssignBranch
// @Router /user/branches [get]
func (cr *UserController) GetUserBranches(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    lx, err := cr.assignBranchService.FindAllByUserId(userId)
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// ChangePassword
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        request          body      dto.PostChangePasswordDto  true  "PostChangePasswordDto"
// @Success 200
// @Router /user/change-password [post]
func (cr *UserController) ChangePassword(c fiber.Ctx) error {
    _, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    data := new(dto.PostChangePasswordDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    valid := cr.applicationUserService.ValidateCredentials(user, data.OldPassword)
    if !valid {
        return fiber.NewError(fiber.StatusBadRequest, "Old password is invalid")
    }

    valid1 := cr.applicationUserService.ValidateCredentials(user, data.NewPassword)
    if valid1 {
        return fiber.NewError(fiber.StatusBadRequest, "New Password is not allowed to be the same with Old Password")
    }

    user.Password = utils.NewNullString(data.NewPassword)
    err = cr.applicationUserService.SavePassword(user)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Password has been updated",
    })
}

// SetNewPasswordUser
//
// @Tags User
// @Produce json
// @Param        request          body      dto.PostUserSetNewPasswordDto  true  "PostUserSetNewPasswordDto"
// @Success 200
// @Router /user/set-new-password [post]
func (cr *UserController) SetNewPasswordUser(c fiber.Ctx) error {
    data := new(dto.PostUserSetNewPasswordDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    user, err := cr.applicationUserService.FindByUsername(data.Email, nil)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    if user.VerificationCode.String != data.VerificationCode {
        return fiber.NewError(fiber.StatusBadRequest, "Verification failed")
    }

    user.Password = utils.NewNullString(data.Password)
    err = cr.applicationUserService.SavePassword(user)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Password has been updated",
    })
}

// VerifyUser
//
// @Tags User
// @Produce json
// @Param        request          body      dto.UserEmailVerificationDto  true  "UserEmailVerificationDto"
// @Success 200
// @Router /user/verify [post]
func (cr *UserController) VerifyUser(c fiber.Ctx) error {
    data := new(dto.UserEmailVerificationDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    user, err := cr.applicationUserService.FindByEmail(data.Email, nil)
    if err != nil {
        return err
    }

    if user == nil {
        return fiber.NewError(fiber.StatusBadRequest, "Verification failed")
    }

    if user.VerificationCode.String != data.VerificationCode {
        return fiber.NewError(fiber.StatusBadRequest, "Verification failed")
    }

    verified, err := cr.applicationUserService.VerifyUser(user)
    if err != nil {
        return err
    }

    if verified {
        emailPrm := mm.MailSignUpSuccess{
            PatientName: user.FirstName.String,
            Username: user.Username.String,
            Email: "",
        }
        if user.Email.Valid {
            emailPrm.Email = user.Email.String
            go func() {
                cr.mailService.SendSignUpSuccess(emailPrm)
            }()
        }
    }

    return c.JSON(fiber.Map{
        "successMessage": "Succuessfully verified",
    })
}

// VerifyUserEmail
//
// @Tags User
// @Produce json
// @Param        request          body      dto.UserEmailVerificationDto  true  "UserEmailVerificationDto"
// @Success 200
// @Router /user/verify-email [post]
func (cr *UserController) VerifyUserEmail(c fiber.Ctx) error {
    data := new(dto.UserEmailVerificationDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    user, err := cr.applicationUserService.FindByUsername(data.Email, nil)
    if err != nil {
        return err
    }

    if user == nil {
        return fiber.NewError(fiber.StatusBadRequest, "Verification failed")
    }

    if user.VerificationCode.String != data.VerificationCode {
        return fiber.NewError(fiber.StatusBadRequest, "Verification failed")
    }

    verified, err := cr.applicationUserService.VerifyUser(user)
    if err != nil {
        return err
    }

    if verified && user.FirstTimeLogin {
        emailPrm := mm.MailSignUpSuccess{
            PatientName: user.FirstName.String,
            Username: user.Username.String,
            Email: "",
        }
        if user.Username.Valid {
            emailPrm.Email = user.Username.String
            go func() {
                cr.mailService.SendSignUpSuccess(emailPrm)
            }()
        }
    }

    return c.JSON(fiber.Map{
        "successMessage": "Succuessfully verified",
    })
}

// VerifyUserSmsTac
//
// @Tags User
// @Produce json
// @Param        request          body      dto.UserMobileVerificationDto  true  "UserMobileVerificationDto"
// @Success 200
// @Router /user/verify-smstac [post]
func (cr *UserController) VerifyUserSmsTac(c fiber.Ctx) error {
    data := new(dto.UserMobileVerificationDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    user, err := cr.applicationUserService.FindByUsername(data.MobileNo, nil)
    if err != nil {
        return err
    }

    if user == nil {
        return fiber.NewError(fiber.StatusBadRequest, "Verification failed")
    }

    if user.VerificationCode.String != data.TAC {
        return fiber.NewError(fiber.StatusBadRequest, "Verification failed")
    }

    verified, err := cr.applicationUserService.VerifyUserSms(user)
    if err != nil {
        return err
    }

    if verified && user.FirstTimeLogin {
        emailPrm := mm.MailSignUpSuccess{
            PatientName: user.FirstName.String,
            Username: user.Username.String,
            Email: "",
        }
        if user.Email.Valid {
            emailPrm.Email = user.Email.String
            go func() {
                cr.mailService.SendSignUpSuccess(emailPrm)
            }()
        }
    }

    token, err := cr.tokenService.GenerateAccessToken(*user)
    if err != nil {
        return err
    }

    refreshToken, err := cr.tokenService.GenerateRefreshToken(*user)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "type":                 "bearer",
        "token":                token,
        "refresh_token":        refreshToken,
        "isFirstTimeLogin":     user.FirstTimeLogin,
        "isFirstTimeBiometric": user.FirstTimeBiometric,
        "role":                 user.Role,
    })
}

// UserDeleteAccount
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /user/delete-account [post]
func (cr *UserController) UserDeleteAccount(c fiber.Ctx) error {
    userId, user, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    err = cr.applicationUserService.DeleteUserAccount(user, nil)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": fmt.Sprintf("User %d has been deleted", userId),
    })
}

// DisableFirstTimeBiometricUser
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /user/disable-firsttime-bio [post]
func (cr *UserController) DisableFirstTimeBiometricUser(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    err = cr.applicationUserService.DisableFirstTimeBiometricUser(userId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "User first time biometric has been disabled",
    })
}

// UpdatePlayerId
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param        playerId              path      string  true  "PlayerId"
// @Success 200
// @Router /user/update-playerid/{playerId} [post]
func (cr *UserController) UpdatePlayerId(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    playerId := c.Params("playerId")
    err = cr.applicationUserService.UpdatePlayerId(playerId, userId, nil)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "PlayerID succuessfully updated",
    })
}

// AddMachineId
//
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param request body dto.PostMachineInfo true "AddMachineId"
// @Success 200
// @Router /user/add-machine-id [post]
func (cr *UserController) AddMachineId(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    data := new(dto.PostMachineInfo)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    err = cr.applicationUserService.UpdateMachineId(data.MachineId, userId, nil)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "MachineID successfully updated",
    })
}
