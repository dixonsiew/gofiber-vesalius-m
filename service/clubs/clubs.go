package clubs

import (
    "context"
    "database/sql"
    "vesaliusm/model/clubs"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    "github.com/nleeper/goment"
)

type ClubService struct {
    db    *sqlx.DB
    dbrs  *sqlx.DB
    ctx   context.Context
    ctxrs context.Context
}

func NewClubService(db *sqlx.DB, ctx context.Context, dbrs *sqlx.DB, ctxrs context.Context) *ClubService {
    return &ClubService{db: db, dbrs: dbrs, ctx: ctx, ctxrs: ctxrs}
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
        sql.Named("kidsMembershipNumber", o.KidsMembershipNumber),
        sql.Named("kidsPrn", prns.kidsPrn),
        sql.Named("kidsName", o.KidsName),
        sql.Named("kidsDob", o.KidsDob),
        sql.Named("kidsDocType", o.KidsDocType),
        sql.Named("kidsDocNumber", o.KidsDocNumber),
        sql.Named("kidsGender", o.KidsGender),
        sql.Named("kidsNationality", o.KidsNationality),
        sql.Named("kidsEmail", o.KidsEmail),
        sql.Named("guardianPrn", prns.guardianPrn),
        sql.Named("guardianName", o.GuardianName),
        sql.Named("guardianDob", o.GuardianDob),
        sql.Named("guardianDocType", o.GuardianDocType),
        sql.Named("guardianDocNumber", o.GuardianDocNumber),
        sql.Named("guardianGender", o.GuardianGender),
        sql.Named("guardianNationality", o.GuardianNationality),
        sql.Named("guardianEmail", o.GuardianEmail),
        sql.Named("guardianHomeContact", o.GuardianHomeContact),
        sql.Named("guardianMobileContact", o.GuardianMobileContact),
        sql.Named("guardianAddress1", o.GuardianAddress1),
        sql.Named("guardianAddress2", o.GuardianAddress2),
        sql.Named("guardianAddress3", o.GuardianAddress3),
        sql.Named("guardianPostCode", o.GuardianPostCode),
        sql.Named("guardianState", o.GuardianState),
        sql.Named("guardianCountryCode", o.GuardianCountryCode),
        sql.Named("relationship", o.Relationship),
        sql.Named("preferredLanguage", o.PreferredLanguage),
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
