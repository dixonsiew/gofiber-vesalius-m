package admin

import (
	"fmt"
	"slices"
	"strconv"
	"vesaliusm/database"
	"vesaliusm/dto"
	"vesaliusm/middleware"
	"vesaliusm/model"
	adminUserService "vesaliusm/service/adminUser"
	applicationuserService "vesaliusm/service/applicationUser"
	assignBranchService "vesaliusm/service/assignBranch"
	"vesaliusm/utils"

	"github.com/gofiber/fiber/v3"
)

var (
    adminUserSvc *adminUserService.AdminUserService = 
        adminUserService.NewAdminUserService(database.GetDb(), database.GetCtx())
    applicationUserSvc *applicationuserService.ApplicationUserService = 
        applicationuserService.NewApplicationUserService(database.GetDb(), database.GetCtx())
    assignBranchSvc *assignBranchService.AssignBranchService = 
        assignBranchService.NewAssignBranchService(database.GetDb(), database.GetCtx())
)

// GetAdmin
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.AdminUser
// @Router /admin [get]
func GetAdmin(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    admin, err := adminUserSvc.FindWithAssignBranchByAdminId(user.AdminID.Int64)
    if err != nil {
        return err
    }

    return c.JSON(admin)
}

// GetAllAdmin
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.AdminUser
// @Router /admin/all [get]
func GetAllAdmin(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := adminUserSvc.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllAdmin
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeywordDto  false  "Search"
// @Success 200 {array} model.AdminUser
// @Router /admin/all [post]
func SearchAllAdmin(c fiber.Ctx) error {
    var data fiber.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    var key string
    if keyword, ok := data["keyword"]; ok {
        key = keyword.(string)
        if key != "" {
            key = "%" + key + "%"
        }
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := adminUserSvc.ListByKeyword(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// GetAllAuditMobileUser
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.MobileUserAuditLog
// @Router /admin/adminportal/mobile-user/log/all [get]
func GetAllAuditMobileUser(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := adminUserSvc.ListMobileUserAuditLog(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllAuditMobileUser
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeyword2Dto  false  "Search"
// @Success 200 {array} model.MobileUserAuditLog
// @Router /admin/adminportal/mobile-user/log/all [post]
func SearchAllAuditMobileUser(c fiber.Ctx) error {
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
    m, err := adminUserSvc.ListMobileUserAuditLogByKeyword(key, key2, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// GetAllAuditLog
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} model.AdminAuditLog
// @Router /admin/adminportal/log/all [get]
func GetAllAuditLog(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := adminUserSvc.ListAuditLog(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllAuditLog
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Param        keyword           body        dto.SearchKeyword2Dto  false  "Search"
// @Success 200 {array} model.AdminAuditLog
// @Router /admin/adminportal/log/all [post]
func SearchAllAuditLog(c fiber.Ctx) error {
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
    m, err := adminUserSvc.ListAuditByKeyword(key, key2, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, fmt.Sprintf("%d", m.Total))
    c.Set(utils.X_TOTAL_PAGE, fmt.Sprintf("%d", m.TotalPages))
    return c.JSON(m.List)
}

// GetUserById
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        adminId           path        string  true  "AdminId"
// @Success 200 {object} model.AdminUser
// @Router /admin/adminId/{adminId} [get]
func GetUserById(c fiber.Ctx) error {
    adminId := c.Params("adminId")
    iadminId, _ := strconv.ParseInt(adminId, 10, 64)
    o, err := adminUserSvc.FindWithAssignBranchByAdminId(iadminId)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }

    return c.JSON(o)
}

// GetUserByEmail
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        email           path        string  true  "Email"
// @Success 200 {object} model.ApplicationUser
// @Router /admin/search-user-email/{email} [get]
func GetUserByEmail(c fiber.Ctx) error {
    email := c.Params("email")
    o, err := applicationUserSvc.FindWithAssignBranchByEmail(email)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }

    return c.JSON(o)
}

// ResetAdminPassword
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        email           path        string  true  "Email"
// @Success 200
// @Router /admin/reset-admin-password/{email} [post]
func ResetAdminPassword(c fiber.Ctx) error {
    email := c.Params("email")
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    if user.Role.String != utils.ROLE_SUPER_ADMIN && user.Role.String != utils.ROLE_ADMIN {
        return middleware.Unauthorized(c)
    }

    o, err := adminUserSvc.FindByEmail(email)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you provide does not exist in our system")
    }

    err = adminUserSvc.SaveResetPassword(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Password Reset successful",
    })
}

// ResetUserPassword
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        email           path        string  true  "Email"
// @Success 200
// @Router /admin/reset-user-password/{email} [post]
func ResetUserPassword(c fiber.Ctx) error {
    email := c.Params("email")
    o, err := applicationUserSvc.FindByEmail(email, nil)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you provide does not exist in our system")
    }

    err = applicationUserSvc.SaveResetPassword(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "New password has been sent to the registered email address",
    })
}

// DeleteUser
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        userId           path        string  true  "UserId"
// @Success 200
// @Router /admin/delete-user/{userId} [post]
func DeleteUser(c fiber.Ctx) error {
    userId := c.Params("userId")
    iuserId, _ := strconv.ParseInt(userId, 10, 64)
    err := adminUserSvc.Delete(iuserId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "User has been deleted",
    })
}

// LinkUserPrn
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param request body dto.PostLinkUserPrnDto true "PostLinkUserPrnDto Request"
// @Success 200
// @Router /admin/link-user-prn [post]
func LinkUserPrn(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    if user.Role.String != utils.ROLE_SUPER_ADMIN && user.Role.String != utils.ROLE_ADMIN {
        return middleware.Unauthorized(c)
    }

    data := new(dto.PostLinkUserPrnDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o, err := applicationUserSvc.FindByEmail(data.Email, nil)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }

    o.Address.String = data.Address
    o.ContactNumber.String = data.ContactNumber
    o.Dob.String = data.Dob
    o.FirstName.String = data.FirstName
    o.LastName.String = data.LastName
    o.MasterPrn.String = data.Prn
    o.MiddleName.String = data.MiddleName
    o.Nationality.String = data.Nationality
    o.Passport.String = data.Passport
    o.Resident.String = data.Resident
    o.Sex.String = data.Sex
    o.Title.String = data.Title
    applicationUserSvc.SaveUserBranch(int64(data.BranchId), o)

    // TODO: Implement link user PRN logic
    // For now, just return success
    return c.JSON(fiber.Map{
        "successMessage": "Hospital has linked successfully",
    })
}

// ChangePassword
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param request body dto.PostChangePasswordDto true "PostLinkUserPrnDto Request"
// @Success 200
// @Router /admin/change-password [post]
func ChangePassword(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
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

    valid := adminUserSvc.ValidateCredentials(*user, data.OldPassword)
    if !valid {
        return fiber.NewError(fiber.StatusBadRequest, "Old password is invalid")
    }

    valid1 := adminUserSvc.ValidateCredentials(*user, data.NewPassword)
    if valid1 {
        return fiber.NewError(fiber.StatusBadRequest, "New Password is not allowed to be the same with Old Password")
    }

    user.Password.String = data.NewPassword
    // err = adminUserSvc.SavePassword(user)
    // if err != nil {
    //     return err
    // }
    
    return c.JSON(fiber.Map{
        "successMessage": "Password has been updated",
    })
}

func AddAdminUser(c fiber.Ctx) error {
    data := new(dto.PostAdminUserDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := adminUserSvc.ExistsByEmail(data.Email)
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, "User with the email already exist")
    }

    o := new(model.AdminUser)
    o.Username.String = data.Email
    o.Email.String = data.Email
    o.FirstName.String = data.FirstName
    o.LastName.String = data.LastName
    o.UserGroupID.Int64 = data.UserGroupId

    return nil
}

// DeleteAdmin
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param        email           path        string  true  "Email"
// @Success 200
// @Router /admin/delete-admin/{email} [post]
func DeleteAdmin(c fiber.Ctx) error {
    email := c.Params("email")
    o, err := adminUserSvc.FindByEmail(email)
    if err != nil {
        return err
    }
    
    if o == nil {
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }

    err = adminUserSvc.Delete(o.AdminID.Int64)
    if err != nil {
        return err
    }
    
    return c.JSON(fiber.Map{
        "successMessage": "Admin has been deleted",
    })
}

// ResetSignUpUserByEmail
//
// @Tags Admin
// @Produce json
// @Param        email           path        string  true  "Email"
// @Success 200
// @Router /admin/reset-signup-email/user/{email} [post]
func ResetSignUpUserByEmail(c fiber.Ctx) error {
    email := c.Params("email")
    b, err := applicationUserSvc.ExistsByEmail(email)
    if err != nil {
        return err
    }

    if !b {
        return fiber.NewError(fiber.StatusBadRequest, "User does not exist to reset signup")
    }

    user, err := applicationUserSvc.FindByUsername(email, nil)
    if err != nil {
        return err
    }

    err = applicationUserSvc.ResetUserSignup(user.UserID.Int64, user.MasterPrn.String)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "User sign up has been reset",
    })
}

// ResetSignUpUserByMobile
//
// @Tags Admin
// @Produce json
// @Param        mobile           path        string  true  "Mobile"
// @Success 200
// @Router /admin/reset-signup-mobile/user/{mobile} [post]
func ResetSignUpUserByMobile(c fiber.Ctx) error {
    mobile := c.Params("mobile")
    b, err := applicationUserSvc.ExistsByMobileNo(mobile)
    if err != nil {
        return err
    }

    if !b {
        return fiber.NewError(fiber.StatusBadRequest, "User does not exist to reset signup")
    }

    user, err := applicationUserSvc.FindByUsername(mobile, nil)
    if err != nil {
        return err
    }

    err = applicationUserSvc.ResetUserSignup(user.UserID.Int64, user.MasterPrn.String)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "User sign up has been reset",
    })
}

// SaveAdminPortalLog
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param request body dto.AdminPortalLogDto true "AdminPortalLogDto Request"
// @Success 200
// @Router /admin/adminportal/save-log [post]
func SaveAdminPortalLog(c fiber.Ctx) error {
    adminId, _, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    data := new(dto.AdminPortalLogDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    err = adminUserSvc.SaveAdminPortalLog(*data, adminId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Audit Log has been inserted",
    })
}

// ChangeUserPassword
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param request body dto.PostChangeUserPasswordDto true "PostChangeUserPasswordDto Request"
// @Success 200
// @Router /admin/change-user-password [post]
func ChangeUserPassword(c fiber.Ctx) error {
    data := new(dto.PostChangeUserPasswordDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o, err := applicationUserSvc.FindByUserId(data.UserId, nil)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusBadRequest, "User does not exist")
    }

    err = adminUserSvc.ChangeUserPassword(data.NewPassword, data.UserId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Change User Password success",
    })
}

// ChangeUserPassword
//
// @Tags Admin
// @Produce json
// @Param        branchId           path        string  true  "BranchId"
// @Param        email              path        string  true  "Email"
// @Success 200
// @Router /admin/self-reset-password/{branchId}/{email} [post]
func SelfResetPassword(c fiber.Ctx) error {
    branchId := c.Params("branchId")
    email := c.Params("email")
    o, err := applicationUserSvc.FindByUsername(email, nil)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you entered does not exist in our system")
    }

    if o.SignInType.Int32 == 1 {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you entered does not exist in our system")
    }

    ab, err := assignBranchSvc.FindAllByUserId(o.UserID.Int64)
    if err != nil {
        return err
    }

    if len(ab) < 1 {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you entered does not exist in our system")
    }

    i := slices.IndexFunc(ab, func(item model.AssignBranch) bool {
        return strconv.FormatInt(item.BranchID.Int64, 10) == branchId
    })

    if i < 0 {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you entered does not exist in our system")
    }

    err = applicationUserSvc.GenerateVerificationCode(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "New password has been sent to your registered email address",
    })
}
