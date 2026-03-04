package model

import (
    "github.com/guregu/null/v6"
)

type MobileUserAuditLog struct {
    AuditID     null.Int64  `json:"audit_id" db:"AUDIT_ID"`
    Prn         null.String `json:"prn" db:"PRN"`
    Username    null.String `json:"username" db:"USERNAME"`
    PatientName null.String `json:"patientName" db:"PATIENT_NAME"`
    Action      null.String `json:"action" db:"ACTION"`
    ActionDesc  null.String `json:"actionDesc" db:"ACTION_DESC"`
    Remarks     null.String `json:"remarks" db:"REMARKS"`
    UserCreate  null.String `json:"userCreate" db:"USER_CREATE"`
    DateCreate  null.String `json:"dateCreate" db:"DATE_CREATE"`
}
