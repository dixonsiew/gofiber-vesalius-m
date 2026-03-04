package model

import (
    "github.com/guregu/null/v6"
)

type Country struct {
    CountryName null.String `json:"countryName" db:"COUNTRY_NAME" swaggertype:"string"`
    TelCode     null.String `json:"telCode" db:"TEL_CODE" swaggertype:"string"`
    CountryCode null.String `json:"countryCode" db:"COUNTRY_CODE" swaggertype:"string"`
}
