package dto

type PackageDto struct {
    PackageCode           string  `json:"packageCode" validate:"required"`
    PackageName           string  `json:"packageName" validate:"required"`
    PackageDesc           string  `json:"packageDesc" validate:"required"`
    PackageImage          string  `json:"packageImage" validate:"required"`
    PackageStartDateTime  string  `json:"packageStartDateTime" validate:"required"`
    PackageEndDateTime    string  `json:"packageEndDateTime" validate:"required"`
    PackageValidityType   string  `json:"packageValidityType" validate:"required"`
    PackageValidity       int     `json:"packageValidity" validate:"required,numeric"`
    PackageValidityDate   string  `json:"packageValidityDate" validate:"required"`
    PackageTnc            string  `json:"packageTnc" validate:"required"`
    PackagePrice          float64 `json:"packagePrice" validate:"required,numeric"`
    PackageMaxPurchase    int     `json:"packageMaxPurchase" validate:"required,numeric"`
    PackageAssignedDoctor int     `json:"packageAssignedDoctor" validate:"required,numeric"`
    PackageAllowAppt      string  `json:"packageAllowAppt" validate:"required"`
    PackageDisplayOrder   int     `json:"packageDisplayOrder" validate:"required,numeric"`
    PackageExtLink        string  `json:"packageExtLink" validate:"required"`
}
