package auth

import (
	"vesaliusm/dto"
	authService "vesaliusm/service/auth"
	tokenService "vesaliusm/service/token"
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

    user, err := authService.AuthenticateUser(data)
    if err != nil {
        return err
    }

    token, err := tokenService.GenerateAccessToken(*user)
    if err != nil {
        return  err
    }

    refreshToken, err := tokenService.GenerateRefreshToken(*user)
    if err != nil {
        return err
    }

    c.Set("Authorization", token)
    return c.JSON(fiber.Map{
        "type": "bearer",
        "token": token,
        "refresh_token": refreshToken,
        "isFirstTimeLogin": user.FirstTimeLogin,
        "isFirstTimeBiometric": user.FirstTimeBiometric,
        "role": user.Role,
    })
}