package model

import (
    "database/sql"
)

type DbGoldenPearlAboutUs struct {
    GoldenClubID      sql.NullInt64  `db:"GOLDEN_CLUB_ID"`
    GoldenClubTitle   sql.NullString `db:"GOLDEN_CLUB_TITLE"`
    GoldenClubDesc    sql.NullString `db:"GOLDEN_CLUB_DESC"`
    GoldenClubImage   sql.NullString `db:"GOLDEN_CLUB_IMG"`
    GoldenClubTnc     sql.NullString `db:"GOLDEN_CLUB_TNC"`
    GoldenClubExtLink sql.NullString `db:"EXTERNAL_LINK"`
    UserCreate        sql.NullInt64  `db:"USER_CREATE"`
    DateCreate        sql.NullString `db:"DATE_CREATE"`
    UserUpdate        sql.NullInt64  `db:"USER_UPDATE"`
    DateUpdate        sql.NullString `db:"DATE_UPDATE"`
}

type GoldenPearlAboutUs struct {
    GoldenClubID      int64  `json:"golden_club_id"`
    GoldenClubTitle   string `json:"goldenClubTitle"`
    GoldenClubDesc    string `json:"goldenClubDesc"`
    GoldenClubImage   string `json:"goldenClubImage"`
    GoldenClubTnc     string `json:"goldenClubTnc"`
    GoldenClubExtLink string `json:"goldenClubExtLink"`
    UserCreate        int64  `json:"-"`
    DateCreate        string `json:"-"`
    UserUpdate        int64  `json:"-"`
    DateUpdate        string `json:"-"`
}

func (o *GoldenPearlAboutUs) FromDbModel(m DbGoldenPearlAboutUs) {
    o.GoldenClubID = m.GoldenClubID.Int64
    o.GoldenClubTitle = m.GoldenClubTitle.String
    o.GoldenClubDesc = m.GoldenClubDesc.String
    o.GoldenClubImage = m.GoldenClubImage.String
    o.GoldenClubTnc = m.GoldenClubTnc.String
    o.GoldenClubExtLink = m.GoldenClubExtLink.String
}
