package model

import (
    "github.com/guregu/null/v6"
)

type MobileUserAuditLog struct {
    AuditID     null.Int64  `json:"audit_id" db:"AUDIT_ID" swaggertype:"integer"`
    Prn         null.String `json:"prn" db:"PRN" swaggertype:"string"`
    Username    null.String `json:"username" db:"USERNAME" swaggertype:"string"`
    PatientName null.String `json:"patientName" db:"PATIENT_NAME" swaggertype:"string"`
    Action      null.String `json:"action" db:"ACTION" swaggertype:"string"`
    ActionDesc  null.String `json:"actionDesc" db:"ACTION_DESC" swaggertype:"string"`
    Remarks     null.String `json:"remarks" db:"REMARKS" swaggertype:"string"`
    UserCreate  null.String `json:"userCreate" db:"USER_CREATE" swaggertype:"string"`
    DateCreate  null.String `json:"dateCreate" db:"DATE_CREATE" swaggertype:"string"`
}
