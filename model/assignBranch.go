package model

import (
    "database/sql"
)

type DbAssignBranch struct {
    AssignBranchID sql.NullInt64  `db:"ASSIGN_BRANCH_ID"`
    Prn            sql.NullString `db:"PRN"`
    UserID         sql.NullInt64  `db:"USER_ID"`
    AdminID        sql.NullInt64  `db:"ADMIN_ID"`
    BranchID       sql.NullInt64  `db:"BRANCH_ID"`
}

type AssignBranch struct {
    AssignBranchID int64  `json:"assignBranchId"`
    Prn            string `json:"prn"`
    UserID         int64  `json:"userId"`
    AdminID        int64  `json:"adminId"`
    BranchID       int64  `json:"branchId"`
    Branch         Branch `json:"branch"`
}

func (o *AssignBranch) FromDbModel(m DbAssignBranch) {
    o.AssignBranchID = m.AssignBranchID.Int64
    o.Prn = m.Prn.String
    o.UserID = m.UserID.Int64
    o.AdminID = m.AdminID.Int64
    o.BranchID = m.BranchID.Int64
}
