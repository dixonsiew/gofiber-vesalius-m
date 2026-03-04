package model

import (
    "github.com/guregu/null/v6"
)

type AppVersion struct {
    LatestVersion null.String `json:"latestVersion" db:"LATEST_VERSION" swaggertype:"string"`
    OSPlatform    null.String `json:"osPlatform" db:"OS_PLATFORM" swaggertype:"string"`
    Status        null.Int64  `json:"status" db:"STATUS" swaggertype:"integer"`
}
