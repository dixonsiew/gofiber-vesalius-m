package dto

type LittleExplorersKidsMembershipDto struct {
    KidsName              string `json:"kidsName" validate:"required"`
    KidsDob               string `json:"kidsDob" validate:"required"`
    KidsDocType           string `json:"kidsDocType" validate:"required"`
    KidsDocNumber         string `json:"kidsDocNumber" validate:"required"`
    KidsGender            string `json:"kidsGender" validate:"required"`
    KidsNationality       string `json:"kidsNationality" validate:"required"`
    KidsEmail             string `json:"kidsEmail"`
    GuardianName          string `json:"guardianName" validate:"required"`
    GuardianDob           string `json:"guardianDob" validate:"required"`
    GuardianDocType       string `json:"guardianDocType" validate:"required"`
    GuardianDocNumber     string `json:"guardianDocNumber" validate:"required"`
    GuardianGender        string `json:"guardianGender" validate:"required"`
    GuardianNationality   string `json:"guardianNationality" validate:"required"`
    GuardianEmail         string `json:"guardianEmail" validate:"required,email"`
    GuardianHomeContact   string `json:"guardianHomeContact"`
    GuardianMobileContact string `json:"guardianMobileContact"`
    GuardianAddress1      string `json:"guardianAddress1"`
    GuardianAddress2      string `json:"guardianAddress2"`
    GuardianAddress3      string `json:"guardianAddress3"`
    GuardianPostCode      string `json:"guardianPostCode"`
    GuardianState         string `json:"guardianState"`
    GuardianCountryCode   string `json:"guardianCountryCode"`
    Relationship          string `json:"relationship" validate:"required"`
    PreferredLanguage     string `json:"preferredLanguage" validate:"required"`
}

type LittleExplorersKidsAboutUsDto struct {
    KidsClubTitle       string `json:"kidsClubTitle" validate:"required"`
    KidsClubDesc        string `json:"kidsClubDesc" validate:"required"`
    KidsClubImage       string `json:"kidsClubImage"`
    KidsClubTnc         string `json:"kidsClubTnc"`
    KidsClubPartnerLink string `json:"kidsClubPartnerLink"`
}

type KidsActvParticipationDto struct {
    KidsActivityId   int    `json:"kidsActivityId" validate:"required,numeric"`
    KidsMembershipId int    `json:"kidsMembershipId" validate:"required,numeric"`
    ActivityDateTime string `json:"activityDateTime" validate:"required"`
}

type LittleExplorersKidsActvParticipationDto struct {
    KidsActvParticipation []KidsActvParticipationDto `json:"kidsActvParticipation" validate:"required"`
}

type LittleExplorersKidsActivityDto struct {
    KidsActivityCode        string `json:"kidsActivityCode" validate:"required"`
    KidsActivityName        string `json:"kidsActivityName" validate:"required"`
    KidsActivityDesc        string `json:"kidsActivityDesc" validate:"required"`
    KidsActivityImage       string `json:"kidsActivityImage"`
    ActivityStartDateTime   string `json:"activityStartDateTime" validate:"required"`
    ActivityEndDateTime     string `json:"activityEndDateTime"`
    ActivityMaxParticipant  string `json:"activityMaxParticipant" validate:"required"`
    ActivityTnc             string `json:"activityTnc"`
    ActivityDisplayOrder    string `json:"activityDisplayOrder" validate:"required"`
}
