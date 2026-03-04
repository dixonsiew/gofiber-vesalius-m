package country

import (
    "database/sql"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"
)

func FindAllCountryTelCode() ([]model.CountryTelCode, error) {
    lx := make([]model.CountryTelCode, 0)
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT COUNTRY_NAME, TEL_CODE FROM NOVA_COUNTRY WHERE TEL_CODE IS NOT NULL ORDER BY COUNTRY_NAME`)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.CountryTelCode{}
        err := rows.StructScan(&o)

        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        lx = append(lx, o)
    }

    return lx, nil
}

func FindCountryCodeByNationality(nationality string) (string, error) {
    s := ""
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT COUNTRY_CODE FROM NOVA_COUNTRY WHERE NATIONALITY = :nationality`, nationality)
    if err != nil {
        utils.LogError(err)
        return s, err
    }

    defer rows.Close()

    if rows.Next() {
        var r sql.NullString
        err := rows.Scan(&r)

        if err != nil {
            utils.LogError(err)
            return s, err
        }

        s = r.String
    }

    return s, nil
}

func FindAllCountries() ([]model.Country, error) {
    lx := make([]model.Country, 0)
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT COUNTRY_NAME, TEL_CODE, COUNTRY_CODE FROM NOVA_COUNTRY WHERE COUNTRY_NAME IS NOT NULL ORDER BY COUNTRY_NAME`)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.Country{}
        err := rows.StructScan(&o)

        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        lx = append(lx, o)
    }

    return lx, nil
}

func FindAllNationalities() ([]model.Nationality, error) {
    lx := make([]model.Nationality, 0)
    db := database.GetDb()
    rows, err := db.Queryx(`SELECT NATIONALITY FROM NOVA_COUNTRY WHERE NATIONALITY IS NOT NULL ORDER BY NATIONALITY`)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.Nationality{}
        err := rows.StructScan(&o)

        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        lx = append(lx, o)
    }

    return lx, nil
}
