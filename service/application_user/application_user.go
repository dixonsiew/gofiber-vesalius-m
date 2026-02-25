package applicationuser

import (
	"database/sql"
	"math/rand/v2"
	"strings"
	"vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func FindAll(offset int, limit int) ([]model.ApplicationUser, error) {
    lx := make([]model.ApplicationUser, 0)
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM APPLICATION_USER WHERE INACTIVE_FLAG = 'N' ORDER BY REGISTRATION_DATE_TIME, MASTER_PRN OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, offset, limit)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.DbApplicationUser{}
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        k := model.ApplicationUser{}
        k.FromDbModel(o)
        lx = append(lx, k)
    }

    return lx, nil
}

func List(page string, limit string) (model.PagedList, error) {
    m := model.PagedList{}
    total, err := Count()
    if err != nil {
        return m, err
    }

    pg := model.GetPager(total, page, limit)
    lx, err := FindAll(pg.GetLowerBound(), pg.PageSize)
    if err != nil {
        return m, err
    }

    m = model.PagedList{
        List: lx,
        Total: total,
        TotalPages: pg.GetTotalPages(),
    }

    return m, nil
}

func Count() (int, error) {
    n := 0
    db := database.GetDb()
    err := db.QueryRowx(`SELECT COUNT(USER_ID) AS COUNT FROM APPLICATION_USER WHERE INACTIVE_FLAG = 'N'`).Scan(&n)
    if err != nil {
        utils.LogError(err)
        return n, err
    }

    return n, nil
}

func FindAllActive(offset int, limit int) ([]model.ApplicationUser, error) {
    lx := make([]model.ApplicationUser, 0)
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM APPLICATION_USER WHERE INACTIVE_FLAG = 'N' ORDER BY REGISTRATION_DATE_TIME, MASTER_PRN OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, offset, limit)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.DbApplicationUser{}
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        k := model.ApplicationUser{}
        k.FromDbModel(o)
        lx = append(lx, k)
    }

    return lx, nil
}

func ListActive(page string, limit string) (model.PagedList, error) {
    m := model.PagedList{}
    total, err := Count()
    if err != nil {
        return m, err
    }

    pg := model.GetPager(total, page, limit)
    lx, err := FindAllActive(pg.GetLowerBound(), pg.PageSize)
    if err != nil {
        return m, err
    }

    m = model.PagedList{
        List: lx,
        Total: total,
        TotalPages: pg.GetTotalPages(),
    }

    return m, nil
}

func CountActive() (int, error) {
    n := 0
    db := database.GetDb()
    err := db.QueryRowx(`SELECT COUNT(USER_ID) AS COUNT FROM APPLICATION_USER WHERE INACTIVE_FLAG = 'N'`).Scan(&n)
    if err != nil {
        utils.LogError(err)
        return n, err
    }

    return n, nil
}

func FindByKeyword(keyword string, offset int, limit int) ([]model.ApplicationUser, error) {
    lx := make([]model.ApplicationUser, 0)
    db := database.GetDb()
    q := `
        SELECT * FROM APPLICATION_USER au
        WHERE (LOWER(au.FIRST_NAME) LIKE :keyword OR LOWER(au.MIDDLE_NAME) LIKE :keyword OR LOWER(au.LAST_NAME) LIKE :keyword
        OR au.MASTER_PRN LIKE :keyword OR LOWER(au.EMAIL) LIKE :keyword)
        AND INACTIVE_FLAG = 'N'
        ORDER BY REGISTRATION_DATE_TIME, MASTER_PRN OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    rows, err := db.Queryx(q, keyword, offset, limit)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.DbApplicationUser{}
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        k := model.ApplicationUser{}
        k.FromDbModel(o)
        lx = append(lx, k)
    }

    return lx, nil
}

func ListByKeyword(keyword string, page string, limit string) (model.PagedList, error) {
    m := model.PagedList{}
    total, err := CountByKeyword(keyword)
    if err != nil {
        return m, err
    }

    pg := model.GetPager(total, page, limit)
    lx, err := FindByKeyword(keyword, pg.GetLowerBound(), pg.PageSize)
    if err != nil {
        return m, err
    }

    m = model.PagedList{
        List: lx,
        Total: total,
        TotalPages: pg.GetTotalPages(),
    }

    return m, nil
}

func CountByKeyword(keyword string) (int, error) {
    n := 0
    db := database.GetDb()
    q := `
        SELECT COUNT(au.USER_ID) AS COUNT FROM APPLICATION_USER au
        WHERE (LOWER(au.FIRST_NAME) LIKE :keyword OR LOWER(au.MIDDLE_NAME) LIKE :keyword OR LOWER(au.LAST_NAME) LIKE :keyword
        OR au.MASTER_PRN LIKE :keyword OR LOWER(au.EMAIL) LIKE :keyword) AND INACTIVE_FLAG = 'N'
    `
    err := db.QueryRowx(q, keyword).Scan(&n)
    if err != nil {
        utils.LogError(err)
        return n, err
    }

    return n, nil
}

func FindByUserIdSessionId(userId int64, sessionId string) (*model.ApplicationUser, error) {
    o := model.DbApplicationUser{}
    k := model.ApplicationUser{}
    var x *model.ApplicationUser
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM APPLICATION_USER WHERE USER_ID = :userId AND SESSION_ID = :sessionId`, userId, sessionId)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    defer rows.Close()

    if rows.Next() {
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return x, err
        }

        k.FromDbModel(o)
        x = &k
    }

    return x, nil
}

func FindByUserId(userId int64) (*model.ApplicationUser, error) {
    o := model.DbApplicationUser{}
    k := model.ApplicationUser{}
    var x *model.ApplicationUser
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM APPLICATION_USER WHERE USER_ID = :userId`, userId)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    defer rows.Close()

    if rows.Next() {
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return x, err
        }

        k.FromDbModel(o)
        x = &k
    }

    return x, nil
}

func FindByUsername(username string) (*model.ApplicationUser, error) {
    o := model.DbApplicationUser{}
    k := model.ApplicationUser{}
    var x *model.ApplicationUser
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM APPLICATION_USER WHERE LOWER(USERNAME) = LOWER(:username) ORDER BY REGISTRATION_DATE_TIME DESC`, username)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    defer rows.Close()

    if rows.Next() {
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return x, err
        }

        k.FromDbModel(o)
        x = &k
    }

    return x, nil
}

func FindByEmail(email string) (*model.ApplicationUser, error) {
    o := model.DbApplicationUser{}
    k := model.ApplicationUser{}
    var x *model.ApplicationUser
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM APPLICATION_USER WHERE LOWER(EMAIL) = LOWER(:email)`, email)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    defer rows.Close()

    if rows.Next() {
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return x, err
        }

        k.FromDbModel(o)
        x = &k
    }

    return x, nil
}

func FindByPRN(prn string) (*model.ApplicationUser, error) {
    o := model.DbApplicationUser{}
    k := model.ApplicationUser{}
    var x *model.ApplicationUser
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return x, nil
    }

    rows, err := db.Queryx(`SELECT * FROM APPLICATION_USER WHERE MASTER_PRN = :prn`, prn)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    defer rows.Close()

    if rows.Next() {
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return x, err
        }

        k.FromDbModel(o)
        x = &k
    }

    return x, nil
}

func FindWithAssignBranchByUserId(userId int64) (*model.ApplicationUser, error) {
    o := model.DbApplicationUser{}
    k := model.ApplicationUser{}
    var x *model.ApplicationUser
    db := database.GetDb()
    q := `SELECT * FROM APPLICATION_USER au 
        LEFT JOIN ASSIGN_BRANCH ab ON au.USER_ID = ab.USER_ID 
        INNER JOIN BRANCH b ON b.BRANCH_ID = ab.BRANCH_ID 
        WHERE au.USER_ID = :userId`
    rows, err := db.Queryx(q, userId)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    defer rows.Close()

    i := 0
    lx := make([]model.AssignBranch, 0)
    for rows.Next() {
        if i == 0 {
            err := rows.StructScan(&o)
            if err != nil {
                utils.LogError(err)
                return x, err
            }

            k.FromDbModel(o)
            k.Password = ""
            x = &k
        }

        ab := model.DbAssignBranch{}
        mab := model.AssignBranch{}
        b := model.DbBranch{}
        mb := model.Branch{}
        err := rows.StructScan(&ab)
        if err != nil {
            utils.LogError(err)
            return x, err
        }
        
        err = rows.StructScan(&b)
        if err != nil {
            utils.LogError(err)
            return x, err
        }

        mab.FromDbModel(ab)
        mb.FromDbModel(b)
        mb.Passcode = ""
        mb.Url = ""
        mab.Branch = mb
        lx = append(lx, mab)

        i++
    }

    k.UserBranches = lx
    x = &k

    return x, nil
}

func FindWithAssignBranchByEmail(email string) (*model.ApplicationUser, error) {
    o := model.DbApplicationUser{}
    k := model.ApplicationUser{}
    var x *model.ApplicationUser
    db := database.GetDb()
    q := `SELECT * FROM APPLICATION_USER au 
        LEFT JOIN ASSIGN_BRANCH ab ON au.USER_ID = ab.USER_ID 
        INNER JOIN BRANCH b ON b.BRANCH_ID = ab.BRANCH_ID 
        WHERE au.EMAIL = :email`
    rows, err := db.Queryx(q, email)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    defer rows.Close()

    i := 0
    lx := make([]model.AssignBranch, 0)
    for rows.Next() {
        if i == 0 {
            err := rows.StructScan(&o)
            if err != nil {
                utils.LogError(err)
                return x, err
            }

            k.FromDbModel(o)
            k.Password = ""
            x = &k
        }

        ab := model.DbAssignBranch{}
        mab := model.AssignBranch{}
        b := model.DbBranch{}
        mb := model.Branch{}
        err := rows.StructScan(&ab)
        if err != nil {
            utils.LogError(err)
            return x, err
        }
        
        err = rows.StructScan(&b)
        if err != nil {
            utils.LogError(err)
            return x, err
        }

        mab.FromDbModel(ab)
        mb.FromDbModel(b)
        mb.Passcode = ""
        mb.Url = ""
        mab.Branch = mb
        lx = append(lx, mab)

        i++
    }

    k.UserBranches = lx
    x = &k

    return x, nil
}

func FindAssignBranchByUserId(userId int64, branchId int64) (*model.AssignBranch, error) {
    o := model.DbAssignBranch{}
    k := model.AssignBranch{}
    var x *model.AssignBranch
    db := database.GetDb()
    err := db.QueryRowx(`SELECT * FROM ASSIGN_BRANCH WHERE BRANCH_ID = :branchId AND USER_ID IN (SELECT USER_ID FROM APPLICATION_USER WHERE USER_ID = :userId)`, branchId, userId).StructScan(&o)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    k.FromDbModel(o)
    x = &k

    return x, nil
}

func FindAssignBranchByEmail(email string, branchId int64) (*model.AssignBranch, error) {
    o := model.DbAssignBranch{}
    k := model.AssignBranch{}
    var x *model.AssignBranch
    db := database.GetDb()
    err := db.QueryRowx(`SELECT * FROM ASSIGN_BRANCH WHERE BRANCH_ID = :branchId AND USER_ID IN (SELECT USER_ID FROM APPLICATION_USER WHERE EMAIL = :email)`, branchId, email).StructScan(&o)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    k.FromDbModel(o)
    x = &k

    return x, nil
}

func FindByOtherPRN(prn string, userId int64) (*model.ApplicationUser, error) {
    o := model.DbApplicationUser{}
    k := model.ApplicationUser{}
    var x *model.ApplicationUser
    db := database.GetDb()
    err := db.QueryRowx(`SELECT * FROM APPLICATION_USER WHERE USER_ID IN (SELECT USER_ID FROM ASSIGN_BRANCH WHERE USER_ID <> :userId AND PRN = :prn)`, prn, userId).StructScan(&o)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    k.FromDbModel(o)
    x = &k

    return x, nil
}

func ExistsByEmail(email string) (bool, error) {
    b := false
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE LOWER(EMAIL) = LOWER(:email))`, email)
    if err != nil {
        utils.LogError(err)
        return b, err
    }

    defer rows.Close()

    if rows.Next() {
        var r sql.NullInt32
        err := rows.Scan(&r)

        if err != nil {
            utils.LogError(err)
            return b, err
        }

        b = r.Int32 == 1
    }

    return b, nil
}

func ExistsByPRN(prn string) (bool, error) {
    b := false
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE MASTER_PRN = :prn)`, prn)
    if err != nil {
        utils.LogError(err)
        return b, err
    }

    defer rows.Close()

    if rows.Next() {
        var r sql.NullInt32
        err := rows.Scan(&r)

        if err != nil {
            utils.LogError(err)
            return b, err
        }

        b = r.Int32 == 1
    }

    return b, nil
}

func ExistsByMobileNo(mobileNo string) (bool, error) {
    b := false
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE USERNAME = :mobileNo)`, mobileNo)
    if err != nil {
        utils.LogError(err)
        return b, err
    }

    defer rows.Close()

    if rows.Next() {
        var r sql.NullInt32
        err := rows.Scan(&r)

        if err != nil {
            utils.LogError(err)
            return b, err
        }

        b = r.Int32 == 1
    }

    return  b, nil
}

func SaveUserBranch(branchId int64, o model.ApplicationUser) error {
    db := database.GetDb()
    tx, err := db.Begin()
    if err != nil {
        utils.LogError(err)
        return err
    }

    q := `UPDATE APPLICATION_USER
        SET ADDRESS = :address, CONTACT_NUMBER = :contact_number, DOB = :dob, FIRST_NAME = :first_name, LAST_NAME = :last_name, 
        MASTER_PRN = :master_prn, MIDDLE_NAME = :middle_name, NATIONALITY = :nationality, PASSPORT = :passport, RESIDENT = :resident, 
        SEX = :sex, TITLE = :title
        WHERE USER_ID = :user_id`
    stmt, err := tx.Prepare(q)
    if err != nil {
        tx.Rollback()
        utils.LogError(err)
        return err
    }

    defer stmt.Close()

    _, err = stmt.Exec(
        sql.Named("address", o.Address), 
        sql.Named("contact_number", o.ContactNumber), 
        sql.Named("dob", o.Dob), 
        sql.Named("first_name", o.FirstName), 
        sql.Named("last_name", o.LastName), 
        sql.Named("master_prn", o.MasterPrn), 
        sql.Named("middle_name", o.MiddleName), 
        sql.Named("nationality", o.Nationality), 
        sql.Named("passport", o.Passport), 
        sql.Named("resident", o.Resident), 
        sql.Named("sex", o.Sex), 
        sql.Named("title", o.Title), 
        sql.Named("user_id", o.UserID),
    )
    if err != nil {
        tx.Rollback()
        utils.LogError(err)
        return err
    }

    q = `INSERT INTO ASSIGN_BRANCH (ASSIGN_BRANCH_ID, ADMIN_ID, PRN, USER_ID, BRANCH_ID) VALUES(USER_BRANCH_SEQ.nextval, NULL, :prn, :user_id, :branchId)`
    stmti, err := tx.Prepare(q)
    if err != nil {
        tx.Rollback()
        utils.LogError(err)
        return err
    }

    defer stmti.Close()

    _, err = stmti.Exec(
        sql.Named("prn", o.MasterPrn), 
        sql.Named("user_id", o.UserID), 
        sql.Named("branchId", branchId),
    )
    if err != nil {
        tx.Rollback()
        utils.LogError(err)
        return err
    }

    err = tx.Commit()
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func UpdateVerificationCode(code string, userId int64) error {
    db := database.GetDb()
    _, err := db.Exec(`UPDATE APPLICATION_USER SET VERIFICATION_CODE = :code WHERE USER_ID = :userId`, code, userId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func UpdateMachineId(id string, userId int64) error {
    mid := []byte(id)
    machine_id_byte, err := bcrypt.GenerateFromPassword(mid, bcrypt.DefaultCost)
    machine_id := string(machine_id_byte)
    if err != nil {
        utils.LogError(err)
        return err
    }

    db := database.GetDb()
    _, err = db.Exec(`UPDATE APPLICATION_USER SET MACHINE_ID = :machine_id WHERE USER_ID = :userId`, machine_id, userId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func UpdatePlayerId(id string, userId int64) error {
    db := database.GetDb()
    _, err := db.Exec(`UPDATE APPLICATION_USER SET PLAYER_ID = NULL WHERE PLAYER_ID = :id`, id)
    if err != nil {
        utils.LogError(err)
        return err
    }

    _, err = db.Exec(`UPDATE APPLICATION_USER SET PLAYER_ID = :id WHERE USER_ID = :userId`, id, userId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func SaveSessionId(userId int64) (string, error) {
    sid := ""
    sessionId := uuid.New().String()
    db := database.GetDb()
    _, err := db.Exec(`UPDATE APPLICATION_USER SET SESSION_ID = :sessionId WHERE USER_ID = :userId`, sessionId, userId)
    if err != nil {
        utils.LogError(err)
        return sid, err
    }

    sid = sessionId
    return sid, nil
}

func DeleteUserAccount() {

}

func InsertDownloadApp(playerId string) error {
    db := database.GetDb()
    q := `MERGE INTO APP_DOWNLOADED_USER apu
        USING (SELECT :playerId AS PLAYER_ID FROM DUAL) src
        ON (apu.PLAYER_ID = src.PLAYER_ID)
        WHEN NOT MATCHED THEN
        INSERT (PLAYER_ID) VALUES (src.PLAYER_ID)`
    _, err := db.Exec(q, playerId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func InsertDownloadAppV2(machineId string, playerId string) error {
    db := database.GetDb()
    q := `MERGE INTO APP_DOWNLOADED_USER apu
        USING (SELECT :machineId AS MACHINE_ID, :playerId AS PLAYER_ID FROM DUAL) src
        ON (apu.MACHINE_ID = src.MACHINE_ID)
         WHEN MATCHED THEN
          UPDATE SET apu.PLAYER_ID = src.PLAYER_ID, DATE_UPDATE = CURRENT_TIMESTAMP
         WHEN NOT MATCHED THEN
          INSERT (MACHINE_ID, PLAYER_ID, DATE_UPDATE) VALUES (src.MACHINE_ID, src.PLAYER_ID, CURRENT_TIMESTAMP)`
    _, err := db.Exec(q, machineId, playerId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func DisableFirstTimeBiometricUser(userId int64) error {
    db := database.GetDb()
    _, err := db.Exec(`UPDATE APPLICATION_USER SET FIRST_TIME_BIOMETRIC = 0 WHERE USER_ID = :userId`, userId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func ResetUserSignup(userId int64, prn string) error {
    db := database.GetDb()
    _, err := db.Exec(`DELETE FROM ASSIGN_BRANCH WHERE USER_ID = :userId AND PRN = :prn`, userId, prn)
    if err != nil {
        utils.LogError(err)
        return err
    }

    _, err = db.Exec(`DELETE FROM APPLICATION_USER WHERE USER_ID = :userId AND MASTER_PRN = :prn`, userId, prn)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func setLogin(userId int64) error {
    db := database.GetDb()
    _, err := db.Exec(`UPDATE APPLICATION_USER SET IS_LOGGED_IN = 1, DATE_LOGGED_IN = CURRENT_TIMESTAMP WHERE USER_ID = :userId`, userId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func SetLogout(userId int64) error {
    db := database.GetDb()
    _, err := db.Exec(`UPDATE APPLICATION_USER SET IS_LOGGED_IN = 0 WHERE USER_ID = :userId`, userId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func GetRandomStr(length int) string {
    ls := make([]string, 0)
    characters := "0123456789"
    for i := 0; i < length; i++ {
        i := rand.IntN(len(characters))
        j := i + 1
        ls = append(ls, characters[i:j])
    }
    
    return strings.Join(ls, "")
}

func ValidateCredentials(user model.ApplicationUser, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    return err == nil
}

func ValidateCredentials2(user model.ApplicationUser, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.MachineID), []byte(password))
    return err == nil
}
