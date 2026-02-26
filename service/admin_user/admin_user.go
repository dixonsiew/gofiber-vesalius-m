package adminuser

import (
	"vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/utils"

    "golang.org/x/crypto/bcrypt"
)

func FindByAdminId(adminId int64) (*model.AdminUser, error) {
    o := model.DbAdminUser{}
    k := model.AdminUser{}
    var x *model.AdminUser
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM ADMIN_USER WHERE ADMIN_ID = :adminId`, adminId)
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

func FindByEmail(email string) (*model.AdminUser, error) {
    o := model.DbAdminUser{}
    k := model.AdminUser{}
    var x *model.AdminUser
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM ADMIN_USER WHERE EMAIL = :email`, email)
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

func FindByUsername(email string) (*model.AdminUser, error) {
    o := model.DbAdminUser{}
    k := model.AdminUser{}
    var x *model.AdminUser
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM ADMIN_USER WHERE USERNAME = :email`, email)
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

func FindByUserGroupId(userGroupId int64) ([]model.AdminUser, error) {
    lx := make([]model.AdminUser, 0)
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM ADMIN_USER WHERE USER_GROUP_ID = :userGroupId`, userGroupId)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.DbAdminUser{}
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        k := model.AdminUser{}
        k.FromDbModel(o)
        lx = append(lx, k)
    }

    return lx, nil
}

func FindWithAssignBranchByAdminId(adminId int64) (*model.AdminUser, error) {
    o := model.DbAdminUser{}
    k := model.AdminUser{}
    var x *model.AdminUser
    db := database.GetDb()
    q := `SELECT * FROM ADMIN_USER au 
        LEFT JOIN ASSIGN_BRANCH ab ON au.ADMIN_ID = ab.ADMIN_ID 
        INNER JOIN BRANCH b ON b.BRANCH_ID = ab.BRANCH_ID 
        WHERE au.ADMIN_ID = :adminId`
    rows, err := db.Queryx(q, adminId)
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

    k.AdminBranches = lx
    x = &k

    return x, nil
}

func FindAssignBranchByAdminId(adminId int64, branchId int64) (*model.AssignBranch, error) {
    o := model.DbAssignBranch{}
    k := model.AssignBranch{}
    var x *model.AssignBranch
    db := database.GetDb()
    err := db.QueryRowx(`SELECT * FROM ASSIGN_BRANCH WHERE BRANCH_ID = :branchId AND ADMIN_ID IN (SELECT ADMIN_ID FROM ADMIN_USER WHERE ADMIN_ID = :adminId)`, branchId, adminId).StructScan(&o)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    k.FromDbModel(o)
    x = &k

    return x, nil
}

func ValidateCredentials(user model.AdminUser, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    return err == nil
}
