package dto

type RefreshTokenDto struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
    FromAdmin    bool   `json:"fromAdmin" default:"false"`
}
