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

func (s *ClubService) FindGoldenPearlMembershipById(goldenMembershipId int64) (*clubs.GoldenPearlMembership, error) {
    query := `SELECT * FROM GOLDEN_CLUB_MEMBERSHIP WHERE GOLDEN_MEMBERSHIP_ID = :goldenMembershipId`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.GoldenPearlMembership{}, ""), 1)
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

func (s *ClubService) FindGoldenPearlMembershipByMembershipId(membershipId int64) (*clubs.GoldenPearlMembership, error) {
    query := `SELECT * FROM GOLDEN_CLUB_MEMBERSHIP WHERE GOLDEN_MEMBERSHIP_ID = :membershipId`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.GoldenPearlMembership{}, ""), 1)
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
        o.GoldenDob = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
    }
    if o.NokDob.Valid {
        g, _ := goment.New(o.NokDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.NokDob = utils.NewNullString(g.Format("DD/MM/YYYY"))
    }
    return &o, nil
}

func (s *ClubService) SaveGoldenPearlMembership(o clubs.GoldenPearlMembership) error {
    type Prns struct {
        goldenPrn string
        nokPrn    string
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
    }
    return err
}

func (s *ClubService) SaveGoldenPearlMembershipViaWebportal(o clubs.GoldenPearlMembership, adminId int64) error {
    type Prns struct {
        goldenPrn string
        nokPrn    string
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
    }
    return err
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
        sql.Named("golden_membership_id", o.GoldenMembershipId.Int64),
    )
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
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

func (s *ClubService) FindGoldenPearlActivitiesByActivityId(activityId int64) (*clubs.GoldenPearlActivity, error) {
    query := `SELECT * FROM GOLDEN_CLUB_ACTIVITY WHERE GOLDEN_ACTIVITY_ID = :activityId`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.GoldenPearlActivity{}, ""), 1)
    var o clubs.GoldenPearlActivity
    err := s.db.GetContext(s.ctx, &o, query, activityId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
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

func (s *ClubService) ParticipateGoldenPearlActivity(o []clubs.GoldenPearlActvParticipation) error {
    query := `
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
            sql.Named("goldenActivityId", ap.GoldenActivityId),
            sql.Named("goldenMembershipId", ap.GoldenMembershipId),
            sql.Named("activityDateTime", ap.ActivityDateTime),
        )
        if err != nil {
            return err
        }
    }
    return tx.Commit()
}

func (s *ClubService) SaveGoldenPearlActivity(o clubs.GoldenPearlActivity, adminId int64) error {
    query := `
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
    }
    return err
}

func (s *ClubService) UpdateGoldenPearlActivity(o clubs.GoldenPearlActivity, adminId int64) error {
    query := `
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
        sql.Named("golden_activity_id", o.GoldenActivityId.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ClubService) SaveGoldenPearlAboutUs(o clubs.GoldenPearlAboutUs, adminId int64) (int64, error) {
    query := `
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

func (s *ClubService) UpdateGoldenPearlAboutUs(o clubs.GoldenPearlAboutUs, adminId int64) error {
    query := `
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
        sql.Named("golden_club_id", o.GoldenClubId.Int64),
    )
    if err != nil {
        utils.LogError(err)
    }
    return err
}

func (s *ClubService) FindGoldenPearlAboutUs() (*clubs.GoldenPearlAboutUs, error) {
    query := `SELECT * FROM GOLDEN_CLUB_INFO`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.GoldenPearlAboutUs{}, ""), 1)
    var o clubs.GoldenPearlAboutUs
    err := s.db.GetContext(s.ctx, &o, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *ClubService) FindAllAppGoldenPearlActivities(offset int, limit int, isHome bool) ([]clubs.GoldenPearlActivity, error) {
    m := map[string]string{
        "gca.ATTENDEES": "",
    }
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
    query = strings.Replace(query, "gca.*", utils.GetDbColsWithReplace(clubs.GoldenPearlActivity{}, "gca.", m), 1)
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

func (s *ClubService) FindAllAppGoldenPearlMemberships(userId int64) ([]clubs.GoldenPearlMembership, error) {
    lx := make([]clubs.GoldenPearlMembership, 0)
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
    q := `SELECT * FROM GOLDEN_CLUB_MEMBERSHIP WHERE NOK_PRN = :prn`
    if len(lid) > 0 {
        q = `
            SELECT * FROM GOLDEN_CLUB_MEMBERSHIP WHERE GOLDEN_PRN IN (%s)
            UNION
            SELECT * FROM GOLDEN_CLUB_MEMBERSHIP WHERE NOK_PRN = :prn
        `
        q = fmt.Sprintf(q, familyPrns)
    }

    q = strings.Replace(q, "*", utils.GetDbCols(clubs.GoldenPearlMembership{}, ""), 1)
    err = s.db.SelectContext(s.ctx, &lx, q, patient.MasterPrn.String)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    for i := range lx {
        lx[i].Set()
        if lx[i].GoldenMembershipJoinDate.Valid {
            g, _ := goment.New(lx[i].GoldenMembershipJoinDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            lx[i].GoldenMembershipJoinDate = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
        }
        if lx[i].GoldenDob.Valid {
            g, _ := goment.New(lx[i].GoldenDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            lx[i].GoldenDob = utils.NewNullString(g.Format("DD/MM/YYYY"))
        }
    }

    return lx, nil
}

func (s *ClubService) FindAllGoldenPearlAttendees(goldenActivityId int64, offset int, limit int) ([]clubs.GoldenPearlMembership, error) {
    m := map[string]string{
        "gcm.ACTIVITY_DATE_TIME": "",
    }
    query := `
        SELECT gcm.*, gcap.ACTIVITY_DATE_TIME
        FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
        JOIN GOLDEN_CLUB_MEMBERSHIP gcm ON gcap.GOLDEN_MEMBERSHIP_ID = gcm.GOLDEN_MEMBERSHIP_ID
        WHERE gcap.GOLDEN_ACTIVITY_ID = :goldenActivityId
        ORDER BY GOLDEN_MEMBERSHIP_NUMBER DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "gcm.*", utils.GetDbColsWithReplace(clubs.GoldenPearlMembership{}, "gcm.", m), 1)
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

func (s *ClubService) FindAllGoldenPearlMemberships(offset int, limit int) ([]clubs.GoldenPearlMembership, error) {
    query := `SELECT * FROM GOLDEN_CLUB_MEMBERSHIP ORDER BY GOLDEN_MEMBERSHIP_NUMBER DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.GoldenPearlMembership{}, ""), 1)
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

func (s *ClubService) FindAllGoldenPearlActivities(offset int, limit int) ([]clubs.GoldenPearlActivity, error) {
    m := map[string]string{
        "gca.ATTENDEES": "",
    }
    query := `
        SELECT gca.*, (SELECT COUNT(*) 
        FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
        WHERE gcap.GOLDEN_ACTIVITY_ID = gca.GOLDEN_ACTIVITY_ID) AS ATTENDEES
        FROM GOLDEN_CLUB_ACTIVITY gca
        ORDER BY ACTIVITY_START_DATETIME DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "gca.*", utils.GetDbColsWithReplace(clubs.GoldenPearlActivity{}, "gca.", m), 1)
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

func (s *ClubService) ListAppGoldenPearlActivities(isHome bool, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountAppGoldenPearlActivities(isHome)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllAppGoldenPearlActivities(pager.GetLowerBound(), pager.PageSize, isHome)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountAppGoldenPearlActivities(isHome bool) (int, error) {
    query := `
        SELECT COUNT(GOLDEN_ACTIVITY_ID) AS COUNT FROM GOLDEN_CLUB_ACTIVITY
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
            SELECT COUNT(GOLDEN_ACTIVITY_ID) AS COUNT FROM GOLDEN_CLUB_ACTIVITY
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

func (s *ClubService) ListGoldenPearlMemberships(page string, limit string) (*model.PagedList, error) {
    total, err := s.CountGoldenPearlMemberships()
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllGoldenPearlMemberships(pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountGoldenPearlMemberships() (int, error) {
    query := `SELECT COUNT(GOLDEN_PEARL_MEMBERSHIP_ID) AS COUNT FROM GOLDEN_PEARL_CLUB_MEMBERSHIP`
    var count int
    err := s.db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) ListGoldenPearlActivityAttendees(goldenActivityId int64, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountGoldenPearlActivityAttendees(goldenActivityId)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllGoldenPearlAttendees(goldenActivityId, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountGoldenPearlActivityAttendees(goldenActivityId int64) (int, error) {
    query := `SELECT COUNT(GOLDEN_ACTV_PARTICIPATION_ID) AS COUNT FROM GOLDEN_CLUB_ACTV_PARTICIPATION WHERE GOLDEN_ACTIVITY_ID = :goldenActivityId`
    var count int
    err := s.db.GetContext(s.ctx, &count, query, goldenActivityId)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) ListGoldenPearlAttendeesByKeyword(goldenActivityId int64, keyword string, keyword2 string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountGoldenPearlAttendeesByKeyword(goldenActivityId, keyword, keyword2)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindGoldenPearlAttendeesByKeyword(goldenActivityId, keyword, keyword2, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountGoldenPearlAttendeesByKeyword(goldenActivityId int64, keyword string, keyword2 string) (int, error) {
    conditions, args := buildGoldenPearlAttendeesKeywordConditions(goldenActivityId, keyword, keyword2)
    base := `SELECT COUNT(gcap.GOLDEN_ACTV_PARTICIPATION_ID) AS COUNT
             FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
             JOIN GOLDEN_CLUB_MEMBERSHIP gcm ON gcap.GOLDEN_MEMBERSHIP_ID = gcm.GOLDEN_MEMBERSHIP_ID`
    query := base + whereClause(conditions)

    var count int
    err := s.db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) ListGoldenPearlMembershipByKeyword(x dto.SearchKeyword2Dto, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountGoldenPearlMembershipByKeyword(x)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindGoldenPearlMembershipByKeyword(x, pager.GetLowerBound(), pager.PageSize)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *ClubService) CountGoldenPearlMembershipByKeyword(x dto.SearchKeyword2Dto) (int, error) {
    conditions, args := buildGoldenPearlMembershipKeywordConditions(x)
    base := `SELECT COUNT(gcm.GOLDEN_MEMBERSHIP_ID) AS COUNT FROM GOLDEN_CLUB_MEMBERSHIP gcm`
    query := base + whereClause(conditions)

    var count int
    err := s.db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *ClubService) FindAllGoldenPearlActivitiesForExcel() ([]clubs.GoldenPearlActivity, error) {
    m := map[string]string{
        "gca.ATTENDEES": "",
    }
    query := `
        SELECT gca.*, (SELECT COUNT(*) 
        FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
        WHERE gcap.GOLDEN_ACTIVITY_ID = gca.GOLDEN_ACTIVITY_ID) AS ATTENDEES
        FROM GOLDEN_CLUB_ACTIVITY gca
        ORDER BY ACTIVITY_START_DATETIME
    `
    query = strings.Replace(query, "gca.*", utils.GetDbColsWithReplace(clubs.GoldenPearlActivity{}, "gca.", m), 1)
    list := make([]clubs.GoldenPearlActivity, 0)
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

func (s *ClubService) FindAllGoldenPearlAttendeesForExcel(goldenActivityId int64) ([]clubs.GoldenPearlMembership, error) {
    m := map[string]string{
        "gcm.ACTIVITY_DATE_TIME": "",
    }
    query := `
        SELECT gcm.*, gcap.ACTIVITY_DATE_TIME
        FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
        JOIN GOLDEN_CLUB_MEMBERSHIP gcm ON gcap.GOLDEN_MEMBERSHIP_ID = gcm.GOLDEN_MEMBERSHIP_ID
        WHERE gcap.GOLDEN_ACTIVITY_ID = :goldenActivityId
        ORDER BY GOLDEN_MEMBERSHIP_NUMBER DESC
    `
    query = strings.Replace(query, "gcm.*", utils.GetDbColsWithReplace(clubs.GoldenPearlMembership{}, "gcm.", m), 1)
    list := make([]clubs.GoldenPearlMembership, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].SetAttendees()
        if list[i].ActivityJoinDate.Valid {
            g, _ := goment.New(list[i].ActivityJoinDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].ActivityJoinDate = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
        }
        if list[i].GoldenMembershipJoinDate.Valid {
            g, _ := goment.New(list[i].GoldenMembershipJoinDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].GoldenMembershipJoinDate = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
        }
        if list[i].GoldenDob.Valid {
            g, _ := goment.New(list[i].GoldenDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].GoldenDob = utils.NewNullString(g.Format("DD/MM/YYYY"))
        }
    }
    return list, nil
}

func (s *ClubService) FindGoldenPearlAttendeesByKeyword(goldenActivityId int64, keyword string, keyword2 string, offset int, limit int) ([]clubs.GoldenPearlMembership, error) {
    conditions, args := buildGoldenPearlAttendeesKeywordConditions(goldenActivityId, keyword, keyword2)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    m := map[string]string{
        "gcm.ACTIVITY_DATE_TIME": "",
    }
    base := `
        SELECT gcm.*, gcap.ACTIVITY_DATE_TIME
        FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
        JOIN GOLDEN_CLUB_MEMBERSHIP gcm ON gcap.GOLDEN_MEMBERSHIP_ID = gcm.GOLDEN_MEMBERSHIP_ID
    `
    base = strings.Replace(base, "gcm.*", utils.GetDbColsWithReplace(clubs.GoldenPearlMembership{}, "gcm.", m), 1)

    query := base + whereClause(conditions) +
        ` ORDER BY gcm.GOLDEN_MEMBERSHIP_NUMBER DESC 
          OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    list := make([]clubs.GoldenPearlMembership, 0)
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

func (s *ClubService) FindGoldenPearlMembershipByKeyword(x dto.SearchKeyword2Dto, offset int, limit int) ([]clubs.GoldenPearlMembership, error) {
    conditions, args := buildGoldenPearlMembershipKeywordConditions(x)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    base := `SELECT * FROM GOLDEN_CLUB_MEMBERSHIP gcm`

    query := base + whereClause(conditions) +
        ` ORDER BY gcm.GOLDEN_MEMBERSHIP_NUMBER DESC 
          OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    list := make([]clubs.GoldenPearlMembership, 0)
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

func (s *ClubService) FindAllGoldenPearlMembershipForExcel() ([]clubs.GoldenPearlMembership, error) {
    query := `SELECT * FROM GOLDEN_CLUB_MEMBERSHIP ORDER BY GOLDEN_MEMBERSHIP_NUMBER DESC`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.GoldenPearlMembership{}, ""), 1)
    list := make([]clubs.GoldenPearlMembership, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
        if list[i].GoldenMembershipJoinDate.Valid {
            g, _ := goment.New(list[i].GoldenMembershipJoinDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].GoldenMembershipJoinDate = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
        }
        if list[i].GoldenDob.Valid {
            g, _ := goment.New(list[i].GoldenDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].GoldenDob = utils.NewNullString(g.Format("DD/MM/YYYY"))
        }
    }
    return list, nil
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
    query := `SELECT COUNT(GOLDEN_ACTIVITY_ID) AS COUNT FROM GOLDEN_CLUB_ACTIVITY`
    var count int
    err := s.db.GetContext(s.ctx, &count, query)
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

func (s *ClubService) FindGoldenPearlActivitiesByKeyword(keyword string, keyword2 string, keyword3 string, offset int, limit int) ([]clubs.GoldenPearlActivity, error) {
    conditions, args := buildGoldenPearlActivitiesKeywordConditions(keyword, keyword2, keyword3)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    m := map[string]string{
        "gca.ATTENDEES": "",
    }
    base := `
        SELECT gca.*, (SELECT COUNT(*)
        FROM GOLDEN_CLUB_ACTV_PARTICIPATION gcap
        WHERE gcap.GOLDEN_ACTIVITY_ID = gca.GOLDEN_ACTIVITY_ID) AS ATTENDEES
        FROM GOLDEN_CLUB_ACTIVITY gca
    `
    base = strings.Replace(base, "gca.*", utils.GetDbColsWithReplace(clubs.GoldenPearlActivity{}, "gca.", m), 1)

    query := base + whereClause(conditions) +
        ` ORDER BY gca.ACTIVITY_START_DATETIME DESC 
          OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    list := make([]clubs.GoldenPearlActivity, 0)
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

func (s *ClubService) GenerateGoldenMembershipNo() (string, error) {
    goldenMembershipNo := ""
    query := `SELECT GENERATE_GOLDEN_MEMBERSHIP() AS GOLDEN_MEMBERSHIP_NUMBER FROM DUAL`
    err := s.db.GetContext(s.ctx, &goldenMembershipNo, query)
    if err != nil {
        utils.LogError(err)
        return "", err
    }
    return goldenMembershipNo, nil
}

func (s *ClubService) ExistsGoldenPearlByPrn(docNumber string) (bool, error) {
    res, ex, err := s.vesaliusGeoService.ClubsGetPatientData(docNumber)
    if err != nil {
        return false, err
    }
    
    if len(res.Patients) > 0 && ex == nil {
        o := res.Patients[0]
        query := `
            SELECT COUNT(*) AS COUNT FROM DUAL 
            WHERE EXISTS (
              SELECT 1 FROM GOLDEN_CLUB_MEMBERSHIP
              WHERE GOLDEN_PRN = :prn
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

func buildGoldenPearlActivitiesKeywordConditions(keyword string, keyword2 string, keyword3 string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `(LOWER(gca.GOLDEN_ACTIVITY_CODE) LIKE :keyword OR LOWER(gca.GOLDEN_ACTIVITY_NAME) LIKE :keyword)`)
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

func buildGoldenPearlAttendeesKeywordConditions(goldenActivityId int64, keyword string, keyword2 string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `(LOWER(gcm.GOLDEN_PRN) LIKE :keyword OR LOWER(gcm.GOLDEN_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        conds = append(conds, `(LOWER(gcm.NOK_PRN) LIKE :keyword2 OR LOWER(gcm.NOK_NAME) LIKE :keyword2)`)
        args = append(args, sql.Named("keyword2", strings.ToLower(keyword2)))
    }
    conds = append(conds, `gcap.GOLDEN_ACTIVITY_ID = :goldenActivityId`)
    args = append(args, sql.Named("goldenActivityId", goldenActivityId))
    return conds, args
}

func buildGoldenPearlMembershipKeywordConditions(x dto.SearchKeyword2Dto) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if x.Keyword != "" {
        conds = append(conds, `(LOWER(gcm.GOLDEN_PRN) LIKE :keyword OR LOWER(gcm.GOLDEN_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(x.Keyword)))
    }
    if x.Keyword2 != "" {
        conds = append(conds, `(LOWER(gcm.NOK_PRN) LIKE :keyword2 OR LOWER(gcm.NOK_NAME) LIKE :keyword2)`)
        args = append(args, sql.Named("keyword2", strings.ToLower(x.Keyword2)))
    }
    return conds, args
}
