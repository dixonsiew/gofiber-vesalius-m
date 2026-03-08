package applicationUser

import (
    "context"
    "database/sql"
    "fmt"
    "math/rand"
    "strings"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "golang.org/x/crypto/bcrypt"
)

type ApplicationUserService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewApplicationUserService(db *sqlx.DB, ctx context.Context) *ApplicationUserService {
    return &ApplicationUserService{db: db, ctx: ctx}
}

const saltRounds = 10

func scanApplicationUser(row *sqlx.Row) (*model.ApplicationUser, error) {
    var u model.ApplicationUser
    err := row.StructScan(&u)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &u, nil
}

func scanApplicationUsers(rows *sqlx.Rows) ([]model.ApplicationUser, error) {
    var users []model.ApplicationUser
    for rows.Next() {
        var u model.ApplicationUser
        if err := rows.StructScan(&u); err != nil {
            utils.LogError(err)
            return nil, err
        }
        users = append(users, u)
    }
    return users, nil
}

func (s *ApplicationUserService) FindAll(offset int, limit int, conn *sqlx.DB) ([]model.ApplicationUser, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM APPLICATION_USER WHERE INACTIVE_FLAG = 'N' ORDER BY REGISTRATION_DATE_TIME, MASTER_PRN OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY`
    var users []model.ApplicationUser
    err := db.SelectContext(s.ctx, &users, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for _, o := range users {
        o.Set()
    }
    return users, nil
}

func (s *ApplicationUserService) List(page string, limit string) (*model.PagedList, error) {
    total, err := s.Count(s.db)
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

func (s *ApplicationUserService) Count(conn *sqlx.DB) (int, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    var count int
    query := `SELECT COUNT(USER_ID) FROM APPLICATION_USER WHERE INACTIVE_FLAG = 'N'`
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ApplicationUserService) FindAllActive(offset int, limit int, conn *sqlx.DB) ([]model.ApplicationUser, error) {
    // same as findAll because condition already includes INACTIVE_FLAG='N'
    return s.FindAll(offset, limit, conn)
}

func (s *ApplicationUserService) ListActive(page string, limit string) (*model.PagedList, error) {
    total, err := s.CountActive(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllActive(pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *ApplicationUserService) CountActive(conn *sqlx.DB) (int, error) {
    return s.Count(conn) // same condition
}

func (s *ApplicationUserService) FindByKeyword(keyword string, offset int, limit int, conn *sqlx.DB) ([]model.ApplicationUser, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    kw := "%" + strings.ToLower(keyword) + "%"
    query := `
        SELECT * FROM APPLICATION_USER au
        WHERE (LOWER(au.FIRST_NAME) LIKE :1 OR LOWER(au.MIDDLE_NAME) LIKE :1 OR LOWER(au.LAST_NAME) LIKE :1
        OR au.MASTER_PRN LIKE :1 OR LOWER(au.EMAIL) LIKE :1)
        AND INACTIVE_FLAG = 'N'
        ORDER BY REGISTRATION_DATE_TIME, MASTER_PRN OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY
    `
    var users []model.ApplicationUser
    err := db.SelectContext(s.ctx, &users, query, kw, offset, limit)
    for _, o := range users {
        o.Set()
    }
    return users, err
}

func (s *ApplicationUserService) ListByKeyword(keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountByKeyword(keyword, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindByKeyword(keyword, pager.GetLowerBound(), pager.PageSize, s.db)
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

func (s *ApplicationUserService) CountByKeyword(keyword string, conn *sqlx.DB) (int, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    kw := "%" + strings.ToLower(keyword) + "%"
    query := `
        SELECT COUNT(au.USER_ID) FROM APPLICATION_USER au
        WHERE (LOWER(au.FIRST_NAME) LIKE :1 OR LOWER(au.MIDDLE_NAME) LIKE :1 OR LOWER(au.LAST_NAME) LIKE :1
        OR au.MASTER_PRN LIKE :1 OR LOWER(au.EMAIL) LIKE :1) AND INACTIVE_FLAG = 'N'
    `
    var count int
    err := db.GetContext(s.ctx, &count, query, kw)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ApplicationUserService) FindByUserIdSessionId(userId int64, sessionId string) (*model.ApplicationUser, error) {
    query := `SELECT * FROM APPLICATION_USER WHERE USER_ID = :1 AND SESSION_ID = :2`
    var u model.ApplicationUser
    err := s.db.GetContext(s.ctx, &u, query, userId, sessionId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    u.Set()
    return &u, err
}

func (s *ApplicationUserService) FindByUserId(userId int64, conn *sqlx.DB) (*model.ApplicationUser, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM APPLICATION_USER WHERE USER_ID = :1`
    var u model.ApplicationUser
    err := db.GetContext(s.ctx, &u, query, userId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    u.Set()
    return &u, err
}

func (s *ApplicationUserService) FindByUsername(username string, conn *sqlx.DB) (*model.ApplicationUser, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM APPLICATION_USER WHERE LOWER(USERNAME) = LOWER(:1) ORDER BY REGISTRATION_DATE_TIME DESC`
    var u model.ApplicationUser
    err := db.GetContext(s.ctx, &u, query, username)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    u.Set()
    return &u, err
}

func (s *ApplicationUserService) FindByEmail(email string, conn *sqlx.DB) (*model.ApplicationUser, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM APPLICATION_USER WHERE LOWER(EMAIL) = LOWER(:1)`
    var u model.ApplicationUser
    err := db.GetContext(s.ctx, &u, query, email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    u.Set()
    return &u, err
}

func (s *ApplicationUserService) FindByPRN(prn string, conn *sqlx.DB) (*model.ApplicationUser, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM APPLICATION_USER WHERE MASTER_PRN = :1`
    var u model.ApplicationUser
    err := db.GetContext(s.ctx, &u, query, prn)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    u.Set()
    return &u, err
}

func (s *ApplicationUserService) FindWithAssignBranchByUserId(userId int64) (*model.ApplicationUser, error) {
    query := `
        SELECT * FROM APPLICATION_USER au 
        LEFT JOIN ASSIGN_BRANCH ab ON au.USER_ID = ab.USER_ID 
        INNER JOIN BRANCH b ON b.BRANCH_ID = ab.BRANCH_ID 
        WHERE au.USER_ID = :1
    `
    rows, err := s.db.QueryxContext(s.ctx, query, userId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    var u model.ApplicationUser
    var user *model.ApplicationUser
    var branches []model.AssignBranch

    for rows.Next() {

        // We need to scan into all fields from three tables. For simplicity, we scan into a map and then construct.
        // Alternatively, we could use sqlx with embedded structs. Let's use a map for clarity.
        if user == nil {
            err = rows.StructScan(&u)
            if err != nil {
                utils.LogError(err)
                return nil, err
            }
            user = &u
            user.Set()
        }
        // Build assign branch
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

        ab.Branch = b
        branches = append(branches, ab)
    }
    if user == nil {
        return nil, nil
    }
    user.Password.String = "" // Clear password for security
    user.UserBranches = branches
    return user, nil
}

func (s *ApplicationUserService) FindWithAssignBranchByEmail(email string) (*model.ApplicationUser, error) {
    // Similar to above but with email condition
    query := `
        SELECT * FROM APPLICATION_USER au 
        LEFT JOIN ASSIGN_BRANCH ab ON au.USER_ID = ab.USER_ID 
        INNER JOIN BRANCH b ON b.BRANCH_ID = ab.BRANCH_ID 
        WHERE au.EMAIL = :1
    `
    rows, err := s.db.QueryxContext(s.ctx, query, email)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()
    
    var u model.ApplicationUser
    var user *model.ApplicationUser
    var branches []model.AssignBranch

    for rows.Next() {

        // We need to scan into all fields from three tables. For simplicity, we scan into a map and then construct.
        // Alternatively, we could use sqlx with embedded structs. Let's use a map for clarity.
        if user == nil {
            err = rows.StructScan(&u)
            if err != nil {
                utils.LogError(err)
                return nil, err
            }
            user = &u
            user.Set()
        }
        // Build assign branch
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

        ab.Branch = b
        branches = append(branches, ab)
    }
    if user == nil {
        return nil, nil
    }
    user.Password.String = "" // Clear password for security
    user.UserBranches = branches
    return user, nil
}

func (s *ApplicationUserService) FindAssignBranchByUserId(userId int64, branchId int64) (*model.AssignBranch, error) {
    query := `SELECT * FROM ASSIGN_BRANCH WHERE BRANCH_ID = :1 AND USER_ID IN (SELECT USER_ID FROM APPLICATION_USER WHERE USER_ID = :2)`
    var ab model.AssignBranch
    err := s.db.GetContext(s.ctx, &ab, query, branchId, userId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &ab, err
}

func (s *ApplicationUserService) FindAssignBranchByEmail(email string, branchId int64) (*model.AssignBranch, error) {
    query := `SELECT * FROM ASSIGN_BRANCH WHERE BRANCH_ID = :1 AND USER_ID IN (SELECT USER_ID FROM APPLICATION_USER WHERE EMAIL = :2)`
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

func (s *ApplicationUserService) FindByOtherPRN(prn string, userId int64) (*model.ApplicationUser, error) {
    // Original query: SELECT * FROM APPLICATION_USER WHERE USER_ID IN (SELECT USER_ID FROM ASSIGN_BRANCH WHERE USER_ID <> :1 AND ab.PRN = :2)
    // Note: 'ab.PRN' likely missing table alias, but we'll replicate.
    query := `SELECT * FROM APPLICATION_USER WHERE USER_ID IN (SELECT USER_ID FROM ASSIGN_BRANCH WHERE USER_ID <> :1 AND PRN = :2)`
    var u model.ApplicationUser
    err := s.db.GetContext(s.ctx, &u, query, userId, prn)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    u.Set()
    return &u, err
}

func (s *ApplicationUserService) ExistsByEmail(email string) (bool, error) {
    query := `SELECT COUNT(*) FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE LOWER(USERNAME) = LOWER(:1))`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, email)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, err
}

func (s *ApplicationUserService) ExistsByPRN(prn string) (bool, error) {
    query := `SELECT COUNT(*) FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE MASTER_PRN = :1)`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, prn)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, err
}

func (s *ApplicationUserService) ExistsByMobileNo(mobileNo string) (bool, error) {
    query := `SELECT COUNT(*) FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE USERNAME = :1)`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, mobileNo)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, err
}

func (s *ApplicationUserService) SaveUserBranch(branchId int64, o *model.ApplicationUser) error {
    tx, err := s.db.Beginx()
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

    // Update APPLICATION_USER
    updateQuery := `
        UPDATE APPLICATION_USER
        SET ADDRESS = :address, CONTACT_NUMBER = :contact_number, DOB = :dob, FIRST_NAME = :first_name, LAST_NAME = :last_name, 
        MASTER_PRN = :master_prn, MIDDLE_NAME = :middle_name, NATIONALITY = :nationality, PASSPORT = :passport, RESIDENT = :resident, 
        SEX = :sex, TITLE = :title
        WHERE USER_ID = :user_id
    `
    _, err = tx.NamedExec(updateQuery, map[string]interface{}{
        "address":        o.Address.String,
        "contact_number": o.ContactNumber.String,
        "dob":            o.Dob.String,
        "first_name":     o.FirstName.String,
        "last_name":      o.LastName.String,
        "master_prn":     o.MasterPrn.String,
        "middle_name":    o.MiddleName.String,
        "nationality":    o.Nationality.String,
        "passport":       o.Passport.String, // Note: Passport field exists? Not in struct, add if needed.
        "resident":       o.Resident.String,
        "sex":            o.Sex.String,
        "title":          o.Title.String,
        "user_id":        o.UserID.Int64,
    })
    if err != nil {
        return err
    }

    // Insert ASSIGN_BRANCH
    insertQuery := `INSERT INTO ASSIGN_BRANCH (ASSIGN_BRANCH_ID, ADMIN_ID, PRN, USER_ID, BRANCH_ID) VALUES(USER_BRANCH_SEQ.nextval, NULL, :1, :2, :3)`
    _, err = tx.ExecContext(s.ctx, insertQuery, o.MasterPrn.String, o.UserID.Int64, branchId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return tx.Commit()
}

func (s *ApplicationUserService) Update(o *model.ApplicationUser) error {
    updateQuery := `
        UPDATE APPLICATION_USER
        SET ADDRESS = :address, CONTACT_NUMBER = :contact_number, DOB = :dob, FIRST_NAME = :first_name, LAST_NAME = :last_name, 
        MASTER_PRN = :master_prn, MIDDLE_NAME = :middle_name, NATIONALITY = :nationality, PASSPORT = :passport, RESIDENT = :resident, 
        SEX = :sex, TITLE = :title
        WHERE USER_ID = :user_id
    `
    _, err := s.db.NamedExecContext(s.ctx, updateQuery, map[string]interface{}{
        "address":        o.Address.String,
        "contact_number": o.ContactNumber.String,
        "dob":            o.Dob.String,
        "first_name":     o.FirstName.String,
        "last_name":      o.LastName.String,
        "master_prn":     o.MasterPrn.String,
        "middle_name":    o.MiddleName.String,
        "nationality":    o.Nationality.String,
        "passport":       o.Passport.String,
        "resident":       o.Resident.String,
        "sex":            o.Sex.String,
        "title":          o.Title.String,
        "user_id":        o.UserID.Int64,
    })
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) SaveSessionId(userId int64, conn *sqlx.DB) (string, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    sessionID := uuid.New().String()
    query := `UPDATE APPLICATION_USER SET SESSION_ID = :1 WHERE USER_ID = :2`
    _, err := db.ExecContext(s.ctx, query, sessionID, userId)
    if err != nil {
        utils.LogError(err)
        return "", err
    }
    return sessionID, nil
}

func (s *ApplicationUserService) SetActive(userId int64) error {
    query := `UPDATE APPLICATION_USER SET INACTIVE_FLAG = 'N' WHERE USER_ID = :1`
    _, err := s.db.ExecContext(s.ctx, query, userId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) SetInactive(userId int64) error {
    query := `UPDATE APPLICATION_USER SET INACTIVE_FLAG = 'Y' WHERE USER_ID = :1`
    _, err := s.db.ExecContext(s.ctx, query, userId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) Delete(userId int64) error {
    tx, err := s.db.Beginx()
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

    _, err = tx.ExecContext(s.ctx, `DELETE FROM ASSIGN_BRANCH WHERE USER_ID = :1`, userId)
    if err != nil {
        utils.LogError(err)
        return err
    }
    _, err = tx.ExecContext(s.ctx, `DELETE FROM APPLICATION_USER WHERE USER_ID = :1`, userId)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return tx.Commit()
}

func (s *ApplicationUserService) SaveSignup(branchId int64, o *model.ApplicationUser) error {
    hashedPwd, err := bcrypt.GenerateFromPassword([]byte(o.Password.String), saltRounds)
    if err != nil {
        utils.LogError(err)
        return err
    }
    verificationCode := getRandomStr(6)

    tx, err := s.db.Beginx()
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

    // Insert APPLICATION_USER with RETURNING USER_ID
    query := `
        INSERT INTO APPLICATION_USER
        (USER_ID, ADDRESS, CONTACT_NUMBER, DOB, EMAIL, 
            FIRST_TIME_LOGIN, FIRST_NAME, LAST_NAME, MASTER_PRN, MIDDLE_NAME, 
            NATIONALITY, PASSPORT, PASSWORD, RESIDENT, "ROLE", 
            SEX, TITLE, USERNAME, VERIFICATION_CODE, BRANCH, RACE, PLAYER_ID)
        VALUES(APP_USER_SEQ.nextval, :address, :contact_number, :dob, :email, 
            1, :first_name, :last_name, :master_prn, :middle_name, 
            :nationality, :passport, :password, :resident, :role, 
            :sex, :title, :username, :verification_code, :branch, :race, :player_id)
        RETURNING USER_ID INTO :user_id
    `
    params := map[string]interface{}{
        "address":           o.Address.String,
        "contact_number":    o.ContactNumber.String,
        "dob":               o.Dob.String,
        "email":             o.Email.String,
        "first_name":        o.FirstName.String,
        "last_name":         o.LastName.String,
        "master_prn":        o.MasterPrn.String,
        "middle_name":       o.MiddleName.String,
        "nationality":       o.Nationality.String,
        "passport":          o.Passport.String,
        "password":          string(hashedPwd),
        "resident":          o.Resident.String,
        "role":              o.Role.String,
        "sex":               o.Sex.String,
        "title":             o.Title.String,
        "username":          o.Username.String,
        "verification_code": verificationCode,
        "branch":            nil,
        "race":              o.Race.String,
        "player_id":         o.PlayerID.String,
        "user_id":           sql.Out{Dest: &o.UserID.Int64},
    }
    _, err = tx.NamedExecContext(s.ctx, query, params)
    if err != nil {
        utils.LogError(err)
        return err
    }

    // Insert into ASSIGN_BRANCH
    _, err = tx.ExecContext(s.ctx, `INSERT INTO ASSIGN_BRANCH (ASSIGN_BRANCH_ID, ADMIN_ID, PRN, USER_ID, BRANCH_ID) VALUES(USER_BRANCH_SEQ.nextval, NULL, :1, :2, :3)`,
        o.MasterPrn.String, o.UserID.Int64, branchId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return tx.Commit()
}

func (s *ApplicationUserService) UpdateInactiveSignup(o *model.ApplicationUser) error {
    var hashedPwd []byte
    var err error
    if o.SignInType.Valid && o.SignInType.Int32 == 2 && o.Password.Valid && o.Password.String != "" {
        hashedPwd, err = bcrypt.GenerateFromPassword([]byte(o.Password.String), saltRounds)
        if err != nil {
            return err
        }
    }
    verificationCode := getRandomStr(6)

    // Convert bool to Y/N
    isGoldenPearl := "N"
    if o.IsGoldenPearl.Valid && o.IsGoldenPearl.String == "Y" {
        isGoldenPearl = "Y"
    }
    isKidsExplorer := "N"
    if o.IsKidsExplorer.Valid && o.IsKidsExplorer.String == "Y" {
        isKidsExplorer = "Y"
    }
    firstTimeLogin := 0
    if o.FirstTimeLoginV.Valid && o.FirstTimeLoginV.Int32 == 1 {
        firstTimeLogin = 1
    }
    firstTimeBiometric := 0
    if o.FirstTimeBiometricV.Valid && o.FirstTimeBiometricV.Int32 == 1 {
        firstTimeBiometric = 1
    }

    tx, err := s.db.Beginx()
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

    // Update APPLICATION_USER
    updateQuery := `
        UPDATE APPLICATION_USER SET
            ADDRESS = :address,
            ADDRESS_1 = :address1,
            ADDRESS_2 = :address2,
            ADDRESS_3 = :address3,
            NATIONALITY = :nationality,
            RACE = :race,
            SEX = :sex,
            TITLE = :title,
            CONTACT_NUMBER = :contact_number,
            DOB = :dob,
            EMAIL = :email, 
            MASTER_PRN = :master_prn, 
            FIRST_NAME = :first_name, 
            MIDDLE_NAME = :middle_name, 
            LAST_NAME = :last_name,
            IS_GOLDEN_PEARL = :isGoldenPearl,
            IS_KIDS_EXPLORER = :isKidsExplorer,
            PASSWORD = :password, 
            RESIDENT = :resident, 
            ROLE = :role, 
            USERNAME = :username, 
            VERIFICATION_CODE = :verification_code,
            REGISTRATION_DATE_TIME = CURRENT_TIMESTAMP,
            INACTIVE_FLAG = :inactive,
            FIRST_TIME_LOGIN = :firstTimeLogin,
            FIRST_TIME_BIOMETRIC = :firstTimeBiometric,
            PLAYER_ID = :player_id, 
            CITYSTATE = :cityState,
            POSTCODE = :postalCode,
            COUNTRY = :country,
            SIGN_IN_TYPE = :signInType, 
            FULLNAME_SIGNUP = :fullNameSignUp, 
            DOC_NO_SIGNUP = :docNoSignUp
        WHERE USER_ID = :user_id
    `
    params := map[string]interface{}{
        "address":            o.Address.String,
        "address1":           o.Address1.String,
        "address2":           o.Address2.String,
        "address3":           o.Address3.String,
        "nationality":        o.Nationality.String,
        "race":               o.Race.String,
        "sex":                o.Sex.String,
        "title":              o.Title.String,
        "contact_number":     o.ContactNumber.String,
        "dob":                o.Dob.String,
        "email":              o.Email.String,
        "master_prn":         o.MasterPrn.String,
        "first_name":         o.FirstName.String,
        "middle_name":        o.MiddleName.String,
        "last_name":          o.LastName.String,
        "isGoldenPearl":      isGoldenPearl,
        "isKidsExplorer":     isKidsExplorer,
        "password":           string(hashedPwd),
        "resident":           o.Resident.String,
        "role":               o.Role.String,
        "username":           o.Username.String,
        "verification_code":  verificationCode,
        "inactive":           o.InactiveFlag.String,
        "firstTimeLogin":     firstTimeLogin,
        "firstTimeBiometric": firstTimeBiometric,
        "player_id":          o.PlayerID.String,
        "cityState":          o.CityState.String,
        "postalCode":         o.Postcode.String,
        "country":            o.Country.String,
        "signInType":         o.SignInType.Int32,
        "fullNameSignUp":     o.FullnameSignup.String,
        "docNoSignUp":        o.DocNoSignup.String,
        "user_id":            o.UserID.Int64,
    }
    _, err = tx.NamedExecContext(s.ctx, updateQuery, params)
    if err != nil {
        utils.LogError(err)
        return err
    }

    // Update ASSIGN_BRANCH
    _, err = tx.ExecContext(s.ctx, `UPDATE ASSIGN_BRANCH SET PRN = :1 WHERE USER_ID = :2`, o.MasterPrn.String, o.UserID.Int64)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return tx.Commit()
}

func (s *ApplicationUserService) SaveNewSignup(branchId int64, o *model.ApplicationUser) (int64, error) {
    var hashedPwd []byte
    var err error
    if o.SignInType.Valid && o.SignInType.Int32 == 2 && o.Password.Valid && o.Password.String != "" {
        hashedPwd, err = bcrypt.GenerateFromPassword([]byte(o.Password.String), saltRounds)
        if err != nil {
            return 0, err
        }
    }
    verificationCode := getRandomStr(6)
    firstTimeLogin := 0
    if o.FirstTimeLoginV.Valid && o.FirstTimeLoginV.Int32 == 1 {
        firstTimeLogin = 1
    }

    tx, err := s.db.Beginx()
    if err != nil {
        utils.LogError(err)
        return -1, err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()

    query := `
        INSERT INTO APPLICATION_USER
        (
            USER_ID, ADDRESS, ADDRESS_1, ADDRESS_2, ADDRESS_3, CITYSTATE,
            POSTCODE, COUNTRY, NATIONALITY, RACE, SEX, TITLE, 
            CONTACT_NUMBER, DOB, EMAIL, MASTER_PRN, FIRST_TIME_LOGIN, FIRST_NAME, MIDDLE_NAME, 
            LAST_NAME, PASSWORD, RESIDENT, ROLE, USERNAME, VERIFICATION_CODE, 
            PLAYER_ID, SIGN_IN_TYPE, FULLNAME_SIGNUP, DOC_NO_SIGNUP
        ) VALUES (
            APP_USER_SEQ.nextval, :address, :address1, :address2, :address3, :cityState,
            :postalCode, :country, :nationality, :race, :sex, :title, 
            :contact_number, :dob, :email, :master_prn, :firstTimeLogin, :first_name, :middle_name, 
            :last_name, :password, :resident, :role, :username, :verification_code, 
            :player_id, :signInType, :fullNameSignUp, :docNoSignUp
        ) RETURNING USER_ID INTO :user_id
    `
    params := map[string]interface{}{
        "address":           o.Address.String,
        "address1":          o.Address1.String,
        "address2":          o.Address2.String,
        "address3":          o.Address3.String,
        "cityState":         o.CityState.String,
        "postalCode":        o.Postcode.String,
        "country":           o.Country.String,
        "nationality":       o.Nationality.String,
        "race":              o.Race.String,
        "sex":               o.Sex.String,
        "title":             o.Title.String,
        "contact_number":    o.ContactNumber.String,
        "dob":               o.Dob.String,
        "email":             o.Email.String,
        "master_prn":        o.MasterPrn.String,
        "firstTimeLogin":    firstTimeLogin,
        "first_name":        o.FirstName.String,
        "middle_name":       o.MiddleName.String,
        "last_name":         o.LastName.String,
        "password":          string(hashedPwd),
        "resident":          o.Resident.String,
        "role":              o.Role.String,
        "username":          o.Username.String,
        "verification_code": verificationCode,
        "player_id":         o.PlayerID.String,
        "signInType":        o.SignInType.Int32,
        "fullNameSignUp":    o.FullnameSignup.String,
        "docNoSignUp":       o.DocNoSignup.String,
        "user_id":           sql.Out{Dest: &o.UserID.Int64},
    }
    _, err = tx.NamedExecContext(s.ctx, query, params)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }

    _, err = tx.ExecContext(s.ctx, `INSERT INTO ASSIGN_BRANCH (ASSIGN_BRANCH_ID, ADMIN_ID, PRN, USER_ID, BRANCH_ID) VALUES(USER_BRANCH_SEQ.nextval, NULL, :1, :2, :3)`,
        o.MasterPrn.String, o.UserID.Int64, branchId)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }

    err = tx.Commit()
    if err != nil {
        utils.LogError(err)
        return -1, err
    }
    return o.UserID.Int64, err
}

func (s *ApplicationUserService) SaveResetPassword(o *model.ApplicationUser) error {
    newPwd := getRandomStr(6)
    hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPwd), saltRounds)
    if err != nil {
        utils.LogError(err)
        return err
    }
    query := `UPDATE APPLICATION_USER SET PASSWORD = :1 WHERE USER_ID = :2`
    _, err = s.db.ExecContext(s.ctx, query, string(hashedPwd), o.UserID.Int64)
    if err != nil {
        utils.LogError(err)
        return err
    }
    // Store new plain password in o.Password? The original sets o.password = newPwd but doesn't save it.
    // We'll just return the new password via o.Password for later use.
    o.Password.String = newPwd
    return err
}

func (s *ApplicationUserService) SavePassword(o *model.ApplicationUser) error {
    hashedPwd, err := bcrypt.GenerateFromPassword([]byte(o.Password.String), saltRounds)
    if err != nil {
        return err
    }
    query := `UPDATE APPLICATION_USER SET PASSWORD = :1 WHERE USER_ID = :2`
    _, err = s.db.ExecContext(s.ctx, query, string(hashedPwd), o.UserID.Int64)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) GenerateVerificationCode(o *model.ApplicationUser) error {
    code := getRandomStr(6)
    o.VerificationCode.String = code
    query := `UPDATE APPLICATION_USER SET VERIFICATION_CODE = :1 WHERE USER_ID = :2`
    _, err := s.db.ExecContext(s.ctx, query, code, o.UserID.Int64)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) UpdateVerificationCode(code string, userId int64) error {
    query := `UPDATE APPLICATION_USER SET VERIFICATION_CODE = :1 WHERE USER_ID = :2`
    _, err := s.db.ExecContext(s.ctx, query, code, userId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) UpdateMachineId(id string, userId int64, conn *sqlx.DB) error {
    db := conn
    if db == nil {
        db = s.db
    }
    hashedID, err := bcrypt.GenerateFromPassword([]byte(id), saltRounds)
    if err != nil {
        utils.LogError(err)
        return err
    }
    query := `UPDATE APPLICATION_USER SET MACHINE_ID = :1 WHERE USER_ID = :2`
    _, err = db.ExecContext(s.ctx, query, string(hashedID), userId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) UpdatePlayerId(id string, userId int64, conn *sqlx.DB) error {
    db := conn
    if db == nil {
        db = s.db
    }
    // First nullify any existing player_id with same id
    _, err := db.ExecContext(s.ctx, `UPDATE APPLICATION_USER SET PLAYER_ID = NULL WHERE PLAYER_ID = :1`, id)
    if err != nil {
        utils.LogError(err)
        return err
    }
    // Then set new player_id
    _, err = db.ExecContext(s.ctx, `UPDATE APPLICATION_USER SET PLAYER_ID = :1 WHERE USER_ID = :2`, id, userId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) InsertDownloadApp(playerId string, conn *sqlx.DB) error {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `
        MERGE INTO APP_DOWNLOADED_USER apu
        USING (SELECT :1 AS PLAYER_ID FROM DUAL) src
        ON (apu.PLAYER_ID = src.PLAYER_ID)
        WHEN NOT MATCHED THEN
        INSERT (PLAYER_ID) VALUES (src.PLAYER_ID)
    `
    _, err := db.ExecContext(s.ctx, query, playerId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) InsertDownloadAppV2(machineId string, playerId string, conn *sqlx.DB) error {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `
        MERGE INTO APP_DOWNLOADED_USER apu
        USING (SELECT :1 AS MACHINE_ID, :2 AS PLAYER_ID FROM DUAL) src
        ON (apu.MACHINE_ID = src.MACHINE_ID)
        WHEN MATCHED THEN
            UPDATE SET apu.PLAYER_ID = src.PLAYER_ID, DATE_UPDATE = CURRENT_TIMESTAMP
        WHEN NOT MATCHED THEN
            INSERT (MACHINE_ID, PLAYER_ID, DATE_UPDATE) VALUES (src.MACHINE_ID, src.PLAYER_ID, CURRENT_TIMESTAMP)
    `
    _, err := db.ExecContext(s.ctx, query, machineId, playerId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) VerifyUserSms(o *model.ApplicationUser) (bool, error) {
    tx, err := s.db.Beginx()
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()

    _, err = tx.ExecContext(s.ctx, `UPDATE APPLICATION_USER SET FIRST_TIME_LOGIN = 0, FIRST_TIME_BIOMETRIC = 0 WHERE USER_ID = :1`, o.UserID.Int64)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    sessionID, err := s.SaveSessionId(o.UserID.Int64, s.db)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    o.SessionID.String = sessionID
    err = tx.Commit()
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return err == nil, err
}

func (s *ApplicationUserService) VerifyUser(o *model.ApplicationUser) (bool, error) {
    query := `UPDATE APPLICATION_USER SET FIRST_TIME_LOGIN = 0 WHERE USER_ID = :1`
    _, err := s.db.ExecContext(s.ctx, query, o.UserID.Int64)
    if err != nil {
        utils.LogError(err)
    }
    return err == nil, err
}

func (s *ApplicationUserService) DeleteUserAccount(user *model.ApplicationUser, admin *model.AdminUser) error {
    tx, err := s.db.Beginx()
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

    updateUserQuery := `
        UPDATE APPLICATION_USER SET 
            INACTIVE_FLAG = 'Y',
            ADDRESS = NULL,
            ADDRESS_1 = NULL,
            ADDRESS_2 = NULL,
            ADDRESS_3 = NULL,
            NATIONALITY = NULL,
            RACE = NULL,
            SEX = NULL,
            TITLE = NULL,
            CONTACT_NUMBER = NULL,
            DOB = NULL,
            EMAIL = NULL,
            MASTER_PRN = NULL,
            FIRST_NAME = NULL,
            MIDDLE_NAME = NULL,
            LAST_NAME = NULL,
            IS_GOLDEN_PEARL = NULL,
            IS_KIDS_EXPLORER = NULL,
            PASSPORT = NULL,
            RESIDENT = NULL,
            ROLE = NULL,
            VERIFICATION_CODE = NULL,
            BRANCH = NULL,
            REGISTRATION_DATE_TIME = NULL,
            FIRST_TIME_LOGIN = NULL,
            FIRST_TIME_BIOMETRIC = NULL,
            IS_LOGGED_IN = NULL,
            DATE_LOGGED_IN = NULL,
            SESSION_ID = NULL,
            CITYSTATE = NULL,
            POSTCODE = NULL,
            COUNTRY = NULL,
            FULLNAME_SIGNUP = NULL,
            DOC_NO_SIGNUP = NULL
        WHERE USER_ID = :1
    `
    _, err = tx.ExecContext(s.ctx, updateUserQuery, user.UserID.Int64)
    if err != nil {
        utils.LogError(err)
        return err
    }

    _, err = tx.ExecContext(s.ctx, `UPDATE ASSIGN_BRANCH SET PRN = NULL WHERE USER_ID = :1`, user.UserID.Int64)
    if err != nil {
        utils.LogError(err)
        return err
    }

    // Insert audit
    patientName := strings.TrimSpace(fmt.Sprintf("%s %s %s",
        user.FirstName.String,
        user.MiddleName.String,
        user.LastName.String))
    actionDesc := fmt.Sprintf("Mobile account: (%s) has been deleted", user.Username.String)
    remarks := "Via Mobile App"
    userCreate := user.Username.String
    if admin != nil {
        remarks = "Via Admin Portal"
        userCreate = admin.Email.String
    }
    auditQuery := `
        INSERT INTO AUDIT_MOBILE_USER 
        (PRN, USERNAME, PATIENT_NAME, ACTION, ACTION_DESC, REMARKS, USER_CREATE)
        VALUES 
        (:1, :2, :3, :4, :5, :6, :7)
    `
    _, err = tx.Exec(auditQuery,
        user.MasterPrn.String,
        user.Username.String,
        patientName,
        "Delete Account",
        actionDesc,
        remarks,
        userCreate)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return tx.Commit()
}

func (s *ApplicationUserService) DisableFirstTimeBiometricUser(userId int64) error {
    query := `UPDATE APPLICATION_USER SET FIRST_TIME_BIOMETRIC = 0 WHERE USER_ID = :1`
    _, err := s.db.ExecContext(s.ctx, query, userId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) ResetUserSignup(userId int64, prn string) error {
    tx, err := s.db.Beginx()
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

    _, err = tx.ExecContext(s.ctx, `DELETE FROM ASSIGN_BRANCH WHERE USER_ID = :1 AND PRN = :2`, userId, prn)
    if err != nil {
        utils.LogError(err)
        return err
    }
    _, err = tx.ExecContext(s.ctx, `DELETE FROM APPLICATION_USER WHERE USER_ID = :1 AND MASTER_PRN = :2`, userId, prn)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return tx.Commit()
}

func (s *ApplicationUserService) SetLogin(userId int64) error {
    query := `UPDATE APPLICATION_USER SET IS_LOGGED_IN = 1, DATE_LOGGED_IN = CURRENT_TIMESTAMP WHERE USER_ID = :1`
    _, err := s.db.ExecContext(s.ctx, query, userId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) SetLogout(userId int64) error {
    query := `UPDATE APPLICATION_USER SET IS_LOGGED_IN = 0 WHERE USER_ID = :1`
    _, err := s.db.ExecContext(s.ctx, query, userId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func getRandomStr(length int) string {
    const charset = "0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}

func (s *ApplicationUserService) ValidateCredentials(user *model.ApplicationUser, password string) (bool, error) {
    if user.Password.String == "" {
        return false, nil
    }
    err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
    if err != nil {
        utils.LogError(err)
    }
    return err == nil, nil
}

func (s *ApplicationUserService) ValidateCredentials2(user *model.ApplicationUser, password string) (bool, error) {
    if user.MachineID.String == "" {
        return false, nil
    }
    err := bcrypt.CompareHashAndPassword([]byte(user.MachineID.String), []byte(password))
    if err != nil {
        utils.LogError(err)
    }
    return err == nil, nil
}
