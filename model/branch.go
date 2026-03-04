package model

import (
    "github.com/guregu/null/v6"
)

type Branch struct {
    BranchId   null.Int64  `json:"branchId" db:"BRANCH_ID"`
    Url        null.String `json:"url" db:"URL"`
    Passcode   null.String `json:"passcode" db:"PASSCODE"`
    BranchName null.String `json:"branchName" db:"BRANCH_NAME"`
}
