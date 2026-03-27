package userPackage

import (
    "vesaliusm/utils"

    "github.com/guregu/null/v6"
)

type UserPackagePaymentEmail struct {
    PatientName       null.String `json:"patientName" db:"PATIENT_NAME" swaggertype:"string"`
    PaymentRequestNo  null.String `json:"paymentRequestNo" db:"PAYMENT_REQUEST_NO" swaggertype:"string"`
    PurchasedDateTime null.String `json:"purchasedDateTime" db:"PURCHASED_DATETIME" swaggertype:"string"`
    ExpiredDateTime   null.String `json:"expiredDateTime" db:"EXPIRED_DATETIME" swaggertype:"string"`
    PackageName       null.String `json:"packageName" db:"PACKAGE_NAME" swaggertype:"string"`
    PackageQuantity   null.Int32  `json:"packageQuantity" db:"PACKAGE_QUANTITY" swaggertype:"integer"`
    PackagePrice      null.Float  `json:"packagePrice" db:"PACKAGE_PRICE" swaggertype:"number"`
    PaymentGatewayV   null.String `json:"-" db:"PAYMENT_GATEWAY" swaggertype:"string"`
    PaymentGateway    string      `json:"paymentGateway"`
    BillingAddress1   null.String `json:"billingAddress1" db:"BILLING_ADDRESS1" swaggertype:"string"`
    BillingAddress2   null.String `json:"billingAddress2" db:"BILLING_ADDRESS2" swaggertype:"string"`
    BillingAddress3   null.String `json:"billingAddress3" db:"BILLING_ADDRESS3" swaggertype:"string"`
    BillingTowncity   null.String `json:"billingTowncity" db:"BILLING_TOWNCITY" swaggertype:"string"`
    BillingState      null.String `json:"billingState" db:"BILLING_STATE" swaggertype:"string"`
    BillingPostcode   null.String `json:"billingPostcode" db:"BILLING_POSTCODE" swaggertype:"string"`
    BillingEmail      null.String `json:"billingEmail" db:"BILLING_EMAIL" swaggertype:"string"`
}

func (o *UserPackagePaymentEmail) Set() {
    if o.PaymentGatewayV.String == "1" {
        o.PaymentGateway = "Wallex"
    } else {
        o.PaymentGateway = "iPay88"
    }
}

type UserPackage struct {
    PurchaseID        null.Int64  `json:"purchase_id" db:"PATIENT_PURCHASE_ID" swaggertype:"integer"`
    PatientPrn        null.String `json:"patientPrn" db:"PATIENT_PRN" swaggertype:"string"`
    PatientName       null.String `json:"patientName" db:"PATIENT_NAME" swaggertype:"string"`
    PackageID         null.Int64  `json:"package_id" db:"PACKAGE_ID" swaggertype:"integer"`
    PackagePurchaseNo null.String `json:"packagePurchaseNo" db:"PACKAGE_PURCHASE_NO" swaggertype:"string"`
    PackageStatus     null.String `json:"packageStatus" db:"PACKAGE_STATUS" swaggertype:"string"`
    OrderedDateTime   null.String `json:"orderedDateTime" db:"ORDERED_DATETIME" swaggertype:"string"`
    BookedDateTime    null.String `json:"bookedDateTime" db:"BOOKED_DATETIME" swaggertype:"string"`
    RedeemedDateTime  null.String `json:"redeemedDateTime" db:"REDEEMED_DATETIME" swaggertype:"string"`
    CancelledDateTime null.String `json:"cancelledDateTime" db:"CANCELLED_DATETIME" swaggertype:"string"`
    PurchasedDateTime null.String `json:"purchasedDateTime" db:"PURCHASED_DATETIME" swaggertype:"string"`
    ExpiredDateTime   null.String `json:"expiredDateTime" db:"EXPIRED_DATETIME" swaggertype:"string"`
    QuantityPurchased int         `json:"quantityPurchased"`

    PackageName      null.String `json:"packageName" db:"PACKAGE_NAME" swaggertype:"string"`
    PackageImage     null.String `json:"packageImage" db:"PACKAGE_IMG" swaggertype:"string"`
    PackageValidity  null.Int32  `json:"packageValidity" db:"PACKAGE_VALIDITY" swaggertype:"integer"`
    PackageAllowAppt null.String `json:"packageAllowAppt" db:"PACKAGE_ALLOW_APPT" swaggertype:"string"`

    AppointmentDateTime null.String `json:"appointmentDateTime" db:"DATE_APPT" swaggertype:"string"`
    DoctorMcr           null.String `json:"doctorMcr" db:"MCR" swaggertype:"string"`

    PaymentGatewayV        null.String `json:"-" db:"PAYMENT_GATEWAY" swaggertype:"string"`
    PaymentGateway         string      `json:"paymentGateway"`
    PaymentRequestNo       null.String `json:"paymentRequestNo" db:"PAYMENT_REQUEST_NO" swaggertype:"string"`
    PaymentRequestCurrency null.String `json:"paymentRequestCurrency" db:"PAYMENT_REQUEST_CURRENCY" swaggertype:"string"`
    PaymentAmount          null.String `json:"paymentAmount" db:"PAYMENT_AMOUNT" swaggertype:"string"`
    PaymentCurrency        null.String `json:"paymentCurrency" db:"PAYMENT_CURRENCY" swaggertype:"string"`
    PaymentAmountCollected null.String `json:"paymentAmountCollected" db:"PAYMENT_AMOUNT_COLLECTED" swaggertype:"string"`
    PaymentStatus          null.String `json:"paymentStatus" db:"PAYMENT_STATUS" swaggertype:"string"`
    PaymentTransDate       null.String `json:"paymentTransDate" db:"PAYMENT_TRANS_DATE" swaggertype:"string"`
    BillingFullname        null.String `json:"billingFullname" db:"BILLING_FULLNAME" swaggertype:"string"`
    BillingContactNo       null.String `json:"billingContactNo" db:"BILLING_CONTACT_NO" swaggertype:"string"`
    BillingContactCode     null.String `json:"billingContactCode" db:"BILLING_CONTACT_CODE" swaggertype:"string"`
    BillingEmail           null.String `json:"billingEmail" db:"BILLING_EMAIL" swaggertype:"string"`
    PaymentUrl             null.String `json:"paymentUrl" db:"PAYMENT_URL" swaggertype:"string"`

    // Excel formatted fields (optional)
    BillingFullContactExcel string `json:"billingFullContactExcel"`
    OrderedDateTimeExcel    string `json:"orderedDateTimeExcel"`
    PurchasedDateTimeExcel  string `json:"purchasedDateTimeExcel"`
    BookedDateTimeExcel     string `json:"bookedDateTimeExcel"`
    RedeemedDateTimeExcel   string `json:"redeemedDateTimeExcel"`
    CancelledDateTimeExcel  string `json:"cancelledDateTimeExcel"`
    ExpiredDateTimeExcel    string `json:"expiredDateTimeExcel"`
    PaymentTransDateExcel   string `json:"paymentTransDateExcel"`
}

func (o *UserPackage) Set() {
    if !o.OrderedDateTime.Valid {
        o.OrderedDateTime = utils.NewNullString("-")
    }

    if !o.PurchasedDateTime.Valid {
        o.PurchasedDateTime = utils.NewNullString("-")
    }

    if !o.BookedDateTime.Valid {
        o.BookedDateTime = utils.NewNullString("-")
    }

    if !o.RedeemedDateTime.Valid {
        o.RedeemedDateTime = utils.NewNullString("-")
    }

    if !o.CancelledDateTime.Valid {
        o.CancelledDateTime = utils.NewNullString("-")
    }

    if !o.ExpiredDateTime.Valid {
        o.ExpiredDateTime = utils.NewNullString("-")
    }
}

func (o *UserPackage) SetMobile() {
    if o.PackageStatus.String == "Ordered" {
        o.PackageStatus = utils.NewNullString("Awaiting Payment")
    }

    if !o.PaymentTransDate.Valid {
        o.PaymentTransDate = utils.NewNullString("-")
    }

    if !o.PaymentAmountCollected.Valid {
        o.PaymentAmountCollected = utils.NewNullString("-")
    }

    if o.PaymentGatewayV.String == "1" {
        o.PaymentGateway = "Wallex"
    } else {
        o.PaymentGateway = "iPay88"
    }

    if !o.AppointmentDateTime.Valid {
        o.AppointmentDateTime = utils.NewNullString("-")
    }

    if !o.PurchasedDateTime.Valid {
        o.PurchasedDateTime = utils.NewNullString("-")
    }

    if !o.RedeemedDateTime.Valid {
        o.RedeemedDateTime = utils.NewNullString("-")
    }

    if !o.CancelledDateTime.Valid {
        o.CancelledDateTime = utils.NewNullString("-")
    }

    if !o.ExpiredDateTime.Valid {
        o.ExpiredDateTime = utils.NewNullString("-")
    }
}

func (o *UserPackage) SetWebAdmin() {
    if !o.OrderedDateTime.Valid {
        o.OrderedDateTime = utils.NewNullString("-")
    }

    if !o.PurchasedDateTime.Valid {
        o.PurchasedDateTime = utils.NewNullString("-")
    }

    if !o.BookedDateTime.Valid {
        o.BookedDateTime = utils.NewNullString("-")
    }

    if !o.RedeemedDateTime.Valid {
        o.RedeemedDateTime = utils.NewNullString("-")
    }

    if !o.CancelledDateTime.Valid {
        o.CancelledDateTime = utils.NewNullString("-")
    }

    if !o.ExpiredDateTime.Valid {
        o.ExpiredDateTime = utils.NewNullString("-")
    }

    if !o.AppointmentDateTime.Valid {
        o.AppointmentDateTime = utils.NewNullString("-")
    }

    if !o.PaymentRequestNo.Valid {
        o.PaymentRequestNo = utils.NewNullString("-")
    }

    if o.PaymentGatewayV.String == "1" {
        o.PaymentGateway = "Wallex"
    } else {
        o.PaymentGateway = "iPay88"
    }

    if !o.PaymentTransDate.Valid {
        o.PaymentTransDate = utils.NewNullString("-")
    }

    if !o.PaymentUrl.Valid {
        o.PaymentUrl = utils.NewNullString("-")
    }
}

type PackagePaymentDetails struct {
    PaymentID              null.Int64  `json:"payment_id" db:"PACKAGE_PAYMENT_ID"`
    PaymentGateway         null.Int32  `json:"paymentGateway" db:"PAYMENT_GATEWAY"`
    PaymentRequestID       null.String `json:"paymentRequestId" db:"PAYMENT_REQUEST_ID"`
    PaymentRequestNo       null.String `json:"paymentRequestNo" db:"PAYMENT_REQUEST_NO"`
    PaymentRefID           null.String `json:"paymentRefId" db:"PAYMENT_REF_ID"`
    PaymentRequestCurrency null.String `json:"paymentRequestCurrency" db:"PAYMENT_REQUEST_CURRENCY"`
    PaymentAmount          null.String `json:"paymentAmount" db:"PAYMENT_AMOUNT"`
    PaymentPurpose         null.String `json:"paymentPurpose" db:"PAYMENT_PURPOSE"`
    PaymentCurrency        null.String `json:"paymentCurrency" db:"PAYMENT_CURRENCY"`
    PaymentAmountCollected null.String `json:"paymentAmountCollected" db:"PAYMENT_AMOUNT_COLLECTED"`
    PaymentRemarks         null.String `json:"paymentRemarks" db:"PAYMENT_REMARKS"`
    PaymentStatus          null.String `json:"paymentStatus" db:"PAYMENT_STATUS"`
    PaymentAuthCode        null.String `json:"paymentAuthCode" db:"PAYMENT_AUTH_CODE"`
    PaymentErrorDesc       null.String `json:"paymentErrorDesc" db:"PAYMENT_ERROR_DESC"`
    PaymentTransDate       null.String `json:"paymentTransDate" db:"PAYMENT_TRANS_DATE"`
    BillingFullname        null.String `json:"billingFullname" db:"BILLING_FULLNAME"`
    BillingAddress1        null.String `json:"billingAddress1" db:"BILLING_ADDRESS1"`
    BillingAddress2        null.String `json:"billingAddress2" db:"BILLING_ADDRESS2"`
    BillingAddress3        null.String `json:"billingAddress3" db:"BILLING_ADDRESS3"`
    BillingTowncity        null.String `json:"billingTowncity" db:"BILLING_TOWNCITY"`
    BillingState           null.String `json:"billingState" db:"BILLING_STATE"`
    BillingPostcode        null.String `json:"billingPostcode" db:"BILLING_POSTCODE"`
    BillingCountryCode     null.String `json:"billingCountryCode" db:"BILLING_COUNTRY_CODE"`
    BillingContactNo       null.String `json:"billingContactNo" db:"BILLING_CONTACT_NO"`
    BillingContactCode     null.String `json:"billingContactCode" db:"BILLING_CONTACT_CODE"`
    BillingEmail           null.String `json:"billingEmail" db:"BILLING_EMAIL"`
    PaymentUrl             string      `json:"paymentUrl"`
    DateCreate             string      `json:"dateCreate"`
}

type ApptDetails struct {
    PatientPrn        null.String `json:"patientPrn" db:"PATIENT_PRN"`
    PackagePurchaseNo null.String `json:"packagePurchaseNo" db:"PACKAGE_PURCHASE_NO"`
    ApptNo            null.String `json:"apptNo" db:"APPT_NO"`
}

type PackageExceedPurchaseStatus struct {
    PurchaseStatus      string `db:"PURCHASE_STATUS"`
    RecommendedQuantity int    `db:"RECOMMENDED_QUANTITY"`
}
