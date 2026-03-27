package dto

type NewSignupUserDto struct {
    BranchId         int    `json:"branchId" validate:"required,numeric"`
    UserPrn          string `json:"userPrn" validate:"required"`
    UserFullName     string `json:"userFullName" validate:"required"`
    UserPersonNumber string `json:"userPersonNumber" validate:"required"`
    UserDOB          string `json:"userDOB" validate:"required"`
    UserMobileNo     string `json:"userMobileNo"`
    UserEmail        string `json:"userEmail" validate:"email"`
    UserPassword     string `json:"userPassword"`
    PlayerId         string `json:"playerId" validate:"required"`
    SignInType       int    `json:"signInType" validate:"required,numeric"`
}
