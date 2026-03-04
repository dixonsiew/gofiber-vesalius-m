package model

import (
    "github.com/guregu/null/v6"
)

type ReleaseVersion struct {
    LatestVersion null.String `json:"latestVersion" db:"LATEST_VERSION"`
    StackPlatform null.String `json:"stackPlatform" db:"STACK_PLATFORM"`
    DateUpdate    null.String `json:"dateUpdate" db:"DATE_UPDATE"`
}
