package clubs

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/model/clubs"
    "vesaliusm/service/applicationUser"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    "github.com/nleeper/goment"
    go_ora "github.com/sijms/go-ora/v2"
)

var ClubSvc *ClubService = NewClubService(database.GetDb(), database.GetCtx(), database.GetDbrs(), database.GetCtxrs())

type ClubService struct {
    db                     *sqlx.DB
    dbrs                   *sqlx.DB
    ctx                    context.Context
    ctxrs                  context.Context
    applicationUserService *applicationUser.ApplicationUserService
}

func NewClubService(db *sqlx.DB, ctx context.Context, dbrs *sqlx.DB, ctxrs context.Context) *ClubService {
    return &ClubService{
        db:                     db,
        dbrs:                   dbrs,
        ctx:                    ctx,
        ctxrs:                  ctxrs,
        applicationUserService: applicationUser.ApplicationUserSvc,
    }
}

func (s *ClubService) FindLittleKidsMembershipById(kidsMembershipId int64) ([]clubs.LittleExplorersKidsMembership, error) {
    query := `SELECT * FROM KIDS_CLUB_MEMBERSHIP WHERE KIDS_MEMBERSHIP_ID = :kidsMembershipId`
    list := make([]clubs.LittleExplorersKidsMembership, 0)
    err := s.db.SelectContext(s.ctx, &list, query, kidsMembershipId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
}

func (s *ClubService) FindGoldenPearlMembershipById(goldenMembershipId int64) (*clubs.GoldenPearlMembership, error) {
    query := `SELECT * FROM GOLDEN_CLUB_MEMBERSHIP WHERE GOLDEN_MEMBERSHIP_ID = :goldenMembershipId`
    var o clubs.GoldenPearlMembership
    err := s.db.GetContext(s.ctx, &o, query, goldenMembershipId)
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

func (s *ClubService) findGuestLittleKidsMembershipByIc(identificationNumber string) (*clubs.LittleExplorersKidsMembership, error) {
    query := `SELECT * FROM KIDS_CLUB_MEMBERSHIP WHERE KIDS_DOC_NUMBER = :identificationNumber`
    var o clubs.LittleExplorersKidsMembership
    err := s.db.GetContext(s.ctx, &o, query, identificationNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    if o.KidsMembershipJoinDate.Valid {
        g, _ := goment.New(o.KidsMembershipJoinDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.KidsMembershipJoinDate.String = g.Format("DD/MM/YYYY HH:mm")
    }
    if o.KidsDob.Valid {
        g, _ := goment.New(o.KidsDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.KidsDob.String = g.Format("DD/MM/YYYY")
    }
    return &o, err
}

func (s *ClubService) FindGuestGoldenPearlMembershipByIc(identificationNumber string) (*clubs.GoldenPearlMembership, error) {
    query := `SELECT * FROM GOLDEN_CLUB_MEMBERSHIP WHERE GOLDEN_DOC_NUMBER = :identificationNumber`
    var o clubs.GoldenPearlMembership
    err := s.db.GetContext(s.ctx, &o, query, identificationNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    if o.GoldenMembershipJoinDate.Valid {
        g, _ := goment.New(o.GoldenMembershipJoinDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.GoldenMembershipJoinDate.String = g.Format("DD/MM/YYYY HH:mm")
    }
    if o.GoldenDob.Valid {
        g, _ := goment.New(o.GoldenDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.GoldenDob.String = g.Format("DD/MM/YYYY")
    }
    return &o, nil
}

func (s *ClubService) FindLittleKidsMembershipByMembershipId(membershipId int64) (*clubs.LittleExplorersKidsMembership, error) {
    query := `SELECT * FROM KIDS_CLUB_MEMBERSHIP WHERE KIDS_MEMBERSHIP_ID = :membershipId`
    var o clubs.LittleExplorersKidsMembership
    err := s.db.GetContext(s.ctx, &o, query, membershipId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.SetWebAdmin()
    if o.KidsDob.Valid {
        g, _ := goment.New(o.KidsDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.KidsDob.String = g.Format("DD/MM/YYYY")
    }
    if o.GuardianDob.Valid {
        g, _ := goment.New(o.GuardianDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.GuardianDob.String = g.Format("DD/MM/YYYY")
    }
    return &o, nil
}

func (s *ClubService) FindGoldenPearlMembershipByMembershipId(membershipId int64) (*clubs.GoldenPearlMembership, error) {
    query := `SELECT * FROM GOLDEN_CLUB_MEMBERSHIP WHERE GOLDEN_MEMBERSHIP_ID = :membershipId`
    var o clubs.GoldenPearlMembership
    err := s.db.GetContext(s.ctx, &o, query, membershipId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.SetWebAdmin()
    if o.GoldenDob.Valid {
        g, _ := goment.New(o.GoldenDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.GoldenDob.String = g.Format("DD/MM/YYYY HH:mm")
    }
    if o.NokDob.Valid {
        g, _ := goment.New(o.NokDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.NokDob.String = g.Format("DD/MM/YYYY")
    }
    return &o, nil
}

func (s *ClubService) SaveLittleKidsMembership(o clubs.LittleExplorersKidsMembership) error {
    type Prns struct {
        kidsPrn     string
        guardianPrn string
    }
    const query = `
        SELECT PRN FROM NOVA_PATIENT_DOCUMENT
        WHERE DOCUMENT_TYPE = 'NRIC / Passport' AND
        DOCUMENT_NUMBER = :doc
    `
    var (
        prns Prns
        prn  string
    )
    err := s.dbrs.GetContext(s.ctxrs, &prn, query, o.KidsDocNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            prn = ""
        } else {
            utils.LogError(err)
            return err
        }
    }
    if prn != "" {
        prns.kidsPrn = prn
    }

    err = s.dbrs.GetContext(s.ctxrs, &prn, query, o.GuardianDocNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            prn = ""
        } else {
            utils.LogError(err)
            return err
        }
    }
    if prn != "" {
        prns.guardianPrn = prn
    }

    q := `
        INSERT INTO KIDS_CLUB_MEMBERSHIP
        (KIDS_MEMBERSHIP_NUMBER, KIDS_PRN, KIDS_NAME, KIDS_DOB, 
         KIDS_DOC_TYPE, KIDS_DOC_NUMBER, KIDS_GENDER, KIDS_NATIONALITY, 
         KIDS_EMAIL, GUARDIAN_PRN, GUARDIAN_NAME, GUARDIAN_DOB, GUARDIAN_DOC_TYPE, 
         GUARDIAN_DOC_NUMBER, GUARDIAN_GENDER, GUARDIAN_NATIONALITY, GUARDIAN_EMAIL, 
         GUARDIAN_HOME_CONTACT, GUARDIAN_MOBILE_CONTACT, GUARDIAN_ADDRESS1, GUARDIAN_ADDRESS2,
         GUARDIAN_ADDRESS3, GUARDIAN_POSTCODE, GUARDIAN_STATE, GUARDIAN_COUNTRY_CODE, 
         RELATIONSHIP, PREFERRED_LANGUAGE
        ) VALUES
        (:kidsMembershipNumber, :kidsPrn, :kidsName, TO_DATE(:kidsDob, 'DD/MM/YYYY'),
         :kidsDocType, :kidsDocNumber, :kidsGender, :kidsNationality,
         :kidsEmail, :guardianPrn, :guardianName, TO_DATE(:guardianDob, 'DD/MM/YYYY'), :guardianDocType,
         :guardianDocNumber, :guardianGender, :guardianNationality, :guardianEmail,
         :guardianHomeContact, :guardianMobileContact, :guardianAddress1, :guardianAddress2,
         :guardianAddress3, :guardianPostCode, :guardianState, :guardianCountryCode,
         :relationship, :preferredLanguage)
    `
    _, err = s.db.ExecContext(s.ctx, q,
        sql.Named("kidsMembershipNumber", o.KidsMembershipNumber.String),
        sql.Named("kidsPrn", prns.kidsPrn),
        sql.Named("kidsName", o.KidsName.String),
        sql.Named("kidsDob", o.KidsDob.String),
        sql.Named("kidsDocType", o.KidsDocType.String),
        sql.Named("kidsDocNumber", o.KidsDocNumber.String),
        sql.Named("kidsGender", o.KidsGender.String),
        sql.Named("kidsNationality", o.KidsNationality.String),
        sql.Named("kidsEmail", o.KidsEmail.String),
        sql.Named("guardianPrn", prns.guardianPrn),
        sql.Named("guardianName", o.GuardianName.String),
        sql.Named("guardianDob", o.GuardianDob.String),
        sql.Named("guardianDocType", o.GuardianDocType.String),
        sql.Named("guardianDocNumber", o.GuardianDocNumber.String),
        sql.Named("guardianGender", o.GuardianGender.String),
        sql.Named("guardianNationality", o.GuardianNationality.String),
        sql.Named("guardianEmail", o.GuardianEmail.String),
        sql.Named("guardianHomeContact", o.GuardianHomeContact.String),
        sql.Named("guardianMobileContact", o.GuardianMobileContact.String),
        sql.Named("guardianAddress1", o.GuardianAddress1.String),
        sql.Named("guardianAddress2", o.GuardianAddress2.String),
        sql.Named("guardianAddress3", o.GuardianAddress3.String),
        sql.Named("guardianPostCode", o.GuardianPostCode.String),
        sql.Named("guardianState", o.GuardianState.String),
        sql.Named("guardianCountryCode", o.GuardianCountryCode.String),
        sql.Named("relationship", o.Relationship.String),
        sql.Named("preferredLanguage", o.PreferredLanguage.String),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *ClubService) SaveLittleKidsMembershipViaWebportal(o clubs.LittleExplorersKidsMembership, adminId int64) error {
    type Prns struct {
        kidsPrn     string
        guardianPrn string
    }
    const query = `
        SELECT PRN FROM NOVA_PATIENT_DOCUMENT
        WHERE DOCUMENT_TYPE = 'NRIC / Passport' AND
        DOCUMENT_NUMBER = :doc
    `
    var (
        prns Prns
        prn  string
    )
    err := s.dbrs.GetContext(s.ctxrs, &prn, query, o.KidsDocNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            prn = ""
        } else {
            utils.LogError(err)
            return err
        }
    }
    if prn != "" {
        prns.kidsPrn = prn
    }

    err = s.dbrs.GetContext(s.ctxrs, &prn, query, o.GuardianDocNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            prn = ""
        } else {
            utils.LogError(err)
            return err
        }
    }
    if prn != "" {
        prns.guardianPrn = prn
    }

    q := `
        INSERT INTO KIDS_CLUB_MEMBERSHIP
        (KIDS_MEMBERSHIP_NUMBER, KIDS_PRN, KIDS_NAME, KIDS_DOB, 
         KIDS_DOC_TYPE, KIDS_DOC_NUMBER, KIDS_GENDER, KIDS_NATIONALITY, 
         KIDS_EMAIL, GUARDIAN_PRN, GUARDIAN_NAME, GUARDIAN_DOB, GUARDIAN_DOC_TYPE, 
         GUARDIAN_DOC_NUMBER, GUARDIAN_GENDER, GUARDIAN_NATIONALITY, GUARDIAN_EMAIL, 
         GUARDIAN_HOME_CONTACT, GUARDIAN_MOBILE_CONTACT, GUARDIAN_ADDRESS1, GUARDIAN_ADDRESS2,
         GUARDIAN_ADDRESS3, GUARDIAN_POSTCODE, GUARDIAN_STATE, GUARDIAN_COUNTRY_CODE, 
         RELATIONSHIP, PREFERRED_LANGUAGE, USER_CREATE
        ) VALUES
        (:kidsMembershipNumber, :kidsPrn, :kidsName, TO_DATE(:kidsDob, 'DD/MM/YYYY'),
         :kidsDocType, :kidsDocNumber, :kidsGender, :kidsNationality,
         :kidsEmail, :guardianPrn, :guardianName, TO_DATE(:guardianDob, 'DD/MM/YYYY'), :guardianDocType,
         :guardianDocNumber, :guardianGender, :guardianNationality, :guardianEmail,
         :guardianHomeContact, :guardianMobileContact, :guardianAddress1, :guardianAddress2,
         :guardianAddress3, :guardianPostCode, :guardianState, :guardianCountryCode,
         :relationship, :preferredLanguage, :adminId)
    `
    _, err = s.db.ExecContext(s.ctx, q,
        sql.Named("kidsMembershipNumber", o.KidsMembershipNumber.String),
        sql.Named("kidsPrn", prns.kidsPrn),
        sql.Named("kidsName", o.KidsName.String),
        sql.Named("kidsDob", o.KidsDob.String),
        sql.Named("kidsDocType", o.KidsDocType.String),
        sql.Named("kidsDocNumber", o.KidsDocNumber.String),
        sql.Named("kidsGender", o.KidsGender.String),
        sql.Named("kidsNationality", o.KidsNationality.String),
        sql.Named("kidsEmail", o.KidsEmail.String),
        sql.Named("guardianPrn", prns.guardianPrn),
        sql.Named("guardianName", o.GuardianName.String),
        sql.Named("guardianDob", o.GuardianDob.String),
        sql.Named("guardianDocType", o.GuardianDocType.String),
        sql.Named("guardianDocNumber", o.GuardianDocNumber.String),
        sql.Named("guardianGender", o.GuardianGender.String),
        sql.Named("guardianNationality", o.GuardianNationality.String),
        sql.Named("guardianEmail", o.GuardianEmail.String),
        sql.Named("guardianHomeContact", o.GuardianHomeContact.String),
        sql.Named("guardianMobileContact", o.GuardianMobileContact.String),
        sql.Named("guardianAddress1", o.GuardianAddress1.String),
        sql.Named("guardianAddress2", o.GuardianAddress2.String),
        sql.Named("guardianAddress3", o.GuardianAddress3.String),
        sql.Named("guardianPostCode", o.GuardianPostCode.String),
        sql.Named("guardianState", o.GuardianState.String),
        sql.Named("guardianCountryCode", o.GuardianCountryCode.String),
        sql.Named("relationship", o.Relationship.String),
        sql.Named("preferredLanguage", o.PreferredLanguage.String),
        sql.Named("adminId", adminId),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *ClubService) UpdateLittleKidsMembershipViaWebportal(o clubs.LittleExplorersKidsMembership, adminId int64) error {
    query := `
        UPDATE KIDS_CLUB_MEMBERSHIP SET
          KIDS_NAME = :kidsName,
          KIDS_DOB = TO_DATE(:kidsDob, 'DD/MM/YYYY'),
          KIDS_DOC_TYPE = :kidsDocType,
          KIDS_DOC_NUMBER = :kidsDocNumber,
          KIDS_GENDER = :kidsGender,
          KIDS_NATIONALITY = :kidsNationality,
          KIDS_EMAIL = :kidsEmail,
          GUARDIAN_NAME = :guardianName,
          GUARDIAN_DOB = TO_DATE(:guardianDob, 'DD/MM/YYYY'),
          GUARDIAN_DOC_TYPE = :guardianDocType,
          GUARDIAN_DOC_NUMBER = :guardianDocNumber,
          GUARDIAN_GENDER = :guardianGender,
          GUARDIAN_NATIONALITY = :guardianNationality,
          GUARDIAN_EMAIL = :guardianEmail,
          GUARDIAN_HOME_CONTACT = :guardianHomeContact,
          GUARDIAN_MOBILE_CONTACT = :guardianMobileContact,
          GUARDIAN_ADDRESS1 = :guardianAddress1,
          GUARDIAN_POSTCODE = :guardianPostCode,
          GUARDIAN_STATE = :guardianState,
          GUARDIAN_COUNTRY_CODE = :guardianCountryCode,
          RELATIONSHIP = :relationship,
          PREFERRED_LANGUAGE = :preferredLanguage,
          USER_UPDATE = :adminId,
          DATE_UPDATE = CURRENT_TIMESTAMP
        WHERE KIDS_MEMBERSHIP_ID = :kids_membership_id
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("kidsName", o.KidsName.String),
        sql.Named("kidsDob", o.KidsDob.String),
        sql.Named("kidsDocType", o.KidsDocType.String),
        sql.Named("kidsDocNumber", o.KidsDocNumber.String),
        sql.Named("kidsGender", o.KidsGender.String),
        sql.Named("kidsNationality", o.KidsNationality.String),
        sql.Named("kidsEmail", o.KidsEmail.String),
        sql.Named("guardianName", o.GuardianName.String),
        sql.Named("guardianDob", o.GuardianDob.String),
        sql.Named("guardianDocType", o.GuardianDocType.String),
        sql.Named("guardianDocNumber", o.GuardianDocNumber.String),
        sql.Named("guardianGender", o.GuardianGender.String),
        sql.Named("guardianNationality", o.GuardianNationality.String),
        sql.Named("guardianEmail", o.GuardianEmail.String),
        sql.Named("guardianHomeContact", o.GuardianHomeContact.String),
        sql.Named("guardianMobileContact", o.GuardianMobileContact.String),
        sql.Named("guardianAddress1", o.GuardianAddress1.String),
        sql.Named("guardianPostCode", o.GuardianPostCode.String),
        sql.Named("guardianState", o.GuardianState.String),
        sql.Named("guardianCountryCode", o.GuardianCountryCode.String),
        sql.Named("relationship", o.Relationship.String),
        sql.Named("preferredLanguage", o.PreferredLanguage.String),
        sql.Named("adminId", adminId),
        sql.Named("kids_membership_id", o.KidsMembershipID.Int64),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *ClubService) SaveGoldenPearlMembership(o clubs.GoldenPearlMembership) error {
    type Prns struct {
        goldenPrn string
        nokPrn    string
    }
    const query = `
        SELECT PRN FROM NOVA_PATIENT_DOCUMENT
        WHERE DOCUMENT_TYPE = 'NRIC / Passport' AND
        DOCUMENT_NUMBER = :doc
    `
    var (
        prns Prns
        prn  string
    )
    err := s.dbrs.GetContext(s.ctxrs, &prn, query, o.GoldenDocNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            prn = ""
        } else {
            utils.LogError(err)
            return err
        }
    }
    if prn != "" {
        prns.goldenPrn = prn
    }

    err = s.dbrs.GetContext(s.ctxrs, &prn, query, o.NokDocNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            prn = ""
        } else {
            utils.LogError(err)
            return err
        }
    }
    if prn != "" {
        prns.nokPrn = prn
    }

    q := `
        INSERT INTO GOLDEN_CLUB_MEMBERSHIP
        (GOLDEN_MEMBERSHIP_NUMBER, GOLDEN_PRN, GOLDEN_NAME, GOLDEN_DOB,
        GOLDEN_DOC_TYPE, GOLDEN_DOC_NUMBER, GOLDEN_GENDER, GOLDEN_NATIONALITY,
        GOLDEN_EMAIL, NOK_PRN, NOK_NAME, NOK_DOB, NOK_DOC_TYPE,
        NOK_DOC_NUMBER, NOK_GENDER, NOK_NATIONALITY, NOK_EMAIL,
        NOK_HOME_CONTACT, NOK_MOBILE_CONTACT, NOK_ADDRESS1, NOK_ADDRESS2,
        NOK_ADDRESS3, NOK_POSTCODE, NOK_STATE, NOK_COUNTRY_CODE,
        RELATIONSHIP, PREFERRED_LANGUAGE
        ) VALUES
        (:goldenMembershipNumber, :goldenPrn, :goldenName, TO_DATE(:goldenDob, 'DD/MM/YYYY'),
        :goldenDocType, :goldenDocNumber, :goldenGender, :goldenNationality,
        :goldenEmail, :nokPrn, :nokName, TO_DATE(:nokDob, 'DD/MM/YYYY'), :nokDocType,
        :nokDocNumber, :nokGender, :nokNationality, :nokEmail,
        :nokHomeContact, :nokMobileContact, :nokAddress1, :nokAddress2,
        :nokAddress3, :nokPostCode, :nokState, :nokCountryCode,
        :relationship, :preferredLanguage)
    `
    _, err = s.db.ExecContext(s.ctx, q,
        sql.Named("goldenMembershipNumber", o.GoldenMembershipNumber.String),
        sql.Named("goldenPrn", prns.goldenPrn),
        sql.Named("goldenName", o.GoldenName.String),
        sql.Named("goldenDob", o.GoldenDob.String),
        sql.Named("goldenDocType", o.GoldenDocType.String),
        sql.Named("goldenDocNumber", o.GoldenDocNumber.String),
        sql.Named("goldenGender", o.GoldenGender.String),
        sql.Named("goldenNationality", o.GoldenNationality.String),
        sql.Named("goldenEmail", o.GoldenEmail.String),
        sql.Named("nokPrn", prns.nokPrn),
        sql.Named("nokName", o.NokName.String),
        sql.Named("nokDob", o.NokDob.String),
        sql.Named("nokDocType", o.NokDocType.String),
        sql.Named("nokDocNumber", o.NokDocNumber.String),
        sql.Named("nokGender", o.NokGender.String),
        sql.Named("nokNationality", o.NokNationality.String),
        sql.Named("nokEmail", o.NokEmail.String),
        sql.Named("nokHomeContact", o.NokHomeContact.String),
        sql.Named("nokMobileContact", o.NokMobileContact.String),
        sql.Named("nokAddress1", o.NokAddress1.String),
        sql.Named("nokAddress2", o.NokAddress2.String),
        sql.Named("nokAddress3", o.NokAddress3.String),
        sql.Named("nokPostCode", o.NokPostCode.String),
        sql.Named("nokState", o.NokState.String),
        sql.Named("nokCountryCode", o.NokCountryCode.String),
        sql.Named("relationship", o.Relationship.String),
        sql.Named("preferredLanguage", o.PreferredLanguage.String),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *ClubService) SaveGoldenPearlMembershipViaWebportal(o clubs.GoldenPearlMembership, adminId int64) error {
    type Prns struct {
        goldenPrn string
        nokPrn    string
    }
    const query = `
        SELECT PRN FROM NOVA_PATIENT_DOCUMENT
        WHERE DOCUMENT_TYPE = 'NRIC / Passport' AND
        DOCUMENT_NUMBER = :doc
    `
    var (
        prns Prns
        prn  string
    )
    err := s.dbrs.GetContext(s.ctxrs, &prn, query, o.GoldenDocNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            prn = ""
        } else {
            utils.LogError(err)
            return err
        }
    }
    if prn != "" {
        prns.goldenPrn = prn
    }

    err = s.dbrs.GetContext(s.ctxrs, &prn, query, o.NokDocNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            prn = ""
        } else {
            utils.LogError(err)
            return err
        }
    }
    if prn != "" {
        prns.nokPrn = prn
    }

    q := `
        INSERT INTO GOLDEN_CLUB_MEMBERSHIP
        (GOLDEN_MEMBERSHIP_NUMBER, GOLDEN_PRN, GOLDEN_NAME, GOLDEN_DOB,
        GOLDEN_DOC_TYPE, GOLDEN_DOC_NUMBER, GOLDEN_GENDER, GOLDEN_NATIONALITY,
        GOLDEN_EMAIL, NOK_PRN, NOK_NAME, NOK_DOB, NOK_DOC_TYPE,
        NOK_DOC_NUMBER, NOK_GENDER, NOK_NATIONALITY, NOK_EMAIL,
        NOK_HOME_CONTACT, NOK_MOBILE_CONTACT, NOK_ADDRESS1, NOK_ADDRESS2,
        NOK_ADDRESS3, NOK_POSTCODE, NOK_STATE, NOK_COUNTRY_CODE,
        RELATIONSHIP, PREFERRED_LANGUAGE, USER_CREATE
        ) VALUES
        (:goldenMembershipNumber, :goldenPrn, :goldenName, TO_DATE(:goldenDob, 'DD/MM/YYYY'),
        :goldenDocType, :goldenDocNumber, :goldenGender, :goldenNationality,
        :goldenEmail, :nokPrn, :nokName, TO_DATE(:nokDob, 'DD/MM/YYYY'), :nokDocType,
        :nokDocNumber, :nokGender, :nokNationality, :nokEmail,
        :nokHomeContact, :nokMobileContact, :nokAddress1, :nokAddress2,
        :nokAddress3, :nokPostCode, :nokState, :nokCountryCode,
        :relationship, :preferredLanguage, :adminId)
    `
    _, err = s.db.ExecContext(s.ctx, q,
        sql.Named("goldenMembershipNumber", o.GoldenMembershipNumber.String),
        sql.Named("goldenPrn", prns.goldenPrn),
        sql.Named("goldenName", o.GoldenName.String),
        sql.Named("goldenDob", o.GoldenDob.String),
        sql.Named("goldenDocType", o.GoldenDocType.String),
        sql.Named("goldenDocNumber", o.GoldenDocNumber.String),
        sql.Named("goldenGender", o.GoldenGender.String),
        sql.Named("goldenNationality", o.GoldenNationality.String),
        sql.Named("goldenEmail", o.GoldenEmail.String),
        sql.Named("nokPrn", prns.nokPrn),
        sql.Named("nokName", o.NokName.String),
        sql.Named("nokDob", o.NokDob.String),
        sql.Named("nokDocType", o.NokDocType.String),
        sql.Named("nokDocNumber", o.NokDocNumber.String),
        sql.Named("nokGender", o.NokGender.String),
        sql.Named("nokNationality", o.NokNationality.String),
        sql.Named("nokEmail", o.NokEmail.String),
        sql.Named("nokHomeContact", o.NokHomeContact.String),
        sql.Named("nokMobileContact", o.NokMobileContact.String),
        sql.Named("nokAddress1", o.NokAddress1.String),
        sql.Named("nokAddress2", o.NokAddress2.String),
        sql.Named("nokAddress3", o.NokAddress3.String),
        sql.Named("nokPostCode", o.NokPostCode.String),
        sql.Named("nokState", o.NokState.String),
        sql.Named("nokCountryCode", o.NokCountryCode.String),
        sql.Named("relationship", o.Relationship.String),
        sql.Named("preferredLanguage", o.PreferredLanguage.String),
        sql.Named("adminId", adminId),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *ClubService) UpdateGoldenPearlMembershipViaWebportal(o clubs.GoldenPearlMembership, adminId int64) error {
    query := `
        UPDATE GOLDEN_CLUB_MEMBERSHIP SET
          GOLDEN_NAME = :goldenName,
          GOLDEN_DOB = TO_DATE(:goldenDob, 'DD/MM/YYYY'),
          GOLDEN_DOC_TYPE = :goldenDocType,
          GOLDEN_DOC_NUMBER = :goldenDocNumber,
          GOLDEN_GENDER = :goldenGender,
          GOLDEN_NATIONALITY = :goldenNationality,
          GOLDEN_EMAIL = :goldenEmail,
          NOK_NAME = :nokName,
          NOK_DOB = TO_DATE(:nokDob, 'DD/MM/YYYY'),
          NOK_DOC_TYPE = :nokDocType,
          NOK_DOC_NUMBER = :nokDocNumber,
          NOK_GENDER = :nokGender,
          NOK_NATIONALITY = :nokNationality,
          NOK_EMAIL = :nokEmail,
          NOK_HOME_CONTACT = :nokHomeContact,
          NOK_MOBILE_CONTACT = :nokMobileContact,
          NOK_ADDRESS1 = :nokAddress1,
          NOK_POSTCODE = :nokPostCode,
          NOK_STATE = :nokState,
          NOK_COUNTRY_CODE = :nokCountryCode,
          RELATIONSHIP = :relationship,
          PREFERRED_LANGUAGE = :preferredLanguage,
          USER_UPDATE = :adminId,
          DATE_UPDATE = CURRENT_TIMESTAMP
        WHERE GOLDEN_MEMBERSHIP_ID = :golden_membership_id
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("goldenName", o.GoldenName.String),
        sql.Named("goldenDob", o.GoldenDob.String),
        sql.Named("goldenDocType", o.GoldenDocType.String),
        sql.Named("goldenDocNumber", o.GoldenDocNumber.String),
        sql.Named("goldenGender", o.GoldenGender.String),
        sql.Named("goldenNationality", o.GoldenNationality.String),
        sql.Named("goldenEmail", o.GoldenEmail.String),
        sql.Named("nokName", o.NokName.String),
        sql.Named("nokDob", o.NokDob.String),
        sql.Named("nokDocType", o.NokDocType.String),
        sql.Named("nokDocNumber", o.NokDocNumber.String),
        sql.Named("nokGender", o.NokGender.String),
        sql.Named("nokNationality", o.NokNationality.String),
        sql.Named("nokEmail", o.NokEmail.String),
        sql.Named("nokHomeContact", o.NokHomeContact.String),
        sql.Named("nokMobileContact", o.NokMobileContact.String),
        sql.Named("nokAddress1", o.NokAddress1.String),
        sql.Named("nokPostCode", o.NokPostCode.String),
        sql.Named("nokState", o.NokState.String),
        sql.Named("nokCountryCode", o.NokCountryCode.String),
        sql.Named("relationship", o.Relationship.String),
        sql.Named("preferredLanguage", o.PreferredLanguage.String),
        sql.Named("adminId", adminId),
        sql.Named("golden_membership_id", o.GoldenMembershipID.Int64),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func (s *ClubService) ExistsLittleKidsByDocTypeDocNo(docType string, docNo string) (bool, error) {
    query := `
        SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS 
         (SELECT 1 FROM KIDS_CLUB_MEMBERSHIP WHERE KIDS_DOC_TYPE = :docType AND KIDS_DOC_NUMBER = :docNo)
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, docType, docNo)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, nil
}

func (s *ClubService) ExistsGoldenPearlByDocTypeDocNo(docType string, docNo string) (bool, error) {
    query := `
        SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS 
         (SELECT 1 FROM GOLDEN_CLUB_MEMBERSHIP WHERE GOLDEN_DOC_TYPE = :docType AND GOLDEN_DOC_NUMBER = :docNo)
    `
    var count int
    err := s.db.GetContext(s.ctx, &count, query, docType, docNo)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, nil
}

func (s *ClubService) ExistsLittleKidsAboutUs() (bool, error) {
    query := `SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM KIDS_CLUB_INFO)`
    var count int
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, nil
}

func (s *ClubService) FindLittleKidsActivityNameById(activityId int64) (*clubs.LittleExplorersKidsActivity, error) {
    query := `SELECT KIDS_ACTIVITY_NAME FROM KIDS_CLUB_ACTIVITY WHERE KIDS_ACTIVITY_ID = :activityId`
    var o clubs.LittleExplorersKidsActivity
    err := s.db.GetContext(s.ctx, &o, query, activityId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *ClubService) FindGoldenPearlActivityNameById(activityId int64) (*clubs.GoldenPearlActivity, error) {
    query := `SELECT GOLDEN_ACTIVITY_NAME FROM GOLDEN_CLUB_ACTIVITY WHERE GOLDEN_ACTIVITY_ID = :activityId`
    var o clubs.GoldenPearlActivity
    err := s.db.GetContext(s.ctx, &o, query, activityId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *ClubService) FindLittleKidsActivitiesByActivityId(activityId int64) (*clubs.LittleExplorersKidsActivity, error) {
    query := `SELECT * FROM KIDS_CLUB_ACTIVITY WHERE KIDS_ACTIVITY_ID = :activityId`
    var o clubs.LittleExplorersKidsActivity
    err := s.db.GetContext(s.ctx, &o, query, activityId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    if o.ActivityStartDateTime.Valid {
        g, _ := goment.New(o.ActivityStartDateTime.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.ActivityStartDateTime.String = g.Format("DD/MM/YYYY HH:mm")
    }
    if o.ActivityEndDateTime.Valid && o.ActivityEndDateTime.String != "-" {
        g, _ := goment.New(o.ActivityEndDateTime.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.ActivityEndDateTime.String = g.Format("DD/MM/YYYY HH:mm")
    }
    return &o, nil
}

func (s *ClubService) FindGoldenPearlActivitiesByActivityId(activityId int64) (*clubs.GoldenPearlActivity, error) {
    query := `SELECT * FROM GOLDEN_CLUB_ACTIVITY WHERE GOLDEN_ACTIVITY_ID = :activityId`
    var o clubs.GoldenPearlActivity
    err := s.db.GetContext(s.ctx, &o, query, activityId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    if o.ActivityStartDateTime.Valid {
        g, _ := goment.New(o.ActivityStartDateTime.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.ActivityStartDateTime.String = g.Format("DD/MM/YYYY HH:mm")
    }
    if o.ActivityEndDateTime.Valid && o.ActivityEndDateTime.String != "-" {
        g, _ := goment.New(o.ActivityEndDateTime.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.ActivityEndDateTime.String = g.Format("DD/MM/YYYY HH:mm")
    }
    return &o, nil
}

func (s *ClubService) ParticipateLittleKidsActivity(o []clubs.LittleExplorersKidsActvParticipation) error {
    const query = `
        INSERT INTO KIDS_CLUB_ACTV_PARTICIPATION
        (KIDS_ACTIVITY_ID, KIDS_MEMBERSHIP_ID, ACTIVITY_DATE_TIME)
        VALUES 
        (:kidsActivityId, :kidsMembershipId, TO_DATE(:activityDateTime, 'DD/MM/YYYY hh24:mi'))
    `
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

    for i := range o {
        ap := o[i]
        _, err = tx.ExecContext(s.ctx, query,
            sql.Named("kidsActivityId", ap.KidsActivityID),
            sql.Named("kidsMembershipId", ap.KidsMembershipID),
            sql.Named("activityDateTime", ap.ActivityDateTime),
        )
        if err != nil {
            return err
        }
    }
    return tx.Commit()
}

func (s *ClubService) ParticipateGoldenPearlActivity(o []clubs.GoldenPearlActvParticipation) error {
    const query = `
        INSERT INTO GOLDEN_CLUB_ACTV_PARTICIPATION
        (GOLDEN_ACTIVITY_ID, GOLDEN_MEMBERSHIP_ID, ACTIVITY_DATE_TIME)
        VALUES
        (:goldenActivityId, :goldenMembershipId, TO_DATE(:activityDateTime, 'DD/MM/YYYY hh24:mi'))
    `
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

    for i := range o {
        ap := o[i]
        _, err = tx.ExecContext(s.ctx, query,
            sql.Named("goldenActivityId", ap.GoldenActivityID),
            sql.Named("goldenMembershipId", ap.GoldenMembershipID),
            sql.Named("activityDateTime", ap.ActivityDateTime),
        )
        if err != nil {
            return err
        }
    }
    return tx.Commit()
}

func (s *ClubService) SaveLittleKidsActivity(o clubs.LittleExplorersKidsActivity, adminId int64) error {
    const query = `
        INSERT INTO KIDS_CLUB_ACTIVITY
        (KIDS_ACTIVITY_CODE, KIDS_ACTIVITY_NAME, KIDS_ACTIVITY_DESC, 
         KIDS_ACTIVITY_IMG, ACTIVITY_START_DATETIME, ACTIVITY_END_DATETIME, 
         ACTIVITY_MAX_PARTICIPANT, ACTIVITY_TNC, ACTIVITY_DISPLAY_ORDER,
         USER_CREATE)
         VALUES
        (:kidsActivityCode, :kidsActivityName, :kidsActivityDesc, 
         :kidsActivityImage, TO_DATE(:activityStartDateTime, 'DD/MM/YYYY hh24:mi'), TO_DATE(:activityEndDateTime, 'DD/MM/YYYY hh24:mi'), 
         :activityMaxParticipant, :activityTnc, :activityDisplayOrder,
         :adminId)
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("kidsActivityCode", o.KidsActivityCode.String),
        sql.Named("kidsActivityName", o.KidsActivityName.String),
        sql.Named("kidsActivityDesc", o.KidsActivityDesc.String),
        sql.Named("kidsActivityImage", o.KidsActivityImage.String),
        sql.Named("activityStartDateTime", o.ActivityStartDateTime.String),
        sql.Named("activityEndDateTime", o.ActivityEndDateTime.String),
        sql.Named("activityMaxParticipant", o.ActivityMaxParticipant.Int32),
        sql.Named("activityTnc", o.ActivityTnc.String),
        sql.Named("activityDisplayOrder", o.ActivityDisplayOrder.String),
        sql.Named("adminId", adminId),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *ClubService) SaveGoldenPearlActivity(o clubs.GoldenPearlActivity, adminId int64) error {
    const query = `
        INSERT INTO GOLDEN_CLUB_ACTIVITY
        (GOLDEN_ACTIVITY_CODE, GOLDEN_ACTIVITY_NAME, GOLDEN_ACTIVITY_DESC, 
         GOLDEN_ACTIVITY_IMG, ACTIVITY_START_DATETIME, ACTIVITY_END_DATETIME, 
         ACTIVITY_MAX_PARTICIPANT, ACTIVITY_TNC, ACTIVITY_DISPLAY_ORDER,
         USER_CREATE)
         VALUES
        (:goldenActivityCode, :goldenActivityName, :goldenActivityDesc, 
         :goldenActivityImage, TO_DATE(:activityStartDateTime, 'DD/MM/YYYY hh24:mi'), TO_DATE(:activityEndDateTime, 'DD/MM/YYYY hh24:mi'), 
         :activityMaxParticipant, :activityTnc, :activityDisplayOrder,
         :adminId)
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("goldenActivityCode", o.GoldenActivityCode.String),
        sql.Named("goldenActivityName", o.GoldenActivityName.String),
        sql.Named("goldenActivityDesc", o.GoldenActivityDesc.String),
        sql.Named("goldenActivityImage", o.GoldenActivityImage.String),
        sql.Named("activityStartDateTime", o.ActivityStartDateTime.String),
        sql.Named("activityEndDateTime", o.ActivityEndDateTime.String),
        sql.Named("activityMaxParticipant", o.ActivityMaxParticipant.Int32),
        sql.Named("activityTnc", o.ActivityTnc.String),
        sql.Named("activityDisplayOrder", o.ActivityDisplayOrder.String),
        sql.Named("adminId", adminId),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *ClubService) UpdateLittleKidsActivity(o clubs.LittleExplorersKidsActivity, adminId int64) error {
    const query = `
        UPDATE KIDS_CLUB_ACTIVITY SET
          KIDS_ACTIVITY_CODE = :kidsActivityCode,
          KIDS_ACTIVITY_NAME = :kidsActivityName,
          KIDS_ACTIVITY_DESC = :kidsActivityDesc,
          KIDS_ACTIVITY_IMG = :kidsActivityImage,
          ACTIVITY_START_DATETIME = TO_DATE(:activityStartDateTime, 'DD/MM/YYYY hh24:mi'),
          ACTIVITY_END_DATETIME = TO_DATE(:activityEndDateTime, 'DD/MM/YYYY hh24:mi'),
          ACTIVITY_MAX_PARTICIPANT = :activityMaxParticipant,
          ACTIVITY_TNC = :activityTnc,
          ACTIVITY_DISPLAY_ORDER = :activityDisplayOrder,
          USER_UPDATE = :adminId,
          DATE_UPDATE = CURRENT_TIMESTAMP
        WHERE KIDS_ACTIVITY_ID = :kids_activity_id
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("kidsActivityCode", o.KidsActivityCode.String),
        sql.Named("kidsActivityName", o.KidsActivityName.String),
        sql.Named("kidsActivityDesc", o.KidsActivityDesc.String),
        sql.Named("kidsActivityImage", o.KidsActivityImage.String),
        sql.Named("activityStartDateTime", o.ActivityStartDateTime.String),
        sql.Named("activityEndDateTime", o.ActivityEndDateTime.String),
        sql.Named("activityMaxParticipant", o.ActivityMaxParticipant.Int32),
        sql.Named("activityTnc", o.ActivityTnc.String),
        sql.Named("activityDisplayOrder", o.ActivityDisplayOrder.String),
        sql.Named("adminId", adminId),
        sql.Named("kids_activity_id", o.KidsActivityID.Int64),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *ClubService) UpdateGoldenPearlActivity(o clubs.GoldenPearlActivity, adminId int64) error {
    const query = `
        UPDATE GOLDEN_CLUB_ACTIVITY SET
          GOLDEN_ACTIVITY_CODE = :goldenActivityCode,
          GOLDEN_ACTIVITY_NAME = :goldenActivityName,
          GOLDEN_ACTIVITY_DESC = :goldenActivityDesc,
          GOLDEN_ACTIVITY_IMG = :goldenActivityImage,
          ACTIVITY_START_DATETIME = TO_DATE(:activityStartDateTime, 'DD/MM/YYYY hh24:mi'),
          ACTIVITY_END_DATETIME = TO_DATE(:activityEndDateTime, 'DD/MM/YYYY hh24:mi'),
          ACTIVITY_MAX_PARTICIPANT = :activityMaxParticipant,
          ACTIVITY_TNC = :activityTnc,
          ACTIVITY_DISPLAY_ORDER = :activityDisplayOrder,
          USER_UPDATE = :adminId,
          DATE_UPDATE = CURRENT_TIMESTAMP
        WHERE GOLDEN_ACTIVITY_ID = :golden_activity_id
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("goldenActivityCode", o.GoldenActivityCode.String),
        sql.Named("goldenActivityName", o.GoldenActivityName.String),
        sql.Named("goldenActivityDesc", o.GoldenActivityDesc.String),
        sql.Named("goldenActivityImage", o.GoldenActivityImage.String),
        sql.Named("activityStartDateTime", o.ActivityStartDateTime.String),
        sql.Named("activityEndDateTime", o.ActivityEndDateTime.String),
        sql.Named("activityMaxParticipant", o.ActivityMaxParticipant.Int32),
        sql.Named("activityTnc", o.ActivityTnc.String),
        sql.Named("activityDisplayOrder", o.ActivityDisplayOrder.String),
        sql.Named("adminId", adminId),
        sql.Named("golden_activity_id", o.GoldenActivityID.Int64),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *ClubService) SaveLittleKidsAboutUs(o clubs.LittleExplorersKidsAboutUs, adminId int64) (int64, error) {
    const query = `
        INSERT INTO KIDS_CLUB_INFO
         (KIDS_CLUB_TITLE, KIDS_CLUB_DESC, KIDS_CLUB_IMG, 
          KIDS_CLUB_TNC, PARTNERS_LINK, USER_CREATE)
         VALUES 
         (:kidsClubTitle, :kidsClubDesc, :kidsClubImage, 
          :kidsClubTnc, :kidsClubPartnerLink, :adminId
        ) RETURNING KIDS_CLUB_ID INTO :kids_club_id
    `
    var kidsClubId go_ora.Number
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("kidsClubTitle", o.KidsClubTitle.String),
        sql.Named("kidsClubDesc", o.KidsClubDesc.String),
        sql.Named("kidsClubImage", o.KidsClubImage.String),
        sql.Named("kidsClubTnc", o.KidsClubTnc.String),
        sql.Named("kidsClubPartnerLink", o.KidsClubPartnerLink.String),
        sql.Named("adminId", adminId),
        go_ora.Out{Dest: &kidsClubId},
    )
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    i, _ := kidsClubId.Int64()
    return i, nil
}

func (s *ClubService) SaveGoldenPearlAboutUs(o clubs.GoldenPearlAboutUs, adminId int64) (int64, error) {
    const query = `
        INSERT INTO GOLDEN_CLUB_INFO
         (GOLDEN_CLUB_TITLE, GOLDEN_CLUB_DESC, GOLDEN_CLUB_IMG, 
          GOLDEN_CLUB_TNC, EXTERNAL_LINK, USER_CREATE)
         VALUES 
         (:goldenClubTitle, :goldenClubDesc, :goldenClubImage, 
          :goldenClubTnc, :goldenClubExtLink, :adminId
        ) RETURNING GOLDEN_CLUB_ID INTO :golden_club_id
    `
    var goldenClubId go_ora.Number
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("goldenClubTitle", o.GoldenClubTitle.String),
        sql.Named("goldenClubDesc", o.GoldenClubDesc.String),
        sql.Named("goldenClubImage", o.GoldenClubImage.String),
        sql.Named("goldenClubTnc", o.GoldenClubTnc.String),
        sql.Named("goldenClubExtLink", o.GoldenClubExtLink.String),
        sql.Named("adminId", adminId),
        go_ora.Out{Dest: &goldenClubId},
    )
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    i, _ := goldenClubId.Int64()
    return i, nil
}

func (s *ClubService) UpdateLittleKidsAboutUs(o clubs.LittleExplorersKidsAboutUs, adminId int64) error {
    const query = `
        UPDATE KIDS_CLUB_INFO SET
          KIDS_CLUB_TITLE = :kidsClubTitle,
          KIDS_CLUB_DESC = :kidsClubDesc,
          KIDS_CLUB_IMG = :kidsClubImage,
          KIDS_CLUB_TNC = :kidsClubTnc,
          PARTNERS_LINK = :kidsClubPartnerLink,
          USER_UPDATE = :adminId,
          DATE_UPDATE = CURRENT_TIMESTAMP
        WHERE KIDS_CLUB_ID = :kids_club_id
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("kidsClubTitle", o.KidsClubTitle.String),
        sql.Named("kidsClubDesc", o.KidsClubDesc.String),
        sql.Named("kidsClubImage", o.KidsClubImage.String),
        sql.Named("kidsClubTnc", o.KidsClubTnc.String),
        sql.Named("kidsClubPartnerLink", o.KidsClubPartnerLink.String),
        sql.Named("adminId", adminId),
        sql.Named("kids_club_id", o.KidsClubID.Int64),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *ClubService) UpdateGoldenPearlAboutUs(o clubs.GoldenPearlAboutUs, adminId int64) error {
    const query = `
        UPDATE GOLDEN_CLUB_INFO SET
          GOLDEN_CLUB_TITLE = :goldenClubTitle,
          GOLDEN_CLUB_DESC = :goldenClubDesc,
          GOLDEN_CLUB_IMG = :goldenClubImage,
          GOLDEN_CLUB_TNC = :goldenClubTnc,
          EXTERNAL_LINK = :goldenClubExtLink,
          USER_UPDATE = :adminId,
          DATE_UPDATE = CURRENT_TIMESTAMP
         WHERE GOLDEN_CLUB_ID = :golden_club_id
    `
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("goldenClubTitle", o.GoldenClubTitle.String),
        sql.Named("goldenClubDesc", o.GoldenClubDesc.String),
        sql.Named("goldenClubImage", o.GoldenClubImage.String),
        sql.Named("goldenClubTnc", o.GoldenClubTnc.String),
        sql.Named("goldenClubExtLink", o.GoldenClubExtLink.String),
        sql.Named("adminId", adminId),
        sql.Named("golden_club_id", o.GoldenClubID.Int64),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *ClubService) FindLittleKidsAboutUs() (*clubs.LittleExplorersKidsAboutUs, error) {
    query := `SELECT ` + utils.GetDbCols(clubs.LittleExplorersKidsAboutUs{}, "") + ` FROM KIDS_CLUB_INFO`
    var o clubs.LittleExplorersKidsAboutUs
    err := s.db.GetContext(s.ctx, &o, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *ClubService) FindGoldenPearlAboutUs() (*clubs.GoldenPearlAboutUs, error) {
    query := `SELECT ` + utils.GetDbCols(clubs.GoldenPearlAboutUs{}, "") + ` FROM GOLDEN_CLUB_INFO`
    var o clubs.GoldenPearlAboutUs
    err := s.db.GetContext(s.ctx, &o, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *ClubService) FindAllAppLittleKidsActivities(offset int, limit int, isHome bool) ([]clubs.LittleExplorersKidsActivity, error) {
    query := `
        SELECT kca.*, (SELECT COUNT(*) 
        FROM KIDS_CLUB_ACTV_PARTICIPATION kcap
        WHERE kcap.KIDS_ACTIVITY_ID = kca.KIDS_ACTIVITY_ID) AS ATTENDEES
        FROM KIDS_CLUB_ACTIVITY kca
        WHERE (
            ACTIVITY_START_DATETIME <= CURRENT_TIMESTAMP AND (
            ACTIVITY_END_DATETIME >= CURRENT_TIMESTAMP OR
            ACTIVITY_END_DATETIME IS NULL)
        )
        ORDER BY ACTIVITY_DISPLAY_ORDER, ACTIVITY_START_DATETIME
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
        
    `
    if isHome {
        query = `
            SELECT kca.*, (SELECT COUNT(*) 
            FROM KIDS_CLUB_ACTV_PARTICIPATION kcap
            WHERE kcap.KIDS_ACTIVITY_ID = kca.KIDS_ACTIVITY_ID) AS ATTENDEES
            FROM KIDS_CLUB_ACTIVITY kca
            WHERE (
            ACTIVITY_START_DATETIME <= CURRENT_TIMESTAMP AND (
            ACTIVITY_END_DATETIME >= CURRENT_TIMESTAMP OR
            ACTIVITY_END_DATETIME IS NULL)
            ) AND ACTIVITY_DISPLAY_ORDER = 1
            ORDER BY ACTIVITY_START_DATETIME
            FETCH FIRST 5 ROWS ONLY
        `
    }
    list := make([]clubs.LittleExplorersKidsActivity, 0)
    var err error
    if isHome {
        err = s.db.SelectContext(s.ctx, &list, query)
    } else {
        err = s.db.SelectContext(s.ctx, &list, query, offset, limit)
    }

    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    for i := range list {
        list[i].Set()
        list[i].ActivitySeatsAvailable = int(list[i].ActivityMaxParticipant.Int32) - int(list[i].ActivityAttendees.Int32)

        if list[i].ActivityEndDateTime.String == "-" {
            futureDate, _ := goment.New()
            futureDate.Add(10, "years")
            list[i].ActivityEndDateTimeCalendar = futureDate.ToISOString()
        } else {
            list[i].ActivityEndDateTimeCalendar = ""
        }
    }
    return list, nil
}

func (s *ClubService) FindAllAppGoldenPearlActivities(offset int, limit int, isHome bool) ([]clubs.GoldenPearlActivity, error) {
    query := `
        SELECT gca.*, (SELECT COUNT(*) 
        FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
        WHERE gcap.GOLDEN_ACTIVITY_ID = gca.GOLDEN_ACTIVITY_ID) AS ATTENDEES
        FROM GOLDEN_CLUB_ACTIVITY gca
        WHERE (
            ACTIVITY_START_DATETIME <= CURRENT_TIMESTAMP AND (
            ACTIVITY_END_DATETIME >= CURRENT_TIMESTAMP OR
            ACTIVITY_END_DATETIME IS NULL)
        )
        ORDER BY ACTIVITY_DISPLAY_ORDER, ACTIVITY_START_DATETIME
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    if isHome {
        query = `
            SELECT gca.*, (SELECT COUNT(*) 
            FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
            WHERE gcap.GOLDEN_ACTIVITY_ID = gca.GOLDEN_ACTIVITY_ID) AS ATTENDEES
            FROM GOLDEN_CLUB_ACTIVITY gca
            WHERE (
            ACTIVITY_START_DATETIME <= CURRENT_TIMESTAMP AND (
            ACTIVITY_END_DATETIME >= CURRENT_TIMESTAMP OR
            ACTIVITY_END_DATETIME IS NULL)
            ) AND ACTIVITY_DISPLAY_ORDER = 1
            ORDER BY ACTIVITY_START_DATETIME
            FETCH FIRST 5 ROWS ONLY
        `
    }
    list := make([]clubs.GoldenPearlActivity, 0)
    var err error
    if isHome {
        err = s.db.SelectContext(s.ctx, &list, query)
    } else {
        err = s.db.SelectContext(s.ctx, &list, query, offset, limit)
    }
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    for i := range list {
        list[i].Set()
        list[i].ActivitySeatsAvailable = int(list[i].ActivityMaxParticipant.Int32) - int(list[i].ActivityAttendees.Int32)

        if list[i].ActivityEndDateTime.String == "-" {
            futureDate, _ := goment.New()
            futureDate.Add(10, "years")
            list[i].ActivityEndDateTimeCalendar = futureDate.ToISOString()
        } else {
            list[i].ActivityEndDateTimeCalendar = ""
        }
    }
    return list, nil
}

func (s *ClubService) FindAllUserLittleKidsActivities(userId int64) ([]clubs.LittleExplorersKidsMyActivity, error) {
    patient, err := s.applicationUserService.FindByUserId(userId, nil)
    if err != nil {
        return nil, err
    }

    query := `
        SELECT KCM.KIDS_NAME, KCM.KIDS_MEMBERSHIP_NUMBER, KCA.KIDS_ACTIVITY_NAME, KCP.ACTIVITY_DATE_TIME
         FROM KIDS_CLUB_ACTV_PARTICIPATION KCP
         INNER JOIN KIDS_CLUB_ACTIVITY KCA ON KCP.KIDS_ACTIVITY_ID = KCA.KIDS_ACTIVITY_ID
         INNER JOIN KIDS_CLUB_MEMBERSHIP KCM ON KCP.KIDS_MEMBERSHIP_ID = KCM.KIDS_MEMBERSHIP_ID
         WHERE KCP.KIDS_MEMBERSHIP_ID IN (
          SELECT KIDS_MEMBERSHIP_ID
          FROM KIDS_CLUB_MEMBERSHIP
          WHERE KIDS_PRN = :prn
          UNION
          SELECT KIDS_MEMBERSHIP_ID
          FROM KIDS_CLUB_MEMBERSHIP
          WHERE GUARDIAN_PRN = :prn
          UNION
          SELECT KIDS_MEMBERSHIP_ID
          FROM KIDS_CLUB_MEMBERSHIP
          WHERE KIDS_PRN IN (
            SELECT NOK_PRN
            FROM APPLICATION_USER_FAMILY
            WHERE PATIENT_PRN = :prn
            AND IS_PATIENT = 'Y'
          )
         )
    `
    list := make([]clubs.LittleExplorersKidsMyActivity, 0)
    err = s.db.SelectContext(s.ctx, &list, query, sql.Named("prn", patient.MasterPrn.String))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *ClubService) FindAllUserGoldenPearlActivities(userId int64) ([]clubs.GoldenPearlMyActivity, error) {
    patient, err := s.applicationUserService.FindByUserId(userId, nil)
    if err != nil {
        return nil, err
    }

    query := `
        SELECT GCM.GOLDEN_NAME, GCM.GOLDEN_MEMBERSHIP_NUMBER, GCA.GOLDEN_ACTIVITY_NAME, GCP.ACTIVITY_DATE_TIME
         FROM GOLDEN_CLUB_ACTV_PARTICIPATION GCP
         INNER JOIN GOLDEN_CLUB_ACTIVITY GCA ON GCP.GOLDEN_ACTIVITY_ID = GCA.GOLDEN_ACTIVITY_ID
         INNER JOIN GOLDEN_CLUB_MEMBERSHIP GCM ON GCP.GOLDEN_MEMBERSHIP_ID = GCM.GOLDEN_MEMBERSHIP_ID
         WHERE GCP.GOLDEN_MEMBERSHIP_ID IN (
          SELECT GOLDEN_MEMBERSHIP_ID
          FROM GOLDEN_CLUB_MEMBERSHIP
          WHERE GOLDEN_PRN = :prn
          UNION
          SELECT GOLDEN_MEMBERSHIP_ID
          FROM GOLDEN_CLUB_MEMBERSHIP
          WHERE NOK_PRN = :prn
          UNION
          SELECT GOLDEN_MEMBERSHIP_ID
          FROM GOLDEN_CLUB_MEMBERSHIP
          WHERE GOLDEN_PRN IN (
            SELECT NOK_PRN
            FROM APPLICATION_USER_FAMILY
            WHERE PATIENT_PRN = :prn
            AND IS_PATIENT = 'Y'
          )
         )
    `
    list := make([]clubs.GoldenPearlMyActivity, 0)
    err = s.db.SelectContext(s.ctx, &list, query, sql.Named("prn", patient.MasterPrn.String))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *ClubService) FindAllLittleKidsAttendees(kidsActivityId int64, offset int, limit int) ([]clubs.LittleExplorersKidsMembership, error) {
    const query = `
        SELECT kcm.*, kcap.ACTIVITY_DATE_TIME
        FROM KIDS_CLUB_ACTV_PARTICIPATION kcap
        JOIN KIDS_CLUB_MEMBERSHIP kcm ON kcap.KIDS_MEMBERSHIP_ID = kcm.KIDS_MEMBERSHIP_ID
        WHERE kcap.KIDS_ACTIVITY_ID = :kidsActivityId
        ORDER BY KIDS_MEMBERSHIP_NUMBER DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    list := make([]clubs.LittleExplorersKidsMembership, 0)
    err := s.db.SelectContext(s.ctx, &list, query, 
        sql.Named("kidsActivityId", kidsActivityId), 
        sql.Named("offset", offset), 
        sql.Named("limit", limit),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].SetAttendees()
    }
    return list, nil
}

func (s *ClubService) FindAllGoldenPearlAttendees(goldenActivityId int64, offset int, limit int) ([]clubs.GoldenPearlMembership, error) {
    const query = `
        SELECT gcm.*, gcap.ACTIVITY_DATE_TIME
        FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
        JOIN GOLDEN_CLUB_MEMBERSHIP gcm ON gcap.GOLDEN_MEMBERSHIP_ID = gcm.GOLDEN_MEMBERSHIP_ID
        WHERE gcap.GOLDEN_ACTIVITY_ID = :goldenActivityId
        ORDER BY GOLDEN_MEMBERSHIP_NUMBER DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    list := make([]clubs.GoldenPearlMembership, 0)
    err := s.db.SelectContext(s.ctx, &list, query, 
        sql.Named("goldenActivityId", goldenActivityId), 
        sql.Named("offset", offset), 
        sql.Named("limit", limit),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].SetAttendees()
    }
    return list, nil
}

func (s *ClubService) FindAllLittleKidsMemberships(offset int, limit int) ([]clubs.LittleExplorersKidsMembership, error) {
    query := `SELECT ` + utils.GetDbCols(clubs.LittleExplorersKidsMembership{}, "") + ` 
        FROM KIDS_CLUB_MEMBERSHIP ORDER BY KIDS_MEMBERSHIP_NUMBER DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    list := make([]clubs.LittleExplorersKidsMembership, 0)
    err := s.db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
}

func (s *ClubService) FindAllGoldenPearlMemberships(offset int, limit int) ([]clubs.GoldenPearlMembership, error) {
    query := `SELECT ` + utils.GetDbCols(clubs.GoldenPearlMembership{}, "") + ` 
        FROM GOLDEN_CLUB_MEMBERSHIP ORDER BY GOLDEN_MEMBERSHIP_NUMBER DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    list := make([]clubs.GoldenPearlMembership, 0)
    err := s.db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
}

func (s *ClubService) FindAllLittleKidsActivities(offset int, limit int) ([]clubs.LittleExplorersKidsActivity, error) {
    m := map[string]string{
        "kca.ATTENDEES": "ATTENDEES",
    }
    query := `
        SELECT ` + utils.GetDbColsWithReplace(clubs.LittleExplorersKidsActivity{}, "kca.", m) + `, (SELECT COUNT(*) 
        FROM KIDS_CLUB_ACTV_PARTICIPATION kcap
        WHERE kcap.KIDS_ACTIVITY_ID = kca.KIDS_ACTIVITY_ID) AS ATTENDEES
        FROM KIDS_CLUB_ACTIVITY kca
        ORDER BY ACTIVITY_START_DATETIME DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    list := make([]clubs.LittleExplorersKidsActivity, 0)
    err := s.db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
}

func (s *ClubService) FindAllGoldenPearlActivities(offset int, limit int) ([]clubs.GoldenPearlActivity, error) {
    m := map[string]string{
        "gca.ATTENDEES": "ATTENDEES",
    }
    query := `
        SELECT ` + utils.GetDbColsWithReplace(clubs.GoldenPearlActivity{}, "gca.", m) + `, (SELECT COUNT(*) 
        FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
        WHERE gcap.GOLDEN_ACTIVITY_ID = gca.GOLDEN_ACTIVITY_ID) AS ATTENDEES
        FROM GOLDEN_CLUB_ACTIVITY gca
        ORDER BY ACTIVITY_START_DATETIME DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    list := make([]clubs.GoldenPearlActivity, 0)
    err := s.db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
}

func (s *ClubService) ListAppLittleKidsActivities(isHome bool, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountAppLittleKidsActivities(isHome)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllAppLittleKidsActivities(pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountAppLittleKidsActivities(isHome bool) (int, error) {
    query := `
        SELECT COUNT(KIDS_ACTIVITY_ID) AS COUNT FROM KIDS_CLUB_ACTIVITY
        WHERE (
            ACTIVITY_START_DATETIME <= CURRENT_TIMESTAMP 
            AND (
                ACTIVITY_END_DATETIME >= CURRENT_TIMESTAMP 
                OR ACTIVITY_END_DATETIME IS NULL
            )
        )
    `
    if isHome {
        query = `
            SELECT COUNT(KIDS_ACTIVITY_ID) AS COUNT FROM KIDS_CLUB_ACTIVITY
             WHERE (
              ACTIVITY_START_DATETIME <= CURRENT_TIMESTAMP 
              AND (
                ACTIVITY_END_DATETIME >= CURRENT_TIMESTAMP 
                OR ACTIVITY_END_DATETIME IS NULL
              )
             ) AND ACTIVITY_DISPLAY_ORDER = 1
        `
    }
    var count int
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) ListLittleKidsActivities(page string, limit string) (*model.PagedList, error) {
    total, err := s.CountLittleKidsActivities()
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllLittleKidsActivities(pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountLittleKidsActivities() (int, error) {
    const query = `SELECT COUNT(KIDS_ACTIVITY_ID) AS COUNT FROM KIDS_CLUB_ACTIVITY`
    var count int
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) ListGoldenPearlActivities(page string, limit string) (*model.PagedList, error) {
    total, err := s.CountGoldenPearlActivities(false)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllGoldenPearlActivities(pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountGoldenPearlActivities(isHome bool) (int, error) {
    const query = `SELECT COUNT(GOLDEN_ACTIVITY_ID) AS COUNT FROM GOLDEN_CLUB_ACTIVITY`
    var count int
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) ListLittleKidsActivitiesByKeyword(keyword string, keyword2 string, keyword3 string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountLittleKidsActivitiesByKeyword(keyword, keyword2, keyword3)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindLittleKidsActivitiesByKeyword(keyword, keyword2, keyword3, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountLittleKidsActivitiesByKeyword(keyword string, keyword2 string, keyword3 string) (int, error) {
    conditions, args := buildLittleKidsActivitiesKeywordConditions(keyword, keyword2, keyword3)
    base := `SELECT COUNT(kca.KIDS_ACTIVITY_ID) AS COUNT FROM KIDS_CLUB_ACTIVITY kca`
    query := base + whereClause(conditions)

    var count int
    err := s.db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) ListGoldenPearlActivitiesByKeyword(keyword string, keyword2 string, keyword3 string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountGoldenPearlActivitiesByKeyword(keyword, keyword2, keyword3)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindGoldenPearlActivitiesByKeyword(keyword, keyword2, keyword3, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountGoldenPearlActivitiesByKeyword(keyword string, keyword2 string, keyword3 string) (int, error) {
    conditions, args := buildGoldenPearlActivitiesKeywordConditions(keyword, keyword2, keyword3)
    base := `SELECT COUNT(gca.GOLDEN_ACTIVITY_ID) AS COUNT FROM GOLDEN_CLUB_ACTIVITY gca`
    query := base + whereClause(conditions)

    var count int
    err := s.db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) FindLittleKidsActivitiesByKeyword(keyword string, keyword2 string, keyword3 string, offset int, limit int) ([]clubs.LittleExplorersKidsActivity, error) {
    conditions, args := buildLittleKidsActivitiesKeywordConditions(keyword, keyword2, keyword3)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    m := map[string]string{
        "kca.ATTENDEES": "ATTENDEES",
    }
    base := `
        SELECT ` + utils.GetDbColsWithReplace(clubs.LittleExplorersKidsActivity{}, "kca.", m) + `, (SELECT COUNT(*)
        FROM KIDS_CLUB_ACTV_PARTICIPATION kcap
        WHERE kcap.KIDS_ACTIVITY_ID = kca.KIDS_ACTIVITY_ID) AS ATTENDEES
        FROM KIDS_CLUB_ACTIVITY kca
    `
    query := base + whereClause(conditions) + 
        ` ORDER BY kca.ACTIVITY_START_DATETIME DESC 
          OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    var list []clubs.LittleExplorersKidsActivity
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *ClubService) FindGoldenPearlActivitiesByKeyword(keyword string, keyword2 string, keyword3 string, offset int, limit int) ([]clubs.GoldenPearlActivity, error) {
    conditions, args := buildGoldenPearlActivitiesKeywordConditions(keyword, keyword2, keyword3)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    m := map[string]string{
        "gca.ATTENDEES": "ATTENDEES",
    }
    base := `
        SELECT ` + utils.GetDbColsWithReplace(clubs.GoldenPearlActivity{}, "gca.", m) + `, (SELECT COUNT(*)
        FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
        WHERE gcap.GOLDEN_ACTIVITY_ID = gca.GOLDEN_ACTIVITY_ID) AS ATTENDEES
        FROM GOLDEN_CLUB_ACTIVITY gca
    `
    query := base + whereClause(conditions) + 
        ` ORDER BY gca.ACTIVITY_START_DATETIME DESC 
          OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    var list []clubs.GoldenPearlActivity
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *ClubService) GenerateKidsMembershipNo() (string, error) {
    kidsMembershipNo := ""
    const query = `SELECT GENERATE_KIDS_MEMBERSHIP() AS KIDS_MEMBERSHIP_NUMBER FROM DUAL`
    err := s.db.GetContext(s.ctx, &kidsMembershipNo, query)
    if err != nil {
        utils.LogError(err)
        return "", err
    }
    return kidsMembershipNo, nil
}

func (s *ClubService) ExistsGoldenPearlAboutUs() (bool, error) {
    query := `SELECT COUNT(*) AS COUNT FROM DUAL WHERE EXISTS (SELECT 1 FROM GOLDEN_CLUB_INFO)`
    var count int
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return false, err
    }
    return count == 1, nil
}

func (s * ClubService) GenerateGoldenMembershipNo() (string, error) {
    goldenMembershipNo := ""
    const query = `SELECT GENERATE_GOLDEN_MEMBERSHIP() AS GOLDEN_MEMBERSHIP_NUMBER FROM DUAL`
    err := s.db.GetContext(s.ctx, &goldenMembershipNo, query)
    if err != nil {
        utils.LogError(err)
        return "", err
    }
    return goldenMembershipNo, nil
}

func (s *ClubService) ExistsLittleKidsByPrn(docNumber string) (bool, error) {
    return false, nil
}



func buildLittleKidsActivitiesKeywordConditions(keyword string, keyword2 string, keyword3 string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `(LOWER(kca.KIDS_ACTIVITY_CODE) LIKE :keyword OR LOWER(kca.KIDS_ACTIVITY_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        conds = append(conds, `TRUNC(kca.ACTIVITY_START_DATETIME) = TO_DATE(:keyword2, 'dd/mm/yyyy')`)
        args = append(args, sql.Named("keyword2", keyword2))
    }
    if keyword3 != "" {
        conds = append(conds, `TRUNC(kca.ACTIVITY_END_DATETIME) = TO_DATE(:keyword3, 'dd/mm/yyyy')`)
        args = append(args, sql.Named("keyword3", keyword3))
    }
    return conds, args
}

func buildGoldenPearlActivitiesKeywordConditions(keyword string, keyword2 string, keyword3 string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `LOWER(gca.GOLDEN_ACTIVITY_CODE) LIKE :keyword OR LOWER(gca.GOLDEN_ACTIVITY_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        conds = append(conds, `TRUNC(gca.ACTIVITY_START_DATETIME) = TO_DATE(:keyword2, 'dd/mm/yyyy')`)
        args = append(args, sql.Named("keyword2", keyword2))
    }
    if keyword3 != "" {
        conds = append(conds, `TRUNC(gca.ACTIVITY_END_DATETIME) = TO_DATE(:keyword3, 'dd/mm/yyyy')`)
        args = append(args, sql.Named("keyword3", keyword3))
    }
    return conds, args
}

func whereClause(conds []string) string {
    if len(conds) == 0 {
        return ""
    }
    return " WHERE " + strings.Join(conds, " AND ")
}
