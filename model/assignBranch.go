package model

import (
    "github.com/guregu/null/v6"
)

type AssignBranch struct {
    AssignBranchID null.Int64  `json:"assignBranchId" db:"ASSIGN_BRANCH_ID" swaggertype:"integer"`
    Prn            null.String `json:"prn" db:"PRN" swaggertype:"string"`
    UserID         null.Int64  `json:"userId" db:"USER_ID" swaggertype:"integer"`
    AdminID        null.Int64  `json:"adminId" db:"ADMIN_ID" swaggertype:"integer"`
    BranchID       null.Int64  `json:"branchId" db:"BRANCH_ID" swaggertype:"integer"`
    Branch         Branch      `json:"branch"`
}
