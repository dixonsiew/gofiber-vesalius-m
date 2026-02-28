package model

import (
    "database/sql"
)

type DbApplicationUser struct {
    UserID               sql.NullInt64  `db:"USER_ID"`
    Username             sql.NullString `db:"USERNAME"`
    Email                sql.NullString `db:"EMAIL"`
    IsKidsExplorer       sql.NullString `db:"IS_KIDS_EXPLORER"`
    IsGoldenPearl        sql.NullString `db:"IS_GOLDEN_PEARL"`
    Password             sql.NullString `db:"PASSWORD"`
    Title                sql.NullString `db:"TITLE"`
    FirstName            sql.NullString `db:"FIRST_NAME"`
    MiddleName           sql.NullString `db:"MIDDLE_NAME"`
    LastName             sql.NullString `db:"LAST_NAME"`
    Resident             sql.NullString `db:"RESIDENT"`
    Dob                  sql.NullString `db:"DOB"`
    Sex                  sql.NullString `db:"SEX"`
    Race                 sql.NullString `db:"RACE"`
    Address              sql.NullString `db:"ADDRESS"`
    Address1             sql.NullString `db:"ADDRESS_1"`
    Address2             sql.NullString `db:"ADDRESS_2"`
    Address3             sql.NullString `db:"ADDRESS_3"`
    Citystate            sql.NullString `db:"CITYSTATE"`
    Postcode             sql.NullString `db:"POSTCODE"`
    Country              sql.NullString `db:"COUNTRY"`
    ContactNumber        sql.NullString `db:"CONTACT_NUMBER"`
    Passport             sql.NullString `db:"PASSPORT"`
    Nationality          sql.NullString `db:"NATIONALITY"`
    VerificationCode     sql.NullString `db:"VERIFICATION_CODE"`
    FirstTimeLogin       sql.NullInt64  `db:"FIRST_TIME_LOGIN"`
    FirstTimeBiometric   sql.NullInt64  `db:"FIRST_TIME_BIOMETRIC"`
    Role                 sql.NullString `db:"ROLE"`
    MasterPrn            sql.NullString `db:"MASTER_PRN"`
    PlayerID             sql.NullString `db:"PLAYER_ID"`
    MachineID            sql.NullString `db:"MACHINE_ID"`
    RegistrationDateTime sql.NullString `db:"REGISTRATION_DATE_TIME"`
    InactiveFlag         sql.NullString `db:"INACTIVE_FLAG"`
    IsLoggedIn           sql.NullInt64  `db:"IS_LOGGED_IN"`
    DateLoggedIn         sql.NullString `db:"DATE_LOGGED_IN"`
    SignInType           sql.NullInt32  `db:"SIGN_IN_TYPE"`
    DocNoSignup          sql.NullString `db:"DOC_NO_SIGNUP"`
    FullnameSignup       sql.NullString `db:"FULLNAME_SIGNUP"`
    SessionID            sql.NullString `db:"SESSION_ID"`
    Branch               sql.NullString `db:"BRANCH"`
}

type ApplicationUser struct {
    UserID               int64          `json:"user_id"`
    Username             string         `json:"username"`
    Email                string         `json:"email"`
    IsKidsExplorer       string         `json:"isKidsExplorer"`
    IsGoldenPearl        string         `json:"isGoldenPearl"`
    Password             string         `json:"-"`
    Title                string         `json:"title"`
    FirstName            string         `json:"firstName"`
    MiddleName           string         `json:"middleName"`
    LastName             string         `json:"lastName"`
    FullName             string         `json:"fullName"`
    Resident             string         `json:"resident"`
    Dob                  string         `json:"dob"`
    Sex                  string         `json:"sex"`
    Race                 string         `json:"race"`
    Address              string         `json:"address"`
    Address1             string         `json:"address1"`
    Address2             string         `json:"address2"`
    Address3             string         `json:"address3"`
    Citystate            string         `json:"cityState"`
    Postcode             string         `json:"postalCode"`
    Country              string         `json:"country"`
    ContactNumber        string         `json:"contactNumber"`
    Passport             string         `json:"passport"`
    Nationality          string         `json:"nationality"`
    VerificationCode     string         `json:"verificationCode"`
    FirstTimeLogin       bool           `json:"firstTimeLogin"`
    FirstTimeBiometric   bool           `json:"firstTimeBiometric"`
    Role                 string         `json:"role"`
    MasterPrn            string         `json:"masterPrn"`
    PlayerID             string         `json:"playerId"`
    MachineID            string         `json:"machineId"`
    RegistrationDateTime string         `json:"registration_date_time"`
    InactiveFlag         string         `json:"inactive"`
    IsLoggedIn           int64          `json:"isLoggedIn"`
    DateLoggedIn         string         `json:"dateLoggedIn"`
    SignInType           int32          `json:"signInType"`
    DocNoSignup          string         `json:"docNoSignUp"`
    FullnameSignup       string         `json:"fullNameSignUp"`
    UserBranches         []AssignBranch `json:"userBranches"`
    SessionID            string         `json:"sessionId"`
    Branch               string         `json:"branch"`
}

func (o *ApplicationUserFamily) FromRsFamilyMember(m DbApplicationUser) {
    o.AufID = 0
    o.UserID = m.UserID.Int64
    o.NokRefNumber = 0
    o.IsPatient = true
    o.Fullname = m.FirstName.String
    o.Relationship = "Self"
    o.NokPrn = m.MasterPrn.String
    o.NricPassport = "-"
    o.DocNumber = "-"
    o.Dob = "-"
    o.Gender = m.Sex.String
    o.Nationality = "-"
    o.ContactNumber = "-"
    o.Address = "-"
    o.IsActive = true
    o.MaritalStatus = "-"
    o.Email = "-"
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

func (o *ApplicationUser) FromDbModel(m DbApplicationUser) {
    o.UserID = m.UserID.Int64
    o.Username = m.Username.String
    if m.Email.Valid {
        o.Email = m.Email.String
    } else {
        o.Email = "-"
    }
    o.Password = m.Password.String
    o.Title = m.Title.String
    o.FirstName = m.FirstName.String
    o.MiddleName = m.MiddleName.String
    o.LastName = m.LastName.String
    o.Resident = m.Resident.String
    o.Dob = m.Dob.String
    o.Sex = m.Sex.String
    if m.Race.Valid {
        o.Race = m.Race.String
    } else {
        o.Race = "-"
    }
    o.Address = m.Address.String
    o.Address1 = m.Address1.String
    o.Address2 = m.Address2.String
    o.Address3 = m.Address3.String
    o.Citystate = m.Citystate.String
    o.Postcode = m.Postcode.String
    o.Country = m.Country.String
    if m.ContactNumber.Valid {
        o.ContactNumber = m.ContactNumber.String
    } else {
        o.ContactNumber = "-"
    }
    o.Passport = m.Passport.String
    o.Nationality = m.Nationality.String
    o.VerificationCode = m.VerificationCode.String
    if m.FirstTimeLogin.Int64 == 1 {
        o.FirstTimeLogin = true
    } else {
        o.FirstTimeLogin = false
    }
    if m.FirstTimeBiometric.Int64 == 1 {
        o.FirstTimeBiometric = true
    } else {
        o.FirstTimeBiometric = false
    }
    o.Role = m.Role.String
    o.MasterPrn = m.MasterPrn.String
    o.PlayerID = m.PlayerID.String
    o.MachineID = m.MachineID.String
    o.RegistrationDateTime = m.RegistrationDateTime.String
    o.InactiveFlag = m.InactiveFlag.String
    o.SessionID = m.SessionID.String
    o.SignInType = m.SignInType.Int32
}
