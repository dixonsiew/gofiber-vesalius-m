package token

import (
    "fmt"
    "strconv"
    "time"
    "vesaliusm/model"
    "vesaliusm/service/applicationUser"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/golang-jwt/jwt/v5"
)

var TokenSvc *TokenService = NewTokenService()

type TokenService struct {
    applicationUserService *applicationUser.ApplicationUserService
}

func NewTokenService() *TokenService {
    return &TokenService{
        applicationUserService: applicationUser.ApplicationUserSvc,
    }
}

func (s *TokenService) GenerateAccessToken(user model.ApplicationUser) (string, error) {
    claims := jwt.MapClaims{
        "username":  user.Email.String,
        "sessionId": user.SessionID.String,
        "type":      "1",
        "subject":   fmt.Sprintf("%d", user.UserID.Int64),
        "exp":       time.Now().Add(time.Hour * 720).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString([]byte(utils.JWT_SECRET))
    if err != nil {
        utils.LogError(err)
        return "", err
    }

    return t, nil
}

func (s *TokenService) GenerateRefreshToken(user model.ApplicationUser) (string, error) {
    claims := jwt.MapClaims{
        "username":  user.Email.String,
        "sessionId": user.SessionID.String,
        "type":      "1",
        "subject":   fmt.Sprintf("%d", user.UserID.Int64),
        "exp":       time.Now().Add(time.Hour * 87600).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString([]byte(utils.JWT_SECRET))
    if err != nil {
        utils.LogError(err)
        return "", err
    }

    return t, nil
}

func (s *TokenService) ResolveRefreshToken(encoded string) (*model.ApplicationUser, error) {
    _, id, err := s.decodeRefreshToken(encoded)
    if id == 0 || err != nil {
        return nil, fmt.Errorf("refresh token not found")
    }

    user, err := s.getUserFromRefreshTokenPayload(id)
    if err != nil {
        return nil, fmt.Errorf("refresh token malformed")
    }

    return user, nil
}

func (s *TokenService) CreateAccessTokenFromRefreshToken(refresh string) (fiber.Map, error) {
    user, err := s.ResolveRefreshToken(refresh)
    if err != nil {
        return nil, err
    }

    token, err := s.GenerateAccessToken(*user)
    if err != nil {
        return nil, err
    }

    return fiber.Map{
        "user":  &user,
        "token": token,
    }, nil
}

func (s *TokenService) decodeRefreshToken(tokenStr string) (string, int64, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(utils.JWT_SECRET), nil
    })

    if err != nil {
        utils.LogError(err)
        return "", 0, err
    }

    claims, ok := token.Claims.(*jwt.MapClaims)
    if !ok {
        return "", 0, fmt.Errorf("could not parse claims")
    }

    sub := (*claims)["subject"].(string)
    username := (*claims)["username"].(string)
    id, _ := strconv.ParseInt(sub, 10, 64)
    return username, id, nil
}

func (s *TokenService) getUserFromRefreshTokenPayload(id int64) (*model.ApplicationUser, error) {
    return s.applicationUserService.FindByUserId(id, nil)
}
