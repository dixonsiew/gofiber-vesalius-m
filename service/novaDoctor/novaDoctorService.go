package novaDoctor

import (
    "context"
    "database/sql"
    "fmt"
    "strconv"
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
    var doctorId int64
    err = s.db.GetContext(s.ctx, &doctorId, "SELECT DOCTOR_ID_SEQ.NEXTVAL FROM DUAL")
    if err != nil {
        utils.LogError(err)
        return err
    }
    doctor.DoctorId = utils.NewInt64(doctorId)

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
        sql.Named("doctorId", doctor.DoctorId.Int64),
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

func (s *NovaDoctorService) ResizeAllDoctorImage(image string, doctorId int64) error {
    _, err := s.db.ExecContext(s.ctx,
        `UPDATE NOVA_DOCTOR SET RESIZE_IMAGE = :image WHERE DOCTOR_ID = :doctorId`,
        image, doctorId,
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
        sql.Named("doctorId", doctor.DoctorId.Int64),
    )
    if err != nil {
        return err
    }

    // Delete existing child records
    if err := s.deleteChildRecordsTx(tx, doctor.DoctorId.Int64); err != nil {
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
    ids := make([]string, 0)
    for i, d := range doctors {
        ids = append(ids, strconv.FormatInt(d.DoctorId.Int64, 10))
        // Set showMakeAppointmentButton based on config
        if config.GetIpayTestEnv() == "Y" && d.AllowAppointment.String == "Y" {
            doctors[i].ShowMakeAppointmentButton = "Y"
        } else {
            doctors[i].ShowMakeAppointmentButton = "N"
        }
    }

    doctorIds := strings.Join(ids, ",")
    // Fetch child data in bulk
    ms, err := s.findAllNovaSpecialtyMap(s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    m3, err := s.FindAllNovaDoctorSpecialities(doctorIds, s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    m4, err := s.FindAllNovaDoctorClinicLocation(doctorIds, s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    m5, err := s.FindAllNovaDoctorContact(doctorIds, s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    m7, err := s.FindAllNovaDoctorSpecialty(doctorIds, ms, s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    // Assign child collections
    for i, doc := range doctors {
        if list, ok := m3[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorSpecialities = list
        }
        if list, ok := m4[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorClinicLocation = list
        }
        if list, ok := m5[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorContact = list
        }
        if list, ok := m7[doc.DoctorId.Int64]; ok {
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
                doctors[i].AllowAppointment = utils.NewNullString("N")
            }
        } else {
            doctors[i].AllowAppointment = utils.NewNullString("N")
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

    ids := make([]string, len(doctors))
    for _, d := range doctors {
        ids = append(ids, strconv.FormatInt(d.DoctorId.Int64, 10))
    }

    doctorIds := strings.Join(ids, ",")
    ms, err := s.findAllNovaSpecialtyMap(s.db)
    if err != nil {
        return nil, err
    }
    m3, err := s.FindAllNovaDoctorSpecialities(doctorIds, s.db)
    if err != nil {
        return nil, err
    }
    m4, err := s.FindAllNovaDoctorClinicLocation(doctorIds, s.db)
    if err != nil {
        return nil, err
    }
    m5, err := s.FindAllNovaDoctorContact(doctorIds, s.db)
    if err != nil {
        return nil, err
    }
    m7, err := s.FindAllNovaDoctorSpecialty(doctorIds, ms, s.db)
    if err != nil {
        return nil, err
    }

    for i, doc := range doctors {
        if list, ok := m3[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorSpecialities = list
        }
        if list, ok := m4[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorClinicLocation = list
        }
        if list, ok := m5[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorContact = list
        }
        if list, ok := m7[doc.DoctorId.Int64]; ok {
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

    ids := []string{strconv.FormatInt(doctorId, 10)}
    doctorIds := strings.Join(ids, ",")

    ms, err := s.findAllNovaSpecialtyMap(s.db)
    if err != nil {
        return nil, err
    }

    // Fetch all child data
    m1, _ := s.FindAllNovaDoctorSpokenLanguage(doctorIds, s.db)
    m2, _ := s.FindAllNovaDoctorQualifications(doctorIds, s.db)
    m3, _ := s.FindAllNovaDoctorSpecialities(doctorIds, s.db)
    m4, _ := s.FindAllNovaDoctorClinicLocation(doctorIds, s.db)
    m5, _ := s.FindAllNovaDoctorClinicHours(doctorIds, s.db)
    m6, _ := s.FindAllNovaDoctorContact(doctorIds, s.db)
    m7, _ := s.FindAllNovaDoctorSpecialty(doctorIds, ms, s.db)
    m8, _ := s.FindAllNovaDoctorAppointments(doctorIds, s.db)

    if list, ok := m1[doctorId]; ok {
        doctor.DoctorSpokenLanguage = list
    }
    if list, ok := m2[doctorId]; ok {
        doctor.DoctorQualifications = list
    }
    if list, ok := m3[doctorId]; ok {
        doctor.DoctorSpecialities = list
    }
    if list, ok := m4[doctorId]; ok {
        doctor.DoctorClinicLocation = list
    }
    if list, ok := m5[doctorId]; ok {
        doctor.DoctorClinicHours = list
    }
    if list, ok := m6[doctorId]; ok {
        doctor.DoctorContact = list
    }
    if list, ok := m7[doctorId]; ok {
        doctor.DoctorSpecialty = list
    }
    if list, ok := m8[doctorId]; ok {
        doctor.DoctorAppointment = list
    }

    return &doctor, nil
}

func (s *NovaDoctorService) FindDoctorNameByDoctorId(doctorId int64, conn *sqlx.DB) (string, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT NAME FROM NOVA_DOCTOR WHERE DOCTOR_ID = :doctorId`
    var name string
    err := db.GetContext(s.ctx, &name, query, doctorId)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", nil
        }
        utils.LogError(err)
        return "", err
    }
    return name, nil
}

func (s *NovaDoctorService) FindDoctorIdByMcr(mcr string, conn *sqlx.DB) (int64, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT DOCTOR_ID FROM NOVA_DOCTOR WHERE MCR = :mcr`
    var id int64
    err := db.GetContext(s.ctx, &id, query, mcr)
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
    query := `SELECT * FROM NOVA_DOCTOR WHERE MCR = :mcr`
    query = strings.Replace(query, "*", getNovaDoctorCols(), 1)
    var doctor model.NovaDoctor
    err := s.db.GetContext(s.ctx, &doctor, query, mcr)
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
    query := `SELECT * FROM NOVA_DOCTOR WHERE MCR = :mcr`
    query = strings.Replace(query, "*", getNovaDoctorCols(), 1)
    doctors := make([]model.NovaDoctor, 0)
    err := s.db.SelectContext(s.ctx, &doctors, query, mcr)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    if len(doctors) == 0 {
        return doctors, nil
    }

    ids := make([]string, 0)
    for _, d := range doctors {
        ids = append(ids, strconv.FormatInt(d.DoctorId.Int64, 10))
    }

    ms, err := s.findAllNovaSpecialtyMap(s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    doctorIds := strings.Join(ids, ",")
    m1, _ := s.FindAllNovaDoctorSpokenLanguage(doctorIds, s.db)
    m2, _ := s.FindAllNovaDoctorQualifications(doctorIds, s.db)
    m3, _ := s.FindAllNovaDoctorSpecialities(doctorIds, s.db)
    m4, _ := s.FindAllNovaDoctorClinicLocation(doctorIds, s.db)
    m5, _ := s.FindAllNovaDoctorClinicHours(doctorIds, s.db)
    m6, _ := s.FindAllNovaDoctorContact(doctorIds, s.db)
    m7, _ := s.FindAllNovaDoctorSpecialtyPrimary(doctorIds, ms, s.db)

    for i, doc := range doctors {
        if list, ok := m1[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorSpokenLanguage = list
        }
        if list, ok := m2[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorQualifications = list
        }
        if list, ok := m3[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorSpecialities = list
        }
        if list, ok := m4[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorClinicLocation = list
        }
        if list, ok := m5[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorClinicHours = list
        }
        if list, ok := m6[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorContact = list
        }
        if list, ok := m7[doc.DoctorId.Int64]; ok {
            doctors[i].DoctorSpecialty = list
            hasPrimary := false
            for _, sp := range list {
                if sp.PrimarySpecialtyV.Int32 == 1 {
                    hasPrimary = true
                    break
                }
            }
            if !hasPrimary {
                doctors[i].AllowAppointment = utils.NewNullString("N")
            }
        } else {
            doctors[i].AllowAppointment = utils.NewNullString("N")
        }
    }

    return doctors, nil
}

func (s *NovaDoctorService) FindAllNovaDoctorSpokenLanguage(doctorIds string, db *sqlx.DB) (map[int64][]model.NovaDoctorSpokenLanguage, error) {
    query := fmt.Sprintf(`SELECT * FROM NOVA_DOCTOR_SPOKEN_LANGUAGE WHERE DOCTOR_ID IN (%s) ORDER BY DISPLAY_SEQUENCE`, doctorIds)
    query = strings.Replace(query, "*", getNovaDoctorSpokenLanguageCols(), 1)
    rows, err := db.QueryxContext(s.ctx, query)
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
        result[item.DoctorId.Int64] = append(result[item.DoctorId.Int64], item)
    }
    return result, rows.Err()
}

func (s *NovaDoctorService) FindAllNovaDoctorQualifications(doctorIds string, db *sqlx.DB) (map[int64][]model.NovaDoctorQualifications, error) {
    query := fmt.Sprintf(`SELECT * FROM NOVA_DOCTOR_QUALIFICATIONS WHERE DOCTOR_ID IN (%s) ORDER BY DISPLAY_SEQUENCE`, doctorIds)
    query = strings.Replace(query, "*", getNovaDoctorQualificationsCols(), 1)
    rows, err := db.QueryxContext(s.ctx, query)
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
        result[item.DoctorId.Int64] = append(result[item.DoctorId.Int64], item)
    }
    return result, err
}

func (s *NovaDoctorService) FindAllNovaDoctorSpecialities(doctorIds string, db *sqlx.DB) (map[int64][]model.NovaDoctorSpecialities, error) {
    query := fmt.Sprintf(`SELECT * FROM NOVA_DOCTOR_SPECIALITIES WHERE DOCTOR_ID IN (%s) ORDER BY DISPLAY_SEQUENCE`, doctorIds)
    query = strings.Replace(query, "*", getNovaDoctorSpecialitiesCols(), 1)
    rows, err := db.QueryxContext(s.ctx, query)
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
        result[item.DoctorId.Int64] = append(result[item.DoctorId.Int64], item)
    }
    return result, rows.Err()
}

func (s *NovaDoctorService) FindAllNovaDoctorClinicLocation(doctorIds string, db *sqlx.DB) (map[int64][]model.NovaDoctorClinicLocation, error) {
    query := fmt.Sprintf(`SELECT * FROM NOVA_DOCTOR_CLINIC_LOCATION WHERE DOCTOR_ID IN (%s)`, doctorIds)
    query = strings.Replace(query, "*", getNovaDoctorClinicLocationCols(), 1)
    rows, err := db.QueryxContext(s.ctx, query)
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
        result[item.DoctorId.Int64] = append(result[item.DoctorId.Int64], item)
    }
    return result, rows.Err()
}

func (s *NovaDoctorService) FindAllNovaDoctorClinicHours(doctorIds string, db *sqlx.DB) (map[int64][]model.NovaDoctorClinicHours, error) {
    query := fmt.Sprintf(`SELECT * FROM NOVA_DOCTOR_CLINIC_HOURS WHERE DOCTOR_ID IN (%s) ORDER BY DISPLAY_SEQUENCE`, doctorIds)
    query = strings.Replace(query, "*", getNovaDoctorClinicHoursCols(), 1)
    rows, err := db.QueryxContext(s.ctx, query)
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
        result[item.DoctorId.Int64] = append(result[item.DoctorId.Int64], item)
    }
    return result, rows.Err()
}

func (s *NovaDoctorService) FindAllNovaDoctorContact(doctorIds string, db *sqlx.DB) (map[int64][]model.NovaDoctorContact, error) {
    query := fmt.Sprintf(`SELECT * FROM NOVA_DOCTOR_CONTACT WHERE DOCTOR_ID IN (%s) ORDER BY DISPLAY_SEQUENCE`, doctorIds)
    query = strings.Replace(query, "*", getNovaDoctorContactCols(), 1)
    rows, err := db.QueryxContext(s.ctx, query)
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
        result[item.DoctorId.Int64] = append(result[item.DoctorId.Int64], item)
    }
    return result, err
}

func (s *NovaDoctorService) FindAllNovaDoctorAppointments(doctorIds string, db *sqlx.DB) (map[int64][]model.NovaDoctorAppointment, error) {
    query := fmt.Sprintf(`SELECT * FROM NOVA_DOCTOR_APPT_SLOT WHERE DOCTOR_ID IN (%s) ORDER BY DISPLAY_SEQUENCE`, doctorIds)
    query = strings.Replace(query, "*", getNovaDoctorAppointmentCols(), 1)
    rows, err := db.QueryxContext(s.ctx, query)
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
        result[item.DoctorId.Int64] = append(result[item.DoctorId.Int64], item)
    }
    return result, err
}

func (s *NovaDoctorService) FindAllNovaDoctorSpecialty(doctorIds string, specialtyMap map[int64]model.NovaSpecialty, db *sqlx.DB) (map[int64][]model.NovaDoctorSpecialty, error) {
    query := fmt.Sprintf(`SELECT * FROM NOVA_DOCTOR_SPECIALTY WHERE DOCTOR_ID IN (%s)`, doctorIds)
    query = strings.Replace(query, "*", getNovaDoctorSpecialtyCols(), 1)
    rows, err := db.QueryxContext(s.ctx, query)
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
        if spec, ok := specialtyMap[item.SpecialtyId.Int64]; ok {
            item.Specialty = &spec
        }
        result[item.DoctorId.Int64] = append(result[item.DoctorId.Int64], item)
    }
    return result, err
}
func (s *NovaDoctorService) FindAllNovaDoctorSpecialtyPrimary(doctorIds string, specialtyMap map[int64]model.NovaSpecialty, db *sqlx.DB) (map[int64][]model.NovaDoctorSpecialty, error) {
    query := fmt.Sprintf(`SELECT * FROM NOVA_DOCTOR_SPECIALTY WHERE DOCTOR_ID IN (%s) AND PRIMARY_SPECIALTY = 1`, doctorIds)
    query = strings.Replace(query, "*", getNovaDoctorSpecialtyCols(), 1)
    rows, err := db.QueryxContext(s.ctx, query)
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
        if spec, ok := specialtyMap[item.SpecialtyId.Int64]; ok {
            item.Specialty = &spec
        }
        result[item.DoctorId.Int64] = append(result[item.DoctorId.Int64], item)
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
        result[spec.SpecialtyId.Int64] = spec
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
            sql.Named("doctorId", doctor.DoctorId.Int64),
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
            sql.Named("doctorId", doctor.DoctorId.Int64),
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
            sql.Named("doctorId", doctor.DoctorId.Int64),
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
            sql.Named("doctorId", doctor.DoctorId.Int64),
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
            sql.Named("doctorId", doctor.DoctorId.Int64),
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
            sql.Named("doctorId", doctor.DoctorId.Int64),
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
            sql.Named("doctorId", doctor.DoctorId.Int64),
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
            sql.Named("doctorId", doctor.DoctorId.Int64),
            sql.Named("primarySpecialty", primary),
            sql.Named("specialtyId", item.SpecialtyId.Int64),
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
            sql.Named("doctorId", doctor.DoctorId.Int64),
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
