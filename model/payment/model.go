package payment

type CustomerInfo struct {
    Name             string `json:"name"`
    ItuTelephoneCode string `json:"ituTelephoneCode"`
    MobileNumber     string `json:"mobileNumber"`
    Email            string `json:"email"`
    Address          string `json:"address"`
}

type PaymentRequest struct {
    ID                      string       `json:"id"`
    CollectionRequestNumber string       `json:"collectionRequestNumber"`
    ReferenceId             string       `json:"referenceId"`
    CustomerInfo            CustomerInfo `json:"customerInfo"`
    Currency                string       `json:"currency"`
    Amount                  float64      `json:"amount"`
    PaymentPurpose          string       `json:"paymentPurpose"`
    PaymentCurrency         string       `json:"paymentCurrency"`
    PaymentAmountCollected  float64      `json:"paymentAmountCollected"`
    PaymentUrl              string       `json:"paymentUrl"`
    Remarks                 string       `json:"remarks"`
    Status                  string       `json:"status"`
    IssuedAt                string       `json:"issuedAt"`
    IssuedBy                string       `json:"issuedBy"`
    ExpiredAt               string       `json:"expiredAt"`
    ExpiredGuaranteeAt      any          `json:"expiredGuaranteeAt"`
    Products                []any        `json:"products"`
    Quotes                  []any        `json:"quotes"`
    CreatedAt               string       `json:"createdAt"`
    CreatedBy               string       `json:"createdBy"`
    UpdatedAt               string       `json:"updatedAt"`
    UpdatedBy               string       `json:"updatedBy"`
    DeletedAt               string       `json:"deletedAt"`
}
