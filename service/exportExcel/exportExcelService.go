package exportExcel

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    hp "vesaliusm/model/hpackage"
    lg "vesaliusm/model/logistic"
    upck "vesaliusm/model/userPackage"
    "vesaliusm/service/novaDoctor"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var ExportExcelSvc *ExportExcelService = NewExportExcelService(database.GetDb(), database.GetCtx())

type ExportExcelService struct {
    db                *sqlx.DB
    ctx               context.Context
    novaDoctorService *novaDoctor.NovaDoctorService
}

func NewExportExcelService(db *sqlx.DB, ctx context.Context) *ExportExcelService {
    return &ExportExcelService{
        db:                db,
        ctx:               ctx,
        novaDoctorService: novaDoctor.NewNovaDoctorService(db, ctx),
    }
}

func (s *ExportExcelService) ExportAllMobileUser() ([]model.ApplicationUser, error) {
    query := `
        SELECT * 
        FROM APPLICATION_USER 
        WHERE INACTIVE_FLAG = 'N' 
        ORDER BY REGISTRATION_DATE_TIME, MASTER_PRN
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.ApplicationUser{}, ""), 1)
    list := make([]model.ApplicationUser, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
        s.setMobileUser(&list[i])
    }
    return list, nil
}

func (s *ExportExcelService) ExportMobileUserByKeyword(keyword string) ([]model.ApplicationUser, error) {
    query := `
        SELECT * FROM APPLICATION_USER au
        WHERE (LOWER(au.FIRST_NAME) LIKE :keyword OR LOWER(au.MIDDLE_NAME) LIKE :keyword OR LOWER(au.LAST_NAME) LIKE :keyword
        OR au.MASTER_PRN LIKE :keyword OR LOWER(au.EMAIL) LIKE :keyword)
        AND INACTIVE_FLAG = 'N'
        ORDER BY REGISTRATION_DATE_TIME, MASTER_PRN
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.ApplicationUser{}, ""), 1)
    list := make([]model.ApplicationUser, 0)
    err := s.db.SelectContext(s.ctx, &list, query, sql.Named("keyword", strings.ToLower(keyword)))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
        s.setMobileUser(&list[i])
    }
    return list, nil
}

func (s *ExportExcelService) ExportAllHospitalPackage() ([]hp.Package, error) {
    query := `
        SELECT hp.*, (SELECT NVL(COUNT(ppd.PACKAGE_ID), 0) 
        FROM PATIENT_PURCHASE_DETAILS ppd 
        WHERE ppd.PACKAGE_ID = hp.PACKAGE_ID) AS PACKAGE_TOTAL_SOLD
        FROM HOSPITAL_PACKAGE hp
        ORDER BY hp.PACKAGE_NAME
    `
    query = strings.Replace(query, "*", utils.GetDbCols(hp.Package{}, "hp."), 1)
    list := make([]hp.Package, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        err := s.setHospitalPackage(&list[i])
        if err != nil {
            return nil, err
        }
    }
    return list, nil
}

func (s *ExportExcelService) ExportHospitalPackageByKeyword(keyword string, keyword2 string, keyword3 string) ([]hp.Package, error) {
    query := `
        SELECT hp.*, (SELECT NVL(COUNT(ppd.PACKAGE_ID), 0) 
        FROM PATIENT_PURCHASE_DETAILS ppd 
        WHERE ppd.PACKAGE_ID = hp.PACKAGE_ID) AS PACKAGE_TOTAL_SOLD
        FROM HOSPITAL_PACKAGE hp
    `
    query = strings.Replace(query, "*", utils.GetDbCols(hp.Package{}, "hp."), 1)
    lq := []string{query}
    lc := []string{}
    args := []any{}

    if keyword != "" {
        lc = append(lc, `(LOWER(hp.PACKAGE_CODE) LIKE :keyword OR LOWER(hp.PACKAGE_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        lc = append(lc, `TRUNC(hp.PACKAGE_START_DATETIME) = TO_DATE(:keyword2, 'dd/mm/yyyy')`)
        args = append(args, sql.Named("keyword2", keyword2))
    }
    if keyword3 != "" {
        lc = append(lc, `TRUNC(hp.PACKAGE_END_DATETIME) = TO_DATE(:keyword3, 'dd/mm/yyyy')`)
        args = append(args, sql.Named("keyword3", keyword3))
    }

    if len(lc) > 0 {
        s := strings.Join(lc, " AND ")
        lq = append(lq, fmt.Sprintf(" WHERE %s", s))
    }

    lq = append(lq, " ORDER BY hp.PACKAGE_NAME")
    query = strings.Join(lq, "")

    list := make([]hp.Package, 0)
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        err := s.setHospitalPackage(&list[i])
        if err != nil {
            return nil, err
        }
    }
    return list, nil
}

func (s *ExportExcelService) ExportAllPurchaseHistory() ([]upck.UserPackage, error) {
    query := `
        SELECT ppd.*, hp.PACKAGE_NAME, ppd2.PAYMENT_REQUEST_NO,
        ppd2.PAYMENT_REQUEST_CURRENCY, ppd2.PAYMENT_AMOUNT, ppd2.PAYMENT_CURRENCY,
        ppd2.PAYMENT_AMOUNT_COLLECTED, ppd2.PAYMENT_STATUS, ppd2.PAYMENT_TRANS_DATE,
        ppd2.BILLING_FULLNAME, ppd2.BILLING_CONTACT_NO, ppd2.BILLING_CONTACT_CODE, 
        ppd2.BILLING_EMAIL, ppd2.PAYMENT_URL
        FROM PATIENT_PURCHASE_DETAILS ppd
        JOIN HOSPITAL_PACKAGE hp ON ppd.PACKAGE_ID = hp.PACKAGE_ID
        JOIN PACKAGE_PAYMENT_DETAILS ppd2 ON ppd.PACKAGE_PAYMENT_ID = ppd2.PACKAGE_PAYMENT_ID
        WHERE ppd.PACKAGE_STATUS = 'Purchased'
        ORDER BY ppd.DATE_CREATE DESC, ppd.PACKAGE_PURCHASE_NO DESC
    `
    query = strings.Replace(query, "ppd.*", getPatientPurchaseDetailsCols("ppd."), 1)
    list := make([]upck.UserPackage, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
        s.setPurchaseHistory(&list[i])
    }
    return list, nil
}

func (s *ExportExcelService) ExportPurchaseHistoryByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string) ([]upck.UserPackage, error) {
    query := `
        SELECT ppd.*, hp.PACKAGE_NAME, ppd2.PAYMENT_REQUEST_NO,
        ppd2.PAYMENT_REQUEST_CURRENCY, ppd2.PAYMENT_AMOUNT, ppd2.PAYMENT_CURRENCY,
        ppd2.PAYMENT_AMOUNT_COLLECTED, ppd2.PAYMENT_STATUS, ppd2.PAYMENT_TRANS_DATE,
        ppd2.BILLING_FULLNAME, ppd2.BILLING_CONTACT_NO, ppd2.BILLING_CONTACT_CODE, 
        ppd2.BILLING_EMAIL, ppd2.PAYMENT_URL
        FROM PATIENT_PURCHASE_DETAILS ppd
        JOIN HOSPITAL_PACKAGE hp ON ppd.PACKAGE_ID = hp.PACKAGE_ID
        JOIN PACKAGE_PAYMENT_DETAILS ppd2 ON ppd.PACKAGE_PAYMENT_ID = ppd2.PACKAGE_PAYMENT_ID
    `
    query = strings.Replace(query, "ppd.*", getPatientPurchaseDetailsCols("ppd."), 1)
    lq := []string{query}
    lc := []string{}
    args := []any{}

    if keyword != "" {
        lc = append(lc, `LOWER(ppd.PATIENT_PRN) LIKE :keyword`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        lc = append(lc, `LOWER(ppd.PACKAGE_PURCHASE_NO) LIKE :keyword2`)
        args = append(args, sql.Named("keyword2", strings.ToLower(keyword2)))
    }
    if keyword3 != "" {
        lc = append(lc, `LOWER(hp.PACKAGE_NAME) LIKE :keyword3`)
        args = append(args, sql.Named("keyword3", strings.ToLower(keyword3)))
    }
    if keyword4 != "" && keyword4 != "All" {
        lc = append(lc, `LOWER(ppd.PACKAGE_STATUS) LIKE :keyword4`)
        args = append(args, sql.Named("keyword4", strings.ToLower(keyword4)))
    }

    if len(lc) > 0 {
        s := strings.Join(lc, " AND ")
        lq = append(lq, fmt.Sprintf(" WHERE %s", s))
    }

    lq = append(lq, " ORDER BY ppd.DATE_CREATE DESC, ppd.PACKAGE_PURCHASE_NO DESC")
    query = strings.Join(lq, "")

    list := make([]upck.UserPackage, 0)
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
        s.setPurchaseHistory(&list[i])
    }
    return list, nil
}

func (s *ExportExcelService) ExportAllLogisticRequest() ([]lg.LogisticRequest, error) {
    query := `
        SELECT lar.*
        FROM LOGISTIC_ARRANGEMENT_REQUESTER lar
        ORDER BY lar.FLIGHT_ARRIVAL_DATE
    `
    query = strings.Replace(query, "lar.*", utils.GetDbCols(lg.LogisticRequest{}, "lar."), 1)
    list := make([]lg.LogisticRequest, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].SetWebAdmin()
        s.setLogisticRequest(&list[i])
    }
    return list, nil
}

func (s *ExportExcelService) ExportLogisticRequestByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string) ([]lg.LogisticRequest, error) {
    query := `
        SELECT lar.* 
        FROM LOGISTIC_ARRANGEMENT_REQUESTER lar
    `
    query = strings.Replace(query, "lar.*", utils.GetDbCols(lg.LogisticRequest{}, "lar."), 1)
    lq := []string{query}
    lc := []string{}
    args := []any{}

    if keyword != "" {
        lc = append(lc, `(LOWER(lar.REQUESTER_PRN) LIKE :keyword OR LOWER(lar.REQUESTER_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        lc = append(lc, `LOWER(lar.PRIMARY_DOCTOR) LIKE :keyword2`)
        args = append(args, sql.Named("keyword2", strings.ToLower(keyword2)))
    }
    if keyword3 != "" {
        lc = append(lc, `TRUNC(lar.REQUESTED_PICKUP_DATE) = TO_DATE(:keyword3, 'dd/mm/yyyy')`)
        args = append(args, sql.Named("keyword3", keyword3))
    }
    if keyword4 != "" {
        lc = append(lc, `LOWER(lar.LOGISTIC_REQUEST_STATUS) LIKE :keyword4`)
        args = append(args, sql.Named("keyword4", strings.ToLower(keyword4)))
    }

    if len(lc) > 0 {
        s := strings.Join(lc, " AND ")
        lq = append(lq, fmt.Sprintf(" WHERE %s", s))
    }

    lq = append(lq, " ORDER BY lar.FLIGHT_ARRIVAL_DATE")
    query = strings.Join(lq, "")

    list := make([]lg.LogisticRequest, 0)
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].SetWebAdmin()
        s.setLogisticRequest(&list[i])
    }
    return list, nil
}

func (s *ExportExcelService) ExportAllMobileUserAuditLog() ([]model.MobileUserAuditLog, error) {
    query := `
        SELECT amu.* FROM AUDIT_MOBILE_USER amu
        WHERE ACTION = 'Delete Account'
        ORDER BY amu.DATE_CREATE DESC
    `
    query = strings.Replace(query, "amu.*", utils.GetDbCols(model.MobileUserAuditLog{}, "amu."), 1)
    list := make([]model.MobileUserAuditLog, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        s.setMobileUserAuditLog(&list[i])
    }
    return list, nil
}

func (s *ExportExcelService) ExportMobileUserAuditLogByKeyword(keyword string, keyword2 string) ([]model.MobileUserAuditLog, error) {
    query := `SELECT amu.* FROM AUDIT_MOBILE_USER amu`
    query = strings.Replace(query, "amu.*", utils.GetDbCols(model.MobileUserAuditLog{}, "amu."), 1)
    lq := []string{query}
    lc := []string{}
    args := []any{}

    if keyword != "" {
        lc = append(lc, `LOWER(amu.PRN) LIKE :keyword`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        lc = append(lc, `LOWER(amu.PATIENT_NAME) LIKE :keyword2`)
        args = append(args, sql.Named("keyword2", strings.ToLower(keyword2)))
    }
    lc = append(lc, `amu.ACTION = :action`)
    args = append(args, sql.Named("action", "Delete Account"))

    if len(lc) > 0 {
        s := strings.Join(lc, " AND ")
        lq = append(lq, fmt.Sprintf(" WHERE %s", s))
    }

    lq = append(lq, " ORDER BY amu.DATE_CREATE DESC")
    query = strings.Join(lq, "")

    list := make([]model.MobileUserAuditLog, 0)
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        s.setMobileUserAuditLog(&list[i])
    }
    return list, nil
}

// func (s *ExportExcelService)

// func (s *ExportExcelService)

// func (s *ExportExcelService)

// func (s *ExportExcelService)

// func (s *ExportExcelService)

// func (s *ExportExcelService)
