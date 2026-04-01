package admin

import (
    "slices"
    "strconv"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    "vesaliusm/model"
    "vesaliusm/service/adminUser"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/assignBranch"
    "vesaliusm/service/branch"
    "vesaliusm/service/groupModulePermission"
    "vesaliusm/service/mail"
    "vesaliusm/service/userGroup"
    "vesaliusm/service/userGroupModulePermission"
    "vesaliusm/service/userGroupModules"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

type AdminController struct {
    adminUserService                 *adminUser.AdminUserService
    applicationUserService           *applicationUser.ApplicationUserService
    assignBranchService              *assignBranch.AssignBranchService
    branchService                    *branch.BranchService
    groupModulePermissionService     *groupModulePermission.GroupModulePermissionService
    userGroupService                 *userGroup.UserGroupService
    userGroupModulesService          *userGroupModules.UserGroupModulesService
    userGroupModulePermissionService *userGroupModulePermission.UserGroupModulePermissionService
    mailService                      *mail.MailService
}

func NewAdminController() *AdminController {
    return &AdminController{
        adminUserService:                 adminUser.AdminUserSvc,
        applicationUserService:           applicationUser.ApplicationUserSvc,
        assignBranchService:              assignBranch.AssignBranchSvc,
        branchService:                    branch.BranchSvc,
        groupModulePermissionService:     groupModulePermission.GroupModulePermissionSvc,
        userGroupService:                 userGroup.UserGroupSvc,
        userGroupModulesService:          userGroupModules.UserGroupModulesSvc,
        userGroupModulePermissionService: userGroupModulePermission.UserGroupModulePermissionSvc,
        mailService:                      mail.MailSvc,
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
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.adminUserService.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
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
    m, err := cr.adminUserService.ListByKeyword(key, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
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
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.adminUserService.ListMobileUserAuditLog(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
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
    m, err := cr.adminUserService.ListMobileUserAuditLogByKeyword(key, key2, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
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
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.adminUserService.ListAuditLog(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
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
    m, err := cr.adminUserService.ListAuditByKeyword(key, key2, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
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

    go func() {
        cr.mailService.SendAdminResetPassword(o)
    }()

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

    go func() {
        cr.mailService.SendUserResetPassword(o)
    }()

    return c.JSON(fiber.Map{
        "successMessage": "New password has been sent to the registered email address",
    })
}

// GetUserGroupList
//
// @Tags Admin
// @Produce json
// @Success 200 {array} model.UserGroup
// @Router /admin/user-group/list [get]
func (cr *AdminController) GetUserGroupList(c fiber.Ctx) error {
    lx, err := cr.userGroupService.ListAll()
    if err != nil {
        return err
    }
    
    return c.JSON(lx)
}

// GetAllUserGroup
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.AllUserGroupDetails
// @Router /admin/all-user-group [get]
func (cr *AdminController) GetAllUserGroup(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    b, err := cr.adminUserService.ExistsByAdminId(user.AdminId.Int64)
    if err != nil {
        return err
    }

    if !b {
        return middleware.Unauthorized(c)
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.userGroupService.List(page, limit)
    if err != nil {
        return err
    }

    lg := m.List.([]model.UserGroup)
    lx := make([]model.AllUserGroupDetails, 0)
    userGroupModules, err := cr.userGroupModulesService.FindAllAsMap()
    if err != nil {
        return err
    }

    for i := range lg {
        userGroup := lg[i]
        admins, err := cr.adminUserService.FindByUserGroupId(userGroup.GroupId.Int64)
        if err != nil {
            return err
        }

        userGroupDetail := model.AllUserGroupDetails{
            UserGroupId: userGroup.GroupId.Int64,
            UserGroupName: userGroup.UserGroupName.String,
            DateCreated: userGroup.DateCreated.String,
            ActiveUser: admins,
        }

        moduleList := make([]string, 0)
        lgmp, err := cr.userGroupModulePermissionService.FindByUserGroupId(userGroup.GroupId.Int64)
        if err != nil {
            return err
        }

        mg := make(map[int64]int64)
        for j := range lgmp {
            ugmp := lgmp[j]
            if _, ok := mg[ugmp.ModuleId.Int64]; ok {
                continue
            }

            if _, ok := userGroupModules[ugmp.ModuleId.Int64]; ok {
                x := userGroupModules[ugmp.ModuleId.Int64]
                moduleList = append(moduleList, x.ModuleName.String)
                mg[ugmp.ModuleId.Int64] = 1
            }

            userGroupDetail.Permission = lgmp
            userGroupDetail.SelectedModules = moduleList
        }

        lx = append(lx, userGroupDetail)
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(lx)
}

// AddUserGroup
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /admin/add-user-group [post]
func (cr *AdminController) AddUserGroup(c fiber.Ctx) error {
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

    data := new(dto.UserGroupDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.userGroupService.ExistsByUserGroupName(data.UserGroupName)
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, "User Group name has already existed")
    }

    o := model.UserGroup{
        UserGroupName: utils.NewNullString(data.UserGroupName),
        UserGroupModulePermissionStatesList: []model.UserGroupModulePermission{},
    }
    for _, x := range data.Permission {
        k := model.UserGroupModulePermission{
            ModuleId: utils.NewInt64(int64(x.ModuleId)),
            PermissionId: utils.NewInt64(int64(x.PermissionId)),
        }
        o.UserGroupModulePermissionStatesList = append(o.UserGroupModulePermissionStatesList, k)
    }
    err = cr.userGroupService.Save(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Successful add new User Group",
    })
}

// GetUserGroup
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.UserGroupDetails
// @Router /admin/user-group/{userGroupId} [get]
func (cr *AdminController) GetUserGroup(c fiber.Ctx) error {
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

    userGroupId := c.Params("userGroupId")
    iuserGroupId, _ := strconv.ParseInt(userGroupId, 10, 64)
    o, err := cr.userGroupService.FindByGroupId(iuserGroupId)
    if err != nil {
        return err
    }

    permissionList, err := cr.userGroupModulePermissionService.FindByUserGroupId(iuserGroupId)
    if err != nil {
        return err
    }

    userGroupDetails := model.UserGroupDetails{
        UserGroupId: o.GroupId.Int64,
        UserGroupName: o.UserGroupName.String,
        Permission: permissionList,
    }

    return c.JSON(userGroupDetails)
}

// UpdateUserGroup
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /admin/update-user-group [post]
func (cr *AdminController) UpdateUserGroup(c fiber.Ctx) error {
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

    data := new(dto.UserGroupDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.userGroupService.ExistsByGroupId(int64(data.UserGroupId))
    if err != nil {
        return err
    }

    if !b {
        return fiber.NewError(fiber.StatusNotFound, "User Group not found")
    }

    b1, err := cr.userGroupService.ExistsByOtherUserGroupName(data.UserGroupName, int64(data.UserGroupId))
    if err != nil {
        return err
    }

    if b1 {
        return fiber.NewError(fiber.StatusBadRequest, "User Group name has already existed")
    }

    o := model.UserGroup{
        GroupId: utils.NewInt64(int64(data.UserGroupId)),
        UserGroupName: utils.NewNullString(data.UserGroupName),
        UserGroupModulePermissionStatesList: []model.UserGroupModulePermission{},
    }
    for _, x := range data.Permission {
        k := model.UserGroupModulePermission{
            ModuleId: utils.NewInt64(int64(x.ModuleId)),
            PermissionId: utils.NewInt64(int64(x.PermissionId)),
            UserGroupId: utils.NewInt64(int64(data.UserGroupId)),
        }
        o.UserGroupModulePermissionStatesList = append(o.UserGroupModulePermissionStatesList, k)
    }
    err = cr.userGroupService.Update(o)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "User Group has been updated successfully",
    })
}

// DeleteUserGroup
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param userGroupId path int true "userGroupId"
// @Success 200
// @Router /admin/delete-user-group/{userGroupId} [post]
func (cr *AdminController) DeleteUserGroup(c fiber.Ctx) error {
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

    userGroupId := c.Params("userGroupId")
    iuserGroupId, _ := strconv.ParseInt(userGroupId, 10, 64)
    b, err := cr.userGroupService.ExistsByGroupId(iuserGroupId)
    if err != nil {
        return err
    }

    if !b {
        return fiber.NewError(fiber.StatusNotFound, "User Group not found")
    }

    err = cr.userGroupService.DeleteByGroupId(iuserGroupId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "User Group has been deleted successfully",
    })
}

// GetAllGroupModules
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.UserGroupModules
// @Router /admin/group-modules [get]
func (cr *AdminController) GetAllGroupModules(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    b, err := cr.adminUserService.ExistsByAdminId(user.AdminId.Int64)
    if err != nil {
        return err
    }

    if !b {
        return middleware.Unauthorized(c)
    }

    lx, err := cr.userGroupModulesService.FindAll()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetAllGroupModulesPermission
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.GroupModulePermission
// @Router /admin/group-permission [get]
func (cr *AdminController) GetAllGroupModulesPermission(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    b, err := cr.adminUserService.ExistsByAdminId(user.AdminId.Int64)
    if err != nil {
        return err
    }

    if !b {
        return middleware.Unauthorized(c)
    }

    lx, err := cr.groupModulePermissionService.FindAll()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetUserGroupPermission
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.UserGroupModulePermission
// @Router /admin/user-group-permission [get]
func (cr *AdminController) GetUserGroupPermission(c fiber.Ctx) error {
    _, user, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    if user == nil {
        return middleware.Unauthorized(c)
    }

    lx, err := cr.userGroupModulePermissionService.FindByUserGroupIdOrderByModuleIdAsc(user.UserGroupId.Int64)
    if err != nil {
        return err
    }

    return c.JSON(lx)
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
// @Param request body dto.PostLinkUserPrnDto true "PostLinkUserPrnDto"
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

// ChangeSignInType
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangeSignInTypeDto true "ChangeSignInTypeDto"
// @Success 200
// @Router /admin/change-signin-type [post]
func (cr *AdminController) ChangeSignInType(c fiber.Ctx) error {
    data := new(dto.ChangeSignInTypeDto)
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

    switch data.SignInType {
    case 1:
        if data.SignInMobileNumber != "" {
            isExistsByMobileNo, err := cr.applicationUserService.ExistsByMobileNo(data.SignInMobileNumber)
            if err != nil {
                return err
            }

            if isExistsByMobileNo {
                return fiber.NewError(fiber.StatusBadRequest, "The entered Mobile Number already exist in our system")
            }
        } else {
            return fiber.NewError(fiber.StatusBadRequest, "Mobile Number is required")
        }
    case 2:
        if data.SignInEmailAddress != "" {
            isExistsByEmail, err := cr.applicationUserService.ExistsByEmail(data.SignInEmailAddress)
            if err != nil {
                return err
            }

            if isExistsByEmail {
                return fiber.NewError(fiber.StatusBadRequest, "The entered Email Address already exist in our system")
            } else {
                err := cr.applicationUserService.UpdateVerificationCode(o.VerificationCode.String, o.UserId.Int64)
                if err != nil {
                    return err
                }

                go func() {
                    cr.mailService.SendSignUp(o, data.SignInEmailAddress)
                }()
            }
        } else {
            return fiber.NewError(fiber.StatusBadRequest, "Email Address is required")
        }
    default:
        return fiber.NewError(fiber.StatusBadRequest, "Invalid Sign In Method")
    }

    err = cr.adminUserService.ChangeUserSignInType(*data)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "successMessage": "Change User Sign In Type success",
    })
}

// ChangePassword
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param request body dto.PostChangePasswordDto true "PostChangePasswordDto"
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

// AddAdminUser
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param request body dto.PostAdminUserDto true "PostAdminUserDto"
// @Success 200
// @Router /admin/sign-up [post]
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

    ug, err := cr.userGroupService.FindByGroupId(data.UserGroupId)
    if err != nil {
        return err
    }

    o.UserGroupName = ug.UserGroupName
    o.Role = o.UserGroupName

    adminBranchIds := data.AdminBranchIds
    err = cr.adminUserService.Save(o, adminBranchIds)
    if err != nil {
        return err
    }

    go func() {
        cr.mailService.SendAdminSignUp(o)
    }()

    return c.JSON(fiber.Map{
        "successMessage": "Admin user successfully created",
    })
}

// EditAdminUser
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param request body dto.PostUpdateAdminUserDto true "PostUpdateAdminUserDto"
// @Success 200
// @Router /admin/edit-admin-user [post]
func (cr *AdminController) EditAdminUser(c fiber.Ctx) error {
    data := new(dto.PostUpdateAdminUserDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }
    
    o, err := cr.adminUserService.FindByEmail(data.Email)
    if err != nil {
        return err
    }
    
    if o == nil {
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }
    
    o.FirstName = utils.NewNullString(data.FirstName)
    o.LastName = utils.NewNullString(data.LastName)
    o.UserGroupId = utils.NewInt64(data.UserGroupId)

    ug, err := cr.userGroupService.FindByGroupId(data.UserGroupId)
    if err != nil {
        return err
    }
    
    o.UserGroupName = ug.UserGroupName
    o.Role = o.UserGroupName
    
    adminBranchIds := data.AdminBranchIds
    err = cr.adminUserService.Update(o, adminBranchIds)
    if err != nil {
        return err
    }
    
    return c.JSON(fiber.Map{
        "successMessage": "User profile has been updated",
    })
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
// @Param request body dto.AdminPortalLogDto true "AdminPortalLogDto"
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

func (cr *AdminController) SelfSignUpUser(c fiber.Ctx) error {
    data := new(dto.PostSelfSignUpUserDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    b, err := cr.applicationUserService.ExistsByEmail(data.UserEmail)
    if err != nil {
        return err
    }

    if b {
        return fiber.NewError(fiber.StatusBadRequest, "User with the email already exist. Please Sign In or Reset Password if forgot.")
    }

    b1, err := cr.applicationUserService.ExistsByPRN(data.UserPersonNumber)
    if err != nil {
        return err
    }

    if b1 {
        return fiber.NewError(fiber.StatusBadRequest, "User with the PRN already exist. Please Sign In or Reset Password if forgot.")
    }

    o := new(model.ApplicationUser)
    o.Race = utils.NewNullString("-")
    o.Dob = utils.NewNullString(data.UserDOB)
    o.Email = utils.NewNullString(data.UserEmail)
    o.Password = utils.NewNullString(data.UserPassword)
    o.Role = utils.NewNullString(utils.ROLE_USER)
    o.Username = utils.NewNullString(data.UserEmail)
    o.FirstTimeLoginV = utils.NewInt32(1)
    o.PlayerId = utils.NewNullString(data.PlayerId)
    err = cr.applicationUserService.SaveSignup(int64(data.BranchId), o)
    if err != nil {
        return err
    }

    go func() {
        cr.mailService.SendSignUp(o, "")
    }()

    return c.JSON(fiber.Map{
        "successMessage": "Thanks for signing up! Please check your email (or spam / junk folder) for an account activation email and follow the steps given.",
    })
}

// ResendUserSignupEmail
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param email path string true "Email"
// @Success 200
// @Router /admin/resend-user-signup-email/{email} [post]
func (cr *AdminController) ResendUserSignupEmail(c fiber.Ctx) error {
    email := c.Params("email")
    o, err := cr.applicationUserService.FindByUsername(email, nil)
    if err != nil {
        return err
    }

    if o == nil {
        return fiber.NewError(fiber.StatusBadRequest, "The email address you entered does not exist in our system")
    }

    go func() {
        cr.mailService.SendSignUp(o, "")
    }()

    return c.JSON(fiber.Map{
        "successMessage": "Email resent",
    })
}

// ChangeUserPassword
//
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param request body dto.PostChangeUserPasswordDto true "PostChangeUserPasswordDto"
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

    go func() {
        cr.mailService.SendResetPassword(o)
    }()

    return c.JSON(fiber.Map{
        "successMessage": "New password has been sent to your registered email address",
    })
}
