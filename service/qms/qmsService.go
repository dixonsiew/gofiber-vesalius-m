package qms

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"vesaliusm/config"
	"vesaliusm/database"
	"vesaliusm/dto"
	"vesaliusm/model"
	"vesaliusm/model/qms"
	"vesaliusm/service/applicationUser"
	"vesaliusm/service/applicationUserFamily"
	"vesaliusm/service/applicationUserNotification"
	"vesaliusm/utils"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/nleeper/goment"
)

var QmsSvc *QmsService = NewQmsService(database.GetDb(), database.GetCtx())

type QmsService struct {
    db                                 *sqlx.DB
    ctx                                context.Context
    url                                string
    wsclientname                       string
    password                           string
    applicationUserService             *applicationUser.ApplicationUserService
    applicationUserFamilyService       *applicationUserFamily.ApplicationUserFamilyService
    applicationUserNotificationService *applicationUserNotification.ApplicationUserNotificationService
}

func NewQmsService(db *sqlx.DB, ctx context.Context) *QmsService {
    return &QmsService{
        db:                                 db,
        ctx:                                ctx,
        url:                                config.Config("queue.qms.url"),
        wsclientname:                       config.Config("queue.qms.wsclientname"),
        password:                           config.Config("queue.qms.password"),
        applicationUserService:             applicationUser.ApplicationUserSvc,
        applicationUserFamilyService:       applicationUserFamily.ApplicationUserFamilySvc,
        applicationUserNotificationService: applicationUserNotification.ApplicationUserNotificationSvc,
    }
}

func (s *QmsService) CallQMS(o []qms.QmsPatient) ([]qms.QueueResult, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()
    results := make([]qms.QueueResult, 0)
    dt, _ := goment.New()
    dtnow := dt.Format("YYYY-MM-DD")
    tnow := dt.Format("hh:mm A")

    m := map[string]string{
        "1": "Issue and Waiting",
        "2": "Calling",
        "3": "Serving Ended",
        "4": "Cancelled",
        "5": "Cancelled",
    }

    for i := range o {
        patient := o[i]
        prm := qms.RequestPayload{
            GetTicketInfoByDateListReq: qms.GetTicketInfoByDateListReq{
                WsAuth: qms.WsAuth{
                    WsClientName: s.wsclientname,
                    Password:     s.password,
                },
                DateRange: qms.DateRange{
                    StartDate: dtnow,
                    EndDate:   dtnow,
                },
                Cust: qms.Cust{
                    MrnNumber: patient.Prn,
                },
            },
        }
        var response qms.ResponseData
        res, err := utils.GetR().
            SetContext(ctx).
            SetBody(prm).
            SetResult(&response).
            Post(s.url)
        if err != nil {
            utils.LogError(err)
            return nil, handleError(err, res, "testing")
        }

        if response.GetTicketInfoByDateListResp.TotalRec == "0" {
            continue
        }

        // Helper function to process ticket info
        processTicketInfo := func(ticketInfo qms.TicketInfo) {
            fieldStatus := m[ticketInfo.TicketStatus]
            if fieldStatus == "" {
                fieldStatus = "Unknown"
            }

            patientName := "-"
            if ticketInfo.Cust.Name != "" {
                patientName = ticketInfo.Cust.Name
            }

            patientsAhead := ticketInfo.TotalWaitingAhead
            if patientsAhead == "" {
                patientsAhead = "0"
            }

            patientPrn := "-"
            if patient.Prn != "" {
                patientPrn = patient.Prn
            }

            relationship := "-"
            if patient.Relationship != "" {
                relationship = patient.Relationship
            }

            results = append(results, qms.QueueResult{
                QueueNumber:        ticketInfo.TicketStr,
                PatientsAheadOfYou: patientsAhead,
                TicketStatus:       fieldStatus,
                PatientPrn:         patientPrn,
                Relationship:       relationship,
                PatientName:        patientName,
                DoctorName:         ticketInfo.IndService.ServiceName,
                RoomNumber:         ticketInfo.IndService.ShortName,
                AsAt:               tnow,
            })
        }

        ticketInfoByDateRes := response.GetTicketInfoByDateListResp.TicketInfoByDateList.TicketInfoByDate
        // Handle different response structures
        switch v := ticketInfoByDateRes.(type) {
        case []interface{}:
            for _, item := range v {
                if ticketByDateMap, ok := item.(map[string]interface{}); ok {
                    s.processTicketInfoFromMap(ticketByDateMap, processTicketInfo)
                }
            }
        case map[string]interface{}:
            s.processTicketInfoFromMap(v, processTicketInfo)
        }
    }

    return results, nil
}

func (s *QmsService) processTicketInfoFromMap(ticketByDateMap map[string]interface{}, processor func(qms.TicketInfo)) {
    ticketInfoList, ok := ticketByDateMap["TicketInfoList"].(map[string]interface{})
    if !ok {
        return
    }

    ticketInfo, ok := ticketInfoList["TicketInfo"]
    if !ok {
        return
    }

    switch ti := ticketInfo.(type) {
    case []interface{}:
        for _, item := range ti {
            if ticketInfoMap, ok := item.(map[string]interface{}); ok {
                ticketInfoStruct := s.mapToTicketInfo(ticketInfoMap)
                processor(ticketInfoStruct)
            }
        }
    case map[string]interface{}:
        ticketInfoStruct := s.mapToTicketInfo(ti)
        processor(ticketInfoStruct)
    }
}

func (s *QmsService) mapToTicketInfo(m map[string]interface{}) qms.TicketInfo {
    ticketInfo := qms.TicketInfo{}
    ticketInfo.TicketStr = utils.GetString("TicketStr", m)
    ticketInfo.TotalWaitingAhead = utils.GetString("TotalWaitingAhead", m)
    ticketInfo.TicketStatus = utils.GetString("TicketStatus", m)

    if cust, ok := m["Cust"].(map[string]interface{}); ok {
        ticketInfo.Cust.Name = utils.GetString("Name", cust)
    }
    if indService, ok := m["IndService"].(map[string]interface{}); ok {
        ticketInfo.IndService.ServiceName = utils.GetString("ServiceName", indService)
        ticketInfo.IndService.ShortName = utils.GetString("ShortName", indService)
    }

    return ticketInfo
}

func (s *QmsService) SaveQMS(data *dto.QMSServerWebhookDto) error {
    notification := &model.OnesignalNotification{
        VisitType:         utils.NewNullString(""),
        AccountNo:         utils.NewNullString(""),
        NotificationTitle: utils.NewNullString(utils.NotificationTitleQueueNotification),
        MsgType:           utils.NewNullString(utils.MsgTypeQueueNotification),
        ShortMessage:      utils.NewNullString(""),
        FullMessage:       utils.NewNullString(""),
    }
    patientName := ""
    user, err := s.applicationUserService.FindByPRN(data.Cust.MrnNumber, s.db)
    if err != nil && err != sql.ErrNoRows {
        return err
    }

    if user != nil {
        notification.UserId = user.UserId
        lname := []string{ strings.TrimSpace(user.FirstName.String) }
        if user.MiddleName.Valid {
            lname = append(lname, user.MiddleName.String)
        }
        if user.LastName.Valid {
            lname = append(lname, user.LastName.String)
        }
        patientName = strings.Join(lname, " ")
    } else {
        family, err := s.applicationUserFamilyService.FindByPRN(data.Cust.MrnNumber, s.db)
        if err != nil {
            return err
        }

        if family != nil {
            userFamily, err := s.applicationUserService.FindByPRN(family.PatientPrn.String, s.db)
            if err != nil {
                return err
            }

            if userFamily != nil {
                notification.UserId = userFamily.UserId
                patientName = family.Fullname.String
            }
        }
    }

    switch data.Method {
    case utils.QmsMethodNew:
        notification.ShortMessage = utils.NewNullString(fmt.Sprintf("A queue number %s has been assigned for your visit. Click for more information.", data.TicketStr))
        notification.FullMessage = utils.NewNullString(fmt.Sprintf("Dear *%s*, we've generated a queue number *%s* for you after your registration. We'll be sending you timely push notifications to keep you updated on your queue status. Alternatively, you can also track your queue's progress using our Queue Tracker feature. Thank you. \n \nDoctor name : *%s* \nRoom number : *%s* \nPatients ahead of you : *%s*", patientName, data.TicketStr, data.IndService.ServiceName, data.IndService.ShortName, data.TicketInfo.TotalWaitingAhead))
    case utils.QmsMethodNear:
        notification.ShortMessage = utils.NewNullString(fmt.Sprintf("Your queue number %s will be calling soon. Click for more information.", data.TicketStr))
        notification.FullMessage = utils.NewNullString(fmt.Sprintf("Dear *%s*, your queue number *%s* is getting closer to the front. We'll be with you shortly. Your patience is appreciated. Thank you. \n \nDoctor name : *%s* \nRoom number : *%s* \nPatients ahead of you : *%s*", patientName, data.TicketStr, data.IndService.ServiceName, data.IndService.ShortName, data.TicketInfo.TotalWaitingAhead))
    case utils.QmsMethodCall:
        notification.ShortMessage = utils.NewNullString(fmt.Sprintf("Your queue number %s is now been called. Click for more information.", data.TicketStr))
        notification.FullMessage = utils.NewNullString(fmt.Sprintf("Dear *%s*, it's time to meet *%s*! Your queue number *%s*, has been called. Kindly make your way to *%s* now. Thank you.", patientName, data.IndService.ServiceName, data.TicketStr, data.IndService.ShortName))
    }

    if notification.UserId.Valid && notification.ShortMessage.Valid && notification.FullMessage.Valid {
        err := s.applicationUserNotificationService.Save(notification)
        if err != nil {
            return err
        }
    }
    
    return nil
}

func handleError(err error, res *resty.Response, ms string) error {
    var e *fiber.Error
    if errors.As(err, &e) {
        code := e.Code
        if code == fiber.StatusNoContent || code == fiber.StatusBadRequest {
            return err
        }
    }

    if res == nil {
        return fiber.NewError(fiber.StatusBadRequest, "Seems like there is an issue to handle your request. Please contact customer service for assistance.")
    }

    if res.StatusCode() == fiber.StatusBadGateway {
        return fiber.NewError(fiber.StatusBadGateway, ms)
    } else if res.StatusCode() == fiber.StatusExpectationFailed {
        return fiber.NewError(fiber.StatusExpectationFailed, ms)
    } else if res.StatusCode() == fiber.StatusNotFound {
        return fiber.NewError(fiber.StatusNotFound, ms)
    }

    return fiber.NewError(fiber.StatusNotFound, ms)
}
