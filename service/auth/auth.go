package auth

import (
	"vesaliusm/dto"
	"vesaliusm/model"
	userService "vesaliusm/service/application_user"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateUser(data dto.LoginDto) (*model.ApplicationUser, error) {
    valid := false
    user, err := userService.FindByUsername(data.Username)
    if err != nil {
        return user, err
    }

    if user != nil {
        if data.FromBiometric == 1 {
            valid = userService.ValidateCredentials2(*user, data.Password)
        } else {
            valid = userService.ValidateCredentials(*user, data.Password)
        }
    }

    if valid == false {
        return user, fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
    }

    if data.PlayerId != "" {
        err := userService.UpdatePlayerId(data.PlayerId, user.UserID)
        if err != nil {
            return user, err
        }
    }

    if data.MachineId != "" {
        err := userService.UpdateMachineId(data.MachineId, user.UserID)
        if err != nil {
            return user, err
        }
    }

    if user.InactiveFlag == "Y" {
        return user, fiber.NewError(fiber.StatusBadRequest, "Your account is not activated")
    }

    sessionId, err := userService.SaveSessionId(user.UserID)
    if sessionId == "" {
        return nil, err
    }

    return user, nil
}