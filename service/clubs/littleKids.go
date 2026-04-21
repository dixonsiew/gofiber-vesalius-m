package clubs

import (
    "database/sql"
    "fmt"
    "strings"
    "vesaliusm/dto"
    "vesaliusm/model"
    "vesaliusm/model/clubs"
    "vesaliusm/utils"

    "github.com/nleeper/goment"
    go_ora "github.com/sijms/go-ora/v2"
)

func (s *ClubService) FindLittleKidsMembershipById(kidsMembershipId int64) (*clubs.LittleExplorersKidsMembership, error) {
    query := `SELECT * FROM KIDS_CLUB_MEMBERSHIP WHERE KIDS_MEMBERSHIP_ID = :kidsMembershipId`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.LittleExplorersKidsMembership{}, ""), 1)
    var o clubs.LittleExplorersKidsMembership
    err := s.db.GetContext(s.ctx, &o, query, kidsMembershipId)
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

func (s *ClubService) FindLittleKidsMembershipByMembershipId(membershipId int64) (*clubs.LittleExplorersKidsMembership, error) {
    query := `SELECT * FROM KIDS_CLUB_MEMBERSHIP WHERE KIDS_MEMBERSHIP_ID = :membershipId`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.LittleExplorersKidsMembership{}, ""), 1)
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
        o.KidsDob = utils.NewNullString(g.Format("DD/MM/YYYY"))
    }
    if o.GuardianDob.Valid {
        g, _ := goment.New(o.GuardianDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.GuardianDob = utils.NewNullString(g.Format("DD/MM/YYYY"))
    }
    return &o, nil
}

func (s *ClubService) SaveLittleKidsMembership(o clubs.LittleExplorersKidsMembership) error {
    type Prns struct {
        kidsPrn     string
        guardianPrn string
    }
    query := `
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
    }
    return err
}

func (s *ClubService) SaveLittleKidsMembershipViaWebportal(o clubs.LittleExplorersKidsMembership, adminId int64) error {
    type Prns struct {
        kidsPrn     string
        guardianPrn string
    }
    query := `
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
    }
    return err
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
        sql.Named("kids_membership_id", o.KidsMembershipId.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
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
    o.Set()
    return &o, nil
}

func (s *ClubService) FindLittleKidsActivitiesByActivityId(activityId int64) (*clubs.LittleExplorersKidsActivity, error) {
    query := `SELECT * FROM KIDS_CLUB_ACTIVITY WHERE KIDS_ACTIVITY_ID = :activityId`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.LittleExplorersKidsActivity{}, ""), 1)
    var o clubs.LittleExplorersKidsActivity
    err := s.db.GetContext(s.ctx, &o, query, activityId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    if o.ActivityStartDateTime.Valid {
        g, _ := goment.New(o.ActivityStartDateTime.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.ActivityStartDateTime = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
    }
    if o.ActivityEndDateTime.Valid && o.ActivityEndDateTime.String != "-" {
        g, _ := goment.New(o.ActivityEndDateTime.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.ActivityEndDateTime = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
    }
    return &o, nil
}

func (s *ClubService) ParticipateLittleKidsActivity(o []clubs.LittleExplorersKidsActvParticipation) error {
    query := `
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
            sql.Named("kidsActivityId", ap.KidsActivityId),
            sql.Named("kidsMembershipId", ap.KidsMembershipId),
            sql.Named("activityDateTime", ap.ActivityDateTime),
        )
        if err != nil {
            return err
        }
    }
    return tx.Commit()
}

func (s *ClubService) SaveLittleKidsActivity(o clubs.LittleExplorersKidsActivity, adminId int64) error {
    query := `
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
    }
    return err
}

func (s *ClubService) UpdateLittleKidsActivity(o clubs.LittleExplorersKidsActivity, adminId int64) error {
    query := `
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
        sql.Named("kids_activity_id", o.KidsActivityId.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ClubService) SaveLittleKidsAboutUs(o clubs.LittleExplorersKidsAboutUs, adminId int64) (int64, error) {
    query := `
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

func (s *ClubService) UpdateLittleKidsAboutUs(o clubs.LittleExplorersKidsAboutUs, adminId int64) error {
    query := `
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
        sql.Named("kids_club_id", o.KidsClubId.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ClubService) FindLittleKidsAboutUs() (*clubs.LittleExplorersKidsAboutUs, error) {
    query := `SELECT * FROM KIDS_CLUB_INFO`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.LittleExplorersKidsAboutUs{}, ""), 1)
    var o clubs.LittleExplorersKidsAboutUs
    err := s.db.GetContext(s.ctx, &o, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *ClubService) FindAllAppLittleKidsActivities(offset int, limit int, isHome bool) ([]clubs.LittleExplorersKidsActivity, error) {
    m := map[string]string{
        "kca.ATTENDEES": "",
    }
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
    query = strings.Replace(query, "kca.*", utils.GetDbColsWithReplace(clubs.LittleExplorersKidsActivity{}, "kca.", m), 1)
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

func (s *ClubService) FindAllAppLittleKidsMemberships(userId int64) ([]clubs.LittleExplorersKidsMembership, error) {
    lx := make([]clubs.LittleExplorersKidsMembership, 0)
    lid := make([]string, 0)
    patient, err := s.applicationUserService.FindByUserId(userId, nil)
    if err != nil {
        return nil, err
    }
    patientFamily, err := s.applicationUserFamilyService.FindAllActiveByUserId(userId, 0, 100, true, true, nil)
    if err != nil {
        return nil, err
    }

    for i := range patientFamily {
        family := patientFamily[i]
        lid = append(lid, family.NokPrn.String)
    }
    
    familyPrns := strings.Join(lid, ",")
    q := `SELECT * FROM KIDS_CLUB_MEMBERSHIP WHERE GUARDIAN_PRN = :prn`
    if len(lid) > 0 {
        q = `
            SELECT * FROM KIDS_CLUB_MEMBERSHIP WHERE KIDS_PRN IN (%s)
            UNION
            SELECT * FROM KIDS_CLUB_MEMBERSHIP WHERE GUARDIAN_PRN = :prn
        `
        q = fmt.Sprintf(q, familyPrns)
    }

    q = strings.Replace(q, "*", utils.GetDbCols(clubs.LittleExplorersKidsMembership{}, ""), 1)
    err = s.db.SelectContext(s.ctx, &lx, q, patient.MasterPrn.String)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    for i := range lx {
        lx[i].Set()
        if lx[i].KidsMembershipJoinDate.Valid {
            g, _ := goment.New(lx[i].KidsMembershipJoinDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            lx[i].KidsMembershipJoinDate = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
        }
        if lx[i].KidsDob.Valid {
            g, _ := goment.New(lx[i].KidsDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            lx[i].KidsDob = utils.NewNullString(g.Format("DD/MM/YYYY"))
        }
    }

    return lx, nil
}

func (s *ClubService) FindAllLittleKidsAttendees(kidsActivityId int64, offset int, limit int) ([]clubs.LittleExplorersKidsMembership, error) {
    m := map[string]string{
        "kcm.ACTIVITY_DATE_TIME": "",
    }
    query := `
        SELECT kcm.*, kcap.ACTIVITY_DATE_TIME
        FROM KIDS_CLUB_ACTV_PARTICIPATION kcap
        JOIN KIDS_CLUB_MEMBERSHIP kcm ON kcap.KIDS_MEMBERSHIP_ID = kcm.KIDS_MEMBERSHIP_ID
        WHERE kcap.KIDS_ACTIVITY_ID = :kidsActivityId
        ORDER BY KIDS_MEMBERSHIP_NUMBER DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "kcm.*", utils.GetDbColsWithReplace(clubs.LittleExplorersKidsMembership{}, "kcm.", m), 1)
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

func (s *ClubService) FindAllLittleKidsMemberships(offset int, limit int) ([]clubs.LittleExplorersKidsMembership, error) {
    query := `SELECT * FROM KIDS_CLUB_MEMBERSHIP ORDER BY KIDS_MEMBERSHIP_NUMBER DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.LittleExplorersKidsMembership{}, ""), 1)
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

func (s *ClubService) FindAllLittleKidsActivities(offset int, limit int) ([]clubs.LittleExplorersKidsActivity, error) {
    m := map[string]string{
        "kca.ATTENDEES": "",
    }
    query := `
        SELECT kca.*, (SELECT COUNT(*) 
        FROM KIDS_CLUB_ACTV_PARTICIPATION kcap
        WHERE kcap.KIDS_ACTIVITY_ID = kca.KIDS_ACTIVITY_ID) AS ATTENDEES
        FROM KIDS_CLUB_ACTIVITY kca
        ORDER BY ACTIVITY_START_DATETIME DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "kca.*", utils.GetDbColsWithReplace(clubs.LittleExplorersKidsActivity{}, "kca.", m), 1)
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

func (s *ClubService) ListAppLittleKidsActivities(isHome bool, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountAppLittleKidsActivities(isHome)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllAppLittleKidsActivities(pager.GetLowerBound(), pager.PageSize, isHome)
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

func (s *ClubService) ListLittleKidsMemberships(page string, limit string) (*model.PagedList, error) {
    total, err := s.CountLittleKidsMemberships()
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllLittleKidsMemberships(pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountLittleKidsMemberships() (int, error) {
    query := `SELECT COUNT(KIDS_MEMBERSHIP_ID) AS COUNT FROM KIDS_CLUB_MEMBERSHIP`
    var count int
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) ListLittleKidsActivityAttendees(kidsActivityId int64, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountLittleKidsActivityAttendees(kidsActivityId)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllLittleKidsAttendees(kidsActivityId, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountLittleKidsActivityAttendees(kidsActivityId int64) (int, error) {
    query := `SELECT COUNT(KIDS_ACTV_PARTICIPATION_ID) AS COUNT FROM KIDS_CLUB_ACTV_PARTICIPATION WHERE KIDS_ACTIVITY_ID = :kidsActivityId`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, kidsActivityId)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) ListLittleKidsAttendeesByKeyword(kidsActivityId int64, x dto.SearchKeyword2Dto, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountLittleKidsAttendeesByKeyword(kidsActivityId, x)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindLittleKidsAttendeesByKeyword(kidsActivityId, x, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountLittleKidsAttendeesByKeyword(kidsActivityId int64, x dto.SearchKeyword2Dto) (int, error) {
    conditions, args := buildLittleKidsAttendeesKeywordConditions(kidsActivityId, x)
    base := `SELECT COUNT(kcap.KIDS_ACTV_PARTICIPATION_ID) AS COUNT
             FROM KIDS_CLUB_ACTV_PARTICIPATION kcap
             JOIN KIDS_CLUB_MEMBERSHIP kcm ON kcap.KIDS_MEMBERSHIP_ID = kcm.KIDS_MEMBERSHIP_ID`
    query := base + whereClause(conditions)

    var count int
    err := s.db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) ListLittleKidsMembershipByKeyword(x dto.SearchKeyword2Dto, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountLittleKidsMembershipByKeyword(x)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindLittleKidsMembershipByKeyword(x, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountLittleKidsMembershipByKeyword(x dto.SearchKeyword2Dto) (int, error) {
    conditions, args := buildLittleKidsMembershipKeywordConditions(x)
    base := `SELECT COUNT(kcm.KIDS_MEMBERSHIP_ID) AS COUNT FROM KIDS_CLUB_MEMBERSHIP kcm`
    query := base + whereClause(conditions)

    var count int
    err := s.db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) FindAllLittleKidsActivitiesForExcel() ([]clubs.LittleExplorersKidsActivity, error) {
    m := map[string]string{
        "kca.ATTENDEES": "",
    }
    query := `
        SELECT kca.*, (SELECT COUNT(*) 
        FROM KIDS_CLUB_ACTV_PARTICIPATION kcap
        WHERE kcap.KIDS_ACTIVITY_ID = kca.KIDS_ACTIVITY_ID) AS ATTENDEES
        FROM KIDS_CLUB_ACTIVITY kca
        ORDER BY ACTIVITY_START_DATETIME
    `
    query = strings.Replace(query, "kca.*", utils.GetDbColsWithReplace(clubs.LittleExplorersKidsActivity{}, "kca.", m), 1)
    list := make([]clubs.LittleExplorersKidsActivity, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
        if list[i].ActivityStartDateTime.Valid {
            g, _ := goment.New(list[i].ActivityStartDateTime.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].ActivityStartDateTime = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
        }
        if list[i].ActivityEndDateTime.Valid && list[i].ActivityEndDateTime.String != "-" {
            g, _ := goment.New(list[i].ActivityEndDateTime.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].ActivityEndDateTime = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
        }
    }
    return list, nil
}

func (s *ClubService) FindLittleKidsAttendeesByKeyword(kidsActivityId int64, x dto.SearchKeyword2Dto, offset int, limit int) ([]clubs.LittleExplorersKidsMembership, error) {
    conditions, args := buildLittleKidsAttendeesKeywordConditions(kidsActivityId, x)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    m := map[string]string{
        "kcm.ACTIVITY_DATE_TIME": "",
    }
    base := `
        SELECT kcm.*, kcap.ACTIVITY_DATE_TIME
        FROM KIDS_CLUB_ACTV_PARTICIPATION kcap
        JOIN KIDS_CLUB_MEMBERSHIP kcm ON kcap.KIDS_MEMBERSHIP_ID = kcm.KIDS_MEMBERSHIP_ID
    `
    base = strings.Replace(base, "kcm.*", utils.GetDbColsWithReplace(clubs.LittleExplorersKidsMembership{}, "kcm.", m), 1)

    query := base + whereClause(conditions) +
        ` ORDER BY kcm.KIDS_MEMBERSHIP_NUMBER DESC 
          OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    list := make([]clubs.LittleExplorersKidsMembership, 0)
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].SetAttendees()
    }
    return list, nil
}

func (s *ClubService) FindLittleKidsMembershipByKeyword(x dto.SearchKeyword2Dto, offset int, limit int) ([]clubs.LittleExplorersKidsMembership, error) {
    conditions, args := buildLittleKidsMembershipKeywordConditions(x)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    base := `SELECT * FROM KIDS_CLUB_MEMBERSHIP kcm`

    query := base + whereClause(conditions) +
        ` ORDER BY kcm.KIDS_MEMBERSHIP_NUMBER DESC 
          OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    list := make([]clubs.LittleExplorersKidsMembership, 0)
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
}

func (s *ClubService) FindAllLittleKidsMembershipForExcel() ([]clubs.LittleExplorersKidsMembership, error) {
    query := `SELECT * FROM KIDS_CLUB_MEMBERSHIP ORDER BY KIDS_MEMBERSHIP_NUMBER DESC`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.LittleExplorersKidsMembership{}, ""), 1)
    list := make([]clubs.LittleExplorersKidsMembership, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
        if list[i].KidsMembershipJoinDate.Valid {
            g, _ := goment.New(list[i].KidsMembershipJoinDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].KidsMembershipJoinDate = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
        }
        if list[i].KidsDob.Valid {
            g, _ := goment.New(list[i].KidsDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].KidsDob = utils.NewNullString(g.Format("DD/MM/YYYY"))
        }
    }
    return list, nil
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
    query := `SELECT COUNT(KIDS_ACTIVITY_ID) AS COUNT FROM KIDS_CLUB_ACTIVITY`
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

func (s *ClubService) FindLittleKidsActivitiesByKeyword(keyword string, keyword2 string, keyword3 string, offset int, limit int) ([]clubs.LittleExplorersKidsActivity, error) {
    conditions, args := buildLittleKidsActivitiesKeywordConditions(keyword, keyword2, keyword3)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    m := map[string]string{
        "kca.ATTENDEES": "",
    }
    base := `
        SELECT kca.*, (SELECT COUNT(*)
        FROM KIDS_CLUB_ACTV_PARTICIPATION kcap
        WHERE kcap.KIDS_ACTIVITY_ID = kca.KIDS_ACTIVITY_ID) AS ATTENDEES
        FROM KIDS_CLUB_ACTIVITY kca
    `
    base = strings.Replace(base, "kca.*", utils.GetDbColsWithReplace(clubs.LittleExplorersKidsActivity{}, "kca.", m), 1)

    query := base + whereClause(conditions) +
        ` ORDER BY kca.ACTIVITY_START_DATETIME DESC 
          OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    list := make([]clubs.LittleExplorersKidsActivity, 0)
    err := s.db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
}

func (s *ClubService) GenerateKidsMembershipNo() (string, error) {
    kidsMembershipNo := ""
    query := `SELECT GENERATE_KIDS_MEMBERSHIP() AS KIDS_MEMBERSHIP_NUMBER FROM DUAL`
    err := s.db.GetContext(s.ctx, &kidsMembershipNo, query)
    if err != nil {
        utils.LogError(err)
        return "", err
    }
    return kidsMembershipNo, nil
}

func (s *ClubService) ExistsLittleKidsByPrn(docNumber string) (bool, error) {
    res, ex, err := s.vesaliusGeoService.ClubsGetPatientData(docNumber)
    if err != nil {
        return false, err
    }
    
    if len(res.Patients) > 0 && ex == nil {
        o := res.Patients[0]
        query := `
            SELECT COUNT(*) AS COUNT FROM DUAL 
            WHERE EXISTS (
              SELECT 1 FROM KIDS_CLUB_MEMBERSHIP
              WHERE KIDS_PRN = :prn
            )
        `
        var count int
        err := s.db.GetContext(s.ctx, &count, query, o.Prn)
        if err != nil {
            utils.LogError(err)
            return false, err
        }
        return count > 0, nil
    }

    return false, err
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

func buildLittleKidsAttendeesKeywordConditions(kidsActivityId int64, x dto.SearchKeyword2Dto) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if x.Keyword != "" {
        conds = append(conds, `(LOWER(kcm.KIDS_PRN) LIKE :keyword OR LOWER(kcm.KIDS_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(x.Keyword)))
    }
    if x.Keyword2 != "" {
        conds = append(conds, `(LOWER(kcm.GUARDIAN_PRN) LIKE :keyword2 OR LOWER(kcm.GUARDIAN_NAME) LIKE :keyword2)`)
        args = append(args, sql.Named("keyword2", strings.ToLower(x.Keyword2)))
    }
    conds = append(conds, `kcap.KIDS_ACTIVITY_ID = :kidsActivityId`)
    args = append(args, sql.Named("kidsActivityId", kidsActivityId))
    return conds, args
}

func buildLittleKidsMembershipKeywordConditions(x dto.SearchKeyword2Dto) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if x.Keyword != "" {
        conds = append(conds, `(LOWER(kcm.KIDS_PRN) LIKE :keyword OR LOWER(kcm.KIDS_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(x.Keyword)))
    }
    if x.Keyword2 != "" {
        conds = append(conds, `(LOWER(kcm.GUARDIAN_PRN) LIKE :keyword2 OR LOWER(kcm.GUARDIAN_NAME) LIKE :keyword2)`)
        args = append(args, sql.Named("keyword2", strings.ToLower(x.Keyword2)))
    }
    return conds, args
}
