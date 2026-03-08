package database

import (
    "context"
    "fmt"
    "time"
    "vesaliusm/config"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    _ "github.com/sijms/go-ora/v2"
)

var (
    dbrsVar *sqlx.DB
    ctxrs = context.Background()
)

func SetDbrs(db *sqlx.DB) {
    dbrsVar = db
}

func GetDbrs() *sqlx.DB {
    if dbrsVar == nil {
        ConnectDBRs()
    }

    return dbrsVar
}

func GetCtxrs() context.Context {
    return ctxrs
}

func ConnectDBRs() {
    username := config.Config("db.rs.username")
    pwd := config.Config("db.rs.pwd")
    url := config.Config("db.rs.url")
    connStr := fmt.Sprintf("oracle://%s:%s@%s", username, pwd, url)
    db, err := sqlx.Open("oracle", connStr)
    if err != nil {
        utils.LogError(err)
    } else {
        db.SetMaxOpenConns(10)
        db.SetMaxIdleConns(5)
        db.SetConnMaxLifetime(5 * time.Minute)
        db.SetConnMaxIdleTime(1 * time.Minute)
        SetDbrs(db)
        utils.LogInfo("Connection Opened to Database")
    }
}

func CloseDBRs() {
    dbrsVar.Close()
}
