package dto

type PostCancelAppointmentDto struct {
    AppointmentNumber string `json:"appointmentNumber" validate:"required"`
    Remark            string `json:"remark"`
}
