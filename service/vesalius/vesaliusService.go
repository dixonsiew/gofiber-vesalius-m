package vesalius

import (
    "context"
    "database/sql"
    "encoding/base64"
    "errors"
    "fmt"
    "slices"
    "strconv"
    "strings"
    "vesaliusm/database"
    "vesaliusm/dto"
    "vesaliusm/model"
    upck "vesaliusm/model/userPackage"
    gm "vesaliusm/model/vesaliusGeo"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/applicationUserFamily"
    "vesaliusm/service/novaDoctor"
    "vesaliusm/service/novaDoctorPatientAppt"
    "vesaliusm/service/patientPurchaseDetails"
    "vesaliusm/service/vesaliusGeo"
    sqx "vesaliusm/sql"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/jmoiron/sqlx"
    "github.com/nleeper/goment"
)

var VesaliusSvc *VesaliusService = NewVesaliusService(database.GetDb(), database.GetCtx())

type VesaliusService struct {
    db                            *sqlx.DB
    ctx                           context.Context
    applicationUserService        *applicationUser.ApplicationUserService
    applicationUserFamilyService  *applicationUserFamily.ApplicationUserFamilyService
    novaDoctorService             *novaDoctor.NovaDoctorService
    novaDoctorPatientApptService  *novaDoctorPatientAppt.NovaDoctorPatientApptService
    patientPurchaseDetailsService *patientPurchaseDetails.PatientPurchaseDetailsService
    vesaliusGeoService            *vesaliusGeo.VesaliusGeoService
}

func NewVesaliusService(db *sqlx.DB, ctx context.Context) *VesaliusService {
    return &VesaliusService{
        db:                            db,
        ctx:                           ctx,
        applicationUserService:        applicationUser.ApplicationUserSvc,
        applicationUserFamilyService:  applicationUserFamily.ApplicationUserFamilySvc,
        novaDoctorService:             novaDoctor.NovaDoctorSvc,
        novaDoctorPatientApptService:  novaDoctorPatientAppt.NovaDoctorPatientApptSvc,
        patientPurchaseDetailsService: patientPurchaseDetails.PatientPurchaseDetailsSvc,
        vesaliusGeoService:            vesaliusGeo.VesaliusGeoSvc,
    }
}

func (s *VesaliusService) VesaliusGetCancelAppointment(prn string, data *dto.PostCancelAppointmentDto) (*gm.AppointmentCancelConfirmation, error) {
    res, ex, err := s.vesaliusGeoService.AppointmentCancelAppointment(prn, data.AppointmentNumber, "No Reason Stated")
    if err != nil {
        if ex != nil {
            if ex.Code == "WS-00041" || ex.Code == "WS-00034" {
                o := gm.AppointmentCancelConfirmation{
                    AppointmentStatus: "CANCELLED",
                    AppointmentNumber: data.AppointmentNumber,
                }
                err := s.handleCancelAppointment(prn, data, o, err)
                return nil, err
            }
        }
        var e *fiber.Error
        if errors.As(err, &e) {
            if e.Code == fiber.StatusBadRequest {
                return nil, err
            }

            return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("vesaliusGetCancelAppointment not found: %s", prn))
        }
        return nil, err
    }

    if res == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    o := res.AppointmentCancelConfirmation
    if o.AppointmentStatus == "CANCELLED" {
        _ = s.handleCancelAppointment(prn, data, o, err)
    } else {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Failed to cancel appointment")
    }

    return &o, nil
}

func (s *VesaliusService) handleCancelAppointment(prn string, data *dto.PostCancelAppointmentDto, o gm.AppointmentCancelConfirmation, err error) error {
    tx, err := s.db.BeginTxx(s.ctx, nil)
    if err != nil {
        utils.LogError(err)
        return err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()
    err = s.novaDoctorPatientApptService.UpdateToCancel(prn, o, tx)
    if err != nil {
        return err
    }
    err = s.patientPurchaseDetailsService.UpdatePackageStatusByPurchaseNo(data.Remark, utils.PackageStatusCancelled, tx)
    if err != nil {
        return err
    }
    err = tx.Commit()
    if err != nil {
        utils.LogError(err)
        return err
    }
    return err
}

func (s *VesaliusService) VesaliusGetChangeAppointment(prn string, data *dto.PostChangeAppointmentDto) (*gm.AppointmentChangeConfirmation, error) {
    res, _, err := s.vesaliusGeoService.AppointmentChangeAppointment(prn, data.SlotNumber, data.AppointmentNumber, "No Reason Stated")
    if err != nil {
        var e *fiber.Error
        if errors.As(err, &e) {
            if e.Code == fiber.StatusBadRequest {
                return nil, err
            }

            return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("vesaliusGetChangeAppointment not found: %s", prn))
        }
        return nil, err
    }

    if res == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    o := res.AppointmentChangeConfirmation
    if o.AppointmentStatus == "CHANGED" {
        err := s.handleChangeAppointment(prn, data, o, err)
        if err != nil {
            return nil, err
        }
    } else {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Failed to reschedule appointment")
    }

    return &o, nil
}

func (s *VesaliusService) handleChangeAppointment(prn string, data *dto.PostChangeAppointmentDto, o gm.AppointmentChangeConfirmation, err error) error {
    tx, err := s.db.BeginTxx(s.ctx, nil)
    if err != nil {
        utils.LogError(err)
        return err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()
    err = s.novaDoctorPatientApptService.UpdateToReschedule(prn, data.ApptSessionType, o, tx)
    if err != nil {
        return err
    }
    err = s.patientPurchaseDetailsService.UpdatePackageStatusByPurchaseNo(data.Remark, utils.PackageStatusBooked, tx)
    if err != nil {
        return err
    }
    err = tx.Commit()
    if err != nil {
        utils.LogError(err)
        return err
    }
    return err
}

func (s *VesaliusService) VesaliusGetMakeAppointment(prn string, data *dto.PostMakeAppointmentDto) (*gm.AppointmentBookingConfirmation, error) {
    res, _, err := s.vesaliusGeoService.AppointmentMakeAppointment(prn, data.SlotNumber, data.CaseType, data.Remark)
    if err != nil {
        var e *fiber.Error
        if errors.As(err, &e) {
            if e.Code == fiber.StatusBadRequest {
                return nil, err
            }

            return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("vesaliusGetMakeAppointment not found: %s", prn))
        }
        return nil, err
    }

    if res == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    o := res.AppointmentBookingConfirmation
    if o.AppointmentStatus == "CONFIRMED" {
        doctorId, err := s.novaDoctorService.FindDoctorIdByMcr(data.DoctorMcr, s.db)
        if err != nil {
            return nil, err
        }

        err = s.handleMakeAppointment(prn, doctorId, data, o, err)
        if err != nil {
            return nil, err
        }
    }
    return &o, nil
}

func (s *VesaliusService) handleMakeAppointment(prn string, doctorId int64, data *dto.PostMakeAppointmentDto, o gm.AppointmentBookingConfirmation, err error) error {
    tx, err := s.db.BeginTxx(s.ctx, nil)
    if err != nil {
        utils.LogError(err)
        return err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()
    var rem *string
    if data.Remark != "" {
        rem = &data.Remark
    }
    err = s.novaDoctorPatientApptService.Save(prn, data.ApptSessionType, doctorId, o, *rem, tx)
    if err != nil {
        return err
    }
    err = s.patientPurchaseDetailsService.UpdatePackageStatusByPurchaseNo(data.Remark, utils.PackageStatusBooked, tx)
    if err != nil {
        return err
    }
    err = tx.Commit()
    if err != nil {
        utils.LogError(err)
        return err
    }
    return err
}

func (s *VesaliusService) getAllPatientFutureAppointments(familyMembers []model.ApplicationUserFamily) ([]gm.Appointment, map[string]model.ApplicationUserFamily) {
    lr := make([]gm.Appointment, 0)
    m := make(map[string]model.ApplicationUserFamily)
    for i := 0; i < len(familyMembers); i++ {
        f := familyMembers[i]
        res, ex, err := s.vesaliusGeoService.AppointmentGetFutureAppointments(f.NokPrn.String)
        if err != nil {
            if ex != nil {
                if ex.Code == "WS-00138" {
                    i = i - 1
                    vesaliusGeo.Sleep()
                    continue
                }
            }
            continue
        }

        if res == nil {
            continue
        }

        la := res.Appointments
        for _, a := range la {
            lr = append(lr, a)
            m[a.AppointmentNumber] = f
        }
    }
    slices.SortFunc(lr, func(a, b gm.Appointment) int {
        dt1, _ := goment.New(a.Date, "DD-MMM-YYYY")
        dt2, _ := goment.New(b.Date, "DD-MMM-YYYY")
        if dt1.IsBefore(dt2) {
            return -1
        }
        if dt1.IsAfter(dt2) {
            return 1
        }
        return 0
    })
    return lr, m
}

func (s *VesaliusService) VesaliusGetPatientFutureAppointments(prn string, isHome bool) ([]model.PatientAppointment, error) {
    lx := make([]model.PatientAppointment, 0)
    lid := make([]string, 0)
    familyMembers := make([]model.ApplicationUserFamily, 0)
    patient, err := s.applicationUserService.FindByPRN(prn, s.db)
    if err != nil {
        return nil, err
    }

    if patient == nil {
        return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("vesaliusGetFutureAppointments not found: %s", prn))
    }

    familyMembers, err = s.applicationUserFamilyService.FindAllByUserPrnAppt(patient.MasterPrn.String, true, true, s.db)
    if err != nil {
        return nil, err
    }

    lrx, mx := s.getAllPatientFutureAppointments(familyMembers)

    var (
        sessionType      = ""
        sessionStartTime = ""
        sessionEndTime   = ""
    )

    if len(lrx) > 0 && isHome {
        firstDate := lrx[0].Date
        la := make([]gm.Appointment, 0)
        for _, x := range lrx {
            if x.Date == firstDate {
                la = append(la, x)
            }
        }
        lrx = la
    }

    for i := range lrx {
        vesAppt := lrx[i]
        auf := mx[vesAppt.AppointmentNumber]
        doc, err := s.novaDoctorService.FindDoctorByMcr(vesAppt.DoctorMCR)
        if err != nil {
            return nil, err
        }

        if doc == nil {
            continue
        }

        patientAppointment := model.PatientAppointment{
            DoctorId: doc.DoctorId.Int64,
            Image:    doc.Image.String,
            MCR:      doc.MCR.String,
            Name:     doc.Name.String,
        }

        query := `
            SELECT hp.PACKAGE_IMG, hp.PACKAGE_NAME, 
              ppd.PACKAGE_PURCHASE_NO, ppd.EXPIRED_DATETIME
              FROM NOVA_DOCTOR_PATIENT_APPT ndpa
              JOIN PATIENT_PURCHASE_DETAILS ppd ON ndpa.PACKAGE_PURCHASE_NO = ppd.PACKAGE_PURCHASE_NO
              LEFT JOIN HOSPITAL_PACKAGE hp ON ppd.PACKAGE_ID = hp.PACKAGE_ID 
            WHERE ndpa.PATIENT_PRN = :prn AND ndpa.APPT_NO = :apptNo`
        list := make([]upck.UserPackage, 0)
        err = s.db.SelectContext(s.ctx, &list, query, auf.NokPrn, vesAppt.AppointmentNumber)
        if err != nil {
            utils.LogError(err)
            return nil, err
        }

        var mobileAppt *upck.UserPackage
        if len(list) > 0 {
            mobileAppt = &list[0]
            patientAppointment.Image = ""
            patientAppointment.Name = ""
        }

        ls := make([]model.NovaDoctorAppointmentLists, 0)
        err = s.db.SelectContext(s.ctx, &ls, sqx.GET_SINGLEDATE_DOCTOR_APPOINTMENTS,
            sql.Named("doctorId", doc.DoctorId.Int64),
            sql.Named("dt", vesAppt.Date),
        )
        if err != nil {
            utils.LogError(err)
            return nil, err
        }

        var appt *model.NovaDoctorAppointmentLists
        if len(ls) > 0 {
            appt = &ls[0]
        }

        if appt == nil {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }

        if appt.NormalStatus.String == "NOT AVAILABLE" {
            docApptRes, err := s.novaDoctorPatientApptService.FindApptSlotsByDoctorId(doc.DoctorId.Int64)
            if err != nil {
                utils.LogError(err)
                return nil, err
            }
            for k := range docApptRes {
                docAppt := docApptRes[k]
                if docAppt.SessionType.String == "MORNING" && appt.MorningStatus.String == "AVAILABLE" {
                    vesStartTime, _ := goment.New(vesAppt.StartTime, "hh:mm")
                    docApptStartTime, _ := goment.New(docAppt.StartTime.String, "hh:mm")
                    docApptEndTime, _ := goment.New(docAppt.EndTime.String, "hh:mm")
                    isWithinMorningRange := !vesStartTime.IsBefore(docApptStartTime) && !vesStartTime.IsAfter(docApptEndTime)

                    if isWithinMorningRange {
                        sessionType = "Morning"
                        sessionStartTime = docAppt.StartTime.String
                        sessionEndTime = docAppt.EndTime.String
                    }
                }

                if docAppt.SessionType.String == "AFTERNOON" && appt.AfternoonStatus.String == "AVAILABLE" {
                    vesStartTime, _ := goment.New(vesAppt.StartTime, "hh:mm")
                    docApptStartTime, _ := goment.New(docAppt.StartTime.String, "hh:mm")
                    docApptEndTime, _ := goment.New(docAppt.EndTime.String, "hh:mm")
                    isWithinAfternoonRange := !vesStartTime.IsBefore(docApptStartTime) && !vesStartTime.IsAfter(docApptEndTime)

                    if isWithinAfternoonRange {
                        sessionType = "Afternoon"
                        sessionStartTime = docAppt.StartTime.String
                        sessionEndTime = docAppt.EndTime.String
                    }
                }
            }
        }

        rel := "Self"
        if auf.Relationship.String != "Self" {
            rel = auf.Fullname.String
        }

        packageName := ""
        packagePurchaseNo := ""
        packageImage := ""
        expDateTime := ""
        if mobileAppt != nil {
            packageName = mobileAppt.PackageName.String
            packagePurchaseNo = mobileAppt.PackagePurchaseNo.String
            packageImage = mobileAppt.PackageImage.String
            expDateTime = mobileAppt.ExpiredDateTime.String
        }

        apptInfo := model.VesaliusApptInfo{
            ApptNo:                vesAppt.AppointmentNumber,
            ApptDate:              vesAppt.Date,
            ApptStartTime:         vesAppt.StartTime,
            ApptEndTime:           vesAppt.EndTime,
            ApptCaseType:          vesAppt.CaseType,
            ApptSessionType:       sessionType,
            ApptPatientPRN:        auf.NokPrn.String,
            ApptPatientName:       rel,
            ApptPackageName:       packageName,
            ApptPackagePurchaseNo: packagePurchaseNo,
            ApptPackageImage:      packageImage,
            PatientPackageExpiry:  expDateTime,
            SessionStartTime:      sessionStartTime,
            SessionEndTime:        sessionEndTime,
            ApptSlotType:          "Normal",
        }

        if sessionType == "Morning" || sessionType == "Afternoon" {
            apptInfo.ApptSlotType = "Session"
        }

        patientAppointment.VesaliusApptInfo = apptInfo
        lid = append(lid, strconv.FormatInt(doc.DoctorId.Int64, 10))
        lx = append(lx, patientAppointment)
        sessionType = ""
        sessionStartTime = ""
        sessionEndTime = ""
    }

    if len(lid) > 0 {
        doctorIds := strings.Join(lid, ",")
        m3, _ := s.novaDoctorService.FindAllNovaDoctorSpecialities(doctorIds, s.db)
        m4, _ := s.novaDoctorService.FindAllNovaDoctorClinicLocation(doctorIds, s.db)
        m5, _ := s.novaDoctorService.FindAllNovaDoctorContact(doctorIds, s.db)

        for i := range lx {
            o := lx[i]
            if list, ok := m3[o.DoctorId]; ok {
                lx[i].DoctorSpecialities = list
            }
            if list, ok := m4[o.DoctorId]; ok {
                lx[i].DoctorClinicLocation = list
            }
            if list, ok := m5[o.DoctorId]; ok {
                lx[i].DoctorContact = list
            }
        }
    }

    return lx, nil
}

func (s *VesaliusService) VesaliusGetFutureAppointments(prn string) ([]gm.Appointment, error) {
    res, _, err := s.vesaliusGeoService.AppointmentGetFutureAppointments(prn)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("vesaliusGetFutureAppointments not found: %s", prn))
    }

    if res == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    lx := res.Appointments
    if lx == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    if len(lx) < 1 {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    return lx, nil
}

func (s *VesaliusService) vesaliusCallNextAvailSlots(data *dto.PostNextAvailableSlotsDto, prn string, conn *sqlx.DB) (bool, bool, []gm.Slot, *gm.VesaliusWSException, error) {
    db := database.GetFromCon(conn, s.db)
    res, ex, err := s.vesaliusGeoService.AppointmentGetNextAvailableSlots(prn, data.SpecialtyCode, data.Mcr, data.StartDate, data.StartTime, data.CaseType)
    if err != nil {
        return false, false, nil, ex, err
    }

    if res == nil {
        return false, false, nil, ex, fiber.NewError(fiber.StatusNoContent)
    }

    lx := res.Slots
    if lx == nil {
        return false, false, nil, ex, fiber.NewError(fiber.StatusNoContent)
    }

    if len(lx) < 1 {
        return false, false, nil, ex, fiber.NewError(fiber.StatusNoContent)
    }

    vesMorningDate := lx[0].Date
    vesAfternoonDate := lx[1].Date

    var morningDuplicateCount bool
    var afternoonDuplicateCount bool
    getMorningSlotAgn := false
    getAfternoonSlotAgn := false

    if vesMorningDate == data.StartDate {
        vesMorningDT := fmt.Sprintf("%s %s", lx[0].Date, lx[0].StartTime)
        morningDuplicateCount, err = s.novaDoctorPatientApptService.ExistsByPrnDateTime(prn, vesMorningDT, db)
        if err != nil {
            return false, false, nil, ex, err
        }
    }

    if vesAfternoonDate == data.StartDate {
        vesAfternoonDT := fmt.Sprintf("%s %s", lx[1].Date, lx[1].StartTime)
        afternoonDuplicateCount, err = s.novaDoctorPatientApptService.ExistsByPrnDateTime(prn, vesAfternoonDT, db)
        if err != nil {
            return false, false, nil, ex, err
        }
    }

    if morningDuplicateCount && !afternoonDuplicateCount {
        if morningDuplicateCount {
            getMorningSlotAgn = true
        }

    } else if !morningDuplicateCount && afternoonDuplicateCount {
        if afternoonDuplicateCount {
            getAfternoonSlotAgn = true
        }
    } else if morningDuplicateCount && afternoonDuplicateCount {
        getMorningSlotAgn = true
        getAfternoonSlotAgn = true
    }

    return getMorningSlotAgn, getAfternoonSlotAgn, lx, nil, nil
}

func (s *VesaliusService) VesaliusGetNextSessionAvailableSlots(prn string, data *dto.PostNextAvailableSlotsDto) ([]gm.Slot, error) {
    getMorningSlotAgn := true
    getAfternoonSlotAgn := true
    lx := make([]gm.Slot, 0)
    var doctorId int64
    var vesData []gm.Slot
    var err error
    var ex *gm.VesaliusWSException

    var prm *dto.PostNextAvailableSlotsDto = data

    for getMorningSlotAgn || getAfternoonSlotAgn {
        getMorningSlotAgn, getAfternoonSlotAgn, vesData, ex, err = s.vesaliusCallNextAvailSlots(prm, prn, s.db)
        if ex != nil {
            if ex.Code != "" && ex.Message != "" {
                titleCase := utils.ToTitleCase(ex.Message)
                ms := fmt.Sprintf("%s: %s", ex.Code, titleCase)
                return nil, fiber.NewError(fiber.StatusBadRequest, ms)
            }

            return nil, fiber.NewError(fiber.StatusBadRequest, ex.ToString())
        }
        if err != nil {
            var e *fiber.Error
            if errors.As(err, &e) {
                if e.Code == fiber.StatusBadRequest {
                    return nil, err
                }

                return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("vesaliusGetNextAvailableSlots not found: %s", prn))
            }
            return nil, err
        }

        if getMorningSlotAgn {
            ts, _ := goment.New(vesData[0].StartTime, "hh:mm")
            prm.StartTime = ts.Add(1, "minute").Format("HH:mm")
        } else if getAfternoonSlotAgn {
            ts, _ := goment.New(vesData[1].StartTime, "hh:mm")
            prm.StartTime = ts.Add(1, "minute").Format("HH:mm")
        }
    }

    if prm.DoctorId != 0 {
        doctorId = int64(prm.DoctorId)
    } else {
        doctorId, err = s.novaDoctorService.FindDoctorIdByMcr(data.Mcr, s.db)
        if err != nil {
            return nil, err
        }
    }

    sd, _ := goment.New(prm.StartDate, "DD-MMM-YYYY")
    dateMonth := sd.Month()
    dateYear := sd.Year()

    doc_appts, _, err := s.novaDoctorPatientApptService.FindAllByDoctorId(doctorId, dateMonth, dateYear, "1")
    if err != nil {
        return nil, err
    }

    if getMorningSlotAgn == false && getAfternoonSlotAgn == false {
        if vesData[0].Date == prm.StartDate {
            i := slices.IndexFunc(doc_appts, func(x model.NovaDoctorAppointment) bool {
                b := x.ApptSlotType.String == "SESSION" && x.ApptDayOfWeek.String == vesData[0].Day
                vesStartTime, _ := goment.New(vesData[0].StartTime, "hh:mm")
                docApptStartTime, _ := goment.New(x.ApptStartTime.String, "hh:mm")
                docApptEndTime, _ := goment.New(x.ApptEndTime.String, "hh:mm")
                isWithinTimeRange := !vesStartTime.IsBefore(docApptStartTime) && !vesStartTime.IsAfter(docApptEndTime)
                return b && isWithinTimeRange
            })
            if i >= 0 {
                a := doc_appts[i]
                doctorMorningSession := a.ApptSessionType
                vesData[0].SessionType = doctorMorningSession.String
                lx = append(lx, vesData[0])
            }
        }

        if vesData[1].Date == prm.StartDate {
            i := slices.IndexFunc(doc_appts, func(x model.NovaDoctorAppointment) bool {
                b := x.ApptSlotType.String == "SESSION" && x.ApptDayOfWeek.String == vesData[1].Day
                vesStartTime, _ := goment.New(vesData[1].StartTime, "hh:mm")
                docApptStartTime, _ := goment.New(x.ApptStartTime.String, "hh:mm")
                docApptEndTime, _ := goment.New(x.ApptEndTime.String, "hh:mm")
                isWithinTimeRange := !vesStartTime.IsBefore(docApptStartTime) && !vesStartTime.IsAfter(docApptEndTime)
                return b && isWithinTimeRange
            })
            if i >= 0 {
                a := doc_appts[i]
                doctorMorningSession := a.ApptSessionType
                vesData[1].SessionType = doctorMorningSession.String
                lx = append(lx, vesData[1])
            }
        }
    }
    return lx, nil
}

func (s *VesaliusService) VesaliusGetNextAvailableSlots(prn string, data *dto.PostNextAvailableSlotsDto) ([]gm.Slot, error) {
    res, _, err := s.vesaliusGeoService.AppointmentGetNextAvailableSlots(prn, data.SpecialtyCode, data.Mcr, data.StartDate, data.StartTime, data.CaseType)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("vesaliusGetNextAvailableSlots not found: %s", prn))
    }

    if res == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    lx := res.Slots
    if lx == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    if len(lx) < 1 {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    return lx, nil
}

func (s *VesaliusService) VesaliusGetPatientOutstandingBillDetails(prn string, billNumber string) ([]byte, error) {
    res, _, err := s.vesaliusGeoService.PatientGetPatientOutstandingBillDetails(prn, billNumber)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusNotFound, "Incorrect PRN or Bill Number: The Patient PRN Number or Bill Number provided does not exist. Please retry")
    }

    if res == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    x := res.BillData
    if x == "" {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    decodedBytes, err := base64.StdEncoding.DecodeString(x)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to decode base64 string")
    }

    return decodedBytes, nil
}

func (s *VesaliusService) VesaliusGetPatientOutstandingBillsByPrn(prn string) (*gm.ResultOutstandingBills, error) {
    res, _, err := s.vesaliusGeoService.PatientGetPatientOutstandingBills(prn)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusNotFound, "Incorrect PRN: The Patient PRN Number provided does not exist. Please retry")
    }

    if res == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    return res, nil
}

func (s *VesaliusService) VesaliusGetPatientDataByPrn(prn string) (*gm.Patient, error) {
    res, _, err := s.vesaliusGeoService.GetPatientDataByPRN(prn)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusNotFound, "Incorrect PRN: The Patient PRN Number provided does not exist. Please retry")
    }

    if res == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    return res, nil
}

func (s *VesaliusService) VesaliusCheckpatientDataByNric(nric string) (*gm.Patient, *gm.VesaliusWSException) {
    return s.vesaliusGeoService.CheckPatientDataByNRIC(nric)
}

func (s *VesaliusService) VesaliusGetPatientDataByNric(nric string) (*gm.Patient, error) {
    res, _, err := s.vesaliusGeoService.GetPatientDataByNRIC(nric)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusNotFound, "The Identification Number provided does not exist in our hospital records. Please retry.")
    }

    if res == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    return res, nil
}

func (s *VesaliusService) vesaliusProcessPersonBiodata(biodata dto.GuestMakeNewPatientDto) (*gm.Person, error) {
    res, _, err := s.vesaliusGeoService.PersonProcessPersonBiodata(biodata)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusNotFound, "Information provided does not match with hospital patient profile. Please retry.")
    }

    if res == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    return &res.Person, nil
}
