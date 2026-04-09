package payment

import (
    ub "vesaliusm/model/userBilling"
    upck "vesaliusm/model/userPackage"
    "vesaliusm/service/packagePaymentDetails"
    "vesaliusm/service/billPaymentDetails"
)

type PaymentService struct {
    packagePaymentDetailsService *packagePaymentDetails.PackagePaymentDetailsService
    billPaymentDetailsService    *billPaymentDetails.BillPaymentDetailsService
}

func NewPaymentService() *PaymentService {
    return &PaymentService{
        packagePaymentDetailsService: packagePaymentDetails.,
        billPaymentDetailsService:    billPaymentDetails.NewBillPaymentDetailsService(),
    }
}

func (s *PaymentService) SaveBillingIPay(o ub.BillingPayment, o2 ub.UserBilling) error {
    return s.billPaymentDetailsService.SaveIPay(o, o2)
}

func (s *PaymentService) SaveIPay(o upck.PackagePaymentDetails, o2 []upck.UserPackage) error {
    return s.packagePaymentDetailsService.SaveIPay(o, o2)
}

