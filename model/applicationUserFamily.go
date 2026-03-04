package model

import (
    "github.com/guregu/null/v6"
    "strconv"
)

type ApplicationUserFamily struct {
    AufID           null.Int64  `json:"auf_id" db:"AUF_ID" swaggertype:"integer"`
    UserID          null.Int64  `json:"user_id" db:"USER_ID" swaggertype:"integer"`
    PatientPrn      null.String `json:"patientPrn" db:"PATIENT_PRN" swaggertype:"string"`
    NokRefNumberV   null.String `json:"-" db:"NOK_REF_NUMBER" swaggertype:"string"`
    IsPatientV      null.String `json:"-" db:"IS_PATIENT" swaggertype:"string"`
    NokRefNumber    int64       `json:"nokRefNumber" swaggertype:"integer"`
    IsPatient       bool        `json:"isPatient"`
    Fullname        null.String `json:"fullName" db:"FULLNAME" swaggertype:"string"`
    Relationship    null.String `json:"relationship" db:"RELATIONSHIP" swaggertype:"string"`
    NokPrn          null.String `json:"prn" db:"NOK_PRN" swaggertype:"string"`
    DocNumber       null.String `json:"docNumber" db:"DOC_NUMBER" swaggertype:"string"`
    NricPassport    null.String `json:"nricPassport" db:"NRIC_PASSPORT" swaggertype:"string"`
    Dob             null.String `json:"dob" db:"DOB" swaggertype:"string"`
    Gender          null.String `json:"gender" db:"GENDER" swaggertype:"string"`
    Nationality     null.String `json:"nationality" db:"NATIONALITY" swaggertype:"string"`
    ContactNumber   null.String `json:"contactNumber" db:"CONTACT_NUMBER" swaggertype:"string"`
    Address         null.String `json:"address" db:"ADDRESS" swaggertype:"string"`
    MaritalStatus   null.String `json:"maritalStatus" db:"MARITAL_STATUS" swaggertype:"string"`
    Email           null.String `json:"email" db:"EMAIL" swaggertype:"string"`
    IsActiveV       null.String `json:"-" db:"IS_ACTIVE" swaggertype:"string"`
    IsGoldenPearlV  null.String `json:"-" db:"IS_GOLDEN_PEARL" swaggertype:"string"`
    IsKidsExplorerV null.String `json:"-" db:"IS_KIDS_EXPLORER" swaggertype:"string"`
    IsActive        bool        `json:"isActive"`
    IsGoldenPearl   bool        `json:"isGoldenPearl"`
    IsKidsExplorer  bool        `json:"isKidsExplorer"`
    DateCreate      null.String `json:"dateCreate" db:"DATE_CREATE" swaggertype:"string"`
    DateLastSync    null.String `json:"dateLastSync" db:"DATE_LAST_SYNC" swaggertype:"string"`
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
