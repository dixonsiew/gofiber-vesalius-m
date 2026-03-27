package dto

type PostMakeAppointmentDto struct {
	CaseType         string `json:"caseType" validate:"required"`
	SlotNumber       string `json:"slotNumber" validate:"required"`
	DoctorMcr        string `json:"doctorMcr" validate:"required"`
	ApptSessionType  string `json:"apptSessionType" validate:"required"`
	Remark           string `json:"remark"`
}

type PostCheckAppointmentDto struct {
	ApptDate        string `json:"apptDate" validate:"required"`
	ApptSessionType string `json:"apptSessionType" validate:"required"`
}
