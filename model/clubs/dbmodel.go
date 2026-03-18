package clubs

import (
    "vesaliusm/utils"

    "github.com/guregu/null/v6"
)

type LittleExplorersKidsMyActivity struct {
    KidsName             null.String `json:"kidsName" db:"KIDS_NAME" swaggertype:"string"`
    KidsMembershipNumber null.String `json:"kidsMembershipNumber" db:"KIDS_MEMBERSHIP_NUMBER" swaggertype:"string"`
    KidsActivityName     null.String `json:"kidsActivityName" db:"KIDS_ACTIVITY_NAME" swaggertype:"string"`
    ActivityDateTime     null.String `json:"activityDateTime" db:"ACTIVITY_DATE_TIME" swaggertype:"string"`
}

type GoldenPearlMyActivity struct {
    GoldenName             null.String `json:"goldenName" db:"GOLDEN_NAME" swaggertype:"string"`
    GoldenMembershipNumber null.String `json:"goldenMembershipNumber" db:"GOLDEN_MEMBERSHIP_NUMBER" swaggertype:"string"`
    GoldenActivityName     null.String `json:"goldenActivityName" db:"GOLDEN_ACTIVITY_NAME" swaggertype:"string"`
    ActivityDateTime       null.String `json:"activityDateTime" db:"ACTIVITY_DATE_TIME" swaggertype:"string"`
}

type LittleExplorersKidsMembership struct {
    KidsMembershipId       null.Int64  `json:"kidsMembershipID" db:"KIDS_MEMBERSHIP_ID" swaggertype:"integer"`
    KidsMembershipNumber   null.String `json:"kidsMembershipNumber" db:"KIDS_MEMBERSHIP_NUMBER" swaggertype:"string"`
    KidsMembershipJoinDate null.String `json:"kidsMembershipJoinDate" db:"KIDS_MEMBERSHIP_JOIN_DATE" swaggertype:"string"`
    KidsPrn                null.String `json:"kidsPrn" db:"KIDS_PRN" swaggertype:"string"`
    KidsName               null.String `json:"kidsName" db:"KIDS_NAME" swaggertype:"string"`
    KidsDob                null.String `json:"kidsDob" db:"KIDS_DOB" swaggertype:"string"`
    KidsDocType            null.String `json:"kidsDocType" db:"KIDS_DOC_TYPE" swaggertype:"string"`
    KidsDocNumber          null.String `json:"kidsDocNumber" db:"KIDS_DOC_NUMBER" swaggertype:"string"`
    KidsGender             null.String `json:"kidsGender" db:"KIDS_GENDER" swaggertype:"string"`
    KidsNationality        null.String `json:"kidsNationality" db:"KIDS_NATIONALITY" swaggertype:"string"`
    KidsEmail              null.String `json:"kidsEmail" db:"KIDS_EMAIL" swaggertype:"string"`
    GuardianPrn            null.String `json:"guardianPrn" db:"GUARDIAN_PRN" swaggertype:"string"`
    GuardianName           null.String `json:"guardianName" db:"GUARDIAN_NAME" swaggertype:"string"`
    GuardianDob            null.String `json:"guardianDob" db:"GUARDIAN_DOB" swaggertype:"string"`
    GuardianDocType        null.String `json:"guardianDocType" db:"GUARDIAN_DOC_TYPE" swaggertype:"string"`
    GuardianDocNumber      null.String `json:"guardianDocNumber" db:"GUARDIAN_DOC_NUMBER" swaggertype:"string"`
    GuardianGender         null.String `json:"guardianGender" db:"GUARDIAN_GENDER" swaggertype:"string"`
    GuardianNationality    null.String `json:"guardianNationality" db:"GUARDIAN_NATIONALITY" swaggertype:"string"`
    GuardianEmail          null.String `json:"guardianEmail" db:"GUARDIAN_EMAIL" swaggertype:"string"`
    GuardianHomeContact    null.String `json:"guardianHomeContact" db:"GUARDIAN_HOME_CONTACT" swaggertype:"string"`
    GuardianMobileContact  null.String `json:"guardianMobileContact" db:"GUARDIAN_MOBILE_CONTACT" swaggertype:"string"`
    GuardianAddress1       null.String `json:"guardianAddress1" db:"GUARDIAN_ADDRESS1" swaggertype:"string"`
    GuardianAddress2       null.String `json:"guardianAddress2" db:"GUARDIAN_ADDRESS2" swaggertype:"string"`
    GuardianAddress3       null.String `json:"guardianAddress3" db:"GUARDIAN_ADDRESS3" swaggertype:"string"`
    GuardianPostCode       null.String `json:"guardianPostCode" db:"GUARDIAN_POSTCODE" swaggertype:"string"`
    GuardianState          null.String `json:"guardianState" db:"GUARDIAN_STATE" swaggertype:"string"`
    GuardianCountryCode    null.String `json:"guardianCountryCode" db:"GUARDIAN_COUNTRY_CODE" swaggertype:"string"`
    Relationship           null.String `json:"relationship" db:"RELATIONSHIP" swaggertype:"string"`
    PreferredLanguage      null.String `json:"preferredLanguage" db:"PREFERRED_LANGUAGE" swaggertype:"string"`
    IsActive               null.String `json:"isActive" db:"IS_ACTIVE" swaggertype:"string"`

    // Fields used in some contexts (e.g., attendee export)
    ActivityJoinDate null.String `json:"activityJoinDate" db:"ACTIVITY_DATE_TIME" swaggertype:"string"`

    // Excel-specific fields (not stored in DB)
    ActivityJoinDateExcel       string `json:"activityJoinDateExcel"`
    KidsMembershipJoinDateExcel string `json:"kidsMembershipJoinDateExcel"`
    KidsDobExcel                string `json:"kidsDobExcel"`
    GuardianDobExcel            string `json:"guardianDobExcel"`
}

func (o *LittleExplorersKidsMembership) SetWebAdmin() {
    if !o.KidsPrn.Valid {
        o.KidsPrn = utils.NewNullString("-")
    }

    if !o.GuardianPrn.Valid {
        o.GuardianPrn = utils.NewNullString("-")
    }
}

func (o *LittleExplorersKidsMembership) Set() {
    if !o.KidsPrn.Valid {
        o.KidsPrn = utils.NewNullString("-")
    }

    if !o.KidsDob.Valid {
        o.KidsDob = utils.NewNullString("-")
    }

    if !o.KidsEmail.Valid {
        o.KidsEmail = utils.NewNullString("-")
    }

    if !o.GuardianDob.Valid {
        o.GuardianDob = utils.NewNullString("-")
    }

    if !o.GuardianHomeContact.Valid {
        o.GuardianHomeContact = utils.NewNullString("-")
    }

    if !o.GuardianMobileContact.Valid {
        o.GuardianMobileContact = utils.NewNullString("-")
    }

    if !o.GuardianAddress1.Valid {
        o.GuardianAddress1 = utils.NewNullString("-")
    }

    if !o.GuardianAddress2.Valid {
        o.GuardianAddress2 = utils.NewNullString("-")
    }

    if !o.GuardianAddress3.Valid {
        o.GuardianAddress3 = utils.NewNullString("-")
    }

    if !o.GuardianPostCode.Valid {
        o.GuardianPostCode = utils.NewNullString("-")
    }

    if !o.GuardianState.Valid {
        o.GuardianState = utils.NewNullString("-")
    }

    if !o.GuardianCountryCode.Valid {
        o.GuardianCountryCode = utils.NewNullString("-")
    }
}

func (o *LittleExplorersKidsMembership) SetAttendees() {
    if !o.KidsPrn.Valid {
        o.KidsPrn = utils.NewNullString("-")
    }

    if !o.KidsDob.Valid {
        o.KidsDob = utils.NewNullString("-")
    }

    if !o.KidsEmail.Valid {
        o.KidsEmail = utils.NewNullString("-")
    }

    if !o.GuardianPrn.Valid {
        o.GuardianPrn = utils.NewNullString("-")
    }

    if !o.GuardianHomeContact.Valid {
        o.GuardianHomeContact = utils.NewNullString("-")
    }

    if !o.GuardianMobileContact.Valid {
        o.GuardianMobileContact = utils.NewNullString("-")
    }

    if !o.GuardianAddress1.Valid {
        o.GuardianAddress1 = utils.NewNullString("-")
    }

    if !o.GuardianAddress2.Valid {
        o.GuardianAddress2 = utils.NewNullString("-")
    }

    if !o.GuardianAddress3.Valid {
        o.GuardianAddress3 = utils.NewNullString("-")
    }

    if !o.GuardianPostCode.Valid {
        o.GuardianPostCode = utils.NewNullString("-")
    }

    if !o.GuardianState.Valid {
        o.GuardianState = utils.NewNullString("-")
    }

    if !o.GuardianCountryCode.Valid {
        o.GuardianCountryCode = utils.NewNullString("-")
    }
}

type GoldenPearlMembership struct {
    GoldenMembershipId       null.Int64  `json:"golden_membership_id" db:"GOLDEN_MEMBERSHIP_ID" swaggertype:"integer"`
    GoldenMembershipNumber   null.String `json:"goldenMembershipNumber" db:"GOLDEN_MEMBERSHIP_NUMBER" swaggertype:"string"`
    GoldenMembershipJoinDate null.String `json:"goldenMembershipJoinDate" db:"GOLDEN_MEMBERSHIP_JOIN_DATE" swaggertype:"string"`
    GoldenPrn                null.String `json:"goldenPrn" db:"GOLDEN_PRN" swaggertype:"string"`
    GoldenName               null.String `json:"goldenName" db:"GOLDEN_NAME" swaggertype:"string"`
    GoldenDob                null.String `json:"goldenDob" db:"GOLDEN_DOB" swaggertype:"string"`
    GoldenDocType            null.String `json:"goldenDocType" db:"GOLDEN_DOC_TYPE" swaggertype:"string"`
    GoldenDocNumber          null.String `json:"goldenDocNumber" db:"GOLDEN_DOC_NUMBER" swaggertype:"string"`
    GoldenGender             null.String `json:"goldenGender" db:"GOLDEN_GENDER" swaggertype:"string"`
    GoldenNationality        null.String `json:"goldenNationality" db:"GOLDEN_NATIONALITY" swaggertype:"string"`
    GoldenEmail              null.String `json:"goldenEmail" db:"GOLDEN_EMAIL" swaggertype:"string"`
    NokPrn                   null.String `json:"nokPrn" db:"NOK_PRN" swaggertype:"string"`
    NokName                  null.String `json:"nokName" db:"NOK_NAME" swaggertype:"string"`
    NokDob                   null.String `json:"nokDob" db:"NOK_DOB" swaggertype:"string"`
    NokDocType               null.String `json:"nokDocType" db:"NOK_DOC_TYPE" swaggertype:"string"`
    NokDocNumber             null.String `json:"nokDocNumber" db:"NOK_DOC_NUMBER" swaggertype:"string"`
    NokGender                null.String `json:"nokGender" db:"NOK_GENDER" swaggertype:"string"`
    NokNationality           null.String `json:"nokNationality" db:"NOK_NATIONALITY" swaggertype:"string"`
    NokEmail                 null.String `json:"nokEmail" db:"NOK_EMAIL" swaggertype:"string"`
    NokHomeContact           null.String `json:"nokHomeContact" db:"NOK_HOME_CONTACT" swaggertype:"string"`
    NokMobileContact         null.String `json:"nokMobileContact" db:"NOK_MOBILE_CONTACT" swaggertype:"string"`
    NokAddress1              null.String `json:"nokAddress1" db:"NOK_ADDRESS1" swaggertype:"string"`
    NokAddress2              null.String `json:"nokAddress2" db:"NOK_ADDRESS2" swaggertype:"string"`
    NokAddress3              null.String `json:"nokAddress3" db:"NOK_ADDRESS3" swaggertype:"string"`
    NokPostCode              null.String `json:"nokPostCode" db:"NOK_POSTCODE" swaggertype:"string"`
    NokState                 null.String `json:"nokState" db:"NOK_STATE" swaggertype:"string"`
    NokCountryCode           null.String `json:"nokCountryCode" db:"NOK_COUNTRY_CODE" swaggertype:"string"`
    Relationship             null.String `json:"relationship" db:"RELATIONSHIP" swaggertype:"string"`
    PreferredLanguage        null.String `json:"preferredLanguage" db:"PREFERRED_LANGUAGE" swaggertype:"string"`
    IsActive                 null.String `json:"isActive" db:"IS_ACTIVE" swaggertype:"string"`

    // Fields used in attendee contexts
    ActivityJoinDate null.String `json:"activityJoinDate" db:"ACTIVITY_DATE_TIME" swaggertype:"string"`

    // Excel-specific fields (not stored in DB)
    ActivityJoinDateExcel         string `json:"activityJoinDateExcel"`
    GoldenMembershipJoinDateExcel string `json:"goldenMembershipJoinDateExcel"`
    GoldenDobExcel                string `json:"goldenDobExcel"`
    NokDobExcel                   string `json:"nokDobExcel"`
}

func (o *GoldenPearlMembership) SetWebAdmin() {
    if !o.GoldenPrn.Valid {
        o.GoldenPrn = utils.NewNullString("-")
    }

    if !o.NokPrn.Valid {
        o.NokPrn = utils.NewNullString("-")
    }
}

func (o *GoldenPearlMembership) Set() {
    if !o.GoldenPrn.Valid {
        o.GoldenPrn = utils.NewNullString("-")
    }

    if !o.GoldenDob.Valid {
        o.GoldenDob = utils.NewNullString("-")
    }

    if !o.GoldenEmail.Valid {
        o.GoldenEmail = utils.NewNullString("-")
    }

    if !o.NokPrn.Valid {
        o.NokPrn = utils.NewNullString("-")
    }

    if !o.NokDob.Valid {
        o.NokDob = utils.NewNullString("-")
    }

    if !o.NokHomeContact.Valid {
        o.NokHomeContact = utils.NewNullString("-")
    }

    if !o.NokMobileContact.Valid {
        o.NokMobileContact = utils.NewNullString("-")
    }

    if !o.NokAddress1.Valid {
        o.NokAddress1 = utils.NewNullString("-")
    }

    if !o.NokAddress2.Valid {
        o.NokAddress2 = utils.NewNullString("-")
    }

    if !o.NokAddress3.Valid {
        o.NokAddress3 = utils.NewNullString("-")
    }

    if !o.NokPostCode.Valid {
        o.NokPostCode = utils.NewNullString("-")
    }

    if !o.NokState.Valid {
        o.NokState = utils.NewNullString("-")
    }

    if !o.NokCountryCode.Valid {
        o.NokCountryCode = utils.NewNullString("-")
    }
}

func (o *GoldenPearlMembership) SetAttendees() {
    if o.GoldenPrn.Valid {
        o.GoldenPrn = utils.NewNullString("-")
    }

    if !o.GoldenDob.Valid {
        o.GoldenDob = utils.NewNullString("-")
    }

    if !o.GoldenEmail.Valid {
        o.GoldenEmail = utils.NewNullString("-")
    }

    if !o.NokPrn.Valid {
        o.NokPrn = utils.NewNullString("-")
    }

    if !o.NokHomeContact.Valid {
        o.NokHomeContact = utils.NewNullString("-")
    }

    if !o.NokMobileContact.Valid {
        o.NokMobileContact = utils.NewNullString("-")
    }

    if !o.NokAddress1.Valid {
        o.NokAddress1 = utils.NewNullString("-")
    }

    if !o.NokAddress2.Valid {
        o.NokAddress2 = utils.NewNullString("-")
    }

    if !o.NokAddress3.Valid {
        o.NokAddress3 = utils.NewNullString("-")
    }

    if !o.NokPostCode.Valid {
        o.NokPostCode = utils.NewNullString("-")
    }

    if !o.NokState.Valid {
        o.NokState = utils.NewNullString("-")
    }

    if !o.NokCountryCode.Valid {
        o.NokCountryCode = utils.NewNullString("-")
    }
}

type LittleExplorersKidsAboutUs struct {
    KidsClubId          null.Int64  `json:"kids_club_id" db:"KIDS_CLUB_ID" swaggertype:"integer"`
    KidsClubTitle       null.String `json:"kidsClubTitle" db:"KIDS_CLUB_TITLE" swaggertype:"string"`
    KidsClubDesc        null.String `json:"kidsClubDesc" db:"KIDS_CLUB_DESC" swaggertype:"string"`
    KidsClubImage       null.String `json:"kidsClubImage" db:"KIDS_CLUB_IMG" swaggertype:"string"`
    KidsClubTnc         null.String `json:"kidsClubTnc" db:"KIDS_CLUB_TNC" swaggertype:"string"`
    KidsClubPartnerLink null.String `json:"kidsClubPartnerLink" db:"PARTNERS_LINK" swaggertype:"string"`
}

type GoldenPearlAboutUs struct {
    GoldenClubId      null.Int64  `json:"golden_club_id" db:"GOLDEN_CLUB_ID" swaggertype:"integer"`
    GoldenClubTitle   null.String `json:"goldenClubTitle" db:"GOLDEN_CLUB_TITLE" swaggertype:"string"`
    GoldenClubDesc    null.String `json:"goldenClubDesc" db:"GOLDEN_CLUB_DESC" swaggertype:"string"`
    GoldenClubImage   null.String `json:"goldenClubImage" db:"GOLDEN_CLUB_IMG" swaggertype:"string"`
    GoldenClubTnc     null.String `json:"goldenClubTnc" db:"GOLDEN_CLUB_TNC" swaggertype:"string"`
    GoldenClubExtLink null.String `json:"goldenClubExtLink" db:"EXTERNAL_LINK" swaggertype:"string"`
}

type LittleExplorersKidsActvParticipation struct {
    KidsActvParticipationId int    `json:"kids_actv_participation_id"`
    KidsActivityId          int    `json:"kidsActivityId"`
    KidsMembershipId        int    `json:"kidsMembershipId"`
    ActivityDateTime        string `json:"activityDateTime"`
}

type GoldenPearlActvParticipation struct {
    GoldenActvParticipationId int    `json:"golden_actv_participation_id"`
    GoldenActivityId          int    `json:"goldenActivityId"`
    GoldenMembershipId        int    `json:"goldenMembershipId"`
    ActivityDateTime          string `json:"activityDateTime"`
}

type LittleExplorersKidsActivity struct {
    KidsActivityId              null.Int64  `json:"kids_activity_id" db:"KIDS_ACTIVITY_ID" swaggertype:"integer"`
    KidsActivityCode            null.String `json:"kidsActivityCode" db:"KIDS_ACTIVITY_CODE" swaggertype:"string"`
    KidsActivityName            null.String `json:"kidsActivityName" db:"KIDS_ACTIVITY_NAME" swaggertype:"string"`
    KidsActivityDesc            null.String `json:"kidsActivityDesc" db:"KIDS_ACTIVITY_DESC" swaggertype:"string"`
    KidsActivityImage           null.String `json:"kidsActivityImage" db:"KIDS_ACTIVITY_IMG" swaggertype:"string"`
    ActivityStartDateTime       null.String `json:"activityStartDateTime" db:"ACTIVITY_START_DATETIME" swaggertype:"string"`
    ActivityEndDateTime         null.String `json:"activityEndDateTime" db:"ACTIVITY_END_DATETIME" swaggertype:"string"`
    ActivityMaxParticipant      null.Int32  `json:"activityMaxParticipant" db:"ACTIVITY_MAX_PARTICIPANT" swaggertype:"integer"`
    ActivityTnc                 null.String `json:"activityTnc" db:"ACTIVITY_TNC" swaggertype:"string"`
    ActivityDisplayOrder        null.String `json:"activityDisplayOrder" db:"ACTIVITY_DISPLAY_ORDER" swaggertype:"string"`
    ActivityAttendees           null.Int32  `json:"activityAttendees" db:"ATTENDEES" swaggertype:"integer"`
    ActivitySeatsAvailable      int         `json:"activitySeatsAvailable"`
    ActivityEndDateTimeCalendar string      `json:"activityEndDateTimeCalendar"`

    // Excel fields - not in DB, used for export
    ActivityStartDateTimeExcel string `json:"activityStartDateTimeExcel"`
    ActivityEndDateTimeExcel   string `json:"activityEndDateTimeExcel"`
}

func (o *LittleExplorersKidsActivity) Set() {
    if !o.ActivityEndDateTime.Valid {
        o.ActivityEndDateTime = utils.NewNullString("-")
    }

    if !o.ActivityTnc.Valid {
        o.ActivityTnc = utils.NewNullString("-")
    }
}

type GoldenPearlActivity struct {
    GoldenActivityId       null.Int64  `json:"golden_activity_id" db:"GOLDEN_ACTIVITY_ID" swaggertype:"integer"`
    GoldenActivityCode     null.String `json:"goldenActivityCode" db:"GOLDEN_ACTIVITY_CODE" swaggertype:"string"`
    GoldenActivityName     null.String `json:"goldenActivityName" db:"GOLDEN_ACTIVITY_NAME" swaggertype:"string"`
    GoldenActivityDesc     null.String `json:"goldenActivityDesc" db:"GOLDEN_ACTIVITY_DESC" swaggertype:"string"`
    GoldenActivityImage    null.String `json:"goldenActivityImage" db:"GOLDEN_ACTIVITY_IMG" swaggertype:"string"`
    ActivityStartDateTime  null.String `json:"activityStartDateTime" db:"ACTIVITY_START_DATETIME" swaggertype:"string"`
    ActivityEndDateTime    null.String `json:"activityEndDateTime" db:"ACTIVITY_END_DATETIME" swaggertype:"string"`
    ActivityMaxParticipant null.Int32  `json:"activityMaxParticipant" db:"ACTIVITY_MAX_PARTICIPANT" swaggertype:"integer"`
    ActivityTnc            null.String `json:"activityTnc" db:"ACTIVITY_TNC" swaggertype:"string"`
    ActivityDisplayOrder   null.String `json:"activityDisplayOrder" db:"ACTIVITY_DISPLAY_ORDER" swaggertype:"string"`
    ActivityAttendees      null.Int32  `json:"activityAttendees" db:"ATTENDEES" swaggertype:"integer"`
    ActivitySeatsAvailable int         `json:"activitySeatsAvailable"`
    ActivityEndDateTimeCalendar string `json:"activityEndDateTimeCalendar"`

    // Excel export fields (not stored in DB)
    ActivityStartDateTimeExcel string `json:"activityStartDateTimeExcel"`
    ActivityEndDateTimeExcel   string `json:"activityEndDateTimeExcel"`
}

func (o *GoldenPearlActivity) Set() {
    if !o.ActivityEndDateTime.Valid {
        o.ActivityEndDateTime = utils.NewNullString("-")
    }

    if !o.ActivityTnc.Valid {
        o.ActivityTnc = utils.NewNullString("-")
    }
}
