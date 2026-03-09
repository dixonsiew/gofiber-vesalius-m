package patientPurchaseDetails

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/model/userPackage"
	applicationuserService "vesaliusm/service/applicationUser"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

var applicationUserSvc *applicationuserService.ApplicationUserService = 
    applicationuserService.NewApplicationUserService(database.GetDb(), database.GetCtx())

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
        return nil, err
    }

    utils.LogInfo("hello")
    pager := model.GetPager(total, page, limit)
    list, err := s.FindByKeyword(keyword, keyword2, keyword3, keyword4, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
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

    var (
        rows *sqlx.Rows
        err error
    )
    if len(conditions) < 1 {
        rows, err = s.db.QueryxContext(s.ctx, query)
    } else {
        rows, err = s.db.NamedQueryContext(s.ctx, query, args)
    }
    
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    defer rows.Close()

    utils.LogInfo("get count")
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
        return nil, err
    }
    prn := user.MasterPrn.String
    total, err := s.CountByPrn(prn)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllByPrn(prn, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
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

func (s *PatientPurchaseDetailsService) FindByKeyword(keyword string, keyword2 string, keyword3 string, keyword4 string, offset int, limit int) ([]userPackage.UserPackage, error) {
    conditions, args := buildKeywordConditions(keyword, keyword2, keyword3, keyword4)
    args["offset"] = offset
    args["limit"] = limit

    base := `
        SELECT ` + getPatientPurchaseDetailsCols() + `, hp.PACKAGE_NAME, ppd2.PAYMENT_REQUEST_NO,
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

    rows, err := s.db.NamedQueryContext(s.ctx, query, args)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    var items []userPackage.UserPackage
    for rows.Next() {
        var item userPackage.UserPackage
        if err = rows.StructScan(&item); err != nil {
            utils.LogError(err)
            return nil, err
        }
        item.SetWebadmin()
        items = append(items, item)
    }
    return items, nil
}

func (s *PatientPurchaseDetailsService) FindAll(offset int, limit int) ([]userPackage.UserPackage, error) {
    query := `
        SELECT ` + getPatientPurchaseDetailsCols() + `, hp.PACKAGE_NAME, ppd2.PAYMENT_REQUEST_NO,
        ppd2.PAYMENT_REQUEST_CURRENCY, ppd2.PAYMENT_AMOUNT, ppd2.PAYMENT_CURRENCY,
        ppd2.PAYMENT_AMOUNT_COLLECTED, ppd2.PAYMENT_STATUS, ppd2.PAYMENT_TRANS_DATE,
        ppd2.BILLING_FULLNAME, ppd2.BILLING_CONTACT_NO, ppd2.BILLING_CONTACT_CODE, 
        ppd2.BILLING_EMAIL, ppd2.PAYMENT_URL
        FROM PATIENT_PURCHASE_DETAILS ppd
        JOIN HOSPITAL_PACKAGE hp ON ppd.PACKAGE_ID = hp.PACKAGE_ID
        JOIN PACKAGE_PAYMENT_DETAILS ppd2 ON ppd.PACKAGE_PAYMENT_ID = ppd2.PACKAGE_PAYMENT_ID
        WHERE ppd.PACKAGE_STATUS = 'Purchased'
        ORDER BY ppd.DATE_CREATE DESC, ppd.PACKAGE_PURCHASE_NO DESC
        OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY
    `
    var list []userPackage.UserPackage
    err := s.db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return list, err
    }
    for _, o := range list {
        o.SetWebadmin()
    }
    return list, nil
}

func (s *PatientPurchaseDetailsService) FindAllByPaymentId(paymentId int64) ([]userPackage.UserPackagePaymentEmail, error) {
    query := `
        SELECT 
          ppd.PACKAGE_ID, 
          COUNT(*) AS PACKAGE_QUANTITY,
          MAX(ppd.PATIENT_NAME) AS PATIENT_NAME,
          MAX(ppd2.PAYMENT_REQUEST_NO) AS PAYMENT_REQUEST_NO,
          MAX(ppd.PURCHASED_DATETIME) AS PURCHASED_DATETIME,
          MAX(hp.PACKAGE_NAME) AS PACKAGE_NAME,
          MAX(hp.PACKAGE_PRICE) AS PACKAGE_PRICE,
          MAX(ppd2.PAYMENT_GATEWAY) AS PAYMENT_GATEWAY,
          MAX(ppd.EXPIRED_DATETIME) AS EXPIRED_DATETIME,
          MAX(ppd2.BILLING_ADDRESS1) AS BILLING_ADDRESS1,
          MAX(ppd2.BILLING_ADDRESS2) AS BILLING_ADDRESS2,
          MAX(ppd2.BILLING_ADDRESS3) AS BILLING_ADDRESS3,
          MAX(ppd2.BILLING_TOWNCITY) AS BILLING_TOWNCITY,
          MAX(ppd2.BILLING_STATE) AS BILLING_STATE,
          MAX(ppd2.BILLING_POSTCODE) AS BILLING_POSTCODE,
          MAX(ppd2.BILLING_EMAIL) AS BILLING_EMAIL
        FROM 
        PATIENT_PURCHASE_DETAILS ppd
        JOIN PACKAGE_PAYMENT_DETAILS ppd2 ON ppd.PACKAGE_PAYMENT_ID = ppd2.PACKAGE_PAYMENT_ID
        JOIN HOSPITAL_PACKAGE hp ON ppd.PACKAGE_ID = hp.PACKAGE_ID
        WHERE ppd.PACKAGE_PAYMENT_ID = :1
        GROUP BY ppd.PACKAGE_ID
    `
    var list []userPackage.UserPackagePaymentEmail
    err := s.db.SelectContext(s.ctx, &list, query, paymentId)
    if err != nil {
        utils.LogError(err)
        return list, err
    }
    for _, o := range list {
        o.Set()
    }
    return list, nil
}

func (s *PatientPurchaseDetailsService) FindAllByPrn(prn string, offset int, limit int) ([]userPackage.UserPackage, error) {
    query := `
        SELECT 
          ppd.PACKAGE_PURCHASE_NO, 
          ppd.PACKAGE_STATUS,
          ppd.REDEEMED_DATETIME, 
          ppd.CANCELLED_DATETIME, 
          ppd.EXPIRED_DATETIME, 
          ppd.PURCHASED_DATETIME,
          hp.PACKAGE_ID, 
          hp.PACKAGE_NAME, 
          hp.PACKAGE_IMG, 
          hp.PACKAGE_ALLOW_APPT, 
          nd.MCR,
          ppd2.BILLING_FULLNAME, 
          ppd2.PAYMENT_TRANS_DATE, 
          ppd2.PAYMENT_AMOUNT_COLLECTED, 
          ppd2.PAYMENT_GATEWAY, 
          ppd2.PAYMENT_REQUEST_NO, 
          ndpa.DATE_APPT
        FROM PATIENT_PURCHASE_DETAILS ppd
        JOIN HOSPITAL_PACKAGE hp ON ppd.PACKAGE_ID = hp.PACKAGE_ID
        JOIN NOVA_DOCTOR nd ON hp.PACKAGE_ASSIGNED_DOCTOR = nd.DOCTOR_ID
        JOIN PACKAGE_PAYMENT_DETAILS ppd2 ON ppd.PACKAGE_PAYMENT_ID = ppd2.PACKAGE_PAYMENT_ID
        LEFT JOIN 
        (
          SELECT PACKAGE_PURCHASE_NO, MAX(DATE_APPT) AS DATE_APPT
          FROM NOVA_DOCTOR_PATIENT_APPT
          WHERE APPT_STATUS IN ('CONFIRMED', 'CHANGED')
          GROUP BY PACKAGE_PURCHASE_NO
        ) 
        ndpa ON ppd.PACKAGE_PURCHASE_NO = ndpa.PACKAGE_PURCHASE_NO
        WHERE ppd.PATIENT_PRN = :1
        ORDER BY ppd.DATE_CREATE DESC, ppd.PACKAGE_PURCHASE_NO DESC 
        OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY
    `
    var list []userPackage.UserPackage
    err := s.db.SelectContext(s.ctx, &list, query, prn, offset, limit)
    if err != nil {
        utils.LogError(err)
        return list, err
    }
    for _, o := range list {
        o.SetMobile()
    }
    return list, nil
}

func (s *PatientPurchaseDetailsService) FindByPurchaseId(purchaseId int64) (*userPackage.UserPackage, error) {
    var o userPackage.UserPackage
    query := `
        SELECT ` + getPatientPurchaseDetailsCols() + `, hp.PACKAGE_NAME, ` + getPackagePaymentDetailsCols() + `, ndpa.DATE_APPT
        FROM PATIENT_PURCHASE_DETAILS ppd
        JOIN HOSPITAL_PACKAGE hp ON ppd.PACKAGE_ID = hp.PACKAGE_ID
        JOIN PACKAGE_PAYMENT_DETAILS ppd2 ON ppd.PACKAGE_PAYMENT_ID = ppd2.PACKAGE_PAYMENT_ID
        LEFT JOIN
        (
        SELECT PACKAGE_PURCHASE_NO, MAX(DATE_APPT) AS DATE_APPT
        FROM NOVA_DOCTOR_PATIENT_APPT
        WHERE APPT_STATUS IN ('CONFIRMED', 'CHANGED')
        GROUP BY PACKAGE_PURCHASE_NO
        )
        ndpa ON ppd.PACKAGE_PURCHASE_NO = ndpa.PACKAGE_PURCHASE_NO
        WHERE ppd.PATIENT_PURCHASE_ID = :1
    `
    err := s.db.GetContext(s.ctx, &o, query, purchaseId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.SetWebadmin()
    return &o, nil
}

func (s *PatientPurchaseDetailsService) Save(payment_id int64, o userPackage.UserPackage) error {
    query := `
        INSERT INTO PATIENT_PURCHASE_DETAILS (
            PATIENT_PRN, PATIENT_NAME, PACKAGE_ID,
            PACKAGE_STATUS, PACKAGE_PAYMENT_ID, ORDERED_DATETIME
        ) VALUES (
            :patientPrn, :patientName, :package_id,
            :packageStatus, :payment_id, CURRENT_TIMESTAMP
        )`
    tx, err := s.db.BeginTxx(s.ctx, nil)
    if err != nil {
        utils.LogError(err)
        return err
    }
    defer tx.Rollback()

    for i := 0; i < o.QuantityPurchased; i++ {
        args := map[string]interface{}{
            "patientPrn":    o.PatientPrn,
            "patientName":   o.PatientName,
            "package_id":    o.PackageID,
            "packageStatus": o.PackageStatus,
            "payment_id":    payment_id,
        }

        _, err = tx.NamedExecContext(s.ctx, query, args)
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return tx.Commit()
}

func (s *PatientPurchaseDetailsService) SaveGuest(payment_id int64, o userPackage.UserPackage) error {
    return s.Save(payment_id, o)
}

func (s *PatientPurchaseDetailsService) UpdatePackageStatusByPurchaseNo(purchaseNo string, status string) error {
    var query string

    if status == utils.PackageStatusCancelled {
        query = `UPDATE PATIENT_PURCHASE_DETAILS SET PACKAGE_STATUS = :1 WHERE PACKAGE_PURCHASE_NO = :2`
        _, err := s.db.ExecContext(s.ctx, query, utils.PackageStatusPurchased, purchaseNo)
        if err != nil {
            utils.LogError(err)
            return err
        }
    } else {
        fieldMap := map[string]string{
            utils.PackageStatusPurchased: "PURCHASED_DATETIME",
            utils.PackageStatusBooked:    "BOOKED_DATETIME",
            utils.PackageStatusRedeemed:  "REDEEMED_DATETIME",
            utils.PackageStatusCancelled: "CANCELLED_DATETIME",
        }
        fieldDt, ok := fieldMap[status]
        if !ok {
            return fmt.Errorf("invalid status: %s", status)
        }

        query = fmt.Sprintf(`UPDATE PATIENT_PURCHASE_DETAILS SET PACKAGE_STATUS = :1, %s = CURRENT_TIMESTAMP WHERE PACKAGE_PURCHASE_NO = :2`, fieldDt)
        _, err := s.db.ExecContext(s.ctx, query, status, purchaseNo)
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func (s *PatientPurchaseDetailsService) UpdatePackageStatusByPaymentId(paymentId int64, status string) error {
    var query string

    if status == utils.PackageStatusCancelled {
        query = `UPDATE PATIENT_PURCHASE_DETAILS SET PACKAGE_STATUS = :1 WHERE PACKAGE_PAYMENT_ID = :2`
        _, err := s.db.ExecContext(s.ctx, query, utils.PackageStatusPurchased, paymentId)
        if err != nil {
            utils.LogError(err)
            return err
        }
    } else {
        fieldMap := map[string]string{
            utils.PackageStatusPurchased: "PURCHASED_DATETIME",
            utils.PackageStatusBooked:    "BOOKED_DATETIME",
            utils.PackageStatusRedeemed:  "REDEEMED_DATETIME",
            utils.PackageStatusCancelled: "CANCELLED_DATETIME",
        }
        fieldDt, ok := fieldMap[status]
        if !ok {
            return fmt.Errorf("invalid status: %s", status)
        }

        query = fmt.Sprintf(`UPDATE PATIENT_PURCHASE_DETAILS SET PACKAGE_STATUS = :1, %s = CURRENT_TIMESTAMP WHERE PACKAGE_PAYMENT_ID = :2`, fieldDt)
        _, err := s.db.ExecContext(s.ctx, query, status, paymentId)
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func (s *PatientPurchaseDetailsService) UpdatePackageStatusByPurchaseId(purchaseId int64, status string) error {
    var query string

    if status == utils.PackageStatusCancelled {
        query = `UPDATE PATIENT_PURCHASE_DETAILS SET PACKAGE_STATUS = :1 WHERE PATIENT_PURCHASE_ID = :2`
        _, err := s.db.ExecContext(s.ctx, query, utils.PackageStatusPurchased, purchaseId)
        if err != nil {
            utils.LogError(err)
            return err
        }
    } else {
        fieldMap := map[string]string{
            utils.PackageStatusPurchased: "PURCHASED_DATETIME",
            utils.PackageStatusBooked:    "BOOKED_DATETIME",
            utils.PackageStatusRedeemed:  "REDEEMED_DATETIME",
            utils.PackageStatusCancelled: "CANCELLED_DATETIME",
        }
        fieldDt, ok := fieldMap[status]
        if !ok {
            return fmt.Errorf("invalid status: %s", status)
        }

        query = fmt.Sprintf(`UPDATE PATIENT_PURCHASE_DETAILS SET PACKAGE_STATUS = :1, %s = CURRENT_TIMESTAMP WHERE PATIENT_PURCHASE_ID = :2`, fieldDt)
        _, err := s.db.ExecContext(s.ctx, query, status, purchaseId)
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func (s *PatientPurchaseDetailsService) GetAppointmentDetailsByPurchaseId(paymentId int64, status string) (*userPackage.ApptDetails, error) {
    var o userPackage.ApptDetails
    query := `
        SELECT ndpa.PATIENT_PRN, ndpa.PACKAGE_PURCHASE_NO, ndpa.APPT_NO FROM PATIENT_PURCHASE_DETAILS ppd
        JOIN NOVA_DOCTOR_PATIENT_APPT ndpa ON ppd.PACKAGE_PURCHASE_NO = ndpa.PACKAGE_PURCHASE_NO
        WHERE ndpa.APPT_STATUS <> 'CANCELLED' AND ppd.PATIENT_PURCHASE_ID = :1
    `
    err := s.db.GetContext(s.ctx, &o, query, paymentId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *PatientPurchaseDetailsService) GetPackageExpiryStatus(packageId int64) (string, error) {
    var r string
    query := `
        SELECT
          CASE
            WHEN PACKAGE_END_DATETIME < CURRENT_TIMESTAMP 
            THEN 'Expired'
            ELSE 'Not Expired'
          END AS PACKAGE_EXPIRY_STATUS
         FROM HOSPITAL_PACKAGE
         WHERE PACKAGE_ID = :1
    `
    err := s.db.GetContext(s.ctx, &r, query, packageId)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", err
        }
        utils.LogError(err)
        return "", err
    }
    return r, nil
}

func (s *PatientPurchaseDetailsService) GetPackageSoldoutStatus(packageId int64) (string, error) {
    var r string
    query := `
        SELECT
        CASE
          WHEN (SELECT COALESCE(COUNT(*), 0)
            FROM PATIENT_PURCHASE_DETAILS 
            WHERE PACKAGE_ID = :1) >=
            (SELECT PACKAGE_MAX_PURCHASE FROM HOSPITAL_PACKAGE WHERE PACKAGE_ID = :2)
          THEN 'Sold Out'
          ELSE 'Available'
        END AS PACKAGE_PURCHASE_AVAILABILITY
        FROM DUAL
    `
    err := s.db.GetContext(s.ctx, &r, query, packageId, packageId)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", err
        }
        utils.LogError(err)
        return "", err
    }
    return r, nil
}

func (s *PatientPurchaseDetailsService) GetPackageExceedPurchaseStatus(packageId int64, quantityPurchased int) (*userPackage.PackageExceedPurchaseStatus, error) {
    var result userPackage.PackageExceedPurchaseStatus
    query := `
        SELECT
          CASE
              WHEN ppd.TOTAL_PURCHASES + :1 > hp.PACKAGE_MAX_PURCHASE 
              THEN 'Exceeded'
              ELSE 'Not Exceeded'
          END AS PURCHASE_STATUS,
          CASE
              WHEN ppd.TOTAL_PURCHASES + :2 > hp.PACKAGE_MAX_PURCHASE 
              THEN GREATEST(hp.PACKAGE_MAX_PURCHASE - ppd.TOTAL_PURCHASES, 0)
              ELSE hp.PACKAGE_MAX_PURCHASE - ppd.TOTAL_PURCHASES
          END AS RECOMMENDED_QUANTITY
         FROM (
          SELECT COALESCE(SUM(1), 0) AS TOTAL_PURCHASES
          FROM PATIENT_PURCHASE_DETAILS 
          WHERE PACKAGE_ID = :3
         ) ppd,
         (
          SELECT PACKAGE_MAX_PURCHASE
          FROM HOSPITAL_PACKAGE
          WHERE PACKAGE_ID = :4
         ) hp
    `
    err := s.db.GetContext(s.ctx, &result, query, quantityPurchased, quantityPurchased, packageId, packageId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &result, nil
}

func (s *PatientPurchaseDetailsService) CheckPackageExpiryMaxPurchase(packageId int64, quantityPurchased int) (*userPackage.PackageCheckResult, error) {
    result := &userPackage.PackageCheckResult{
        PackageID: packageId,
        Expired:   0,
        Soldout:   0,
        ExceedPurchase: 0,
        RecommendedQuantity: 0,
    }
    expiry, err := s.GetPackageExpiryStatus(packageId)
    if err != nil {
        return nil, err
    }
    if expiry == "Expired" {
        result.Expired = 1
    }

    soldout, err := s.GetPackageSoldoutStatus(packageId)
    if err != nil {
        return nil, err
    }
    if soldout == "Sold Out" {
        result.Soldout = 1
    }

    exceedRes, err := s.GetPackageExceedPurchaseStatus(packageId, quantityPurchased)
    if err != nil {
        return nil, err
    }
    if exceedRes.PurchaseStatus == "Exceeded" {
        result.ExceedPurchase = 1
        result.RecommendedQuantity = exceedRes.RecommendedQuantity
    }

    return result, nil
}

func buildKeywordConditions(keyword string, keyword2 string, keyword3 string, keyword4 string) ([]string, map[string]interface{}) {
    var conds []string
    args := make(map[string]interface{})

    if keyword != "" {
        conds = append(conds, `LOWER(ppd.PATIENT_PRN) LIKE :keyword`)
        args["keyword"] = keyword
    }
    if keyword2 != "" {
        conds = append(conds, `LOWER(ppd.PACKAGE_PURCHASE_NO) LIKE :keyword2`)
        args["keyword2"] = keyword2
    }
    if keyword3 != "" {
        conds = append(conds, `LOWER(hp.PACKAGE_NAME) LIKE :keyword3`)
        args["keyword3"] = keyword3
    }
    if keyword4 != "" && keyword4 != "All" {
        conds = append(conds, `LOWER(ppd.PACKAGE_STATUS) LIKE :keyword4`)
        args["keyword4"] = keyword4
    }
    return conds, args
}

func whereClause(conds []string) string {
    if len(conds) == 0 {
        return ""
    }
    return " WHERE " + strings.Join(conds, " AND ")
}

func getPackagePaymentDetailsCols() string {
    return `
        ppd2.PAYMENT_GATEWAY,
        ppd2.PAYMENT_REQUEST_NO,
        ppd2.PAYMENT_REQUEST_CURRENCY,
        ppd2.PAYMENT_AMOUNT,
        ppd2.PAYMENT_CURRENCY,
        ppd2.PAYMENT_AMOUNT_COLLECTED,
        ppd2.PAYMENT_STATUS,
        ppd2.PAYMENT_TRANS_DATE,
        ppd2.BILLING_FULLNAME,
        ppd2.BILLING_CONTACT_NO,
        ppd2.BILLING_CONTACT_CODE,
        ppd2.BILLING_EMAIL,
        ppd2.PAYMENT_URL
    `
}

func getPatientPurchaseDetailsCols() string {
    return `
        ppd.PATIENT_PURCHASE_ID,
        ppd.PATIENT_PRN,
        ppd.PATIENT_NAME,
        ppd.PACKAGE_ID,
        ppd.PACKAGE_PURCHASE_NO,
        ppd.PACKAGE_STATUS,
        ppd.ORDERED_DATETIME,
        ppd.BOOKED_DATETIME,
        ppd.REDEEMED_DATETIME,
        ppd.CANCELLED_DATETIME,
        ppd.PURCHASED_DATETIME,
        ppd.EXPIRED_DATETIME
    `
}
