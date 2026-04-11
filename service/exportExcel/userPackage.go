package exportExcel

import (
    "fmt"
    model "vesaliusm/model/userPackage"

    "github.com/nleeper/goment"
)

func (s *ExportExcelService) setPurchaseHistory(o *model.UserPackage) {
    if o.BillingContactCode.Valid && o.BillingContactNo.Valid {
        o.BillingFullContactExcel = fmt.Sprintf("%s%s", o.BillingContactCode.String, o.BillingContactNo.String) 
    }
    if o.PurchasedDateTime.Valid && o.PurchasedDateTime.String != "-" {
        g, _ := goment.New(o.PurchasedDateTime.String)
        o.PurchasedDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.PurchasedDateTimeExcel = o.PurchasedDateTime.String
    }
    if o.ExpiredDateTime.Valid && o.ExpiredDateTime.String != "-" {
        g, _ := goment.New(o.ExpiredDateTime.String)
        o.ExpiredDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.ExpiredDateTimeExcel = o.ExpiredDateTime.String
    }
    if o.OrderedDateTime.Valid && o.OrderedDateTime.String != "-" {
        g, _ := goment.New(o.OrderedDateTime.String)
        o.OrderedDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.OrderedDateTimeExcel = o.OrderedDateTime.String
    }
    if o.BookedDateTime.Valid && o.BookedDateTime.String != "-" {
        g, _ := goment.New(o.BookedDateTime.String)
        o.BookedDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.BookedDateTimeExcel = o.BookedDateTime.String
    }
    if o.RedeemedDateTime.Valid && o.RedeemedDateTime.String != "-" {
        g, _ := goment.New(o.RedeemedDateTime.String)
        o.RedeemedDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.RedeemedDateTimeExcel = o.RedeemedDateTime.String
    }
    if o.CancelledDateTime.Valid && o.CancelledDateTime.String != "-" {
        g, _ := goment.New(o.CancelledDateTime.String)
        o.CancelledDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.CancelledDateTimeExcel = o.CancelledDateTime.String
    }
    if o.PaymentTransDate.Valid && o.PaymentTransDate.String != "-" {
        g, _ := goment.New(o.PaymentTransDate.String)
        o.PaymentTransDateExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.PaymentTransDateExcel = o.PaymentTransDate.String
    }
}

func getPatientPurchaseDetailsCols(prefix string) string {
    lx := []string{
        "PATIENT_PURCHASE_ID",
        "PATIENT_PRN",
        "PATIENT_NAME",
        "PACKAGE_ID",
        "PACKAGE_PURCHASE_NO",
        "PACKAGE_STATUS",
        "ORDERED_DATETIME",
        "BOOKED_DATETIME",
        "REDEEMED_DATETIME",
        "CANCELLED_DATETIME",
        "PURCHASED_DATETIME",
        "EXPIRED_DATETIME",
    }
    ls := make([]string, 0)
    for _, v := range lx {
        ls = append(ls, prefix+v)
    }
    return strings.Join(ls, ", ")
}
