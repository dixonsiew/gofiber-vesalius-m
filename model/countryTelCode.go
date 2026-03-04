package model

import (
    "github.com/guregu/null/v6"
)

type CountryTelCode struct {
    CountryName null.String `json:"countryName" db:"COUNTRY_NAME" swaggertype:"string"`
    TelCode     null.String `json:"telCode" db:"TEL_CODE" swaggertype:"string"`
}
