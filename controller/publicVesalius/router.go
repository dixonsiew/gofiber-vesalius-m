package publicVesalius

import (
    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    publicVesaliusController := NewPublicVesaliusController()
    publicVesaliusController.registerRoutes(router)
}

func (c *PublicVesaliusController) registerRoutes(router fiber.Router) {
    api := router.Group("/public/vesalius")
    api.Get("/patient-data/:branchId/:prn", c.GetPatientData)
    api.Post("/get-next-session-available-slots/:branchId/:prn", c.GetNextSessionAvailableSlots)
    api.Post("/get-next-available-slots/:branchId/:prn", c.GetNextAvailableSlots)
    api.Post("/make-appointment/:branchId/:prn", c.GetMakeAppointment)
    api.Get("/getAllDoctorInformation/:branchId/:webadmin", c.GetAllDoctorInformation)
    api.Get("/getAllHSDoctorInformation/:branchId", c.GetAllHSDoctorInformation)
    api.Post("/getAllDoctorInformation/:branchId", c.SearchAllDoctorInformation)
    api.Get("/getDoctorInformationByMCR/:branchId/:mcr", c.GetDoctorInformationByMCR)
    api.Get("/lookup/specialty", c.GetSpecialtyLookup)
    api.Get("/check-file-size/:size", c.CheckFileSize)
}
