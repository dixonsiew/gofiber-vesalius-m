package auth

import (
    "vesaliusm/dto"
    "vesaliusm/service/adminUser"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/auth"
    "vesaliusm/service/token"
    "vesaliusm/service/tokenAdmin"
    "vesaliusm/utils"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v3"
)

type AuthController struct {
    adminUserService       *adminUser.AdminUserService
    applicationUserService *applicationUser.ApplicationUserService
    authService            *auth.AuthService
    tokenService           *token.TokenService
    tokenAdminService      *tokenAdmin.TokenAdminService
}

func NewAuthController() *AuthController {
    return &AuthController{
        adminUserService:       adminUser.AdminUserSvc,
        applicationUserService: applicationUser.ApplicationUserSvc,
        authService:            auth.AuthSvc,
        tokenService:           token.TokenSvc,
        tokenAdminService:      tokenAdmin.TokenAdminSvc,
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
        user, err := cr.adminUserService.FindByEmail(data.Username)
        if err != nil {
            return err
        }

        valid := false
        if user != nil {
            valid = cr.adminUserService.ValidateCredentials(*user, data.Password)
        }

        if valid == false {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
        }

        token, err := cr.tokenAdminService.GenerateAccessToken(*user)
        if err != nil {
            return err
        }

        refreshToken, err := cr.tokenAdminService.GenerateRefreshToken(*user)
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
        user, err := cr.authService.AuthenticateUser(*data)
        if err != nil {
            return err
        }

        token, err := cr.tokenService.GenerateAccessToken(*user)
        if err != nil {
            return err
        }

        refreshToken, err := cr.tokenService.GenerateRefreshToken(*user)
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
