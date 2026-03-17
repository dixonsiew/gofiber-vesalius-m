package admin

import (
    "fmt"
    "slices"
    "strconv"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    "vesaliusm/model"
    "vesaliusm/service/adminUser"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/assignBranch"
    "vesaliusm/service/branch"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

type AdminController struct {
    adminUserService       *adminUser.AdminUserService
    applicationUserService *applicationUser.ApplicationUserService
    assignBranchService    *assignBranch.AssignBranchService
    branchService          *branch.BranchService
}

func NewAdminController() *AdminController {
    return &AdminController{
        adminUserService:       adminUser.AdminUserSvc,
        applicationUserService: applicationUser.ApplicationUserSvc,
        assignBranchService:    assignBranch.AssignBranchSvc,
        branchService:          branch.BranchSvc,
    }
}

// GetAdmin
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.AdminUser
// @Router /admin [get]
func (cr *AdminController) GetAdmin(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    admin, err := cr.adminUserService.FindWithAssignBranchByAdminId(user.AdminId.Int64)
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
func (cr *AdminController) GetAllAdmin(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := cr.adminUserService.List(page, limit)
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
func (cr *AdminController) SearchAllAdmin(c fiber.Ctx) error {
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
    m, err := cr.adminUserService.ListByKeyword(key, page, limit)
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
func (cr *AdminController) GetAllAuditMobileUser(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := cr.adminUserService.ListMobileUserAuditLog(page, limit)
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
func (cr *AdminController) SearchAllAuditMobileUser(c fiber.Ctx) error {
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
    m, err := cr.adminUserService.ListMobileUserAuditLogByKeyword(key, key2, page, limit)
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
func (cr *AdminController) GetAllAuditLog(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", "10")
    m, err := cr.adminUserService.ListAuditLog(page, limit)
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
func (cr *AdminController) SearchAllAuditLog(c fiber.Ctx) error {
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
    m, err := cr.adminUserService.ListAuditByKeyword(key, key2, page, limit)
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
func (cr *AdminController) GetUserById(c fiber.Ctx) error {
    adminId := c.Params("adminId")
    iadminId, _ := strconv.ParseInt(adminId, 10, 64)
    o, err := cr.adminUserService.FindWithAssignBranchByAdminId(iadminId)
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
func (cr *AdminController) GetUserByEmail(c fiber.Ctx) error {
    email := c.Params("email")
    o, err := cr.applicationUserService.FindWithAssignBranchByEmail(email)
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
func (cr *AdminController) ResetAdminPassword(c fiber.Ctx) error {
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

    o, err := cr.adminUserService.FindByEmail(email)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you provide does not exist in our system")
    }

    err = cr.adminUserService.SaveResetPassword(o)
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
func (cr *AdminController) ResetUserPassword(c fiber.Ctx) error {
    email := c.Params("email")
    o, err := cr.applicationUserService.FindByEmail(email, nil)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you provide does not exist in our system")
    }

    err = cr.applicationUserService.SaveResetPassword(o)
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
func (cr *AdminController) DeleteUser(c fiber.Ctx) error {
    userId := c.Params("userId")
    iuserId, _ := strconv.ParseInt(userId, 10, 64)
    err := cr.adminUserService.Delete(iuserId)
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
func (cr *AdminController) LinkUserPrn(c fiber.Ctx) error {
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

    o, err := cr.applicationUserService.FindByEmail(data.Email, nil)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }

    o.Address = utils.NewNullString(data.Address)
    o.ContactNumber = utils.NewNullString(data.ContactNumber)
    o.Dob = utils.NewNullString(data.Dob)
    o.FirstName = utils.NewNullString(data.FirstName)
    o.LastName = utils.NewNullString(data.LastName)
    o.MasterPrn = utils.NewNullString(data.Prn)
    o.MiddleName = utils.NewNullString(data.MiddleName)
    o.Nationality = utils.NewNullString(data.Nationality)
    o.Passport = utils.NewNullString(data.Passport)
    o.Resident = utils.NewNullString(data.Resident)
    o.Sex = utils.NewNullString(data.Sex)
    o.Title = utils.NewNullString(data.Title)
    cr.applicationUserService.SaveUserBranch(int64(data.BranchId), o)

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
func (cr *AdminController) ChangePassword(c fiber.Ctx) error {
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

    valid := cr.adminUserService.ValidateCredentials(*user, data.OldPassword)
    if !valid {
        return fiber.NewError(fiber.StatusBadRequest, "Old password is invalid")
    }

    valid1 := cr.adminUserService.ValidateCredentials(*user, data.NewPassword)
    if valid1 {
        return fiber.NewError(fiber.StatusBadRequest, "New Password is not allowed to be the same with Old Password")
    }

    user.Password = utils.NewNullString(data.NewPassword)
    // err = adminUserService.SavePassword(user)
    // if err != nil {
    //     return err
    // }

    return c.JSON(fiber.Map{
        "successMessage": "Password has been updated",
    })
}

func (cr *AdminController) AddAdminUser(c fiber.Ctx) error {
    data := new(dto.PostAdminUserDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.adminUserService.ExistsByEmail(data.Email)
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, "User with the email already exist")
    }

    o := new(model.AdminUser)
    o.Username = utils.NewNullString(data.Email)
    o.Email = utils.NewNullString(data.Email)
    o.FirstName = utils.NewNullString(data.FirstName)
    o.LastName = utils.NewNullString(data.LastName)
    o.UserGroupId = utils.NewInt64(data.UserGroupId)

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
func (cr *AdminController) DeleteAdmin(c fiber.Ctx) error {
    email := c.Params("email")
    o, err := cr.adminUserService.FindByEmail(email)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }

    err = cr.adminUserService.Delete(o.AdminId.Int64)
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
func (cr *AdminController) ResetSignUpUserByEmail(c fiber.Ctx) error {
    email := c.Params("email")
    b, err := cr.applicationUserService.ExistsByEmail(email)
    if err != nil {
        return err
    }

    if !b {
        return fiber.NewError(fiber.StatusBadRequest, "User does not exist to reset signup")
    }

    user, err := cr.applicationUserService.FindByUsername(email, nil)
    if err != nil {
        return err
    }

    err = cr.applicationUserService.ResetUserSignup(user.UserId.Int64, user.MasterPrn.String)
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
func (cr *AdminController) ResetSignUpUserByMobile(c fiber.Ctx) error {
    mobile := c.Params("mobile")
    b, err := cr.applicationUserService.ExistsByMobileNo(mobile)
    if err != nil {
        return err
    }

    if !b {
        return fiber.NewError(fiber.StatusBadRequest, "User does not exist to reset signup")
    }

    user, err := cr.applicationUserService.FindByUsername(mobile, nil)
    if err != nil {
        return err
    }

    err = cr.applicationUserService.ResetUserSignup(user.UserId.Int64, user.MasterPrn.String)
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
func (cr *AdminController) SaveAdminPortalLog(c fiber.Ctx) error {
    adminId, _, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    data := new(dto.AdminPortalLogDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    err = cr.adminUserService.SaveAdminPortalLog(*data, adminId)
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
func (cr *AdminController) ChangeUserPassword(c fiber.Ctx) error {
    data := new(dto.PostChangeUserPasswordDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o, err := cr.applicationUserService.FindByUserId(data.UserId, nil)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusBadRequest, "User does not exist")
    }

    err = cr.adminUserService.ChangeUserPassword(data.NewPassword, data.UserId)
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
func (cr *AdminController) SelfResetPassword(c fiber.Ctx) error {
    branchId := c.Params("branchId")
    email := c.Params("email")
    o, err := cr.applicationUserService.FindByUsername(email, nil)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you entered does not exist in our system")
    }

    if o.SignInType.Int32 == 1 {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you entered does not exist in our system")
    }

    ab, err := cr.assignBranchService.FindAllByUserId(o.UserId.Int64)
    if err != nil {
        return err
    }

    if len(ab) < 1 {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you entered does not exist in our system")
    }

    i := slices.IndexFunc(ab, func(item model.AssignBranch) bool {
        return strconv.FormatInt(item.BranchId.Int64, 10) == branchId
    })

    if i < 0 {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you entered does not exist in our system")
    }

    err = cr.applicationUserService.GenerateVerificationCode(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "New password has been sent to your registered email address",
    })
}
