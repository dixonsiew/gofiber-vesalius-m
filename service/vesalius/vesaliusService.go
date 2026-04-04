package vesalius

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	"vesaliusm/sql"
	"vesaliusm/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
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
                ms := err.Error()
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
            return  nil, err
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

func (s *VesaliusService) getAllPatientFutureAppointments(familyMembers []model.ApplicationUserFamily)  {
    lr := make([]gm.Appointment, 0)
    m := make(map[string]model.ApplicationUserFamily)
    for _, f := range familyMembers {
        res, err := s.vesaliusGeoService.AppointmentGetFutureAppointments(f.NokPrn.String)
        if err != nil || res == nil {
            continue
        }
        
        la := res.Appointments
        for _, a := range la {
            lr = append(lr, a)
            m[a.AppointmentNumber] = f
        }
    }
    return lr
}
