package clubs

import (
    "github.com/guregu/null/v6"
)

type LittleExplorersKidsMyActivity struct {
    KidsName             null.String `json:"kidsName" db:"KIDS_NAME"`
    KidsMembershipNumber null.String `json:"kidsMembershipNumber" db:"KIDS_MEMBERSHIP_NUMBER"`
    KidsActivityName     null.String `json:"kidsActivityName" db:"KIDS_ACTIVITY_NAME"`
    ActivityDateTime     null.String `json:"activityDateTime" db:"ACTIVITY_DATE_TIME"`
}

type GoldenPearlMyActivity struct {
    GoldenName             null.String `json:"goldenName" db:"GOLDEN_NAME"`
    GoldenMembershipNumber null.String `json:"goldenMembershipNumber" db:"GOLDEN_MEMBERSHIP_NUMBER"`
    GoldenActivityName     null.String `json:"goldenActivityName" db:"GOLDEN_ACTIVITY_NAME"`
    ActivityDateTime       null.String `json:"activityDateTime" db:"ACTIVITY_DATE_TIME"`
}

type LittleExplorersKidsMembership struct {
    KidsMembershipID       null.Int64  `json:"kidsMembershipID" db:"KIDS_MEMBERSHIP_ID"`
    KidsMembershipNumber   null.String `json:"kidsMembershipNumber" db:"KIDS_MEMBERSHIP_NUMBER"`
    KidsMembershipJoinDate null.String `json:"kidsMembershipJoinDate" db:"KIDS_MEMBERSHIP_JOIN_DATE"`
    KidsPrn                null.String `json:"kidsPrn" db:"KIDS_PRN"`
    KidsName               null.String `json:"kidsName" db:"KIDS_NAME"`
    KidsDob                null.String `json:"kidsDob" db:"KIDS_DOB"`
    KidsDocType            null.String `json:"kidsDocType" db:"KIDS_DOC_TYPE"`
    KidsDocNumber          null.String `json:"kidsDocNumber" db:"KIDS_DOC_NUMBER"`
    KidsGender             null.String `json:"kidsGender" db:"KIDS_GENDER"`
    KidsNationality        null.String `json:"kidsNationality" db:"KIDS_NATIONALITY"`
    KidsEmail              null.String `json:"kidsEmail" db:"KIDS_EMAIL"`
    GuardianPrn            null.String `json:"guardianPrn" db:"GUARDIAN_PRN"`
    GuardianName           null.String `json:"guardianName" db:"GUARDIAN_NAME"`
    GuardianDob            null.String `json:"guardianDob" db:"GUARDIAN_DOB"`
    GuardianDocType        null.String `json:"guardianDocType" db:"GUARDIAN_DOC_TYPE"`
    GuardianDocNumber      null.String `json:"guardianDocNumber" db:"GUARDIAN_DOC_NUMBER"`
    GuardianGender         null.String `json:"guardianGender" db:"GUARDIAN_GENDER"`
    GuardianNationality    null.String `json:"guardianNationality" db:"GUARDIAN_NATIONALITY"`
    GuardianEmail          null.String `json:"guardianEmail" db:"GUARDIAN_EMAIL"`
    GuardianHomeContact    null.String `json:"guardianHomeContact" db:"GUARDIAN_HOME_CONTACT"`
    GuardianMobileContact  null.String `json:"guardianMobileContact" db:"GUARDIAN_MOBILE_CONTACT"`
    GuardianAddress1       null.String `json:"guardianAddress1" db:"GUARDIAN_ADDRESS1"`
    GuardianAddress2       null.String `json:"guardianAddress2" db:"GUARDIAN_ADDRESS2"`
    GuardianAddress3       null.String `json:"guardianAddress3" db:"GUARDIAN_ADDRESS3"`
    GuardianPostCode       null.String `json:"guardianPostCode" db:"GUARDIAN_POSTCODE"`
    GuardianState          null.String `json:"guardianState" db:"GUARDIAN_STATE"`
    GuardianCountryCode    null.String `json:"guardianCountryCode" db:"GUARDIAN_COUNTRY_CODE"`
    Relationship           null.String `json:"relationship" db:"RELATIONSHIP"`
    PreferredLanguage      null.String `json:"preferredLanguage" db:"PREFERRED_LANGUAGE"`
    IsActive               null.String `json:"isActive" db:"IS_ACTIVE"`

    // Fields used in some contexts (e.g., attendee export)
    ActivityJoinDate null.String `json:"activityJoinDate" db:"ACTIVITY_DATE_TIME"`

    // Excel-specific fields (not stored in DB)
    ActivityJoinDateExcel       string `json:"activityJoinDateExcel"`
    KidsMembershipJoinDateExcel string `json:"kidsMembershipJoinDateExcel"`
    KidsDobExcel                string `json:"kidsDobExcel"`
    GuardianDobExcel            string `json:"guardianDobExcel"`
}

func (o *LittleExplorersKidsMembership) SetWebAdmin() {
    if !o.KidsPrn.Valid {
        o.KidsPrn.String = "-"
    }

    if !o.GuardianPrn.Valid {
        o.GuardianPrn.String = "-"
    }
}

func (o *LittleExplorersKidsMembership) Set() {
    if !o.KidsPrn.Valid {
        o.KidsPrn.String = "-"
    }

    if !o.KidsDob.Valid {
        o.KidsDob.String = "-"
    }

    if !o.KidsEmail.Valid {
        o.KidsEmail.String = "-"
    }

    if !o.GuardianDob.Valid {
        o.GuardianDob.String = "-"
    }

    if !o.GuardianHomeContact.Valid {
        o.GuardianHomeContact.String = "-"
    }

    if !o.GuardianMobileContact.Valid {
        o.GuardianMobileContact.String = "-"
    }

    if !o.GuardianAddress1.Valid {
        o.GuardianAddress1.String = "-"
    }

    if !o.GuardianAddress2.Valid {
        o.GuardianAddress2.String = "-"
    }

    if !o.GuardianAddress3.Valid {
        o.GuardianAddress3.String = "-"
    }

    if !o.GuardianPostCode.Valid {
        o.GuardianPostCode.String = "-"
    }

    if !o.GuardianState.Valid {
        o.GuardianState.String = "-"
    }

    if !o.GuardianCountryCode.Valid {
        o.GuardianCountryCode.String = "-"
    }
}

func (o *LittleExplorersKidsMembership) SetAttendees() {
    if !o.KidsPrn.Valid {
        o.KidsPrn.String = "-"
    }

    if !o.KidsDob.Valid {
        o.KidsDob.String = "-"
    }

    if !o.KidsEmail.Valid {
        o.KidsEmail.String = "-"
    }

    if !o.GuardianPrn.Valid {
        o.GuardianPrn.String = "-"
    }

    if !o.GuardianHomeContact.Valid {
        o.GuardianHomeContact.String = "-"
    }

    if !o.GuardianMobileContact.Valid {
        o.GuardianMobileContact.String = "-"
    }

    if !o.GuardianAddress1.Valid {
        o.GuardianAddress1.String = "-"
    }

    if !o.GuardianAddress2.Valid {
        o.GuardianAddress2.String = "-"
    }

    if !o.GuardianAddress3.Valid {
        o.GuardianAddress3.String = "-"
    }

    if !o.GuardianPostCode.Valid {
        o.GuardianPostCode.String = "-"
    }

    if !o.GuardianState.Valid {
        o.GuardianState.String = "-"
    }

    if !o.GuardianCountryCode.Valid {
        o.GuardianCountryCode.String = "-"
    }
}

type GoldenPearlMembership struct {
    GoldenMembershipID       null.Int64  `json:"golden_membership_id" db:"GOLDEN_MEMBERSHIP_ID"`
    GoldenMembershipNumber   null.String `json:"goldenMembershipNumber" db:"GOLDEN_MEMBERSHIP_NUMBER"`
    GoldenMembershipJoinDate null.String `json:"goldenMembershipJoinDate" db:"GOLDEN_MEMBERSHIP_JOIN_DATE"`
    GoldenPrn                null.String `json:"goldenPrn" db:"GOLDEN_PRN"`
    GoldenName               null.String `json:"goldenName" db:"GOLDEN_NAME"`
    GoldenDob                null.String `json:"goldenDob" db:"GOLDEN_DOB"`
    GoldenDocType            null.String `json:"goldenDocType" db:"GOLDEN_DOC_TYPE"`
    GoldenDocNumber          null.String `json:"goldenDocNumber" db:"GOLDEN_DOC_NUMBER"`
    GoldenGender             null.String `json:"goldenGender" db:"GOLDEN_GENDER"`
    GoldenNationality        null.String `json:"goldenNationality" db:"GOLDEN_NATIONALITY"`
    GoldenEmail              null.String `json:"goldenEmail" db:"GOLDEN_EMAIL"`
    NokPrn                   null.String `json:"nokPrn" db:"NOK_PRN"`
    NokName                  null.String `json:"nokName" db:"NOK_NAME"`
    NokDob                   null.String `json:"nokDob" db:"NOK_DOB"`
    NokDocType               null.String `json:"nokDocType" db:"NOK_DOC_TYPE"`
    NokDocNumber             null.String `json:"nokDocNumber" db:"NOK_DOC_NUMBER"`
    NokGender                null.String `json:"nokGender" db:"NOK_GENDER"`
    NokNationality           null.String `json:"nokNationality" db:"NOK_NATIONALITY"`
    NokEmail                 null.String `json:"nokEmail" db:"NOK_EMAIL"`
    NokHomeContact           null.String `json:"nokHomeContact" db:"NOK_HOME_CONTACT"`
    NokMobileContact         null.String `json:"nokMobileContact" db:"NOK_MOBILE_CONTACT"`
    NokAddress1              null.String `json:"nokAddress1" db:"NOK_ADDRESS1"`
    NokAddress2              null.String `json:"nokAddress2" db:"NOK_ADDRESS2"`
    NokAddress3              null.String `json:"nokAddress3" db:"NOK_ADDRESS3"`
    NokPostCode              null.String `json:"nokPostCode" db:"NOK_POSTCODE"`
    NokState                 null.String `json:"nokState" db:"NOK_STATE"`
    NokCountryCode           null.String `json:"nokCountryCode" db:"NOK_COUNTRY_CODE"`
    Relationship             null.String `json:"relationship" db:"RELATIONSHIP"`
    PreferredLanguage        null.String `json:"preferredLanguage" db:"PREFERRED_LANGUAGE"`
    IsActive                 null.String `json:"isActive" db:"IS_ACTIVE"`

    // Fields used in attendee contexts
    ActivityJoinDate null.String `json:"activityJoinDate" db:"ACTIVITY_DATE_TIME"`

    // Excel-specific fields (not stored in DB)
    ActivityJoinDateExcel         string `json:"activityJoinDateExcel"`
    GoldenMembershipJoinDateExcel string `json:"goldenMembershipJoinDateExcel"`
    GoldenDobExcel                string `json:"goldenDobExcel"`
    NokDobExcel                   string `json:"nokDobExcel"`
}

func (o *GoldenPearlMembership) Set() {
    if !o.GoldenPrn.Valid {
        o.GoldenPrn.String = "-"
    }

    if !o.NokPrn.Valid {
        o.NokPrn.String = "-"
    }
}

func (o *GoldenPearlMembership) Set() {
    if !o.GoldenPrn.Valid {
        o.GoldenPrn.String = "-"
    }

    if !o.GoldenDob.Valid {
        o.GoldenDob.String = "-"
    }

    if !o.GoldenEmail.Valid {
        o.GoldenEmail.String = "-"
    }

    if !o.NokPrn.Valid {
        o.NokPrn.String = "-"
    }

    if !o.NokDob.Valid {
        o.NokDob.String = "-"
    }

    if !o.NokHomeContact.Valid {
        o.GoldenDob.String = "-"
    }

    if !o.NokMobileContact.Valid {
        o.GoldenDob.String = "-"
    }

    if !o.NokAddress1.Valid {
        o.NokAddress1.String = "-"
    }

    if !o.NokAddress2.Valid {
        o.NokAddress2.String = "-"
    }

    if !o.NokAddress3.Valid {
        o.NokAddress3.String = "-"
    }

    if !o.NokPostCode.Valid {
        o.NokPostCode.String = "-"
    }

    if !o.NokState.Valid {
        o.NokState.String = "-"
    }

    if !o.NokCountryCode.Valid {
        o.NokCountryCode.String = "-"
    }
}
