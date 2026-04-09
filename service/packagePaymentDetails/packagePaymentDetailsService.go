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
    "github.com/nleeper/goment"
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

func (s *PackagePaymentDetailsService) packagePaymentFindByRequestNo(conn *sqlx.DB, paymentRequestNo string) ([]upck.PackagePaymentDetails, error) {
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

func (s *PackagePaymentDetailsService) paymentFindByRequestId(conn *sqlx.DB, paymentRequestId string) ([]upck.PackagePaymentDetails, error) {
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

func (s *PackagePaymentDetailsService) SaveGuestIPay(o upck.PackagePaymentDetails, o2 []upck.UserPackage) error {
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
        err = s.patientPurchaseDetailsService.SaveGuest(ipaymentId, up, tx)
        if err != nil {
            return err
        }
    }
    return tx.Commit()
}

func (s *PackagePaymentDetailsService) SaveWallex(o upck.PackagePaymentDetails, o2 []upck.UserPackage) error {
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
            PAYMENT_GATEWAY, PAYMENT_REQUEST_ID, PAYMENT_REQUEST_NO,
            PAYMENT_REF_ID, PAYMENT_REQUEST_CURRENCY, PAYMENT_AMOUNT, PAYMENT_PURPOSE,
            PAYMENT_CURRENCY, PAYMENT_AMOUNT_COLLECTED, PAYMENT_REMARKS, PAYMENT_STATUS,
            PAYMENT_AUTH_CODE, PAYMENT_ERROR_DESC, PAYMENT_TRANS_DATE,
            BILLING_FULLNAME, BILLING_ADDRESS1, BILLING_ADDRESS2, BILLING_ADDRESS3,
            BILLING_TOWNCITY, BILLING_STATE, BILLING_POSTCODE, BILLING_COUNTRY_CODE,
            BILLING_CONTACT_NO, BILLING_CONTACT_CODE, BILLING_EMAIL, PAYMENT_URL
           ) VALUES (
            :paymentGateway, :paymentRequestId, :paymentRequestNo,
            :paymentRefId, :paymentRequestCurrency, :paymentAmount, :paymentPurpose,
            :paymentCurrency, :paymentAmountCollected, :paymentRemarks, :paymentStatus,
            :paymentAuthCode, :paymentErrorDesc, TO_DATE(:paymentTransDate, 'YYYY-MM-DD HH24:MI:SS'),
            :billingFullname, :billingAddress1, :billingAddress2, :billingAddress3,
            :billingTowncity, :billingState, :billingPostcode, :billingCountryCode,
            :billingContactNo, :billingContactCode, :billingEmail, :paymentUrl
        ) RETURNING PACKAGE_PAYMENT_ID INTO :payment_id
    `
    var paymentId go_ora.Number
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

func (s *PackagePaymentDetailsService) SaveGuestWallex(o upck.PackagePaymentDetails, o2 []upck.UserPackage) error {
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
            PAYMENT_GATEWAY, PAYMENT_REQUEST_ID, PAYMENT_REQUEST_NO,
            PAYMENT_REF_ID, PAYMENT_REQUEST_CURRENCY, PAYMENT_AMOUNT, PAYMENT_PURPOSE,
            PAYMENT_CURRENCY, PAYMENT_AMOUNT_COLLECTED, PAYMENT_REMARKS, PAYMENT_STATUS,
            PAYMENT_AUTH_CODE, PAYMENT_ERROR_DESC, PAYMENT_TRANS_DATE,
            BILLING_FULLNAME, BILLING_ADDRESS1, BILLING_ADDRESS2, BILLING_ADDRESS3,
            BILLING_TOWNCITY, BILLING_STATE, BILLING_POSTCODE, BILLING_COUNTRY_CODE,
            BILLING_CONTACT_NO, BILLING_CONTACT_CODE, BILLING_EMAIL, PAYMENT_URL
           ) VALUES (
            :paymentGateway, :paymentRequestId, :paymentRequestNo,
            :paymentRefId, :paymentRequestCurrency, :paymentAmount, :paymentPurpose,
            :paymentCurrency, :paymentAmountCollected, :paymentRemarks, :paymentStatus,
            :paymentAuthCode, :paymentErrorDesc, TO_DATE(:paymentTransDate, 'YYYY-MM-DD HH24:MI:SS'),
            :billingFullname, :billingAddress1, :billingAddress2, :billingAddress3,
            :billingTowncity, :billingState, :billingPostcode, :billingCountryCode,
            :billingContactNo, :billingContactCode, :billingEmail, :paymentUrl
        ) RETURNING PACKAGE_PAYMENT_ID INTO :payment_id
    `
    var paymentId go_ora.Number
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
        go_ora.Out{Dest: &paymentId},
    )
    if err != nil {
        return err
    }

    ipaymentId, _ := paymentId.Int64()
    for _, up := range o2 {
        err = s.patientPurchaseDetailsService.SaveGuest(ipaymentId, up, tx)
        if err != nil {
            return err
        }
    }
    return tx.Commit()
}

func (s *PackagePaymentDetailsService) SetWallexPaymentStatusToPaid(tx *sqlx.Tx, paymentAmount string, packagePaymentId int64, paymentRequestId string) error {
    query := `
        UPDATE PACKAGE_PAYMENT_DETAILS SET
        PAYMENT_STATUS = 'paid',
        PAYMENT_AMOUNT_COLLECTED = :amt,
        PAYMENT_TRANS_DATE = CURRENT_TIMESTAMP
        WHERE PACKAGE_PAYMENT_ID = :packagePaymentId
        AND PAYMENT_REQUEST_ID = :paymentRequestId
    `
    _, err := tx.ExecContext(s.ctx, query, paymentAmount, packagePaymentId, paymentRequestId)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *PackagePaymentDetailsService) SetIPayPackagePaymentStatusToPaid(tx *sqlx.Tx, paymentAmount string, packagePaymentId int64, paymentRequestNo string) error {
    query := `
        UPDATE PACKAGE_PAYMENT_DETAILS SET
        PAYMENT_STATUS = 'paid',
        PAYMENT_AMOUNT_COLLECTED = :amt,
        PAYMENT_TRANS_DATE = CURRENT_TIMESTAMP
        WHERE PACKAGE_PAYMENT_ID = :packagePaymentId
        AND PAYMENT_REQUEST_NO = :paymentRequestNo
    `
    _, err := tx.ExecContext(s.ctx, query, paymentAmount, packagePaymentId, paymentRequestNo)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *PackagePaymentDetailsService) UpdateWallexPaymentStatus(paymentRequestId string) error {
    lx, err := s.paymentFindByRequestId(s.db, paymentRequestId)
    if err != nil {
        return err
    }
    
    if len(lx) < 1 {
        return nil
    }

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

    r := lx[0]
    if r.PaymentRequestId.String == paymentRequestId {
        err = s.SetWallexPaymentStatusToPaid(tx, utils.GetAmount(r.PaymentAmount.Float64), r.PaymentId.Int64, r.PaymentRequestId.String)
        if err != nil {
            return nil
        }
        err = s.patientPurchaseDetailsService.UpdatePackageStatusByPaymentId(tx, r.PaymentId.Int64, utils.PackageStatusPurchased)
        if err != nil {
            return err
        }
        err = tx.Commit()
        if err != nil {
            utils.LogError(err)
            return err
        }
        err = s.sendPackagePayment(r.PaymentId.Int64)
        if err != nil {
            return err
        }
    }
    return nil
}

func (s *PackagePaymentDetailsService) UpdateIPayPaymentStatus(paymentRequestNo string) error {
    lx, err := s.packagePaymentFindByRequestNo(s.db, paymentRequestNo)
    if err != nil {
        return err
    }
    
    if len(lx) < 1 {
        return nil
    }
    
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
    
    r := lx[0]
    if r.PaymentRequestNo.String == paymentRequestNo && r.PaymentStatus.String == "unpaid" {
        err = s.SetIPayPackagePaymentStatusToPaid(tx, utils.GetAmount(r.PaymentAmount.Float64), r.PaymentId.Int64, r.PaymentRequestNo.String)
        if err != nil {
            return err
        }
        err = s.patientPurchaseDetailsService.UpdatePackageStatusByPaymentId(tx, r.PaymentId.Int64, utils.PackageStatusPurchased)
        if err != nil {
            return err
        }
        err = tx.Commit()
        if err != nil {
            utils.LogError(err)
            return err
        }
        err = s.sendPackagePayment(r.PaymentId.Int64)
        if err != nil {
            return err
        }
    }
    return nil
}

func (s *PackagePaymentDetailsService) sendPackagePayment(paymentId int64) error {
    patientPaymentPurchaseRes, err := s.patientPurchaseDetailsService.FindAllByPaymentId(paymentId)
    if err != nil {
        return err
    }
    
    for _, r2 := range patientPaymentPurchaseRes {
        lsaddr := make([]string, 0)
        if r2.BillingAddress1.Valid {
            lsaddr = append(lsaddr, r2.BillingAddress1.String)
        }
        if r2.BillingAddress2.Valid {
            lsaddr = append(lsaddr, r2.BillingAddress2.String)
        }
        if r2.BillingAddress3.Valid {
            lsaddr = append(lsaddr, r2.BillingAddress3.String)
        }
        if r2.BillingTowncity.Valid {
            lsaddr = append(lsaddr, r2.BillingTowncity.String)
        }
        if r2.BillingState.Valid {
            lsaddr = append(lsaddr, r2.BillingState.String)
        }
        if r2.BillingPostcode.Valid {
            lsaddr = append(lsaddr, r2.BillingPostcode.String)
        }
        
        addr := strings.Join(lsaddr, ", ")
        subtotalPrice := r2.PackagePrice.Float64 * float64(r2.PackageQuantity.Int32)

        dateOfPur, _ := goment.New(r2.PurchasedDateTime, "YYYY-MM-DD[T]HH:mm:ssZ")
        dateOfPurStr := dateOfPur.Format("DD-MMM-YYYY")

        packageExpiryDate, _ := goment.New(r2.ExpiredDateTime, "YYYY-MM-DD[T]HH:mm:ssZ")
        packageExpiryDateStr := packageExpiryDate.Format("DD-MMM-YYYY")

        email := ""
        if r2.BillingEmail.Valid {
            email = r2.BillingEmail.String
        }
        
        emailPrm := utils.Map{
            "patientName": r2.PatientName.String,
            "orderNumber": r2.PaymentRequestNo.String,
            "dateOfPurchase": dateOfPurStr,
            "productName": r2.PackageName.String,
            "productQuantity": r2.PackageQuantity.Int32,
            "productPrice": utils.GetAmount(r2.PackagePrice.Float64),
            "subtotalPrice": utils.GetAmount(subtotalPrice),
            "paymentMethod": r2.PaymentGateway,
            "totalPrice": utils.GetAmount(subtotalPrice),
            "packageExpiryDate": packageExpiryDateStr,
            "billingAddress": addr,
            "email": email,
        }
        go func() {
            s.mailService.SendSuccessPackagePayment(emailPrm)
        }()
    }
    return nil
}
