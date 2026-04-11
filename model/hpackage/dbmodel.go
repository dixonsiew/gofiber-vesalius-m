package hpackage

import (
	"github.com/guregu/null/v6"
)

type Package struct {
    PackageId                 null.Int64  `json:"package_id" db:"PACKAGE_ID" swaggertype:"integer"`
    PackageCode               null.String `json:"packageCode" db:"PACKAGE_CODE" swaggertype:"string"`
    PackageName               null.String `json:"packageName" db:"PACKAGE_NAME" swaggertype:"string"`
    PackageDesc               null.String `json:"packageDesc" db:"PACKAGE_DESC" swaggertype:"string"`
    PackageImage              null.String `json:"packageImage" db:"PACKAGE_IMG" swaggertype:"string"`
    ResizePackageImage        string      `json:"resizePackageImage"`
    PackageStartDateTime      null.String `json:"packageStartDateTime" db:"PACKAGE_START_DATETIME" swaggertype:"string"`
    PackageEndDateTime        null.String `json:"packageEndDateTime" db:"PACKAGE_END_DATETIME" swaggertype:"string"`
    PackageValidityType       null.String `json:"packageValidityType" db:"PACKAGE_VALIDITY_TYPE" swaggertype:"string"`
    PackageValidity           null.Int64  `json:"packageValidity" db:"PACKAGE_VALIDITY" swaggertype:"integer"`
    PackageValidityDate       null.String `json:"packageValidityDate" db:"PACKAGE_VALIDITY_DATETIME" swaggertype:"string"`
    PackageTnc                null.String `json:"packageTnc" db:"PACKAGE_TNC" swaggertype:"string"`
    PackagePrice              null.Float  `json:"packagePrice" db:"PACKAGE_PRICE" swaggertype:"number"`
    PackageMaxPurchase        null.Int32  `json:"packageMaxPurchase" db:"PACKAGE_MAX_PURCHASE" swaggertype:"integer"`
    PackageAssignedDoctor     null.Int64  `json:"packageAssignedDoctor" db:"PACKAGE_ASSIGNED_DOCTOR" swaggertype:"integer"`
    PackageAssignedDoctorName string      `json:"packageAssignedDoctorName"`
    PackageAllowAppt          null.String `json:"packageAllowAppt" db:"PACKAGE_ALLOW_APPT" swaggertype:"string"`
    PackageDisplayOrder       null.Int32  `json:"packageDisplayOrder" db:"PACKAGE_DISPLAY_ORDER" swaggertype:"integer"`
    PackageExtLink            null.String `json:"packageExtLink" db:"PACKAGE_EXT_LINK" swaggertype:"string"`
    PackageTotalSold          null.Int32  `json:"packageTotalSold" db:"PACKAGE_TOTAL_SOLD" swaggertype:"integer"`
    PackageStartDateTimeExcel string      `json:"packageStartDateTimeExcel"`
    PackageEndDateTimeExcel   string      `json:"packageEndDateTimeExcel"`
    PackageValidityDateExcel  string      `json:"packageValidityDateExcel"`
    PackageValidityExcel      string      `json:"packageValidityExcel"`
}
