package admin

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    adminController := NewAdminController()
    adminController.registerRoutes(router)
}

func (c *AdminController) registerRoutes(router fiber.Router) {
	api := router.Group("/admin")
    api.Post("/reset-signup-email/user/:email", c.ResetSignUpUserByEmail)
    api.Post("/self-reset-password/:branchId/:email", c.SelfResetPassword)
    api.Get("/user-group/list", c.GetUserGroupList)
    api.Post("/self-sign-up", c.SelfSignUpUser)
    api.Post("/self-sign-up/v2", c.MobileSignUpUser)
    
    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Get("/", c.GetAdmin)
    api.Get("/all", c.GetAllAdmin)
    api.Post("/all", c.SearchAllAdmin)
    api.Get("/adminportal/mobile-user/log/all", c.GetAllAuditMobileUser)
    api.Post("/adminportal/mobile-user/log/all", c.SearchAllAuditMobileUser)
    api.Get("/adminportal/log/all", c.GetAllAuditLog)
    api.Post("/adminportal/log/all", c.SearchAllAuditLog)
    api.Get("/adminId/:adminId", c.GetUserById)
    api.Post("/reset-admin-password/:email", c.ResetAdminPassword)
    api.Post("/reset-user-password/:email", c.ResetUserPassword)
    api.Get("/all-user-group", c.GetAllUserGroup)
    api.Post("/add-user-group", c.AddUserGroup)
    api.Get("/user-group/:userGroupId", c.GetUserGroup)
    api.Post("/update-user-group", c.UpdateUserGroup)
    api.Post("/delete-user-group/:userGroupId", c.DeleteUserGroup)
    api.Get("/group-modules", c.GetAllGroupModules)
    api.Get("/group-permission", c.GetAllGroupModulesPermission)
    api.Get("/user-group-permission", c.GetUserGroupPermission)
    api.Post("/delete-user/:userId", c.DeleteUser)
    api.Post("/link-user-prn", c.LinkUserPrn)
    api.Post("/unlink-user-prn", c.UnlinkUserPrn)
    api.Post("/set-master-profile", c.SetMasterPrn)
    api.Post("/change-password", c.ChangePassword)
    api.Post("/sign-up", c.AddAdminUser)
    api.Post("/edit-admin-user", c.EditAdminUser)
    api.Post("/delete-admin/:email", c.DeleteAdmin)
    api.Post("/adminportal/save-log", c.SaveAdminPortalLog)

    api.Post("/resend-user-signup-email/:email", c.ResendUserSignupEmail)
    api.Post("/change-signin-type", c.ChangeSignInType)
    api.Post("/change-user-password", c.ChangeUserPassword)
}
