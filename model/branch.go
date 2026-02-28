package model

import (
    "database/sql"
)

type DbBranch struct {
    BranchID   sql.NullInt64  `db:"BRANCH_ID"`
    Url        sql.NullString `db:"URL"`
    Passcode   sql.NullString `db:"PASSCODE"`
    BranchName sql.NullString `db:"BRANCH_NAME"`
}

type Branch struct {
    BranchId   int64  `json:"branchId"`
    Url        string `json:"url"`
    Passcode   string `json:"passcode"`
    BranchName string `json:"branchName"`
}

func (o *Branch) FromDbModel(m DbBranch) {
    o.BranchId = m.BranchID.Int64
    o.Url = m.Url.String
    o.Passcode = m.Passcode.String
    o.BranchName = m.BranchName.String
}
