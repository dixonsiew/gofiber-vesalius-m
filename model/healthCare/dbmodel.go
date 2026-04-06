package healthCare

import (
    "github.com/guregu/null/v6"
)

type Document struct {
    Code        string `json:"code"`
    Description string `json:"description"`
    Value       string `json:"value"`
    ExpireDate  string `json:"expireDate"`
}

type Address struct {
    Address1   string `json:"address1"`
    Address2   string `json:"address2"`
    Address3   string `json:"address3"`
    CityState  string `json:"cityState"`
    PostalCode string `json:"postalCode"`
    Country    string `json:"country"`
}

type Name struct {
    Title      string `json:"title"`
    FirstName  string `json:"firstName"`
    MiddleName string `json:"middleName"`
    LastName   string `json:"lastName"`
}

type ContactNumber struct {
    Home  string `json:"home"`
    Email string `json:"email"`
}

type Nationality struct {
    Code        string `json:"code"`
    Description string `json:"description"`
}

type Sex struct {
    Code        string `json:"code"`
    Description string `json:"description"`
}

type Patient struct {
    PRN           string         `json:"prn"`
    Name          *Name          `json:"name"`
    Resident      string         `json:"resident"`
    DOB           string         `json:"dob"`
    ContactNumber *ContactNumber `json:"contactNumber"`
    HomeAddress   *Address       `json:"homeAddress"`
    Sex           *Sex           `json:"sex"`
    Documents     []Document     `json:"documents"`
    Nationality   *Nationality   `json:"nationality"`
    Race          *string        `json:"race"`
    ErrorCode     *string        `json:"errorCode"`
    ErrorMessage  *string        `json:"errorMessage"`
}

type NovaPatient struct {
    PRN                null.String `json:"prn" db:"PRN" swaggertype:"string"`
    PatientName        null.String `json:"patientName" db:"PATIENT_NAME" swaggertype:"string"`
    Title              null.String `json:"title" db:"TITLE" swaggertype:"string"`
    Gender             null.String `json:"gender" db:"GENDER" swaggertype:"string"`
    DOB                null.String `json:"dob" db:"DOB" swaggertype:"string"`
    Religion           null.String `json:"religion" db:"RELIGION" swaggertype:"string"`
    Nationality        null.String `json:"nationality" db:"NATIONALITY" swaggertype:"string"`
    EthnicGroup        null.String `json:"ethnicGroup" db:"ETHNIC_GROUP" swaggertype:"string"`
    MaritalStatus      null.String `json:"maritalStatus" db:"MARITAL_STATUS" swaggertype:"string"`
    Occupation         null.String `json:"occupation" db:"OCCUPATION" swaggertype:"string"`
    HomePhone          null.String `json:"homePhone" db:"HOME_PHONE" swaggertype:"string"`
    MobilePhone        null.String `json:"mobilePhone" db:"MOBILE_PHONE" swaggertype:"string"`
    HomeAddress        null.String `json:"homeAddress" db:"HOME_ADDRESS" swaggertype:"string"`
    Email              null.String `json:"email" db:"EMAIL" swaggertype:"string"`
    SpokenLanguage     null.String `json:"spokenLanguage" db:"SPOKEN_LANGUAGE" swaggertype:"string"`
    WrittenLanguage    null.String `json:"writtenLanguage" db:"WRITTEN_LANGUAGE" swaggertype:"string"`
    CountryOfBirth     null.String `json:"countryOfBirth" db:"COUNTRY_OF_BIRTH" swaggertype:"string"`
    PlaceOfBirth       null.String `json:"placeOfBirth" db:"PLACE_OF_BIRTH" swaggertype:"string"`
    DeathDate          null.String `json:"deathDate" db:"DEATH_DATE" swaggertype:"string"`
    ChargeCategoryCode null.String `json:"chargeCategoryCode" db:"CHARGE_CATEGORY_CODE" swaggertype:"string"`
    AccessLevel        null.String `json:"accessLevel" db:"ACCESS_LEVEL" swaggertype:"string"`
    ReminderPreference null.String `json:"reminderPreference" db:"REMINDER_PREFERENCE" swaggertype:"string"`
    ResidentFlag       null.String `json:"residentFlag" db:"RESIDENT_FLAG" swaggertype:"string"`
    Remark             null.String `json:"remark" db:"REMARK" swaggertype:"string"`
    FirstOrganization  null.String `json:"firstOrganization" db:"FIRST_ORGANIZATION" swaggertype:"string"`
    EmployerCode       null.String `json:"employerCode" db:"EMPLOYER_CODE" swaggertype:"string"`
    EmployerName       null.String `json:"employerName" db:"EMPLOYER_NAME" swaggertype:"string"`
    CreationDate       null.String `json:"creationDate" db:"CREATION_DATE" swaggertype:"string"`
    CreatedBy          null.String `json:"createdBy" db:"CREATED_BY" swaggertype:"string"`
    Weight             null.String `json:"weight" db:"WEIGHT" swaggertype:"string"`
    Height             null.String `json:"height" db:"HEIGHT" swaggertype:"string"`
    PreferredDoctor    null.String `json:"preferredDoctor" db:"PREFERRED_DOCTOR" swaggertype:"string"`
    ImageURL           null.String `json:"imageUrl" db:"IMAGE_URL" swaggertype:"string"`
    SyncDate           null.String `json:"syncDate" db:"SYNC_DATE" swaggertype:"string"`
    TransferFlag       null.String `json:"transferFlag" db:"TRANSFER_FLAG" swaggertype:"string"`
    TransferDateTime   null.String `json:"transferDateTime" db:"TRANSFER_DATE_TIME" swaggertype:"string"`
    TransferSystem     null.String `json:"transferSystem" db:"TRANSFER_SYSTEM" swaggertype:"string"`
}

type NovaPatientRx struct {
    OrderNo       null.String `json:"orderNo" db:"ORDER_NO" swaggertype:"string"`
    RxRefNo       null.String `json:"rxRefNo" db:"RX_REF_NO" swaggertype:"string"`
    AccountNo     null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
    Event         null.String `json:"event" db:"EVENT" swaggertype:"string"`
    Code          null.String `json:"code" db:"CODE" swaggertype:"string"`
    Description   null.String `json:"description" db:"DESCRIPTION" swaggertype:"string"`
    EventDateTime null.String `json:"eventDateTime" db:"EVENT_DATE_TIME" swaggertype:"string"`
    Quantity      null.String `json:"quantity" db:"QUANTITY" swaggertype:"string"`
    UOM           null.String `json:"uom" db:"UOM" swaggertype:"string"`
    Instruction   null.String `json:"instruction" db:"INSTRUCTION" swaggertype:"string"`
    DoctorName    null.String `json:"doctorName" db:"DOCTOR_NAME" swaggertype:"string"`
    DoctorMCR     null.String `json:"doctorMCR" db:"DOCTOR_MCR" swaggertype:"string"`
}

type NovaReferralLetter struct {
    ReferralRefNo            null.String `json:"referralRefNo" db:"REFERRAL_REF_NO" swaggertype:"string"`
    PRN                      null.String `json:"prn" db:"PRN" swaggertype:"string"`
    AccountNo                null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
    ReferralDateTime         null.String `json:"referralDateTime" db:"REFERRAL_DATE_TIME" swaggertype:"string"`
    ReferrerDoctorMCR        null.String `json:"referrerDoctorMCR" db:"REFERRER_DOCTOR_MCR" swaggertype:"string"`
    ReferralTo               null.String `json:"referralTo" db:"REFERRAL_TO" swaggertype:"string"`
    ReferralTitleDept        null.String `json:"referralTitleDept" db:"TITLE_DEPARTMENT" swaggertype:"string"`
    ReferralAddressOrSubject null.String `json:"referralAddressOrSubject" db:"ADDRESS" swaggertype:"string"`
    ReferralLetter           null.String `json:"referralLetter" db:"REFERRAL_LETTER" swaggertype:"string"`
    ReferralType             null.String `json:"referralType" db:"REFERRAL_TYPE" swaggertype:"string"`
}

type NovaPatientInvestigation struct {
    InvestigationRefNo null.String `json:"investigationRefNo" db:"INVESTIGATION_REF_NO" swaggertype:"string"`
    InvestigationType  null.String `json:"investigationType" db:"INVESTIGATION_TYPE" swaggertype:"string"`
    AccountNo          null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
    Code               null.String `json:"code" db:"CODE" swaggertype:"string"`
    Description        null.String `json:"description" db:"DESCRIPTION" swaggertype:"string"`
    PanelCode          null.String `json:"panelCode" db:"PANEL_CODE" swaggertype:"string"`
    PanelDescription   null.String `json:"panelDescription" db:"PANEL_DESCRIPTION" swaggertype:"string"`
}

type NovaPatientInvestigationDetail struct {
    InvestigationRefNo null.String `json:"investigationRefNo" db:"INVESTIGATION_REF_NO" swaggertype:"string"`
    Code               null.String `json:"code" db:"CODE" swaggertype:"string"`
    ResultValue        null.String `json:"resultValue" db:"RESULT_VALUE" swaggertype:"string"`
    ResultUnit         null.String `json:"resultUnit" db:"RESULT_UNIT" swaggertype:"string"`
    ReferenceRange     null.String `json:"referenceRange" db:"REFERENCE_RANGE" swaggertype:"string"`
    RangeType          null.String `json:"rangeType" db:"RANGE_TYPE" swaggertype:"string"`
    ResultClob         null.String `json:"resultClob" db:"RESULT_CLOB" swaggertype:"string"`
    RecordedDate       null.String `json:"recordedDate" db:"SYNC_DATE" swaggertype:"string"`
}

type NovaLabHistoryDashboard struct {
    LabCode string                           `json:"labCode"`
    LabData []NovaPatientInvestigationDetail `json:"labData"`
}

type NovaPatientVitalSignsDetail struct {
	RefNo        null.String `json:"refNo" db:"REF_NO" swaggertype:"string"`
	Code         null.String `json:"code" db:"CODE" swaggertype:"string"`
	Description  null.String `json:"description" db:"DESCRIPTION" swaggertype:"string"`
	Value1       null.String `json:"value1" db:"VALUE1" swaggertype:"string"`
	Value2       null.String `json:"value2" db:"VALUE2" swaggertype:"string"`
	Unit         null.String `json:"unit" db:"UNIT" swaggertype:"string"`
	RecordedDate null.String `json:"recordedDate" db:"DATE_TIME" swaggertype:"string"`
}

type NovaPatientVitalSignsDetailDto struct {
	NovaPatientVitalSignsDetail *NovaPatientVitalSignsDetail `json:"novaPatientVitalSignsDetail"`
	Value1HighValue             string                       `json:"value1_high_value"`
	Value1LowValue              string                       `json:"value1_low_value"`
	Value2HighValue             string                       `json:"value2_high_value"`
	Value2LowValue              string                       `json:"value2_low_value"`
}

type NovaPatientVitalSigns struct {
	RefNo     null.String `json:"RefNo" db:"REF_NO" swaggertype:"string"`
	PRN       null.String `json:"prn" db:"PRN" swaggertype:"string"`
	AccountNo null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
	DateTime  null.String `json:"dateTime" db:"DATE_TIME" swaggertype:"string"`
}



type NovaVisit struct {
    PRN                   null.String `json:"prn" db:"PRN" swaggertype:"string"`
    AccountNo             null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
    RegistrationDate      null.String `json:"registrationDate" db:"REGISTRATION_DATE" swaggertype:"string"`
    RegistrationTime      null.String `json:"registrationTime" db:"REGISTRATION_TIME" swaggertype:"string"`
    VisitType             null.String `json:"visitType" db:"VISIT_TYPE" swaggertype:"string"`
    VisitStatus           null.String `json:"visitStatus" db:"VISIT_STATUS" swaggertype:"string"`
    ClinicName            null.String `json:"clinicName" db:"CLINIC_NAME" swaggertype:"string"`
    WardNo                null.String `json:"wardNo" db:"WARD_NO" swaggertype:"string"`
    RoomNo                null.String `json:"roomNo" db:"ROOM_NO" swaggertype:"string"`
    BedNo                 null.String `json:"bedNo" db:"BED_NO" swaggertype:"string"`
    PrimaryDoctor         null.String `json:"primaryDoctor" db:"PRIMARY_DOCTOR" swaggertype:"string"`
    AdmittingDoctor       null.String `json:"admittingDoctor" db:"ADMITTING_DOCTOR" swaggertype:"string"`
    AdmissionStatus       null.String `json:"admissionStatus" db:"ADMISSION_STATUS" swaggertype:"string"`
    AdmissionDate         null.String `json:"admissionDate" db:"ADMISSION_DATE" swaggertype:"string"`
    AdmissionTime         null.String `json:"admissionTime" db:"ADMISSION_TIME" swaggertype:"string"`
    PrelimDischargeDate   null.String `json:"prelimDischargeDate" db:"PRELIM_DISCHARGE_DATE" swaggertype:"string"`
    PrelimDischargeTime   null.String `json:"prelimDischargeTime" db:"PRELIM_DISCHARGE_TIME" swaggertype:"string"`
    DischargeDate         null.String `json:"dischargeDate" db:"DISCHARGE_DATE" swaggertype:"string"`
    DischargeTime         null.String `json:"dischargeTime" db:"DISCHARGE_TIME" swaggertype:"string"`
    ReferrerCode          null.String `json:"referrerCode" db:"REFERRER_CODE" swaggertype:"string"`
    Referral              null.String `json:"referral" db:"REFERRAL" swaggertype:"string"`
    ReferralText          null.String `json:"referralText" db:"REFERRAL_TEXT" swaggertype:"string"`
    HospitalCode          null.String `json:"hospitalCode" db:"HOSPITAL_CODE" swaggertype:"string"`
    McrNo                 null.String `json:"mcrNo" db:"MCR_NO" swaggertype:"string"`
    ChargeCategoryCode    null.String `json:"chargeCategoryCode" db:"CHARGE_CATEGORY_CODE" swaggertype:"string"`
    TransferChargeAccount null.String `json:"transferChargeAccountNo" db:"TRANSFER_CHARGE_ACCOUNT_NO" swaggertype:"string"`
    CaseType              null.String `json:"caseType" db:"CASE_TYPE" swaggertype:"string"`
    PrimarySpecialty      null.String `json:"primarySpecialty" db:"PRIMARY_SPECIALTY" swaggertype:"string"`
    HealthCheck           null.String `json:"healthCheck" db:"HEALTH_CHECK" swaggertype:"string"`
    PaymentClassCode      null.String `json:"paymentClassCode" db:"PAYMENT_CLASS_CODE" swaggertype:"string"`
    DischargeDiagnosis    null.String `json:"dischargeDiagnosis" db:"DISCHARGE_DIAGNOSIS" swaggertype:"string"`
    DischargeReason       null.String `json:"dischargeReason" db:"DISCHARGE_REASON" swaggertype:"string"`
    DischargeDoctor       null.String `json:"dischargeDoctor" db:"DISCHARGE_DOCTOR" swaggertype:"string"`
    DispositionStatus     null.String `json:"dispositionStatus" db:"DISPOSITION_STATUS" swaggertype:"string"`
    DispositionRemark     null.String `json:"dispositionRemark" db:"DISPOSITION_REMARK" swaggertype:"string"`
    CreatedBy             null.String `json:"createdBy" db:"CREATED_BY" swaggertype:"string"`
    SyncDate              null.String `json:"syncDate" db:"SYNC_DATE" swaggertype:"string"`
    TransferFlag          null.String `json:"transferFlag" db:"TRANSFER_FLG" swaggertype:"string"`
    TransferDateTime      null.String `json:"transferDateTime" db:"TRANSFER_DATE_TIME" swaggertype:"string"`
    TransferSystem        null.String `json:"transferSystem" db:"TRANSFER_SYSTEM" swaggertype:"string"`
}
