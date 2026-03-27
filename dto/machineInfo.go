package dto

type PostMachineInfo struct {
    MachineId string `json:"machineId" validate:"required"`
}
