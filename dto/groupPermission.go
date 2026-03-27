package dto

type GroupPermissionDto struct {
	ModuleId     int `json:"moduleId" validate:"required,numeric"`
	PermissionId int `json:"permissionId" validate:"required,numeric"`
}
