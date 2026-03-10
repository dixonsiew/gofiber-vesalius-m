package dto

type SearchPurchaseHistoryDto struct {
    Keyword  string `json:"keyword" default:""`
    Keyword2 string `json:"keyword2" default:""`
    Keyword3 string `json:"keyword3" default:""`
    Keyword4 string `json:"keyword4" default:""`
}

type CheckPackageDto struct {
    PackageID         int64 `json:"package_id" validate:"required,number"`
    QuantityPurchased int   `json:"quantityPurchased" validate:"required,number"`
}

type UserPackagePaymentDto struct {
    BillingFullname    string `json:"billingFullname" validate:"required"`
    BillingAddress1    string `json:"billingAddress1" validate:"required"`
    BillingAddress2    string `json:"billingAddress2"`
    BillingAddress3    string `json:"billingAddress3"`
    BillingTowncity    string `json:"billingTowncity"`
    BillingState       string `json:"billingState"`
    BillingPostcode    string `json:"billingPostcode"`
    BillingCountryCode string `json:"billingCountryCode" validate:"required"`
    BillingContactNo   string `json:"billingContactNo" validate:"required,max=20"`
    BillingContactCode string `json:"billingContactCode" validate:"required,email"`
    BillingEmail       string `json:"billingEmail" validate:"required,email"`
}

type GuestPackagePaymentDto struct {
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

type UserPackageDto struct {
    PatientPrn        string  `json:"patientPrn" validate:"required"`
    PatientName       string  `json:"patientName" validate:"required"`
    PackageID         int64   `json:"package_id" validate:"required,number"`
    PackageName       string  `json:"packageName" validate:"required"`
    PackagePrice      float64 `json:"packagePrice" validate:"required,number"`
    QuantityPurchased int     `json:"quantityPurchased" validate:"required,min=0,number"`
}

type GuestPackageDto struct {
    PatientPrn        string  `json:"patientPrn" validate:"required"`
    PatientName       string  `json:"patientName" validate:"required"`
    PackageID         int64   `json:"package_id" validate:"required,number"`
    PackageName       string  `json:"packageName" validate:"required"`
    PackagePrice      float64 `json:"packagePrice" validate:"required"`
    QuantityPurchased int     `json:"quantityPurchased" validate:"required,min=0,number"`
}

type CheckPackageExpiryMaxpurchaseDto struct {
    Package []CheckPackageDto `json:"package" validate:"required"`
}

type UserPackageStatusDto struct {
    Status string `json:"status" validate:"required"`
}

type CreateUserPackageDto struct {
    UserPackage        []UserPackageDto      `json:"userPackage" validate:"required"`
    UserPackagePayment UserPackagePaymentDto `json:"userPackagePayment" validate:"required"`
}

type CreateGuestPackageDto struct {
    GuestPackage        []GuestPackageDto      `json:"guestPackage" validate:"required"`
    GuestPackagePayment GuestPackagePaymentDto `json:"guestPackagePayment" validate:"required"`
}
