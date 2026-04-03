package vesaliusGeo

import (
    "context"
    "encoding/xml"
    "fmt"
    "slices"
    "strings"
    "time"
    "vesaliusm/config"
    "vesaliusm/database"
    model "vesaliusm/model/vesaliusGeo"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/jmoiron/sqlx"
    "github.com/nleeper/goment"
)

var VesaliusGeoSvc *VesaliusGeoService = NewVesaliusGeoService(database.GetDb(), database.GetCtx())

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
    db                           *sqlx.DB
    ctx                          context.Context
}

func NewVesaliusGeoService(db *sqlx.DB, ctx context.Context) *VesaliusGeoService {
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
        db:                           db,
        ctx:                          ctx,
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

func (s *VesaliusGeoService) PatientProcessPatientBillPayment(
    prn string, billNumber string, paymentMethod string, paymentAmount string,
    paymentRequestNumber string, remark string, payorName string) (*model.ResultBillPaymentConfirmation, error) {
    result, ex := s.patientProcessPatientBillPaymentResult(prn, billNumber, paymentMethod, paymentAmount, paymentRequestNumber, remark, payorName)
    if ex != nil {
        errorMessage := fmt.Sprintf("%s (%s)", ex.Message, ex.Code)
        requestNo := paymentRequestNumber
        err := s.updateVesaliusWSLog(errorMessage, requestNo)
        if err != nil {
            return result, err
        }
    }

    return result, nil
}

func (s *VesaliusGeoService) patientProcessPatientBillPaymentResult(
    prn string, billNumber string, paymentMethod string, paymentAmount string,
    paymentRequestNumber string, remark string, payorName string,
) (*model.ResultBillPaymentConfirmation, *model.VesaliusWSException) {
    var (
        result *model.ResultBillPaymentConfirmation = new(model.ResultBillPaymentConfirmation)
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
    url := fmt.Sprintf("%sPATIENT/ProcessPatientBillPayment.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("ProcessPatientBillPayment")
    envelope :=
        `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:ProcessPatientBillPayment>
                <axi:prn>%s</axi:prn>
                <axi:bill_no>%s</axi:bill_no>
                <axi:payment_mode>%s</axi:payment_mode>
                <axi:payment_amount>%s</axi:payment_amount>
                <axi:extermal_payment_ref>%s</axi:extermal_payment_ref>
                <axi:remark>%s</axi:remark>
                <axi:cashier_name>%s</axi:cashier_name>
                <axi:cashier_location>%s</axi:cashier_location>
                <axi:payer_name>%s</axi:payer_name>
                <axi:company_code>%s</axi:company_code>
                <axi:token_number>%s</axi:token_number>
            </axi:ProcessPatientBillPayment>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, prn, billNumber, paymentMethod, paymentAmount, paymentRequestNumber,
        remark, s.vesaliusCashierName, s.vesaliusCashierLocation, payorName, s.vesaliusServerCompanyCode, localToken)
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

    result.Set()
    return result, ex
}

func (s *VesaliusGeoService) PersonProcessPersonBiodata(
    fullname string, dob string, gender string, contactNumber string, maritalStatus string,
    address1 string, address2 string, postcode string, city string, state string,
    nationality string, email string,
) (*model.ResultPerson, error) {
    result, ex := s.personProcessPersonBiodataResult(
        fullname, dob, gender, contactNumber, maritalStatus,
        address1, address2, postcode, city, state,
        nationality, email,
    )
    if ex != nil {
        return result, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error encountered. Please contact Island Hospital. (Error Code: %s)", ex.Code))
    }

    return result, nil
}

func (s *VesaliusGeoService) personProcessPersonBiodataResult(
    fullname string, dob string, gender string, contactNumber string, maritalStatus string,
    address1 string, address2 string, postcode string, city string, state string,
    nationality string, email string,
) (*model.ResultPerson, *model.VesaliusWSException) {
    var (
        result *model.ResultPerson = new(model.ResultPerson)
        ex     *model.VesaliusWSException
    )
    var m = map[string]string{
        "Others":  "00",
        "Single":  "01",
        "Married": "02",
    }
    maritalMap := m[maritalStatus]
    g, _ := goment.New(dob, "DD/MM/YYYY")
    dobs := g.Format("DD-MMM-YYYY")
    sex := "F"
    if gender == "Male" {
        sex = "M"
    }
    localToken, resx, ex := s.Login()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer s.defLogout(localToken)

    sleep()
    url := fmt.Sprintf("%sPATIENT/ProcessPersonBiodata.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("ProcessPersonBiodata")
    envelope := `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:ProcessPersonBiodata>
                <axi:token_number>%s</axi:token_number>
                <axi:person_id></axi:person_id>
                <axi:title></axi:title>
                <axi:first_name>%s</axi:first_name>
                <axi:middle_name></axi:middle_name>
                <axi:last_name></axi:last_name>
                <axi:resident>Y</axi:resident>
                <axi:dob>%s</axi:dob>
                <axi:sex>%s</axi:sex>
                <axi:contact>%s</axi:contact>
                <axi:marital_status>%s</axi:marital_status>
                <axi:address1>%s</axi:address1>
                <axi:address2>%s</axi:address2>
                <axi:address3></axi:address3>
                <axi:postal_code>%s</axi:postal_code>
                <axi:city_state>%s %s</axi:city_state>
                <axi:nationality>%s</axi:nationality>
                <axi:operation_flag>I</axi:operation_flag>
                <axi:company_code>%s</axi:company_code>
                <axi:workstation_code></axi:workstation_code>
                <axi:country_code></axi:country_code>
                <axi:email>%s</axi:email>
            </axi:ProcessPersonBiodata>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, localToken, fullname, dobs, sex, contactNumber,
        maritalMap, address1, address2, postcode, city,
        state, nationality, s.vesaliusServerCompanyCode, email)
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

func (s *VesaliusGeoService) patientCheckPatientData(prn string) (*model.ResultListPatient, *model.VesaliusWSException) {
    result, ex := s.patientGetPatientDataResult(prn)
    if ex != nil {
        return result, ex
    }

    if len(result.Patients) > 1 {
        exx := model.VesaliusWSException{
            Code:    "99",
            Message: "More than 1 Patient Found",
        }
        return result, &exx
    }

    if len(result.Patients) < 1 {
        exx := model.VesaliusWSException{
            Code:    "99",
            Message: "No Patient Found",
        }
        return result, &exx
    }

    return result, nil
}

func (s *VesaliusGeoService) ClubsGetPatientData(docNumber string) (*model.ResultListPatient, error) {
    result, ex := s.patientGetPatientDataResult(docNumber)
    if ex != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, ex.Message)
    }

    if len(result.Patients) > 1 {
        return result, fiber.NewError(fiber.StatusBadRequest, "More than 1 patient found with the Identification Type / Number provided. Please reach out to our customer service at +604-238 3388 for further action")
    }

    return result, nil
}

func (s *VesaliusGeoService) patientGetPatientData(prn string) (*model.ResultListPatient, error) {
    result, ex := s.patientGetPatientDataResult(prn)
    if ex != nil {
        if ex.Code == "WS-00014" {
            return result, fiber.NewError(fiber.StatusBadRequest, "Incorrect PRN: The Patient PRN Number provided does not exist. Please retry")
        }

        return result, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error encountered. Please contact Island Hospital. (Error Code: %s)", ex.Code))
    }

    if len(result.Patients) > 1 {
        return result, fiber.NewError(fiber.StatusBadRequest, "More than 1 Patient Found")
    }

    if len(result.Patients) < 1 {
        return result, fiber.NewError(fiber.StatusBadRequest, "No Patient Found")
    }

    return result, nil
}

func (s *VesaliusGeoService) patientGetPatientDataResult(prn string) (*model.ResultListPatient, *model.VesaliusWSException) {
    var (
        result *model.ResultListPatient = new(model.ResultListPatient)
        ex     *model.VesaliusWSException
    )
    url := fmt.Sprintf("%sPATIENT/GetPatientData.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("getPatientData")
    envelope :=
        `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:getPatientData>
                <axi:prn>%s</axi:prn>
            </axi:getPatientData>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, prn)
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

func (s *VesaliusGeoService) PatientGetPatientOutstandingBills(prn string) (*model.ResultOutstandingBills, error) {
    result, ex := s.patientGetPatientOutstandingBillsResult(prn)
    if ex != nil {
        if ex.Code == "NH-00294" {
            return result, fiber.NewError(fiber.StatusBadRequest, ex.Message)
        }

        return result, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error encountered. Please contact Island Hospital. (Error Code: %s)", ex.Code))
    }

    return result, nil
}

func (s *VesaliusGeoService) patientGetPatientOutstandingBillsResult(prn string) (*model.ResultOutstandingBills, *model.VesaliusWSException) {
    var (
        result *model.ResultOutstandingBills = new(model.ResultOutstandingBills)
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
    url := fmt.Sprintf("%sPATIENT/GetPatientOutstandingBills.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("GetPatientOutstandingBills")
    envelope :=
        `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:GetPatientOutstandingBills>
            <axi:prn>%s</axi:prn>
            <axi:json_format>N</axi:json_format>
            <axi:company_code>%s</axi:company_code>
            <axi:token_number>%s</axi:token_number>
            </axi:GetPatientOutstandingBills>
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

func (s *VesaliusGeoService) PatientGetPatientOutstandingBillDetails(prn string, billNumber string) (*model.ResultPDFOutstandingBill, error) {
    result, ex := s.patientGetPatientOutstandingBillDetailsResult(prn, billNumber)
    if ex != nil {
        return result, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Error encountered. Please contact Island Hospital. (Error Code: %s)", ex.Code))
    }

    return result, nil
}

func (s *VesaliusGeoService) patientGetPatientOutstandingBillDetailsResult(prn string, billNumber string) (*model.ResultPDFOutstandingBill, *model.VesaliusWSException) {
    var (
        result *model.ResultPDFOutstandingBill = new(model.ResultPDFOutstandingBill)
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
    url := fmt.Sprintf("%sPATIENT/GetPatientBill.cfc", s.wsVesaliusServerBaseUrl)
    r := utils.GetRWs("GetPatientBill")
    envelope :=
        `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:GetPatientBill>
            <axi:prn>%s</axi:prn>
            <axi:bill_no>%s</axi:bill_no>
            <axi:json_format>N</axi:json_format>
            <axi:company_code>%s</axi:company_code>
            <axi:token_number>%s</axi:token_number>
            </axi:GetPatientBill>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, prn, billNumber, s.vesaliusServerCompanyCode, localToken)
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

func (s *VesaliusGeoService) getPatientFromDocuments(o model.Patient, prn string) (model.Patient, error) {
    if o.Documents == nil {
        return o, fiber.NewError(fiber.StatusNoContent)
    }

    if len(o.Documents) < 1 {
        return o, fiber.NewError(fiber.StatusNoContent)
    }

    patientDocumentCode := s.patientDocumentCode
    i := slices.IndexFunc(o.Documents, func(doc model.Document) bool {
        return strings.EqualFold(strings.TrimSpace(doc.Code), patientDocumentCode)
    })
    if i < 0 {
        return o, fiber.NewError(fiber.StatusNoContent)
    }

    patientDoc := o.Documents[i]
    if patientDoc.Value != "" {
        if strings.EqualFold(strings.TrimSpace(patientDoc.Value), prn) {
            return o, nil
        }
    }

    return o, fiber.NewError(fiber.StatusNoContent)
}

func (s *VesaliusGeoService) GetPatientDataByPRN(prn string) (*model.Patient, error) {
    result, err := s.patientGetPatientData(prn)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    if len(result.Patients) < 1 {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    return &result.Patients[0], nil
}

func (s *VesaliusGeoService) CheckPatientDataByNRIC(nric string) (*model.Patient, *model.VesaliusWSException) {
    nricWSWithDash := strings.ToLower(s.wsVesaliusNRICWithDash)
    prn := nric
    if nricWSWithDash != "y" {
        prn = strings.ReplaceAll(nric, "-", "")
    }

    result, ex := s.patientCheckPatientData(prn)
    if ex != nil {
        return nil, ex
    }

    o := result.Patients[0]
    if o.Prn == prn {
        return &o, nil
    } else {
        patient, err := s.getPatientFromDocuments(o, prn)
        if err != nil {
            return nil, &model.VesaliusWSException{
                Code:    "99",
                Message: "Patient not found",
            }
        }
        return &patient, nil
    }
}

func (s *VesaliusGeoService) GetPatientDataByNRIC(nric string) (*model.Patient, error) {
    nricWSWithDash := strings.ToLower(s.wsVesaliusNRICWithDash)
    prn := nric
    if nricWSWithDash != "y" {
        prn = strings.ReplaceAll(nric, "-", "")
    }

    result, err := s.patientGetPatientData(prn)
    if err != nil {
        return nil, err
    }
    
    o := result.Patients[0]
    if o.Prn == prn {
        return &o, nil
    } else {
        patient, err := s.getPatientFromDocuments(o, prn)
        if err != nil {
            return nil, err
        }
        return &patient, nil
    }
}

func (s *VesaliusGeoService) updateVesaliusWSLog(errorMessage string, requestNo string) error {
    query := `
        UPDATE OUTSTANDING_BILL_PAYMENT SET
          WS_CALL_FLAG = 'Y',
          WS_CALL_ERROR = :errorMessage
        WHERE PAYMENT_REQUEST_NO = :requestNo
    `
    _, err := s.db.ExecContext(s.ctx, query, errorMessage, requestNo)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func (s *VesaliusGeoService) defLogout(token string) {
    _, _ = s.Logout(token)
}

func sleep() {
    time.Sleep(1200 * time.Millisecond)
}
