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

type NovaVitalSignsDashboard struct {
    VitalSignCode  string                           `json:"vitalSignCode"`
    VitalSignsData []NovaPatientVitalSignsDetailDto `json:"vitalSignsData"`
}

type NovaPatientAlert struct {
    AlertRefNo       null.Int64  `json:"alertRefNo" db:"ALERT_REF_NO" swaggertype:"integer"`
    PRN              null.String `json:"prn" db:"PRN" swaggertype:"string"`
    AlertType        null.String `json:"alertType" db:"ALERT_TYPE" swaggertype:"string"`
    AllergyType      null.String `json:"allergyType" db:"ALLERGY_TYPE" swaggertype:"string"`
    Description      null.String `json:"description" db:"DESCRIPTION" swaggertype:"string"`
    System           null.String `json:"system" db:"SYSTEM" swaggertype:"string"`
    Route            null.String `json:"route" db:"ROUTE" swaggertype:"string"`
    Probability      null.String `json:"probability" db:"PROBABILITY" swaggertype:"string"`
    Reaction         null.String `json:"reaction" db:"REACTION" swaggertype:"string"`
    CreatedBy        null.String `json:"createdBy" db:"CREATED_BY" swaggertype:"string"`
    CreationDate     null.String `json:"creationDate" db:"CREATION_DATE" swaggertype:"string"`
    InactiveUser     null.String `json:"inactiveUser" db:"INACTIVE_USER" swaggertype:"string"`
    InactiveDateTime null.String `json:"inactiveDateTime" db:"INACTIVE_DATE_TIME" swaggertype:"string"`
    InactiveReason   null.String `json:"inactiveReason" db:"INACTIVE_REASON" swaggertype:"string"`
    SyncDate         string      `json:"syncDate"`
    TransferFlag     string      `json:"transferFlag"`
    TransferDateTime string      `json:"transferDateTime"`
    TransferSystem   string      `json:"transferSystem"`
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

type NovaBill struct {
    InvoiceNumber      null.String `json:"invoiceNumber" db:"INVOICE_NO" swaggertype:"string"`
    BillNumber         null.String `json:"billNumber" db:"BILL_NO" swaggertype:"string"`
    AccountNumber      null.String `json:"accountNumber" db:"ACCOUNT_NO" swaggertype:"string"`
    PRN                null.String `json:"prn" db:"PRN" swaggertype:"string"`
    BillDate           null.String `json:"billDate" db:"BILL_DATE" swaggertype:"string"`
    BillTime           null.String `json:"billTime" db:"BILL_TIME" swaggertype:"string"`
    BillType           null.String `json:"billType" db:"BILL_TYPE" swaggertype:"string"`
    InvoiceType        null.String `json:"invoiceType" db:"INVOICE_TYPE" swaggertype:"string"`
    Payer              null.String `json:"payer" db:"PAYER" swaggertype:"string"`
    Amount             null.Float  `json:"amount" db:"AMOUNT" swaggertype:"number"`
    OsAmount           null.Float  `json:"osAmount" db:"OS_AMOUNT" swaggertype:"number"`
    RoundingAmount     null.Float  `json:"roundingAmount" db:"ROUNDING_AMOUNT" swaggertype:"number"`
    OriginalBillNumber null.String `json:"originalBillNumber" db:"ORIGINAL_BILL_NO" swaggertype:"string"`
    MemberId           null.String `json:"memberId" db:"MEMBER_ID" swaggertype:"string"`
    SyncDate           null.String `json:"syncDate" db:"SYNC_DATE" swaggertype:"string"`
    CnAmount           null.Float  `json:"cnAmount" db:"CN_AMOUNT" swaggertype:"number"`
    DnAmount           null.Float  `json:"dnAmount" db:"DN_AMOUNT" swaggertype:"number"`
    BadDebtStatus      null.String `json:"badDebtStatus" db:"BAD_DEBT_STATUS" swaggertype:"string"`
    TransferFlag       null.String `json:"transferFlag" db:"TRANSFER_FLG" swaggertype:"string"`
    TransferDateTime   null.String `json:"transferDateTime" db:"TRANSFER_DATE_TIME" swaggertype:"string"`
    TransferSystem     null.String `json:"transferSystem" db:"TRANSFER_SYSTEM" swaggertype:"string"`
    FinalizeUser       null.String `json:"finalizeUser" db:"FINALIZE_USER" swaggertype:"string"`
    Remark             null.String `json:"remark" db:"REMARK" swaggertype:"string"`
    FinalizeCounter    null.String `json:"finalizeCounter" db:"FINALIZE_COUNTER" swaggertype:"string"`
    FinalizeLocation   null.String `json:"finalizeLocation" db:"FINALIZE_LOCATION" swaggertype:"string"`
    ExclusiveTax       null.Float  `json:"exclusiveTax" db:"EXCLUSIVE_TAX" swaggertype:"number"`
    InclusiveTax       null.Float  `json:"inclusiveTax" db:"INCLUSIVE_TAX" swaggertype:"number"`
    PaymentTerm        null.Int    `json:"paymentTerm" db:"PAYMENT_TERM" swaggertype:"integer"`
}

type NovaVisitSummary struct {
    AccountNo        null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
    Category         null.String `json:"category" db:"CATEGORY" swaggertype:"string"`
    CategoryData     null.String `json:"categoryData" db:"CATEGORY_DATA" swaggertype:"string"`
    SyncDate         null.String `json:"syncDate" db:"SYNC_DATE" swaggertype:"string"`
    TransferFlag     null.String `json:"transferFlag" db:"TRANSFER_FLG" swaggertype:"string"`
    TransferDateTime null.String `json:"transferDateTime" db:"TRANSFER_DATE_TIME" swaggertype:"string"` // note: original TS has "transferDateTIme" (typo)
    TransferSystem   null.String `json:"transferSystem" db:"TRANSFER_SYSTEM" swaggertype:"string"`
}

type NovaPatientDiagnosis struct {
    DiagnosisRefNo      null.Int    `json:"diagnosisRefNo" db:"DIAGNOSIS_REF_NO" swaggertype:"integer"`
    DoctorName          null.String `json:"doctorName" db:"DOCTOR_NAME" swaggertype:"string"`
    MCR                 null.String `json:"mcr" db:"MCR" swaggertype:"string"`
    DiagnosisDesc       null.String `json:"diagnonsisDesc" db:"DIAGNOSIS_DESC" swaggertype:"string"`
    Status              null.String `json:"status" db:"STATUS" swaggertype:"string"`
    AccountNo           null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
    DiagnosisDate       null.String `json:"diagnosisDate" db:"DIAGNOSIS_DATE" swaggertype:"string"`
    DiagnosisCloseDate  null.String `json:"diagnonsisCloseDate" db:"DIAGNOSIS_CLOSE_DATE" swaggertype:"string"`
    DiagnosisCancelDate null.String `json:"diagnosisCancelDate" db:"DIAGNOSIS_CANCEL_DATE" swaggertype:"string"`
    PrimaryDiagnosis    null.String `json:"primaryDiagnosis" db:"PRIMARY_DIAGNOSIS" swaggertype:"string"`
    SyncDate            null.String `json:"syncDate" db:"SYNC_DATE" swaggertype:"string"`
    TransferFlag        null.String `json:"transferFlag" db:"TRANSFER_FLG" swaggertype:"string"`
    TransferDateTime    null.String `json:"transferDateTime" db:"TRANSFER_DATE_TIME" swaggertype:"string"`
    TransferSystem      null.String `json:"transferSystem" db:"TRANSFER_SYSTEM" swaggertype:"string"`
}

type NovaPatientDiagnosisInfo struct {
    DiagnosisRefNo      null.Int    `json:"diagnosisRefNo" db:"DIAGNOSIS_REF_NO" swaggertype:"integer"`
    DiagnosisType       null.String `json:"diagnosisType" db:"DIAGNOSIS_TYPE" swaggertype:"string"`
    DiagnosisSeverity   null.String `json:"diagnosisSeverity" db:"DIAGNOSIS_SEVERITY" swaggertype:"string"`
    DiagnosisLaterality null.String `json:"diagnosisLaterality" db:"DIAGNOSIS_LATERALITY" swaggertype:"string"`
    SyncDate            null.String `json:"syncDate" db:"SYNC_DATE" swaggertype:"string"`
    TransferFlag        null.String `json:"transferFlag" db:"TRANSFER_FLG" swaggertype:"string"`
    TransferDateTime    null.String `json:"transferDateTime" db:"TRANSFER_DATE_TIME" swaggertype:"string"`
    TransferSystem      null.String `json:"transferSystem" db:"TRANSFER_SYSTEM" swaggertype:"string"`
}

type NovaPatientDiagnosisDetails struct {
    NovaPatientDiagnosis     *NovaPatientDiagnosis     `json:"novaPatientDiagnosis"`
    NovaPatientDiagnosisInfo *NovaPatientDiagnosisInfo `json:"novaPatientDiagnosisInfo"`
}

type NovaVisitInvestigationPanelDetail struct {
    Code           string `json:"code"`
    Description    string `json:"description"`
    ResultValue    string `json:"resultValue"`
    ResultUnit     string `json:"resultUnit"`
    ReferenceRange string `json:"referenceRange"`
    RangeType      string `json:"rangeType"`
    ResultClob     string `json:"resultClob"`
}

type NovaVisitInvestigationDetail struct {
    InvestigationRefNo string                              `json:"investigationRefNo"`
    InvestigationType  string                              `json:"investigationType"`
    AccountNo          string                              `json:"accountNo"`
    Code               string                              `json:"code"`
    Description        string                              `json:"description"`
    ResultValue        string                              `json:"resultValue"`
    ResultUnit         string                              `json:"resultUnit"`
    ReferenceRange     string                              `json:"referenceRange"`
    RangeType          string                              `json:"rangeType"`
    ResultClob         string                              `json:"resultClob"`
    PanelCode          string                              `json:"panelCode"`
    PanelDescription   string                              `json:"panelDescription"`
    PanelDetail        []NovaVisitInvestigationPanelDetail `json:"panelDetail"`
}

type NovaVisitVitalSignsDetail struct {
    Code   string `json:"code"`
    Desc   string `json:"desc"`
    Value1 string `json:"value1"`
    Value2 string `json:"value2"`
    Unit   string `json:"unit"`
}

type NovaVisitReferralLetter struct {
    ReferralRefNo            string `json:"referralRefNo"`
    ReferralType             string `json:"referralType"`
    PRN                      string `json:"prn"`
    AccountNo                string `json:"accountNo"`
    ReferralDateTime         string `json:"referralDateTime"`
    ReferrerDoctor           string `json:"referrerDoctor"`
    ReferralDoctor           string `json:"referralDoctor"`
    ReferralTitleDept        string `json:"referralTitleDept"`
    ReferralAddressOrSubject string `json:"referralAddressOrSubject"`
    ReferralLetter           string `json:"referralLetter"`
}

type NovaHealthScreeningRpt struct {
    HsrRefNo   null.String `json:"hsrRefNo" db:"HSR_REF_NO" swaggertype:"string"`
    AccountNo  null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
    ReportDate null.String `json:"reportDate" db:"REPORT_DATE" swaggertype:"string"`
    ReportUser null.String `json:"reportUser" db:"REPORT_USER" swaggertype:"string"`
}

type NovaHealthScreeningRptDetail struct {
    HsrRefNo     string `json:"hsrRefNo"`
    HsrClobValue string `json:"hsrClobValue"`
}

type NovaVisitDetails struct {
    NovaVisit                        *NovaVisit                     `json:"novaVisit"`
    NovaBills                        []NovaBill                     `json:"novaBills"`
    NovaVisitSummaries               []NovaVisitSummary             `json:"novaVisitSummaries"`
    NovaVisitPatientRxList           []NovaPatientRx                `json:"novaVisitPatientRxList"`
    NovaVisitPrescriptionList        []NovaVisitSummary             `json:"novaVisitPrescriptionList"`
    NovaPatientDiagnosisDetails      []NovaPatientDiagnosisDetails  `json:"novaPatientDiagnosisDetails"`
    NovaVisitInvestigationDetailList []NovaVisitInvestigationDetail `json:"novaVisitInvestigationDetailList"`
    NovaVisitVitalSignsDetailList    []NovaVisitVitalSignsDetail    `json:"novaVisitVitalSignsDetailList"`
    NovaVisitReferralLetterList      []NovaVisitReferralLetter      `json:"novaVisitReferralLetterList"`
    NovaHealthScreeningRptList       []NovaHealthScreeningRpt       `json:"novaHealthScreeningRptList"`
}

type Slot struct {
    SlotNumber string `json:"slotNumber"`
    Date       string `json:"date"`
    Day        string `json:"day"`
    StartTime  string `json:"startTime"`
    EndTime    string `json:"endTime"`
    DoctorName string `json:"doctorName"`
    Speciality string `json:"speciality"`
    Clinic     string `json:"clinic"`
    Room       string `json:"room"`
    CaseType   string `json:"caseType"`
}

type Appointment struct {
    AppointmentNumber string `json:"appointmentNumber"`
    Date              string `json:"date"`
    Day               string `json:"day"`
    StartTime         string `json:"startTime"`
    EndTime           string `json:"endTime"`
    DoctorMcr         string `json:"doctorMcr"`
    DoctorName        string `json:"doctorName"`
    SpecialtyCode     string `json:"specialtyCode"`
    Specialty         string `json:"specialty"`
    Clinic            string `json:"clinic"`
    Room              string `json:"room"`
    CaseType          string `json:"caseType"`
}

type PastAppointment struct {
    AppointmentRefNo  null.Int    `json:"appointmentRefNo" db:"APPOINTMENT_REF_NO" swaggertype:"integer"`
    HospitalCode      null.String `json:"hospitalCode" db:"HOSPITAL_CODE" swaggertype:"string"`
    PatientPrn        null.String `json:"patientPrn" db:"PRN" swaggertype:"string"`
    PatientName       null.String `json:"patientName" db:"NAME" swaggertype:"string"`
    PatientDob        null.String `json:"patientDob" db:"DOB" swaggertype:"string"`
    HomeContactNo     null.String `json:"homeContactNo" db:"HOME_CONTACT_NO" swaggertype:"string"`
    MobileContactNo   null.String `json:"mobileContactNo" db:"MOBILE_CONTACT_NO" swaggertype:"string"`
    HomeAddress       null.String `json:"homeAddress" db:"HOME_ADDRESS" swaggertype:"string"`
    AppointmentDate   null.String `json:"appointmentDate" db:"APPOINTMENT_DATE" swaggertype:"string"`
    AppointmentTime   null.String `json:"appointmentTime" db:"APPOINTMENT_TIME" swaggertype:"string"`
    DurationMins      null.Int    `json:"durationMins" db:"DURATION_MINS" swaggertype:"integer"`
    CaseType          null.String `json:"caseType" db:"CASE_TYPE" swaggertype:"string"`
    Status            null.String `json:"status" db:"STATUS" swaggertype:"string"`
    Reason            null.String `json:"reason" db:"REASON" swaggertype:"string"`
    Mcr               null.String `json:"mcr" db:"MCR" swaggertype:"string"`
    DoctorName        null.String `json:"doctorName" db:"DOCTOR_NAME" swaggertype:"string"`
    RoomNo            null.String `json:"roomNo" db:"ROOM_NO" swaggertype:"string"`
    ClinicName        null.String `json:"clinicName" db:"CLINIC_NAME" swaggertype:"string"`
    AccountNo         null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
    AppointmentSource null.String `json:"appointmentSource" db:"APPOINTMENT_SOURCE" swaggertype:"string"`
    SourceOfReferral  null.String `json:"sourceOfReferral" db:"SOURCE_OF_REFERRAL" swaggertype:"string"`
    AppointmentType   null.String `json:"appointmentType" db:"APPOINTMENT_TYPE" swaggertype:"string"`
    SyncDate          null.String `json:"syncDate" db:"SYNC_DATE" swaggertype:"string"`
    TransferFlg       null.String `json:"transferFlg" db:"TRANSFER_FLG" swaggertype:"string"`
    TransferDateTime  null.String `json:"transferDateTime" db:"TRANSFER_DATE_TIME" swaggertype:"string"`
    TransferSystem    null.String `json:"transferSystem" db:"TRANSFER_SYSTEM" swaggertype:"string"`
}

type NovaHealthScreeningDetailRpt struct {
    HsrRefNo     string `json:"hsrRefNo"`
    HsrClobValue string `json:"hsrClobValue"`
}
