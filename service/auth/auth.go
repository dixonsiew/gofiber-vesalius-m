package auth

import (
	"vesaliusm/dto"
    "vesaliusm/database"
	"vesaliusm/model"
	applicationuserService "vesaliusm/service/applicationUser"

	"github.com/gofiber/fiber/v3"
)

var applicationUserSvc *applicationuserService.ApplicationUserService = 
    applicationuserService.NewApplicationUserService(database.GetDb(), database.GetCtx())

func AuthenticateUser(data dto.LoginDto) (*model.ApplicationUser, error) {
    valid := false
    user, err := applicationUserSvc.FindByUsername(data.Username, nil)
    if err != nil {
        return user, err
    }

    if user != nil {
        if data.FromBiometric == 1 {
            valid, _ = applicationUserSvc.ValidateCredentials2(user, data.Password)
        } else {
            valid, _ = applicationUserSvc.ValidateCredentials(user, data.Password)
        }
    }

    if valid == false {
        return user, fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
    }

    if data.PlayerId != "" {
        err := applicationUserSvc.UpdatePlayerId(data.PlayerId, user.UserID.Int64, nil)
        if err != nil {
            return user, err
        }

        err = applicationUserSvc.InsertDownloadAppV2(data.MachineId, data.PlayerId, nil)
        if err != nil {
            return user, err
        }
    }

    if data.MachineId != "" {
        err := applicationUserSvc.UpdateMachineId(data.MachineId, user.UserID.Int64, nil)
        if err != nil {
            return user, err
        }
    }

    if user.InactiveFlag.String == "Y" {
        return user, fiber.NewError(fiber.StatusBadRequest, "Your account is not activated")
    }

    sessionId, err := applicationUserSvc.SaveSessionId(user.UserID.Int64, nil)
    if sessionId == "" {
        return nil, err
    }

    user.SessionID.String = sessionId
    return user, nil
}