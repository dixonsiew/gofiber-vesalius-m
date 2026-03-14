package auth

import (
    "vesaliusm/database"
    adminUserService "vesaliusm/service/adminUser"
    applicationuserService "vesaliusm/service/applicationUser"
    authService "vesaliusm/service/auth"
    tokenService "vesaliusm/service/token"
    tokenAdminService "vesaliusm/service/tokenAdmin"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    var (
        adminUserSvc *adminUserService.AdminUserService = 
            adminUserService.NewAdminUserService(database.GetDb(), database.GetCtx())
        applicationUserSvc *applicationuserService.ApplicationUserService = 
            applicationuserService.NewApplicationUserService(database.GetDb(), database.GetCtx())
        authSvc *authService.AuthService = 
            authService.NewAuthService(applicationUserSvc)
        tokenSvc *tokenService.TokenService = 
            tokenService.NewTokenService(applicationUserSvc)
        tokenAdminSvc *tokenAdminService.TokenAdminService = 
            tokenAdminService.NewTokenAdminService(adminUserSvc)
    )

    authController := NewAuthController(
        adminUserSvc, 
        applicationUserSvc, 
        authSvc, 
        tokenSvc, 
        tokenAdminSvc)
    authController.registerRoutes(router)
}

func (c *AuthController) registerRoutes(router fiber.Router) {
    router.Post("/login", c.Login)
}
