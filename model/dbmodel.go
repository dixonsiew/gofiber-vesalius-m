package model

import (
    "strconv"
    "vesaliusm/utils"

    "github.com/guregu/null/v6"
    "github.com/nleeper/goment"
)

type Wayfinding struct {
    MapFrom   null.String `json:"mapFrom" db:"MAP_FROM" swaggertype:"string"`
    MapTo     null.String `json:"mapTo" db:"MAP_TO" swaggertype:"string"`
    MapImage  null.String `json:"mapImage" db:"MAP_IMAGE" swaggertype:"string"`
    MapImage2 null.String `json:"mapImage2" db:"MAP_IMAGE2" swaggertype:"string"`
    MapFloor  null.String `json:"mapFloor" db:"MAP_FLOOR" swaggertype:"string"`
    MapFloor2 null.String `json:"mapFloor2" db:"MAP_FLOOR2" swaggertype:"string"`
}

type WayFindingBuildings struct {
    BuildingCode         null.String `json:"buildingCode" db:"BUILDING_CODE" swaggertype:"string"`
    BuildingName         null.String `json:"buildingName" db:"BUILDING_NAME" swaggertype:"string"`
    BuildingDisplayOrder null.Int32  `json:"buildingDisplayOrder" db:"BUILDING_DISPLAY_ORDER" swaggertype:"integer"`
}

type WayFindingFloors struct {
    FloorId           null.Int64  `json:"floorId" db:"FLOOR_DISPLAY_ORDER" swaggertype:"integer"`
    FloorDisplayOrder null.Int32  `json:"floorDisplayOrder" db:"FLOOR_DISPLAY_ORDER" swaggertype:"integer"`
    FloorCode         null.String `json:"floorCode" db:"FLOOR_CODE" swaggertype:"string"`
    FloorName         null.String `json:"floorName" db:"FLOOR_NAME" swaggertype:"string"`
    FloorImageRaw     null.String `json:"floorImageRaw" db:"FLOOR_IMAGE_RAW" swaggertype:"string"`
}

type WayFindingLocations struct {
    LocationId           null.Int64  `json:"location_id" db:"LOCATION_ID" swaggertype:"integer"`
    LocationIDMobile     null.Int64  `json:"locationId" db:"LOCATION_ID" swaggertype:"integer"`
    LocationBuilding     null.String `json:"locationBuilding" db:"LOCATION_BUILDING_CODE" swaggertype:"string"`
    LocationBuildingCode null.String `json:"locationBuildingCode" db:"LOCATION_BUILDING_CODE" swaggertype:"string"`
    LocationBuildingName null.String `json:"locationBuildingName" db:"LOCATION_BUILDING_NAME" swaggertype:"string"`
    LocationFloorCode    null.String `json:"locationFloorCode" db:"LOCATION_FLOOR_CODE" swaggertype:"string"`
    LocationFloorName    null.String `json:"locationFloorName" db:"LOCATION_FLOOR_NAME" swaggertype:"string"`
    LocationTypeCode     null.String `json:"locationTypeCode" db:"LOCATION_TYPE_CODE" swaggertype:"string"`
    LocationTypeName     null.String `json:"locationTypeName" db:"LOCATION_TYPE_NAME" swaggertype:"string"`
    LocationCode         null.String `json:"locationCode" db:"LOCATION_CODE" swaggertype:"string"`
    LocationName         null.String `json:"locationName" db:"LOCATION_NAME" swaggertype:"string"`
}

type WayFindingLocationTypes struct {
    LocationTypeId           null.Int64  `json:"locationTypeId" db:"LOCATION_TYPE_DISPLAY_ORDER" swaggertype:"integer"`
    LocationTypeDisplayOrder null.Int32  `json:"locationTypeDisplayOrder" db:"LOCATION_TYPE_DISPLAY_ORDER" swaggertype:"integer"`
    LocationTypeCode         null.String `json:"locationTypeCode" db:"LOCATION_TYPE_CODE" swaggertype:"string"`
    LocationTypeName         null.String `json:"locationTypeName" db:"LOCATION_TYPE_NAME" swaggertype:"string"`
}

type WayFindingRoutes struct {
    RouteId               null.Int64  `json:"route_id" db:"ROUTE_ID" swaggertype:"integer"`
    RouteIDMobile         null.Int64  `json:"routeId" db:"ROUTE_ID" swaggertype:"integer"`
    RouteFromLocationId   null.Int64  `json:"routeFromLocationId" db:"ROUTE_FROM_LOC_ID" swaggertype:"integer"`
    FromFloorCode         null.String `json:"fromFloorCode" db:"FROM_FLOOR_CODE" swaggertype:"string"`
    FromTypeCode          null.String `json:"fromTypeCode" db:"FROM_TYPE_CODE" swaggertype:"string"`
    FromCode              null.String `json:"fromCode" db:"FROM_CODE" swaggertype:"string"`
    FromName              null.String `json:"fromName" db:"FROM_NAME" swaggertype:"string"`
    RouteFromBuildingCode null.String `json:"routeFromBuildingCode" db:"FROM_BUILDING_CODE" swaggertype:"string"`
    RouteFromBuildingName null.String `json:"routeFromBuildingName" db:"FROM_BUILDING_NAME" swaggertype:"string"`
    RouteToLocationId     null.Int64  `json:"routeToLocationId" db:"ROUTE_TO_LOC_ID" swaggertype:"integer"`
    ToFloorCode           null.String `json:"toFloorCode" db:"TO_FLOOR_CODE" swaggertype:"string"`
    ToTypeCode            null.String `json:"toTypeCode" db:"TO_TYPE_CODE" swaggertype:"string"`
    ToCode                null.String `json:"toCode" db:"TO_CODE" swaggertype:"string"`
    ToName                null.String `json:"toName" db:"TO_NAME" swaggertype:"string"`
    RouteToBuildingCode   null.String `json:"routeToBuildingCode" db:"TO_BUILDING_CODE" swaggertype:"string"`
    RouteToBuildingName   null.String `json:"routeToBuildingName" db:"TO_BUILDING_NAME" swaggertype:"string"`
    RouteFromImageRaw     null.String `json:"routeFromImageRaw" db:"ROUTE_FROM_IMAGE_RAW" swaggertype:"string"`
    RouteToImageRaw       null.String `json:"routeToImageRaw" db:"ROUTE_TO_IMAGE_RAW" swaggertype:"string"`
}

type TempGuest struct {
    GuestPRN   null.String `json:"guestPrn" db:"GUEST_PRN" swaggertype:"string"`
    GuestName  null.String `json:"guestName" db:"GUEST_NAME" swaggertype:"string"`
    GuestEmail null.String `json:"guestEmail" db:"GUEST_EMAIL" swaggertype:"string"`
}

type HospitalProfile struct {
    ProfileDesc  null.String `json:"profileDesc" db:"PROFILE_DESC" swaggertype:"string"`
    ProfileValue null.String `json:"profileValue" db:"PROFILE_VALUE" swaggertype:"string"`
}

type ParamSetting struct {
    ParamCode  null.String `json:"paramCode" db:"PARAM_CODE" swaggertype:"string"`
    ParamDesc  null.String `json:"paramDesc" db:"PARAM_DESC" swaggertype:"string"`
    ParamValue null.Int32  `json:"paramValue" db:"PARAM_VALUE" swaggertype:"integer"`
}

type NotificationSetting struct {
    NotificationCode   null.String `json:"notificationCode" db:"NOTIFICATION_CODE" swaggertype:"string"`
    NotificationDesc   null.String `json:"notificationDesc" db:"NOTIFICATION_DESC" swaggertype:"string"`
    NotificationParam1 null.Int32  `json:"notificationParam1" db:"NOTIFICATION_PARAM_1" swaggertype:"integer"`
    NotificationParam2 null.Int32  `json:"notificationParam2" db:"NOTIFICATION_PARAM_2" swaggertype:"integer"`
}

type CronjobHistory struct {
    CronjobName           null.String `json:"cronjobName" db:"CRONJOB_NAME" swaggertype:"string"`
    CronjobNameDesc       null.String `json:"cronjobNameDesc" db:"CRONJOB_NAME_DESC" swaggertype:"string"`
    CronjobExpression     null.String `json:"cronjobExpression" db:"CRONJOB_EXPRESSION" swaggertype:"string"`
    CronjobExpressionDesc null.String `json:"cronjobExpressionDesc" db:"CRONJOB_EXPRESSION_DESC" swaggertype:"string"`
    StartDate             null.String `json:"startDate" db:"START_DATE" swaggertype:"string"`
    EndDate               null.String `json:"endDate" db:"END_DATE" swaggertype:"string"`
    CronjobStatus         null.String `json:"cronjobStatus" db:"CRONJOB_STATUS" swaggertype:"string"`
    Remarks               null.String `json:"remarks" db:"REMARKS" swaggertype:"string"`
}

type DynamicEmailMaster struct {
    EmailFunctionName null.String `json:"emailFunctionName" db:"EMAIL_FUNCTION_NAME" swaggertype:"string"`
    EmailModule       null.String `json:"emailModule" db:"EMAIL_MODULE" swaggertype:"string"`
    EmailFor          null.String `json:"emailFor" db:"EMAIL_FOR" swaggertype:"string"`
    EmailSubject      null.String `json:"emailSubject" db:"EMAIL_SUBJECT" swaggertype:"string"`
    EmailRecipient    null.String `json:"emailRecipient" db:"EMAIL_RECIPIENT" swaggertype:"string"`
    EmailSender       null.String `json:"emailSender" db:"EMAIL_SENDER" swaggertype:"string"`
    EmailSenderName   null.String `json:"emailSenderName" db:"EMAIL_SENDER_NAME" swaggertype:"string"`
    EmailTemplate     null.String `json:"emailTemplate" db:"EMAIL_TEMPLATE" swaggertype:"string"`
}

type AppVersion struct {
    LatestVersion null.String `json:"latestVersion" db:"LATEST_VERSION" swaggertype:"string"`
    OSPlatform    null.String `json:"osPlatform" db:"OS_PLATFORM" swaggertype:"string"`
    Status        null.Int64  `json:"status" db:"STATUS" swaggertype:"integer"`
}

type ReleaseVersion struct {
    LatestVersion null.String `json:"latestVersion" db:"LATEST_VERSION" swaggertype:"string"`
    StackPlatform null.String `json:"stackPlatform" db:"STACK_PLATFORM" swaggertype:"string"`
    DateUpdate    null.String `json:"dateUpdate" db:"DATE_UPDATE" swaggertype:"string"`
}

type AppServices struct {
    ServiceName  null.String `json:"serviceName" db:"SERVICE_NAME" swaggertype:"string"`
    ServiceImage null.String `json:"serviceImage" db:"SERVICE_IMAGE" swaggertype:"string"`
}

type RequestDeleteAccount struct {
    PRN            string `json:"prn"`
    Fullname       string `json:"fullname"`
    DocumentNumber string `json:"documentNumber"`
    DOB            string `json:"dob"`
    ContactNumber  string `json:"contactNumber"`
    Email          string `json:"email"`
}

type CountryTelCode struct {
    CountryName null.String `json:"countryName" db:"COUNTRY_NAME" swaggertype:"string"`
    TelCode     null.String `json:"telCode" db:"TEL_CODE" swaggertype:"string"`
}

type Country struct {
    CountryName null.String `json:"countryName" db:"COUNTRY_NAME" swaggertype:"string"`
    TelCode     null.String `json:"telCode" db:"TEL_CODE" swaggertype:"string"`
    CountryCode null.String `json:"countryCode" db:"COUNTRY_CODE" swaggertype:"string"`
}

type Nationality struct {
    Nationality null.String `json:"nationality" db:"NATIONALITY" swaggertype:"string"`
}

type Branch struct {
    BranchId   null.Int64  `json:"branchId" db:"BRANCH_ID" swaggertype:"integer"`
    BranchName null.String `json:"branchName" db:"BRANCH_NAME" swaggertype:"string"`
    Passcode   null.String `json:"passcode" db:"PASSCODE" swaggertype:"string"`
    Url        null.String `json:"url" db:"URL" swaggertype:"string"`
}

type AssignBranch struct {
    AssignBranchId null.Int64  `json:"assignBranchId" db:"ASSIGN_BRANCH_ID" swaggertype:"integer"`
    AdminId        null.Int64  `json:"adminId" db:"ADMIN_ID" swaggertype:"integer"`
    Prn            null.String `json:"prn" db:"PRN" swaggertype:"string"`
    UserId         null.Int64  `json:"userId" db:"USER_ID" swaggertype:"integer"`
    BranchId       null.Int64  `json:"branchId" db:"BRANCH_ID" swaggertype:"integer"`
    Branch         *Branch     `json:"branch"`
}

type AdminUser struct {
    AdminId       null.Int64     `json:"admin_id" db:"ADMIN_ID" swaggertype:"integer"`
    Username      null.String    `json:"username" db:"USERNAME" swaggertype:"string"`
    Email         null.String    `json:"email" db:"EMAIL" swaggertype:"string"`
    Password      null.String    `json:"-" db:"PASSWORD" swaggertype:"string"`
    Title         null.String    `json:"title" db:"TITLE" swaggertype:"string"`
    FirstName     null.String    `json:"firstName" db:"FIRST_NAME" swaggertype:"string"`
    MiddleName    null.String    `json:"middleName" db:"MIDDLE_NAME" swaggertype:"string"`
    LastName      null.String    `json:"lastName" db:"LAST_NAME" swaggertype:"string"`
    Resident      null.String    `json:"resident" db:"RESIDENT" swaggertype:"string"`
    Dob           null.String    `json:"dob" db:"DOB" swaggertype:"string"`
    Sex           null.String    `json:"sex" db:"SEX" swaggertype:"string"`
    Address       null.String    `json:"address" db:"ADDRESS" swaggertype:"string"`
    ContactNumber null.String    `json:"contactNumber" db:"CONTACT_NUMBER" swaggertype:"string"`
    Passport      null.String    `json:"passport" db:"PASSPORT" swaggertype:"string"`
    Nationality   null.String    `json:"nationality" db:"NATIONALITY" swaggertype:"string"`
    Role          null.String    `json:"role" db:"ROLE" swaggertype:"string"`
    UserGroupId   null.Int64     `json:"userGroupId" db:"USER_GROUP_ID" swaggertype:"integer"`
    UserGroupName null.String    `json:"userGroupName" db:"USER_GROUP_NAME" swaggertype:"string"`
    AdminBranches []AssignBranch `json:"adminBranches"`
}

type PatientNOK struct {
    PatientPRN      null.String `json:"patientPrn" db:"PATIENT_PRN" swaggertype:"string"`
    NOKId           null.String `json:"nokId" db:"NOK_ID" swaggertype:"string"`
    IsPatient       null.String `json:"isPatient" db:"IS_PATIENT" swaggertype:"string"`
    NOKFullname     null.String `json:"nokFullname" db:"NOK_FULLNAME" swaggertype:"string"`
    NOKRelationship null.String `json:"nokRelationship" db:"NOK_RELATIONSHIP" swaggertype:"string"`
    NOKPRN          null.String `json:"nokPrn" db:"NOK_PRN" swaggertype:"string"`
    NOKGender       null.String `json:"nokGender" db:"NOK_GENDER" swaggertype:"string"`
    NOKDocNumber    null.String `json:"nokDocNumber" db:"NOK_DOC_NUMBER" swaggertype:"string"`
    NOKDOB          null.String `json:"nokDob" db:"NOK_DOB" swaggertype:"string"`
    NOKNationality  null.String `json:"nokNationality" db:"NOK_NATIONALITY" swaggertype:"string"`
    NOKContact      null.String `json:"nokContact" db:"NOK_CONTACT" swaggertype:"string"`
    NOKAddress      null.String `json:"nokAddress" db:"NOK_ADDRESS" swaggertype:"string"`
    RefNo           null.Int64  `json:"refNo" db:"REF_NO" swaggertype:"integer"`
    MaritalStatus   null.String `json:"maritalStatus" db:"NOK_MARITAL" swaggertype:"string"`
    Email           null.String `json:"email" db:"NOK_EMAIL" swaggertype:"string"`
    NRICPassport    null.String `json:"nricPassport" swaggertype:"string"`
}

func (o *PatientNOK) Set(s string) {
    o.NRICPassport = utils.NewNullString(s)
}

type ApplicationUserFamily struct {
    AufId           null.Int64  `json:"auf_id" db:"AUF_ID" swaggertype:"integer"`
    UserId          null.Int64  `json:"user_id" db:"USER_ID" swaggertype:"integer"`
    PatientPrn      null.String `json:"patientPrn" db:"PATIENT_PRN" swaggertype:"string"`
    NokRefNumberV   null.String `json:"-" db:"NOK_REF_NUMBER" swaggertype:"string"`
    IsPatientV      null.String `json:"-" db:"IS_PATIENT" swaggertype:"string"`
    NokRefNumber    int64       `json:"nokRefNumber"`
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
        o.Relationship = utils.NewNullString("-")
    }
    if !o.NokPrn.Valid {
        o.NokPrn = utils.NewNullString("-")
    }
    if !o.NricPassport.Valid {
        o.NricPassport = utils.NewNullString("-")
    }
    if !o.Dob.Valid {
        o.Dob = utils.NewNullString("-")
    }
    if !o.Nationality.Valid {
        o.Nationality = utils.NewNullString("-")
    }
    if !o.ContactNumber.Valid {
        o.ContactNumber = utils.NewNullString("-")
    }
    if !o.Address.Valid {
        o.Address = utils.NewNullString("-")
    }
    if o.IsActiveV.String == "Y" {
        o.IsActive = true
    } else {
        o.IsActive = false
    }
    if !o.MaritalStatus.Valid {
        o.MaritalStatus = utils.NewNullString("-")
    }
    if !o.Email.Valid {
        o.Email = utils.NewNullString("-")
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

type ApplicationUser struct {
    UserId              null.Int64     `json:"user_id" db:"USER_ID" swaggertype:"integer"`
    Username            null.String    `json:"username" db:"USERNAME" swaggertype:"string"`
    Email               null.String    `json:"email" db:"EMAIL" swaggertype:"string"`
    IsKidsExplorer      null.String    `json:"isKidsExplorer" db:"IS_KIDS_EXPLORER" swaggertype:"string"`
    IsGoldenPearl       null.String    `json:"isGoldenPearl" db:"IS_GOLDEN_PEARL" swaggertype:"string"`
    Password            null.String    `json:"-" db:"PASSWORD" swaggertype:"string"`
    Title               null.String    `json:"title" db:"TITLE" swaggertype:"string"`
    FirstName           null.String    `json:"firstName" db:"FIRST_NAME" swaggertype:"string"`
    MiddleName          null.String    `json:"middleName" db:"MIDDLE_NAME" swaggertype:"string"`
    LastName            null.String    `json:"lastName" db:"LAST_NAME" swaggertype:"string"`
    FullName            null.String    `json:"fullName" db:"FIRST_NAME" swaggertype:"string"`
    Resident            null.String    `json:"resident" db:"RESIDENT" swaggertype:"string"`
    Dob                 null.String    `json:"dob" db:"DOB" swaggertype:"string"`
    Sex                 null.String    `json:"sex" db:"SEX" swaggertype:"string"`
    Race                null.String    `json:"race" db:"RACE" swaggertype:"string"`
    Address             null.String    `json:"address" db:"ADDRESS" swaggertype:"string"`
    Address1            null.String    `json:"address1" db:"ADDRESS_1" swaggertype:"string"`
    Address2            null.String    `json:"address2" db:"ADDRESS_2" swaggertype:"string"`
    Address3            null.String    `json:"address3" db:"ADDRESS_3" swaggertype:"string"`
    CityState           null.String    `json:"cityState" db:"CITYSTATE" swaggertype:"string"`
    Postcode            null.String    `json:"postalCode" db:"POSTCODE" swaggertype:"string"`
    Country             null.String    `json:"country" db:"COUNTRY" swaggertype:"string"`
    ContactNumber       null.String    `json:"contactNumber" db:"CONTACT_NUMBER" swaggertype:"string"`
    Passport            null.String    `json:"passport" db:"PASSPORT" swaggertype:"string"`
    Nationality         null.String    `json:"nationality" db:"NATIONALITY" swaggertype:"string"`
    VerificationCode    null.String    `json:"verificationCode" db:"VERIFICATION_CODE" swaggertype:"string"`
    FirstTimeLoginV     null.Int32     `json:"-" db:"FIRST_TIME_LOGIN" swaggertype:"integer"`
    FirstTimeBiometricV null.Int32     `json:"-" db:"FIRST_TIME_BIOMETRIC" swaggertype:"integer"`
    FirstTimeLogin      bool           `json:"firstTimeLogin"`
    FirstTimeBiometric  bool           `json:"firstTimeBiometric"`
    Role                null.String    `json:"role" db:"ROLE" swaggertype:"string"`
    MasterPrn           null.String    `json:"masterPrn" db:"MASTER_PRN" swaggertype:"string"`
    PlayerId            null.String    `json:"playerId" db:"PLAYER_ID" swaggertype:"string"`
    MachineId           null.String    `json:"machineId" db:"MACHINE_ID" swaggertype:"string"`
    RegisterDate        null.String    `json:"registerDate" db:"REGISTRATION_DATE_TIME" swaggertype:"string"`
    InactiveFlag        null.String    `json:"inactive" db:"INACTIVE_FLAG" swaggertype:"string"`
    IsLoggedIn          null.Int64     `json:"isLoggedIn" db:"IS_LOGGED_IN" swaggertype:"integer"`
    DateLoggedIn        null.String    `json:"dateLoggedIn" db:"DATE_LOGGED_IN" swaggertype:"string"`
    SignInType          null.Int32     `json:"signInType" db:"SIGN_IN_TYPE" swaggertype:"integer"`
    DocNoSignup         null.String    `json:"docNoSignUp" db:"DOC_NO_SIGNUP" swaggertype:"string"`
    FullnameSignup      null.String    `json:"fullNameSignUp" db:"FULLNAME_SIGNUP" swaggertype:"string"`
    UserBranches        []AssignBranch `json:"userBranches"`
    SessionId           null.String    `json:"sessionId" db:"SESSION_ID" swaggertype:"string"`
    RegisterDateExcel   string         `json:"registerDateExcel"`
}

func (o *ApplicationUserFamily) SetFromFamilyMember(m ApplicationUser) {
    o.AufId = utils.NewInt64(0)
    o.UserId = utils.NewInt64(m.UserId.Int64)
    o.NokRefNumber = 0
    o.IsPatient = true
    o.Fullname = utils.NewNullString(m.FirstName.String)
    o.Relationship = utils.NewNullString("Self")
    o.NokPrn = utils.NewNullString(m.MasterPrn.String)
    o.NricPassport = utils.NewNullString("-")
    o.DocNumber = utils.NewNullString("-")
    o.Dob = utils.NewNullString("-")
    o.Gender = utils.NewNullString(m.Sex.String)
    o.Nationality = utils.NewNullString("-")
    o.ContactNumber = utils.NewNullString("-")
    o.Address = utils.NewNullString("-")
    o.IsActive = true
    o.MaritalStatus = utils.NewNullString("-")
    o.Email = utils.NewNullString("-")
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
        o.Email = utils.NewNullString("-")
    }
    if !o.Race.Valid {
        o.Race = utils.NewNullString("-")
    }
    if !o.ContactNumber.Valid {
        o.ContactNumber = utils.NewNullString("-")
    }
    if o.FirstTimeLoginV.Int32 == 1 {
        o.FirstTimeLogin = true
    } else {
        o.FirstTimeLogin = false
    }
    if o.FirstTimeBiometricV.Int32 == 1 {
        o.FirstTimeBiometric = true
    } else {
        o.FirstTimeBiometric = false
    }

    o.RegisterDateExcel = o.RegisterDate.String
}

type DoctorPatientAppointment struct {
    DoctorPatientApptId null.Int64  `json:"doctor_patient_appt_id" db:"DOCTOR_PATIENT_APPT_ID" swaggertype:"integer"`
    DoctorId            null.Int64  `json:"doctor_id" db:"DOCTOR_ID" swaggertype:"integer"`
    DoctorName          null.String `json:"doctorName" db:"DOCTOR_NAME" swaggertype:"string"`
    DoctorSpecialty     null.String `json:"doctorSpecialty" db:"DOCTOR_SPECIALTY" swaggertype:"string"`
    ApptStatus          null.String `json:"apptStatus" db:"APPT_STATUS" swaggertype:"string"`
    ApptNo              null.String `json:"apptNo" db:"APPT_NO" swaggertype:"string"`
    ApptDay             null.String `json:"apptDay" db:"APPT_DAY" swaggertype:"string"`
    ApptSessionType     null.String `json:"apptSessionType" db:"APPT_SESSIONTYPE" swaggertype:"string"`
    ApptClinic          null.String `json:"apptClinic" db:"APPT_CLINIC" swaggertype:"string"`
    ApptRoom            null.String `json:"apptRoom" db:"APPT_ROOM" swaggertype:"string"`
    ApptCasetype        null.String `json:"apptCasetype" db:"APPT_CASETYPE" swaggertype:"string"`
    DateAppt            null.String `json:"dateAppt" db:"DATE_APPT" swaggertype:"string"`
    PatientId           null.Int64  `json:"patientId" db:"PATIENT_ID" swaggertype:"integer"`
    PatientPrn          null.String `json:"patientPrn" db:"PATIENT_PRN" swaggertype:"string"`
    PatientName         null.String `json:"patientName" db:"PATIENT_NAME" swaggertype:"string"`
}

type GroupModulePermission struct {
    PermissionId   null.Int64  `json:"permissionId" db:"PERMISSION_ID" swaggertype:"integer"`
    PermissionName null.String `json:"permissionName" db:"PERMISSION_NAME" swaggertype:"string"`
}

type UserGroupModules struct {
    ModuleId   null.Int64  `json:"moduleId" db:"MODULE_ID" swaggertype:"integer"`
    ModuleName null.String `json:"moduleName" db:"MODULE_NAME" swaggertype:"string"`
}

type UserGroupModulePermission struct {
    UserGroupModulePermissionId null.Int64 `json:"userGroupModulePermissionId" db:"USR_GRP_MOD_PERM_ID" swaggertype:"integer"`
    UserGroupId                 null.Int64 `json:"userGroupId" db:"USER_GROUP_ID" swaggertype:"integer"`
    ModuleId                    null.Int64 `json:"moduleId" db:"MODULE_ID" swaggertype:"integer"`
    PermissionId                null.Int64 `json:"permissionId" db:"PERMISSION_ID" swaggertype:"integer"`
}

type UserGroup struct {
    GroupId                             null.Int64                  `json:"groupId" db:"GROUP_ID" swaggertype:"integer"`
    UserGroupName                       null.String                 `json:"userGroupName" db:"USER_GROUP_NAME" swaggertype:"string"`
    DateCreated                         null.String                 `json:"dateCreated" db:"DATE_CREATED" swaggertype:"string"`
    UserGroupModulePermissionStatesList []UserGroupModulePermission `json:"userGroupModulePermissionStatesList"`
}

type UserGroupDetails struct {
    UserGroupId   int64                       `json:"userGroupId"`
    UserGroupName string                      `json:"userGroupName"`
    Permission    []UserGroupModulePermission `json:"permission"`
}

type AllUserGroupDetails struct {
    UserGroupId     int64                       `json:"userGroupId"`
    UserGroupName   string                      `json:"userGroupName"`
    DateCreated     string                      `json:"dateCreated"`
    SelectedModules []string                    `json:"selectedModules"`
    ActiveUser      []AdminUser                 `json:"activeUser"`
    Permission      []UserGroupModulePermission `json:"permission"`
}

type NovaDoctorSpokenLanguage struct {
    SpokenLanguageId null.Int64  `json:"spokenLanguageId" db:"SPOKEN_LANGUAGE_ID" swaggertype:"integer"`
    DoctorId         null.Int64  `json:"doctorId" db:"DOCTOR_ID" swaggertype:"integer"`
    DisplaySequence  null.Int32  `json:"displaySequence" db:"DISPLAY_SEQUENCE" swaggertype:"integer"`
    SpokenLanguage   null.String `json:"spokenLanguage" db:"SPOKEN_LANGUAGE" swaggertype:"string"`
}

type NovaDoctorQualifications struct {
    QualificationId null.Int64  `json:"qualificationId" db:"QUALIFICATION_ID" swaggertype:"integer"`
    DoctorId        null.Int64  `json:"doctorId" db:"DOCTOR_ID" swaggertype:"integer"`
    DisplaySequence null.Int32  `json:"displaySequence" db:"DISPLAY_SEQUENCE" swaggertype:"integer"`
    Qualification   null.String `json:"qualification" db:"QUALIFICATION" swaggertype:"string"`
}

type NovaDoctorSpecialities struct {
    SpecialitiesId  null.Int64  `json:"specialitiesId" db:"SPECIALITIES_ID" swaggertype:"integer"`
    DoctorId        null.Int64  `json:"doctorId" db:"DOCTOR_ID" swaggertype:"integer"`
    DisplaySequence null.Int32  `json:"displaySequence" db:"DISPLAY_SEQUENCE" swaggertype:"integer"`
    Specialities    null.String `json:"specialities" db:"SPECIALITIES" swaggertype:"string"`
    Subspecialty    null.String `json:"subspecialty" db:"SUBSPECIALTY" swaggertype:"string"`
}

type NovaDoctorClinicLocation struct {
    ClinicLocationId null.Int64  `json:"clinicLocationId" db:"CLINIC_LOCATION_ID" swaggertype:"integer"`
    DoctorId         null.Int64  `json:"doctorId" db:"DOCTOR_ID" swaggertype:"integer"`
    Location         null.String `json:"location" db:"LOCATION" swaggertype:"string"`
    Building         null.String `json:"building" db:"BUILDING" swaggertype:"string"`
}

type NovaDoctorClinicHours struct {
    ClinicHourId       null.Int64  `json:"clinicHourId" db:"CLINIC_HOUR_ID" swaggertype:"integer"`
    DoctorId           null.Int64  `json:"doctorId" db:"DOCTOR_ID" swaggertype:"integer"`
    DisplaySequence    null.Int32  `json:"displaySequence" db:"DISPLAY_SEQUENCE" swaggertype:"integer"`
    DayOfTheWeek       null.String `json:"dayOfTheWeek" db:"DAY_OF_THE_WEEK" swaggertype:"string"`
    DayStartTime       null.String `json:"dayStartTime" db:"DAY_START_TIME" swaggertype:"string"`
    DayEndTime         null.String `json:"dayEndTime" db:"DAY_END_TIME" swaggertype:"string"`
    ByAppointmentOnlyV null.Int32  `json:"-" db:"BY_APPOINTMENT_ONLY" swaggertype:"integer"`
    ByAppointmentOnly  bool        `json:"byAppointmentOnly"`
}

func (o *NovaDoctorClinicHours) Set() {
    o.ByAppointmentOnly = o.ByAppointmentOnlyV.Int32 == 1
}

type NovaDoctorAppointment struct {
    DoctorApptSlotId null.Int64  `json:"doctorApptSlotId" db:"DOCTOR_APPT_SLOT_ID" swaggertype:"integer"`
    DoctorId         null.Int64  `json:"doctorId" db:"DOCTOR_ID" swaggertype:"integer"`
    ApptDayOfWeek    null.String `json:"apptDayOfWeek" db:"DAY_OF_WEEK" swaggertype:"string"`
    ApptSlotType     null.String `json:"apptSlotType" db:"SLOT_TYPE" swaggertype:"string"`
    ApptSessionType  null.String `json:"apptSessionType" db:"SESSION_TYPE" swaggertype:"string"`
    ApptStartTime    null.String `json:"apptStartTime" db:"START_TIME" swaggertype:"string"`
    ApptEndTime      null.String `json:"apptEndTime" db:"END_TIME" swaggertype:"string"`
    ApptMaxSlots     null.Int32  `json:"apptMaxSlots" db:"MAX_SLOTS" swaggertype:"integer"`
    DisplaySequence  null.Int32  `json:"displaySequence" db:"DISPLAY_SEQUENCE" swaggertype:"integer"`
}

type NovaDoctorContact struct {
    ContactId       null.Int64  `json:"contactId" db:"CONTACT_ID" swaggertype:"integer"`
    DoctorId        null.Int64  `json:"doctorId" db:"DOCTOR_ID" swaggertype:"integer"`
    DisplaySequence null.Int32  `json:"displaySequence" db:"DISPLAY_SEQUENCE" swaggertype:"integer"`
    ContactType     null.String `json:"contactType" db:"CONTACT_TYPE" swaggertype:"string"`
    ContactValue    null.String `json:"contactValue" db:"CONTACT_VALUE" swaggertype:"string"`
}

type NovaSpecialty struct {
    SpecialtyId   null.Int64  `json:"specialtyId" db:"SPECIALTY_ID" swaggertype:"integer"`
    SpecialtyCode null.String `json:"specialtyCode" db:"SPECIALTY_CODE" swaggertype:"string"`
    SpecialtyDesc null.String `json:"specialtyDesc" db:"SPECIALTY_DESC" swaggertype:"string"`
}

type NovaDoctorSpecialty struct {
    DoctorSpecialtyId null.Int64     `json:"doctorSpecialtyId" db:"DOCTOR_SPECIALTY_ID" swaggertype:"integer"`
    DoctorId          null.Int64     `json:"doctorId" db:"DOCTOR_ID" swaggertype:"integer"`
    SpecialtyId       null.Int64     `json:"specialtyId" db:"SPECIALTY_ID" swaggertype:"integer"`
    PrimarySpecialtyV null.Int32     `json:"-" db:"PRIMARY_SPECIALTY" swaggertype:"integer"`
    PrimarySpecialty  bool           `json:"primarySpecialty"`
    Specialty         *NovaSpecialty `json:"specialty"`
}

func (o *NovaDoctorSpecialty) Set() {
    o.PrimarySpecialty = o.PrimarySpecialtyV.Int32 == 1
}

type VesaliusApptInfo struct {
    ApptNo                string `json:"apptNo"`
    ApptDate              string `json:"apptDate"`
    ApptStartTime         string `json:"apptStartTime"`
    ApptEndTime           string `json:"apptEndTime"`
    ApptCaseType          string `json:"apptCaseType"`
    ApptSlotType          string `json:"apptSlotType"`
    ApptSessionType       string `json:"apptSessionType"`
    ApptPatientPRN        string `json:"apptPatientPrn"`
    ApptPatientName       string `json:"apptPatientName"`
    ApptPackageName       string `json:"apptPackageName"`
    ApptPackageImage      string `json:"apptPackageImage"`
    ApptPackagePurchaseNo string `json:"apptPackagePurchaseNo"`
    PatientPackageExpiry  string `json:"patientPackageExpiry"`
    SessionStartTime      string `json:"sessionStartTime"`
    SessionEndTime        string `json:"sessionEndTime"`
}

type VesaliusPastApptInfo struct {
    AppointmentRefNo      int    `json:"appointmentRefNo"`
    HospitalCode          string `json:"hospitalCode"`
    PatientPRN            string `json:"patientPrn"`
    PatientName           string `json:"patientName"`
    PatientDOB            string `json:"patientDob"`
    HomeContactNo         string `json:"homeContactNo"`
    MobileContactNo       string `json:"mobileContactNo"`
    HomeAddress           string `json:"homeAddress"`
    AppointmentDate       string `json:"appointmentDate"`
    AppointmentTime       string `json:"appointmentTime"`
    DurationMins          int    `json:"durationMins"`
    CaseType              string `json:"caseType"`
    Status                string `json:"status"`
    Reason                string `json:"reason"`
    MCR                   string `json:"mcr"`
    DoctorName            string `json:"doctorName"`
    RoomNo                string `json:"roomNo"`
    ClinicName            string `json:"clinicName"`
    AccountNo             string `json:"accountNo"`
    AppointmentSource     string `json:"appointmentSource"`
    SourceOfReferral      string `json:"sourceOfReferral"`
    AppointmentType       string `json:"appointmentType"`
    SyncDate              string `json:"syncDate"`
    TransferFlg           string `json:"transferFlg"`
    TransferDateTime      string `json:"transferDateTime"`
    TransferSystem        string `json:"transferSystem"`
    ApptSlotType          string `json:"apptSlotType"`
    ApptSessionType       string `json:"apptSessionType"`
    ApptPatientPRN        string `json:"apptPatientPrn"`
    ApptPatientName       string `json:"apptPatientName"`
    ApptPackageName       string `json:"apptPackageName"`
    ApptPackageImage      string `json:"apptPackageImage"`
    ApptPackagePurchaseNo string `json:"apptPackagePurchaseNo"`
    SessionStartTime      string `json:"sessionStartTime"`
    SessionEndTime        string `json:"sessionEndTime"`
}

type PatientAppointment struct {
    DoctorId             int64                      `json:"doctor_id"`
    MCR                  string                     `json:"mcr"`
    Name                 string                     `json:"name"`
    Image                string                     `json:"image"`
    DoctorSpecialities   []NovaDoctorSpecialities   `json:"doctorSpecialities"`
    DoctorClinicLocation []NovaDoctorClinicLocation `json:"doctorClinicLocation"`
    DoctorContact        []NovaDoctorContact        `json:"doctorContact"`
    VesaliusApptInfo     VesaliusApptInfo           `json:"vesaliusApptInfo"`
}

type PatientPastAppointment struct {
    DoctorId             int64                      `json:"doctor_id"`
    MCR                  string                     `json:"mcr"`
    Name                 string                     `json:"name"`
    Image                string                     `json:"image"`
    DoctorSpecialities   []NovaDoctorSpecialities   `json:"doctorSpecialities"`
    DoctorClinicLocation []NovaDoctorClinicLocation `json:"doctorClinicLocation"`
    DoctorContact        []NovaDoctorContact        `json:"doctorContact"`
    VesaliusPastApptInfo VesaliusPastApptInfo       `json:"vesaliusPastApptInfo"`
}

type NovaDoctorAppointmentLists struct {
    CalendarDate    null.String `json:"calendarDate" db:"CALENDAR_DATE" swaggertype:"string"`
    NormalStatus    null.String `json:"normalStatus" db:"NORMAL_AVAILABILITY_STATUS" swaggertype:"string"`
    MorningStatus   null.String `json:"morningStatus" db:"MORNING_AVAILABILITY_STATUS" swaggertype:"string"`
    AfternoonStatus null.String `json:"afternoonStatus" db:"AFTERNOON_AVAILABILITY_STATUS" swaggertype:"string"`
    NightStatus     null.String `json:"nightStatus" db:"NIGHT_AVAILABILITY_STATUS" swaggertype:"string"`
    DailyStatus     null.String `json:"dailyStatus" db:"DAILY_STATUS" swaggertype:"string"`
}

type NovaDoctor struct {
    DoctorId                  null.Int64                 `json:"doctor_id" db:"DOCTOR_ID" swaggertype:"integer"`
    MCR                       null.String                `json:"mcr" db:"MCR" swaggertype:"string"`
    Name                      null.String                `json:"name" db:"NAME" swaggertype:"string"`
    Gender                    null.String                `json:"gender" db:"GENDER" swaggertype:"string"`
    Nationality               null.String                `json:"nationality" db:"NATIONALITY" swaggertype:"string"`
    DisplaySequence           null.Int32                 `json:"displaySequence" db:"DISPLAY_SEQUENCE" swaggertype:"integer"`
    AllowAppointment          null.String                `json:"allowAppointment" db:"ALLOW_APPOINTMENT" swaggertype:"string"`
    ConsultantType            null.String                `json:"consultantType" db:"CONSULTANT_TYPE" swaggertype:"string"`
    IsForPackage              null.String                `json:"isForPackage" db:"IS_FOR_PACKAGE" swaggertype:"string"`
    Qualifications            null.String                `json:"qualifications" db:"QUALIFICATIONS_SHORT" swaggertype:"string"`
    RegistrationNum           null.String                `json:"registrationNum" db:"REGISTRATION_NO" swaggertype:"string"`
    Image                     null.String                `json:"image" db:"IMAGE" swaggertype:"string"`
    ResizeImage               string                     `json:"resizeImage"`
    ShowMakeAppointmentButton string                     `json:"showMakeAppointmentButton"`
    DoctorSpokenLanguage      []NovaDoctorSpokenLanguage `json:"doctorSpokenLanguage"`
    DoctorQualifications      []NovaDoctorQualifications `json:"doctorQualifications"`
    DoctorSpecialities        []NovaDoctorSpecialities   `json:"doctorSpecialities"`
    DoctorClinicLocation      []NovaDoctorClinicLocation `json:"doctorClinicLocation"`
    DoctorClinicHours         []NovaDoctorClinicHours    `json:"doctorClinicHours"`
    DoctorAppointment         []NovaDoctorAppointment    `json:"doctorAppointment"`
    DoctorContact             []NovaDoctorContact        `json:"doctorContact"`
    DoctorSpecialty           []NovaDoctorSpecialty      `json:"doctorSpecialty"`
}

type NovaDoctorApptSlot struct {
    DoctorApptSlotId null.Int64  `json:"doctorApptSlotId" db:"DOCTOR_APPT_SLOT_ID" swaggertype:"integer"`
    DoctorId         null.Int64  `json:"doctorId" db:"DOCTOR_ID" swaggertype:"integer"`
    DayOfWeek        null.String `json:"dayOfWeek" db:"DAY_OF_WEEK" swaggertype:"string"`
    SlotType         null.String `json:"slotType" db:"SLOT_TYPE" swaggertype:"string"`
    SessionType      null.String `json:"sessionType" db:"SESSION_TYPE" swaggertype:"string"`
    StartTime        null.String `json:"startTime" db:"START_TIME" swaggertype:"string"`
    EndTime          null.String `json:"endTime" db:"END_TIME" swaggertype:"string"`
    MaxSlots         null.Int32  `json:"maxSlots" db:"MAX_SLOTS" swaggertype:"integer"`
    DisplaySequence  null.Int32  `json:"displaySequence" db:"DISPLAY_SEQUENCE" swaggertype:"integer"`
}

type ColumnMap struct {
    Field string `json:"field"`
    Text  string `json:"text"`
}

type OnesignalNotification struct {
    NotificationId          null.Int64  `json:"notification_id" db:"NOTIFICATION_ID" swaggertype:"integer"`
    UserId                  null.Int64  `json:"user_id" db:"USER_ID" swaggertype:"integer"`
    VisitType               null.String `json:"visitType" db:"VISIT_TYPE" swaggertype:"string"`
    AccountNo               null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
    NotificationTitle       null.String `json:"notificationTitle" db:"NOTIFICATION_TITLE" swaggertype:"string"`
    MsgType                 null.String `json:"msgType" db:"MSG_TYPE" swaggertype:"string"`
    ShortMessage            null.String `json:"shortMessage" db:"SHORT_MESSAGE" swaggertype:"string"`
    FullMessage             null.String `json:"fullMessage" db:"FULL_MESSAGE" swaggertype:"string"`
    MasterId                null.Int64  `json:"master_id" db:"MASTER_ID" swaggertype:"integer"`
    IsSeenV                 null.String `json:"-" db:"IS_SEEN" swaggertype:"string"`
    IsSeen                  bool        `json:"isSeen"`
    DateCreate              null.String `json:"dateCreate" db:"DATE_SENT" swaggertype:"string"`
    DateSent                string      `json:"dateSent"`
    GuestPlayerId           string      `json:"guestPlayerId"`
    DateCreate2             null.String `json:"-" db:"DATE_CREATE" swaggertype:"string"`
    DateSeen                null.String `json:"-" db:"DATE_SEEN" swaggertype:"string"`
    OneSignalMsg            null.String `json:"-" db:"ONESIGNAL_MSG" swaggertype:"string"`
    OneSignalNotificationId null.String `json:"-" db:"ONESIGNAL_NOTIFICATION_ID" swaggertype:"string"`
}

func (o *OnesignalNotification) Set() {
    if o.IsSeenV.String == "Y" {
        o.IsSeen = true
    } else {
        o.IsSeen = false
    }
}

type GeneralNotification struct {
    NotificationMasterId null.Int64  `json:"notification_master_id" db:"NOTIFICATION_MASTER_ID" swaggertype:"integer"`
    NotificationTitle    null.String `json:"notificationTitle" db:"NOTIFICATION_TITLE" swaggertype:"string"`
    ShortMessage         null.String `json:"shortMessage" db:"SHORT_MESSAGE" swaggertype:"string"`
    FullMessage          null.String `json:"fullMessage" db:"FULL_MESSAGE" swaggertype:"string"`
    StartDate            null.String `json:"startDate" db:"START_DATE_TIME" swaggertype:"string"`
    EndDate              null.String `json:"endDate" db:"END_DATE_TIME" swaggertype:"string"`
    TargetAgeFrom        null.Int64  `json:"targetAgeFrom" db:"TARGET_AGE_FROM" swaggertype:"integer"`
    TargetAgeTo          null.Int64  `json:"targetAgeTo" db:"TARGET_AGE_TO" swaggertype:"integer"`
    TargetGender         null.String `json:"targetGender" db:"TARGET_GENDER" swaggertype:"string"`
    TargetNationality    null.String `json:"targetNationality" db:"TARGET_NATIONALITY" swaggertype:"string"`
    TargetCity           null.String `json:"targetCity" db:"TARGET_CITY" swaggertype:"string"`
    TargetState          null.String `json:"targetState" db:"TARGET_STATE" swaggertype:"string"`
}

func (o *GeneralNotification) Set() {
    if o.StartDate.Valid {
        g, _ := goment.New(o.StartDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.StartDate = utils.NewNullString(g.Format("DD/MM/YYYY"))
    }

    if o.EndDate.Valid {
        g, _ := goment.New(o.EndDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.EndDate = utils.NewNullString(g.Format("DD/MM/YYYY"))
    }
}

type MobileUserAuditLog struct {
    AuditId         null.Int64  `json:"audit_id" db:"AUDIT_ID" swaggertype:"integer"`
    Prn             null.String `json:"prn" db:"PRN" swaggertype:"string"`
    Username        null.String `json:"username" db:"USERNAME" swaggertype:"string"`
    PatientName     null.String `json:"patientName" db:"PATIENT_NAME" swaggertype:"string"`
    Action          null.String `json:"action" db:"ACTION" swaggertype:"string"`
    ActionDesc      null.String `json:"actionDesc" db:"ACTION_DESC" swaggertype:"string"`
    Remarks         null.String `json:"remarks" db:"REMARKS" swaggertype:"string"`
    UserCreate      null.String `json:"userCreate" db:"USER_CREATE" swaggertype:"string"`
    DateCreate      null.String `json:"dateCreate" db:"DATE_CREATE" swaggertype:"string"`
    DateCreateExcel string      `json:"dateCreateExcel"`
}

type EmailMaster struct {
    EmailFunctionName null.String `json:"emailFunctionName" db:"EMAIL_FUNCTION_NAME" swaggertype:"string"`
    EmailModule       null.String `json:"emailModule" db:"EMAIL_MODULE" swaggertype:"string"`
    EmailFor          null.String `json:"emailFor" db:"EMAIL_FOR" swaggertype:"string"`
    EmailSubject      null.String `json:"emailSubject" db:"EMAIL_SUBJECT" swaggertype:"string"`
    EmailRecipient    null.String `json:"emailRecipient" db:"EMAIL_RECIPIENT" swaggertype:"string"`
    EmailSender       null.String `json:"emailSender" db:"EMAIL_SENDER" swaggertype:"string"`
    EmailSenderName   null.String `json:"emailSenderName" db:"EMAIL_SENDER_NAME" swaggertype:"string"`
    EmailTemplate     null.String `json:"emailTemplate" db:"EMAIL_TEMPLATE" swaggertype:"string"`
}

type AdminAuditLog struct {
    EventId          null.Int64  `json:"event_id" db:"EVENT_ID" swaggertype:"integer"`
    EventDateTime    null.String `json:"eventDateTime" db:"EVENT_DATE_TIME" swaggertype:"string"`
    EventAdminId     null.Int64  `json:"eventAdminId" db:"EVENT_ADMIN_ID" swaggertype:"integer"`
    EventAdminEmail  null.String `json:"eventAdminEmail" db:"EVENT_ADMIN_EMAIL" swaggertype:"string"`
    EventModule      null.String `json:"eventModule" db:"EVENT_MODULE" swaggertype:"string"`
    EventFunction    null.String `json:"eventFunction" db:"EVENT_FUNCTION" swaggertype:"string"`
    EventAction      null.String `json:"eventAction" db:"EVENT_ACTION" swaggertype:"string"`
    EventKeyword     null.String `json:"eventKeyword" db:"EVENT_KEYWORD" swaggertype:"string"`
    EventDescription null.String `json:"eventDescription" db:"EVENT_DESCRIPTION" swaggertype:"string"`
    PatientPRN       null.String `json:"patientPrn" db:"PATIENT_PRN" swaggertype:"string"`
}

type StatisticAppointment struct {
    Month          null.String `json:"month" db:"MONTH" swaggertype:"string"`
    Year           null.String `json:"year" db:"YEAR" swaggertype:"string"`
    TotalChanged   null.String `json:"totalChanged" db:"TOTAL_CHANGED" swaggertype:"string"`
    TotalCancelled null.String `json:"totalCancelled" db:"TOTAL_CANCELLED" swaggertype:"string"`
    TotalConfirmed null.String `json:"totalConfirmed" db:"TOTAL_CONFIRMED" swaggertype:"string"`
}

type StatisticMobileRegistration struct {
    Month              null.String `json:"month" db:"MONTH" swaggertype:"string"`
    Year               null.String `json:"year" db:"YEAR" swaggertype:"string"`
    TotalRegistrations null.String `json:"totalRegistrations" db:"TOTAL_REGISTRATIONS" swaggertype:"string"`
}

type StatisticMobileFeedback struct {
    Month          null.String `json:"month" db:"MONTH" swaggertype:"string"`
    Year           null.String `json:"year" db:"YEAR" swaggertype:"string"`
    TotalFeedbacks null.String `json:"totalFeedbacks" db:"TOTAL_FEEDBACKS" swaggertype:"string"`
}

type StatisticMobilePackage struct {
    Month     null.String `json:"month" db:"MONTH" swaggertype:"string"`
    Year      null.String `json:"year" db:"YEAR" swaggertype:"string"`
    Purchased null.String `json:"purchased" db:"PURCHASED" swaggertype:"string"`
    Redeemed  null.String `json:"redeemed" db:"REDEEMED" swaggertype:"string"`
    Expired   null.String `json:"expired" db:"EXPIRED" swaggertype:"string"`
}

type StatisticMobileClubs struct {
    Month     null.String `json:"month" db:"MONTH" swaggertype:"string"`
    Year      null.String `json:"year" db:"YEAR" swaggertype:"string"`
    TotalJoin null.String `json:"totalJoin" db:"TOTAL_JOINS" swaggertype:"string"`
}
