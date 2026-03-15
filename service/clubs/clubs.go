package clubs

import (
    "context"
    "database/sql"
    "vesaliusm/database"
    "vesaliusm/model/clubs"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    "github.com/nleeper/goment"
)

var ClubSvc *ClubService = NewClubService(database.GetDb(), database.GetCtx(), database.GetDbrs(), database.GetCtxrs())

type ClubService struct {
    db    *sqlx.DB
    dbrs  *sqlx.DB
    ctx   context.Context
    ctxrs context.Context
}

func NewClubService(db *sqlx.DB, ctx context.Context, dbrs *sqlx.DB, ctxrs context.Context) *ClubService {
    return &ClubService{
        db: db, 
        dbrs: dbrs, 
        ctx: ctx, 
        ctxrs: ctxrs,
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
        prn string
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
        prn string
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
    q := `
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
    _, err := s.db.ExecContext(s.ctx, q,
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
    const query := `
        SELECT PRN FROM NOVA_PATIENT_DOCUMENT
        WHERE DOCUMENT_TYPE = 'NRIC / Passport' AND
        DOCUMENT_NUMBER = Ldoc
    `
    var (
        prns Prns
        prn string
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
