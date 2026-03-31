package userGroup

import (
	"context"
    "database/sql"
    "strings"
    "vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
	go_ora "github.com/sijms/go-ora/v2"
)

var UserGroupSvc *UserGroupService = NewUserGroupService(database.GetDb(), database.GetCtx())

type UserGroupService struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewUserGroupService(db *sqlx.DB, ctx context.Context) *UserGroupService {
	return &UserGroupService{
		db:  db,
		ctx: ctx,
	}
}

func (s *UserGroupService) Save(o model.UserGroup) error {
	tx, err := s.db.BeginTxx(s.ctx, nil)
    if err != nil {
        utils.LogError(err)
        return err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()

	query := `INSERT INTO USER_GROUP (DATE_CREATED, USER_GROUP_NAME) VALUES(SYSDATE, :userGroupName) RETURNING GROUP_ID INTO :group_id`
	var groupId go_ora.Number
	_, err = tx.ExecContext(s.ctx, query, sql.Named("userGroupName", o.UserGroupName.String), go_ora.Out{Dest: &groupId})
	if err != nil {
		utils.LogError(err)
		return err
	}

	q := `
	    INSERT INTO USR_GRP_MOD_PERMERSSION (USR_GRP_MOD_PERM_ID, MODULE_ID, PERMISSION_ID, USER_GROUP_ID) 
		VALUES(USR_GRP_MOD_PERM_SEQ.nextval, :moduleId, :permId, :userGroupId)
	`
	for _, r := range o.UserGroupModulePermissionStatesList {
		_, err = tx.ExecContext(s.ctx, q, 
			sql.Named("moduleId", r.ModuleId.Int64), 
			sql.Named("permId", r.PermissionId.Int64), 
			sql.Named("userGroupId", o.GroupId.Int64),
		)
		if err != nil {
			utils.LogError(err)
			return err
		}
	}
	
	err = tx.Commit()
    if err != nil {
        utils.LogError(err)
        return err
    }
    return err
}

func (s *UserGroupService) Update(o model.UserGroup) error {
    tx, err := s.db.BeginTxx(s.ctx, nil)
    if err != nil {
        utils.LogError(err)
        return err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()

    query := `DELETE FROM USR_GRP_MOD_PERMERSSION WHERE USER_GROUP_ID = :uGroupId`
    _, err = tx.ExecContext(s.ctx, query, o.GroupId.Int64)
    if err != nil {
        utils.LogError(err)
        return err
    }

	query = `UPDATE USER_GROUP SET USER_GROUP_NAME = :uGroupName WHERE GROUP_ID = :groupId`
	_, err = tx.ExecContext(s.ctx, query, o.UserGroupName.String, o.GroupId.Int64)
	if err != nil {
		utils.LogError(err)
		return err
	}

    q := `
        INSERT INTO USR_GRP_MOD_PERMERSSION (USR_GRP_MOD_PERM_ID, MODULE_ID, PERMISSION_ID, USER_GROUP_ID) 
		VALUES(USR_GRP_MOD_PERM_SEQ.nextval, :moduleId, :permId, :userGroupId)
    `
    for _, r := range o.UserGroupModulePermissionStatesList {
        _, err = tx.ExecContext(s.ctx, q, 
            sql.Named("moduleId", r.ModuleId.Int64), 
            sql.Named("permId", r.PermissionId.Int64), 
            sql.Named("userGroupId", o.GroupId.Int64),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }

	q = `UPDATE ADMIN_USER SET USER_GROUP_NAME = :uGroupName WHERE USER_GROUP_ID = :userGroupId`
	_, err = tx.ExecContext(s.ctx, q, o.UserGroupName.String, o.GroupId.Int64)
	if err != nil {
		utils.LogError(err)
		return err
	}
    
    err = tx.Commit()
    if err != nil {
        utils.LogError(err)
        return err
    }
    return err
}

func (s *UserGroupService) DeleteByGroupId(groupId int64) error {
	tx, err := s.db.BeginTxx(s.ctx, nil)
    if err != nil {
        utils.LogError(err)
        return err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()

	q := `DELETE FROM USR_GRP_MOD_PERMERSSION WHERE USER_GROUP_ID = :userGroupId`
	_, err = tx.ExecContext(s.ctx, q, groupId)
	if err != nil {
		utils.LogError(err)
		return err
	}

	q = `DELETE FROM USER_GROUP WHERE GROUP_ID = :groupId`
	_, err = tx.ExecContext(s.ctx, q, groupId)
	if err != nil {
		utils.LogError(err)
		return err
	}

	q = `UPDATE ADMIN_USER SET USER_GROUP_ID = NULL, USER_GROUP_NAME = NULL WHERE USER_GROUP_ID = :userGroupId`
	_, err = tx.ExecContext(s.ctx, q, groupId)
	if err != nil {
		utils.LogError(err)
		return err
	}
	
	err = tx.Commit()
	if err != nil {
		utils.LogError(err)
		return err
	}
	return err
}

func (s *UserGroupService) ListAll() ([]model.UserGroup, error) {
	query := `SELECT * FROM USER_GROUP ORDER BY USER_GROUP_NAME`
	query = strings.Replace(query, "*", utils.GetDbCols(model.UserGroup{}, ""), 1)
	list := make([]model.UserGroup, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *UserGroupService) List(page string, limit string) (*model.PagedList, error) {
    total, err := s.Count()
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAll(pager.GetLowerBound(), pager.PageSize, s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *UserGroupService) FindAll(offset int, limit int, conn *sqlx.DB) ([]model.UserGroup, error) {
	db := conn
    if db == nil {
        db = s.db
    }
	query := `SELECT * FROM USER_GROUP ORDER BY GROUP_ID OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
	query = strings.Replace(query, "*", utils.GetDbCols(model.UserGroup{}, ""), 1)
	list := make([]model.UserGroup, 0)
    err := conn.SelectContext(context.Background(), &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *UserGroupService) Count() (int, error) {
    var count int
    query := `SELECT COUNT(GROUP_ID) AS COUNT FROM USER_GROUP`
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *UserGroupService) ExistsByOtherUserGroupName(userGroupName string, groupId int64) (bool, error) {
    query := `SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM USER_GROUP WHERE USER_GROUP_NAME = :uGroupName AND GROUP_ID <> :groupId)`
	var count int
    err := s.db.GetContext(s.ctx, &count, query, userGroupName, groupId)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, err
}

func (s *UserGroupService) ExistsByUserGroupName(userGroupName string) (bool, error) {
    query := `SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM USER_GROUP WHERE USER_GROUP_NAME = :uGroupName)`
	var count int
    err := s.db.GetContext(s.ctx, &count, query, userGroupName)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, err
}

func (s *UserGroupService) ExistsByGroupId(groupId int64) (bool, error) {
    query := `SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM USER_GROUP WHERE GROUP_ID = :groupId)`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, groupId)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, err
}

func (s *UserGroupService) FindByGroupId(groupId int64) (*model.UserGroup, error) {
    query := `SELECT * FROM USER_GROUP WHERE GROUP_ID = :groupId`
    var o model.UserGroup
    err := s.db.GetContext(s.ctx, &o, query, groupId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *UserGroupService) FindByUserGroupName(userGroupName string) (*model.UserGroup, error) {
    query := `SELECT * FROM USER_GROUP WHERE USER_GROUP_NAME = :uGroupName`
    var o model.UserGroup
    err := s.db.GetContext(s.ctx, &o, query, userGroupName)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}
