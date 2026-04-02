package vesaliusGeo

import (
    "encoding/xml"
    "fmt"
    "strings"
    "time"
    "vesaliusm/config"
    model "vesaliusm/model/vesaliusGeo"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
)

var VesaliusGeoSvc *VesaliusGeoService = NewVesaliusGeoService()

type VesaliusGeoService struct {
    patientDocumentCode          string
    wsVesaliusMakeApptReasonCode string
    wsVesaliusNRICWithDash       string
    wsVesaliusServerBaseUrl      string
    vesaliusServerCompanyCode    string
    vesaliusServerSystemCode     string
    vesaliusServerPassword       string
    vesaliusChargeCategory       string
    vesaliusDocumentType         string
    vesaliusPaymentClass         string
    vesaliusCashierName          string
    vesaliusCashierLocation      string
}

func NewVesaliusGeoService() *VesaliusGeoService {
    return &VesaliusGeoService{
        patientDocumentCode:          config.Config("patient.document.code"),
        wsVesaliusMakeApptReasonCode: config.Config("ws.vesalius.makeappt.reasoncode"),
        wsVesaliusNRICWithDash:       config.Config("ws.vesalius.nric.withDash"),
        wsVesaliusServerBaseUrl:      config.Config("ws.vesalius.server_baseurl"),
        vesaliusServerCompanyCode:    config.Config("ws.vesalius.server_companycode"),
        vesaliusServerSystemCode:     config.Config("ws.vesalius.server_systemcode"),
        vesaliusServerPassword:       config.Config("ws.vesalius.server_password"),
        vesaliusChargeCategory:       config.Config("ws.vesalius.charge_category"),
        vesaliusDocumentType:         config.Config("ws.vesalius.document_type"),
        vesaliusPaymentClass:         config.Config("ws.vesalius.payment_class"),
        vesaliusCashierName:          config.Config("ws.vesalius.cashier_name"),
        vesaliusCashierLocation:      config.Config("ws.vesalius.cashier_location"),
    }
}

func (s *VesaliusGeoService) Login() (string, *model.ResultToken, *model.VesaliusWSException) {
    var (
        token  string             = ""
        result *model.ResultToken = new(model.ResultToken)
        ex     *model.VesaliusWSException
    )
    url := fmt.Sprintf("%sAUTHENTICATION/Login.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("login")
    envelope :=
        `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:login>
                <axi:company_code>%s</axi:company_code>
                <axi:system_code>%s</axi:system_code>
                <axi:password>%s</axi:password>
            </axi:login>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, s.vesaliusServerCompanyCode, s.vesaliusServerSystemCode, s.vesaliusServerPassword)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return token, result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return token, result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := model.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return token, result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    if result.Error.ErrorCode != "" {
        return token, result, &model.VesaliusWSException{
            Code:    result.Error.ErrorCode,
            Message: result.Error.ErrorMessage,
        }
    }

    token = result.Token.TokenNumber
    return strings.TrimSpace(token), result, ex
}

func (s *VesaliusGeoService) Logout(token string) (*model.ResultLogout, *model.VesaliusWSException) {
    var (
        result *model.ResultLogout = new(model.ResultLogout)
        ex     *model.VesaliusWSException
    )
    url := fmt.Sprintf("%sAUTHENTICATION/Logout.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("Logout")
    envelope :=
        `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:Logout>
                <axi:token_number>%s</axi:token_number>
            </axi:Logout>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, token)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := model.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    if result.Error.ErrorCode != "" {
        return result, &model.VesaliusWSException{
            Code:    result.Error.ErrorCode,
            Message: result.Error.ErrorMessage,
        }
    }

    return result, ex
}

func (s *VesaliusGeoService) AppointmentCancelAppointment(prn string, appointmentNumber string, reason string) (*model.ResultAppointmentCancelConfirmation, error) {
    result, ex := s.appointmentCancelAppointmentResult(prn, appointmentNumber, reason)
    if ex != nil {
        if ex.Code == "WS-00041" || ex.Code == "WS-00034" {
            return result, fiber.NewError(fiber.StatusForbidden, ex.Code)
        }

        if ex.Code == "WS-00009" || ex.Code == "WS-00005" {
            res, e := s.appointmentCancelAppointmentResult(prn, appointmentNumber, reason)
            if e != nil {
                return res, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to Vesalius at the moment.\nPlease try again later.")
            }
        }

        return result, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to Vesalius at the moment.\nPlease try again later.")
    }

    return result, nil
}

func (s *VesaliusGeoService) appointmentCancelAppointmentResult(prn string, appointmentNumber string, reason string) (*model.ResultAppointmentCancelConfirmation, *model.VesaliusWSException) {
    var (
        result *model.ResultAppointmentCancelConfirmation = new(model.ResultAppointmentCancelConfirmation)
        ex     *model.VesaliusWSException
    )
    localToken, resx, ex := s.Login()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer s.defLogout(localToken)

    sleep()
    url := fmt.Sprintf("%sAPPOINTMENT/cancelAppointment.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("CancelAppointment")
    envelope :=
        `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:CancelAppointment>
                <axi:prn>%s</axi:prn>
                <axi:appointment_number>%s</axi:appointment_number>
                <axi:company_code>%s</axi:company_code>
                <axi:token_number>%s</axi:token_number>
                <axi:reason>%s</axi:reason>
            </axi:CancelAppointment>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, prn, appointmentNumber, s.vesaliusServerCompanyCode, localToken, reason)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := model.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    if result.Error.ErrorCode != "" {
        return result, &model.VesaliusWSException{
            Code:    result.Error.ErrorCode,
            Message: result.Error.ErrorMessage,
        }
    }

    return result, ex
}

func (s *VesaliusGeoService) AppointmentChangeAppointment(prn string, slotNumber string, appointmentNumber string, reason string) (*model.ResultAppointmentChangeConfirmation, error) {
    result, ex := s.appointmentChangeAppointmentResult(prn, slotNumber, appointmentNumber, reason)
    if ex != nil {
        if ex.Code == "WS-00009" || ex.Code == "WS-00005" {
            res, e := s.appointmentChangeAppointmentResult(prn, slotNumber, appointmentNumber, reason)
            if e != nil {
                return res, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to Vesalius at the moment.\nPlease try again later.")
            }
        }

        return result, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to Vesalius at the moment.\nPlease try again later.")
    }

    return result, nil
}

func (s *VesaliusGeoService) appointmentChangeAppointmentResult(prn string, slotNumber string, appointmentNumber string, reason string) (*model.ResultAppointmentChangeConfirmation, *model.VesaliusWSException) {
    var (
        result *model.ResultAppointmentChangeConfirmation = new(model.ResultAppointmentChangeConfirmation)
        ex     *model.VesaliusWSException
    )
    localToken, resx, ex := s.Login()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer s.defLogout(localToken)

    sleep()
    url := fmt.Sprintf("%sAPPOINTMENT/changeAppointment.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("ChangeAppointment")
    envelope :=
        `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:ChangeAppointment>
                <axi:prn>%s</axi:prn>
                <axi:appointment_number>%s</axi:appointment_number>
                <axi:company_code>%s</axi:company_code>
                <axi:slot_number>%s</axi:slot_number>
                <axi:token_number>%s</axi:token_number>
                <axi:reason>%s</axi:reason>
                <axi:mobile_app>Y</axi:mobile_app>
            </axi:ChangeAppointment>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, prn, appointmentNumber, s.vesaliusServerCompanyCode, slotNumber, localToken, reason)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := model.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    if result.Error.ErrorCode != "" {
        return result, &model.VesaliusWSException{
            Code:    result.Error.ErrorCode,
            Message: result.Error.ErrorMessage,
        }
    }

    return result, ex
}

func (s *VesaliusGeoService) AppointmentMakeAppointment(prn string, slotNumber string, caseType string, remark string) (*model.ResultAppointmentBookingConfirmation, error) {
    result, ex := s.appointmentMakeAppointmentResult(prn, slotNumber, caseType, remark)
    if ex != nil {
        if ex.Code == "WS-00009" || ex.Code == "WS-00005" {
            res, e := s.appointmentMakeAppointmentResult(prn, slotNumber, caseType, remark)
            if e != nil {
                return res, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to Vesalius at the moment.\nPlease try again later.")
            }
        }

        return result, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to Vesalius at the moment.\nPlease try again later.")
    }

    return result, nil
}

func (s *VesaliusGeoService) appointmentMakeAppointmentResult(prn string, slotNumber string, caseType string, remark string) (*model.ResultAppointmentBookingConfirmation, *model.VesaliusWSException) {
    var (
        result *model.ResultAppointmentBookingConfirmation = new(model.ResultAppointmentBookingConfirmation)
        ex     *model.VesaliusWSException
    )
    localToken, resx, ex := s.Login()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer s.defLogout(localToken)

    sleep()
    sremark := remark
    reasonCode := s.wsVesaliusMakeApptReasonCode
    url := fmt.Sprintf("%sAPPOINTMENT/makeAppointment.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("MakeAppointment")
    envelope :=
        `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:MakeAppointment>
                <axi:prn>%s</axi:prn>
                <axi:slot_number>%s</axi:slot_number>
                <axi:company_code>%s</axi:company_code>
                <axi:case_type>%s</axi:case_type>
                <axi:token_number>%s</axi:token_number>
                <axi:reason_code>%s</axi:reason_code>
                <axi:remark>%s</axi:remark>
                <axi:mobile_app>Y</axi:mobile_app>
            </axi:MakeAppointment>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, prn, slotNumber, s.vesaliusServerCompanyCode, caseType, localToken, reasonCode, sremark)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := model.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    if result.Error.ErrorCode != "" {
        return result, &model.VesaliusWSException{
            Code:    result.Error.ErrorCode,
            Message: result.Error.ErrorMessage,
        }
    }

    return result, ex
}

func (s *VesaliusGeoService) AppointmentGetNextAvailableSlots(
    prn string, specialtyCode string, mcr string, 
    startDate string, startTime string, caseType string) (*model.ResultListSlot, error) {
    result, ex := s.appointmentGetNextAvailableSlotsResult(prn, specialtyCode, mcr, startDate, startTime, caseType)
    if ex != nil {
        if ex.Code == "WS-00009" || ex.Code == "WS-00005" {
            res, e := s.appointmentGetNextAvailableSlotsResult(prn, specialtyCode, mcr, startDate, startTime, caseType)
            if e != nil {
                return res, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to Vesalius at the moment.\nPlease try again later.")
            }
        }

        return result, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to Vesalius at the moment.\nPlease try again later.")
    }

    return result, nil
}

func (s *VesaliusGeoService) appointmentGetNextAvailableSlotsResult(
    prn string, specialtyCode string, mcr string, 
    startDate string, startTime string, caseType string) (*model.ResultListSlot, *model.VesaliusWSException) {
    var (
        result *model.ResultListSlot = new(model.ResultListSlot)
        ex     *model.VesaliusWSException
    )
    localToken, resx, ex := s.Login()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer s.defLogout(localToken)

    sleep()
    url := fmt.Sprintf("%sAPPOINTMENT/getNextAvailableSlots.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("getNextAvailableSlots")
    envelope :=
        `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:getNextAvailableSlots>
                <axi:token_number>%s</axi:token_number>
                <axi:prn>%s</axi:prn>
                <axi:personid>%s</axi:personid>
                <axi:company_code>%s</axi:company_code>
                <axi:specialty_code>%s</axi:specialty_code>
                <axi:start_date>%s</axi:start_date>
                <axi:start_time>%s</axi:start_time>
                <axi:case_type>%s</axi:case_type>
                <axi:mcr>%s</axi:mcr>
            </axi:getNextAvailableSlots>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, localToken, prn, "", s.vesaliusServerCompanyCode, specialtyCode, startDate, startTime, caseType, mcr)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := model.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    if result.Error.ErrorCode != "" {
        return result, &model.VesaliusWSException{
            Code:    result.Error.ErrorCode,
            Message: result.Error.ErrorMessage,
        }
    }

    return result, ex
}

func (s *VesaliusGeoService) AppointmentGetFutureAppointments(prn string) (*model.ResultListAppointment, error) {
    result, ex := s.appointmentGetFutureAppointmentsResult(prn)
    if ex != nil {
        if ex.Code == "WS-00009" || ex.Code == "WS-00005" {
            res, e := s.appointmentGetFutureAppointmentsResult(prn)
            if e != nil {
                return res, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to Vesalius at the moment.\nPlease try again later.")
            }
        }

        return result, fiber.NewError(fiber.StatusBadRequest, "Encountered connection issue to Vesalius at the moment.\nPlease try again later.")
    }

    return result, nil
}

func (s *VesaliusGeoService) appointmentGetFutureAppointmentsResult(prn string) (*model.ResultListAppointment, *model.VesaliusWSException) {
    var (
        result *model.ResultListAppointment = new(model.ResultListAppointment)
        ex     *model.VesaliusWSException
    )
    localToken, resx, ex := s.Login()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer s.defLogout(localToken)

    sleep()
    url := fmt.Sprintf("%sAPPOINTMENT/getFutureAppointment.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("getFutureAppointment")
    envelope :=
        `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:getFutureAppointment>
                <axi:prn>%s</axi:prn>
                <axi:company_code>%s</axi:company_code>
                <axi:token_number>%s</axi:token_number>
            </axi:getFutureAppointment>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, prn, s.vesaliusServerCompanyCode, localToken)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := model.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &model.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    if result.Error.ErrorCode != "" {
        return result, &model.VesaliusWSException{
            Code:    result.Error.ErrorCode,
            Message: result.Error.ErrorMessage,
        }
    }

    return result, ex
}

func (s *VesaliusGeoService) defLogout(token string) {
    _, _ = s.Logout(token)
}

func sleep() {
    time.Sleep(1200 * time.Millisecond)
}
