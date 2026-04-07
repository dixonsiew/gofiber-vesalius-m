package app

import (
    "context"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    "github.com/nleeper/goment"
)

var AppSvc *AppService = NewAppService(database.GetDb(), database.GetCtx())

type AppService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewAppService(db *sqlx.DB, ctx context.Context) *AppService {
    return &AppService{
        db:  db,
        ctx: ctx,
    }
}

func (s *AppService) FindAllAppHospitalProfile() ([]model.HospitalProfile, error) {
    query := `SELECT * FROM HOSPITAL_PROFILE`
    query = strings.Replace(query, "*", utils.GetDbCols(model.HospitalProfile{}, ""), 1)
    list := make([]model.HospitalProfile, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AppService) FindAllAppVersion() ([]model.AppVersion, error) {
    query := `SELECT * FROM APP_VERSION`
    query = strings.Replace(query, "*", utils.GetDbCols(model.AppVersion{}, ""), 1)
    list := make([]model.AppVersion, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AppService) FindAllReleaseVersion() ([]model.ReleaseVersion, error) {
    query := `SELECT * FROM RELEASE_VERSION`
    query = strings.Replace(query, "*", utils.GetDbCols(model.ReleaseVersion{}, ""), 1)
    list := make([]model.ReleaseVersion, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        if list[i].DateUpdate.Valid {
            g, _ := goment.New(list[i].DateUpdate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].DateUpdate = utils.NewNullString(g.Format("DD/MM/YYYY"))
        }
    }
    return list, nil
}

func (s *AppService) GetCurrentReleaseVersion(stackPlatform string) (string, error) {
    query := `SELECT LATEST_VERSION FROM RELEASE_VERSION WHERE STACK_PLATFORM = :stackPlatform`
    var version string
    err := s.db.GetContext(s.ctx, &version, query, stackPlatform)
    if err != nil {
        utils.LogError(err)
        return "", err
    }
    return version, nil
}

func (s *AppService) UpdateReleaseVersion(latestVersion string, stackPlatform string) error {
    query := `
        UPDATE RELEASE_VERSION SET
          LATEST_VERSION = :ver,
          DATE_UPDATE = CURRENT_TIMESTAMP
        WHERE STACK_PLATFORM = :platform
    `
    _, err := s.db.ExecContext(s.ctx, query, latestVersion, stackPlatform)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *AppService) UpdateAppVersion(latestVersion string, osPlatform string, status int) error {
    query := `
        UPDATE APP_VERSION SET
          LATEST_VERSION = :ver,
          DATE_UPDATE = CURRENT_TIMESTAMP,
          STATUS = :status
        WHERE OS_PLATFORM = :platform
    `
    _, err := s.db.ExecContext(s.ctx, query, latestVersion, status, osPlatform)
    if err != nil {
        utils.LogError(err)
        return err
    }
    return nil
}

func (s *AppService) FindAllGuestMode(conn *sqlx.DB) ([]model.AppServices, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM APP_SERVICE WHERE IS_GUEST_ONLY = 1 AND IS_ENABLED = 1 ORDER BY SERVICE_DISPLAY_ORDER`
    query = strings.Replace(query, "*", utils.GetDbCols(model.AppServices{}, ""), 1)
    list := make([]model.AppServices, 0)
    err := db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AppService) ListByGuestMode() (*model.PagedList, error) {
    total, err := s.CountByGuestMode(s.db)
    if err != nil {
        return nil, err
    }

    list, err := s.FindAllGuestMode(s.db)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: 1,
    }, nil
}

func (s *AppService) CountByGuestMode(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(SERVICE_NAME) AS COUNT FROM APP_SERVICE WHERE IS_GUEST_ONLY = 1 AND IS_ENABLED = 1`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *AppService) FindAllAuthMode(conn *sqlx.DB) ([]model.AppServices, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM APP_SERVICE WHERE IS_ENABLED = 1 ORDER BY SERVICE_DISPLAY_ORDER`
    query = strings.Replace(query, "*", utils.GetDbCols(model.AppServices{}, ""), 1)
    list := make([]model.AppServices, 0)
    err := db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *AppService) ListByAuthMode() (*model.PagedList, error) {
    total, err := s.CountByAuthMode(s.db)
    if err != nil {
        return nil, err
    }

    list, err := s.FindAllAuthMode(s.db)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: 1,
    }, nil
}

func (s *AppService) CountByAuthMode(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(SERVICE_NAME) AS COUNT FROM APP_SERVICE WHERE IS_ENABLED = 1`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}
