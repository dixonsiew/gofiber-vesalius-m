package novaPatientDiagnosis

import (
    "context"
    "strings"
    "vesaliusm/database"
    model "vesaliusm/model/healthCare"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var NovaPatientDiagnosisSvc *NovaPatientDiagnosisService = NewNovaPatientDiagnosisService(database.GetDbrs(), database.GetCtx())

type NovaPatientDiagnosisService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaPatientDiagnosisService(db *sqlx.DB, ctx context.Context) *NovaPatientDiagnosisService {
    return &NovaPatientDiagnosisService{db: db, ctx: ctx}
}

func (s *NovaPatientDiagnosisService) FindByAccountNo(accountNo string, conn *sqlx.DB) (*model.NovaPatientDiagnosis, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM NOVA_PATIENT_DIAGNOSIS WHERE ACCOUNT_NO = :accountNo`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaPatientDiagnosis{}, ""), 1)
    var o model.NovaPatientDiagnosis
    err := db.GetContext(s.ctx, &o, query, accountNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}
