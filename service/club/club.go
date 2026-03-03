package club

import (
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"
)

func FindGoldenPearlAboutUs() (*model.GoldenPearlAboutUs, error) {
    var x *model.GoldenPearlAboutUs
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT * FROM GOLDEN_CLUB_INFO`)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    defer rows.Close()

    if rows.Next() {
        o := model.DbGoldenPearlAboutUs{}
        err := rows.StructScan(&o)

        if err != nil {
            utils.LogError(err)
            return x, err
        }

        k := model.GoldenPearlAboutUs{}
        k.FromDbModel(o)
        x = &k
    }

    return x, nil
}
