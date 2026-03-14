package futureOrder

import (
    "context"
    "vesaliusm/model"
    "vesaliusm/model/futureOrder"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

type FutureOrderService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewFutureOrderService(db *sqlx.DB, ctx context.Context) *FutureOrderService {
    return &FutureOrderService{db: db, ctx: ctx}
}

func (s *FutureOrderService) FindAll(patientPrn string, offset int, limit int) ([]futureOrder.FutureOrder, error) {
    query := `
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
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    list := make([]futureOrder.FutureOrder, 0)
    err := s.db.SelectContext(s.ctx, &list, query, patientPrn, offset, limit)
    if err != nil {
        utils.LogError(err)
        return list, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
}

func (s *FutureOrderService) List(patientPrn string, page string, limit string) (*model.PagedList, error) {
    total, err := s.Count(patientPrn)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAll(patientPrn, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *FutureOrderService) Count(patientPrn string) (int, error) {
    query := `
        SELECT COUNT(*) AS COUNT FROM (
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
        )
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, patientPrn)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}
