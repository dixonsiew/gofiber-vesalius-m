package model

import (
    "database/sql"
)

type DbCountryTelCode struct {
    CountryName sql.NullString `db:"COUNTRY_NAME"`
    TelCode     sql.NullString `db:"TEL_CODE"`
}

type CountryTelCode struct {
    CountryName string `json:"countryName"`
    TelCode     string `json:"telCode"`
}

func (o *CountryTelCode) FromDbModel(m DbCountryTelCode) {
    o.CountryName = m.CountryName.String
    o.TelCode = m.TelCode.String
}
