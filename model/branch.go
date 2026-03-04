package model

import (
    "github.com/guregu/null/v6"
)

type Branch struct {
    BranchId   null.Int64  `json:"branchId" db:"BRANCH_ID" swaggertype:"integer"`
    Url        null.String `json:"url" db:"URL" swaggertype:"string"`
    Passcode   null.String `json:"passcode" db:"PASSCODE" swaggertype:"string"`
    BranchName null.String `json:"branchName" db:"BRANCH_NAME" swaggertype:"string"`
}
