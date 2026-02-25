package auth

import (
	"vesaliusm/dto"
	"vesaliusm/model"
	applicationuserService "vesaliusm/service/application_user"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateUser(data dto.LoginDto) (*model.ApplicationUser, error) {
    valid := false
    user, err := applicationuserService.FindByUsername(data.Username)
    if err != nil {
        return user, err
    }

    if user != nil {
        if data.FromBiometric == 1 {
            valid = applicationuserService.ValidateCredentials2(*user, data.Password)
        } else {
            valid = applicationuserService.ValidateCredentials(*user, data.Password)
        }
    }

    if valid == false {
        return user, fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
    }

    if data.PlayerId != "" {
        err := applicationuserService.UpdatePlayerId(data.PlayerId, user.UserID)
        if err != nil {
            return user, err
        }
    }

    if data.MachineId != "" {
        err := applicationuserService.UpdateMachineId(data.MachineId, user.UserID)
        if err != nil {
            return user, err
        }
    }

    if user.InactiveFlag == "Y" {
        return user, fiber.NewError(fiber.StatusBadRequest, "Your account is not activated")
    }

    sessionId, err := applicationuserService.SaveSessionId(user.UserID)
    if sessionId == "" {
        return nil, err
    }

    return user, nil
}