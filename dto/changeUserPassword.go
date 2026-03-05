package dto

type PostChangeUserPasswordDto struct {
    UserId      int    `json:"user_id" validate:"required,numeric"`
    NewPassword string `json:"new_password" validate:"required"`
}
