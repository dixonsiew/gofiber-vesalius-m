package model

import (
    "database/sql"
)

type DbAppServices struct {
    ServiceName  sql.NullString `db:"SERVICE_NAME"`
    ServiceImage sql.NullString `db:"SERVICE_IMAGE"`
}

type AppServices struct {
    ServiceName  string `json:"serviceName"`
    ServiceImage string `json:"serviceImage"`
}

func (o *AppServices) FromDbModel(m DbAppServices) {
    o.ServiceName = m.ServiceName.String
    o.ServiceImage = m.ServiceImage.String
}
