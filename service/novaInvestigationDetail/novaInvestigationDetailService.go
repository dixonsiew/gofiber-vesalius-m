package novaInvestigationDetail

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    model "vesaliusm/model/healthCare"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var NovaInvestigationDetailSvc *NovaInvestigationDetailService = NewNovaInvestigationDetailService(database.GetDbrs(), database.GetCtx())

type NovaInvestigationDetailService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaInvestigationDetailService(db *sqlx.DB, ctx context.Context) *NovaInvestigationDetailService {
    return &NovaInvestigationDetailService{db: db, ctx: ctx}
}

func (s *NovaInvestigationDetailService) GetInvestigationRefNoAndCode(investigationRefNo string, code string, conn *sqlx.DB) ([]model.NovaPatientInvestigationDetail, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT dbms_random.string('U', 30) as id, i.* 
        FROM NOVA_PATIENT_INVESTIGATION_DET i 
        WHERE i.INVESTIGATION_REF_NO = :investigationRefNo 
        AND i.CODE = :code 
        ORDER BY i.CODE, TO_NUMBER(i.INVESTIGATION_REF_NO) ASC
    `
    query = strings.Replace(query, "i.*", utils.GetDbCols(model.NovaPatientInvestigationDetail{}, "i."), 1)
    list := make([]model.NovaPatientInvestigationDetail, 0)
    err := db.SelectContext(s.ctx, &list, query, investigationRefNo, code)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaInvestigationDetailService) GetLabHistoryTrendingForDashboard(
    prn string, labInvestigationType string, labCodeHDL string, labCodeLDL string, labCodeGlucose string, 
    labCodeHemoglobin string, conn *sqlx.DB) ([]model.NovaPatientInvestigationDetail, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        WITH ranked_results AS ( 
          SELECT i.STATUS_DATE_TIME, 
          d.INVESTIGATION_REF_NO, 
          d.CODE, 
          d.RESULT_VALUE, 
          d.RESULT_UNIT, 
          d.REFERENCE_RANGE, 
          d.RANGE_TYPE, 
          d.RESULT_CLOB, 
          dbms_random.string('U', 30) as id, 
          ROW_NUMBER() OVER (PARTITION BY d.CODE ORDER BY i.INVESTIGATION_REF_NO DESC) as rn 
          FROM NOVA_PATIENT_INVESTIGATION i 
          JOIN NOVA_VISIT v ON i.ACCOUNT_NO = v.ACCOUNT_NO 
          JOIN NOVA_PATIENT_INVESTIGATION_DET d ON i.INVESTIGATION_REF_NO = d.INVESTIGATION_REF_NO AND i.CODE = d.CODE 
          LEFT JOIN NOVA_PATIENT_MERGE_DETAIL m ON v.PRN = m.OLD_PRN 
          WHERE i.INVESTIGATION_TYPE = :labInvestigationType 
          AND d.RESULT_VALUE IS NOT NULL 
          AND (d.RESULT_CLOB IS NULL OR NOT REGEXP_LIKE (TO_CHAR(d.RESULT_CLOB), '[<>]')) 
          AND (v.PRN = :prn OR m.NEW_PRN = :prn) 
          AND i.CODE IN (:labCodeHDL, :labCodeLDL, :labCodeGlucose, :labCodeHemoglobin) 
              ) 
              SELECT TO_CHAR(STATUS_DATE_TIME, 'DD Mon YYYY') as SYNC_DATE, 
              INVESTIGATION_REF_NO, 
              CODE, 
              RESULT_VALUE, 
              RESULT_UNIT, 
              REFERENCE_RANGE, 
              RANGE_TYPE, 
              RESULT_CLOB, 
              id 
              FROM ranked_results 
              WHERE rn <= 5 
              ORDER BY CODE, INVESTIGATION_REF_NO ASC
    `
    list := make([]model.NovaPatientInvestigationDetail, 0)
    err := db.SelectContext(s.ctx, &list, query, 
        sql.Named("prn", prn),
        sql.Named("labInvestigationType", labInvestigationType),
        sql.Named("labCodeHDL", labCodeHDL),
        sql.Named("labCodeLDL", labCodeLDL),
        sql.Named("labCodeGlucose", labCodeGlucose),
        sql.Named("labCodeHemoglobin", labCodeHemoglobin),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
