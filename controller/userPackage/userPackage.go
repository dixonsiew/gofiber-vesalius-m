package userPackage

import (
    "vesaliusm/database"
    patientPurchaseDetailsService "vesaliusm/service/patientPurchaseDetails"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

var patientPurchaseDetailsSvc *patientPurchaseDetailsService.PatientPurchaseDetailsService = patientPurchaseDetailsService.NewPatientPurchaseDetailsService(database.GetDb(), database.GetCtx())

func searchAllPurchaseHistory(c fiber.Ctx) error {
    return c.SendString("Hello, World!")
}