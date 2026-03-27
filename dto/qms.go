package dto

type CustDto struct {
    CustId            string `json:"CustId" validate:"required"`
    MrnNumber         string `json:"MrnNumber" validate:"required"`
    IcNumber          string `json:"IcNumber" validate:"required"`
    Name              string `json:"Name" validate:"required"`
    MobilePhoneNumber string `json:"MobilePhoneNumber" validate:"required"`
}

type IndServiceDto struct {
    ServiceNo   string `json:"ServiceNo" validate:"required"`
    ServiceName string `json:"ServiceName" validate:"required"`
    ShortName   string `json:"ShortName" validate:"required"`
}

type TicketInfoDto struct {
    IssuedTime        string `json:"IssuedTime" validate:"required"`
    TotalWaitingAhead string `json:"TotalWaitingAhead" validate:"required"`
}

type QMSServerWebhookDto struct {
    Method     string        `json:"Method" validate:"required"`
    TicketStr  string        `json:"TicketStr" validate:"required"`
    Cust       CustDto       `json:"Cust" validate:"required"`
    IndService IndServiceDto `json:"IndService" validate:"required"`
    TicketInfo TicketInfoDto `json:"TicketInfo" validate:"required"`
}

type QMSPrnDto struct {
    Prn string `json:"prn" validate:"required"`
}
