package model

import (
    "github.com/guregu/null/v6"
)

type Nationality struct {
    Nationality null.String `json:"nationality" db:"NATIONALITY"`
}
