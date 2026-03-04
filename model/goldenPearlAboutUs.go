package model

import (
    "github.com/guregu/null/v6"
)

type GoldenPearlAboutUs struct {
    GoldenClubID      null.Int64  `json:"golden_club_id" db:"GOLDEN_CLUB_ID"`
    GoldenClubTitle   null.String `json:"goldenClubTitle" db:"GOLDEN_CLUB_TITLE"`
    GoldenClubDesc    null.String `json:"goldenClubDesc" db:"GOLDEN_CLUB_DESC"`
    GoldenClubImage   null.String `json:"goldenClubImage" db:"GOLDEN_CLUB_IMG"`
    GoldenClubTnc     null.String `json:"goldenClubTnc" db:"GOLDEN_CLUB_TNC"`
    GoldenClubExtLink null.String `json:"goldenClubExtLink" db:"EXTERNAL_LINK"`
    UserCreate        null.Int64  `json:"-" db:"USER_CREATE"`
    DateCreate        null.String `json:"-" db:"DATE_CREATE"`
    UserUpdate        null.Int64  `json:"-" db:"USER_UPDATE"`
    DateUpdate        null.String `json:"-" db:"DATE_UPDATE"`
}
