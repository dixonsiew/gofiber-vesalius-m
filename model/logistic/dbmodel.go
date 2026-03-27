package logistic

import (
	"vesaliusm/utils"

	"github.com/guregu/null/v6"
)

type LogisticSetup struct {
	LogisticSetupId    null.Int64  `json:"logistic_setup_id" db:"LOGISTIC_SETUP_ID" swaggertype:"string"`
	LogisticSetupCode  string      `json:"logisticSetupCode" swaggertype:"string"`
	LogisticSetupValue null.String `json:"logisticSetupValue" db:"LOGISTIC_SETUP_VALUE" swaggertype:"string"`
}

type LogisticSlot struct {
	LogisticSlotId  null.Int64  `json:"logistic_slot_id" db:"LOGISTIC_SLOT_ID" swaggertype:"integer"`
	DayOfWeek       null.String `json:"dayOfWeek" db:"DAY_OF_WEEK" swaggertype:"string"`
	PickUpDate      null.String `json:"pickUpDate" db:"REQUESTED_PICKUP_DATE" swaggertype:"string"`
	PickUpTime      null.String `json:"pickUpTime" db:"PICKUP_TIME" swaggertype:"string"`
	MaxSlots        null.Int32  `json:"maxSlots" db:"MAX_SLOTS" swaggertype:"integer"`
	DisplaySequence null.Int32  `json:"displaySequence" db:"DISPLAY_SEQUENCE" swaggertype:"integer"`
	AvailableSlots  null.Int32  `json:"availableSlots" db:"AVAILABLE_SLOTS" swaggertype:"integer"`
}

type LogisticSlots struct {
	LogisticSlots []LogisticSlot `json:"logisticSlots"`
}

type LogisticRequest struct {
	LogisticRequestId        null.Int64  `json:"logistic_request_id" db:"LOGISTIC_REQUEST_ID" swaggertype:"integer"`
	LogisticRequestNumber    null.String `json:"logisticRequestNumber" db:"LOGISTIC_REQUEST_NUMBER" swaggertype:"string"`
	LogisticRequestStatus    null.String `json:"logisticRequestStatus" db:"LOGISTIC_REQUEST_STATUS" swaggertype:"string"`
	DisplayCancelBtn         bool        `json:"displayCancelBtn"`
	RequesterPrn             null.String `json:"requesterPrn" db:"REQUESTER_PRN" swaggertype:"string"`
	RequesterName            null.String `json:"requesterName" db:"REQUESTER_NAME" swaggertype:"string"`
	RequesterDob             null.String `json:"requesterDob" db:"REQUESTER_DOB" swaggertype:"string"`
	RequesterDocType         null.String `json:"requesterDocType" db:"REQUESTER_DOC_TYPE" swaggertype:"string"`
	RequesterDocNumber       null.String `json:"requesterDocNumber" db:"REQUESTER_DOC_NUMBER" swaggertype:"string"`
	RequesterNationality     null.String `json:"requesterNationality" db:"REQUESTER_NATIONALITY" swaggertype:"string"`
	RequesterEmail           null.String `json:"requesterEmail" db:"REQUESTER_EMAIL" swaggertype:"string"`
	PrimaryDoctor            null.String `json:"primaryDoctor" db:"PRIMARY_DOCTOR" swaggertype:"string"`
	PrimaryDoctorName        string      `json:"primaryDoctorName"`
	VisitWithCompanion       null.String `json:"visitWithCompanion" db:"VISIT_WITH_COMPANION" swaggertype:"string"`
	CompanionName            null.String `json:"companionName" db:"COMPANION_NAME" swaggertype:"string"`
	CompanionDob             null.String `json:"companionDob" db:"COMPANION_DOB" swaggertype:"string"`
	CompanionDocType         null.String `json:"companionDocType" db:"COMPANION_DOC_TYPE" swaggertype:"string"`
	CompanionDocNumber       null.String `json:"companionDocNumber" db:"COMPANION_DOC_NUMBER" swaggertype:"string"`
	RelationshipToRequester  null.String `json:"relationshipToRequester" db:"RELATIONSHIP_TO_REQUESTER" swaggertype:"string"`
	FlightAirlineName        null.String `json:"flightAirlineName" db:"FLIGHT_AIRLINE_NAME" swaggertype:"string"`
	FlightNumber             null.String `json:"flightNumber" db:"FLIGHT_NUMBER" swaggertype:"string"`
	FlightArrivalDate        null.String `json:"flightArrivalDate" db:"FLIGHT_ARRIVAL_DATE" swaggertype:"string"`
	FlightArrivalTime        null.String `json:"flightArrivalTime" db:"FLIGHT_ARRIVAL_TIME" swaggertype:"string"`
	RequestedPickupDate      null.String `json:"requestedPickupDate" db:"REQUESTED_PICKUP_DATE" swaggertype:"string"`
	RequestedPickupTime      null.String `json:"requestedPickupTime" db:"REQUESTED_PICKUP_TIME" swaggertype:"string"`
	RequestedPickupDay       null.String `json:"requestedPickupDay" db:"REQUESTED_PICKUP_DAY" swaggertype:"string"`
	DateCreate               null.String `json:"dateCreate" db:"DATE_CREATE" swaggertype:"string"`
	UserUpdate               null.String `json:"userUpdate" db:"USER_UPDATE" swaggertype:"string"`
	UserDateUpdate           null.String `json:"userDateUpdate" db:"USER_DATE_UPDATE" swaggertype:"string"`
	AdminUpdate              null.String `json:"adminUpdate" db:"ADMIN_UPDATE" swaggertype:"string"`
	AdminDateUpdate          null.String `json:"adminDateUpdate" db:"ADMIN_DATE_UPDATE" swaggertype:"string"`
	FlightArrivalDateExcel   string      `json:"flightArrivalDateExcel"`
	RequestedPickupDateExcel string      `json:"requestedPickupDateExcel"`
	AdminDateUpdateExcel     string      `json:"adminDateUpdateExcel"`
	UserDateUpdateExcel      string      `json:"userDateUpdateExcel"`
	DateCreateExcel          string      `json:"dateCreateExcel"`
}

func (o *LogisticRequest) Set() {
	if o.LogisticRequestStatus.String == "Requested" || o.LogisticRequestStatus.String == "Confirmed" {
		o.DisplayCancelBtn = true
	} else {
		o.DisplayCancelBtn = false
	}

	if !o.PrimaryDoctor.Valid {
		o.PrimaryDoctorName = "-"
	} else {
		o.PrimaryDoctorName = o.PrimaryDoctor.String
	}
}

func (o *LogisticRequest) SetWebAdmin() {
	if !o.RequesterPrn.Valid {
		o.RequesterPrn = utils.NewNullString("-")
	}

	if !o.PrimaryDoctor.Valid {
		o.PrimaryDoctor = utils.NewNullString("-")
	}

	if !o.CompanionName.Valid {
		o.CompanionName = utils.NewNullString("-")
	}

	if !o.UserUpdate.Valid {
		o.UserUpdate = utils.NewNullString("-")
	}

	if !o.UserDateUpdate.Valid {
		o.UserDateUpdate = utils.NewNullString("-")
	}

	if !o.AdminUpdate.Valid {
		o.AdminUpdate = utils.NewNullString("-")
	}

	if !o.AdminDateUpdate.Valid {
		o.AdminDateUpdate = utils.NewNullString("-")
	}
}
