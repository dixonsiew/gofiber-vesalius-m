package country

import (
	"context"
	"database/sql"
	"vesaliusm/model"
	"vesaliusm/utils"

	"github.com/guregu/null/v6"
	"github.com/jmoiron/sqlx"
)

type CountryService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewCountryService(db *sqlx.DB, ctx context.Context) *CountryService {
    return &CountryService{db: db, ctx: ctx}
}

func (s *CountryService) FindAllCountryTelCode() ([]model.CountryTelCode, error) {
    query := `SELECT COUNTRY_NAME, TEL_CODE FROM NOVA_COUNTRY WHERE TEL_CODE IS NOT NULL ORDER BY COUNTRY_NAME`
    list := make([]model.CountryTelCode, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *CountryService) FindCountryCodeByNationality(nationality string) (string, error) {
    query := `SELECT COUNTRY_CODE FROM NOVA_COUNTRY WHERE NATIONALITY = :nationality`
    v := ""
    var r null.String
    err := s.db.GetContext(s.ctx, &r, query, nationality)
    if err != nil {
        if err == sql.ErrNoRows {
            return v, err
        }
        utils.LogError(err)
        return v, err
    }
    v = r.String
    return v, nil
}

func (s *CountryService) FindAllCountries() ([]model.Country, error) {
    query := `SELECT COUNTRY_NAME, TEL_CODE, COUNTRY_CODE FROM NOVA_COUNTRY WHERE COUNTRY_NAME IS NOT NULL ORDER BY COUNTRY_NAME`
    list := make([]model.Country, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *CountryService) FindAllNationalities() ([]model.Nationality, error) {
    query := `SELECT NATIONALITY FROM NOVA_COUNTRY WHERE NATIONALITY IS NOT NULL ORDER BY NATIONALITY`
    list := make([]model.Nationality, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
