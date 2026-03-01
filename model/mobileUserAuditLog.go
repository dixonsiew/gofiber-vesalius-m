package model

import (
    "database/sql"
)

type DbMobileUserAuditLog struct {
    AuditID     sql.NullInt64  `db:"AUDIT_ID"`
    Prn         sql.NullString `db:"PRN"`
    Username    sql.NullString `db:"USERNAME"`
    PatientName sql.NullString `db:"PATIENT_NAME"`
    Action      sql.NullString `db:"ACTION"`
    ActionDesc  sql.NullString `db:"ACTION_DESC"`
    Remarks     sql.NullString `db:"REMARKS"`
    UserCreate  sql.NullString `db:"USER_CREATE"`
    DateCreate  sql.NullString `db:"DATE_CREATE"`
}

type MobileUserAuditLog struct {
    AuditID     int64  `json:"audit_id"`
    Prn         string `json:"prn"`
    Username    string `json:"username"`
    PatientName string `json:"patientName"`
    Action      string `json:"action"`
    ActionDesc  string `json:"actionDesc"`
    Remarks     string `json:"remarks"`
    UserCreate  string `json:"userCreate"`
    DateCreate  string `json:"dateCreate"`
}

func (o *MobileUserAuditLog) FromDbModel(m DbMobileUserAuditLog) {
    o.AuditID = m.AuditID.Int64
    o.Prn = m.Prn.String
    o.Username = m.Username.String
    o.PatientName = m.PatientName.String
    o.Action = m.Action.String
    o.ActionDesc = m.ActionDesc.String
    o.Remarks = m.Remarks.String
    o.UserCreate = m.UserCreate.String
    o.DateCreate = m.DateCreate.String
}
