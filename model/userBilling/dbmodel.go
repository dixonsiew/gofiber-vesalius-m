package userBilling

import (
    "github.com/guregu/null/v6"
)

type UserBillingPayment struct {
    PatientPrn             null.String `json:"patientPrn" db:"PATIENT_PRN" swaggertype:"string"`
    PatientName            null.String `json:"patientName" db:"PATIENT_NAME" swaggertype:"string"`
    PatientDocNumber       null.String `json:"patientDocNumber" db:"PATIENT_DOC_NUMBER" swaggertype:"string"`
    VesBillNumber          null.String `json:"vesBillNumber" db:"VES_BILL_NUMBER" swaggertype:"string"`
    VesInvoiceDateTime     null.String `json:"vesInvoiceDateTime" db:"INVOICE_DATE_TIME" swaggertype:"string"`
    VesInvoiceAmount       null.String `json:"vesInvoiceAmount" db:"INVOICE_AMOUNT" swaggertype:"string"`
    VesOutstandingAmount   null.String `json:"vesOutstandingAmount" db:"OUTSTANDING_AMOUNT" swaggertype:"string"`
    PaymentGatewayV        null.String `json:"-" db:"PAYMENT_GATEWAY" swaggertype:"string"`
    PaymentGateway         string      `json:"paymentGateway"`
    PaymentRequestNo       null.String `json:"paymentRequestNo" db:"PAYMENT_REQUEST_NO" swaggertype:"string"`
    PaymentAmountCollected null.String `json:"paymentAmountCollected" db:"PAYMENT_AMOUNT_COLLECTED" swaggertype:"string"`
    PaymentTransDate       null.String `json:"paymentTransDate" db:"PAYMENT_TRANS_DATE" swaggertype:"string"`
    BillingFullname        null.String `json:"billingFullname" db:"BILLING_FULLNAME" swaggertype:"string"`
    VesPaymentReceiptNo    null.String `json:"vesPaymentReceiptNo" db:"PAYMENT_RECEIPT_NO" swaggertype:"string"`
}

func(o *UserBillingPayment) Set() {
    if o.PaymentGatewayV.String == "1" {
        o.PaymentGateway = "Wallex"
    } else {
        o.PaymentGateway = "iPay88"
    }
}

type UserBilling struct {
    OutstandingBillId        int64  `json:"outstanding_bill_id"`
    OutstandingBillPaymentId int64  `json:"outstandingBillPaymentId"`
    PatientPrn               string `json:"patientPrn"`
    PatientName              string `json:"patientName"`
    PatientDocNumber         string `json:"patientDocNumber"`
    VesRegDateTime           string `json:"vesRegDateTime"`
    VesBillNumber            string `json:"vesBillNumber"`
    VesInvoiceNumber         string `json:"vesInvoiceNumber"`
    VesInvoiceDateTime       string `json:"vesInvoiceDateTime"`
    VesInvoiceAmount         string `json:"vesInvoiceAmount"`
    VesOutstandingAmount     string `json:"vesOutstandingAmount"`
    PaymentStatus            string `json:"paymentStatus"`
    SubmittedDateTime        string `json:"submittedDateTime"`
    PaidDateTime             string `json:"paidDateTime"`
}

type BillingPayment struct {
    OutstandingBillPaymentId null.Int64  `json:"outstanding_bill_payment_id" db:"OUTSTANDING_BILL_PAYMENT_ID" swaggertype:"integer"`
    PaymentGateway           null.Int32  `json:"paymentGateway" db:"PAYMENT_GATEWAY" swaggertype:"integer"`
    PaymentRequestId         null.String `json:"paymentRequestId" db:"PAYMENT_REQUEST_ID" swaggertype:"string"`
    PaymentRequestNo         null.String `json:"paymentRequestNo" db:"PAYMENT_REQUEST_NO" swaggertype:"string"`
    PaymentRefId             null.String `json:"paymentRefId" db:"PAYMENT_REF_ID" swaggertype:"string"`
    PaymentRequestCurrency   null.String `json:"paymentRequestCurrency" db:"PAYMENT_REQUEST_CURRENCY" swaggertype:"string"`
    PaymentAmount            null.Float  `json:"paymentAmount" db:"PAYMENT_AMOUNT" swaggertype:"number"`
    PaymentPurpose           null.String `json:"paymentPurpose" db:"PAYMENT_PURPOSE" swaggertype:"string"`
    PaymentCurrency          null.String `json:"paymentCurrency" db:"PAYMENT_CURRENCY" swaggertype:"string"`
    PaymentAmountCollected   null.Float  `json:"paymentAmountCollected" db:"PAYMENT_AMOUNT_COLLECTED" swaggertype:"number"`
    PaymentRemarks           null.String `json:"paymentRemarks" db:"PAYMENT_REMARKS" swaggertype:"string"`
    PaymentStatus            null.String `json:"paymentStatus" db:"PAYMENT_STATUS" swaggertype:"string"`
    PaymentAuthCode          null.String `json:"paymentAuthCode" db:"PAYMENT_AUTH_CODE" swaggertype:"string"`
    PaymentErrorDesc         null.String `json:"paymentErrorDesc" db:"PAYMENT_ERROR_DESC" swaggertype:"string"`
    PaymentTransDate         null.String `json:"paymentTransDate" db:"PAYMENT_TRANS_DATE" swaggertype:"string"`
    BillingFullname          null.String `json:"billingFullname" db:"BILLING_FULLNAME" swaggertype:"string"`
    BillingAddress1          null.String `json:"billingAddress1" db:"BILLING_ADDRESS1" swaggertype:"string"`
    BillingAddress2          null.String `json:"billingAddress2" db:"BILLING_ADDRESS2" swaggertype:"string"`
    BillingAddress3          null.String `json:"billingAddress3" db:"BILLING_ADDRESS3" swaggertype:"string"`
    BillingTowncity          null.String `json:"billingTowncity" db:"BILLING_TOWNCITY" swaggertype:"string"`
    BillingState             null.String `json:"billingState" db:"BILLING_STATE" swaggertype:"string"`
    BillingPostcode          null.String `json:"billingPostcode" db:"BILLING_POSTCODE" swaggertype:"string"`
    BillingCountryCode       null.String `json:"billingCountryCode" db:"BILLING_COUNTRY_CODE" swaggertype:"string"`
    BillingContactNo         null.String `json:"billingContactNo" db:"BILLING_CONTACT_NO" swaggertype:"string"`
    BillingContactCode       null.String `json:"billingContactCode" db:"BILLING_CONTACT_CODE" swaggertype:"string"`
    BillingEmail             null.String `json:"billingEmail" db:"BILLING_EMAIL" swaggertype:"string"`
    PaymentUrl               string      `json:"paymentUrl"`
    VesPaymentReceiptNo      string      `json:"vesPaymentReceiptNo"`
    DateCreate               string      `json:"dateCreate"`
}

type PatientOutstandingBill struct {
    PatientPRN    null.String `json:"patient_prn" db:"PATIENT_PRN" swaggertype:"string"`
    BillNumber    null.String `json:"bill_number" db:"BILL_NUMBER" swaggertype:"string"`
    InvoiceNumber null.String `json:"invoice_number" db:"INVOICE_NUMBER" swaggertype:"string"`
    PaymentStatus null.String `json:"payment_status" db:"PAYMENT_STATUS" swaggertype:"string"`
}
