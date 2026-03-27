package dto

type PostSelfSignUpUserDto struct {
    BranchId         int    `json:"branchId" validate:"required,numeric"`
    UserFullName     string `json:"userFullName" validate:"required"`
    UserPersonNumber string `json:"userPersonNumber" validate:"required"`
    UserDOB          string `json:"userDOB" validate:"required"`
    UserEmail        string `json:"userEmail" validate:"required,email"`
    UserPassword     string `json:"userPassword" validate:"required"`
    PlayerId         string `json:"playerId" validate:"required"`
}
