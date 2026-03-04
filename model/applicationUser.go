package model

import (
    "github.com/guregu/null/v6"
)

type ApplicationUser struct {
    UserID               null.Int64     `json:"user_id" db:"USER_ID" swaggertype:"integer"`
    Username             null.String    `json:"username" db:"USERNAME" swaggertype:"string"`
    Email                null.String    `json:"email" db:"EMAIL" swaggertype:"string"`
    IsKidsExplorer       null.String    `json:"isKidsExplorer" db:"IS_KIDS_EXPLORER" swaggertype:"string"`
    IsGoldenPearl        null.String    `json:"isGoldenPearl" db:"IS_GOLDEN_PEARL" swaggertype:"string"`
    Password             null.String    `json:"-" db:"PASSWORD" swaggertype:"string"`
    Title                null.String    `json:"title" db:"TITLE" swaggertype:"string"`
    FirstName            null.String    `json:"firstName" db:"FIRST_NAME" swaggertype:"string"`
    MiddleName           null.String    `json:"middleName" db:"MIDDLE_NAME" swaggertype:"string"`
    LastName             null.String    `json:"lastName" db:"LAST_NAME" swaggertype:"string"`
    FullName             null.String    `json:"fullName" db:"FULL_NAME" swaggertype:"string"`
    Resident             null.String    `json:"resident" db:"RESIDENT" swaggertype:"string"`
    Dob                  null.String    `json:"dob" db:"DOB" swaggertype:"string"`
    Sex                  null.String    `json:"sex" db:"SEX" swaggertype:"string"`
    Race                 null.String    `json:"race" db:"RACE" swaggertype:"string"`
    Address              null.String    `json:"address" db:"ADDRESS" swaggertype:"string"`
    Address1             null.String    `json:"address1" db:"ADDRESS_1" swaggertype:"string"`
    Address2             null.String    `json:"address2" db:"ADDRESS_2" swaggertype:"string"`
    Address3             null.String    `json:"address3" db:"ADDRESS_3" swaggertype:"string"`
    Citystate            null.String    `json:"cityState" db:"CITYSTATE" swaggertype:"string"`
    Postcode             null.String    `json:"postalCode" db:"POSTCODE" swaggertype:"string"`
    Country              null.String    `json:"country" db:"COUNTRY" swaggertype:"string"`
    ContactNumber        null.String    `json:"contactNumber" db:"CONTACT_NUMBER" swaggertype:"string"`
    Passport             null.String    `json:"passport" db:"PASSPORT" swaggertype:"string"`
    Nationality          null.String    `json:"nationality" db:"NATIONALITY" swaggertype:"string"`
    VerificationCode     null.String    `json:"verificationCode" db:"VERIFICATION_CODE" swaggertype:"string"`
    FirstTimeLoginV      null.Int64     `json:"-" db:"FIRST_TIME_LOGIN" swaggertype:"integer"`
    FirstTimeBiometricV  null.Int64     `json:"-" db:"FIRST_TIME_BIOMETRIC" swaggertype:"integer"`
    FirstTimeLogin       bool           `json:"firstTimeLogin"`
    FirstTimeBiometric   bool           `json:"firstTimeBiometric"`
    Role                 null.String    `json:"role" db:"ROLE" swaggertype:"string"`
    MasterPrn            null.String    `json:"masterPrn" db:"MASTER_PRN" swaggertype:"string"`
    PlayerID             null.String    `json:"playerId" db:"PLAYER_ID" swaggertype:"string"`
    MachineID            null.String    `json:"machineId" db:"MACHINE_ID" swaggertype:"string"`
    RegistrationDateTime null.String    `json:"registration_date_time" db:"REGISTRATION_DATE_TIME" swaggertype:"string"`
    InactiveFlag         null.String    `json:"inactive" db:"INACTIVE_FLAG" swaggertype:"string"`
    IsLoggedIn           null.Int64     `json:"isLoggedIn" db:"IS_LOGGED_IN" swaggertype:"integer"`
    DateLoggedIn         null.String    `json:"dateLoggedIn" db:"DATE_LOGGED_IN" swaggertype:"string"`
    SignInType           null.Int32     `json:"signInType" db:"SIGN_IN_TYPE" swaggertype:"integer"`
    DocNoSignup          null.String    `json:"docNoSignUp" db:"DOC_NO_SIGNUP" swaggertype:"string"`
    FullnameSignup       null.String    `json:"fullNameSignUp" db:"FULLNAME_SIGNUP" swaggertype:"string"`
    UserBranches         []AssignBranch `json:"userBranches"`
    SessionID            null.String    `json:"sessionId" db:"SESSION_ID" swaggertype:"string"`
    Branch               null.String    `json:"branch" db:"BRANCH" swaggertype:"string"`
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
