package dto

type ChangeSignInTypeDto struct {
    UserId              int    `json:"user_id" validate:"required,numeric"`
    SignInType          int    `json:"signInType" validate:"required,numeric"`
    SignInMobileNumber  string `json:"signInMobileNumber"`
    SignInEmailAddress  string `json:"signInEmailAddress"`
    SignInEmailPassword string `json:"signInEmailPassword"`
}
