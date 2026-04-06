package healthCare

import (
    "context"
    "vesaliusm/config"
    "vesaliusm/database"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/jmoiron/sqlx"
)

type HealthCareService struct {
    db                              *sqlx.DB
    ctx                             context.Context
    vitalCodeWeight                 string
    vitalCodeHeight                 string
    vitalCodeBMI                    string
    vitalCodeBP                     string
    vitalCodePulse                  string
    vitalCodeClinicaltemplateWeight string
    vitalCodeClinicaltemplateHeight string
    vitalCodeClinicaltemplateBMI    string
    vitalCodeClinicaltemplatePulse  string
    labInvestigationType            string
    labCodeHDL                      string
    labCodeLDL                      string
    labCodeGlucose                  string
    labCodeHemoglobin               string
}

func NewHealthCareService(db *sqlx.DB, ctx context.Context) *HealthCareService {
    return &HealthCareService{
        db: db,
        ctx: ctx,
        vitalCodeWeight: config.Config("vital.code.weight"),
        vitalCodeHeight: config.Config("vital.code.height"),
        vitalCodeBMI: config.Config("vital.code.bmi"),
        vitalCodeBP: config.Config("vital.code.bp"),
        vitalCodePulse: config.Config("vital.code.pulse"),
        vitalCodeClinicaltemplateWeight: config.Config("vital.code.weight.clinicaltemplate"),
        vitalCodeClinicaltemplateHeight: config.Config("vital.code.height.clinicaltemplate"),
        vitalCodeClinicaltemplateBMI: config.Config("vital.code.bmi.clinicaltemplate"),
        vitalCodeClinicaltemplatePulse: config.Config("vital.code.pulse.clinicaltemplate"),
        labInvestigationType: config.Config("lab.investigation.type"),
        labCodeHDL: config.Config("lab.code.hdl"),
        labCodeLDL: config.Config("lab.code.ldl"),
        labCodeGlucose: config.Config("lab.code.glucose"),
        labCodeHemoglobin: config.Config("lab.code.hemoglobin"),
    }
}

// func (s *NovaDoctorService)

// func (s *NovaDoctorService)

// func (s *NovaDoctorService)
