package clubs

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model/clubs"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/applicationUserFamily"
    "vesaliusm/service/vesaliusGeo"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    "github.com/nleeper/goment"
)

var ClubSvc *ClubService = NewClubService(database.GetDb(), database.GetCtx(), database.GetDbrs(), database.GetCtxrs())

type ClubService struct {
    db                           *sqlx.DB
    dbrs                         *sqlx.DB
    ctx                          context.Context
    ctxrs                        context.Context
    applicationUserService       *applicationUser.ApplicationUserService
    applicationUserFamilyService *applicationUserFamily.ApplicationUserFamilyService
    vesaliusGeoService           *vesaliusGeo.VesaliusGeoService
}

func NewClubService(db *sqlx.DB, ctx context.Context, dbrs *sqlx.DB, ctxrs context.Context) *ClubService {
    return &ClubService{
        db:                           db,
        dbrs:                         dbrs,
        ctx:                          ctx,
        ctxrs:                        ctxrs,
        applicationUserService:       applicationUser.ApplicationUserSvc,
        applicationUserFamilyService: applicationUserFamily.ApplicationUserFamilySvc,
        vesaliusGeoService:           vesaliusGeo.VesaliusGeoSvc,
    }
}

func (s *ClubService) FindGuestLittleKidsMembershipByIc(identificationNumber string) (*clubs.LittleExplorersKidsMembership, error) {
    query := `SELECT * FROM KIDS_CLUB_MEMBERSHIP WHERE KIDS_DOC_NUMBER = :identificationNumber`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.LittleExplorersKidsMembership{}, ""), 1)
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
        o.KidsMembershipJoinDate = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
    }
    if o.KidsDob.Valid {
        g, _ := goment.New(o.KidsDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.KidsDob = utils.NewNullString(g.Format("DD/MM/YYYY"))
    }
    return &o, err
}

func (s *ClubService) FindGuestGoldenPearlMembershipByIc(identificationNumber string) (*clubs.GoldenPearlMembership, error) {
    query := `SELECT * FROM GOLDEN_CLUB_MEMBERSHIP WHERE GOLDEN_DOC_NUMBER = :identificationNumber`
    query = strings.Replace(query, "*", utils.GetDbCols(clubs.GoldenPearlMembership{}, ""), 1)
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
        o.GoldenMembershipJoinDate = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
    }
    if o.GoldenDob.Valid {
        g, _ := goment.New(o.GoldenDob.String, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.GoldenDob = utils.NewNullString(g.Format("DD/MM/YYYY"))
    }
    return &o, nil
}

func whereClause(conds []string) string {
    if len(conds) == 0 {
        return ""
    }
    return " WHERE " + strings.Join(conds, " AND ")
}
