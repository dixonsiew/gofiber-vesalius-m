package applicationUserFamily

import (
	"context"
	"database/sql"
	"strings"
	"vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/service/applicationUser"
	sqx "vesaliusm/sql"
	"vesaliusm/utils"

	"github.com/guregu/null/v6"
	"github.com/jmoiron/sqlx"
	"github.com/nleeper/goment"
)

var ApplicationUserFamilySvc *ApplicationUserFamilyService = NewApplicationUserFamilyService(database.GetDb(), database.GetCtx(), database.GetDbrs())

type ApplicationUserFamilyService struct {
    db                     *sqlx.DB
    ctx                    context.Context
    dbrs                   *sqlx.DB
    applicationUserService *applicationUser.ApplicationUserService
}

func NewApplicationUserFamilyService(db *sqlx.DB, ctx context.Context, dbrs *sqlx.DB) *ApplicationUserFamilyService {
    return &ApplicationUserFamilyService{
        db:                     db,
        ctx:                    ctx,
        dbrs:                   dbrs,
        applicationUserService: applicationUser.ApplicationUserSvc,
    }
}

func (s *ApplicationUserFamilyService) FindByPRN(prn string, conn *sqlx.DB) (*model.ApplicationUserFamily, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM APPLICATION_USER_FAMILY WHERE NOK_PRN = :prn`
    query = strings.Replace(query, "*", utils.GetDbCols(model.ApplicationUserFamily{}, ""), 1)
    var o model.ApplicationUserFamily
    err := db.GetContext(s.ctx, &o, query, prn)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    return &o, err
}

func (s *ApplicationUserFamilyService) FindAllByUserId(userId int64, offset int, limit int, includeSelf bool, isForAppt bool, conn *sqlx.DB) ([]model.ApplicationUserFamily, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := ""
    if isForAppt {
        query = `
            SELECT * FROM APPLICATION_USER_FAMILY 
            WHERE PATIENT_PRN IN (SELECT MASTER_PRN FROM APPLICATION_USER WHERE USER_ID = :userId) AND IS_PATIENT = 'Y'
            ORDER BY NOK_REF_NUMBER DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
        `
    } else {
        query = `
            SELECT * FROM APPLICATION_USER_FAMILY 
            WHERE PATIENT_PRN IN (SELECT MASTER_PRN FROM APPLICATION_USER WHERE USER_ID = :userId)
            ORDER BY NOK_REF_NUMBER DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
        `
    }
    query = strings.Replace(query, "*", utils.GetDbCols(model.ApplicationUserFamily{}, ""), 1)
    list := make([]model.ApplicationUserFamily, 0)
    err := db.SelectContext(s.ctx, &list, query, userId, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }

    if includeSelf {
        o, err := s.applicationUserService.FindByUserId(userId, conn)
        if err != nil {
            return nil, err
        }

        if o != nil {
            f := new(model.ApplicationUserFamily)
            f.SetFromFamilyMember(*o)
            list = append([]model.ApplicationUserFamily{*f}, list...)
        }
    }

    return list, nil
}

func (s *ApplicationUserFamilyService) ListByUserId(userId int64, includeSelf bool, isForAppt bool, self string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountByUserId(userId, isForAppt, s.db)
    if err != nil {
        return nil, err
    }
    if includeSelf {
        total = total + 1
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllByUserId(userId, pager.GetLowerBound(), pager.PageSize, includeSelf, isForAppt, s.db)
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

func (s *ApplicationUserFamilyService) CountByUserId(userId int64, isForAppt bool, conn *sqlx.DB) (int, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    var count int
    query := `SELECT COUNT(AUF_ID) AS COUNT FROM APPLICATION_USER_FAMILY WHERE PATIENT_PRN IN (SELECT MASTER_PRN FROM APPLICATION_USER WHERE USER_ID = :userId)`
    if isForAppt {
        query = `SELECT COUNT(AUF_ID) AS COUNT FROM APPLICATION_USER_FAMILY WHERE PATIENT_PRN IN (SELECT MASTER_PRN FROM APPLICATION_USER WHERE USER_ID = :userId) AND IS_PATIENT = 'Y'`
    }

    err := db.GetContext(s.ctx, &count, query, userId)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ApplicationUserFamilyService) FindAllActiveByUserId(userId int64, offset int, limit int, includeSelf bool, isForAppt bool, conn *sqlx.DB) ([]model.ApplicationUserFamily, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := ""
    if isForAppt {
        query = `
            SELECT * FROM APPLICATION_USER_FAMILY 
            WHERE PATIENT_PRN IN (SELECT MASTER_PRN FROM APPLICATION_USER WHERE USER_ID = :userId) AND IS_PATIENT = 'Y' AND IS_ACTIVE = 'Y'
            ORDER BY NOK_REF_NUMBER DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
        `
    } else {
        query = `
            SELECT * FROM APPLICATION_USER_FAMILY 
            WHERE PATIENT_PRN IN (SELECT MASTER_PRN FROM APPLICATION_USER WHERE USER_ID = :userId) AND IS_ACTIVE = 'Y'
            ORDER BY NOK_REF_NUMBER DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
        `
    }
    query = strings.Replace(query, "*", utils.GetDbCols(model.ApplicationUserFamily{}, ""), 1)
    list := make([]model.ApplicationUserFamily, 0)
    err := db.SelectContext(s.ctx, &list, query, userId, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
        if list[i].Dob.Valid && list[i].Dob.String != "-" {
            g, _ := goment.New(list[i].Dob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].Dob = utils.NewNullString(g.Format("DD/MM/YYYY"))
        }
    }

    if includeSelf {
        o, err := s.applicationUserService.FindByUserId(userId, conn)
        if err != nil {
            return nil, err
        }

        if o != nil {
            f := new(model.ApplicationUserFamily)
            f.SetFromFamilyMember(*o)
            list = append([]model.ApplicationUserFamily{*f}, list...)
        }
    }

    return list, nil
}

func (s *ApplicationUserFamilyService) ListActiveByUserId(userId int64, includeSelf bool, isForAppt bool, self string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountActiveByUserId(userId, isForAppt, s.db)
    if err != nil {
        return nil, err
    }
    if includeSelf {
        total = total + 1
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllActiveByUserId(userId, pager.GetLowerBound(), pager.PageSize, includeSelf, isForAppt, s.db)
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

func (s *ApplicationUserFamilyService) CountActiveByUserId(userId int64, isForAppt bool, conn *sqlx.DB) (int, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    var count int
    query := `SELECT COUNT(AUF_ID) AS COUNT FROM APPLICATION_USER_FAMILY WHERE PATIENT_PRN IN (SELECT MASTER_PRN FROM APPLICATION_USER WHERE USER_ID = :userId) AND IS_ACTIVE = 'Y'`
    if isForAppt {
        query = `SELECT COUNT(AUF_ID) AS COUNT FROM APPLICATION_USER_FAMILY WHERE PATIENT_PRN IN (SELECT MASTER_PRN FROM APPLICATION_USER WHERE USER_ID = :userId) AND IS_PATIENT = 'Y' AND IS_ACTIVE = 'Y'`
    }

    err := db.GetContext(s.ctx, &count, query, userId)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ApplicationUserFamilyService) FindAllByUserIdAppt(userId int64, includeSelf bool, isForAppt bool, conn *sqlx.DB) ([]model.ApplicationUserFamily, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := ""
    if isForAppt {
        query = `
            SELECT * FROM APPLICATION_USER_FAMILY 
            WHERE PATIENT_PRN IN (SELECT MASTER_PRN FROM APPLICATION_USER WHERE USER_ID = :userId) 
            AND IS_PATIENT = 'Y'
            AND IS_ACTIVE = 'Y'
            ORDER BY NOK_REF_NUMBER DESC
        `
    } else {
        query = `
            SELECT * FROM APPLICATION_USER_FAMILY 
            WHERE PATIENT_PRN IN (SELECT MASTER_PRN FROM APPLICATION_USER WHERE USER_ID = :userId)
            AND IS_ACTIVE = 'Y'
            ORDER BY NOK_REF_NUMBER
        `
    }
    query = strings.Replace(query, "*", utils.GetDbCols(model.ApplicationUserFamily{}, ""), 1)
    list := make([]model.ApplicationUserFamily, 0)
    err := db.SelectContext(s.ctx, &list, query, userId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }

    if includeSelf {
        o, err := s.applicationUserService.FindByUserId(userId, conn)
        if err != nil {
            return nil, err
        }

        if o != nil {
            f := new(model.ApplicationUserFamily)
            f.SetFromFamilyMember(*o)
            list = append([]model.ApplicationUserFamily{*f}, list...)
        }
    }

    return list, nil
}

func (s *ApplicationUserFamilyService) FindAllByUserPrnAppt(prn string, includeSelf bool, isForAppt bool, conn *sqlx.DB) ([]model.ApplicationUserFamily, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := ""
    if isForAppt {
        query = `
            SELECT * FROM APPLICATION_USER_FAMILY 
            WHERE PATIENT_PRN = :prn 
            AND IS_PATIENT = 'Y'
            AND IS_ACTIVE = 'Y'
            ORDER BY NOK_REF_NUMBER DESC
        `
    } else {
        query = `
            SELECT * FROM APPLICATION_USER_FAMILY 
            WHERE PATIENT_PRN = :prn
            AND IS_ACTIVE = 'Y'
            ORDER BY NOK_REF_NUMBER
        `
    }
    query = strings.Replace(query, "*", utils.GetDbCols(model.ApplicationUserFamily{}, ""), 1)
    list := make([]model.ApplicationUserFamily, 0)
    err := db.SelectContext(s.ctx, &list, query, prn)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }

    if includeSelf {
        o, err := s.applicationUserService.FindByPRN(prn, conn)
        if err != nil {
            return nil, err
        }

        if o != nil {
            f := new(model.ApplicationUserFamily)
            f.SetFromFamilyMember(*o)
            list = append([]model.ApplicationUserFamily{*f}, list...)
        }
    }

    return list, nil
}

func (s *ApplicationUserFamilyService) FindByFamilyId(familyId int64) (*model.ApplicationUserFamily, error) {
    query := `SELECT * FROM APPLICATION_USER_FAMILY WHERE AUF_ID = :familyId`
    query = strings.Replace(query, "*", utils.GetDbCols(model.ApplicationUserFamily{}, ""), 1)
    var o model.ApplicationUserFamily
    err := s.db.GetContext(s.ctx, &o, query, familyId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    return &o, err
}

func (s *ApplicationUserFamilyService) SignupSync(prn string, userId int64) error {
    list := make([]model.PatientNOK, 0)
    err := s.dbrs.SelectContext(s.ctx, &list, sqx.GET_PATIENT_NOK_FAMILY, prn)
    if err != nil {
        utils.LogError(err)
        return err
    }
    for i := range list {
        var doc null.String
        err := s.dbrs.GetContext(s.ctx, &doc, sqx.GET_PATIENT_NOK_NRICPASSPORT, list[i].NOKPRN)
        if err != nil {
            if err == sql.ErrNoRows {
                continue
            }
            utils.LogError(err)
            return err
        }
        
        if doc.Valid {
            list[i].Set(doc.String)
        } else {
            list[i].Set("")
        }

        prm := []any{
            sql.Named("user_id", userId),
            sql.Named("patient_prn", prn),
            sql.Named("nok_ref_number", list[i].RefNo.Int64),
            sql.Named("is_patient", list[i].IsPatient.String),
            sql.Named("fullname", list[i].NOKFullname.String),
            sql.Named("relationship", list[i].NOKRelationship.String),
            sql.Named("nok_prn", list[i].NOKPRN.String),
            sql.Named("doc_number", list[i].NOKDocNumber.String),
            sql.Named("nric_passport", list[i].NRICPassport.String),
            sql.Named("dob", list[i].NOKDOB.String),
            sql.Named("gender", list[i].NOKGender.String),
            sql.Named("nationality", list[i].NOKNationality.String),
            sql.Named("contact_number", list[i].NOKContact.String),
            sql.Named("address", list[i].NOKAddress.String),
            sql.Named("is_active", "Y"),
            sql.Named("marital_status", list[i].MaritalStatus.String),
            sql.Named("email", list[i].Email.String),
        }
        q := `
            DECLARE v_count NUMBER;
             BEGIN
              SELECT COUNT(*) INTO v_count FROM APPLICATION_USER_FAMILY
              WHERE PATIENT_PRN = :patient_prn
              AND DOC_NUMBER = :doc_number;

              IF v_count = 0 THEN
                INSERT INTO APPLICATION_USER_FAMILY 
                (
                  USER_ID, PATIENT_PRN, NOK_REF_NUMBER, IS_PATIENT, FULLNAME, RELATIONSHIP, 
                  NOK_PRN, DOC_NUMBER, NRIC_PASSPORT, DOB, GENDER, NATIONALITY, 
                  CONTACT_NUMBER, ADDRESS, MARITAL_STATUS, EMAIL, IS_ACTIVE
                )
                VALUES 
                (
                  :user_id, :patient_prn, :nok_ref_number, :is_patient, :fullname, :relationship, 
                  :nok_prn, :doc_number, :nric_passport, :dob, :gender, :nationality, 
                  :contact_number, :address, :marital_status, :email, :is_active
                );

              ELSIF v_count = 1 THEN
                UPDATE APPLICATION_USER_FAMILY SET
                  USER_ID = :user_id
                WHERE PATIENT_PRN = :patient_prn
                AND DOC_NUMBER = :doc_number;
                
              END IF;
             END;
        `
        _, err = s.db.ExecContext(s.ctx, q, prm...)
        if err != nil {
            utils.LogError(err)
            return err
        }
    }
    return nil
}
