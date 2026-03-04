package model

import (
    "github.com/guregu/null/v6"
)

type DoctorPatientAppointment struct {
    DoctorPatientApptId null.Int64  `json:"doctor_patient_appt_id" db:"DOCTOR_PATIENT_APPT_ID"`
    DoctorID            null.Int64  `json:"doctor_id" db:"DOCTOR_ID"`
    DoctorName          null.String `json:"doctorName" db:"DOCTOR_NAME"`
    DoctorSpecialty     null.String `json:"doctorSpecialty" db:"DOCTOR_SPECIALTY"`
    ApptStatus          null.String `json:"apptStatus" db:"APPT_STATUS"`
    ApptNo              null.String `json:"apptNo" db:"APPT_NO"`
    ApptDay             null.String `json:"apptDay" db:"APPT_DAY"`
    ApptSessionType     null.String `json:"apptSessionType" db:"APPT_SESSIONTYPE"`
    ApptClinic          null.String `json:"apptClinic" db:"APPT_CLINIC"`
    ApptRoom            null.String `json:"apptRoom" db:"APPT_ROOM"`
    ApptCasetype        null.String `json:"apptCasetype" db:"APPT_CASETYPE"`
    DateAppt            null.String `json:"dateAppt" db:"DATE_APPT"`
    PatientID           null.Int64  `json:"patientId" db:"PATIENT_ID"`
    PatientPrn          null.String `json:"patientPrn" db:"PATIENT_PRN"`
    PatientName         null.String `json:"patientName" db:"PATIENT_NAME"`
}
