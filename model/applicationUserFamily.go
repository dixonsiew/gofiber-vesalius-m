package model

import (
    "database/sql"
    "strconv"
)

type DbApplicationUserFamily struct {
    AufID          sql.NullInt64  `db:"AUF_ID"`
    UserID         sql.NullInt64  `db:"USER_ID"`
    PatientPrn     sql.NullString `db:"PATIENT_PRN"`
    NokRefNumber   sql.NullString `db:"NOK_REF_NUMBER"`
    IsPatient      sql.NullString `db:"IS_PATIENT"`
    Fullname       sql.NullString `db:"FULLNAME"`
    Relationship   sql.NullString `db:"RELATIONSHIP"`
    NokPrn         sql.NullString `db:"NOK_PRN"`
    DocNumber      sql.NullString `db:"DOC_NUMBER"`
    NricPassport   sql.NullString `db:"NRIC_PASSPORT"`
    Dob            sql.NullString `db:"DOB"`
    Gender         sql.NullString `db:"GENDER"`
    Nationality    sql.NullString `db:"NATIONALITY"`
    ContactNumber  sql.NullString `db:"CONTACT_NUMBER"`
    Address        sql.NullString `db:"ADDRESS"`
    MaritalStatus  sql.NullString `db:"MARITAL_STATUS"`
    Email          sql.NullString `db:"EMAIL"`
    IsActive       sql.NullString `db:"IS_ACTIVE"`
    IsGoldenPearl  sql.NullString `db:"IS_GOLDEN_PEARL"`
    IsKidsExplorer sql.NullString `db:"IS_KIDS_EXPLORER"`
    DateCreate     sql.NullString `db:"DATE_CREATE"`
    DateLastSync   sql.NullString `db:"DATE_LAST_SYNC"`
}

type ApplicationUserFamily struct {
    AufID          int64  `json:"auf_id"`
    UserID         int64  `json:"user_id"`
    PatientPrn     string `json:"patientPrn"`
    NokRefNumber   int64  `json:"nokRefNumber"`
    IsPatient      bool   `json:"isPatient"`
    Fullname       string `json:"fullName"`
    Relationship   string `json:"relationship"`
    NokPrn         string `json:"prn"`
    DocNumber      string `json:"docNumber"`
    NricPassport   string `json:"nricPassport"`
    Dob            string `json:"dob"`
    Gender         string `json:"gender"`
    Nationality    string `json:"nationality"`
    ContactNumber  string `json:"contactNumber"`
    Address        string `json:"address"`
    MaritalStatus  string `json:"maritalStatus"`
    Email          string `json:"email"`
    IsActive       bool   `json:"isActive"`
    IsGoldenPearl  bool   `json:"isGoldenPearl"`
    IsKidsExplorer bool   `json:"isKidsExplorer"`
    DateCreate     string `json:"dateCreate"`
    DateLastSync   string `json:"dateLastSync"`
}

func (o *ApplicationUserFamily) FromDbModel(m DbApplicationUserFamily) {
    i, _ := strconv.ParseInt(m.NokRefNumber.String, 10, 64)
    o.AufID = m.AufID.Int64
    o.UserID = m.UserID.Int64
    o.PatientPrn = m.PatientPrn.String
    o.NokRefNumber = i
    if m.IsPatient.String == "Y" {
        o.IsPatient = true
    } else {
        o.IsPatient = false
    }
    o.Fullname = m.Fullname.String
    if m.Relationship.Valid {
        o.Relationship = m.Relationship.String
    } else {
        o.Relationship = "-"
    }
    if m.NokPrn.Valid {
        o.NokPrn = m.NokPrn.String
    } else {
        o.NokPrn = "-"
    }
    if m.NricPassport.Valid {
        o.NricPassport = m.NricPassport.String
    } else {
        o.NricPassport = "-"
    }
    o.DocNumber = m.DocNumber.String
    if m.Dob.Valid {
        o.Dob = m.Dob.String
    } else {
        o.Dob = "-"
    }
    o.Gender = m.Gender.String
    if m.Nationality.Valid {
        o.Nationality = m.Nationality.String
    } else {
        o.Nationality = "-"
    }
    if m.ContactNumber.Valid {
        o.ContactNumber = m.ContactNumber.String
    } else {
        o.ContactNumber = "-"
    }
    if m.Address.Valid {
        o.Address = m.Address.String
    } else {
        o.Address = "-"
    }
    if m.IsActive.String == "Y" {
        o.IsActive = true
    } else {
        o.IsActive = false
    }
    if m.MaritalStatus.Valid {
        o.MaritalStatus = m.MaritalStatus.String
    } else {
        o.MaritalStatus = "-"
    }
    if m.Email.Valid {
        o.Email = m.Email.String
    } else {
        o.Email = "-"
    }
    if m.IsKidsExplorer.String == "Y" {
        o.IsKidsExplorer = true
    } else {
        o.IsKidsExplorer = false
    }
    if m.IsGoldenPearl.String == "Y" {
        o.IsGoldenPearl = true
    } else {
        o.IsGoldenPearl = false
    }
}
