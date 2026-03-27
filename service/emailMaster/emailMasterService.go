package emailMaster

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var EmailMasterSvc *EmailMasterService = NewEmailMasterService(database.GetDb(), database.GetCtx())

type EmailMasterService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewEmailMasterService(db *sqlx.DB, ctx context.Context) *EmailMasterService {
    return &EmailMasterService{
        db:  db,
        ctx: ctx,
    }
}

func (s *EmailMasterService) FindByEmailFunctionName(emailFunctionName string) (*model.EmailMaster, error) {
    query := `SELECT * FROM EMAIL_MASTER WHERE EMAIL_FUNCTION_NAME = :emailFunctionName`
    query = strings.Replace(query, "*", utils.GetDbCols(model.EmailMaster{}, ""), 1)
    var o model.EmailMaster
    err := s.db.GetContext(s.ctx, &o, query, emailFunctionName)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}
