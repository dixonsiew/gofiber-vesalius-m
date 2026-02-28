package model

import (
    "database/sql"
)

type DbReleaseVersion struct {
    LatestVersion sql.NullString `db:"LATEST_VERSION"`
    StackPlatform sql.NullString `db:"STACK_PLATFORM"`
    DateUpdate    sql.NullString `db:"DATE_UPDATE"`
}

type ReleaseVersion struct {
    LatestVersion string `json:"latestVersion"`
    StackPlatform string `json:"stackPlatform"`
    DateUpdate    string `json:"dateUpdate"`
}

func (o *ReleaseVersion) FromDbModel(m DbReleaseVersion) {
    o.LatestVersion = m.LatestVersion.String
    o.StackPlatform = m.StackPlatform.String
    o.DateUpdate = m.DateUpdate.String
}
