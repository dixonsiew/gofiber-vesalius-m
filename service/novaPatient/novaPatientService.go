package novaPatient

import (
	"context"
	"strings"
	"vesaliusm/database"
	model "vesaliusm/model/healthCare"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

var NovaPatientSvc *NovaPatientService = NewNovaPatientService(database.GetDbrs(), database.GetCtx())

type NovaPatientService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaPatientService(db *sqlx.DB, ctx context.Context) *NovaPatientService {
    return &NovaPatientService{db: db, ctx: ctx}
}

func (s *NovaPatientService) FindByPrn(prn string, conn *sqlx.DB) (*model.NovaPatient, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM NOVA_PATIENT WHERE PRN = :prn`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaPatient{}, ""), 1)
    var o model.NovaPatient
    err := db.GetContext(s.ctx, &o, query, prn)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}
