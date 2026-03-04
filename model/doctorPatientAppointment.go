package model

import (
    "github.com/guregu/null/v6"
)

type DoctorPatientAppointment struct {
    DoctorPatientApptId null.Int64  `json:"doctor_patient_appt_id" db:"DOCTOR_PATIENT_APPT_ID" swaggertype:"integer"`
    DoctorID            null.Int64  `json:"doctor_id" db:"DOCTOR_ID" swaggertype:"integer"`
    DoctorName          null.String `json:"doctorName" db:"DOCTOR_NAME" swaggertype:"string"`
    DoctorSpecialty     null.String `json:"doctorSpecialty" db:"DOCTOR_SPECIALTY" swaggertype:"string"`
    ApptStatus          null.String `json:"apptStatus" db:"APPT_STATUS" swaggertype:"string"`
    ApptNo              null.String `json:"apptNo" db:"APPT_NO" swaggertype:"string"`
    ApptDay             null.String `json:"apptDay" db:"APPT_DAY" swaggertype:"string"`
    ApptSessionType     null.String `json:"apptSessionType" db:"APPT_SESSIONTYPE" swaggertype:"string"`
    ApptClinic          null.String `json:"apptClinic" db:"APPT_CLINIC" swaggertype:"string"`
    ApptRoom            null.String `json:"apptRoom" db:"APPT_ROOM" swaggertype:"string"`
    ApptCasetype        null.String `json:"apptCasetype" db:"APPT_CASETYPE" swaggertype:"string"`
    DateAppt            null.String `json:"dateAppt" db:"DATE_APPT" swaggertype:"string"`
    PatientID           null.Int64  `json:"patientId" db:"PATIENT_ID" swaggertype:"integer"`
    PatientPrn          null.String `json:"patientPrn" db:"PATIENT_PRN" swaggertype:"string"`
    PatientName         null.String `json:"patientName" db:"PATIENT_NAME" swaggertype:"string"`
}
