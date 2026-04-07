package billPaymentDetails

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    gm "vesaliusm/model"
    model "vesaliusm/model/userBilling"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/mail"
    "vesaliusm/service/patientOutstandingBill"
    "vesaliusm/service/vesaliusGeo"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    go_ora "github.com/sijms/go-ora/v2"
)

var BillPaymentDetailsSvc *BillPaymentDetailsService = NewBillPaymentDetailsService(database.GetDb(), database.GetCtx())

type BillPaymentDetailsService struct {
    db                            *sqlx.DB
    ctx                           context.Context
    applicationUserService        *applicationUser.ApplicationUserService
    mailService                   *mail.MailService
    patientOutstandingBillService *patientOutstandingBill.PatientOutstandingBillService
    vesaliusGeoService            *vesaliusGeo.VesaliusGeoService
}

func NewBillPaymentDetailsService(db *sqlx.DB, ctx context.Context) *BillPaymentDetailsService {
    return &BillPaymentDetailsService{
        db:                            db,
        ctx:                           ctx,
        applicationUserService:        applicationUser.ApplicationUserSvc,
        mailService:                   mail.MailSvc,
        patientOutstandingBillService: patientOutstandingBill.PatientOutstandingBillSvc,
        vesaliusGeoService:            vesaliusGeo.VesaliusGeoSvc,
    }
}

func (s *BillPaymentDetailsService) UpdateReceiptNoByRequestNo(conn *sqlx.DB, receiptNo string, requestNo string) error {
    db := database.GetFromCon(conn, s.db)
    query := `
        UPDATE OUTSTANDING_BILL_PAYMENT SET
        PAYMENT_RECEIPT_NO = :receiptNo,
        WS_CALL_FLAG = 'Y'
        WHERE PAYMENT_REQUEST_NO = :requestNo
    `
    _, err := db.ExecContext(s.ctx, query, receiptNo, requestNo)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *BillPaymentDetailsService) FindByRequestNo(requestNo string) ([]model.BillingPayment, error) {
    query := `SELECT * FROM OUTSTANDING_BILL_PAYMENT WHERE PAYMENT_REQUEST_NO = :requestNo`
    query = strings.Replace(query, "*", utils.GetDbCols(model.BillingPayment{}, ""), 1)
    list := make([]model.BillingPayment, 0)
    err := s.db.SelectContext(s.ctx, &list, query, requestNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *BillPaymentDetailsService) BillPaymentFindByRequestNo(conn *sqlx.DB, paymentRequestNo string) ([]model.BillingPayment, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT 
        OUTSTANDING_BILL_PAYMENT_ID, 
        PAYMENT_GATEWAY, 
        PAYMENT_REQUEST_NO, 
        PAYMENT_AMOUNT, 
        PAYMENT_STATUS, 
        BILLING_EMAIL,
        BILLING_FULLNAME
        FROM OUTSTANDING_BILL_PAYMENT
        WHERE PAYMENT_REQUEST_NO = :paymentRequestNo
    `
    list := make([]model.BillingPayment, 0)
    err := db.SelectContext(s.ctx, &list, query, paymentRequestNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *BillPaymentDetailsService) SetIPayBillPaymentStatusToPaid(tx *sqlx.Tx, paymentAmount float64, outstandingBillPaymentId int64, paymentRequestNo string) error {
    query := `
        UPDATE OUTSTANDING_BILL_PAYMENT SET
        PAYMENT_STATUS = 'paid',
        PAYMENT_AMOUNT_COLLECTED = :paymentAmount,
        PAYMENT_TRANS_DATE = CURRENT_TIMESTAMP
        WHERE OUTSTANDING_BILL_PAYMENT_ID = :outstandingBillPaymentId
        AND PAYMENT_REQUEST_NO = :paymentRequestNo
    `
    args := []any{
        sql.Named("paymentAmount", paymentAmount),
        sql.Named("outstandingBillPaymentId", outstandingBillPaymentId),
        sql.Named("paymentRequestNo", paymentRequestNo),
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

func (s *BillPaymentDetailsService) PaymentFindByRequestId(conn *sqlx.DB, paymentRequestId string) ([]model.BillingPayment, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT 
        OUTSTANDING_BILL_PAYMENT_ID, 
        PAYMENT_GATEWAY, 
        PAYMENT_REQUEST_ID, 
        PAYMENT_REQUEST_NO, 
        PAYMENT_AMOUNT, 
        PAYMENT_STATUS, 
        BILLING_EMAIL,
        BILLING_FULLNAME
        FROM OUTSTANDING_BILL_PAYMENT
        WHERE PAYMENT_REQUEST_ID = :paymentRequestId
    `
    list := make([]model.BillingPayment, 0)
    err := db.SelectContext(s.ctx, &list, query, paymentRequestId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *BillPaymentDetailsService) SetWallexPaymentStatusToPaid(tx *sqlx.Tx, paymentAmount float64, outstandingBillPaymentId int64, paymentRequestId string) error {
    query := `
        UPDATE OUTSTANDING_BILL_PAYMENT SET
        PAYMENT_STATUS = 'paid',
        PAYMENT_AMOUNT_COLLECTED = :paymentAmount,
        PAYMENT_TRANS_DATE = CURRENT_TIMESTAMP
        WHERE OUTSTANDING_BILL_PAYMENT_ID = :outstandingBillPaymentId
        AND PAYMENT_REQUEST_ID = :paymentRequestId
    `
    args := []any{
        sql.Named("paymentAmount", paymentAmount),
        sql.Named("outstandingBillPaymentId", outstandingBillPaymentId),
        sql.Named("paymentRequestId", paymentRequestId),
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

func (s *BillPaymentDetailsService) SaveIPay(o model.BillingPayment, o2 model.UserBilling) error {
    tx, err := s.db.BeginTxx(s.ctx, nil)
    if err != nil {
        utils.LogError(err)
        return err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()

    query := `
        INSERT INTO OUTSTANDING_BILL_PAYMENT
           (
            PAYMENT_GATEWAY, PAYMENT_REQUEST_NO, PAYMENT_REQUEST_CURRENCY, PAYMENT_AMOUNT,
            PAYMENT_PURPOSE, PAYMENT_CURRENCY, PAYMENT_REMARKS, PAYMENT_STATUS,
            PAYMENT_AUTH_CODE, PAYMENT_ERROR_DESC, PAYMENT_TRANS_DATE,
            BILLING_FULLNAME, BILLING_ADDRESS1, BILLING_ADDRESS2, BILLING_ADDRESS3,
            BILLING_TOWNCITY, BILLING_STATE, BILLING_POSTCODE, BILLING_COUNTRY_CODE,
            BILLING_CONTACT_NO, BILLING_CONTACT_CODE, BILLING_EMAIL, PAYMENT_RECEIPT_NO
           ) VALUES (
            :paymentGateway, :paymentRequestNo, :paymentRequestCurrency, :paymentAmount,
            :paymentPurpose, :paymentCurrency, :paymentRemarks, :paymentStatus,
            :paymentAuthCode, :paymentErrorDesc, TO_DATE(:paymentTransDate, 'YYYY-MM-DD HH24:MI:SS'),
            :billingFullname, :billingAddress1, :billingAddress2, :billingAddress3,
            :billingTowncity, :billingState, :billingPostcode, :billingCountryCode,
            :billingContactNo, :billingContactCode, :billingEmail, :vesPaymentReceiptNo
        ) RETURNING OUTSTANDING_BILL_PAYMENT_ID INTO :bill_payment_id
    `
    var billPaymentId go_ora.Number
    _, err = tx.ExecContext(s.ctx, query,
        sql.Named("paymentGateway", o.PaymentGateway.Int32),
        sql.Named("paymentRequestNo", o.PaymentRequestNo.String),
        sql.Named("paymentRequestCurrency", o.PaymentRequestCurrency.String),
        sql.Named("paymentAmount", o.PaymentAmount.Float64),
        sql.Named("paymentPurpose", o.PaymentPurpose.String),
        sql.Named("paymentCurrency", o.PaymentCurrency.String),
        sql.Named("paymentRemarks", o.PaymentRemarks.String),
        sql.Named("paymentStatus", o.PaymentStatus.String),
        sql.Named("paymentAuthCode", o.PaymentAuthCode.String),
        sql.Named("paymentErrorDesc", o.PaymentErrorDesc.String),
        sql.Named("paymentTransDate", o.PaymentTransDate.String),
        sql.Named("billingFullname", o.BillingFullname.String),
        sql.Named("billingAddress1", o.BillingAddress1.String),
        sql.Named("billingAddress2", o.BillingAddress2.String),
        sql.Named("billingAddress3", o.BillingAddress3.String),
        sql.Named("billingTowncity", o.BillingTowncity.String),
        sql.Named("billingState", o.BillingState.String),
        sql.Named("billingPostcode", o.BillingPostcode.String),
        sql.Named("billingCountryCode", o.BillingCountryCode.String),
        sql.Named("billingContactNo", o.BillingContactNo.String),
        sql.Named("billingContactCode", o.BillingContactCode.String),
        sql.Named("billingEmail", o.BillingEmail.String),
        sql.Named("vesPaymentReceiptNo", o.VesPaymentReceiptNo),
        go_ora.Out{Dest: &billPaymentId},
    )
    if err != nil {
        utils.LogError(err)
        return err
    }

    ibillPaymentId, _ := billPaymentId.Int64()
    err = s.patientOutstandingBillService.Save(ibillPaymentId, o2, tx)
    if err != nil {
        return err
    }

    err = tx.Commit()
    if err != nil {
        utils.LogError(err)
        return err
    }
    return err
}

func (s *BillPaymentDetailsService) SaveWallex(o model.BillingPayment, o2 model.UserBilling) error {
    tx, err := s.db.BeginTxx(s.ctx, nil)
    if err != nil {
        utils.LogError(err)
        return err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()

    query := `
        INSERT INTO OUTSTANDING_BILL_PAYMENT
        (
            PAYMENT_GATEWAY, PAYMENT_REQUEST_ID, PAYMENT_REQUEST_NO,
            PAYMENT_REF_ID, PAYMENT_REQUEST_CURRENCY, PAYMENT_AMOUNT, PAYMENT_PURPOSE,
            PAYMENT_CURRENCY, PAYMENT_AMOUNT_COLLECTED, PAYMENT_REMARKS, PAYMENT_STATUS,
            PAYMENT_AUTH_CODE, PAYMENT_ERROR_DESC, PAYMENT_TRANS_DATE,
            BILLING_FULLNAME, BILLING_ADDRESS1, BILLING_ADDRESS2, BILLING_ADDRESS3,
            BILLING_TOWNCITY, BILLING_STATE, BILLING_POSTCODE, BILLING_COUNTRY_CODE,
            BILLING_CONTACT_NO, BILLING_CONTACT_CODE, BILLING_EMAIL, PAYMENT_URL, PAYMENT_RECEIPT_NO
        ) VALUES (
            :paymentGateway, :paymentRequestId, :paymentRequestNo,
            :paymentRefId, :paymentRequestCurrency, :paymentAmount, :paymentPurpose,
            :paymentCurrency, :paymentAmountCollected, :paymentRemarks, :paymentStatus,
            :paymentAuthCode, :paymentErrorDesc, TO_DATE(:paymentTransDate, 'YYYY-MM-DD HH24:MI:SS'),
            :billingFullname, :billingAddress1, :billingAddress2, :billingAddress3,
            :billingTowncity, :billingState, :billingPostcode, :billingCountryCode,
            :billingContactNo, :billingContactCode, :billingEmail, :paymentUrl, :vesPaymentReceiptNo
        ) RETURNING OUTSTANDING_BILL_PAYMENT_ID INTO :bill_payment_id
    `
    var billPaymentId go_ora.Number
    _, err = tx.ExecContext(s.ctx, query,
        sql.Named("paymentGateway", o.PaymentGateway.Int32),
        sql.Named("paymentRequestId", o.PaymentRequestId.String),
        sql.Named("paymentRequestNo", o.PaymentRequestNo.String),
        sql.Named("paymentRefId", o.PaymentRefId.String),
        sql.Named("paymentRequestCurrency", o.PaymentRequestCurrency.String),
        sql.Named("paymentAmount", o.PaymentAmount.Float64),
        sql.Named("paymentPurpose", o.PaymentPurpose.String),
        sql.Named("paymentCurrency", o.PaymentCurrency.String),
        sql.Named("paymentAmountCollected", o.PaymentAmountCollected.Float64),
        sql.Named("paymentRemarks", o.PaymentRemarks.String),
        sql.Named("paymentStatus", o.PaymentStatus.String),
        sql.Named("paymentAuthCode", o.PaymentAuthCode.String),
        sql.Named("paymentErrorDesc", o.PaymentErrorDesc.String),
        sql.Named("paymentTransDate", o.PaymentTransDate.String),
        sql.Named("billingFullname", o.BillingFullname.String),
        sql.Named("billingAddress1", o.BillingAddress1.String),
        sql.Named("billingAddress2", o.BillingAddress2.String),
        sql.Named("billingAddress3", o.BillingAddress3.String),
        sql.Named("billingTowncity", o.BillingTowncity.String),
        sql.Named("billingState", o.BillingState.String),
        sql.Named("billingPostcode", o.BillingPostcode.String),
        sql.Named("billingCountryCode", o.BillingCountryCode.String),
        sql.Named("billingContactNo", o.BillingContactNo.String),
        sql.Named("billingContactCode", o.BillingContactCode.String),
        sql.Named("billingEmail", o.BillingEmail.String),
        sql.Named("paymentUrl", o.PaymentUrl),
        sql.Named("vesPaymentReceiptNo", o.VesPaymentReceiptNo),
        go_ora.Out{Dest: &billPaymentId},
    )
    if err != nil {
        utils.LogError(err)
        return err
    }

    ibillPaymentId, _ := billPaymentId.Int64()
    err = s.patientOutstandingBillService.Save(ibillPaymentId, o2, tx)
    if err != nil {
        return err
    }

    err = tx.Commit()
    if err != nil {
        utils.LogError(err)
        return err
    }
    return err
}

func (s *BillPaymentDetailsService) UpdateWallexPaymentStatus(paymentRequestId string) error {
    billPaymentRes, err := s.PaymentFindByRequestId(s.db, paymentRequestId)
    if err != nil {
        return err
    }

    if len(billPaymentRes) > 0 {
        r := billPaymentRes[0]
        if r.PaymentRequestId.String == paymentRequestId && r.PaymentStatus.String == "unpaid" {
            tx, err := s.db.BeginTxx(s.ctx, nil)
            if err != nil {
                utils.LogError(err)
                return err
            }
            defer func() {
                if err != nil {
                    utils.LogError(err)
                    tx.Rollback()
                }
            }()

            err = s.SetWallexPaymentStatusToPaid(tx, r.PaymentAmount.Float64, r.OutstandingBillPaymentID.Int64, r.PaymentRequestId.String)
            if err != nil {
                return err
            }

            err = s.patientOutstandingBillService.UpdatePaymentStatusByBillPaymentId(r.OutstandingBillPaymentID.Int64, utils.PaymentStatusPaid, tx)
            if err != nil {
                return err
            }

            err = tx.Commit()
            if err != nil {
                utils.LogError(err)
                return err
            }

            patientOutstandingBillRes, err := s.patientOutstandingBillService.GetOutstandingBillByBillPaymentId(s.db, r.OutstandingBillPaymentID.Int64)
            if err != nil {
                return err
            }

            r2 := patientOutstandingBillRes[0]

            emailPrm := utils.Map{
                "amount":        r.PaymentAmount.Float64,
                "paymentMethod": "Wallex",
                "billNumber":    "",
                "invoiceNumber": "",
                "email":         "",
            }

            if r.BillingEmail.Valid {
                emailPrm["email"] = r.BillingEmail.String
                if r2.PaymentStatus.String == utils.PaymentStatusPaid {
                    emailPrm["billNumber"] = r2.BillNumber.String
                    emailPrm["invoiceNumber"] = r2.InvoiceNumber.String
                    go func() {
                        s.mailService.SendSuccessOutstandingBillPayment(emailPrm)
                    }()
                }
            }

            wsRes, _, err := s.vesaliusGeoService.PatientProcessPatientBillPayment(r2.PatientPRN.String, r2.BillNumber.String, "WALLEX", 
                utils.GetAmount(r.PaymentAmount.Float64), r.PaymentRequestNo.String, "Wallex Online Payment from Bill", 
                r.BillingFullname.String)
            if err != nil {
                return err
            }

            if r2.BillNumber.String == wsRes.Bill {
                err = s.UpdateReceiptNoByRequestNo(s.db, wsRes.ReceiptNumber, r.PaymentRequestNo.String)
                if err != nil {
                    return err
                }
            }

            return err
        }
    }
    return nil
}

func (s *BillPaymentDetailsService) UpdateIPayPaymentStatus(paymentRequestNo string) error {
    billPaymentRes, err := s.BillPaymentFindByRequestNo(s.db, paymentRequestNo)
    if err != nil {
        return err
    }

    if len(billPaymentRes) > 0 {
        r := billPaymentRes[0]
        if r.PaymentRequestNo.String == paymentRequestNo && r.PaymentStatus.String == "unpaid" {
            tx, err := s.db.BeginTxx(s.ctx, nil)
            if err != nil {
                utils.LogError(err)
                return err
            }
            defer func() {
                if err != nil {
                    utils.LogError(err)
                    tx.Rollback()
                }
            }()

            err = s.SetIPayBillPaymentStatusToPaid(tx, r.PaymentAmount.Float64, r.OutstandingBillPaymentID.Int64, r.PaymentRequestNo.String)
            if err != nil {
                return err
            }

            err = s.patientOutstandingBillService.UpdatePaymentStatusByBillPaymentId(r.OutstandingBillPaymentID.Int64, utils.PaymentStatusPaid, tx)
            if err != nil {
                return err
            }

            err = tx.Commit()
            if err != nil {
                utils.LogError(err)
                return err
            }

            patientOutstandingBillRes, err := s.patientOutstandingBillService.GetOutstandingBillByBillPaymentId(s.db, r.OutstandingBillPaymentID.Int64)
            if err != nil {
                return err
            }

            r2 := patientOutstandingBillRes[0]

            if r2.PaymentStatus.String == utils.PaymentStatusPaid {
                emailPrm := utils.Map{
                    "amount":        r.PaymentAmount.Float64,
                    "paymentMethod": "iPay88",
                    "billNumber":    "",
                    "invoiceNumber": "",
                    "email":         "",
                }

                if r.BillingEmail.Valid {
                    emailPrm["email"] = r.BillingEmail.String
                    if r2.PaymentStatus.String == utils.PaymentStatusPaid {
                        emailPrm["billNumber"] = r2.BillNumber.String
                        emailPrm["invoiceNumber"] = r2.InvoiceNumber.String
                        go func() {
                            s.mailService.SendSuccessOutstandingBillPayment(emailPrm)
                        }()
                    }
                }

                wsRes, _, err := s.vesaliusGeoService.PatientProcessPatientBillPayment(r2.PatientPRN.String, r2.BillNumber.String, "IPAY88", 
                    utils.GetAmount(r.PaymentAmount.Float64), r.PaymentRequestNo.String, "iPay88 Online Payment from Bill", 
                    r.BillingFullname.String)
                if err != nil {
                    return err
                }

                if r2.BillNumber.String == wsRes.Bill {
                    err = s.UpdateReceiptNoByRequestNo(s.db, wsRes.ReceiptNumber, r.PaymentRequestNo.String)
                    if err != nil {
                        return err
                    }
                }
            }
        }
    }
    return nil
}

func (s *BillPaymentDetailsService) FindAllPaidByPrn(prn string, offset int, limit int, conn *sqlx.DB) ([]model.UserBillingPayment, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT pob.PATIENT_PRN, pob.PATIENT_NAME, pob.PATIENT_DOC_NUMBER, pob.BILL_NUMBER,
        pob.INVOICE_DATE_TIME, pob.INVOICE_AMOUNT, pob.OUTSTANDING_AMOUNT,
        obp.PAYMENT_GATEWAY, obp.PAYMENT_REQUEST_NO, obp.PAYMENT_AMOUNT_COLLECTED,
        obp.PAYMENT_TRANS_DATE, obp.BILLING_FULLNAME, obp.PAYMENT_RECEIPT_NO
        FROM PATIENT_OUTSTANDING_BILL pob
        JOIN OUTSTANDING_BILL_PAYMENT obp ON pob.OUTSTANDING_BILL_PAYMENT_ID = obp.OUTSTANDING_BILL_PAYMENT_ID
        WHERE pob.PATIENT_PRN = :prn AND pob.PAYMENT_STATUS = 'Paid'
        ORDER BY pob.DATE_CREATE DESC, obp.PAYMENT_REQUEST_NO DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    list := make([]model.UserBillingPayment, 0)
    err := db.SelectContext(s.ctx, &list, query, prn, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
}

func (s *BillPaymentDetailsService) ListPaidByPrn(userId int64, page string, limit string) (*gm.PagedList, error) {
    user, err := s.applicationUserService.FindByUserId(userId, s.db)
    if err != nil {
        return nil, err
    }
    
    prn := user.MasterPrn.String
    total, err := s.CountPaidByPrn(prn, s.db)
    if err != nil {
        return nil, err
    }

    pager := gm.GetPager(total, page, limit)
    list, err := s.FindAllPaidByPrn(prn, pager.GetLowerBound(), pager.PageSize, s.db)
    if err != nil {
        return nil, err
    }

    return &gm.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *BillPaymentDetailsService) CountPaidByPrn(prn string, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(PATIENT_OUTSTANDING_BILL_ID) AS COUNT FROM PATIENT_OUTSTANDING_BILL WHERE PATIENT_PRN = :prn AND PAYMENT_STATUS = 'Paid'`
    var count int
    err := db.GetContext(s.ctx, &count, query, prn)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}
