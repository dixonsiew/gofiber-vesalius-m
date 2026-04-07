package patientOutstandingBill

import (
    "context"
    "database/sql"
    "vesaliusm/database"
    ub "vesaliusm/model/userBilling"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var PatientOutstandingBillSvc *PatientOutstandingBillService = NewPatientOutstandingBillService(database.GetDb(), database.GetCtx())

type PatientOutstandingBillService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewPatientOutstandingBillService(db *sqlx.DB, ctx context.Context) *PatientOutstandingBillService {
    return &PatientOutstandingBillService{db: db, ctx: ctx}
}

func (s *PatientOutstandingBillService) FindPaidCountByBillInvoiceNumber(billNumber string, invoiceNumber string) (int, error) {
    query := `
        SELECT COUNT(*) AS COUNT
         FROM PATIENT_OUTSTANDING_BILL
         WHERE BILL_NUMBER = :billNumber
         AND INVOICE_NUMBER = :invoiceNumber
        AND PAYMENT_STATUS = 'Paid'
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, billNumber, invoiceNumber)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *PatientOutstandingBillService) GetOutstandingBillByBillPaymentId(conn *sqlx.DB, billPaymentId int64) ([]ub.PatientOutstandingBill, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT PATIENT_PRN, BILL_NUMBER, INVOICE_NUMBER, PAYMENT_STATUS
         FROM PATIENT_OUTSTANDING_BILL
        WHERE OUTSTANDING_BILL_PAYMENT_ID = :billPaymentId
    `
    list := make([]ub.PatientOutstandingBill, 0)
    err := db.SelectContext(s.ctx, &list, query, billPaymentId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *PatientOutstandingBillService) Save(billPaymentId int64, o ub.UserBilling, tx *sqlx.Tx) error {
    query := `
        INSERT INTO PATIENT_OUTSTANDING_BILL
        (
        OUTSTANDING_BILL_PAYMENT_ID, PATIENT_PRN, PATIENT_NAME, PATIENT_DOC_NUMBER, 
        REG_DATE_TIME, BILL_NUMBER, INVOICE_NUMBER, INVOICE_DATE_TIME, 
        INVOICE_AMOUNT, OUTSTANDING_AMOUNT, PAYMENT_STATUS, SUBMITTED_DATETIME
        ) VALUES (
        :bill_payment_id, :patientPrn, :patientName, :patientDocNumber,
        TO_DATE(:vesRegDateTime, 'DD-MM-YYYY HH24:MI'), :vesBillNumber, :vesInvoiceNumber, TO_DATE(:vesInvoiceDateTime, 'DD-MM-YYYY HH24:MI'),
        :vesInvoiceAmount, :vesOutstandingAmount, :paymentStatus, CURRENT_TIMESTAMP
        )
    `
    args := []any{
        sql.Named("bill_payment_id", billPaymentId),
        sql.Named("patientPrn", o.PatientPrn),
        sql.Named("patientName", o.PatientName),
        sql.Named("patientDocNumber", o.PatientDocNumber),
        sql.Named("vesRegDateTime", o.VesRegDateTime),
        sql.Named("vesBillNumber", o.VesBillNumber),
        sql.Named("vesInvoiceNumber", o.VesInvoiceNumber),
        sql.Named("vesInvoiceDateTime", o.VesInvoiceDateTime),
        sql.Named("vesInvoiceAmount", o.VesInvoiceAmount),
        sql.Named("vesOutstandingAmount", o.VesOutstandingAmount),
        sql.Named("paymentStatus", o.PaymentStatus),
    }
    var err error
    if tx == nil {
        _, err = s.db.ExecContext(s.ctx, query, args...)
    } else {
        _, err = tx.ExecContext(s.ctx, query, args...)
    }
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *PatientOutstandingBillService) UpdatePaymentStatusByBillPaymentId(billPaymentId int64, status string, tx *sqlx.Tx) error {
    query := `UPDATE PATIENT_OUTSTANDING_BILL SET PAYMENT_STATUS = :status, PAID_DATETIME = CURRENT_TIMESTAMP WHERE OUTSTANDING_BILL_PAYMENT_ID = :billPaymentId`
    var err error
    if tx == nil {
        _, err = s.db.ExecContext(s.ctx, query, status, billPaymentId)
    } else {
        _, err = tx.ExecContext(s.ctx, query, status, billPaymentId)
    }
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}
