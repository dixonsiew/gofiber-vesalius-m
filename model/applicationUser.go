package model

import (
    "github.com/guregu/null/v6"
)

type ApplicationUser struct {
    UserID               null.Int64     `json:"user_id" db:"USER_ID"`
    Username             null.String    `json:"username" db:"USERNAME"`
    Email                null.String    `json:"email" db:"EMAIL"`
    IsKidsExplorer       null.String    `json:"isKidsExplorer" db:"IS_KIDS_EXPLORER"`
    IsGoldenPearl        null.String    `json:"isGoldenPearl" db:"IS_GOLDEN_PEARL"`
    Password             null.String    `json:"-" db:"PASSWORD"`
    Title                null.String    `json:"title" db:"TITLE"`
    FirstName            null.String    `json:"firstName" db:"FIRST_NAME"`
    MiddleName           null.String    `json:"middleName" db:"MIDDLE_NAME"`
    LastName             null.String    `json:"lastName" db:"LAST_NAME"`
    FullName             null.String    `json:"fullName" db:"FULL_NAME"`
    Resident             null.String    `json:"resident" db:"RESIDENT"`
    Dob                  null.String    `json:"dob" db:"DOB"`
    Sex                  null.String    `json:"sex" db:"SEX"`
    Race                 null.String    `json:"race" db:"RACE"`
    Address              null.String    `json:"address" db:"ADDRESS"`
    Address1             null.String    `json:"address1" db:"ADDRESS_1"`
    Address2             null.String    `json:"address2" db:"ADDRESS_2"`
    Address3             null.String    `json:"address3" db:"ADDRESS_3"`
    Citystate            null.String    `json:"cityState" db:"CITYSTATE"`
    Postcode             null.String    `json:"postalCode" db:"POSTCODE"`
    Country              null.String    `json:"country" db:"COUNTRY"`
    ContactNumber        null.String    `json:"contactNumber" db:"CONTACT_NUMBER"`
    Passport             null.String    `json:"passport" db:"PASSPORT"`
    Nationality          null.String    `json:"nationality" db:"NATIONALITY"`
    VerificationCode     null.String    `json:"verificationCode" db:"VERIFICATION_CODE"`
    FirstTimeLoginV      null.Int64     `json:"-" db:"FIRST_TIME_LOGIN"`
    FirstTimeBiometricV  null.Int64     `json:"-" db:"FIRST_TIME_BIOMETRIC"`
    FirstTimeLogin       bool           `json:"firstTimeLogin"`
    FirstTimeBiometric   bool           `json:"firstTimeBiometric"`
    Role                 null.String    `json:"role" db:"ROLE"`
    MasterPrn            null.String    `json:"masterPrn" db:"MASTER_PRN"`
    PlayerID             null.String    `json:"playerId" db:"PLAYER_ID"`
    MachineID            null.String    `json:"machineId" db:"MACHINE_ID"`
    RegistrationDateTime null.String    `json:"registration_date_time" db:"REGISTRATION_DATE_TIME"`
    InactiveFlag         null.String    `json:"inactive" db:"INACTIVE_FLAG"`
    IsLoggedIn           null.Int64     `json:"isLoggedIn" db:"IS_LOGGED_IN"`
    DateLoggedIn         null.String    `json:"dateLoggedIn" db:"DATE_LOGGED_IN"`
    SignInType           null.Int32     `json:"signInType" db:"SIGN_IN_TYPE"`
    DocNoSignup          null.String    `json:"docNoSignUp" db:"DOC_NO_SIGNUP"`
    FullnameSignup       null.String    `json:"fullNameSignUp" db:"FULLNAME_SIGNUP"`
    UserBranches         []AssignBranch `json:"userBranches"`
    SessionID            null.String    `json:"sessionId" db:"SESSION_ID"`
    Branch               null.String    `json:"branch" db:"BRANCH"`
}

func (o *ApplicationUserFamily) FromRsFamilyMember(m ApplicationUser) {
    o.AufID.Int64 = 0
    o.UserID.Int64 = m.UserID.Int64
    o.NokRefNumber = 0
    o.IsPatient = true
    o.Fullname.String = m.FirstName.String
    o.Relationship.String = "Self"
    o.NokPrn.String = m.MasterPrn.String
    o.NricPassport.String = "-"
    o.DocNumber.String = "-"
    o.Dob.String = "-"
    o.Gender.String = m.Sex.String
    o.Nationality.String = "-"
    o.ContactNumber.String = "-"
    o.Address.String = "-"
    o.IsActive = true
    o.MaritalStatus.String = "-"
    o.Email.String = "-"
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

func (o *ApplicationUser) Set() {
    if !o.Email.Valid {
        o.Email.String = "-"
    }
    if !o.Race.Valid {
        o.Race.String = "-"
    }
    if !o.ContactNumber.Valid {
        o.ContactNumber.String = "-"
    }
    if o.FirstTimeLoginV.Int64 == 1 {
        o.FirstTimeLogin = true
    } else {
        o.FirstTimeLogin = false
    }
    if o.FirstTimeBiometricV.Int64 == 1 {
        o.FirstTimeBiometric = true
    } else {
        o.FirstTimeBiometric = false
    }
}
