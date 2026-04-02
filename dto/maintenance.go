package dto

type NotificationSettingDto struct {
    NotificationCode   string `json:"notificationCode" validate:"required"`
    NotificationParam1 int    `json:"notificationParam1" validate:"required,numeric"`
    NotificationParam2 int    `json:"notificationParam2" validate:"required,numeric"`
}

type ParamSettingDto struct {
    ParamCode  string `json:"paramCode" validate:"required"`
    ParamValue int    `json:"paramValue" validate:"required,numeric"`
}

type HospitalProfileDto struct {
    ProfileDesc  string `json:"profileDesc" validate:"required"`
    ProfileValue string `json:"profileValue" validate:"required"`
}
type DynamicEmailMasterDto struct {
    EmailFunctionName string `json:"emailFunctionName" validate:"required"`
    EmailModule       string `json:"emailModule" validate:"required"`
    EmailFor          string `json:"emailFor" validate:"required"`
    EmailSubject      string `json:"emailSubject" validate:"required"`
    EmailRecipient    string `json:"emailRecipient"`
    EmailSender       string `json:"emailSender" validate:"required"`
    EmailSenderName   string `json:"emailSenderName" validate:"required"`
    EmailTemplate     string `json:"emailTemplate" validate:"required"`
}
