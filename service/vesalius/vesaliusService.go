package vesalius

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"vesaliusm/database"
	"vesaliusm/dto"
	"vesaliusm/model"
	gm "vesaliusm/model/vesaliusGeo"
	"vesaliusm/service/applicationUser"
	"vesaliusm/service/applicationUserFamily"
	"vesaliusm/service/novaDoctor"
	"vesaliusm/service/novaDoctorPatientAppt"
	"vesaliusm/service/patientPurchaseDetails"
	"vesaliusm/service/vesaliusGeo"
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

func (s *VesaliusService) VesaliusGetCancelAppointment(prn string, data dto.PostCancelAppointmentDto) (*gm.AppointmentCancelConfirmation, error) {
	res, err := s.vesaliusGeoService.AppointmentCancelAppointment(prn, data.AppointmentNumber, "No Reason Stated")
	if err != nil {
		var e *fiber.Error
		if errors.As(err, &e) {
			if e.Code == fiber.StatusForbidden {
				ms := e.Message
				if ms == "WS-00041" || ms == "WS-00034" {
					o := gm.AppointmentCancelConfirmation{
						AppointmentStatus: "CANCELLED",
						AppointmentNumber: data.AppointmentNumber,
					}
					err := s.handleCancelAppointment(prn, data, o, err)
					return nil, err
				}
			}

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

func (s *VesaliusService) handleCancelAppointment(prn string, data dto.PostCancelAppointmentDto, o gm.AppointmentCancelConfirmation, err error) error {
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

func (s *VesaliusService) VesaliusGetChangeAppointment(prn string, data dto.PostChangeAppointmentDto) (*gm.AppointmentChangeConfirmation, error) {
	res, err := s.vesaliusGeoService.AppointmentChangeAppointment(prn, data.SlotNumber, data.AppointmentNumber, "No Reason Stated")
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

func (s *VesaliusService) handleChangeAppointment(prn string, data dto.PostChangeAppointmentDto, o gm.AppointmentChangeConfirmation, err error) error {
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

func (s *VesaliusService) VesaliusGetMakeAppointment(prn string, data dto.PostMakeAppointmentDto) (*gm.AppointmentBookingConfirmation, error) {
	res, err := s.vesaliusGeoService.AppointmentMakeAppointment(prn, data.SlotNumber, data.CaseType, data.Remark)
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

func (s *VesaliusService) handleMakeAppointment(prn string, doctorId int64, data dto.PostMakeAppointmentDto, o gm.AppointmentBookingConfirmation, err error) error {
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
		res, err := s.vesaliusGeoService.AppointmentGetFutureAppointments(f.NokPrn.String)
		if err != nil {
			var e *fiber.Error
			if errors.As(err, &e) {
				ms := e.Message
				if ms == "WS-00138" {
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

func (s *VesaliusService) VesaliusGetFutureAppointments(prn string) ([]gm.Appointment, error) {
	res, err := s.vesaliusGeoService.AppointmentGetFutureAppointments(prn)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("vesaliusGetFutureAppointments not found: %s", prn))
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

func (s *VesaliusService) vesaliusCallNextAvailSlots(data dto.PostNextAvailableSlotsDto, prn string, conn *sqlx.DB) (bool, bool, []gm.Slot, error) {
	db := conn
	if db == nil {
		db = s.db
	}
	res, err := s.vesaliusGeoService.AppointmentGetNextAvailableSlots(prn, data.SpecialtyCode, data.Mcr, data.StartDate, data.StartTime, data.CaseType)
	if err != nil {
		return false, false, nil, err
	}

	lx := res.Slots
	if lx == nil {
		return false, false, nil, fiber.NewError(fiber.StatusNoContent)
	}

	if len(lx) < 1 {
		return false, false, nil, fiber.NewError(fiber.StatusNoContent)
	}

	vesMorningDate := lx[0].Date
	vesAfternoonDate := lx[1].Date

	var morningDuplicateCount bool
	var afternoonDuplicateCount bool
	getMorningSlotAgn := false
	getAfternoonSlotAgn := false

	if vesMorningDate == data.StartDate {
		vesMorningDT := fmt.Sprintf("%s %s", lx[0].Date, lx[0].StartTime)
		morningDuplicateCount, err = s.novaDoctorPatientApptService.ExistsByPrnDateTime(prn, vesMorningDT, s.db)
		if err != nil {
			return false, false, nil, err
		}
	}

	if vesAfternoonDate == data.StartDate {
		vesAfternoonDT := fmt.Sprintf("%s %s", lx[1].Date, lx[1].StartTime)
		afternoonDuplicateCount, err = s.novaDoctorPatientApptService.ExistsByPrnDateTime(prn, vesAfternoonDT, s.db)
		if err != nil {
			return false, false, nil, err
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

	return getMorningSlotAgn, getAfternoonSlotAgn, lx, nil
}

func (s *VesaliusService) VesaliusGetNextSessionAvailableSlots(prn string, data dto.PostNextAvailableSlotsDto) ([]gm.Slot, error) {
	getMorningSlotAgn := true
	getAfternoonSlotAgn := true
	lx := make([]gm.Slot, 0)
	var doctorId int64
	var vesData []gm.Slot
	var err error

	var prm *dto.PostNextAvailableSlotsDto = &data

	for getMorningSlotAgn || getAfternoonSlotAgn {
		getMorningSlotAgn, getAfternoonSlotAgn, vesData, err = s.vesaliusCallNextAvailSlots(*prm, prn, s.db)
		if err != nil {
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

	doc_appts, _, err := s.novaDoctorPatientApptService.FindAllByDoctorId(doctorId, dateMonth, dateYear, 1)
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

func convertToTitleCase(text string) string {
	if len(text) == 0 {
		return text
	}
	return strings.ToUpper(string(text[0])) + strings.ToLower(text[1:])
}
