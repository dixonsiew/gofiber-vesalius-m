package model

import (
    "github.com/guregu/null/v6"
    "strconv"
)

type ApplicationUserFamily struct {
    AufID           null.Int64  `json:"auf_id" db:"AUF_ID"`
    UserID          null.Int64  `json:"user_id" db:"USER_ID"`
    PatientPrn      null.String `json:"patientPrn" db:"PATIENT_PRN"`
    NokRefNumberV   null.String `json:"-" db:"NOK_REF_NUMBER"`
    IsPatientV      null.String `json:"-" db:"IS_PATIENT"`
    NokRefNumber    int64       `json:"nokRefNumber"`
    IsPatient       bool        `json:"isPatient"`
    Fullname        null.String `json:"fullName" db:"FULLNAME"`
    Relationship    null.String `json:"relationship" db:"RELATIONSHIP"`
    NokPrn          null.String `json:"prn" db:"NOK_PRN"`
    DocNumber       null.String `json:"docNumber" db:"DOC_NUMBER"`
    NricPassport    null.String `json:"nricPassport" db:"NRIC_PASSPORT"`
    Dob             null.String `json:"dob" db:"DOB"`
    Gender          null.String `json:"gender" db:"GENDER"`
    Nationality     null.String `json:"nationality" db:"NATIONALITY"`
    ContactNumber   null.String `json:"contactNumber" db:"CONTACT_NUMBER"`
    Address         null.String `json:"address" db:"ADDRESS"`
    MaritalStatus   null.String `json:"maritalStatus" db:"MARITAL_STATUS"`
    Email           null.String `json:"email" db:"EMAIL"`
    IsActiveV       null.String `json:"-" db:"IS_ACTIVE"`
    IsGoldenPearlV  null.String `json:"-" db:"IS_GOLDEN_PEARL"`
    IsKidsExplorerV null.String `json:"-" db:"IS_KIDS_EXPLORER"`
    IsActive        bool        `json:"isActive"`
    IsGoldenPearl   bool        `json:"isGoldenPearl"`
    IsKidsExplorer  bool        `json:"isKidsExplorer"`
    DateCreate      null.String `json:"dateCreate" db:"DATE_CREATE"`
    DateLastSync    null.String `json:"dateLastSync" db:"DATE_LAST_SYNC"`
}

func (o *ApplicationUserFamily) Set() {
    i, _ := strconv.ParseInt(o.NokRefNumberV.String, 10, 64)
    o.NokRefNumber = i
    if o.IsPatientV.String == "Y" {
        o.IsPatient = true
    } else {
        o.IsPatient = false
    }
    if !o.Relationship.Valid {
        o.Relationship.String = "-"
    }
    if !o.NokPrn.Valid {
        o.NokPrn.String = "-"
    }
    if !o.NricPassport.Valid {
        o.NricPassport.String = "-"
    }
    if !o.Dob.Valid {
        o.Dob.String = "-"
    }
    if !o.Nationality.Valid {
        o.Nationality.String = "-"
    }
    if !o.ContactNumber.Valid {
        o.ContactNumber.String = "-"
    }
    if !o.Address.Valid {
        o.Address.String = "-"
    }
    if o.IsActiveV.String == "Y" {
        o.IsActive = true
    } else {
        o.IsActive = false
    }
    if !o.MaritalStatus.Valid {
        o.MaritalStatus.String = "-"
    }
    if !o.Email.Valid {
        o.Email.String = "-"
    }
    if o.IsKidsExplorerV.String == "Y" {
        o.IsKidsExplorer = true
    } else {
        o.IsKidsExplorer = false
    }
    if o.IsGoldenPearlV.String == "Y" {
        o.IsGoldenPearl = true
    } else {
        o.IsGoldenPearl = false
    }
}
