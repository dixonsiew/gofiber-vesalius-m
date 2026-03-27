package dto

type DeleteAccountDto struct {
    PRN            string `json:"prn" validate:"required"`
    FullName       string `json:"fullname" validate:"required"`
    DocumentNumber string `json:"documentNumber" validate:"required"`
    DOB            string `json:"dob" validate:"required"`
    ContactNumber  string `json:"contactNumber"`
    Email          string `json:"email"`
}
