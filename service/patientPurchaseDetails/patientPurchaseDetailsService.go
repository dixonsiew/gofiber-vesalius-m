package patientPurchaseDetails

import (
    "context"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    applicationuserService "vesaliusm/service/applicationUser"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var applicationUserSvc *applicationuserService.ApplicationUserService = applicationuserService.NewApplicationUserService(database.GetDb(), database.GetCtx())

type PatientPurchaseDetailsService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewPatientPurchaseDetailsService(db *sqlx.DB, ctx context.Context) *PatientPurchaseDetailsService {
    return &PatientPurchaseDetailsService{db: db, ctx: ctx}
}

func (s *PatientPurchaseDetailsService) ListByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountByKeyword(keyword, keyword2, keyword3, keyword4)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindByKeyword(keyword, keyword2, keyword3, keyword4, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *PatientPurchaseDetailsService) CountByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string) (int, error) {
    conditions, args := buildKeywordConditions(keyword, keyword2, keyword3, keyword4)
    base := `SELECT COUNT(ppd.PATIENT_PURCHASE_ID) FROM PATIENT_PURCHASE_DETAILS ppd
             JOIN HOSPITAL_PACKAGE hp ON ppd.PACKAGE_ID = hp.PACKAGE_ID`
    query := base + whereClause(conditions)

    rows, err := s.db.NamedQueryContext(s.ctx, query, args)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    defer rows.Close()

    var count int
    if rows.Next() {
        if err = rows.Scan(&count); err != nil {
            utils.LogError(err)
            return 0, err
        }
    }
    return count, nil
}

func (s *PatientPurchaseDetailsService) ListByPrn(userId int64, page string, limit string) (*model.PagedList, error) {
    user, err := applicationUserSvc.FindByUserId(userId, s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    prn := user.MasterPrn
    total, err := s.CountByPrn(prn)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllByPrn(prn, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *PatientPurchaseDetailsService) CountByPrn(prn string) (int, error) {
    query := `SELECT COUNT(PATIENT_PURCHASE_ID) FROM PATIENT_PURCHASE_DETAILS WHERE PATIENT_PRN = :1`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, prn)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *PatientPurchaseDetailsService) List(page string, limit string) (*model.PagedList, error) {
    total, err := s.Count()
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAll(pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *PatientPurchaseDetailsService) Count() (int, error) {
    query := `SELECT COUNT(PATIENT_PURCHASE_ID) FROM PATIENT_PURCHASE_DETAILS WHERE PACKAGE_STATUS = 'Purchased'`
    var count int
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *PatientPurchaseDetailsService) findByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string, offset int, limit int) ([]model.u, error) {
    conditions, args := buildKeywordConditions(keyword, keyword2, keyword3, keyword4)
    args["offset"] = offset
    args["limit"] = limit

    base := `
        SELECT ppd.*, hp.PACKAGE_NAME, ppd2.PAYMENT_REQUEST_NO,
               ppd2.PAYMENT_REQUEST_CURRENCY, ppd2.PAYMENT_AMOUNT, ppd2.PAYMENT_CURRENCY,
               ppd2.PAYMENT_AMOUNT_COLLECTED, ppd2.PAYMENT_STATUS, ppd2.PAYMENT_TRANS_DATE,
               ppd2.BILLING_FULLNAME, ppd2.BILLING_CONTACT_NO, ppd2.BILLING_CONTACT_CODE, 
               ppd2.BILLING_EMAIL, ppd2.PAYMENT_URL
        FROM PATIENT_PURCHASE_DETAILS ppd
        JOIN HOSPITAL_PACKAGE hp ON ppd.PACKAGE_ID = hp.PACKAGE_ID
        JOIN PACKAGE_PAYMENT_DETAILS ppd2 ON ppd.PACKAGE_PAYMENT_ID = ppd2.PACKAGE_PAYMENT_ID
    `
    query := base + whereClause(conditions) +
        ` ORDER BY ppd.DATE_CREATE DESC, ppd.PACKAGE_PURCHASE_NO DESC
          OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    rows, err := sqlx.NamedQuery(q, query, args)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.UserPackage
    for rows.Next() {
        var item models.UserPackage
        if err = rows.StructScan(&item); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, nil
}



func buildKeywordConditions(keyword string, keyword2 string, keyword3 string, keyword4 string) ([]string, map[string]interface{}) {
    var conds []string
    args := make(map[string]interface{})

    if keyword != "" {
        conds = append(conds, `LOWER(ppd.PATIENT_PRN) LIKE LOWER(:keyword)`)
        args["keyword"] = "%" + keyword + "%"
    }
    if keyword2 != "" {
        conds = append(conds, `LOWER(ppd.PACKAGE_PURCHASE_NO) LIKE LOWER(:keyword2)`)
        args["keyword2"] = "%" + keyword2 + "%"
    }
    if keyword3 != "" {
        conds = append(conds, `LOWER(hp.PACKAGE_NAME) LIKE LOWER(:keyword3)`)
        args["keyword3"] = "%" + keyword3 + "%"
    }
    if keyword4 != "" && keyword4 != "All" {
        conds = append(conds, `LOWER(ppd.PACKAGE_STATUS) LIKE LOWER(:keyword4)`)
        args["keyword4"] = "%" + keyword4 + "%"
    }
    return conds, args
}

func whereClause(conds []string) string {
    if len(conds) == 0 {
        return ""
    }
    return " WHERE " + strings.Join(conds, " AND ")
}
