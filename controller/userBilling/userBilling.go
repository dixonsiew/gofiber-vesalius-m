package userBilling

import (
    "vesaliusm/service/billPaymentDetails"
    "vesaliusm/service/patientOutstandingBill"

    //"github.com/gofiber/fiber/v3"
)

type UserBillingController struct {
    billPaymentDetailsService *billPaymentDetails.BillPaymentDetailsService
    patientOutstandingBill    *patientOutstandingBill.PatientOutstandingBillService
}

func NewUserBillingController() *UserBillingController {
    return &UserBillingController{
        billPaymentDetailsService: billPaymentDetails.BillPaymentDetailsSvc,
        patientOutstandingBill:    patientOutstandingBill.PatientOutstandingBillSvc,
    }
}
