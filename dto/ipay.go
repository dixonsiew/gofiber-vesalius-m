package dto

type IPayPaymentResponseDto struct {
    MerchantCode string `json:"MerchantCode" validate:"required"`
    PaymentId    string `json:"PaymentId" validate:"required"`
    RefNo        string `json:"RefNo" validate:"required"`
    Amount       string `json:"Amount" validate:"required"`
    Currency     string `json:"Currency" validate:"required"`
    Remark       string `json:"Remark"`
    TransId      string `json:"TransId"`
    AuthCode     string `json:"AuthCode"`
    Status       string `json:"Status" validate:"required"`
    ErrDesc      string `json:"ErrDesc"`
    Signature    string `json:"Signature" validate:"required"`
    CCName       string `json:"CCName"`
    CCNo         string `json:"CCNo"`
    S_bankname   string `json:"S_bankname"`
    S_country    string `json:"S_country"`
}
