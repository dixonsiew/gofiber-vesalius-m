package payment

import (
    "context"
    "vesaliusm/database"
    ub "vesaliusm/model/userBilling"
    upck "vesaliusm/model/userPackage"
    "vesaliusm/service/packagePaymentDetails"
    "vesaliusm/service/billPaymentDetails"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var PaymentSvc *PaymentService = NewPaymentService(database.GetDb(), database.GetCtx())

type PaymentService struct {
    db  *sqlx.DB
    ctx context.Context
    packagePaymentDetailsService *packagePaymentDetails.PackagePaymentDetailsService
    billPaymentDetailsService    *billPaymentDetails.BillPaymentDetailsService
}

func NewPaymentService(db *sqlx.DB, ctx context.Context) *PaymentService {
    return &PaymentService{
        db:                            db,
        ctx:                           ctx,
        packagePaymentDetailsService: packagePaymentDetails.PackagePaymentDetailsSvc,
        billPaymentDetailsService:    billPaymentDetails.BillPaymentDetailsSvc,
    }
}

func (s *PaymentService) SaveBillingIPay(o ub.BillingPayment, o2 ub.UserBilling) error {
    return s.billPaymentDetailsService.SaveIPay(o, o2)
}

func (s *PaymentService) SaveIPay(o upck.PackagePaymentDetails, o2 []upck.UserPackage) error {
    return s.packagePaymentDetailsService.SaveIPay(o, o2)
}

func (s *PaymentService) SaveGuestIPay(o upck.PackagePaymentDetails, o2 []upck.UserPackage) error {
    return s.packagePaymentDetailsService.SaveGuestIPay(o, o2)
}

func (s *PaymentService) SaveBillingWallex(o ub.BillingPayment, o2 ub.UserBilling) error {
    return s.billPaymentDetailsService.SaveWallex(o, o2)
}

func (s *PaymentService) SaveWallex(o upck.PackagePaymentDetails, o2 []upck.UserPackage) error {
    return s.packagePaymentDetailsService.SaveWallex(o, o2)
}

func (s *PaymentService) SaveGuestWallex(o upck.PackagePaymentDetails, o2 []upck.UserPackage) error {
    return s.packagePaymentDetailsService.SaveGuestWallex(o, o2)
}

func (s *PaymentService) GeneratePaymentRefNo() (string, error) {
    var refId string
    query := `SELECT GENERATE_PAYMENT_REF_ID() AS PAYMENT_REF_ID FROM DUAL`
    err := s.db.GetContext(s.ctx, &refId, query)
    if err != nil {
        utils.LogError(err)
        return "", err
    }
    return refId, nil
}

func (s *PaymentService) GenerateBillPaymentRefNo() (string, error) {
    var refId string
    query := `SELECT GENERATE_BILL_PAYMENT_REF_ID() AS BILL_PAYMENT_REF_ID FROM DUAL`
    err := s.db.GetContext(s.ctx, &refId, query)
    if err != nil {
        utils.LogError(err)
        return "", err
    }
    return refId, nil
}
