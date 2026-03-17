package futureOrder

import (
    "vesaliusm/utils"

    "github.com/guregu/null/v6"
)

type FutureOrder struct {
    PatientName  null.String `json:"patientName" db:"PATIENT_NAME" swaggertype:"string"`
    Prn          null.String `json:"prn" db:"PRN" swaggertype:"string"`
    PlanType     null.String `json:"planType" db:"PLAN_TYPE" swaggertype:"string"`
    PerformDate  null.String `json:"performDate" db:"START_DATE_TIME" swaggertype:"string"`
    OrderDoctor  null.String `json:"orderDoctor" db:"ORDER_DOCTOR" swaggertype:"string"`
    Description  null.String `json:"description" db:"DESCRIPTION" swaggertype:"string"`
}

func (o *FutureOrder) Set() {
    if o.PlanType.String != "DATE" {
        o.PerformDate = utils.NewNullString("Next Visit")
    }
}
