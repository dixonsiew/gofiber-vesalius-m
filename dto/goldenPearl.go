package dto

type GoldenPearlMembershipDto struct {
    GoldenName        string `json:"goldenName"        validate:"required"`
    GoldenDob         string `json:"goldenDob"         validate:"required"`
    GoldenDocType     string `json:"goldenDocType"     validate:"required"`
    GoldenDocNumber   string `json:"goldenDocNumber"   validate:"required"`
    GoldenGender      string `json:"goldenGender"      validate:"required"`
    GoldenNationality string `json:"goldenNationality" validate:"required"`
    GoldenEmail       string `json:"goldenEmail"`

    NokName        string `json:"nokName"           validate:"required"`
    NokDob         string `json:"nokDob"            validate:"required"`
    NokDocType     string `json:"nokDocType"        validate:"required"`
    NokDocNumber   string `json:"nokDocNumber"      validate:"required"`
    NokGender      string `json:"nokGender"         validate:"required"`
    NokNationality string `json:"nokNationality"    validate:"required"`
    NokEmail       string `json:"nokEmail"          validate:"required"`

    NokHomeContact   string `json:"nokHomeContact"`
    NokMobileContact string `json:"nokMobileContact"`
    NokAddress1      string `json:"nokAddress1"`
    NokAddress2      string `json:"nokAddress2"`
    NokAddress3      string `json:"nokAddress3"`
    NokPostCode      string `json:"nokPostCode"`
    NokState         string `json:"nokState"`
    NokCountryCode   string `json:"nokCountryCode"`

    Relationship      string `json:"relationship"      validate:"required"`
    PreferredLanguage string `json:"preferredLanguage" validate:"required"`
}

type GoldenPearlAboutUsDto struct {
    GoldenClubTitle   string `json:"goldenClubTitle"   validate:"required"`
    GoldenClubDesc    string `json:"goldenClubDesc"    validate:"required"`
    GoldenClubImage   string `json:"goldenClubImage"`
    GoldenClubTnc     string `json:"goldenClubTnc"`
    GoldenClubExtLink string `json:"goldenClubExtLink"`
}

type GoldenActvParticipationDto struct {
    GoldenActivityId   int    `json:"goldenActivityId"   validate:"required"`
    GoldenMembershipId int    `json:"goldenMembershipId" validate:"required"`
    ActivityDateTime   string `json:"activityDateTime"   validate:"required"`
}

type GoldenPearlActvParticipationDto struct {
    GoldenActvParticipation []GoldenActvParticipationDto `json:"goldenActvParticipation"`
}

type GoldenPearlActivityDto struct {
    GoldenActivityCode     string `json:"goldenActivityCode"        validate:"required"`
    GoldenActivityName     string `json:"goldenActivityName"        validate:"required"`
    GoldenActivityDesc     string `json:"goldenActivityDesc"        validate:"required"`
    GoldenActivityImage    string `json:"goldenActivityImage"`
    ActivityStartDateTime  string `json:"activityStartDateTime"     validate:"required"`
    ActivityEndDateTime    string `json:"activityEndDateTime"`
    ActivityMaxParticipant string `json:"activityMaxParticipant"    validate:"required"`
    ActivityTnc            string `json:"activityTnc"`
    ActivityDisplayOrder   string `json:"activityDisplayOrder"      validate:"required"`
}
