package qms

import (
    "vesaliusm/dto"
    "vesaliusm/middleware"
    model "vesaliusm/model/qms"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/applicationUserFamily"
    "vesaliusm/service/qms"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

type QmsController struct {
    applicationUserService       *applicationUser.ApplicationUserService
    applicationUserFamilyService *applicationUserFamily.ApplicationUserFamilyService
    qmsService                   *qms.QmsService
}

func NewQmsController() *QmsController {
    return &QmsController{
        applicationUserService:       applicationUser.ApplicationUserSvc,
        applicationUserFamilyService: applicationUserFamily.ApplicationUserFamilySvc,
        qmsService:                   qms.QmsSvc,
    }
}

// QmsServerWebhook
//
// @Tags QMS
// @Produce json
// @Success 200
// @Router /qms/backend/qms_response [post]
func (cr *QmsController) QmsServerWebhook(c fiber.Ctx) error {
    data := new(dto.QMSServerWebhookDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    err := cr.qmsService.SaveQMS(data)
    if err != nil {
        return err
    }

    return c.SendString("0")
}

// QmsClientWebhook
//
// @Tags QMS
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.QueueResult
// @Router /qms/backend/qms_request [post]
func (cr *QmsController) QmsClientWebhook(c fiber.Ctx) error {
    userId, _, err := middleware.ValidateToken(c)
    if err != nil {
        return err
    }

    patients := make([]model.QmsPatient, 0)
    selfAndFamily, err := cr.applicationUserFamilyService.FindAllByUserId(userId, 0, 50, true, true, nil)
    if err != nil {
        return err
    }

    for i := range selfAndFamily {
        patient := selfAndFamily[i]
        rel := patient.Relationship.String
        if rel == "Self" {
            rel = "self"
        }
        patients = append(patients, model.QmsPatient{
            Prn:          patient.PatientPrn.String,
            Relationship: rel,
        })
    }
    lx, err := cr.qmsService.CallQMS(patients)
    if err != nil {
        return err
    }
    return c.JSON(lx)
}
