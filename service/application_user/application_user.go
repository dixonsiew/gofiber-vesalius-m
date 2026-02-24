package applicationuser

import (
	"database/sql"
	"vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/utils"

	"golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
)

func FindByUserId(userId int64) (*model.ApplicationUser, error) {
    o := model.DbApplicationUser{}
    k := model.ApplicationUser{}
    var x *model.ApplicationUser
    db := database.GetDb()
    q := `SELECT * FROM APPLICATION_USER WHERE USER_ID = :userId`
    rows, err := db.Queryx(q, userId)
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
    q := `SELECT * FROM APPLICATION_USER WHERE LOWER(USERNAME) = LOWER(:username) ORDER BY REGISTRATION_DATE_TIME DESC`
    rows, err := db.Queryx(q, username)
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
    q := `SELECT * FROM APPLICATION_USER WHERE LOWER(EMAIL) = LOWER(:email)`
    rows, err := db.Queryx(q, email)
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

    q := `SELECT * FROM APPLICATION_USER WHERE MASTER_PRN = :prn`
    rows, err := db.Queryx(q, prn)
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
    q := `SELECT * FROM ASSIGN_BRANCH WHERE BRANCH_ID = :branchId AND USER_ID IN (SELECT USER_ID FROM APPLICATION_USER WHERE USER_ID = :userId)`
    err := db.QueryRowx(q, branchId, userId).StructScan(&o)
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
    q := `SELECT * FROM ASSIGN_BRANCH WHERE BRANCH_ID = :branchId AND USER_ID IN (SELECT USER_ID FROM APPLICATION_USER WHERE EMAIL = :email)`
    err := db.QueryRowx(q, branchId, email).StructScan(&o)
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
    q := `SELECT * FROM APPLICATION_USER WHERE USER_ID IN (SELECT USER_ID FROM ASSIGN_BRANCH WHERE USER_ID <> :userId AND PRN = :prn)`
    err := db.QueryRowx(q, prn, userId).StructScan(&o)
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
    q := `SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE LOWER(EMAIL) = LOWER(:email))`
    rows, err := db.Queryx(q, email)
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
    q := `SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE MASTER_PRN = :prn)`
    rows, err := db.Queryx(q, prn)
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
    q := `SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM APPLICATION_USER WHERE USERNAME = :mobileNo)`
    rows, err := db.Queryx(q, mobileNo)
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

func UpdateMachineId(id string, userId int64) error {
    mid := []byte(id)
    machine_id, err := bcrypt.GenerateFromPassword(mid, bcrypt.DefaultCost)
    if err != nil {
        utils.LogError(err)
        return err
    }

    db := database.GetDb()
    _, err = db.Exec(`UPDATE APPLICATION_USER SET MACHINE_ID = :mid WHERE USER_ID = :userId`, machine_id, userId)
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
    q := `UPDATE APPLICATION_USER SET SESSION_ID = :sessionId WHERE USER_ID = :userId`
    _, err := db.Exec(q, sessionId, userId)
    if err != nil {
        utils.LogError(err)
        return sid, err
    }

    sid = sessionId
    return sid, nil
}

func ValidateCredentials(user model.ApplicationUser, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    return err == nil
}

func ValidateCredentials2(user model.ApplicationUser, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.MachineID), []byte(password))
    return err == nil
}
