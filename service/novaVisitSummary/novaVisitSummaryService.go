package novaVisitSummary

import (
	"context"
	"strings"
	"vesaliusm/database"
	model "vesaliusm/model/healthCare"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

var NovaVisitSummarySvc *NovaVisitSummaryService = NewNovaVisitSummaryService(database.GetDbrs(), database.GetCtx())

type NovaVisitSummaryService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaVisitSummaryService(db *sqlx.DB, ctx context.Context) *NovaVisitSummaryService {
    return &NovaVisitSummaryService{db: db, ctx: ctx}
}

func (s *NovaVisitSummaryService) FindByAccountNo(accountNo string, conn *sqlx.DB) ([]model.NovaVisitSummary, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM NOVA_VISIT_SUMMARY WHERE ACCOUNT_NO = :accountNo`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaVisitSummary{}, ""), 1)
    list := make([]model.NovaVisitSummary, 0)
    err := db.SelectContext(s.ctx, &list, query, accountNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVisitSummaryService) FindByCategory(category string, conn *sqlx.DB) ([]model.NovaVisitSummary, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM NOVA_VISIT_SUMMARY WHERE CATEGORY = :category`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaVisitSummary{}, ""), 1)
    list := make([]model.NovaVisitSummary, 0)
    err := db.SelectContext(s.ctx, &list, query, category)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVisitSummaryService) FindByCategoryAndAccountNo(category string, accountNo string, conn *sqlx.DB) ([]model.NovaVisitSummary, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM NOVA_VISIT_SUMMARY WHERE CATEGORY = :category AND ACCOUNT_NO = :accountNo`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaVisitSummary{}, ""), 1)
    list := make([]model.NovaVisitSummary, 0)
    err := db.SelectContext(s.ctx, &list, query, category, accountNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVisitSummaryService) FindByAccountNoAndCategory(accountNo string, category string, conn *sqlx.DB) ([]model.NovaVisitSummary, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM NOVA_VISIT_SUMMARY WHERE ACCOUNT_NO = :accountNo AND CATEGORY = :category`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaVisitSummary{}, ""), 1)
    list := make([]model.NovaVisitSummary, 0)
    err := db.SelectContext(s.ctx, &list, query, accountNo, category)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVisitSummaryService) GetByAccountNoAndNotInCategories(accountNo string, firstCategory string, secondCategory string, conn *sqlx.DB) ([]model.NovaVisitSummary, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM NOVA_VISIT_SUMMARY WHERE ACCOUNT_NO = :accountNo AND CATEGORY = :firstCategory`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaVisitSummary{}, ""), 1)
    list := make([]model.NovaVisitSummary, 0)
    err := db.SelectContext(s.ctx, &list, query, accountNo, firstCategory)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
