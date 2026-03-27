package maintenance

import (
	"context"
	"strings"
	"vesaliusm/database"

	"github.com/jmoiron/sqlx"
)

var MaintenanceSvc *MaintenanceService = NewMaintenanceService(database.GetDb(), database.GetCtx())

type MaintenanceService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewMaintenanceService(db *sqlx.DB, ctx context.Context) *MaintenanceService {
    return &MaintenanceService{
        db:  db,
        ctx: ctx,
    }
}

func whereClause(conds []string) string {
    if len(conds) == 0 {
        return ""
    }
    return " WHERE " + strings.Join(conds, " AND ")
}
