package wallex

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "vesaliusm/config"
    "vesaliusm/model/payment"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

var WallexSvc *WallexService = NewWallexService()

type WallexService struct {
    url             string
    xApiKey         string
    accessKeyId     string
    secretAccessKey string
}

func NewWallexService() *WallexService {
    return &WallexService{
        url:             config.Config("payment.wallex.url"),
        xApiKey:         config.Config("payment.wallex.x_api_key"),
        accessKeyId:     config.Config("payment.wallex.accesskeyid"),
        secretAccessKey: config.Config("payment.wallex.secretaccesskey"),
    }
}

func (s *WallexService) authenticate() (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    prm := map[string]string{
        "accessKeyId":     s.accessKeyId,
        "secretAccessKey": s.secretAccessKey,
    }
    headers := map[string]string{
        "X-Api-Key":    s.xApiKey,
        "Content-Type": "application/json",
    }
    url := fmt.Sprintf("%s/v2/authenticate", s.url)
    sleep(100)
    token := ""
    res, err := utils.GetR().
        SetContext(ctx).
        SetHeaders(headers).
        SetBody(prm).
        Post(url)
    if err != nil {
        utils.LogError(err)
        return token, err
    }

    var mx fiber.Map
    err = json.Unmarshal(res.Body(), &mx)
    if err != nil {
        return token, err
    }

    token = utils.GetString("token", mx)
    return token, nil
}

func (s* WallexService) SubmitPaymentRequest(prm utils.Map) (*payment.PaymentRequest, error) {
    var o *payment.PaymentRequest
    token, err := s.authenticate()
    if err != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to payment gateway Wallex at the moment.\nPlease try again later.")
    }

    sleep(100)
    o, err = s.submitPaymentRequest(token, prm)
    if err != nil {
        if err.Error() == "Authorization token is invalid" {
            stoken, err := s.authenticate()
            if err != nil {
                return nil, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to payment gateway Wallex at the moment.\nPlease try again later.")
            }

            sleep(100)
            o, err = s.submitPaymentRequest(stoken, prm)
            if err != nil {
                return nil, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to payment gateway Wallex at the moment.\nPlease try again later.")
            }
        }

        return nil, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to payment gateway Wallex at the moment.\nPlease try again later.")
    }
    return o, nil
}

func (s *WallexService) submitPaymentRequest(token string, prm utils.Map) (*payment.PaymentRequest, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    headers := map[string]string{
        "X-Api-Key":     s.xApiKey,
        "Authorization": token,
        "Content-Type":  "application/json",
    }
    url := fmt.Sprintf("%s/v2/collections/request/", s.url)
    sleep(200)
    var response payment.PaymentRequest
    res, err := utils.GetR().
        SetContext(ctx).
        SetHeaders(headers).
        SetBody(prm).
        SetResult(&response).
        Post(url)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    if res == nil {
        return nil, fmt.Errorf("no content exception")
    }
    return &response, nil
}

func sleep(ms int) {
    time.Sleep(time.Duration(ms) * time.Millisecond)
}
