package novaHealthScreenRpt

import (
    "context"
    "vesaliusm/database"
    model "vesaliusm/model/healthCare"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var NovaHealthScreenRptSvc *NovaHealthScreenRptService = NewNovaHealthScreenRptService(database.GetDbrs(), database.GetCtx())

type NovaHealthScreenRptService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaHealthScreenRptService(db *sqlx.DB, ctx context.Context) *NovaHealthScreenRptService {
    return &NovaHealthScreenRptService{db: db, ctx: ctx}
}

func (s *NovaHealthScreenRptService) FindHealthScreeningRptByPrnAndAccountNo(prn string, accountNo string, conn *sqlx.DB) ([]model.NovaHealthScreeningRpt, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT r.HSR_REF_NO, r.ACCOUNT_NO, r.REPORT_DATE, r.REPORT_USER 
        FROM NOVA_VISIT v, NOVA_ACCOUNT_HSR r, NOVA_ACCOUNT_HSR_CLOB c 
        WHERE v.PRN = :prn 
        AND v.ACCOUNT_NO = :accountNo 
        AND v.ACCOUNT_NO = r.ACCOUNT_NO 
        AND r.HSR_REF_NO = c.HSR_REF_NO 
        AND c.HSR_TYPE = 'HEALTH-SCREENING-REPORT'
    `
    list := make([]model.NovaHealthScreeningRpt, 0)
    err := db.SelectContext(s.ctx, &list, query, prn, accountNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
