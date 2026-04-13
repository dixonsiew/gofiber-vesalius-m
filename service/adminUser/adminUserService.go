package adminUser

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/dto"
    "vesaliusm/model"
    "vesaliusm/service/assignBranch"
    "vesaliusm/service/branch"
    "vesaliusm/service/groupModulePermission"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    go_ora "github.com/sijms/go-ora/v2"
    "golang.org/x/crypto/bcrypt"
)

var AdminUserSvc *AdminUserService = NewAdminUserService(database.GetDb(), database.GetCtx())

type AdminUserService struct {
    db                       *sqlx.DB
    ctx                      context.Context
    branchSvc                *branch.BranchService
    assignBranchSvc          *assignBranch.AssignBranchService
    groupModulePermissionSvc *groupModulePermission.GroupModulePermissionService
}

func NewAdminUserService(db *sqlx.DB, ctx context.Context) *AdminUserService {
    return &AdminUserService{
        db:                       db,
        ctx:                      ctx,
        branchSvc:                branch.BranchSvc,
        assignBranchSvc:          assignBranch.AssignBranchSvc,
        groupModulePermissionSvc: groupModulePermission.GroupModulePermissionSvc,
    }
}

const saltRounds = 10

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

func (s *AdminUserService) FindAll(offset int, limit int) ([]model.AdminUser, error) {
    query := `SELECT * FROM ADMIN_USER ORDER BY FIRST_NAME OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(model.AdminUser{}, ""), 1)
    list := make([]model.AdminUser, 0)
    err := s.db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
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

func (s *AdminUserService) ListMobileUserAuditLog(page string, limit string) (*model.PagedList, error) {
    total, err := s.MobileUserAuditLogCount()
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllMobileUserAuditLog(pager.GetLowerBound(), pager.PageSize)
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

func (s *AdminUserService) MobileUserAuditLogCount() (int, error) {
    var count int
    query := `SELECT COUNT(AUDIT_ID) AS COUNT FROM AUDIT_MOBILE_USER WHERE ACTION = 'Delete Account'`
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *AdminUserService) FindAllMobileUserAuditLog(offset int, limit int) ([]model.MobileUserAuditLog, error) {
    query := `
        SELECT amu.* FROM AUDIT_MOBILE_USER amu
         WHERE ACTION = 'Delete Account'
         ORDER BY amu.DATE_CREATE DESC 
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "amu.*", utils.GetDbCols(model.MobileUserAuditLog{}, "amu."), 1)
    list := make([]model.MobileUserAuditLog, 0)
    err := s.db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AdminUserService) ListMobileUserAuditLogByKeyword(keyword string, keyword2 string, page string, limit string) (*model.PagedList, error) {
    total, err := s.MobileUserAuditLogCountByKeyword(keyword, keyword2)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.MobileUserAuditLogFindByKeyword(keyword, keyword2, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *AdminUserService) MobileUserAuditLogCountByKeyword(keyword string, keyword2 string) (int, error) {
    conds, args := buildKeywordConditions(keyword, keyword2)
    base := `SELECT COUNT(amu.AUDIT_ID) AS COUNT FROM AUDIT_MOBILE_USER amu`
    query := base + whereClause(conds)

    var count int
    err := s.db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *AdminUserService) MobileUserAuditLogFindByKeyword(keyword string, keyword2 string, offset int, limit int) ([]model.MobileUserAuditLog, error) {
    conditions, args := buildKeywordConditions(keyword, keyword2)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    base := `SELECT amu.* FROM AUDIT_MOBILE_USER amu`
    query := base + whereClause(conditions) +
        ` ORDER BY amu.DATE_CREATE DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "amu.*", utils.GetDbCols(model.MobileUserAuditLog{}, "amu."), 1)

    list := make([]model.MobileUserAuditLog, 0)
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AdminUserService) ListAuditLog(page string, limit string) (*model.PagedList, error) {
    total, err := s.AuditLogCount()
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllAuditLog(pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *AdminUserService) FindAllAuditLog(offset int, limit int) ([]model.AdminAuditLog, error) {
    m := map[string]string{
        "apal.EVENT_ADMIN_EMAIL": "",
    }
    query := `
        SELECT apal.*, au.EMAIL AS EVENT_ADMIN_EMAIL
         FROM ADMIN_PORTAL_AUDIT_LOG apal
         JOIN ADMIN_USER au ON apal.EVENT_ADMIN_ID = au.ADMIN_ID
         ORDER BY apal.EVENT_DATE_TIME DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "apal.*", utils.GetDbColsWithReplace(model.AdminAuditLog{}, "apal.", m), 1)
    list := make([]model.AdminAuditLog, 0)
    err := s.db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AdminUserService) AuditLogCount() (int, error) {
    var count int
    query := `SELECT COUNT(EVENT_ID) AS COUNT FROM ADMIN_PORTAL_AUDIT_LOG`
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *AdminUserService) ListAuditByKeyword(keyword string, keyword2 string, page string, limit string) (*model.PagedList, error) {
    total, err := s.AuditCountByKeyword(keyword, keyword2)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.AuditFindByKeyword(keyword, keyword2, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *AdminUserService) AuditCountByKeyword(keyword string, keyword2 string) (int, error) {
    conds, args := buildAdminAuditLogConditions(keyword, keyword2)
    base := `
        SELECT COUNT(apal.EVENT_ID) AS COUNT FROM ADMIN_PORTAL_AUDIT_LOG apal
        JOIN ADMIN_USER au ON apal.EVENT_ADMIN_ID = au.ADMIN_ID`
    query := base + whereClause(conds)

    var count int
    err := s.db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *AdminUserService) AuditFindByKeyword(keyword string, keyword2 string, offset int, limit int) ([]model.AdminAuditLog, error) {
    conditions, args := buildAdminAuditLogConditions(keyword, keyword2)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    m := map[string]string{
        "apal.EVENT_ADMIN_EMAIL": "",
    }
    base := `
        SELECT apal.*, au.EMAIL AS EVENT_ADMIN_EMAIL 
        FROM ADMIN_PORTAL_AUDIT_LOG apal
        JOIN ADMIN_USER au ON apal.EVENT_ADMIN_ID = au.ADMIN_ID
    `
    base = strings.Replace(base, "apal.*", utils.GetDbColsWithReplace(model.AdminAuditLog{}, "apal.", m), 1)
    query := base + whereClause(conditions) +
        ` ORDER BY apal.EVENT_DATE_TIME DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    list := make([]model.AdminAuditLog, 0)
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AdminUserService) ListByKeyword(keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountByKeyword(keyword)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindByKeyword(keyword, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *AdminUserService) FindByKeyword(keyword string, offset int, limit int) ([]model.AdminUser, error) {
    query := `
        SELECT au.* FROM ADMIN_USER au
        WHERE LOWER(au.FIRST_NAME) LIKE :keyword OR LOWER(au.MIDDLE_NAME) LIKE :keyword OR LOWER(au.LAST_NAME) LIKE :keyword
        OR LOWER(au.EMAIL) LIKE :keyword
        ORDER BY au.FIRST_NAME OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "au.*", utils.GetDbCols(model.AdminUser{}, "au."), 1)
    list := make([]model.AdminUser, 0)
    err := s.db.SelectContext(s.ctx, &list, query,
        sql.Named("keyword", strings.ToLower(keyword)),
        sql.Named("offset", offset),
        sql.Named("limit", limit),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AdminUserService) CountByKeyword(keyword string) (int, error) {
    var count int
    query := `
        SELECT COUNT(au.ADMIN_ID) AS COUNT FROM ADMIN_USER au
        WHERE LOWER(au.FIRST_NAME) LIKE :keyword OR LOWER(au.MIDDLE_NAME) LIKE :keyword OR LOWER(au.LAST_NAME) LIKE :keyword
        OR LOWER(au.EMAIL) LIKE :keyword
    `
    err := s.db.GetContext(s.ctx, &count, query, sql.Named("keyword", strings.ToLower(keyword)))
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *AdminUserService) FindByAdminId(adminId int64) (*model.AdminUser, error) {
    query := `SELECT * FROM ADMIN_USER WHERE ADMIN_ID = :adminId`
    query = strings.Replace(query, "*", utils.GetDbCols(model.AdminUser{}, ""), 1)
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
    query = strings.Replace(query, "*", utils.GetDbCols(model.AdminUser{}, ""), 1)
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
    query = strings.Replace(query, "*", utils.GetDbCols(model.AdminUser{}, ""), 1)
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
    query := `SELECT * FROM ADMIN_USER WHERE USER_GROUP_ID = :userGroupId`
    query = strings.Replace(query, "*", utils.GetDbCols(model.AdminUser{}, ""), -1)
    list := make([]model.AdminUser, 0)
    err := s.db.SelectContext(s.ctx, &list, query, userGroupId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AdminUserService) FindWithAssignBranchByAdminId(adminId int64) (*model.AdminUser, error) {
    o, err := s.FindByAdminId(adminId)
    if err != nil {
        return nil, err
    }

    ablist, err := s.assignBranchSvc.FindAllByAdminId(o.AdminId.Int64)
    if err != nil {
        return nil, err
    }

    for i := range ablist {
        b, err := s.branchSvc.FindByBranchId(ablist[i].BranchId.Int64)
        if err != nil {
            return nil, err
        }
        b.Passcode = utils.NewNullString("")
        b.Url = utils.NewNullString("")
        ablist[i].Branch = b
    }
    o.Password = utils.NewNullString("")
    o.AdminBranches = ablist
    return o, nil
}

func (s *AdminUserService) FindAssignBranchByAdminId(adminId int64, branchId int64) (*model.AssignBranch, error) {
    query := `SELECT * FROM ASSIGN_BRANCH WHERE BRANCH_ID = :branchId AND ADMIN_ID IN (SELECT ADMIN_ID FROM ADMIN_USER WHERE ADMIN_ID = :adminId)`
    query = strings.Replace(query, "*", utils.GetDbCols(model.AssignBranch{}, ""), 1)
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

func (s *AdminUserService) FindAssignBranchByEmail(email string, branchId int64) (*model.AssignBranch, error) {
    query := `SELECT * FROM ASSIGN_BRANCH WHERE BRANCH_ID = :branchId AND ADMIN_ID IN (SELECT ADMIN_ID FROM ADMIN_USER WHERE EMAIL = :email)`
    var ab model.AssignBranch
    err := s.db.GetContext(s.ctx, &ab, query, branchId, email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &ab, err
}

func (s *AdminUserService) ExistsByAdminId(adminId int64) (bool, error) {
    query := `SELECT COUNT(ADMIN_ID) AS COUNT FROM ADMIN_USER WHERE ADMIN_ID = :adminId`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, adminId)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, err
}

func (s *AdminUserService) ExistsByEmail(email string) (bool, error) {
    query := `SELECT COUNT(ADMIN_ID) AS COUNT FROM ADMIN_USER WHERE EMAIL = :email`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, email)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, err
}

func (s *AdminUserService) SaveResetPassword(o *model.AdminUser) error {
    newPwd := utils.GetRandomStr(6)
    hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPwd), saltRounds)
    if err != nil {
        utils.LogError(err)
        return err
    }
    query := `UPDATE ADMIN_USER SET PASSWORD = :pw WHERE ADMIN_ID = :adminId`
    _, err = s.db.ExecContext(s.ctx, query,
        sql.Named("pw", string(hashedPwd)),
        sql.Named("adminId", o.AdminId.Int64),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    o.Password = utils.NewNullString(newPwd)
    return err
}

func (s *AdminUserService) SavePassword(o *model.AdminUser) error {
    hashedPwd, err := bcrypt.GenerateFromPassword([]byte(o.Password.String), saltRounds)
    if err != nil {
        utils.LogError(err)
        return err
    }
    query := `UPDATE ADMIN_USER SET PASSWORD = :pw WHERE ADMIN_ID = :adminId`
    _, err = s.db.ExecContext(s.ctx, query,
        sql.Named("pw", string(hashedPwd)),
        sql.Named("adminId", o.AdminId.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *AdminUserService) Save(o *model.AdminUser, adminBranchIds []int64) error {
    newPwd := utils.GetRandomStr(6)
    o.Password = utils.NewNullString(newPwd)
    hashedPwd, err := bcrypt.GenerateFromPassword([]byte(o.Password.String), saltRounds)
    if err != nil {
        utils.LogError(err)
        return err
    }

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

    query := `
        INSERT INTO ADMIN_USER
        (ADMIN_ID, EMAIL, FIRST_NAME, LAST_NAME, PASSWORD, "ROLE", USER_GROUP_ID, USER_GROUP_NAME, USERNAME)
        VALUES(ADMIN_USER_SEQ.nextval, :email, :first_name, :last_name, :password, :role, :user_group_id, :user_group_name, :username) RETURNING ADMIN_ID INTO :admin_id
    `
    var adminId go_ora.Number
    _, err = tx.ExecContext(s.ctx, query,
        sql.Named("email", o.Email.String),
        sql.Named("first_name", o.FirstName.String),
        sql.Named("last_name", o.LastName.String),
        sql.Named("password", string(hashedPwd)),
        sql.Named("role", o.Role.String),
        sql.Named("user_group_id", o.UserGroupId.Int64),
        sql.Named("user_group_name", o.UserGroupName.String),
        sql.Named("username", o.Username.String),
        go_ora.Out{Dest: &adminId},
    )
    if err != nil {
        utils.LogError(err)
        return err
    }

    iadminId, _ := adminId.Int64()
    o.AdminId.Int64 = iadminId

    for _, adminBranchId := range adminBranchIds {
        _, err = tx.ExecContext(s.ctx, `INSERT INTO ASSIGN_BRANCH (ASSIGN_BRANCH_ID, ADMIN_ID, BRANCH_ID) VALUES(USER_BRANCH_SEQ.nextval, :adminId, :branchId)`,
            sql.Named("adminId", o.AdminId.Int64),
            sql.Named("branchId", adminBranchId),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }

    return tx.Commit()
}

func (s *AdminUserService) Update(o *model.AdminUser, adminBranchIds []int64) error {
    _, err := s.db.ExecContext(s.ctx, `DELETE FROM ASSIGN_BRANCH WHERE ADMIN_ID = :adminId`, o.AdminId.Int64)
    if err != nil {
        utils.LogError(err)
        return err
    }

    for _, adminBranchId := range adminBranchIds {
        _, err = s.db.ExecContext(s.ctx, `INSERT INTO ASSIGN_BRANCH (ASSIGN_BRANCH_ID, ADMIN_ID, BRANCH_ID) VALUES(USER_BRANCH_SEQ.nextval, :adminId, :branchId)`,
            sql.Named("adminId", o.AdminId.Int64),
            sql.Named("branchId", adminBranchId),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }

    _, err = s.db.ExecContext(s.ctx, `UPDATE ADMIN_USER SET FIRST_NAME = :first_name, LAST_NAME = :last_name,  "ROLE" = :role, USER_GROUP_ID = :user_group_id, USER_GROUP_NAME = :user_group_name WHERE ADMIN_ID = :admin_id`,
        sql.Named("first_name", o.FirstName.String),
        sql.Named("last_name", o.LastName.String),
        sql.Named("role", o.Role.String),
        sql.Named("user_group_id", o.UserGroupId.Int64),
        sql.Named("user_group_name", o.UserGroupName.String),
        sql.Named("admin_id", o.AdminId.Int64),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func (s *AdminUserService) Delete(adminId int64) error {
    _, err := s.db.ExecContext(s.ctx, `DELETE FROM ASSIGN_BRANCH WHERE ADMIN_ID = :adminId`, adminId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    _, err = s.db.ExecContext(s.ctx, `DELETE FROM ADMIN_USER WHERE ADMIN_ID = :adminId`, adminId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func (s *AdminUserService) ChangeUserPassword(newPw string, adminId int64) error {
    hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPw), saltRounds)
    if err != nil {
        utils.LogError(err)
        return err
    }

    _, err = s.db.ExecContext(s.ctx, `UPDATE APPLICATION_USER SET PASSWORD = :pw WHERE USER_ID = :adminId`,
        sql.Named("pw", string(hashedPwd)),
        sql.Named("adminId", adminId),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func (s *AdminUserService) SaveAdminPortalLog(o dto.AdminPortalLogDto, adminId int64) error {
    query := `
        INSERT INTO ADMIN_PORTAL_AUDIT_LOG
        (EVENT_ADMIN_ID, EVENT_MODULE, EVENT_FUNCTION, EVENT_ACTION, EVENT_KEYWORD, EVENT_DESCRIPTION, PATIENT_PRN)
        VALUES
        (:admin_id, :eventModule, :eventFunction, :eventAction, :eventKeyword, :eventDesc, :patientPrn)
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("admin_id", adminId),
        sql.Named("eventModule", o.EventModule),
        sql.Named("eventFunction", o.EventFunction),
        sql.Named("eventAction", o.EventAction),
        sql.Named("eventKeyword", o.EventKeyword),
        sql.Named("eventDesc", o.EventDesc),
        sql.Named("patientPrn", o.PatientPrn),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func (s *AdminUserService) ChangeUserSignInType(prm dto.ChangeSignInTypeDto) error {
    var err error
    switch prm.SignInType {
    case 1:
        query := `
            UPDATE APPLICATION_USER SET 
              SIGN_IN_TYPE = :signInType,
              USERNAME = :signInMobileNumber,
              FIRST_TIME_LOGIN = 1,
              FIRST_TIME_BIOMETRIC = 1
            WHERE USER_ID = :user_id
        `
        _, err = s.db.ExecContext(s.ctx, query,
            sql.Named("signInType", prm.SignInType),
            sql.Named("signInMobileNumber", prm.SignInMobileNumber),
            sql.Named("user_id", prm.UserId),
        )
    case 2:
        hashedPwd, err := bcrypt.GenerateFromPassword([]byte(prm.SignInEmailPassword), saltRounds)
        if err != nil {
            utils.LogError(err)
            return err
        }

        query := `
            UPDATE APPLICATION_USER SET 
              SIGN_IN_TYPE = :signInType,
              USERNAME = :signInEmailAddress,
              PASSWORD = :signInEmailPassword,
              FIRST_TIME_LOGIN = 1,
              FIRST_TIME_BIOMETRIC = 1
            WHERE USER_ID = :user_id
        `
        _, err = s.db.ExecContext(s.ctx, query,
            sql.Named("signInType", prm.SignInType),
            sql.Named("signInEmailAddress", prm.SignInEmailAddress),
            sql.Named("signInEmailPassword", string(hashedPwd)),
            sql.Named("user_id", prm.UserId),
        )
    }

    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func (s *AdminUserService) ValidateCredentials(user model.AdminUser, password string) bool {
    if user.Password.String == "" {
        return false
    }
    err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
    return err == nil
}

func buildKeywordConditions(keyword string, keyword2 string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `LOWER(amu.PRN) LIKE :keyword`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        conds = append(conds, `LOWER(amu.PATIENT_NAME) LIKE :keyword2`)
        args = append(args, sql.Named("keyword2", strings.ToLower(keyword2)))
    }

    conds = append(conds, `amu.ACTION = :action`)
    args = append(args, sql.Named("action", "Delete Account"))

    return conds, args
}

func buildAdminAuditLogConditions(keyword string, keyword2 string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `LOWER(au.EMAIL) LIKE :keyword`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        conds = append(conds, `LOWER(apal.PATIENT_PRN) LIKE :keyword2`)
        args = append(args, sql.Named("keyword2", strings.ToLower(keyword2)))
    }

    return conds, args
}

func whereClause(conds []string) string {
    if len(conds) == 0 {
        return ""
    }
    return " WHERE " + strings.Join(conds, " AND ")
}

/* func getAdminUserCols() string {
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

func getAuditMobileUserCols() string {
    return `
        AUDIT_ID,
        PRN,
        USERNAME,
        PATIENT_NAME,
        ACTION,
        ACTION_DESC,
        REMARKS,
        USER_CREATE,
        DATE_CREATE
    `
} */
