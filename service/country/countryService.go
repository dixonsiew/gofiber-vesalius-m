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
    lx := make([]model.CountryTelCode, 0)
    err := s.db.SelectContext(s.ctx, &lx, `SELECT COUNTRY_NAME, TEL_CODE FROM NOVA_COUNTRY WHERE TEL_CODE IS NOT NULL ORDER BY COUNTRY_NAME`)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return lx, nil
}

func (s *CountryService) FindCountryCodeByNationality(nationality string) (string, error) {
    v := ""
    var r null.String
    err := s.db.GetContext(s.ctx, &r, `SELECT COUNTRY_CODE FROM NOVA_COUNTRY WHERE NATIONALITY = :nationality`, nationality)
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
    lx := make([]model.Country, 0)
    err := s.db.SelectContext(s.ctx, &lx, `SELECT COUNTRY_NAME, TEL_CODE, COUNTRY_CODE FROM NOVA_COUNTRY WHERE COUNTRY_NAME IS NOT NULL ORDER BY COUNTRY_NAME`)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return lx, nil
}

func (s *CountryService) FindAllNationalities() ([]model.Nationality, error) {
    lx := make([]model.Nationality, 0)
    err := s.db.SelectContext(s.ctx, &lx, `SELECT NATIONALITY FROM NOVA_COUNTRY WHERE NATIONALITY IS NOT NULL ORDER BY NATIONALITY`)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return lx, nil
}
