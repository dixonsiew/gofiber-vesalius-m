package novaInvestigation

import (
	"context"
	"strings"
	"vesaliusm/database"
	model "vesaliusm/model/healthCare"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

var NovaInvestigationSvc *NovaInvestigationService = NewNovaInvestigationService(database.GetDbrs(), database.GetCtx())

type NovaInvestigationService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaInvestigationService(db *sqlx.DB, ctx context.Context) *NovaInvestigationService {
    return &NovaInvestigationService{db: db, ctx: ctx}
}

func (s *NovaInvestigationService) FindNonPanelInvestiationByAccountNo(accountNo string, conn *sqlx.DB) ([]model.NovaPatientInvestigation, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM NOVA_PATIENT_INVESTIGATION WHERE ACCOUNT_NO = :accountNo AND PANEL_CODE IS NULL ORDER BY TO_NUMBER(INVESTIGATION_REF_NO)`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaPatientInvestigation{}, ""), 1)
    list := make([]model.NovaPatientInvestigation, 0)
    err := db.SelectContext(s.ctx, &list, query, accountNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaInvestigationService) FindPanelInvestigationByAccountNoPanelCode(accountNo string, panelCode string, conn *sqlx.DB) ([]model.NovaPatientInvestigation, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM NOVA_PATIENT_INVESTIGATION WHERE ACCOUNT_NO = :accountNo AND PANEL_CODE = :panelCode ORDER BY TO_NUMBER(INVESTIGATION_REF_NO)`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaPatientInvestigation{}, ""), 1)
    list := make([]model.NovaPatientInvestigation, 0)
    err := db.SelectContext(s.ctx, &list, query, accountNo, panelCode)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaInvestigationService) FindUniquePanelInvestigationByAccountNo(accountNo string, conn *sqlx.DB) ([]string, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT UNIQUE PANEL_CODE FROM NOVA_PATIENT_INVESTIGATION WHERE ACCOUNT_NO = :accountNo AND PANEL_CODE IS NOT NULL ORDER BY PANEL_CODE`
    list := make([]string, 0)
    err := db.SelectContext(s.ctx, &list, query, accountNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
