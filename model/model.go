package model

import (
    "database/sql"
    "strconv"
)

type DbApplicationUser struct {
    UserID               sql.NullInt64  `db:"USER_ID"`
    Address              sql.NullString `db:"ADDRESS"`
    ContactNumber        sql.NullString `db:"CONTACT_NUMBER"`
    Dob                  sql.NullString `db:"DOB"`
    Email                sql.NullString `db:"EMAIL"`
    FirstTimeLogin       sql.NullInt64  `db:"FIRST_TIME_LOGIN"`
    FirstName            sql.NullString `db:"FIRST_NAME"`
    LastName             sql.NullString `db:"LAST_NAME"`
    MasterPrn            sql.NullString `db:"MASTER_PRN"`
    MiddleName           sql.NullString `db:"MIDDLE_NAME"`
    Nationality          sql.NullString `db:"NATIONALITY"`
    Passport             sql.NullString `db:"PASSPORT"`
    Password             sql.NullString `db:"PASSWORD"`
    Resident             sql.NullString `db:"RESIDENT"`
    Role                 sql.NullString `db:"ROLE"`
    Sex                  sql.NullString `db:"SEX"`
    Title                sql.NullString `db:"TITLE"`
    Username             sql.NullString `db:"USERNAME"`
    VerificationCode     sql.NullString `db:"VERIFICATION_CODE"`
    Branch               sql.NullString `db:"BRANCH"`
    PlayerID             sql.NullString `db:"PLAYER_ID"`
    RegistrationDateTime sql.NullString `db:"REGISTRATION_DATE_TIME"`
    InactiveFlag         sql.NullString `db:"INACTIVE_FLAG"`
    MachineID            sql.NullString `db:"MACHINE_ID"`
    Race                 sql.NullString `db:"RACE"`
    FirstTimeBiometric   sql.NullInt64  `db:"FIRST_TIME_BIOMETRIC"`
    IsLoggedIn           sql.NullInt64  `db:"IS_LOGGED_IN"`
    DateLoggedIn         sql.NullString `db:"DATE_LOGGED_IN"`
    SessionID            sql.NullString `db:"SESSION_ID"`
    IsGoldenPearl        sql.NullString `db:"IS_GOLDEN_PEARL"`
    IsKidsExplorer       sql.NullString `db:"IS_KIDS_EXPLORER"`
    Address1             sql.NullString `db:"ADDRESS_1"`
    Address2             sql.NullString `db:"ADDRESS_2"`
    Address3             sql.NullString `db:"ADDRESS_3"`
    Country              sql.NullString `db:"COUNTRY"`
    Postcode             sql.NullString `db:"POSTCODE"`
    Citystate            sql.NullString `db:"CITYSTATE"`
    SignInType           sql.NullInt32  `db:"SIGN_IN_TYPE"`
    DocNoSignup          sql.NullString `db:"DOC_NO_SIGNUP"`
    FullnameSignup       sql.NullString `db:"FULLNAME_SIGNUP"`
}

type ApplicationUser struct {
    UserID               int64          `json:"user_id"`
    Address              string         `json:"address"`
    ContactNumber        string         `json:"contactNumber"`
    Dob                  string         `json:"dob"`
    Email                string         `json:"email"`
    FirstTimeLogin       bool           `json:"firstTimeLogin"`
    FirstName            string         `json:"firstName"`
    LastName             string         `json:"lastName"`
    FullName             string         `json:"fullName"`
    MasterPrn            string         `json:"masterPrn"`
    MiddleName           string         `json:"middleName"`
    Nationality          string         `json:"nationality"`
    Passport             string         `json:"passport"`
    Password             string         `json:"password"`
    Resident             string         `json:"resident"`
    Role                 string         `json:"role"`
    Sex                  string         `json:"sex"`
    Title                string         `json:"title"`
    Username             string         `json:"username"`
    VerificationCode     string         `json:"verificationCode"`
    Branch               string         `json:"branch"`
    PlayerID             string         `json:"playerId"`
    RegistrationDateTime string         `json:"registration_date_time"`
    InactiveFlag         string         `json:"inactive"`
    MachineID            string         `json:"machineId"`
    Race                 string         `json:"race"`
    FirstTimeBiometric   bool           `json:"firstTimeBiometric"`
    IsLoggedIn           int64          `json:"isLoggedIn"`
    DateLoggedIn         string         `json:"dateLoggedIn"`
    UserBranches         []AssignBranch `json:"userBranches"`
    SessionID            string         `json:"sessionId"`
    IsGoldenPearl        string         `json:"isGoldenPearl"`
    IsKidsExplorer       string         `json:"isKidsExplorer"`
    Address1             string         `json:"address1"`
    Address2             string         `json:"address2"`
    Address3             string         `json:"address3"`
    Country              string         `json:"country"`
    Postcode             string         `json:"postalCode"`
    Citystate            string         `json:"cityState"`
    SignInType           int32          `json:"signInType"`
    DocNoSignup          string         `json:"docNoSignUp"`
    FullnameSignup       string         `json:"fullNameSignUp"`
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

type DbAppVersion struct {
    LatestVersion sql.NullString `db:"LATEST_VERSION"`
    OSPlatform    sql.NullString `db:"OS_PLATFORM"`
    Status        sql.NullInt64  `db:"STATUS"`
}

type AppVersion struct {
    LatestVersion string `json:"latestVersion"`
    OSPlatform    string `json:"osPlatform"`
    Status        int64  `json:"status"`
}

func (o *AppVersion) FromDbModel(m DbAppVersion) {
    o.LatestVersion = m.LatestVersion.String
    o.OSPlatform = m.OSPlatform.String
    o.Status = m.Status.Int64
}

type DbReleaseVersion struct {
    LatestVersion sql.NullString `db:"LATEST_VERSION"`
    StackPlatform sql.NullString `db:"STACK_PLATFORM"`
    DateUpdate    sql.NullString `db:"DATE_UPDATE"`
}

type ReleaseVersion struct {
    LatestVersion string `json:"latestVersion"`
    StackPlatform string `json:"stackPlatform"`
    DateUpdate    string `json:"dateUpdate"`
}

func (o *ReleaseVersion) FromDbModel(m DbReleaseVersion) {
    o.LatestVersion = m.LatestVersion.String
    o.StackPlatform = m.StackPlatform.String
    o.DateUpdate = m.DateUpdate.String
}

type DbAppServices struct {
    ServiceName  sql.NullString `db:"SERVICE_NAME"`
    ServiceImage sql.NullString `db:"SERVICE_IMAGE"`
}

type AppServices struct {
    ServiceName  string `json:"serviceName"`
    ServiceImage string `json:"serviceImage"`
}

func (o *AppServices) FromDbModel(m DbAppServices) {
    o.ServiceName = m.ServiceName.String
    o.ServiceImage = m.ServiceImage.String
}

type DbCountryTelCode struct {
    CountryName sql.NullString `db:"COUNTRY_NAME"`
    TelCode     sql.NullString `db:"TEL_CODE"`
}

type CountryTelCode struct {
    CountryName string `json:"countryName"`
    TelCode     string `json:"telCode"`
}

func (o *CountryTelCode) FromDbModel(m DbCountryTelCode) {
    o.CountryName = m.CountryName.String
    o.TelCode = m.TelCode.String
}

type DbCountry struct {
    CountryName sql.NullString `db:"COUNTRY_NAME"`
    TelCode     sql.NullString `db:"TEL_CODE"`
    CountryCode sql.NullString `db:"COUNTRY_CODE"`
}

type Country struct {
    CountryName string `json:"countryName"`
    TelCode     string `json:"telCode"`
    CountryCode string `json:"countryCode"`
}

func (o *Country) FromDbModel(m DbCountry) {
    o.CountryName = m.CountryName.String
    o.TelCode = m.TelCode.String
    o.CountryCode = m.CountryCode.String
}

type DbNationality struct {
    Nationality sql.NullString `db:"NATIONALITY"`
}

type Nationality struct {
    Nationality string `json:"nationality"`
}

func (o *Nationality) FromDbModel(m DbNationality) {
    o.Nationality = m.Nationality.String
}

type DbBranch struct {
    BranchID   sql.NullInt64  `db:"BRANCH_ID"`
    Url        sql.NullString `db:"URL"`
    Passcode   sql.NullString `db:"PASSCODE"`
    BranchName sql.NullString `db:"BRANCH_NAME"`
}

type Branch struct {
    BranchId   int64  `json:"branchId"`
    Url        string `json:"url"`
    Passcode   string `json:"passcode"`
    BranchName string `json:"branchName"`
}

func (o *Branch) FromDbModel(m DbBranch) {
    o.BranchId = m.BranchID.Int64
    o.Url = m.Url.String
    o.Passcode = m.Passcode.String
    o.BranchName = m.BranchName.String
}

type DbAssignBranch struct {
    AssignBranchID sql.NullInt64  `db:"ASSIGN_BRANCH_ID"`
    Prn            sql.NullString `db:"PRN"`
    UserID         sql.NullInt64  `db:"USER_ID"`
    AdminID        sql.NullInt64  `db:"ADMIN_ID"`
    BranchID       sql.NullInt64  `db:"BRANCH_ID"`
}

type AssignBranch struct {
    AssignBranchID int64  `json:"assignBranchId"`
    Prn            string `json:"prn"`
    UserID         int64  `json:"userId"`
    AdminID        int64  `json:"adminId"`
    BranchID       int64  `json:"branchId"`
    Branch         Branch `json:"branch"`
}

func (o *AssignBranch) FromDbModel(m DbAssignBranch) {
    o.AssignBranchID = m.AssignBranchID.Int64
    o.Prn = m.Prn.String
    o.UserID = m.UserID.Int64
    o.AdminID = m.AdminID.Int64
    o.BranchID = m.BranchID.Int64
}

type DbAdminUser struct {
    AdminID       sql.NullInt64  `db:"ADMIN_ID"`
    Address       sql.NullString `db:"ADDRESS"`
    ContactNumber sql.NullString `db:"CONTACT_NUMBER"`
    Dob           sql.NullString `db:"DOB"`
    Email         sql.NullString `db:"EMAIL"`
    FirstName     sql.NullString `db:"FIRST_NAME"`
    LastName      sql.NullString `db:"LAST_NAME"`
    MiddleName    sql.NullString `db:"MIDDLE_NAME"`
    Nationality   sql.NullString `db:"NATIONALITY"`
    Passport      sql.NullString `db:"PASSPORT"`
    Password      sql.NullString `db:"PASSWORD"`
    Resident      sql.NullString `db:"RESIDENT"`
    Role          sql.NullString `db:"ROLE"`
    Sex           sql.NullString `db:"SEX"`
    Title         sql.NullString `db:"TITLE"`
    UserGroupID   sql.NullInt64  `db:"USER_GROUP_ID"`
    UserGroupName sql.NullString `db:"USER_GROUP_NAME"`
    Username      sql.NullString `db:"USERNAME"`
}

type AdminUser struct {
    AdminID       int64  `json:"admin_id"`
    Address       string `json:"address"`
    ContactNumber string `json:"contactNumber"`
    Dob           string `json:"dob"`
    Email         string `json:"email"`
    FirstName     string `json:"firstName"`
    LastName      string `json:"lastName"`
    MiddleName    string `json:"middleName"`
    Nationality   string `json:"nationality"`
    Passport      string `json:"passport"`
    Password      string `json:"password"`
    Resident      string `json:"resident"`
    Role          string `json:"role"`
    Sex           string `json:"sex"`
    Title         string `json:"title"`
    UserGroupID   int64  `json:"userGroupId"`
    UserGroupName string `json:"userGroupName"`
    Username      string `json:"username"`
}

func (o *AdminUser) FromDbModel(m DbAdminUser) {
    o.AdminID = m.AdminID.Int64
    o.Username = m.Username.String
    o.Email = m.Email.String
    o.Password = m.Password.String
    o.Title = m.Title.String
    o.FirstName = m.FirstName.String
    o.MiddleName = m.MiddleName.String
    o.LastName = m.LastName.String
    o.Resident = m.Resident.String
    o.Dob = m.Dob.String
    o.Sex = m.Sex.String
    o.Address = m.Address.String
    o.ContactNumber = m.ContactNumber.String
    o.Passport = m.Passport.String
    o.Nationality = m.Nationality.String
    o.Role = m.Role.String
    o.UserGroupID = m.UserGroupID.Int64
    o.UserGroupName = m.UserGroupName.String
}

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
