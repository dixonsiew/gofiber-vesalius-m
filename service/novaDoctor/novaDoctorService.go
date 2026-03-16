package novaDoctor

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"vesaliusm/config"
	"vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

var NovaDoctorSvc *NovaDoctorService = NewNovaDoctorService(database.GetDb(), database.GetCtx())

type NovaDoctorService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaDoctorService(db *sqlx.DB, ctx context.Context) *NovaDoctorService {
    return &NovaDoctorService{db: db, ctx: ctx}
}

func (s *NovaDoctorService) Save(doctor *model.NovaDoctor) error {
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

    // Get next sequence value
    var doctorID int64
    err = s.db.GetContext(s.ctx, &doctorID, "SELECT DOCTOR_ID_SEQ.NEXTVAL FROM DUAL")
    if err != nil {
        utils.LogError(err)
        return err
    }
    doctor.DoctorID.Int64 = doctorID

    // Insert main doctor
    _, err = tx.ExecContext(s.ctx, `
        INSERT INTO NOVA_DOCTOR (
            DOCTOR_ID, GENDER, IMAGE, RESIZE_IMAGE, MCR, NAME,
            NATIONALITY, DISPLAY_SEQUENCE, ALLOW_APPOINTMENT,
            CONSULTANT_TYPE, IS_FOR_PACKAGE, QUALIFICATIONS_SHORT,
            REGISTRATION_NO
        ) VALUES (
            :doctorId, :gender, :image, :resizeImage, :mcr, :name, :nationality, :disSeq, :allowAppt, :cType, :isForPackage, :qualifications, :reg
        )`,
        sql.Named("doctorId", doctor.DoctorID.Int64),
        sql.Named("gender", doctor.Gender.String),
        sql.Named("image", doctor.Image.String),
        sql.Named("resizeImage", doctor.ResizeImage),
        sql.Named("mcr", doctor.MCR.String),
        sql.Named("name", doctor.Name.String),
        sql.Named("nationality", doctor.Nationality.String),
        sql.Named("disSeq", doctor.DisplaySequence.Int32),
        sql.Named("allowAppt", doctor.AllowAppointment.String),
        sql.Named("cType", doctor.ConsultantType.String),
        sql.Named("isForPackage", doctor.IsForPackage.String),
        sql.Named("qualifications", doctor.Qualifications.String),
        sql.Named("reg", doctor.RegistrationNum.String),
    )
    if err != nil {
        return err
    }

    // Save child records
    if err := s.saveDoctorClinicHoursTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorClinicLocationTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorContactsTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorQualificationsTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorSpecialitiesTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorSpecialtyTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorSpokenLanguageTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorAppointmentTx(tx, doctor); err != nil {
        return err
    }

    return tx.Commit()
}

func (s *NovaDoctorService) ResizeAllDoctorImage(image string, doctorID int) error {
    _, err := s.db.ExecContext(s.ctx,
        `UPDATE NOVA_DOCTOR SET RESIZE_IMAGE = :1 WHERE DOCTOR_ID = :2`,
        image, doctorID,
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *NovaDoctorService) Update(doctor *model.NovaDoctor) error {
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

    // Update main doctor
    _, err = tx.ExecContext(s.ctx, `
        UPDATE NOVA_DOCTOR SET
            GENDER = :gender,
            IMAGE = :image,
            RESIZE_IMAGE = :resizeImage,
            MCR = :mcr,
            NAME = :name,
            NATIONALITY = :nationality,
            DISPLAY_SEQUENCE = :disSeq,
            ALLOW_APPOINTMENT = :allowAppt,
            CONSULTANT_TYPE = :cType,
            IS_FOR_PACKAGE = :isForPackage,
            QUALIFICATIONS_SHORT = :qualifications,
            REGISTRATION_NO = :reg
        WHERE DOCTOR_ID = :doctorId
        `,
        sql.Named("gender", doctor.Gender.String),
        sql.Named("image", doctor.Image.String),
        sql.Named("resizeImage", doctor.ResizeImage),
        sql.Named("mcr", doctor.MCR.String),
        sql.Named("name", doctor.Name.String),
        sql.Named("nationality", doctor.Nationality.String),
        sql.Named("disSeq", doctor.DisplaySequence.Int32),
        sql.Named("allowAppt", doctor.AllowAppointment.String),
        sql.Named("cType", doctor.ConsultantType.String),
        sql.Named("isForPackage", doctor.IsForPackage.String),
        sql.Named("qualifications", doctor.Qualifications.String),
        sql.Named("reg", doctor.RegistrationNum.String),
        sql.Named("doctorId", doctor.DoctorID.Int64),
    )
    if err != nil {
        return err
    }

    // Delete existing child records
    if err := s.deleteChildRecordsTx(tx, doctor.DoctorID.Int64); err != nil {
        return err
    }

    // Re-insert child records
    if err := s.saveDoctorClinicHoursTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorClinicLocationTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorContactsTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorQualificationsTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorSpecialitiesTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorSpecialtyTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorSpokenLanguageTx(tx, doctor); err != nil {
        return err
    }
    if err := s.saveDoctorAppointmentTx(tx, doctor); err != nil {
        return err
    }

    return tx.Commit()
}

func (s *NovaDoctorService) deleteChildRecordsTx(tx *sqlx.Tx, doctorId int64) error {
    tables := []string{
        "NOVA_DOCTOR_CLINIC_HOURS",
        "NOVA_DOCTOR_CLINIC_LOCATION",
        "NOVA_DOCTOR_CONTACT",
        "NOVA_DOCTOR_QUALIFICATIONS",
        "NOVA_DOCTOR_SPECIALITIES",
        "NOVA_DOCTOR_SPECIALTY",
        "NOVA_DOCTOR_SPOKEN_LANGUAGE",
        "NOVA_DOCTOR_APPT_SLOT",
    }
    for _, table := range tables {
        q := fmt.Sprintf("DELETE FROM %s WHERE DOCTOR_ID = :doctorId", table)
        if _, err := tx.ExecContext(s.ctx, q, doctorId); err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func (s *NovaDoctorService) DeleteByDoctorId(doctorId int64) error {
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

    if err := s.deleteChildRecordsTx(tx, doctorId); err != nil {
        return err
    }

    _, err = tx.ExecContext(s.ctx, `DELETE FROM NOVA_DOCTOR WHERE DOCTOR_ID = :doctorId`, doctorId)
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (s *NovaDoctorService) DeleteImageById(doctorId int64) error {
    _, err := s.db.ExecContext(s.ctx, `UPDATE NOVA_DOCTOR SET IMAGE = NULL WHERE DOCTOR_ID = :doctorId`, doctorId)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *NovaDoctorService) List(page string, limit string, isWebadmin bool) (*model.PagedList, error) {
    total, err := s.Count(isWebadmin)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAll(pager.GetLowerBound(), pager.PageSize, isWebadmin)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *NovaDoctorService) FindAll(offset int, limit int, isWebadmin bool) ([]model.NovaDoctor, error) {
    var query string
    if isWebadmin {
        query = `
            SELECT ` + getNovaDoctorCols() + ` FROM NOVA_DOCTOR
            ORDER BY DISPLAY_SEQUENCE, UPPER(NAME)
            OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
        `
    } else {
        query = `
            SELECT DOCTOR_ID, GENDER, RESIZE_IMAGE AS IMAGE, MCR, NAME,
                   NATIONALITY, DISPLAY_SEQUENCE, ALLOW_APPOINTMENT,
                   QUALIFICATIONS_SHORT, REGISTRATION_NO, CONSULTANT_TYPE,
                   IS_FOR_PACKAGE
            FROM NOVA_DOCTOR
            WHERE IS_FOR_PACKAGE = 'N'
            ORDER BY DISPLAY_SEQUENCE, UPPER(NAME)
            OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
        `
    }

    doctors := make([]model.NovaDoctor, 0)
    err := s.db.SelectContext(s.ctx, &doctors, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return doctors, err
    }

    if len(doctors) == 0 {
        return doctors, nil
    }

    // Collect IDs for child fetches
    ids := make([]int64, 0)
    for i, d := range doctors {
        ids = append(ids, d.DoctorID.Int64)
        // Set showMakeAppointmentButton based on config
        if config.GetIpayTestEnv() == "Y" && d.AllowAppointment.String == "Y" {
            doctors[i].ShowMakeAppointmentButton = "Y"
        } else {
            doctors[i].ShowMakeAppointmentButton = "N"
        }
    }

    // Fetch child data in bulk
    specialtyMap, err := s.findAllNovaSpecialtyMap(s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    specMap, err := s.findAllNovaDoctorSpecialities(s.db, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    locMap, err := s.findAllNovaDoctorClinicLocation(s.db, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    contactMap, err := s.findAllNovaDoctorContact(s.db, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    doctorSpecialtyMap, err := s.findAllNovaDoctorSpecialty(s.db, ids, specialtyMap)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    // Assign child collections
    for i, doc := range doctors {
        if list, ok := specMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorSpecialities = list
        }
        if list, ok := locMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorClinicLocation = list
        }
        if list, ok := contactMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorContact = list
        }
        if list, ok := doctorSpecialtyMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorSpecialty = list
            // Check if any primary specialty exists; if not, disable appointment
            hasPrimary := false
            for _, sp := range list {
                if sp.PrimarySpecialtyV.Int32 == 1 {
                    hasPrimary = true
                    break
                }
            }
            if !hasPrimary {
                doctors[i].AllowAppointment.String = "N"
            }
        } else {
            doctors[i].AllowAppointment.String = "N"
        }
    }

    return doctors, nil
}

func (s *NovaDoctorService) FindAllHSMcrAndName() ([]model.NovaDoctor, error) {
    query := `
        SELECT DOCTOR_ID, MCR, NAME
        FROM NOVA_DOCTOR
        WHERE IS_FOR_PACKAGE = 'Y'
        ORDER BY DISPLAY_SEQUENCE, UPPER(NAME)
    `
    list := make([]model.NovaDoctor, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return list, err
    }
    return list, nil
}

func (s *NovaDoctorService) Count(isWebadmin bool) (int, error) {
    var query string
    if isWebadmin {
        query = `SELECT COUNT(DOCTOR_ID) FROM NOVA_DOCTOR`
    } else {
        query = `SELECT COUNT(DOCTOR_ID) FROM NOVA_DOCTOR WHERE IS_FOR_PACKAGE = 'N'`
    }
    var count int
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *NovaDoctorService) FindByKeyword(keyword string, offset, limit int) ([]model.NovaDoctor, error) {
    query := `
        SELECT ` + getNovaDoctorCols() + ` FROM NOVA_DOCTOR nd
        WHERE nd.IS_FOR_PACKAGE = 'N'
          AND (
            nd.DOCTOR_ID IN (
                SELECT nds.DOCTOR_ID
                FROM NOVA_DOCTOR_SPECIALITIES nds
                WHERE LOWER(nds.SPECIALITIES) LIKE :key OR LOWER(nds.SUBSPECIALTY) LIKE :key
            )
            OR LOWER(nd.NAME) LIKE :key
          )
        ORDER BY nd.DISPLAY_SEQUENCE, nd.NAME
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    key := strings.ToLower(keyword)
    doctors := make([]model.NovaDoctor, 0)
    err := s.db.SelectContext(s.ctx, &doctors, query, sql.Named("key", key), sql.Named("offset", offset), sql.Named("limit", limit))
    if err != nil {
        utils.LogError(err)
        return doctors, err
    }

    if len(doctors) == 0 {
        return doctors, nil
    }

    ids := make([]int64, len(doctors))
    for _, d := range doctors {
        ids = append(ids, d.DoctorID.Int64)
    }

    specialtyMap, err := s.findAllNovaSpecialtyMap(s.db)
    if err != nil {
        return nil, err
    }
    specMap, err := s.findAllNovaDoctorSpecialities(s.db, ids)
    if err != nil {
        return nil, err
    }
    locMap, err := s.findAllNovaDoctorClinicLocation(s.db, ids)
    if err != nil {
        return nil, err
    }
    contactMap, err := s.findAllNovaDoctorContact(s.db, ids)
    if err != nil {
        return nil, err
    }
    doctorSpecialtyMap, err := s.findAllNovaDoctorSpecialty(s.db, ids, specialtyMap)
    if err != nil {
        return nil, err
    }

    for i, doc := range doctors {
        if list, ok := specMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorSpecialities = list
        }
        if list, ok := locMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorClinicLocation = list
        }
        if list, ok := contactMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorContact = list
        }
        if list, ok := doctorSpecialtyMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorSpecialty = list
        }
    }

    return doctors, nil
}

func (s *NovaDoctorService) GuestFindByKeyword(keyword string, offset, limit int) ([]model.NovaDoctor, error) {
    doctors, err := s.FindByKeyword(keyword, offset, limit)
    if err != nil {
        return nil, err
    }
    for i := range doctors {
        if config.GetIpayTestEnv() == "Y" && doctors[i].AllowAppointment.String == "Y" {
            doctors[i].ShowMakeAppointmentButton = "Y"
        } else {
            doctors[i].ShowMakeAppointmentButton = "N"
        }
    }
    return doctors, nil
}

func (s *NovaDoctorService) ListByKeyword(keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountByKeyword(keyword)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindByKeyword(keyword, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }
    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *NovaDoctorService) ListByKeywordGuest(keyword string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountByKeyword(keyword)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.GuestFindByKeyword(keyword, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }
    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *NovaDoctorService) CountByKeyword(keyword string) (int, error) {
    query := `
        SELECT COUNT(nd.DOCTOR_ID)
        FROM NOVA_DOCTOR nd
        WHERE nd.IS_FOR_PACKAGE = 'N'
          AND (
            nd.DOCTOR_ID IN (
                SELECT nds.DOCTOR_ID
                FROM NOVA_DOCTOR_SPECIALITIES nds
                WHERE LOWER(nds.SPECIALITIES) LIKE :key OR LOWER(nds.SUBSPECIALTY) LIKE :key
            )
            OR LOWER(nd.NAME) LIKE :key
          )
    `
    key := strings.ToLower(keyword)
    var count int
    err := s.db.GetContext(s.ctx, &count, query, sql.Named("key", key))
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *NovaDoctorService) ExistsByOtherMcr(mcr string, doctorId int) (bool, error) {
    query := `SELECT COUNT(DOCTOR_ID) FROM NOVA_DOCTOR WHERE MCR = :mcr AND DOCTOR_ID <> :doctorId`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, mcr, doctorId)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, nil
}

func (s *NovaDoctorService) ExistsByMcr(mcr string) (bool, error) {
    query := `SELECT COUNT(DOCTOR_ID) FROM NOVA_DOCTOR WHERE MCR = :mcr`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, mcr)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count > 0, nil
}

func (s *NovaDoctorService) FindAllByDoctorId(doctorId int64) (*model.NovaDoctor, error) {
    var doctor model.NovaDoctor
    err := s.db.GetContext(s.ctx, &doctor, `SELECT ` + getNovaDoctorCols() + ` FROM NOVA_DOCTOR WHERE DOCTOR_ID = :doctorId`, doctorId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }

    ids := []int64{doctorId}

    specialtyMap, err := s.findAllNovaSpecialtyMap(s.db)
    if err != nil {
        return nil, err
    }

    // Fetch all child data
    spokenMap, _ := s.findAllNovaDoctorSpokenLanguage(s.db, ids)
    qualMap, _ := s.findAllNovaDoctorQualifications(s.db, ids)
    specMap, _ := s.findAllNovaDoctorSpecialities(s.db, ids)
    locMap, _ := s.findAllNovaDoctorClinicLocation(s.db, ids)
    hoursMap, _ := s.findAllNovaDoctorClinicHours(s.db, ids)
    contactMap, _ := s.findAllNovaDoctorContact(s.db, ids)
    doctorSpecialtyMap, _ := s.findAllNovaDoctorSpecialty(s.db, ids, specialtyMap)
    apptMap, _ := s.findAllNovaDoctorAppointments(s.db, ids)

    if list, ok := spokenMap[doctorId]; ok {
        doctor.DoctorSpokenLanguage = list
    }
    if list, ok := qualMap[doctorId]; ok {
        doctor.DoctorQualifications = list
    }
    if list, ok := specMap[doctorId]; ok {
        doctor.DoctorSpecialities = list
    }
    if list, ok := locMap[doctorId]; ok {
        doctor.DoctorClinicLocation = list
    }
    if list, ok := hoursMap[doctorId]; ok {
        doctor.DoctorClinicHours = list
    }
    if list, ok := contactMap[doctorId]; ok {
        doctor.DoctorContact = list
    }
    if list, ok := doctorSpecialtyMap[doctorId]; ok {
        doctor.DoctorSpecialty = list
    }
    if list, ok := apptMap[doctorId]; ok {
        doctor.DoctorAppointment = list
    }

    return &doctor, nil
}

func (s *NovaDoctorService) FindDoctorNameByDoctorId(doctorId int64) (string, error) {
    var name string
    err := s.db.GetContext(s.ctx, &name, `SELECT NAME FROM NOVA_DOCTOR WHERE DOCTOR_ID = :doctorId`, doctorId)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", nil
        }
        utils.LogError(err)
        return "", err
    }
    return name, nil
}

func (s *NovaDoctorService) FindDoctorIdByMcr(mcr string) (int64, error) {
    var id int64
    err := s.db.GetContext(s.ctx, &id, `SELECT DOCTOR_ID FROM NOVA_DOCTOR WHERE MCR = :mcr`, mcr)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, err
        }
        utils.LogError(err)
        return 0, err
    }
    return id, nil
}

func (s *NovaDoctorService) FindDoctorByMcr(mcr string) (*model.NovaDoctor, error) {
    var doctor model.NovaDoctor
    err := s.db.GetContext(s.ctx, &doctor, `SELECT ` + getNovaDoctorCols() + ` FROM NOVA_DOCTOR WHERE MCR = :mcr`, mcr)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &doctor, nil
}

func (s *NovaDoctorService) FindAllByMcr(mcr string) ([]model.NovaDoctor, error) {
    doctors := make([]model.NovaDoctor, 0)
    err := s.db.SelectContext(s.ctx, &doctors, `SELECT ` + getNovaDoctorCols() + ` FROM NOVA_DOCTOR WHERE MCR = :mcr`, mcr)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    if len(doctors) == 0 {
        return doctors, nil
    }

    ids := make([]int64, 0)
    for _, d := range doctors {
        ids = append(ids, d.DoctorID.Int64)
    }

    specialtyMap, err := s.findAllNovaSpecialtyMap(s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    spokenMap, _ := s.findAllNovaDoctorSpokenLanguage(s.db, ids)
    qualMap, _ := s.findAllNovaDoctorQualifications(s.db, ids)
    specMap, _ := s.findAllNovaDoctorSpecialities(s.db, ids)
    locMap, _ := s.findAllNovaDoctorClinicLocation(s.db, ids)
    hoursMap, _ := s.findAllNovaDoctorClinicHours(s.db, ids)
    contactMap, _ := s.findAllNovaDoctorContact(s.db, ids)
    doctorSpecialtyMap, _ := s.findAllNovaDoctorSpecialtyPrimary(s.db, ids, specialtyMap)

    for i, doc := range doctors {
        if list, ok := spokenMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorSpokenLanguage = list
        }
        if list, ok := qualMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorQualifications = list
        }
        if list, ok := specMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorSpecialities = list
        }
        if list, ok := locMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorClinicLocation = list
        }
        if list, ok := hoursMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorClinicHours = list
        }
        if list, ok := contactMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorContact = list
        }
        if list, ok := doctorSpecialtyMap[doc.DoctorID.Int64]; ok {
            doctors[i].DoctorSpecialty = list
            hasPrimary := false
            for _, sp := range list {
                if sp.PrimarySpecialtyV.Int32 == 1 {
                    hasPrimary = true
                    break
                }
            }
            if !hasPrimary {
                doctors[i].AllowAppointment.String = "N"
            }
        } else {
            doctors[i].AllowAppointment.String = "N"
        }
    }

    return doctors, nil
}

func (s *NovaDoctorService) findAllNovaDoctorSpokenLanguage(db *sqlx.DB, ids []int64) (map[int64][]model.NovaDoctorSpokenLanguage, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    query, args, err := sqlx.In(`SELECT ` + getNovaDoctorSpokenLanguageCols() + ` FROM NOVA_DOCTOR_SPOKEN_LANGUAGE WHERE DOCTOR_ID IN (?) ORDER BY DISPLAY_SEQUENCE`, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    query = db.Rebind(query)
    rows, err := db.QueryxContext(s.ctx, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    result := make(map[int64][]model.NovaDoctorSpokenLanguage)
    for rows.Next() {
        var item model.NovaDoctorSpokenLanguage
        if err := rows.StructScan(&item); err != nil {
            utils.LogError(err)
            return nil, err
        }
        result[item.DoctorID.Int64] = append(result[item.DoctorID.Int64], item)
    }
    return result, rows.Err()
}

func (s *NovaDoctorService) findAllNovaDoctorQualifications(db *sqlx.DB, ids []int64) (map[int64][]model.NovaDoctorQualifications, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    query, args, err := sqlx.In(`SELECT ` + getNovaDoctorQualificationsCols() + ` FROM NOVA_DOCTOR_QUALIFICATIONS WHERE DOCTOR_ID IN (?) ORDER BY DISPLAY_SEQUENCE`, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    query = db.Rebind(query)
    rows, err := db.QueryxContext(s.ctx, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    result := make(map[int64][]model.NovaDoctorQualifications)
    for rows.Next() {
        var item model.NovaDoctorQualifications
        if err := rows.StructScan(&item); err != nil {
            utils.LogError(err)
            return nil, err
        }
        result[item.DoctorID.Int64] = append(result[item.DoctorID.Int64], item)
    }
    return result, err
}

func (s *NovaDoctorService) findAllNovaDoctorSpecialities(db *sqlx.DB, ids []int64) (map[int64][]model.NovaDoctorSpecialities, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    query, args, err := sqlx.In(`SELECT ` + getNovaDoctorSpecialitiesCols() + ` FROM NOVA_DOCTOR_SPECIALITIES WHERE DOCTOR_ID IN (?) ORDER BY DISPLAY_SEQUENCE`, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    query = db.Rebind(query)
    rows, err := db.QueryxContext(s.ctx, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    result := make(map[int64][]model.NovaDoctorSpecialities)
    for rows.Next() {
        var item model.NovaDoctorSpecialities
        if err := rows.StructScan(&item); err != nil {
            utils.LogError(err)
            return nil, err
        }
        result[item.DoctorID.Int64] = append(result[item.DoctorID.Int64], item)
    }
    return result, rows.Err()
}

func (s *NovaDoctorService) findAllNovaDoctorClinicLocation(db *sqlx.DB, ids []int64) (map[int64][]model.NovaDoctorClinicLocation, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    query, args, err := sqlx.In(`SELECT ` + getNovaDoctorClinicLocationCols() + ` FROM NOVA_DOCTOR_CLINIC_LOCATION WHERE DOCTOR_ID IN (?)`, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    query = db.Rebind(query)
    rows, err := db.QueryxContext(s.ctx, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    result := make(map[int64][]model.NovaDoctorClinicLocation)
    for rows.Next() {
        var item model.NovaDoctorClinicLocation
        if err := rows.StructScan(&item); err != nil {
            utils.LogError(err)
            return nil, err
        }
        result[item.DoctorID.Int64] = append(result[item.DoctorID.Int64], item)
    }
    return result, rows.Err()
}

func (s *NovaDoctorService) findAllNovaDoctorClinicHours(db *sqlx.DB, ids []int64) (map[int64][]model.NovaDoctorClinicHours, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    query, args, err := sqlx.In(`SELECT ` + getNovaDoctorClinicHoursCols() + ` FROM NOVA_DOCTOR_CLINIC_HOURS WHERE DOCTOR_ID IN (?) ORDER BY DISPLAY_SEQUENCE`, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    query = db.Rebind(query)
    rows, err := db.QueryxContext(s.ctx, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    result := make(map[int64][]model.NovaDoctorClinicHours)
    for rows.Next() {
        var item model.NovaDoctorClinicHours
        if err := rows.StructScan(&item); err != nil {
            utils.LogError(err)
            return nil, err
        }
        result[item.DoctorID.Int64] = append(result[item.DoctorID.Int64], item)
    }
    return result, rows.Err()
}

func (s *NovaDoctorService) findAllNovaDoctorContact(db *sqlx.DB, ids []int64) (map[int64][]model.NovaDoctorContact, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    query, args, err := sqlx.In(`SELECT `+getNovaDoctorContactCols()+` FROM NOVA_DOCTOR_CONTACT WHERE DOCTOR_ID IN (?) ORDER BY DISPLAY_SEQUENCE`, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    query = db.Rebind(query)
    rows, err := db.QueryxContext(s.ctx, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    result := make(map[int64][]model.NovaDoctorContact)
    for rows.Next() {
        var item model.NovaDoctorContact
        if err := rows.StructScan(&item); err != nil {
            utils.LogError(err)
            return nil, err
        }
        result[item.DoctorID.Int64] = append(result[item.DoctorID.Int64], item)
    }
    return result, err
}

func (s *NovaDoctorService) findAllNovaDoctorAppointments(db *sqlx.DB, ids []int64) (map[int64][]model.NovaDoctorAppointment, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    query, args, err := sqlx.In(`SELECT `+getNovaDoctorAppointmentCols()+` FROM NOVA_DOCTOR_APPT_SLOT WHERE DOCTOR_ID IN (?) ORDER BY DISPLAY_SEQUENCE`, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    query = db.Rebind(query)
    rows, err := db.QueryxContext(s.ctx, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    result := make(map[int64][]model.NovaDoctorAppointment)
    for rows.Next() {
        var item model.NovaDoctorAppointment
        if err := rows.StructScan(&item); err != nil {
            utils.LogError(err)
            return nil, err
        }
        result[item.DoctorID.Int64] = append(result[item.DoctorID.Int64], item)
    }
    return result, err
}

func (s *NovaDoctorService) findAllNovaDoctorSpecialty(db *sqlx.DB, ids []int64, specialtyMap map[int64]model.NovaSpecialty) (map[int64][]model.NovaDoctorSpecialty, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    query, args, err := sqlx.In(`SELECT `+getNovaDoctorSpecialtyCols()+` FROM NOVA_DOCTOR_SPECIALTY WHERE DOCTOR_ID IN (?)`, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    query = db.Rebind(query)
    rows, err := db.QueryxContext(s.ctx, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    result := make(map[int64][]model.NovaDoctorSpecialty)
    for rows.Next() {
        var item model.NovaDoctorSpecialty
        if err := rows.StructScan(&item); err != nil {
            utils.LogError(err)
            return nil, err
        }
        if spec, ok := specialtyMap[item.SpecialtyID.Int64]; ok {
            item.Specialty = &spec
        }
        result[item.DoctorID.Int64] = append(result[item.DoctorID.Int64], item)
    }
    return result, err
}

func (s *NovaDoctorService) findAllNovaDoctorSpecialtyPrimary(db *sqlx.DB, ids []int64, specialtyMap map[int64]model.NovaSpecialty) (map[int64][]model.NovaDoctorSpecialty, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    query, args, err := sqlx.In(`SELECT `+getNovaDoctorSpecialtyCols()+` FROM NOVA_DOCTOR_SPECIALTY WHERE DOCTOR_ID IN (?) AND PRIMARY_SPECIALTY = 1`, ids)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    query = db.Rebind(query)
    rows, err := db.QueryxContext(s.ctx, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    result := make(map[int64][]model.NovaDoctorSpecialty)
    for rows.Next() {
        var item model.NovaDoctorSpecialty
        if err := rows.StructScan(&item); err != nil {
            utils.LogError(err)
            return nil, err
        }
        if spec, ok := specialtyMap[item.SpecialtyID.Int64]; ok {
            item.Specialty = &spec
        }
        result[item.DoctorID.Int64] = append(result[item.DoctorID.Int64], item)
    }
    return result, rows.Err()
}

func (s *NovaDoctorService) findAllNovaSpecialtyMap(q *sqlx.DB) (map[int64]model.NovaSpecialty, error) {
    query := `SELECT ` + getNovaSpecialtyCols() + ` FROM NOVA_SPECIALTY`
    rows, err := q.QueryxContext(s.ctx, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    result := make(map[int64]model.NovaSpecialty)
    for rows.Next() {
        var spec model.NovaSpecialty
        if err := rows.StructScan(&spec); err != nil {
            utils.LogError(err)
            return nil, err
        }
        result[spec.SpecialtyID.Int64] = spec
    }
    return result, rows.Err()
}

func (s *NovaDoctorService) saveDoctorClinicHoursTx(tx *sqlx.Tx, doctor *model.NovaDoctor) error {
    query := `
        INSERT INTO NOVA_DOCTOR_CLINIC_HOURS (
            CLINIC_HOUR_ID, BY_APPOINTMENT_ONLY, DAY_END_TIME,
            DAY_OF_THE_WEEK, DAY_START_TIME, DISPLAY_SEQUENCE, DOCTOR_ID
        ) VALUES (DR_CLINIC_HOURS_SEQ.NEXTVAL, :byAppt, :dayEndTime, :dayOfTheWeek, :dayStartTime, :displaySequence, :doctorId)
    `
    for _, item := range doctor.DoctorClinicHours {
        byAppt := 0
        if item.ByAppointmentOnlyV.Int32 == 1 {
            byAppt = 1
        }
        _, err := tx.ExecContext(s.ctx, query,
            sql.Named("byAppt", byAppt),
            sql.Named("dayEndTime", item.DayEndTime.String),
            sql.Named("dayOfTheWeek", item.DayOfTheWeek.String),
            sql.Named("dayStartTime", item.DayStartTime.String),
            sql.Named("displaySequence", item.DisplaySequence.Int32),
            sql.Named("doctorId", doctor.DoctorID.Int64),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func (s *NovaDoctorService) saveDoctorClinicLocationTx(tx *sqlx.Tx, doctor *model.NovaDoctor) error {
    query := `
        INSERT INTO NOVA_DOCTOR_CLINIC_LOCATION (
            CLINIC_LOCATION_ID, DOCTOR_ID, LOCATION, BUILDING
        ) VALUES (DR_CLINIC_LOCATION_SEQ.NEXTVAL, :doctorId, :location, :building)
    `
    for _, item := range doctor.DoctorClinicLocation {
        _, err := tx.ExecContext(s.ctx, query,
            sql.Named("doctorId", doctor.DoctorID.Int64),
            sql.Named("location", item.Location.String),
            sql.Named("building", item.Building.String),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func (s *NovaDoctorService) saveDoctorAppointmentTx(tx *sqlx.Tx, doctor *model.NovaDoctor) error {
    // First check for duplicates
    const checkQuery = `
        SELECT COUNT(*) FROM NOVA_DOCTOR_APPT_SLOT
        WHERE DOCTOR_ID = :doctorId AND DAY_OF_WEEK = :dayOfWeek AND SLOT_TYPE = :slotType AND SESSION_TYPE = :sessionType
    `
    const insertQuery = `
        INSERT INTO NOVA_DOCTOR_APPT_SLOT (
            DOCTOR_APPT_SLOT_ID, DOCTOR_ID, DAY_OF_WEEK, SLOT_TYPE,
            SESSION_TYPE, START_TIME, END_TIME, MAX_SLOTS, DISPLAY_SEQUENCE
        ) VALUES (DR_APPT_SLOT_ID_SEQ.NEXTVAL, :doctorId, :dayOfWeek, :slotType, :sessionType, :startTime, :endTime, :maxSlots, :displaySequence)
    `
    for _, item := range doctor.DoctorAppointment {
        var count int
        err := tx.GetContext(s.ctx, &count, checkQuery,
            sql.Named("doctorId", doctor.DoctorID.Int64),
            sql.Named("dayOfWeek", item.ApptDayOfWeek.String),
            sql.Named("slotType", item.ApptSlotType.String),
            sql.Named("sessionType", item.ApptSessionType.String),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
        if count > 0 {
            return fiber.NewError(fiber.StatusBadRequest, "Existing records with same appointment setup found")
        }
        _, err = tx.ExecContext(s.ctx, insertQuery,
            sql.Named("doctorId", doctor.DoctorID.Int64),
            sql.Named("dayOfWeek", item.ApptDayOfWeek.String),
            sql.Named("slotType", item.ApptSlotType.String),
            sql.Named("sessionType", item.ApptSessionType.String),
            sql.Named("startTime", item.ApptStartTime.String),
            sql.Named("endTime", item.ApptEndTime.String),
            sql.Named("maxSlots", item.ApptMaxSlots.Int32),
            sql.Named("displaySequence", item.DisplaySequence.Int32),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func (s *NovaDoctorService) saveDoctorContactsTx(tx *sqlx.Tx, doctor *model.NovaDoctor) error {
    query := `
        INSERT INTO NOVA_DOCTOR_CONTACT (
            CONTACT_ID, CONTACT_TYPE, CONTACT_VALUE, DISPLAY_SEQUENCE, DOCTOR_ID
        ) VALUES (DR_CONTACT_SEQ.NEXTVAL, :contactType, :contactValue, :displaySequence, :doctorId)
    `
    for _, item := range doctor.DoctorContact {
        _, err := tx.ExecContext(s.ctx, query,
            sql.Named("contactType", item.ContactType.String),
            sql.Named("contactValue", item.ContactValue.String),
            sql.Named("displaySequence", item.DisplaySequence.Int32),
            sql.Named("doctorId", doctor.DoctorID.Int64),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func (s *NovaDoctorService) saveDoctorQualificationsTx(tx *sqlx.Tx, doctor *model.NovaDoctor) error {
    query := `
        INSERT INTO NOVA_DOCTOR_QUALIFICATIONS (
            QUALIFICATION_ID, DISPLAY_SEQUENCE, DOCTOR_ID, QUALIFICATION
        ) VALUES (DR_QUALIFICATIONS_SEQ.NEXTVAL, :displaySequence, :doctorId, :qualification)
    `
    for _, item := range doctor.DoctorQualifications {
        _, err := tx.ExecContext(s.ctx, query,
            sql.Named("displaySequence", item.DisplaySequence.Int32),
            sql.Named("doctorId", doctor.DoctorID.Int64),
            sql.Named("qualification", item.Qualification.String),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func (s *NovaDoctorService) saveDoctorSpecialitiesTx(tx *sqlx.Tx, doctor *model.NovaDoctor) error {
    query := `
        INSERT INTO NOVA_DOCTOR_SPECIALITIES (
            SPECIALITIES_ID, DISPLAY_SEQUENCE, DOCTOR_ID, SPECIALITIES, SUBSPECIALTY
        ) VALUES (DOCTOR_SPECIALITIES_ID_SEQ.NEXTVAL, :displaySequence, :doctorId, :specialities, :subspecialty)
    `
    for _, item := range doctor.DoctorSpecialities {
        _, err := tx.ExecContext(s.ctx, query,
            sql.Named("displaySequence", item.DisplaySequence.Int32),
            sql.Named("doctorId", doctor.DoctorID.Int64),
            sql.Named("specialities", item.Specialities.String),
            sql.Named("subspecialty", item.Subspecialty.String),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func (s *NovaDoctorService) saveDoctorSpecialtyTx(tx *sqlx.Tx, doctor *model.NovaDoctor) error {
    query := `
        INSERT INTO NOVA_DOCTOR_SPECIALTY (
            DOCTOR_SPECIALTY_ID, DOCTOR_ID, PRIMARY_SPECIALTY, SPECIALTY_ID
        ) VALUES (DR_SPECIALTIES_SEQ.NEXTVAL, :doctorId, :primarySpecialty, :specialtyId)
    `
    for _, item := range doctor.DoctorSpecialty {
        primary := 0
        if item.PrimarySpecialtyV.Int32 == 1 {
            primary = 1
        }
        _, err := tx.ExecContext(s.ctx, query,
            sql.Named("doctorId", doctor.DoctorID.Int64),
            sql.Named("primarySpecialty", primary),
            sql.Named("specialtyId", item.SpecialtyID.Int64),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func (s *NovaDoctorService) saveDoctorSpokenLanguageTx(tx *sqlx.Tx, doctor *model.NovaDoctor) error {
    query := `
        INSERT INTO NOVA_DOCTOR_SPOKEN_LANGUAGE (
            SPOKEN_LANGUAGE_ID, DISPLAY_SEQUENCE, DOCTOR_ID, SPOKEN_LANGUAGE
        ) VALUES (DR_SPOKEN_LANGUAGE_SEQ.NEXTVAL, :displaySequence, :doctorId, :spokenLanguage)
    `
    for _, item := range doctor.DoctorSpokenLanguage {
        _, err := tx.ExecContext(s.ctx, query,
            sql.Named("displaySequence", item.DisplaySequence.Int32),
            sql.Named("doctorId", doctor.DoctorID.Int64),
            sql.Named("spokenLanguage", item.SpokenLanguage.String),
        )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}

func getNovaDoctorSpokenLanguageCols() string {
    return `
        SPOKEN_LANGUAGE_ID,
        DOCTOR_ID,
        DISPLAY_SEQUENCE,
        SPOKEN_LANGUAGE
    `
}

func getNovaDoctorQualificationsCols() string {
    return `
        QUALIFICATION_ID,
        DOCTOR_ID,
        DISPLAY_SEQUENCE,
        QUALIFICATION
    `
}

func getNovaDoctorSpecialitiesCols() string {
    return `
        SPECIALITIES_ID,
        DOCTOR_ID,
        DISPLAY_SEQUENCE,
        SPECIALITIES,
        SUBSPECIALTY
    `
}

func getNovaDoctorClinicLocationCols() string {
    return `
        CLINIC_LOCATION_ID,
        DOCTOR_ID,
        LOCATION,
        BUILDING
    `
}

func getNovaDoctorClinicHoursCols() string {
    return `
        CLINIC_HOUR_ID,
        DOCTOR_ID,
        DISPLAY_SEQUENCE,
        DAY_OF_THE_WEEK,
        DAY_START_TIME,
        DAY_END_TIME,
        BY_APPOINTMENT_ONLY
    `
}

func getNovaDoctorAppointmentCols() string {
    return `
        DOCTOR_APPT_SLOT_ID,
        DOCTOR_ID,
        DAY_OF_WEEK,
        SLOT_TYPE,
        SESSION_TYPE,
        START_TIME,
        END_TIME,
        MAX_SLOTS,
        DISPLAY_SEQUENCE
    `
}

func getNovaDoctorContactCols() string {
    return `
        CONTACT_ID,
        DOCTOR_ID,
        DISPLAY_SEQUENCE,
        CONTACT_TYPE,
        CONTACT_VALUE
    `
}

func getNovaSpecialtyCols() string {
    return `
        SPECIALTY_ID,
        SPECIALTY_CODE,
        SPECIALTY_DESC
    `
}

func getNovaDoctorSpecialtyCols() string {
    return `
        DOCTOR_SPECIALTY_ID,
        DOCTOR_ID,
        SPECIALTY_ID,
        PRIMARY_SPECIALTY
    `
}

func getNovaDoctorCols() string {
    return `
        DOCTOR_ID,
        MCR,
        NAME,
        GENDER,
        NATIONALITY,
        IMAGE,
        DISPLAY_SEQUENCE,
        ALLOW_APPOINTMENT,
        CONSULTANT_TYPE,
        IS_FOR_PACKAGE,
        QUALIFICATIONS_SHORT,
        REGISTRATION_NO
    `
}
