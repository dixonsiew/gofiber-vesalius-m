package dto

type LoginDto struct {
    Username      string `json:"username" validate:"required" default:""`
    Password      string `json:"password" validate:"required" default:""`
    PlayerId      string `json:"playerId"`
    MachineId     string `json:"machineId"`
    FromAdmin     bool   `json:"fromAdmin" default:"false"`
    FromBiometric int    `json:"fromBiometric" default:"0"`
}

type NewLoginDto struct {
    SignInType    int    `json:"signInType" validate:"required"`
    Username      string `json:"username" validate:"required" default:""`
    Password      string `json:"password" default:""`
    PlayerId      string `json:"playerId"`
    MachineId     string `json:"machineId"`
    FromBiometric int    `json:"fromBiometric" default:"0"`
}

type RefreshTokenDto struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
    FromAdmin    bool   `json:"fromAdmin" default:"false"`
}
