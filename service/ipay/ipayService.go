package ipay

import (
    "context"
    "crypto/hmac"
    "crypto/sha512"
    "encoding/hex"
    "fmt"
    "strings"
    "time"
    "vesaliusm/config"
    "vesaliusm/utils"
)

var IpaySvc *IpayService = NewIpayService()

type IpayService struct {
    merchantKey  string
    MerchantCode string
    TestEnv      string
}

func NewIpayService() *IpayService {
    return &IpayService{
        merchantKey:  config.Config("payment.ipay.merchant_key"),
        MerchantCode: config.Config("payment.ipay.merchant_code"),
        TestEnv:      config.Config("payment.ipay.testenv"),
    }
}

func (s *IpayService) Requery(refno string, amount string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    url := "https://payment.ipay88.com.my/epayment/enquiry.asp"
    prm := map[string]string{
        "MerchantCode": s.MerchantCode,
        "RefNo":        refno,
        "Amount":       amount,
    }
    data := ""
    res, err := utils.GetR().
        SetContext(ctx).
        SetQueryParams(prm).
        Get(url)
    if err != nil {
        utils.LogError(err)
        return data, err
    }

    data = res.String()
    return data, nil
}

func (s *IpayService) BuildSignature(refno string, amount string, currency string) string {
    // Remove dots and commas from amount
    amt := strings.ReplaceAll(amount, ".", "")
    amt = strings.ReplaceAll(amt, ",", "")
    
    // Build the string to sign
    v := fmt.Sprintf("%s%s%s%s%s", s.merchantKey, s.MerchantCode, refno, amount, currency)
    
    // Create HMAC-SHA512 hash
    h := hmac.New(sha512.New, []byte(s.merchantKey))
    h.Write([]byte(v))
    
    // Return hex encoded string
    return hex.EncodeToString(h.Sum(nil))
}

func (s *IpayService) BuildResponseSignature(paymentid string, refno string, amount string, currency string, status string) string {
    // Remove dots and commas from amount
    amt := strings.ReplaceAll(amount, ".", "")
    amt = strings.ReplaceAll(amt, ",", "")
    
    // Build the string to sign (order matters!)
    v := fmt.Sprintf("%s%s%s%s%s%s%s", s.merchantKey, s.MerchantCode, paymentid, refno, amt, currency, status)
    
    // Create HMAC-SHA512 hash
    h := hmac.New(sha512.New, []byte(s.merchantKey))
    h.Write([]byte(v))
    
    // Return hex encoded string
    return hex.EncodeToString(h.Sum(nil))
}
