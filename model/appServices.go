package model

import (
    "github.com/guregu/null/v6"
)

type AppServices struct {
    ServiceName  null.String `json:"serviceName" db:"SERVICE_NAME"`
    ServiceImage null.String `json:"serviceImage" db:"SERVICE_IMAGE"`
}
