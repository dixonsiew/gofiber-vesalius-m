package vesaliusGeo

import "encoding/xml"

type OutstandingBill struct {
    XMLName             xml.Name `xml:"Bill" json:"-"`
    RegDateTime         string   `xml:"RegDateTime" json:"regDateTime"`
    BillNumber          string   `xml:"BillNumber" json:"billNumber"`
    InvoiceNumber       string   `xml:"InvoiceNumber" json:"invoiceNumber"`
    BillInvoiceDateTime string   `xml:"BillInvoiceDateTime" json:"billInvoiceDateTime"`
    BillAmount          string   `xml:"BillAmount" json:"billAmount"`
    InvoiceAmount       string   `xml:"InvoiceAmount" json:"invoiceAmount"`
    OutstandingAmount   string   `xml:"OutstandingAmount" json:"outstandingAmount"`
}
