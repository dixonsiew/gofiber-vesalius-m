package dto

type WallexWebhookDto struct {
	Resource   string `json:"resource" validate:"required"`
	ResourceId string `json:"resourceId" validate:"required"`
	AccountId  string `json:"accountId" validate:"required"`
	Status     string `json:"status" validate:"required"`
	Remarks    string `json:"remarks"`
	Reason     string `json:"reason"`
}
