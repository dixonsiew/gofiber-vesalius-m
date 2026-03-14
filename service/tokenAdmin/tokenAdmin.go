package tokenadmin

import (
    "fmt"
    "strconv"
    "time"
    "vesaliusm/model"
    adminUserService "vesaliusm/service/adminUser"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/golang-jwt/jwt/v5"
)

type TokenAdminService struct {
    adminUserSvc *adminUserService.AdminUserService
}

func NewTokenAdminService(adminUserSvc *adminUserService.AdminUserService) *TokenAdminService {
    return &TokenAdminService{adminUserSvc: adminUserSvc}
}

func (s *TokenAdminService) GenerateAccessToken(user model.AdminUser) (string, error) {
    claims := jwt.MapClaims{
        "username": user.Email.String,
        "type":     "0",
        "subject":  fmt.Sprintf("%d", user.AdminID.Int64),
        "exp":      time.Now().Add(time.Hour * 720).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString([]byte(utils.JWT_SECRET))
    if err != nil {
        utils.LogError(err)
        return "", err
    }

    return t, nil
}

func (s *TokenAdminService) GenerateRefreshToken(user model.AdminUser) (string, error) {
    claims := jwt.MapClaims{
        "username": user.Email.String,
        "type":     "0",
        "subject":  fmt.Sprintf("%d", user.AdminID.Int64),
        "exp":      time.Now().Add(time.Hour * 87600).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString([]byte(utils.JWT_SECRET))
    if err != nil {
        utils.LogError(err)
        return "", err
    }

    return t, nil
}

func (s *TokenAdminService) ResolveRefreshToken(encoded string) (*model.AdminUser, error) {
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

func (s *TokenAdminService) CreateAccessTokenFromRefreshToken(refresh string) (fiber.Map, error) {
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

func (s *TokenAdminService) decodeRefreshToken(tokenStr string) (string, int64, error) {
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

func (s *TokenAdminService) getUserFromRefreshTokenPayload(id int64) (*model.AdminUser, error) {
    return s.adminUserSvc.FindByAdminId(id)
}
