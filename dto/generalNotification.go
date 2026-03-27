package dto

type GeneralNotificationDto struct {
	NotificationTitle  string `json:"notificationTitle" validate:"required"`
	ShortMessage       string `json:"shortMessage" validate:"required"`
	FullMessage        string `json:"fullMessage" validate:"required"`
	StartDate          string `json:"startDate" validate:"required"`
	EndDate            string `json:"endDate" validate:"required"`
	TargetAgeFrom      int    `json:"targetAgeFrom" validate:"required,numeric"`
	TargetAgeTo        int    `json:"targetAgeTo" validate:"required,numeric"`
	TargetGender       string `json:"targetGender" validate:"required"`
	TargetNationality  string `json:"targetNationality" validate:"required"`
	TargetCity         string `json:"targetCity" validate:"required"`
	TargetState        string `json:"targetState" validate:"required"`
}
