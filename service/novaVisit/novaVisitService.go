package novaVisit

import (
	"context"
	"database/sql"
	"strings"
	"vesaliusm/database"
	model "vesaliusm/model/healthCare"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
)

var NovaVisitSvc *NovaVisitService = NewNovaVisitService(database.GetDbrs(), database.GetCtx())

type NovaVisitService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaVisitService(db *sqlx.DB, ctx context.Context) *NovaVisitService {
    return &NovaVisitService{db: db, ctx: ctx}
}

func (s *NovaVisitService) FindByPrn(prn string, conn *sqlx.DB) ([]model.NovaVisit, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `SELECT * FROM NOVA_VISIT WHERE PRN = :prn`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaVisit{}, ""), 1)
    list := make([]model.NovaVisit, 0)
    err := db.SelectContext(s.ctx, &list, query, prn)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVisitService) GetSpecificPatientVisit(prn string, conn *sqlx.DB) ([]model.NovaVisit, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `
        SELECT * FROM NOVA_VISIT WHERE (PRN = :prn OR PRN IN (SELECT OLD_PRN FROM NOVA_PATIENT_MERGE_DETAIL WHERE NEW_PRN = :prn)) 
        AND ((VISIT_TYPE IN ('DAY-SURGERY', 'INPATIENT') AND VISIT_STATUS <> 'CANCEL' AND ADMISSION_STATUS = 'DISCHARGED') 
        OR (VISIT_TYPE IN ('OUTPATIENT', 'EXTERNAL', 'ACCIDENT-EMERGENCY') AND VISIT_STATUS <> 'CANCEL')) 
        ORDER BY TO_DATE(REGISTRATION_DATE, 'DD/MM/YYYY HH24:MI:SS') DESC`
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaVisit{}, ""), 1)
    list := make([]model.NovaVisit, 0)
    err := db.SelectContext(s.ctx, &list, query, sql.Named("prn", prn))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVisitService) GetPatientVisitWithVitalSign(prn string, conn *sqlx.DB) ([]model.NovaVisit, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `
        SELECT * FROM ( 
          SELECT v.* 
          FROM NOVA_PATIENT_VITALS p, NOVA_VISIT v 
          WHERE p.ACCOUNT_NO = v.ACCOUNT_NO 
          AND (v.PRN = :prn OR v.PRN IN (SELECT OLD_PRN FROM NOVA_PATIENT_MERGE_DETAIL WHERE NEW_PRN = :prn)) 
          AND ((v.VISIT_TYPE IN ('DAY-SURGERY', 'INPATIENT') AND v.VISIT_STATUS <> 'CANCEL') 
          OR (v.VISIT_TYPE IN ('OUTPATIENT', 'EXTERNAL', 'ACCIDENT-EMERGENCY') AND v.VISIT_STATUS <> 'CANCEL')) 
          AND p.STATUS <> 'VOID' 
          UNION 
          SELECT v.* 
          FROM NOVA_VISIT v, NOVA_PATIENT_TEMPLATE t 
          WHERE t.ACCOUNT_NO = v.ACCOUNT_NO 
          AND (v.PRN = :prn OR v.PRN IN (SELECT OLD_PRN FROM NOVA_PATIENT_MERGE_DETAIL WHERE NEW_PRN = :prn)) 
          AND ((v.VISIT_TYPE IN ('DAY-SURGERY', 'INPATIENT') AND v.VISIT_STATUS <> 'CANCEL') 
          OR (v.VISIT_TYPE IN ('OUTPATIENT', 'EXTERNAL', 'ACCIDENT-EMERGENCY') AND v.VISIT_STATUS <> 'CANCEL')) 
          AND t.TEMPLATE_REF_NO IN ( 
            SELECT TEMPLATE_REF_NO FROM NOVA_PATIENT_TEMPLATE_DETAIL 
            WHERE TEMPLATE_REF_NO IN (SELECT TEMPLATE_REF_NO FROM NOVA_PATIENT_TEMPLATE WHERE TEMPLATE_NAME = 'PHYSICAL EXAMINATION') 
            AND SECTION_NAME = 'GENERAL CONDITION AND VITAL SIGNS' 
          )
        ) 
        ORDER BY TO_DATE(REGISTRATION_DATE, 'DD/MM/YYYY HH24:MI:SS') DESC 
        FETCH FIRST 1 ROW ONLY 
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaVisit{}, ""), 1)
    list := make([]model.NovaVisit, 0)
    err := db.SelectContext(s.ctx, &list, query, sql.Named("prn", prn))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVisitService) GetPatientVisitWithReferralLetter(prn string, conn *sqlx.DB) ([]model.NovaVisit, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `
        SELECT * FROM ( 
          SELECT v.* FROM NOVA_VISIT v, NOVA_EXT_REFERRAL_LETTER e 
          WHERE (v.PRN = :prn OR v.PRN IN (SELECT OLD_PRN FROM NOVA_PATIENT_MERGE_DETAIL WHERE NEW_PRN = :prn)) 
          AND v.PRN = e.PRN 
          AND v.ACCOUNT_NO = e.ACCOUNT_NO 
          AND ((v.VISIT_TYPE IN ('DAY-SURGERY', 'INPATIENT') AND v.VISIT_STATUS NOT IN ('OPEN', 'CANCEL')) 
          OR (v.VISIT_TYPE IN ('OUTPATIENT', 'EXTERNAL', 'ACCIDENT-EMERGENCY') AND v.VISIT_STATUS NOT IN ('OPEN', 'CANCEL'))) 
          UNION 
          SELECT v.* FROM NOVA_VISIT v, NOVA_INT_REFERRAL_LETTER e 
          WHERE (v.PRN = :prn OR v.PRN IN (SELECT OLD_PRN FROM NOVA_PATIENT_MERGE_DETAIL WHERE NEW_PRN = :prn)) 
          AND v.PRN = e.PRN 
          AND v.ACCOUNT_NO = e.ACCOUNT_NO 
          AND ((v.VISIT_TYPE IN ('DAY-SURGERY', 'INPATIENT') AND v.VISIT_STATUS NOT IN ('OPEN', 'CANCEL')) 
          OR (v.VISIT_TYPE IN ('OUTPATIENT', 'EXTERNAL', 'ACCIDENT-EMERGENCY') AND v.VISIT_STATUS NOT IN ('OPEN', 'CANCEL'))) 
        ) ORDER BY 5
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaVisit{}, ""), 1)
    list := make([]model.NovaVisit, 0)
    err := db.SelectContext(s.ctx, &list, query, sql.Named("prn", prn))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVisitService) GetPatientVisitWithHealthScreeningReport(prn string, conn *sqlx.DB) ([]model.NovaVisit, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `
        SELECT v.* FROM NOVA_VISIT v, NOVA_ACCOUNT_HSR r, NOVA_ACCOUNT_HSR_CLOB c 
        WHERE (v.PRN = :prn OR v.PRN IN (SELECT OLD_PRN FROM NOVA_PATIENT_MERGE_DETAIL WHERE NEW_PRN = :prn)) 
        AND v.ACCOUNT_NO = r.ACCOUNT_NO 
        AND ((v.VISIT_TYPE IN ('DAY-SURGERY', 'INPATIENT') AND v.VISIT_STATUS <> 'CANCEL') 
        OR (v.VISIT_TYPE IN ('OUTPATIENT', 'EXTERNAL', 'ACCIDENT-EMERGENCY') AND v.VISIT_STATUS <> 'CANCEL'))
        AND r.HSR_REF_NO = c.HSR_REF_NO 
        AND c.HSR_TYPE = 'HEALTH-SCREENING-REPORT'
    `
    query = strings.Replace(query, "v.*", utils.GetDbCols(model.NovaVisit{}, ""), 1)
    list := make([]model.NovaVisit, 0)
    err := db.SelectContext(s.ctx, &list, query, sql.Named("prn", prn))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
