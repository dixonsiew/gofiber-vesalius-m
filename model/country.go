package model

import (
    "database/sql"
)

type DbCountry struct {
    CountryName sql.NullString `db:"COUNTRY_NAME"`
    TelCode     sql.NullString `db:"TEL_CODE"`
    CountryCode sql.NullString `db:"COUNTRY_CODE"`
}

type Country struct {
    CountryName string `json:"countryName"`
    TelCode     string `json:"telCode"`
    CountryCode string `json:"countryCode"`
}

func (o *Country) FromDbModel(m DbCountry) {
    o.CountryName = m.CountryName.String
    o.TelCode = m.TelCode.String
    o.CountryCode = m.CountryCode.String
}
