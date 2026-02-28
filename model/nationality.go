package model

import (
    "database/sql"
)

type DbNationality struct {
    Nationality sql.NullString `db:"NATIONALITY"`
}

type Nationality struct {
    Nationality string `json:"nationality"`
}

func (o *Nationality) FromDbModel(m DbNationality) {
    o.Nationality = m.Nationality.String
}
