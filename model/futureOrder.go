package model

import (
    "github.com/guregu/null/v6"
)

type FutureOrder struct {
    PatientName  null.String `json:"patientName" db:"PATIENT_NAME"`
    Prn          null.String `json:"prn" db:"PRN"`
    PlanType     null.String `json:"planType" db:"PLAN_TYPE"`
    PerformDate  null.String `json:"performDate" db:"START_DATE_TIME"`
    OrderDoctor  null.String `json:"orderDoctor" db:"ORDER_DOCTOR"`
    Description  null.String `json:"description" db:"DESCRIPTION"`
}

func (o *FutureOrder) Set() {
    if o.PlanType.String != "DATE" {
        o.PerformDate.String = "Next Visit"
    }
}
