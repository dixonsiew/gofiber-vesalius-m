package model

import (
    "database/sql"
)

type DbFutureOrder struct {
    PatientName  sql.NullString `db:"PATIENT_NAME"`
    Prn          sql.NullString `db:"PRN"`
    PlanType     sql.NullString `db:"PLAN_TYPE"`
    PerformDate  sql.NullString `db:"START_DATE_TIME"`
    OrderDoctor  sql.NullString `db:"ORDER_DOCTOR"`
    Description  sql.NullString `db:"DESCRIPTION"`
}

type FutureOrder struct {
    PatientName  string `json:"patientName"`
    Prn          string `json:"prn"`
    PlanType     string `json:"planType"`
    PerformDate  string `json:"performDate"`
    OrderDoctor  string `json:"orderDoctor"`
    Description  string `json:"description"`
}

func (o *FutureOrder) FromDbModel(m DbFutureOrder) {
    o.PatientName = m.PatientName.String
    o.Prn = m.Prn.String
    o.PlanType = m.PlanType.String
    
    if o.PlanType == "DATE" {
        o.PerformDate = m.PerformDate.String
    } else {
        o.PerformDate = "Next Visit"
    }

    o.OrderDoctor = m.OrderDoctor.String
    o.Description = m.Description.String
}
