package applicationUser

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
    "vesaliusm/database"
	"vesaliusm/model"
    assignBranchService "vesaliusm/service/assignBranch"
    branchService "vesaliusm/service/branch"
	"vesaliusm/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	go_ora "github.com/sijms/go-ora/v2"
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

var (
    assignBranchSvc *assignBranchService.AssignBranchService = 
        assignBranchService.NewAssignBranchService(database.GetDb(), database.GetCtx())
    branchSvc *branchService.BranchService = 
        branchService.NewBranchService(database.GetDb(), database.GetCtx())
)

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
    query := `SELECT ` + utils.GetDbCols(model.ApplicationUser{}, "") + 
        ` FROM APPLICATION_USER WHERE INACTIVE_FLAG = 'N' ORDER BY REGISTRATION_DATE_TIME, MASTER_PRN OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    list := make([]model.ApplicationUser, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
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
    query := `
        SELECT ` + utils.GetDbCols(model.ApplicationUser{}, "au.") + 
        ` FROM APPLICATION_USER au
        WHERE (LOWER(au.FIRST_NAME) LIKE :key OR LOWER(au.MIDDLE_NAME) LIKE :key OR LOWER(au.LAST_NAME) LIKE :key
        OR au.MASTER_PRN LIKE :key OR LOWER(au.EMAIL) LIKE :key)
        AND INACTIVE_FLAG = 'N'
        ORDER BY REGISTRATION_DATE_TIME, MASTER_PRN OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    users := make([]model.ApplicationUser, 0)
    err := db.SelectContext(s.ctx, &users, query, 
        sql.Named("key", strings.ToLower(keyword)),
        sql.Named("offset", offset),
        sql.Named("limit", limit),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range users {
        users[i].Set()
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
    query := `
        SELECT COUNT(au.USER_ID) FROM APPLICATION_USER au
        WHERE (LOWER(au.FIRST_NAME) LIKE :key OR LOWER(au.MIDDLE_NAME) LIKE :key OR LOWER(au.LAST_NAME) LIKE :key
        OR au.MASTER_PRN LIKE :key OR LOWER(au.EMAIL) LIKE :key) AND INACTIVE_FLAG = 'N'
    `
    var count int
    err := db.GetContext(s.ctx, &count, query, 
        sql.Named("key", strings.ToLower(keyword)),
    )
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ApplicationUserService) FindByUserIdSessionId(userId int64, sessionId string) (*model.ApplicationUser, error) {
    query := `SELECT ` + utils.GetDbCols(model.ApplicationUser{}, "") + ` FROM APPLICATION_USER WHERE USER_ID = :userId AND SESSION_ID = :sessionId`
    var o model.ApplicationUser
    err := s.db.GetContext(s.ctx, &o, query, userId, sessionId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    return &o, err
}

func (s *ApplicationUserService) FindByUserId(userId int64, conn *sqlx.DB) (*model.ApplicationUser, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT ` + utils.GetDbCols(model.ApplicationUser{}, "") + ` FROM APPLICATION_USER WHERE USER_ID = :userId`
    var o model.ApplicationUser
    err := db.GetContext(s.ctx, &o, query, userId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    return &o, err
}

func (s *ApplicationUserService) FindByUsername(username string, conn *sqlx.DB) (*model.ApplicationUser, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT ` + utils.GetDbCols(model.ApplicationUser{}, "") + ` FROM APPLICATION_USER WHERE LOWER(USERNAME) = LOWER(:username) ORDER BY REGISTRATION_DATE_TIME DESC`
    var o model.ApplicationUser
    err := db.GetContext(s.ctx, &o, query, username)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    return &o, err
}

func (s *ApplicationUserService) FindByEmail(email string, conn *sqlx.DB) (*model.ApplicationUser, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT ` + utils.GetDbCols(model.ApplicationUser{}, "") + ` FROM APPLICATION_USER WHERE LOWER(EMAIL) = LOWER(:email)`
    var o model.ApplicationUser
    err := db.GetContext(s.ctx, &o, query, email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    return &o, err
}

func (s *ApplicationUserService) FindByPRN(prn string, conn *sqlx.DB) (*model.ApplicationUser, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT ` + utils.GetDbCols(model.ApplicationUser{}, "") + ` FROM APPLICATION_USER WHERE MASTER_PRN = :prn`
    var o model.ApplicationUser
    err := db.GetContext(s.ctx, &o, query, prn)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    return &o, err
}

func (s *ApplicationUserService) FindWithAssignBranchByUserId(userId int64) (*model.ApplicationUser, error) {
    o, err := s.FindByUserId(userId, s.db)
    if err != nil {
        return nil, err
    }

    ablist, err := assignBranchSvc.FindAllByUserId(userId)
    if err != nil {
        return nil, err
    }

    for i := range ablist {
        b, err := branchSvc.FindByBranchId(ablist[i].BranchID.Int64)
        if err != nil {
            return nil, err
        }
        b.Passcode.String = ""
        b.Url.String = ""
        ablist[i].Branch = b
    }
    o.Password.String = ""
    o.UserBranches = ablist
    return o, nil
}

func (s *ApplicationUserService) FindWithAssignBranchByEmail(email string) (*model.ApplicationUser, error) {
    o, err := s.FindByEmail(email, s.db)
    if err != nil {
        return nil, err
    }

    ablist, err := assignBranchSvc.FindAllByUserId(o.UserID.Int64)
    if err != nil {
        return nil, err
    }

    for i := range ablist {
        b, err := branchSvc.FindByBranchId(ablist[i].BranchID.Int64)
        if err != nil {
            return nil, err
        }
        b.Passcode.String = ""
        b.Url.String = ""
        ablist[i].Branch = b
    }
    o.Password.String = ""
    o.UserBranches = ablist
    return o, nil
}

func (s *ApplicationUserService) FindAssignBranchByUserId(userId int64, branchId int64) (*model.AssignBranch, error) {
    query := `SELECT ` + utils.GetDbCols(model.AssignBranch{}, "") + 
        ` FROM ASSIGN_BRANCH WHERE BRANCH_ID = :branchId AND USER_ID IN (SELECT USER_ID FROM APPLICATION_USER WHERE USER_ID = :userId)`
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
    query := `SELECT ` + utils.GetDbCols(model.AssignBranch{}, "") + 
        ` FROM ASSIGN_BRANCH WHERE BRANCH_ID = :branchId AND USER_ID IN (SELECT USER_ID FROM APPLICATION_USER WHERE EMAIL = :email)`
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
    query := `SELECT ` + utils.GetDbCols(model.ApplicationUser{}, "") + ` FROM APPLICATION_USER WHERE USER_ID IN (SELECT USER_ID FROM ASSIGN_BRANCH WHERE USER_ID <> :userId AND PRN = :prn)`
    var o model.ApplicationUser
    err := s.db.GetContext(s.ctx, &o, query, userId, prn)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    return &o, err
}

func (s *ApplicationUserService) ExistsByEmail(email string) (bool, error) {
    query := `SELECT COUNT(*) FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE LOWER(USERNAME) = LOWER(:email))`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, email)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, err
}

func (s *ApplicationUserService) ExistsByPRN(prn string) (bool, error) {
    query := `SELECT COUNT(*) FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE MASTER_PRN = :prn)`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, prn)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, err
}

func (s *ApplicationUserService) ExistsByMobileNo(mobileNo string) (bool, error) {
    query := `SELECT COUNT(*) FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE USERNAME = :mobileNo)`
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
    _, err = tx.ExecContext(s.ctx, updateQuery,
        sql.Named("address", o.Address.String),
        sql.Named("contact_number", o.ContactNumber.String),
        sql.Named("dob", o.Dob.String),
        sql.Named("first_name", o.FirstName.String),
        sql.Named("last_name", o.LastName.String),
        sql.Named("master_prn", o.MasterPrn.String),
        sql.Named("middle_name", o.MiddleName.String),
        sql.Named("nationality", o.Nationality.String),
        sql.Named("passport", o.Passport.String),
        sql.Named("resident", o.Resident.String),
        sql.Named("sex", o.Sex.String),
        sql.Named("title", o.Title.String),
        sql.Named("user_id", o.UserID.Int64),
    )
    if err != nil {
        return err
    }

    // Insert ASSIGN_BRANCH
    insertQuery := `INSERT INTO ASSIGN_BRANCH (ASSIGN_BRANCH_ID, ADMIN_ID, PRN, USER_ID, BRANCH_ID) VALUES(USER_BRANCH_SEQ.nextval, NULL, :prn, :userId, :branchId)`
    _, err = tx.ExecContext(s.ctx, insertQuery, 
        sql.Named("prn", o.MasterPrn.String),
        sql.Named("userId", o.UserID.Int64),
        sql.Named("branchId", branchId),
    )
    if err != nil {
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
    _, err := s.db.ExecContext(s.ctx, updateQuery,
        sql.Named("address", o.Address.String),
        sql.Named("contact_number", o.ContactNumber.String),
        sql.Named("dob", o.Dob.String),
        sql.Named("first_name", o.FirstName.String),
        sql.Named("last_name", o.LastName.String),
        sql.Named("master_prn", o.MasterPrn.String),
        sql.Named("middle_name", o.MiddleName.String),
        sql.Named("nationality", o.Nationality.String),
        sql.Named("passport", o.Passport.String),
        sql.Named("resident", o.Resident.String),
        sql.Named("sex", o.Sex.String),
        sql.Named("title", o.Title.String),
        sql.Named("user_id", o.UserID.Int64),
    )
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
    sessionId := uuid.New().String()
    query := `UPDATE APPLICATION_USER SET SESSION_ID = :sessionId WHERE USER_ID = :userId`
    _, err := db.ExecContext(s.ctx, query, sessionId, userId)
    if err != nil {
        utils.LogError(err)
        return "", err
    }
    return sessionId, nil
}

func (s *ApplicationUserService) SetActive(userId int64) error {
    query := `UPDATE APPLICATION_USER SET INACTIVE_FLAG = 'N' WHERE USER_ID = :userId`
    _, err := s.db.ExecContext(s.ctx, query, userId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) SetInactive(userId int64) error {
    query := `UPDATE APPLICATION_USER SET INACTIVE_FLAG = 'Y' WHERE USER_ID = :userId`
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

    _, err = tx.ExecContext(s.ctx, `DELETE FROM ASSIGN_BRANCH WHERE USER_ID = :userId`, userId)
    if err != nil {
        return err
    }
    _, err = tx.ExecContext(s.ctx, `DELETE FROM APPLICATION_USER WHERE USER_ID = :userId`, userId)
    if err != nil {
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
    var userId go_ora.Number
    _, err = tx.ExecContext(s.ctx, query,
        sql.Named("address", o.Address.String),
        sql.Named("contact_number", o.ContactNumber.String),
        sql.Named("dob", o.Dob.String),
        sql.Named("email", o.Email.String),
        sql.Named("first_name", o.FirstName.String),
        sql.Named("last_name", o.LastName.String),
        sql.Named("master_prn", o.MasterPrn.String),
        sql.Named("middle_name", o.MiddleName.String),
        sql.Named("nationality", o.Nationality.String),
        sql.Named("passport", o.Passport.String),
        sql.Named("password", string(hashedPwd)),
        sql.Named("resident", o.Resident.String),
        sql.Named("role", o.Role.String),
        sql.Named("sex", o.Sex.String),
        sql.Named("title", o.Title.String),
        sql.Named("username", o.Username.String),
        sql.Named("verification_code", verificationCode),
        nil,
        sql.Named("race", o.Race.String),
        sql.Named("player_id", o.PlayerID.String),
        go_ora.Out{Dest: &userId},
    )
    if err != nil {
        return err
    }

    o.UserID.Int64, _ = userId.Int64()

    // Insert into ASSIGN_BRANCH
    _, err = tx.ExecContext(s.ctx, `INSERT INTO ASSIGN_BRANCH (ASSIGN_BRANCH_ID, ADMIN_ID, PRN, USER_ID, BRANCH_ID) VALUES(USER_BRANCH_SEQ.nextval, NULL, :prn, :userId, :branchId)`,
        sql.Named("prn", o.MasterPrn.String),
        sql.Named("userId", o.UserID.Int64),
        sql.Named("branchId", branchId),
    )
    if err != nil {
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
            utils.LogError(err)
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
    _, err = tx.ExecContext(s.ctx, updateQuery,
        sql.Named("address", o.Address.String),
        sql.Named("address1", o.Address1.String),
        sql.Named("address2", o.Address2.String),
        sql.Named("address3", o.Address3.String),
        sql.Named("nationality", o.Nationality.String),
        sql.Named("race", o.Race.String),
        sql.Named("sex", o.Sex.String),
        sql.Named("title", o.Title.String),
        sql.Named("contact_number", o.ContactNumber.String),
        sql.Named("dob", o.Dob.String),
        sql.Named("email", o.Email.String),
        sql.Named("master_prn", o.MasterPrn.String),
        sql.Named("first_name", o.FirstName.String),
        sql.Named("middle_name", o.MiddleName.String),
        sql.Named("last_name", o.LastName.String),
        sql.Named("isGoldenPearl", isGoldenPearl),
        sql.Named("isKidsExplorer", isKidsExplorer),
        sql.Named("password", string(hashedPwd)),
        sql.Named("resident", o.Resident.String),
        sql.Named("role", o.Role.String),
        sql.Named("username", o.Username.String),
        sql.Named("verification_code", verificationCode),
        sql.Named("inactive", o.InactiveFlag.String),
        sql.Named("firstTimeLogin", firstTimeLogin),
        sql.Named("firstTimeBiometric", firstTimeBiometric),
        sql.Named("player_id", o.PlayerID.String),
        sql.Named("cityState", o.CityState.String),
        sql.Named("postalCode", o.Postcode.String),
        sql.Named("country", o.Country.String),
        sql.Named("signInType", o.SignInType.Int32),
        sql.Named("fullNameSignUp", o.FullnameSignup.String),
        sql.Named("docNoSignUp", o.DocNoSignup.String),
        sql.Named("user_id", o.UserID.Int64),
    )
    if err != nil {
        return err
    }

    // Update ASSIGN_BRANCH
    _, err = tx.ExecContext(s.ctx, `UPDATE ASSIGN_BRANCH SET PRN = :prn WHERE USER_ID = :userId`, o.MasterPrn.String, o.UserID.Int64)
    if err != nil {
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
    var userId go_ora.Number
    _, err = tx.ExecContext(s.ctx, query,
        sql.Named("address", o.Address.String),
        sql.Named("address1", o.Address1.String),
        sql.Named("address2", o.Address2.String),
        sql.Named("address3", o.Address3.String),
        sql.Named("cityState", o.CityState.String),
        sql.Named("postalCode", o.Postcode.String),
        sql.Named("country", o.Country.String),
        sql.Named("nationality", o.Nationality.String),
        sql.Named("race", o.Race.String),
        sql.Named("sex", o.Sex.String),
        sql.Named("title", o.Title.String),
        sql.Named("contact_number", o.ContactNumber.String),
        sql.Named("dob", o.Dob.String),
        sql.Named("email", o.Email.String),
        sql.Named("master_prn", o.MasterPrn.String),
        sql.Named("firstTimeLogin", firstTimeLogin),
        sql.Named("first_name", o.FirstName.String),
        sql.Named("middle_name", o.MiddleName.String),
        sql.Named("last_name", o.LastName.String),
        sql.Named("password", string(hashedPwd)),
        sql.Named("resident", o.Resident.String),
        sql.Named("role", o.Role.String),
        sql.Named("username", o.Username.String),
        sql.Named("verification_code", verificationCode),
        sql.Named("player_id", o.PlayerID.String),
        sql.Named("signInType", o.SignInType.Int32),
        sql.Named("fullNameSignUp", o.FullnameSignup.String),
        sql.Named("docNoSignUp", o.DocNoSignup.String),
        go_ora.Out{Dest: &userId},
    )
    if err != nil {
        return 0, err
    }

    o.UserID.Int64, _ = userId.Int64()

    _, err = tx.ExecContext(s.ctx, `INSERT INTO ASSIGN_BRANCH (ASSIGN_BRANCH_ID, ADMIN_ID, PRN, USER_ID, BRANCH_ID) VALUES(USER_BRANCH_SEQ.nextval, NULL, :prn, :userId, :branchId)`,
        sql.Named("prn", o.MasterPrn.String),
        sql.Named("userId", o.UserID.Int64),
        sql.Named("branchId", branchId),
    )
    if err != nil {
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
    query := `UPDATE APPLICATION_USER SET PASSWORD = :pw WHERE USER_ID = :userId`
    _, err = s.db.ExecContext(s.ctx, query, 
        sql.Named("pw", string(hashedPwd)),
        sql.Named("userId", o.UserID.Int64),
    )
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
    query := `UPDATE APPLICATION_USER SET PASSWORD = :pw WHERE USER_ID = :userId`
    _, err = s.db.ExecContext(s.ctx, query, 
        sql.Named("pw", string(hashedPwd)),
        sql.Named("userId", o.UserID.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) GenerateVerificationCode(o *model.ApplicationUser) error {
    code := getRandomStr(6)
    o.VerificationCode.String = code
    query := `UPDATE APPLICATION_USER SET VERIFICATION_CODE = :verificationCode WHERE USER_ID = :userId`
    _, err := s.db.ExecContext(s.ctx, query, 
        sql.Named("verificationCode", code),
        sql.Named("userId", o.UserID.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) UpdateVerificationCode(code string, userId int64) error {
    query := `UPDATE APPLICATION_USER SET VERIFICATION_CODE = :verificationCode WHERE USER_ID = :userId`
    _, err := s.db.ExecContext(s.ctx, query, 
        sql.Named("verificationCode", code),
        sql.Named("userId", userId),
    )
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
    query := `UPDATE APPLICATION_USER SET MACHINE_ID = :machineId WHERE USER_ID = :userId`
    _, err = db.ExecContext(s.ctx, query, 
        sql.Named("machineId", string(hashedID)),
        sql.Named("userId", userId),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) UpdatePlayerId(playerId string, userId int64, conn *sqlx.DB) error {
    db := conn
    if db == nil {
        db = s.db
    }
    // First nullify any existing player_id with same id
    _, err := db.ExecContext(s.ctx, `UPDATE APPLICATION_USER SET PLAYER_ID = NULL WHERE PLAYER_ID = :playerId`, playerId)
    if err != nil {
        utils.LogError(err)
        return err
    }
    // Then set new player_id
    _, err = db.ExecContext(s.ctx, `UPDATE APPLICATION_USER SET PLAYER_ID = :playerId WHERE USER_ID = :userId`, playerId, userId)
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
        USING (SELECT :playerId AS PLAYER_ID FROM DUAL) src
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
        USING (SELECT :machineId AS MACHINE_ID, :playerId AS PLAYER_ID FROM DUAL) src
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

    query := `UPDATE APPLICATION_USER SET FIRST_TIME_LOGIN = 0, FIRST_TIME_BIOMETRIC = 0 WHERE USER_ID = :userId`
    _, err = tx.ExecContext(s.ctx, query, o.UserID.Int64)
    if err != nil {
        return false, err
    }
    sessionID, err := s.SaveSessionId(o.UserID.Int64, s.db)
    if err != nil {
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
    query := `UPDATE APPLICATION_USER SET FIRST_TIME_LOGIN = 0 WHERE USER_ID = :userId`
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
        WHERE USER_ID = :userId
    `
    _, err = tx.ExecContext(s.ctx, updateUserQuery, user.UserID.Int64)
    if err != nil {
        return err
    }

    _, err = tx.ExecContext(s.ctx, `UPDATE ASSIGN_BRANCH SET PRN = NULL WHERE USER_ID = :userId`, user.UserID.Int64)
    if err != nil {
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
        (:prn, :username, :patientName, :action, :actionDesc, :remarks, :userCreate)
    `
    _, err = tx.ExecContext(s.ctx, auditQuery,
        sql.Named("prn", user.MasterPrn.String),
        sql.Named("username", user.Username.String),
        sql.Named("patientName", patientName),
        sql.Named("action", "Delete Account"),
        sql.Named("actionDesc", actionDesc),
        sql.Named("remarks", remarks),
        sql.Named("userCreate", userCreate),
    )
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (s *ApplicationUserService) DisableFirstTimeBiometricUser(userId int64) error {
    query := `UPDATE APPLICATION_USER SET FIRST_TIME_BIOMETRIC = 0 WHERE USER_ID = :userId`
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

    _, err = tx.ExecContext(s.ctx, `DELETE FROM ASSIGN_BRANCH WHERE USER_ID = :userId AND PRN = :prn`, userId, prn)
    if err != nil {
        return err
    }
    _, err = tx.ExecContext(s.ctx, `DELETE FROM APPLICATION_USER WHERE USER_ID = :userId AND MASTER_PRN = :prn`, userId, prn)
    if err != nil {
        return err
    }
    return tx.Commit()
}

func (s *ApplicationUserService) SetLogin(userId int64) error {
    query := `UPDATE APPLICATION_USER SET IS_LOGGED_IN = 1, DATE_LOGGED_IN = CURRENT_TIMESTAMP WHERE USER_ID = :userId`
    _, err := s.db.ExecContext(s.ctx, query, userId)
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ApplicationUserService) SetLogout(userId int64) error {
    query := `UPDATE APPLICATION_USER SET IS_LOGGED_IN = 0 WHERE USER_ID = :userId`
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

func (s *ApplicationUserService) ValidateCredentials(user *model.ApplicationUser, password string) bool {
    if user.Password.String == "" {
        return false
    }
    err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
    return err == nil
}

func (s *ApplicationUserService) ValidateCredentials2(user *model.ApplicationUser, password string) bool {
    if user.MachineID.String == "" {
        return false
    }
    err := bcrypt.CompareHashAndPassword([]byte(user.MachineID.String), []byte(password))
    return err == nil
}

/* func getApplicationUserCols() string {
    return `
        USER_ID,
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
        ADDRESS_1,
        ADDRESS_2,
        ADDRESS_3,
        CITYSTATE,
        POSTCODE,
        COUNTRY,
        CONTACT_NUMBER,
        PASSPORT,
        NATIONALITY,
        VERIFICATION_CODE,
        FIRST_TIME_LOGIN,
        FIRST_TIME_BIOMETRIC,
        ROLE,
        MASTER_PRN,
        PLAYER_ID,
        MACHINE_ID,
        REGISTRATION_DATE_TIME,
        INACTIVE_FLAG,
        SESSION_ID,
        SIGN_IN_TYPE
    `
}

func getApplicationUserFamilyCols() string {
    return `
        AUF_ID,
        USER_ID,
        PATIENT_PRN,
        NOK_REF_NUMBER,
        IS_PATIENT,
        FULLNAME,
        RELATIONSHIP,
        NOK_PRN,
        NRIC_PASSPORT,
        DOC_NUMBER,
        DOB,
        GENDER,
        NATIONALITY,
        CONTACT_NUMBER,
        ADDRESS,
        IS_ACTIVE,
        MARITAL_STATUS,
        EMAIL,
        IS_KIDS_EXPLORER,
        IS_GOLDEN_PEARL
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
} */
