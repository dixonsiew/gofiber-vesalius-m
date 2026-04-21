package novaPatientAlert

import (
    "context"
    "strings"
    "vesaliusm/database"
    model "vesaliusm/model/healthCare"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var NovaPatientAlertSvc *NovaPatientAlertService = NewNovaPatientAlertService(database.GetDbrs(), database.GetCtx())

type NovaPatientAlertService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaPatientAlertService(db *sqlx.DB, ctx context.Context) *NovaPatientAlertService {
    return &NovaPatientAlertService{db: db, ctx: ctx}
}

func (s *NovaPatientAlertService) FindNovaPatientAlertByPrn(prn string, conn *sqlx.DB) ([]model.NovaPatientAlert, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM NOVA_PATIENT_ALERT WHERE PRN = :prn`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaPatientAlert{}, ""), 1)
    list := make([]model.NovaPatientAlert, 0)
    err := db.GetContext(s.ctx, &list, query, prn)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaPatientAlertService) FindPatientActiveAlertByPrn(prn string, conn *sqlx.DB) ([]model.NovaPatientAlert, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM NOVA_PATIENT_ALERT WHERE PRN = :prn AND INACTIVE_DATE_TIME IS NULL ORDER BY ALERT_TYPE`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaPatientAlert{}, ""), 1)
    list := make([]model.NovaPatientAlert, 0)
    err := db.GetContext(s.ctx, &list, query, prn)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
