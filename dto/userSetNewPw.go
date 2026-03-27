package dto

type PostUserSetNewPasswordDto struct {
    Email            string `json:"email" validate:"required,email"`
    VerificationCode string `json:"verificationCode" validate:"required"`
    Password         string `json:"password" validate:"required"`
}
