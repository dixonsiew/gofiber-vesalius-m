package dto

type PostChangeAppointmentDto struct {
    SlotNumber        string `json:"slotNumber" validate:"required"`
    AppointmentNumber string `json:"appointmentNumber" validate:"required"`
    ApptDate          string `json:"apptDate" validate:"required"`
    ApptSessionType   string `json:"apptSessionType" validate:"required"`
    Remark            string `json:"remark"`
}
