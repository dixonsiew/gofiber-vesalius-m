package dto

type PostAdminUserDto struct {
    Username       string `json:"username" validate:"required"`
    Title          string `json:"title"`
    FirstName      string `json:"first_name" validate:"required"`
    MiddleName     string `json:"middle_name"`
    LastName       string `json:"last_name" validate:"required"`
    Resident       string `json:"resident"`
    Dob            string `json:"dob"`
    Sex            string `json:"sex"`
    Address        string `json:"address"`
    ContactNumber  string `json:"contact_number"`
    Passport       string `json:"passport"`
    Nationality    string `json:"nationality"`
    Email          string `json:"email" validate:"required,email"`
    UserGroupId    int    `json:"userGroupId" validate:"required,numeric"`
    AdminBranchIds []int  `json:"adminBranchIds"`
}
