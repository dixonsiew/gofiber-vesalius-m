package vesalius

import (
    "context"
    "database/sql"
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
        
        return nil, err
    }

    if res == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    o := res.AppointmentCancelConfirmation
    if o.AppointmentStatus == "CANCELLED" {
        tx, err := s.db.BeginTxx(s.ctx, nil)
        if err != nil {
            utils.LogError(err)
            return nil, err
        }
        defer func() {
            if err != nil {
                utils.LogError(err)
                tx.Rollback()
            }
        }()
        err = s.novaDoctorPatientApptService.UpdateToCancel(prn, o, tx)
        if err != nil {
            return nil, err
        }
        err = s.patientPurchaseDetailsService.UpdatePackageStatusByPurchaseNo(data.Remark, utils.PackageStatusCancelled, tx)
        if err != nil {
            return nil, err
        }
        err = tx.Commit()
        if err != nil {
            return nil, err
        }
    } else {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Failed to cancel appointment")
    }

    return &o, nil
}
