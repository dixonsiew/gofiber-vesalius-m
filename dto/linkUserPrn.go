package dto

type PostLinkUserPrnDto struct {
    Email         string `json:"email" validate:"required,email"`
    Prn           string `json:"prn" validate:"required"`
    BranchId      int    `json:"branchId" validate:"required,numeric"`
    Title         string `json:"title" validate:"required"`
    FirstName     string `json:"first_name" validate:"required"`
    MiddleName    string `json:"middle_name"`
    LastName      string `json:"last_name"`
    Resident      string `json:"resident" validate:"required"`
    Dob           string `json:"dob" validate:"required"`
    Sex           string `json:"sex" validate:"required"`
    Address       string `json:"address"`
    ContactNumber string `json:"contact_number"`
    Passport      string `json:"passport"`
    Nationality   string `json:"nationality"`
}
