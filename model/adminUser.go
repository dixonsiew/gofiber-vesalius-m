package model

import (
    "database/sql"

    "github.com/guregu/null/v6"
)

type AdminUser struct {
    AdminID       null.Int64     `json:"admin_id" db:"ADMIN_ID"`
    Username      sql.NullString `json:"username" db:"USERNAME"`
    Email         null.String    `json:"email" db:"EMAIL"`
    Password      null.String    `json:"-" db:"PASSWORD"`
    Title         null.String    `json:"title" db:"TITLE"`
    FirstName     null.String    `json:"firstName" db:"FIRST_NAME"`
    MiddleName    null.String    `json:"middleName" db:"MIDDLE_NAME"`
    LastName      null.String    `json:"lastName" db:"LAST_NAME"`
    Resident      null.String    `json:"resident" db:"RESIDENT"`
    Dob           null.String    `json:"dob" db:"DOB"`
    Sex           null.String    `json:"sex" db:"SEX"`
    Address       null.String    `json:"address" db:"ADDRESS"`
    ContactNumber null.String    `json:"contactNumber" db:"CONTACT_NUMBER"`
    Passport      null.String    `json:"passport" db:"PASSPORT"`
    Nationality   null.String    `json:"nationality" db:"NATIONALITY"`
    Role          null.String    `json:"role" db:"ROLE"`
    UserGroupID   null.Int64     `json:"userGroupId" db:"USER_GROUP_ID"`
    UserGroupName null.String    `json:"userGroupName"`
    AdminBranches []AssignBranch `json:"adminBranches"`
}
