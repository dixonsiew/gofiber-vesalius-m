package model

import (
    "github.com/guregu/null/v6"
)

type ReleaseVersion struct {
    LatestVersion null.String `json:"latestVersion" db:"LATEST_VERSION" swaggertype:"string"`
    StackPlatform null.String `json:"stackPlatform" db:"STACK_PLATFORM" swaggertype:"string"`
    DateUpdate    null.String `json:"dateUpdate" db:"DATE_UPDATE" swaggertype:"string"`
}
