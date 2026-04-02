package requestDeleteAccount

import (
    "context"
    "database/sql"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var RequestDeleteAccountSvc *RequestDeleteAccountService = NewRequestDeleteAccountService(database.GetDb(), database.GetCtx())

type RequestDeleteAccountService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewRequestDeleteAccountService(db *sqlx.DB, ctx context.Context) *RequestDeleteAccountService {
    return &RequestDeleteAccountService{
        db:  db,
        ctx: ctx,
    }
}

func (s *RequestDeleteAccountService) SaveRequest(o model.RequestDeleteAccount) error {
    query := `
        INSERT INTO REQUEST_ACCOUNT_DELETE
        (PRN, FULLNAME, DOCUMENT_NUMBER, DOB, CONTACT_NUMBER, EMAIL)
         VALUES
        (:prn, :fullname, :documentNumber, :dob, :contactNumber, :email)
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("prn", o.PRN),
        sql.Named("fullname", o.Fullname),
        sql.Named("documentNumber", o.DocumentNumber),
        sql.Named("dob", o.DOB),
        sql.Named("contactNumber", o.ContactNumber),
        sql.Named("email", o.Email),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    
    return nil
}
