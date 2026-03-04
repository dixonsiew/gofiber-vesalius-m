package generalNotificationMaster

import (
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/nleeper/goment"
)

func FindAll(offset int, limit int) ([]model.GeneralNotification, error) {
    lx := make([]model.GeneralNotification, 0)
    db := database.GetDb()
    rows, err := db.Queryx(
        `
        SELECT * FROM GENERAL_NOTIFICATION_MASTER
        ORDER BY DATE_CREATE DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, offset, limit)
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

        if o.StartDate.Valid {
            g, _ := goment.New(o.StartDate.String, "YYYY-MM-DD[T]HH:mm:ss")
            o.StartDate.String = g.Format("DD/MM/YYYY")
        }

        if o.EndDate.Valid {
            g, _ := goment.New(o.EndDate.String, "YYYY-MM-DD[T]HH:mm:ss")
            o.EndDate.String = g.Format("DD/MM/YYYY")
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
    q := `SELECT COUNT(NOTIFICATION_MASTER_ID) AS COUNT FROM GENERAL_NOTIFICATION_MASTER`
    err := db.QueryRowx(q).Scan(&n)
    if err != nil {
        utils.LogError(err)
        return n, err
    }

    return n, nil
}
