package adminUser

import (
    "context"
    "database/sql"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    "golang.org/x/crypto/bcrypt"
)

type AdminUserService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewAdminUserService(db *sqlx.DB, ctx context.Context) *AdminUserService {
    return &AdminUserService{db: db, ctx: ctx}
}

const saltRounds = 10

func (s *AdminUserService) FindAll(offset int, limit int) ([]model.AdminUser, error) {
    list := make([]model.AdminUser, 0)
    err := s.db.SelectContext(s.ctx, &list, `SELECT * FROM ADMIN_USER ORDER BY FIRST_NAME OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AdminUserService) List(page string, limit string) (*model.PagedList, error) {
    total, err := s.Count()
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAll(pager.GetLowerBound(), pager.PageSize)
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

func (s *AdminUserService) Count() (int, error) {
    var count int
    query := `SELECT COUNT(ADMIN_ID) AS COUNT FROM ADMIN_USER`
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *AdminUserService) FindByAdminId(adminId int64) (*model.AdminUser, error) {
    query := `SELECT * FROM ADMIN_USER WHERE ADMIN_ID = :adminId`
    var o model.AdminUser
    err := s.db.GetContext(s.ctx, &o, query, adminId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, err
}

func (s *AdminUserService) FindByEmail(email string) (*model.AdminUser, error) {
    query := `SELECT * FROM ADMIN_USER WHERE EMAIL = :email`
    var o model.AdminUser
    err := s.db.GetContext(s.ctx, &o, query, email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, err
}

func (s *AdminUserService) FindByUsername(email string) (*model.AdminUser, error) {
    query := `SELECT * FROM ADMIN_USER WHERE USERNAME = :email`
    var o model.AdminUser
    err := s.db.GetContext(s.ctx, &o, query, email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, err
}

func (s *AdminUserService) FindByUserGroupId(userGroupId int64) ([]model.AdminUser, error) {
    list := make([]model.AdminUser, 0)
    err := s.db.SelectContext(s.ctx, &list, `SELECT * FROM ADMIN_USER WHERE USER_GROUP_ID = :userGroupId`, userGroupId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AdminUserService) FindWithAssignBranchByAdminId(adminId int64) (*model.AdminUser, error) {
    query := `
        SELECT ` + getAdminUserCols() + `, ` + getAssignBranchCols() + `, ` + getBranchCols() + ` 
        FROM ADMIN_USER au 
        LEFT JOIN ASSIGN_BRANCH ab ON au.ADMIN_ID = ab.ADMIN_ID 
        INNER JOIN BRANCH b ON b.BRANCH_ID = ab.BRANCH_ID 
        WHERE au.ADMIN_ID = :adminId
    `
    rows, err := s.db.QueryxContext(s.ctx, query, adminId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    var o model.AdminUser
    var user *model.AdminUser
    var branches []model.AssignBranch

    for rows.Next() {
        if user == nil {
            err = rows.StructScan(&o)
            if err != nil {
                utils.LogError(err)
                return nil, err
            }
            user = &o
        }

        ab := model.AssignBranch{}
        err = rows.StructScan(&ab)
        if err != nil {
            utils.LogError(err)
            return nil, err
        }

        b := model.Branch{}
        err = rows.StructScan(&b)
        if err != nil {
            utils.LogError(err)
            return nil, err
        }

        b.Passcode.String = ""
        b.Url.String = ""
        ab.Branch = &b
        branches = append(branches, ab)
    }
    if user == nil {
        return nil, nil
    }
    user.Password.String = ""
    user.AdminBranches = branches
    return user, nil
}

func (s *AdminUserService) FindAssignBranchByAdminId(adminId int64, branchId int64) (*model.AssignBranch, error) {
    query := `SELECT ` + getAssignBranchCols() + ` FROM ASSIGN_BRANCH WHERE BRANCH_ID = :branchId AND ADMIN_ID IN (SELECT ADMIN_ID FROM ADMIN_USER WHERE ADMIN_ID = :adminId)`
    var ab model.AssignBranch
    err := s.db.GetContext(s.ctx, &ab, query, branchId, adminId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &ab, err
}

func (s *AdminUserService) ValidateCredentials(user model.AdminUser, password string) bool {
    if user.Password.String == "" {
        return false
    }
    err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
    return err == nil
}

func getAdminUserCols() string {
    return `
        ADMIN_ID,
        USERNAME,
        EMAIL,
        PASSWORD,
        TITLE,
        FIRST_NAME,
        MIDDLE_NAME,
        LAST_NAME,
        RESIDENT,
        DOB,
        SEX,
        RACE,
        ADDRESS,
        CONTACT_NUMBER,
        PASSPORT,
        NATIONALITY,
        ROLE,
        USER_GROUP_ID,
        USER_GROUP_NAME
    `
}

func getBranchCols() string {
    return `
        BRANCH_ID,
        URL,
        PASSCODE,
        BRANCH_NAME
    `
}

func getAssignBranchCols() string {
    return `
        ASSIGN_BRANCH_ID,
        PRN,
        USER_ID,
        ADMIN_ID,
        BRANCH_ID
    `
}
