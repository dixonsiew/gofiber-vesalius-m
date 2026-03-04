package model

import (
    "github.com/guregu/null/v6"
)

type AdminUser struct {
    AdminID       null.Int64     `json:"admin_id" db:"ADMIN_ID" swaggertype:"integer"`
    Username      null.String    `json:"username" db:"USERNAME" swaggertype:"string"`
    Email         null.String    `json:"email" db:"EMAIL" swaggertype:"string"`
    Password      null.String    `json:"-" db:"PASSWORD" swaggertype:"string"`
    Title         null.String    `json:"title" db:"TITLE" swaggertype:"string"`
    FirstName     null.String    `json:"firstName" db:"FIRST_NAME" swaggertype:"string"`
    MiddleName    null.String    `json:"middleName" db:"MIDDLE_NAME" swaggertype:"string"`
    LastName      null.String    `json:"lastName" db:"LAST_NAME" swaggertype:"string"`
    Resident      null.String    `json:"resident" db:"RESIDENT" swaggertype:"string"`
    Dob           null.String    `json:"dob" db:"DOB" swaggertype:"string"`
    Sex           null.String    `json:"sex" db:"SEX" swaggertype:"string"`
    Address       null.String    `json:"address" db:"ADDRESS" swaggertype:"string"`
    ContactNumber null.String    `json:"contactNumber" db:"CONTACT_NUMBER" swaggertype:"string"`
    Passport      null.String    `json:"passport" db:"PASSPORT" swaggertype:"string"`
    Nationality   null.String    `json:"nationality" db:"NATIONALITY" swaggertype:"string"`
    Role          null.String    `json:"role" db:"ROLE" swaggertype:"string"`
    UserGroupID   null.Int64     `json:"userGroupId" db:"USER_GROUP_ID" swaggertype:"integer"`
    UserGroupName null.String    `json:"userGroupName" swaggertype:"string"`
    AdminBranches []AssignBranch `json:"adminBranches"`
}
