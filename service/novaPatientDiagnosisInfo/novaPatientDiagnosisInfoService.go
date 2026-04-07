package novaPatientDiagnosisInfo

import (
	"context"
	"strings"
	"vesaliusm/database"
	model "vesaliusm/model/healthCare"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

var NovaPatientDiagnosisInfoSvc *NovaPatientDiagnosisInfoService = NewNovaPatientDiagnosisInfoService(database.GetDbrs(), database.GetCtx())

type NovaPatientDiagnosisInfoService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaPatientDiagnosisInfoService(db *sqlx.DB, ctx context.Context) *NovaPatientDiagnosisInfoService {
    return &NovaPatientDiagnosisInfoService{db: db, ctx: ctx}
}

func (s *NovaPatientDiagnosisInfoService) FindByDiagnosisRefNo(diagnosisRefNo string, conn *sqlx.DB) (*model.NovaPatientDiagnosisInfo, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM NOVA_PATIENT_DIAGNOSIS_INFO WHERE DIAGNOSIS_REF_NO = :diagnosisRefNo`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaPatientDiagnosis{}, ""), 1)
    var o model.NovaPatientDiagnosisInfo
    err := db.GetContext(s.ctx, &o, query, diagnosisRefNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}
