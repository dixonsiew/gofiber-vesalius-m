package userGroupModules

import (
    "context"
    "strings"
    "vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

var UserGroupModulesSvc *UserGroupModulesService = NewUserGroupModulesService(database.GetDb(), database.GetCtx())

type UserGroupModulesService struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewUserGroupModulesService(db *sqlx.DB, ctx context.Context) *UserGroupModulesService {
	return &UserGroupModulesService{
		db:  db,
		ctx: ctx,
	}
}

func (s *UserGroupModulesService) FindAllAsMap() (map[int64]model.UserGroupModules, error) {
	query := `SELECT * FROM USER_GROUP_MODULES`
    query = strings.Replace(query, "*", utils.GetDbCols(model.UserGroupModules{}, ""), 1)
    list := make([]model.UserGroupModules, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    result := make(map[int64]model.UserGroupModules)
    for i := range list {
        result[list[i].ModuleId.Int64] = list[i]
    }
    return result, nil
}

func (s *UserGroupModulesService) FindAll() ([]model.UserGroupModules, error) {
    query := `SELECT * FROM USER_GROUP_MODULES`
    query = strings.Replace(query, "*", utils.GetDbCols(model.UserGroupModules{}, ""), 1)
    list := make([]model.UserGroupModules, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
