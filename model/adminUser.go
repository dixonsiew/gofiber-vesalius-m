package model

import (
    "database/sql"
)

type DbAdminUser struct {
    AdminID       sql.NullInt64  `db:"ADMIN_ID"`
    Username      sql.NullString `db:"USERNAME"`
    Email         sql.NullString `db:"EMAIL"`
    Password      sql.NullString `db:"PASSWORD"`
    Title         sql.NullString `db:"TITLE"`
    FirstName     sql.NullString `db:"FIRST_NAME"`
    MiddleName    sql.NullString `db:"MIDDLE_NAME"`
    LastName      sql.NullString `db:"LAST_NAME"`
    Resident      sql.NullString `db:"RESIDENT"`
    Dob           sql.NullString `db:"DOB"`
    Sex           sql.NullString `db:"SEX"`
    Address       sql.NullString `db:"ADDRESS"`
    ContactNumber sql.NullString `db:"CONTACT_NUMBER"`
    Passport      sql.NullString `db:"PASSPORT"`
    Nationality   sql.NullString `db:"NATIONALITY"`
    Role          sql.NullString `db:"ROLE"`
    UserGroupID   sql.NullInt64  `db:"USER_GROUP_ID"`
    UserGroupName sql.NullString `db:"USER_GROUP_NAME"`
    
}

type AdminUser struct {
    AdminID       int64  `json:"admin_id"`
    Username      string `json:"username"`
    Email         string `json:"email"`
    Password      string `json:"-"`
    Title         string `json:"title"`
    FirstName     string `json:"firstName"`
    MiddleName    string `json:"middleName"`
    LastName      string `json:"lastName"`
    Resident      string `json:"resident"`
    Dob           string `json:"dob"`
    Sex           string `json:"sex"`
    Address       string `json:"address"`
    ContactNumber string `json:"contactNumber"`
    Passport      string `json:"passport"`
    Nationality   string `json:"nationality"`
    Role          string `json:"role"`
    UserGroupID   int64  `json:"userGroupId"`
    UserGroupName string `json:"userGroupName"`
    AdminBranches []AssignBranch `json:"adminBranches"`
}

func (o *AdminUser) FromDbModel(m DbAdminUser) {
    o.AdminID = m.AdminID.Int64
    o.Username = m.Username.String
    o.Email = m.Email.String
    o.Password = m.Password.String
    o.Title = m.Title.String
    o.FirstName = m.FirstName.String
    o.MiddleName = m.MiddleName.String
    o.LastName = m.LastName.String
    o.Resident = m.Resident.String
    o.Dob = m.Dob.String
    o.Sex = m.Sex.String
    o.Address = m.Address.String
    o.ContactNumber = m.ContactNumber.String
    o.Passport = m.Passport.String
    o.Nationality = m.Nationality.String
    o.Role = m.Role.String
    o.UserGroupID = m.UserGroupID.Int64
    o.UserGroupName = m.UserGroupName.String
}
