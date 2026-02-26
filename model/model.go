package model

import (
    "database/sql"
    "math"
    "strconv"
)

type PagedList struct {
    List       any `json:"list"`
    Total      int `json:"total"`
    TotalPages int `json:"totalPages"`
}

type Pager struct {
    Total    int
    PageNum  int
    PageSize int
}

func (o *Pager) SetPageSize(pageSize int) {
    if (o.Total < pageSize || pageSize < 1) && o.Total > 0 {
        o.PageSize = o.Total
    } else {
        o.PageSize = pageSize
    }

    if o.GetTotalPages() < o.PageNum {
        o.PageNum = o.GetTotalPages()
    }

    if o.PageNum < 1 {
        o.PageNum = 1
    }
}

func (o *Pager) GetLowerBound() int {
    return (o.PageNum - 1) * o.PageSize
}

func (o *Pager) GetUpperBound() int {
    x := o.PageNum * o.PageSize
    if o.Total < x {
        x = o.Total
    }

    return x
}

func (o *Pager) GetTotalPages() int {
    v := float64(o.Total) / float64(o.PageSize)
    x := math.Ceil(v)
    return int(x)
}

func GetPager(total int, page string, limit string) Pager {
    pageNum, _ := strconv.Atoi(page)
    pageSize, _ := strconv.Atoi(limit)
    pg := Pager{
        Total:    total,
        PageNum:  pageNum,
        PageSize: pageSize,
    }
    pg.SetPageSize(pageSize)
    return pg
}

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
    Username      sql.NullString `db:"USERNAME"`
    Email         sql.NullString `db:"EMAIL"`
    Password      sql.NullString `db:"PASSWORD"`
    Title         sql.NullString `db:"TITLE"`
    FirstName     sql.NullString `db:"FIRST_NAME"`
    MiddleName    sql.NullString `db:"MIDDLE_NAME"`
    LastName      sql.NullString `db:"LAST_NAME"`
    Resident      sql.NullString `db:"RESIDENT"`
    Dob           sql.NullString `db:"DOB"`
    Sex           sql.NullString `db:"SEX"`
    Address       sql.NullString `db:"ADDRESS"`
    ContactNumber sql.NullString `db:"CONTACT_NUMBER"`
    Passport      sql.NullString `db:"PASSPORT"`
    Nationality   sql.NullString `db:"NATIONALITY"`
    Role          sql.NullString `db:"ROLE"`
    UserGroupID   sql.NullInt64  `db:"USER_GROUP_ID"`
    UserGroupName sql.NullString `db:"USER_GROUP_NAME"`
    
}

type AdminUser struct {
    AdminID       int64  `json:"admin_id"`
    Username      string `json:"username"`
    Email         string `json:"email"`
    Password      string `json:"-"`
    Title         string `json:"title"`
    FirstName     string `json:"firstName"`
    MiddleName    string `json:"middleName"`
    LastName      string `json:"lastName"`
    Resident      string `json:"resident"`
    Dob           string `json:"dob"`
    Sex           string `json:"sex"`
    Address       string `json:"address"`
    ContactNumber string `json:"contactNumber"`
    Passport      string `json:"passport"`
    Nationality   string `json:"nationality"`
    Role          string `json:"role"`
    UserGroupID   int64  `json:"userGroupId"`
    UserGroupName string `json:"userGroupName"`
    AdminBranches []AssignBranch `json:"adminBranches"`
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
