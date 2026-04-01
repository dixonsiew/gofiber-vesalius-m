package userGroupModulePermission

import (
    "context"
    "strings"
    "vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

var UserGroupModulePermissionSvc *UserGroupModulePermissionService = NewUserGroupModulePermissionService(database.GetDb(), database.GetCtx())

type UserGroupModulePermissionService struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewUserGroupModulePermissionService(db *sqlx.DB, ctx context.Context) *UserGroupModulePermissionService {
	return &UserGroupModulePermissionService{
		db:  db,
		ctx: ctx,
	}
}

func (s *UserGroupModulePermissionService) FindByUserGroupIdOrderByModuleIdAsc(userGroupId int64) ([]model.UserGroupModulePermission, error) {
	query := `SELECT * FROM USR_GRP_MOD_PERMERSSION WHERE USER_GROUP_ID = :userGroupId ORDER BY MODULE_ID ASC`
	query = strings.Replace(query, "*", utils.GetDbCols(model.UserGroupModulePermission{}, ""), 1)
    list := make([]model.UserGroupModulePermission, 0)
    err := s.db.SelectContext(s.ctx, &list, query, userGroupId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *UserGroupModulePermissionService) FindByUserGroupId(userGroupId int64) ([]model.UserGroupModulePermission, error) {
	query := `SELECT * FROM USR_GRP_MOD_PERMERSSION WHERE USER_GROUP_ID = :userGroupId`
	query = strings.Replace(query, "*", utils.GetDbCols(model.UserGroupModulePermission{}, ""), 1)
    list := make([]model.UserGroupModulePermission, 0)
    err := s.db.SelectContext(s.ctx, &list, query, userGroupId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
