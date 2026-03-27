package dto

type UserGroupDto struct {
	UserGroupId   int                  `json:"userGroupId" validate:"required"`
	UserGroupName string               `json:"userGroupName" validate:"required"`
	Permission    []GroupPermissionDto `json:"permission" validate:"required,min=1"`
}
