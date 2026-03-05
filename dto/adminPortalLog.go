package dto

type AdminPortalLogDto struct {
    EventModule   string `json:"eventModule" validate:"required"`
    EventFunction string `json:"eventFunction" validate:"required"`
    EventAction   string `json:"eventAction" validate:"required"`
    EventKeyword  string `json:"eventKeyword"`
    EventDesc     string `json:"eventDesc"`
    PatientPrn    string `json:"patientPrn"`
}
