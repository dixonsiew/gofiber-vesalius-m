package model

import (
    "github.com/guregu/null/v6"
)

type GoldenPearlAboutUs struct {
    GoldenClubID      null.Int64  `json:"golden_club_id" db:"GOLDEN_CLUB_ID" swaggertype:"integer"`
    GoldenClubTitle   null.String `json:"goldenClubTitle" db:"GOLDEN_CLUB_TITLE" swaggertype:"string"`
    GoldenClubDesc    null.String `json:"goldenClubDesc" db:"GOLDEN_CLUB_DESC" swaggertype:"string"`
    GoldenClubImage   null.String `json:"goldenClubImage" db:"GOLDEN_CLUB_IMG" swaggertype:"string"`
    GoldenClubTnc     null.String `json:"goldenClubTnc" db:"GOLDEN_CLUB_TNC" swaggertype:"string"`
    GoldenClubExtLink null.String `json:"goldenClubExtLink" db:"EXTERNAL_LINK" swaggertype:"string"`
    UserCreate        null.Int64  `json:"-" db:"USER_CREATE" swaggertype:"integer"`
    DateCreate        null.String `json:"-" db:"DATE_CREATE" swaggertype:"string"`
    UserUpdate        null.Int64  `json:"-" db:"USER_UPDATE" swaggertype:"integer"`
    DateUpdate        null.String `json:"-" db:"DATE_UPDATE" swaggertype:"string"`
}
