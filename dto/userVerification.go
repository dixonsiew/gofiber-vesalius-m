package dto

type UserEmailVerificationDto struct {
    Email            string `json:"email" validate:"required,email"`
    VerificationCode string `json:"verificationCode" validate:"required"`
}

type UserMobileVerificationDto struct {
    MobileNo string `json:"mobileNo" validate:"required"`
    TAC      string `json:"tac" validate:"required"`
}
