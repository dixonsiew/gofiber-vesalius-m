package model

import (
    "database/sql"
)

type DbDoctorPatientAppointment struct {
    DoctorPatientApptId sql.NullInt64  `db:"DOCTOR_PATIENT_APPT_ID"`
    DoctorID            sql.NullInt64  `db:"DOCTOR_ID"`
    DoctorName          sql.NullString `db:"DOCTOR_NAME"`
    DoctorSpecialty     sql.NullString `db:"DOCTOR_SPECIALTY"`
    ApptStatus          sql.NullString `db:"APPT_STATUS"`
    ApptNo              sql.NullString `db:"APPT_NO"`
    ApptDay             sql.NullString `db:"APPT_DAY"`
    ApptSessionType     sql.NullString `db:"APPT_SESSIONTYPE"`
    ApptClinic          sql.NullString `db:"APPT_CLINIC"`
    ApptRoom            sql.NullString `db:"APPT_ROOM"`
    ApptCasetype        sql.NullString `db:"APPT_CASETYPE"`
    DateAppt            sql.NullString `db:"DATE_APPT"`
    PatientID           sql.NullInt64  `db:"PATIENT_ID"`
    PatientPrn          sql.NullString `db:"PATIENT_PRN"`
    PatientName         sql.NullString `db:"PATIENT_NAME"`
}

type DoctorPatientAppointment struct {
    DoctorPatientApptId int64  `json:"doctor_patient_appt_id"`
    DoctorID            int64  `json:"doctor_id"`
    DoctorName          string `json:"doctorName"`
    DoctorSpecialty     string `json:"doctorSpecialty"`
    ApptStatus          string `json:"apptStatus"`
    ApptNo              string `json:"apptNo"`
    ApptDay             string `json:"apptDay"`
    ApptSessionType     string `json:"apptSessionType"`
    ApptClinic          string `json:"apptClinic"`
    ApptRoom            string `json:"apptRoom"`
    ApptCasetype        string `json:"apptCasetype"`
    DateAppt            string `json:"dateAppt"`
    PatientID           int64  `json:"patientId"`
    PatientPrn          string `json:"patientPrn"`
    PatientName         string `json:"patientName"`
}

func (o *DoctorPatientAppointment) FromDbModel(m DbDoctorPatientAppointment) {
    o.DoctorPatientApptId = m.DoctorPatientApptId.Int64
    o.DoctorID = m.DoctorID.Int64
    o.DoctorName = m.DoctorName.String
    o.DoctorSpecialty = m.DoctorSpecialty.String
    o.ApptStatus = m.ApptStatus.String
    o.ApptNo = m.ApptNo.String
    o.ApptDay = m.ApptDay.String
    o.ApptSessionType = m.ApptSessionType.String
    o.ApptClinic = m.ApptClinic.String
    o.ApptRoom = m.ApptRoom.String
    o.ApptCasetype = m.ApptCasetype.String
    o.DateAppt = m.DateAppt.String
    o.PatientID = m.PatientID.Int64
    o.PatientPrn = m.PatientPrn.String
    o.PatientName = m.PatientName.String
}
