package dto

type ReleaseVersionDto struct {
	LatestVersion string `json:"latestVersion" validate:"required"`
	StackPlatform string `json:"stackPlatform" validate:"required"`
}
