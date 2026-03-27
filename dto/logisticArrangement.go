package dto

import (
	"vesaliusm/model/logistic"
	"vesaliusm/utils"
)

type LogisticSetupDto struct {
    LogisticSetupValue string `json:"logisticSetupValue" validate:"required"`
}

type LogisticSlotMobileDto struct {
    FlightArrivalDate string `json:"flightArrivalDate" validate:"required"`
    FlightArrivalTime string `json:"flightArrivalTime" validate:"required"`
    WithCompanion     bool   `json:"withCompanion" validate:"required"`
}

type LogisticSlotDto struct {
    DayOfWeek       string `json:"dayOfWeek" validate:"required"`
    PickUpTime      string `json:"pickUpTime" validate:"required"`
    MaxSlots        int    `json:"maxSlots" validate:"required,numeric"`
    DisplaySequence int    `json:"displaySequence" validate:"numeric"`
}

func (o LogisticSlotDto) ToDbModel() logistic.LogisticSlot {
    return logistic.LogisticSlot{
        DayOfWeek:       utils.NewNullString(o.DayOfWeek),
        PickUpTime:      utils.NewNullString(o.PickUpTime),
        MaxSlots:        utils.NewInt32(int32(o.MaxSlots)),
        DisplaySequence: utils.NewInt32(int32(o.DisplaySequence)),
    }
}

type LogisticSlotsDto struct {
    LogisticSlots []LogisticSlotDto `json:"logisticSlots"`
}

type LogisticRequestStatusDto struct {
    Status        string `json:"status" validate:"required"`
    RequestNumber string `json:"requestNumber" validate:"required"`
}

type LogisticRequestDto struct {
    RequesterPrn            string `json:"requesterPrn"`
    RequesterName           string `json:"requesterName" validate:"required"`
    RequesterDob            string `json:"requesterDob"`
    RequesterDocType        string `json:"requesterDocType" validate:"required"`
    RequesterDocNumber      string `json:"requesterDocNumber" validate:"required"`
    RequesterNationality    string `json:"requesterNationality" validate:"required"`
    RequesterEmail          string `json:"requesterEmail" validate:"required"`
    PrimaryDoctor           string `json:"primaryDoctor" validate:"required"`
    VisitWithCompanion      string `json:"visitWithCompanion" validate:"required"`
    CompanionName           string `json:"companionName"`
    CompanionDob            string `json:"companionDob"`
    CompanionDocType        string `json:"companionDocType"`
    CompanionDocNumber      string `json:"companionDocNumber"`
    RelationshipToRequester string `json:"relationshipToRequester"`
    FlightAirlineName       string `json:"flightAirlineName" validate:"required"`
    FlightNumber            string `json:"flightNumber" validate:"required"`
    FlightArrivalDate       string `json:"flightArrivalDate" validate:"required"`
    FlightArrivalTime       string `json:"flightArrivalTime" validate:"required"`
    RequestedPickupDate     string `json:"requestedPickupDate" validate:"required"`
    RequestedPickupTime     string `json:"requestedPickupTime" validate:"required"`
    DateCreate              string `json:"dateCreate"`
    DateUpdate              string `json:"dateUpdate"`
}

type LogisticRequestExportDto struct {
    LogisticRequests []LogisticRequestDto `json:"logisticRequests"`
}
