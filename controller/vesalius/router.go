package vesalius

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    vesaliusController := NewVesaliusController()
    vesaliusController.registerRoutes(router)
}

func (c *VesaliusController) registerRoutes(router fiber.Router) {
    api := router.Group("/vesalius")
    api.Post("/process-resize-image", c.ProcessResizeImage)

    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Get("/doctor/appointment/all", c.GetAllAppointments)
    api.Post("/doctor/appointment/all", c.SearchAllAppointments)
    api.Post("/doctor", c.CreateDoctor)
    api.Get("/doctor/:doctorId", c.GetDoctorInformationByDoctorId)
    api.Put("/doctor/:doctorId", c.UpdateDoctor)
    api.Put("/doctor-image-delete/:doctorId", c.RemoveDoctorImage)
    api.Delete("/doctor/:doctorId", c.Remove)
    api.Get("/get-doctor-appointments/:doctorId/:month/:year/:needAppt", c.GetDoctorAppointments)
    api.Get("/outstanding-bill/:branchId/:prn/:billNumber", c.GetPatientOutstandingBillData)
    api.Get("/patient-data/:branchId/:prn", c.GetPatientData)
}
