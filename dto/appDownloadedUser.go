package dto

type AppDownloadedUserDto struct {
    PlayerId  string `json:"playerId" validate:"required"`
    MachineId string `json:"machineId" validate:"required"`
}
