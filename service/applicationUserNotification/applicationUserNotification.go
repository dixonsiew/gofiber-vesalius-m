package applicationusernotification

import (
	"vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/utils"
)

func CountUnseenByUserId(userId int64) (int, error) {
    n := 0
    db := database.GetDb()
    q := `SELECT COUNT(USER_ID) AS COUNT FROM APPLICATION_USER_NOTIFICATION WHERE USER_ID = :userId AND DATE_SENT IS NOT NULL AND IS_SEEN = 'N'`
    err := db.QueryRowx(q, userId).Scan(&n)
    if err != nil {
        utils.LogError(err)
        return n, err
    }

    return n, nil
}

func FindAllByUserId(userId int64, offset int, limit int) ([]model.OneSignalNotification, error) {
    lx := make([]model.OneSignalNotification, 0)
    db := database.GetDb()
    q := `
        SELECT * FROM APPLICATION_USER_NOTIFICATION WHERE USER_ID = :userId AND DATE_SENT IS NOT NULL
        ORDER BY DATE_SENT DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    rows, err := db.Queryx(q, userId, offset, limit)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.OneSignalNotification{}
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        o.Set()
        lx = append(lx, o)
    }

    return lx, nil
}

func ListByUserId(userId int64, page string, limit string) (model.PagedList, error) {
    m := model.PagedList{}
    total, err := CountByUserId(userId)
    if err != nil {
        return m, err
    }

    pg := model.GetPager(total, page, limit)
    lx, err := FindAllByUserId(userId, pg.GetLowerBound(), pg.PageSize)
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

func CountByUserId(userId int64) (int, error) {
    n := 0
    db := database.GetDb()
    q := `SELECT COUNT(USER_ID) AS COUNT FROM APPLICATION_USER_NOTIFICATION WHERE USER_ID = :userId AND DATE_SENT IS NOT NULL`
    err := db.QueryRowx(q, userId).Scan(&n)
    if err != nil {
        utils.LogError(err)
        return n, err
    }

    return n, nil
}
