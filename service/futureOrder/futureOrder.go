package country

import (
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"
)

func FindAll(patientPrn string, offset int, limit int) ([]model.FutureOrder, error) {
    lx := make([]model.FutureOrder, 0)
    db := database.GetDb()
    rows, err := db.Queryx(
        `SELECT np.PATIENT_NAME, npfo.PRN, npfo.PLAN_TYPE, 
        npfo.START_DATE_TIME, npfo.ORDER_DOCTOR,
        COALESCE(nsm_service.DESCRIPTION, nsm_package.DESCRIPTION) AS DESCRIPTION
        FROM NOVA_PATIENT_FUTURE_ORDER npfo
        JOIN NOVA_PATIENT np ON npfo.PRN = np.PRN
        LEFT JOIN NOVA_SERVICE_MASTER nsm_service ON nsm_service.CODE = npfo.SERVICE_CODE
        LEFT JOIN NOVA_SERVICE_MASTER nsm_package ON nsm_package.CODE = npfo.PACKAGE_CODE
        WHERE npfo.FUTURE_ORDER_TYPE = 'ORDER'
        AND (npfo.START_DATE_TIME IS NULL OR TRUNC(npfo.START_DATE_TIME) > TRUNC(SYSDATE))
        AND npfo.PRN = :patientPrn

        UNION 

        SELECT np.PATIENT_NAME, npfo.PRN, npfo.PLAN_TYPE, 
        npfo.START_DATE_TIME, npfo.ORDER_DOCTOR,
        COALESCE(nsm_service.DESCRIPTION, nsm_package.DESCRIPTION) AS DESCRIPTION
        FROM NOVA_PATIENT_FUTURE_ORDER npfo
        JOIN NOVA_PATIENT np ON npfo.PRN = np.PRN
        LEFT JOIN NOVA_SERVICE_MASTER nsm_service ON nsm_service.CODE = npfo.SERVICE_CODE
        LEFT JOIN NOVA_SERVICE_MASTER nsm_package ON nsm_package.CODE = npfo.PACKAGE_CODE
        WHERE npfo.FUTURE_ORDER_TYPE = 'ORDER'
        AND (npfo.START_DATE_TIME IS NULL OR TRUNC(npfo.START_DATE_TIME) > TRUNC(SYSDATE))
        AND npfo.PRN IN (
        SELECT PRN FROM NOVA_PATIENT_NOK WHERE PATIENT_PRN = :patientPrn AND PRN IS NOT NULL
        )
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, patientPrn, offset, limit)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.DbFutureOrder{}
        err := rows.StructScan(&o)
        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        k := model.FutureOrder{}
        k.FromDbModel(o)
        lx = append(lx, k)
    }

    return lx, nil
}

func List(patientPrn string, page string, limit string) (model.PagedList, error) {
    m := model.PagedList{}
    total, err := Count(patientPrn)
    if err != nil {
        return m, err
    }

    pg := model.GetPager(total, page, limit)
    lx, err := FindAll(patientPrn, pg.GetLowerBound(), pg.PageSize)
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

func Count(patientPrn string) (int, error) {
    n := 0
    db := database.GetDb()
    err := db.QueryRowx(
        `SELECT COUNT(*) AS COUNT FROM (
        SELECT np.PATIENT_NAME, npfo.PRN, npfo.PLAN_TYPE, 
        npfo.START_DATE_TIME, npfo.ORDER_DOCTOR,
        COALESCE(nsm_service.DESCRIPTION, nsm_package.DESCRIPTION) AS DESCRIPTION
        FROM NOVA_PATIENT_FUTURE_ORDER npfo
        JOIN NOVA_PATIENT np ON npfo.PRN = np.PRN
        LEFT JOIN NOVA_SERVICE_MASTER nsm_service ON nsm_service.CODE = npfo.SERVICE_CODE
        LEFT JOIN NOVA_SERVICE_MASTER nsm_package ON nsm_package.CODE = npfo.PACKAGE_CODE
        WHERE npfo.FUTURE_ORDER_TYPE = 'ORDER'
        AND (npfo.START_DATE_TIME IS NULL OR TRUNC(npfo.START_DATE_TIME) > TRUNC(SYSDATE))
        AND npfo.PRN = :patientPrn

        UNION 

        SELECT np.PATIENT_NAME, npfo.PRN, npfo.PLAN_TYPE, 
        npfo.START_DATE_TIME, npfo.ORDER_DOCTOR,
        COALESCE(nsm_service.DESCRIPTION, nsm_package.DESCRIPTION) AS DESCRIPTION
        FROM NOVA_PATIENT_FUTURE_ORDER npfo
        JOIN NOVA_PATIENT np ON npfo.PRN = np.PRN
        LEFT JOIN NOVA_SERVICE_MASTER nsm_service ON nsm_service.CODE = npfo.SERVICE_CODE
        LEFT JOIN NOVA_SERVICE_MASTER nsm_package ON nsm_package.CODE = npfo.PACKAGE_CODE
        WHERE npfo.FUTURE_ORDER_TYPE = 'ORDER'
        AND (npfo.START_DATE_TIME IS NULL OR TRUNC(npfo.START_DATE_TIME) > TRUNC(SYSDATE))
        AND npfo.PRN IN (
        SELECT PRN FROM NOVA_PATIENT_NOK WHERE PATIENT_PRN = :patientPrn AND PRN IS NOT NULL
        )
        )`, patientPrn).Scan(&n)
    if err != nil {
        utils.LogError(err)
        return n, err
    }

    return n, nil
}
