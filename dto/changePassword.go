package dto

type PostChangePasswordDto struct {
    OldPassword string `json:"oldPassword" validate:"required"`
    NewPassword string `json:"newPassword" validate:"required"`
}
