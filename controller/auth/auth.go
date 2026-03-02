package auth

import (
    "vesaliusm/dto"
    adminUserService "vesaliusm/service/adminUser"
    authService "vesaliusm/service/auth"
    tokenService "vesaliusm/service/token"
    tokenAdminService "vesaliusm/service/tokenAdmin"
    "vesaliusm/utils"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v3"
)

// Login
//
// @Tags Auth
// @Produce json
// @Param request body dto.LoginDto true "Login Request"
// @Success 200
// @Router /login [post]
func Login(c fiber.Ctx) error {
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
        user, err := adminUserService.FindByEmail(data.Username)
        if err != nil {
            return err
        }

        valid := false
        if user != nil {
            valid = adminUserService.ValidateCredentials(*user, data.Password)
        }

        if valid == false {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
        }

        token, err := tokenAdminService.GenerateAccessToken(*user)
        if err != nil {
            return err
        }

        refreshToken, err := tokenAdminService.GenerateRefreshToken(*user)
        if err != nil {
            return err
        }

        c.Set(fiber.HeaderAuthorization, token)
        return c.JSON(fiber.Map{
            "type": "bearer",
            "token": token,
            "refresh_token": refreshToken,
            "isFirstTimeLogin": false,
            "role": user.Role,
        })
    } else {
        user, err := authService.AuthenticateUser(*data)
        if err != nil {
            return err
        }

        token, err := tokenService.GenerateAccessToken(*user)
        if err != nil {
            return err
        }

        refreshToken, err := tokenService.GenerateRefreshToken(*user)
        if err != nil {
            return err
        }

        c.Set(fiber.HeaderAuthorization, token)
        return c.JSON(fiber.Map{
            "type": "bearer",
            "token": token,
            "refresh_token": refreshToken,
            "isFirstTimeLogin": user.FirstTimeLogin,
            "isFirstTimeBiometric": user.FirstTimeBiometric,
            "role": user.Role,
        })
    }
}