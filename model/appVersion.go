package model

import (
    "github.com/guregu/null/v6"
)

type AppVersion struct {
    LatestVersion null.String `json:"latestVersion" db:"LATEST_VERSION"`
    OSPlatform    null.String `json:"osPlatform" db:"OS_PLATFORM"`
    Status        null.Int64  `json:"status" db:"STATUS"`
}
