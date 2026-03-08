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
    dbVar *sqlx.DB
    ctx = context.Background()
)

func SetDb(db *sqlx.DB) {
    dbVar = db
}

func GetDb() *sqlx.DB {
    if dbVar == nil {
        ConnectDB()
    }

    return dbVar
}

func GetCtx() context.Context {
    return ctx
}

func ConnectDB() {
    username := config.Config("db.username")
    pwd := config.Config("db.pwd")
    url := config.Config("db.url")
    connStr := fmt.Sprintf("oracle://%s:%s@%s", username, pwd, url)
    db, err := sqlx.Open("oracle", connStr)
    if err != nil {
        utils.LogError(err)
    } else {
        db.SetMaxOpenConns(10)
        db.SetMaxIdleConns(5)
        db.SetConnMaxLifetime(5 * time.Minute)
        db.SetConnMaxIdleTime(1 * time.Minute)
        SetDb(db)
        utils.LogInfo("Connection Opened to Database")
    }
}

func CloseDB() {
    dbVar.Close()
}
