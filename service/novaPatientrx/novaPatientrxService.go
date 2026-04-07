package novaPatientrx

import (
	"context"
    "database/sql"
	"strings"
	"vesaliusm/database"
	model "vesaliusm/model/healthCare"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

var NovaPatientrxSvc *NovaPatientrxService = NewNovaPatientrxService(database.GetDbrs(), database.GetCtx())

type NovaPatientrxService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaPatientrxService(db *sqlx.DB, ctx context.Context) *NovaPatientrxService {
    return &NovaPatientrxService{db: db, ctx: ctx}
}

func (s *NovaPatientrxService) FindPatientRxByAccountNo(accountNo string, conn *sqlx.DB) ([]model.NovaPatientRx, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `
        SELECT * FROM 
        ( 
        SELECT * FROM NOVA_PATIENT_RX WHERE ACCOUNT_NO = :accountNo AND UPPER(EVENT) IN ('INTERVENE') 
        UNION ALL 
        SELECT * FROM NOVA_PATIENT_RX WHERE ACCOUNT_NO = :accountNo 
        AND UPPER(EVENT) NOT IN ('DISPENSE', 'CLOSE', 'HOME DRUG', 'RETURN') 
        AND ORDER_NO NOT IN (SELECT ORDER_NO FROM NOVA_PATIENT_RX WHERE ACCOUNT_NO = 'AN20-000142' AND UPPER(EVENT) IN ('INTERVENE')) 
        ) 
        ORDER BY ORDER_NO ASC
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaPatientRx{}, ""), 1)
    list := make([]model.NovaPatientRx, 0)
    err := db.SelectContext(s.ctx, &list, query, sql.Named("accountNo", accountNo))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
