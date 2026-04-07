package novaBill

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    model "vesaliusm/model/healthCare"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var NovaBillSvc *NovaBillService = NewNovaBillService(database.GetDbrs(), database.GetCtx())

type NovaBillService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaBillService(db *sqlx.DB, ctx context.Context) *NovaBillService {
    return &NovaBillService{db: db, ctx: ctx}
}

func (s *NovaBillService) GetNovaBillByPrnAndAccountNo(prn string, accountNo string, conn *sqlx.DB) ([]model.NovaBill, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT * FROM NOVA_BILL WHERE PRN = :prn AND ACCOUNT_NO = :accountNo
        AND BILL_NO NOT IN (SELECT ORIGINAL_BILL_NO FROM NOVA_BILL WHERE PRN = :prn AND ACCOUNT_NO = :accountNo 
        AND BILL_TYPE = 'CANCEL') AND BILL_TYPE <> 'CANCEL'
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaBill{}, ""), 1)
    list := make([]model.NovaBill, 0)
    err := db.SelectContext(s.ctx, &list, query, sql.Named("prn", prn), sql.Named("accountNo", accountNo))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
