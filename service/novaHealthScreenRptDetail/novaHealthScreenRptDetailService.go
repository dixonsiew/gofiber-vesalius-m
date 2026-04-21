package novaHealthScreenRptDetail

import (
    "context"
    "vesaliusm/database"
    model "vesaliusm/model/healthCare"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var NovaHealthScreenRptDetailSvc *NovaHealthScreenRptDetailService = NewNovaHealthScreenRptDetailService(database.GetDbrs(), database.GetCtx())

type NovaHealthScreenRptDetailService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaHealthScreenRptDetailService(db *sqlx.DB, ctx context.Context) *NovaHealthScreenRptDetailService {
    return &NovaHealthScreenRptDetailService{db: db, ctx: ctx}
}

func (s *NovaHealthScreenRptDetailService) FindEachHealthScreeningRptByHSRRefNo(hsrRefNo string, conn *sqlx.DB) ([]model.NovaHealthScreeningDetailRpt, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT dbms_random.string('U', 30) as HSR_REF_NO, c.HSR_CLOB_VALUE 
        FROM NOVA_ACCOUNT_HSR_CLOB c 
        WHERE c.HSR_REF_NO = :hsrRefNo 
        AND c.HSR_TYPE = 'HEALTH-SCREENING-REPORT'
    `
    list := make([]model.NovaHealthScreeningDetailRpt, 0)
    err := db.SelectContext(s.ctx, &list, query, hsrRefNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
