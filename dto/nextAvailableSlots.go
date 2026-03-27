package dto

type PostNextAvailableSlotsDto struct {
    SpecialtyCode string `json:"specialtyCode" validate:"required"`
    Mcr           string `json:"mcr" validate:"required"`
    DoctorId      int    `json:"doctorId" validate:"numeric"`
    StartDate     string `json:"startDate" validate:"required"`
    StartTime     string `json:"startTime" validate:"required"`
    CaseType      string `json:"caseType" validate:"required"`
}
