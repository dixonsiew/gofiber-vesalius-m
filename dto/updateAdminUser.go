package dto

type PostUpdateAdminUserDto struct {
    Title          string  `json:"title"`
    FirstName      string  `json:"first_name" validate:"required"`
    MiddleName     string  `json:"middle_name"`
    LastName       string  `json:"last_name" validate:"required"`
    Resident       string  `json:"resident"`
    Dob            string  `json:"dob"`
    Sex            string  `json:"sex"`
    Address        string  `json:"address"`
    ContactNumber  string  `json:"contact_number"`
    Passport       string  `json:"passport"`
    Nationality    string  `json:"nationality"`
    Email          string  `json:"email" validate:"required,email"`
    UserGroupID    int64   `json:"userGroupId" validate:"required"`
    UserGroupName  string  `json:"userGroupName,omitempty"`
    AdminBranchIDs []int64 `json:"adminBranchIds"`
}
