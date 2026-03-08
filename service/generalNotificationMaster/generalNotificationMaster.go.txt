package generalNotificationMaster

import (
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"
)

func FindByNotificationMasterId(notificationMasterId int64) (*model.GeneralNotification, error) {
    o := model.GeneralNotification{}
    var x *model.GeneralNotification
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM GENERAL_NOTIFICATION_MASTER WHERE NOTIFICATION_MASTER_ID = :notificationMasterId`, notificationMasterId)
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

        o.Set()
        x = &o
    }

    return x, nil
}

func FindAll(offset int, limit int) ([]model.GeneralNotification, error) {
    lx := make([]model.GeneralNotification, 0)
    db := database.GetDb()
    q := `
        SELECT * FROM GENERAL_NOTIFICATION_MASTER
        ORDER BY DATE_CREATE DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    rows, err := db.Queryx(q, offset, limit)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.GeneralNotification{}
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
    err := db.QueryRowx(`SELECT COUNT(NOTIFICATION_MASTER_ID) AS COUNT FROM GENERAL_NOTIFICATION_MASTER`).Scan(&n)
    if err != nil {
        utils.LogError(err)
        return n, err
    }

    return n, nil
}
