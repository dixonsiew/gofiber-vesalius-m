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
    api.Get("/patient-future-appointments/:branchId/:prn/:isHome", c.GetPatientFutureAppointments)
    api.Get("/past-appointments/:branchId/:prn", c.GetPastAppointments)

    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Get("/doctor/appointment/all", c.GetAllAppointments)
    api.Post("/doctor/appointment/all", c.SearchAllAppointments)
    api.Post("/doctor", c.CreateDoctor)
    api.Get("/doctor/:doctorId", c.GetDoctorInformationByDoctorId)
    api.Put("/doctor/:doctorId", c.UpdateDoctor)
    api.Put("/doctor-image-delete/:doctorId", c.RemoveDoctorImage)
    api.Delete("/doctor/:doctorId", c.Remove)
    api.Get("/get-doctor-appointments/:doctorId/:month/:year/:needAppt", c.GetDoctorAppointments)
    api.Get("/get-lab-history/:prn", c.GetLabHistoryForDashboard)
    api.Get("/get-vital-signs-history/:prn", c.GetVitalSignsHistoryForDashboard)
    api.Get("/get-vital-signs-history/:prn/:visitDate/:vitalSignsCode", c.GetVitalSignsHistory)
    api.Get("/patient-allergy/:prn", c.GetPatientAllergy)
    api.Get("/patient-visit/:prn/:pageId", c.GetPatientVisit)
    api.Get("/health-screening-report/:hsrRefNo", c.GetHealthScreeningRpt)
    api.Get("/patient/:prn", c.GetPatientFromReportSchemaByPrn)
    api.Get("/outstanding-bills/:branchId/:prn", c.GetPatientOutstandingBills)
    api.Get("/outstanding-bill/:branchId/:prn/:billNumber", c.GetPatientOutstandingBillData)
    api.Get("/patient-data/:branchId/:prn", c.GetPatientData)
    api.Get("/search-patient-data/:prn/:branchId", c.SearchPatientData)
    api.Get("/future-appointments/:branchId/:prn", c.GetFutureAppointments)
    api.Post("/get-next-available-slots/:branchId/:prn", c.GetNextAvailableSlots)
    api.Post("/check-make-appointment/:branchId/:prn", c.CheckPatientAppointment)
    api.Post("/make-appointment/:branchId/:prn", c.GetMakeAppointment)
    api.Post("/change-appointment/:branchId/:prn", c.GetChangeAppointment)
    api.Post("/cancel-appointment/:branchId/:prn", c.GetCancelAppointment)
    
}
