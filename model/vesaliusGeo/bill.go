package vesaliusGeo

import "encoding/xml"

type ResultOutstandingBills struct {
    XMLName xml.Name          `xml:"Result"`
    Name    string            `xml:"Name" json:"name"`
    PRN     string            `xml:"PRN" json:"prn"`
    Bills   []OutstandingBill `xml:"Bills"`
    Success Success           `xml:"Success"`
    Error   Error             `xml:"Error"`
}

type ResultPDFOutstandingBill struct {
    XMLName  xml.Name `xml:"Result"`
    BillData string   `xml:"BillData" json:"billData"`
    Success  Success  `xml:"Success"`
    Error    Error    `xml:"Error"`
}

type ResultBillPaymentConfirmation struct {
    XMLName       xml.Name `xml:"Result"`
    Name          string   `xml:"Name" json:"name"`
    PRN           string   `xml:"PRN" json:"prn"`
    Bill          string   `xml:"Bill" json:"billNumber"`
    Receipt       []string `xml:"Receipt" json:"receipt"`
    ReceiptNumber string   `json:"receiptNumber"`
    ReceiptAmount string   `json:"receiptAmount"`
    Success       Success  `xml:"Success"`
    Error         Error    `xml:"Error"`
}

func (o *ResultBillPaymentConfirmation) Set() {
    o.ReceiptNumber = o.Receipt[0]
    o.ReceiptAmount = o.Receipt[1]
}

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
