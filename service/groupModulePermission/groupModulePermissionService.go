package groupModulePermission

import (
    "context"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var GroupModulePermissionSvc *GroupModulePermissionService = NewGroupModulePermissionService(database.GetDb(), database.GetCtx())

type GroupModulePermissionService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewGroupModulePermissionService(db *sqlx.DB, ctx context.Context) *GroupModulePermissionService {
    return &GroupModulePermissionService{
        db:  db,
        ctx: ctx,
    }
}

func (s *GroupModulePermissionService) FindAll() ([]model.GroupModulePermission, error) {
    query := `SELECT * FROM GROUP_MODULE_PERMISSION`
    query = strings.Replace(query, "*", utils.GetDbCols(model.GroupModulePermission{}, ""), 1)
    list := make([]model.GroupModulePermission, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
