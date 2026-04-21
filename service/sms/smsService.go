package sms

import (
    "fmt"
    "vesaliusm/config"
    "vesaliusm/utils"
)

var SmsSvc *SmsService

type SmsService struct {
    url         string
    apiusername string
    apipassword string
}

func NewSmsService() *SmsService {
    return &SmsService{
        url:         config.Config("auth.onewaysms.url"),
        apiusername: config.Config("auth.onewaysms.apiusername"),
        apipassword: config.Config("auth.onewaysms.apipassword"),
    }
}

func (s *SmsService) SendSignIn(mobileNo string) (string, error) {
    tac := utils.GetRandomStr(6)
    prm := map[string]string{
        "apiusername":  s.apiusername,
        "apipassword":  s.apipassword,
        "mobileno":     mobileNo,
        "senderid":     "INFO",
        "languagetype": "1",
        "message":      fmt.Sprintf("Island Hospital: Your OTP (One Time Password) is %s for Sign In", tac),
    }
    _, err := utils.GetR().SetQueryParams(prm).Get(s.url)
    if err != nil {
        return "", err
    }

    return tac, nil
}
