package auth

import (
    "database/sql"
    "vesaliusm/dto"
    "vesaliusm/model"
    "vesaliusm/service/applicationUser"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

var AuthSvc *AuthService = NewAuthService()

type AuthService struct {
    applicationUserService *applicationUser.ApplicationUserService
}

func NewAuthService() *AuthService {
    return &AuthService{
        applicationUserService: applicationUser.ApplicationUserSvc,
    }
}

func (s *AuthService) AuthenticateUserV2(data dto.NewLoginDto) (*model.ApplicationUser, error) {
    valid := false
    user, err := s.applicationUserService.FindByUsername(data.Username, nil)
    if err != nil {
        if err == sql.ErrNoRows {
            return user, fiber.NewError(fiber.StatusBadRequest, "The email address / password provided does not match with our system. Please retry.")
        }
        
        return user, err
    }

    if user != nil {
        if data.FromBiometric == 1 {
            valid = s.applicationUserService.ValidateCredentials2(user, data.Password)
        } else {
            valid = s.applicationUserService.ValidateCredentials(user, data.Password)
        }
    }

    if valid == false {
        return user, fiber.NewError(fiber.StatusBadRequest, "The email address / password provided does not match with our system. Please retry.")
    }

    if data.PlayerId != "" {
        err := s.applicationUserService.UpdatePlayerId(data.PlayerId, user.UserId.Int64, nil)
        if err != nil {
            return user, err
        }

        err = s.applicationUserService.InsertDownloadAppV2(data.MachineId, data.PlayerId, nil)
        if err != nil {
            return user, err
        }
    }

    if data.MachineId != "" {
        err := s.applicationUserService.UpdateMachineId(data.MachineId, user.UserId.Int64, nil)
        if err != nil {
            return user, err
        }
    }

    sessionId, err := s.applicationUserService.SaveSessionId(user.UserId.Int64, nil)
    if sessionId == "" {
        return nil, err
    }

    user.SessionId = utils.NewNullString(sessionId)
    return user, nil
}

func (s *AuthService) AuthenticateUser(data dto.LoginDto) (*model.ApplicationUser, error) {
    valid := false
    user, err := s.applicationUserService.FindByUsername(data.Username, nil)
    if err != nil {
        if err == sql.ErrNoRows {
            return user, fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
        }
        
        return user, err
    }

    if user != nil {
        if data.FromBiometric == 1 {
            valid = s.applicationUserService.ValidateCredentials2(user, data.Password)
        } else {
            valid = s.applicationUserService.ValidateCredentials(user, data.Password)
        }
    }

    if valid == false {
        return user, fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
    }

    if data.PlayerId != "" {
        err := s.applicationUserService.UpdatePlayerId(data.PlayerId, user.UserId.Int64, nil)
        if err != nil {
            return user, err
        }

        err = s.applicationUserService.InsertDownloadAppV2(data.MachineId, data.PlayerId, nil)
        if err != nil {
            return user, err
        }
    }

    if data.MachineId != "" {
        err := s.applicationUserService.UpdateMachineId(data.MachineId, user.UserId.Int64, nil)
        if err != nil {
            return user, err
        }
    }

    if user.InactiveFlag.String == "Y" {
        return user, fiber.NewError(fiber.StatusBadRequest, "Your account is not activated")
    }

    sessionId, err := s.applicationUserService.SaveSessionId(user.UserId.Int64, nil)
    if sessionId == "" {
        return nil, err
    }

    user.SessionId = utils.NewNullString(sessionId)
    return user, nil
}