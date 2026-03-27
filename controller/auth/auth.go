package auth

import (
    "database/sql"
    "vesaliusm/dto"
    "vesaliusm/service/adminUser"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/auth"
    "vesaliusm/service/sms"
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
    smsService             *sms.SmsService
    tokenService           *token.TokenService
    tokenAdminService      *tokenAdmin.TokenAdminService
}

func NewAuthController() *AuthController {
    return &AuthController{
        adminUserService:       adminUser.AdminUserSvc,
        applicationUserService: applicationUser.ApplicationUserSvc,
        authService:            auth.AuthSvc,
        smsService:             sms.SmsSvc,
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
            if err == sql.ErrNoRows {
                return fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
            }
            
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

// NewLogin
//
// @Tags Auth
// @Produce json
// @Param request body dto.NewLoginDto true "New Login Request"
// @Success 200
// @Router /login/v2 [post]
func (cr *AuthController) NewLogin(c fiber.Ctx) error {
    data := new(dto.NewLoginDto)
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

    if data.SignInType == 1 {
        user, err := cr.applicationUserService.FindByUsername(data.Username, nil)
        if err != nil {
            if err == sql.ErrNoRows {
                return fiber.NewError(fiber.StatusBadRequest, "The mobile number does not exist in our system. Please retry.")
            }

            return err
        }
        
        if user == nil {
            return fiber.NewError(fiber.StatusBadRequest, "The mobile number does not exist in our system. Please retry.")
        }

        if user.SignInType.Int32 != 1 {
            return fiber.NewError(fiber.StatusBadRequest, "Your account does not support this login method. Please contact customer service for assistance.")
        }

        if user.InactiveFlag.String == "Y" {
            return fiber.NewError(fiber.StatusBadRequest, "The account you are attempting to sign in has been deleted. Please Sign Up again.")
        }

        tac, err := cr.smsService.SendSignIn(data.Username)
        if err != nil {
            return err
        }

        err = cr.applicationUserService.UpdateVerificationCode(tac, user.UserId.Int64)
        if err != nil {
            return err
        }
        
        return c.JSON(fiber.Map{
            "smsTacSent": true,
        })
    }

    if data.SignInType == 2 {
        user, err := cr.authService.AuthenticateUserV2(*data)
        if err != nil {
            return err
        }
        
        if user == nil {
            return fiber.NewError(fiber.StatusBadRequest, "The email address / password provided does not match with our system. Please retry.")
        }

        if user.SignInType.Int32 != 2 {
            return fiber.NewError(fiber.StatusBadRequest, "Your account does not support this login method. Please contact customer service for assistance.")
        }

        if user.InactiveFlag.String == "Y" {
            return fiber.NewError(fiber.StatusBadRequest, "The account you are attempting to sign in has been deleted. Please Sign Up again.")
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

    return fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
}
