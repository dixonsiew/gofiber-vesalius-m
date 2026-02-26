package auth

import (
	"vesaliusm/dto"
	authService "vesaliusm/service/auth"
    adminUserService "vesaliusm/service/admin_user"
	tokenService "vesaliusm/service/token"
    tokenAdminService "vesaliusm/service/token-admin"
	"vesaliusm/utils"

	"github.com/gofiber/fiber/v2"
)

// Login
//
// @Tags Auth
// @Produce json
// @Param request body dto.LoginDto true "Login Request"
// @Success 200
// @Router /login [post]
func Login(c *fiber.Ctx) error {
    data := dto.LoginDto{}
    mx := fiber.Map{
        "statusCode": fiber.StatusUnauthorized,
        "message":    "Invalid Credentials",
    }
    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(mx)
    }

    errs := utils.ValidatePayload(data, c)
    if errs != nil {
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
        user, err := authService.AuthenticateUser(data)
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