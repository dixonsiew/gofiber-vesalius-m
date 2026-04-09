package packagePaymentDetails

import (
    "context"
    "database/sql"
	"strings"
    "vesaliusm/database"
	upck "vesaliusm/model/userPackage"
	"vesaliusm/service/mail"
	"vesaliusm/service/patientPurchaseDetails"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
    go_ora "github.com/sijms/go-ora/v2"
)

var PackagePaymentDetailsSvc *PackagePaymentDetailsService = NewPackagePaymentDetailsService(database.GetDb(), database.GetCtx())

type PackagePaymentDetailsService struct {
    db  *sqlx.DB
    ctx context.Context
    mailService                   *mail.MailService
    patientPurchaseDetailsService *patientPurchaseDetails.PatientPurchaseDetailsService
}

func NewPackagePaymentDetailsService(db *sqlx.DB, ctx context.Context) *PackagePaymentDetailsService {
    return &PackagePaymentDetailsService{
        db:                            db,
        ctx:                           ctx,
        mailService:                   mail.MailSvc,
        patientPurchaseDetailsService: patientPurchaseDetails.PatientPurchaseDetailsSvc,
    }
}

func (s *PackagePaymentDetailsService) FindByRequestNo(requestNo string) (*upck.PackagePaymentDetails, error) {
    query := `SELECT * FROM PACKAGE_PAYMENT_DETAILS WHERE PAYMENT_REQUEST_NO = :requestNo`
    query = strings.Replace(query, "*", utils.GetDbCols(upck.PackagePaymentDetails{}, ""), 1)
    list := make([]upck.PackagePaymentDetails, 0)
    err := s.db.SelectContext(s.ctx, &list, query, requestNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &list[0], nil
}

func (s *PackagePaymentDetailsService) PackagePaymentFindByRequestNo(conn *sqlx.DB, paymentRequestNo string) ([]upck.PackagePaymentDetails, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT PACKAGE_PAYMENT_ID, PAYMENT_REQUEST_NO, PAYMENT_AMOUNT, PAYMENT_STATUS
        FROM PACKAGE_PAYMENT_DETAILS
        WHERE PAYMENT_REQUEST_NO = :paymentRequestNo
    `
    list := make([]upck.PackagePaymentDetails, 0)
    err := db.SelectContext(s.ctx, &list, query, paymentRequestNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *PackagePaymentDetailsService) PaymentFindByRequestId(conn *sqlx.DB, paymentRequestId string) ([]upck.PackagePaymentDetails, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT PACKAGE_PAYMENT_ID, PAYMENT_REQUEST_ID, PAYMENT_AMOUNT
        FROM PACKAGE_PAYMENT_DETAILS
        WHERE PAYMENT_REQUEST_ID = :paymentRequestId
    `
    list := make([]upck.PackagePaymentDetails, 0)
    err := db.SelectContext(s.ctx, &list, query, paymentRequestId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *PackagePaymentDetailsService) SaveIPay(o upck.PackagePaymentDetails, o2 []upck.UserPackage) error {
    var err error
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
        INSERT INTO PACKAGE_PAYMENT_DETAILS
           (
            PAYMENT_GATEWAY, PAYMENT_REQUEST_NO, PAYMENT_REQUEST_CURRENCY, PAYMENT_AMOUNT,
            PAYMENT_PURPOSE, PAYMENT_CURRENCY, PAYMENT_REMARKS, PAYMENT_STATUS,
            PAYMENT_AUTH_CODE, PAYMENT_ERROR_DESC, PAYMENT_TRANS_DATE,
            BILLING_FULLNAME, BILLING_ADDRESS1, BILLING_ADDRESS2, BILLING_ADDRESS3,
            BILLING_TOWNCITY, BILLING_STATE, BILLING_POSTCODE, BILLING_COUNTRY_CODE,
            BILLING_CONTACT_NO, BILLING_CONTACT_CODE, BILLING_EMAIL
           ) VALUES (
            :paymentGateway, :paymentRequestNo, :paymentRequestCurrency, :paymentAmount,
            :paymentPurpose, :paymentCurrency, :paymentRemarks, :paymentStatus,
            :paymentAuthCode, :paymentErrorDesc, TO_DATE(:paymentTransDate, 'YYYY-MM-DD HH24:MI:SS'),
            :billingFullname, :billingAddress1, :billingAddress2, :billingAddress3,
            :billingTowncity, :billingState, :billingPostcode, :billingCountryCode,
            :billingContactNo, :billingContactCode, :billingEmail
           ) RETURNING PACKAGE_PAYMENT_ID INTO :payment_id
        `
    var paymentId go_ora.Number
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
        go_ora.Out{Dest: &paymentId},
    )
    if err != nil {
        return err
    }
    
    ipaymentId, _ := paymentId.Int64()
    for _, up := range o2 {
        err = s.patientPurchaseDetailsService.Save(ipaymentId, up, tx)
        if err != nil {
            return err
        }
    }
    return tx.Commit()
}


