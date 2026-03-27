package dto

type LoginDto struct {
    Username      string `json:"username" validate:"required" default:"eugene.lim@nova-hub.com"`
    Password      string `json:"password" validate:"required" default:"Abcd1234"`
    PlayerId      string `json:"playerId"`
    MachineId     string `json:"machineId"`
    FromAdmin     bool   `json:"fromAdmin" default:"false"`
    FromBiometric int    `json:"fromBiometric" validate:"numeric" default:"0"`
}
