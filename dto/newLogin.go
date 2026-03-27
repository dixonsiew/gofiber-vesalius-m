package dto

type NewLoginDto struct {
    SignInType    int    `json:"signInType" validate:"required,numeric" default:"2"`
    Username      string `json:"username" validate:"required" default:"eugene.lim@nova-hub.com"`
    Password      string `json:"password" default:"Abcd1234"`
    PlayerId      string `json:"playerId" default:""`
    MachineId     string `json:"machineId" default:""`
    FromBiometric int    `json:"fromBiometric" validate:"numeric" default:"0"`
}
