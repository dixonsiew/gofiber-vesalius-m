package novaVisit

import (
    "context"
    "vesaliusm/database"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/jmoiron/sqlx"
)

type NovaVisitService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaVisitService(db *sqlx.DB, ctx context.Context) *NovaVisitService {
    return &NovaVisitService{db: db, ctx: ctx}
}

func (s *NovaVisitService) findByPrn(prn string, conn *sqlx.DB) (model.Nov)
