package model

import (
    "github.com/guregu/null/v6"
)

type AssignBranch struct {
    AssignBranchID null.Int64  `json:"assignBranchId" db:"ASSIGN_BRANCH_ID"`
    Prn            null.String `json:"prn" db:"PRN"`
    UserID         null.Int64  `json:"userId" db:"USER_ID"`
    AdminID        null.Int64  `json:"adminId" db:"ADMIN_ID"`
    BranchID       null.Int64  `json:"branchId" db:"BRANCH_ID"`
    Branch         Branch      `json:"branch"`
}
