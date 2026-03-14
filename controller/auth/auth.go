package auth

import (
    "vesaliusm/dto"
    adminUserService "vesaliusm/service/adminUser"
    applicationuserService "vesaliusm/service/applicationUser"
    authService "vesaliusm/service/auth"
    tokenService "vesaliusm/service/token"
    tokenAdminService "vesaliusm/service/tokenAdmin"
    "vesaliusm/utils"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v3"
)

type AuthController struct {
    adminUserSvc       *adminUserService.AdminUserService
    applicationUserSvc *applicationuserService.ApplicationUserService
    authSvc            *authService.AuthService
    tokenSvc           *tokenService.TokenService
    tokenAdminSvc      *tokenAdminService.TokenAdminService
}

func NewAuthController(
    adminUserSvc *adminUserService.AdminUserService,
    applicationUserSvc *applicationuserService.ApplicationUserService,
    authSvc *authService.AuthService,
    tokenSvc *tokenService.TokenService,
    tokenAdminSvc *tokenAdminService.TokenAdminService) *AuthController {
    return &AuthController{
        adminUserSvc:       adminUserSvc,
        applicationUserSvc: applicationUserSvc,
        authSvc:            authSvc,
        tokenSvc:           tokenSvc,
        tokenAdminSvc:      tokenAdminSvc,
    }
}

// Login
//
// @Tags Auth
// @Produce json
// @Param request body dto.LoginDto true "Login Request"
// @Success 200
// @Router /login [post]
func (cr *AuthController) Login(c fiber.Ctx) error {
    data := new(dto.LoginDto)
    mx := fiber.Map{
        "statusCode": fiber.StatusUnauthorized,
        "message":    "Invalid Credentials",
    }
    if err := c.Bind().Body(data); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            errs := utils.GetValidationErrors(validationErrors)
            if errs != nil {
                return c.Status(fiber.StatusUnauthorized).JSON(mx)
            }
        }

        return c.Status(fiber.StatusUnauthorized).JSON(mx)
    }

    if data.FromAdmin {
        user, err := cr.adminUserSvc.FindByEmail(data.Username)
        if err != nil {
            return err
        }

        valid := false
        if user != nil {
            valid = cr.adminUserSvc.ValidateCredentials(*user, data.Password)
        }

        if valid == false {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
        }

        token, err := cr.tokenAdminSvc.GenerateAccessToken(*user)
        if err != nil {
            return err
        }

        refreshToken, err := cr.tokenAdminSvc.GenerateRefreshToken(*user)
        if err != nil {
            return err
        }

        c.Set(fiber.HeaderAuthorization, token)
        return c.JSON(fiber.Map{
            "type":             "bearer",
            "token":            token,
            "refresh_token":    refreshToken,
            "isFirstTimeLogin": false,
            "role":             user.Role,
        })
    } else {
        user, err := cr.authSvc.AuthenticateUser(*data)
        if err != nil {
            return err
        }

        token, err := cr.tokenSvc.GenerateAccessToken(*user)
        if err != nil {
            return err
        }

        refreshToken, err := cr.tokenSvc.GenerateRefreshToken(*user)
        if err != nil {
            return err
        }

        c.Set(fiber.HeaderAuthorization, token)
        return c.JSON(fiber.Map{
            "type":                 "bearer",
            "token":                token,
            "refresh_token":        refreshToken,
            "isFirstTimeLogin":     user.FirstTimeLogin,
            "isFirstTimeBiometric": user.FirstTimeBiometric,
            "role":                 user.Role,
        })
    }
}
