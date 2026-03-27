package dto

type UserBillingPaymentStatusDto struct {
    Status string `json:"status" validate:"required"`
}

type UserBillingPaymentDto struct {
    BillingFullname    string `json:"billingFullname" validate:"required"`
    BillingAddress1    string `json:"billingAddress1" validate:"required"`
    BillingAddress2    string `json:"billingAddress2"`
    BillingAddress3    string `json:"billingAddress3"`
    BillingTowncity    string `json:"billingTowncity"`
    BillingState       string `json:"billingState"`
    BillingPostcode    string `json:"billingPostcode"`
    BillingCountryCode string `json:"billingCountryCode" validate:"required"`
    BillingContactNo   string `json:"billingContactNo" validate:"required"`
    BillingContactCode string `json:"billingContactCode" validate:"required"`
    BillingEmail       string `json:"billingEmail" validate:"required,email"`
}

type UserBillingDto struct {
    PatientPrn           string `json:"patientPrn" validate:"required"`
    PatientName          string `json:"patientName" validate:"required"`
    PatientDocNumber     string `json:"patientDocNumber" validate:"required"`
    VesRegDateTime       string `json:"vesRegDateTime" validate:"required"`
    VesBillNumber        string `json:"vesBillNumber" validate:"required"`
    VesInvoiceNumber     string `json:"vesInvoiceNumber" validate:"required"`
    VesInvoiceDateTime   string `json:"vesInvoiceDateTime" validate:"required"`
    VesInvoiceAmount     string `json:"vesInvoiceAmount" validate:"required"`
    VesOutstandingAmount string `json:"vesOutstandingAmount" validate:"required"`
}

type CreateUserBillingDto struct {
    UserBilling        UserBillingDto        `json:"userBilling" validate:"required"`
    UserBillingPayment UserBillingPaymentDto `json:"userBillingPayment" validate:"required"`
}
