package adminuser

import (
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "golang.org/x/crypto/bcrypt"
)

func FindAll(offset int, limit int) ([]model.AdminUser, error) {
    lx := make([]model.AdminUser, 0)
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM ADMIN_USER ORDER BY FIRST_NAME OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, offset, limit)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.AdminUser{}
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        lx = append(lx, o)
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
    err := db.QueryRowx(`SELECT COUNT(ADMIN_ID) AS COUNT FROM ADMIN_USER`).Scan(&n)
    if err != nil {
        utils.LogError(err)
        return n, err
    }

    return n, nil
}

func FindByAdminId(adminId int64) (*model.AdminUser, error) {
    o := model.AdminUser{}
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

        x = &o
    }

    return x, nil
}

func FindByEmail(email string) (*model.AdminUser, error) {
    o := model.AdminUser{}
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

        x = &o
    }

    return x, nil
}

func FindByUsername(email string) (*model.AdminUser, error) {
    o := model.AdminUser{}
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

        x = &o
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
        o := model.AdminUser{}
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        lx = append(lx, o)
    }

    return lx, nil
}

func FindWithAssignBranchByAdminId(adminId int64) (*model.AdminUser, error) {
    o := model.AdminUser{}
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

            o.Password.String = ""
            x = &o
        }

        ab := model.AssignBranch{}
        b := model.Branch{}
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

        b.Passcode.String = ""
        b.Url.String = ""
        ab.Branch = b
        lx = append(lx, ab)

        i++
    }

    k.AdminBranches = lx
    x = &k

    return x, nil
}

func FindAssignBranchByAdminId(adminId int64, branchId int64) (*model.AssignBranch, error) {
    o := model.AssignBranch{}
    var x *model.AssignBranch
    db := database.GetDb()
    err := db.QueryRowx(`SELECT * FROM ASSIGN_BRANCH WHERE BRANCH_ID = :branchId AND ADMIN_ID IN (SELECT ADMIN_ID FROM ADMIN_USER WHERE ADMIN_ID = :adminId)`, branchId, adminId).StructScan(&o)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    x = &o

    return x, nil
}

func ValidateCredentials(user model.AdminUser, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
    return err == nil
}
