package auth

import (
	"vesaliusm/dto"
	"vesaliusm/model"
	applicationuserService "vesaliusm/service/applicationUser"

	"github.com/gofiber/fiber/v3"
)

type AuthService struct {
    applicationUserSvc *applicationuserService.ApplicationUserService
}

func NewAuthService(applicationUserSvc *applicationuserService.ApplicationUserService) *AuthService {
    return &AuthService{applicationUserSvc: applicationUserSvc}
}

func (s *AuthService) AuthenticateUser(data dto.LoginDto) (*model.ApplicationUser, error) {
    valid := false
    user, err := s.applicationUserSvc.FindByUsername(data.Username, nil)
    if err != nil {
        return user, err
    }

    if user != nil {
        if data.FromBiometric == 1 {
            valid = s.applicationUserSvc.ValidateCredentials2(user, data.Password)
        } else {
            valid = s.applicationUserSvc.ValidateCredentials(user, data.Password)
        }
    }

    if valid == false {
        return user, fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
    }

    if data.PlayerId != "" {
        err := s.applicationUserSvc.UpdatePlayerId(data.PlayerId, user.UserID.Int64, nil)
        if err != nil {
            return user, err
        }

        err = s.applicationUserSvc.InsertDownloadAppV2(data.MachineId, data.PlayerId, nil)
        if err != nil {
            return user, err
        }
    }

    if data.MachineId != "" {
        err := s.applicationUserSvc.UpdateMachineId(data.MachineId, user.UserID.Int64, nil)
        if err != nil {
            return user, err
        }
    }

    if user.InactiveFlag.String == "Y" {
        return user, fiber.NewError(fiber.StatusBadRequest, "Your account is not activated")
    }

    sessionId, err := s.applicationUserSvc.SaveSessionId(user.UserID.Int64, nil)
    if sessionId == "" {
        return nil, err
    }

    user.SessionID.String = sessionId
    return user, nil
}