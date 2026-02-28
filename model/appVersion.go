package model

import (
    "database/sql"
)

type DbAppVersion struct {
    LatestVersion sql.NullString `db:"LATEST_VERSION"`
    OSPlatform    sql.NullString `db:"OS_PLATFORM"`
    Status        sql.NullInt64  `db:"STATUS"`
}

type AppVersion struct {
    LatestVersion string `json:"latestVersion"`
    OSPlatform    string `json:"osPlatform"`
    Status        int64  `json:"status"`
}

func (o *AppVersion) FromDbModel(m DbAppVersion) {
    o.LatestVersion = m.LatestVersion.String
    o.OSPlatform = m.OSPlatform.String
    o.Status = m.Status.Int64
}
