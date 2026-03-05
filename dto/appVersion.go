package dto

type AppVersionDto struct {
	LatestVersion string `json:"latestVersion" validate:"required"`
	OsPlatform    string `json:"osPlatform" validate:"required"`
	Status        int    `json:"status" validate:"required,numeric"`
}
