package middleware

import (
    "fmt"
    "strconv"
    "strings"
    "vesaliusm/model"
    applicationuserService "vesaliusm/service/applicationUser"
    adminUserService "vesaliusm/service/adminUser"
    "vesaliusm/utils"

    jwtware "github.com/gofiber/contrib/v3/jwt"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/extractors"
    "github.com/golang-jwt/jwt/v5"
)

func JWTProtected(c fiber.Ctx) error {
    return jwtware.New(jwtware.Config{
        SigningKey: jwtware.SigningKey{Key: []byte(utils.JWT_SECRET)},
        //ContextKey: "jwt",
        Extractor:  extractors.FromAuthHeader("Bearer"),
        ErrorHandler: func(c fiber.Ctx, err error) error {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "statusCode": fiber.StatusUnauthorized,
                "message":    "Unauthorized",
            })
        },
    })(c)
}

func DecodeToken(c fiber.Ctx) (string, int, string, string, error) {
    tokenStr := c.Get("Authorization")
    tokenStr = strings.ReplaceAll(tokenStr, "Bearer ", "")
    token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(utils.JWT_SECRET), nil
    })

    if err != nil {
        utils.LogError(err)
        return "", 0, "", "", err
    }

    claims, ok := token.Claims.(*jwt.MapClaims)
    if !ok {
        return "", 0, "", "", fmt.Errorf("could not parse claims")
    }

    sub := (*claims)["subject"].(string)
    username := (*claims)["username"].(string)
    sessionId := (*claims)["sessionId"].(string)
    types := (*claims)["type"].(string)
    id, _ := strconv.Atoi(sub)
    return username, id, types, sessionId, nil
}

func ValidateToken(c fiber.Ctx) (int, *model.ApplicationUser, error) {
    _, id, _, _, err := DecodeToken(c)
    if err != nil {
        return id, nil, err
    }

    user, err := applicationuserService.FindByUserId(int64(id))
    if err != nil || user == nil {
        return id, user, err
    }

    return id, user, nil
}

func ValidateAdminToken(c fiber.Ctx) (int, *model.AdminUser, error) {
    _, id, _, _, err := DecodeToken(c)
    if err != nil {
        return id, nil, err
    }

    user, err := adminUserService.FindByAdminId(int64(id))
    if err != nil || user == nil {
        return id, user, err
    }

    return id, user, nil
}

func ValidateAppUser(c fiber.Ctx) error {
    _, id, types, sessionId, err := DecodeToken(c)
    if err != nil || types != "1" {
        return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
    }

    user, err := applicationuserService.FindByUserId(int64(id))
    if err != nil || user == nil {
        return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
    }

    if user.SessionID == "" {
        return fiber.NewError(fiber.StatusUnauthorized, "The system has detected your account is no longer valid. Please sign in again.")
    }

    if sessionId != "" {
        userSession, err := applicationuserService.FindByUserIdSessionId(user.UserID, sessionId)
        if err != nil || userSession == nil {
            return fiber.NewError(fiber.StatusUnauthorized, "The system has detected you have signed in using another device. Please sign in again.")
        }
    }

    return c.Next()
}

func ValidateAdminUser(c fiber.Ctx) error {
    _, _, types, _, err := DecodeToken(c)
    if err != nil || types != "0" {
        return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
    }

    return c.Next()
}

func NoContent(c fiber.Ctx) error {
    return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
        "statusCode": fiber.StatusNoContent,
        "message":    "",
    })
}

func Unauthorized(c fiber.Ctx) error {
    mx := fiber.Map{
        "statusCode": fiber.StatusUnauthorized,
        "message":    "User Not Found",
    }
    return c.Status(fiber.StatusUnauthorized).JSON(mx)
}
