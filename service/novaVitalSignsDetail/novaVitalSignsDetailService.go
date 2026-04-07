package novaVitalSignsDetail

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    model "vesaliusm/model/healthCare"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var NovaVitalSignsDetailSvc *NovaVitalSignsDetailService = NewNovaVitalSignsDetailService(database.GetDbrs(), database.GetCtx())

type NovaVitalSignsDetailService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaVitalSignsDetailService(db *sqlx.DB, ctx context.Context) *NovaVitalSignsDetailService {
    return &NovaVitalSignsDetailService{db: db, ctx: ctx}
}

func (s *NovaVitalSignsDetailService) GetLocalPatientVitalSignsDetailsByRefNo(
    refNo string, vitalCodeHeight string, vitalCodeWeight string, vitalCodeBP string, vitalCodeBMI string, 
    vitalCodePulse string, conn *sqlx.DB,
) ([]model.NovaPatientVitalSignsDetail, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT dbms_random.string('U', 30) as id, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, 
        to_char(SYNC_DATE, 'DD Mon YYYY') as DATE_TIME 
        FROM NOVA_PATIENT_VITALS_DET WHERE REF_NO = :refNo 
        AND (CODE IN (:vitalCodeBP, :vitalCodeBMI, :vitalCodePulse) OR (LOWER(DESCRIPTION) LIKE :vitalCodeHeight OR LOWER(DESCRIPTION) LIKE :vitalCodeWeight))
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaPatientVitalSignsDetail{}, ""), 1)
    list := make([]model.NovaPatientVitalSignsDetail, 0)
    err := db.SelectContext(s.ctx, &list, query, refNo, vitalCodeBP, vitalCodeBMI, vitalCodePulse, vitalCodeHeight, vitalCodeWeight)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVitalSignsDetailService) GetPatientVitalSignsHistoryHeight(prn string, vitalSignsCode string, conn *sqlx.DB) ([]model.NovaPatientVitalSignsDetail, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
          SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
          SELECT dbms_random.string('U', 30) as ID, d.REF_NO, d.CODE, d.DESCRIPTION, d.VALUE1, d.VALUE2, d.UNIT, to_char(vi.date_time, 'DD Mon YYYY') as DATE_TIME 
          FROM NOVA_PATIENT_VITALS_DET d, NOVA_PATIENT_VITALS vi, NOVA_VISIT v 
          WHERE LOWER(d.DESCRIPTION)= :vitalSignsCode 
          AND d.REF_NO = vi.REF_NO 
          AND vi.ACCOUNT_NO = v.ACCOUNT_NO 
          AND d.REF_NO IN (SELECT REF_NO FROM NOVA_PATIENT_VITALS WHERE STATUS <> 'VOID' 
          AND (PRN = :prn OR PRN IN (SELECT OLD_PRN FROM NOVA_PATIENT_MERGE_DETAIL WHERE NEW_PRN = :prn))) 
        ) ORDER BY 8 DESC FETCH FIRST 5 ROWS ONLY) ORDER BY TO_DATE(DATE_TIME, 'DD Mon YYYY') ASC
    `
    list := make([]model.NovaPatientVitalSignsDetail, 0)
    err := db.SelectContext(s.ctx, &list, query, sql.Named("prn", prn), sql.Named("vitalSignsCode", vitalSignsCode))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVitalSignsDetailService) GetPatientVitalSignsHistoryWeight(prn string, vitalSignsCode string, conn *sqlx.DB) ([]model.NovaPatientVitalSignsDetail, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
          SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
          SELECT dbms_random.string('U', 30) as ID, d.REF_NO, d.CODE, d.DESCRIPTION, d.VALUE1, d.VALUE2, d.UNIT, to_char(vi.date_time, 'DD Mon YYYY') as DATE_TIME 
          FROM NOVA_PATIENT_VITALS_DET d, NOVA_PATIENT_VITALS vi, NOVA_VISIT v 
          WHERE LOWER(d.DESCRIPTION) = :vitalSignsCode 
          AND d.REF_NO = vi.REF_NO 
          AND vi.ACCOUNT_NO = v.ACCOUNT_NO 
          AND d.REF_NO IN (SELECT REF_NO FROM NOVA_PATIENT_VITALS WHERE STATUS <> 'VOID' 
          AND (PRN = :prn OR PRN IN (SELECT OLD_PRN FROM NOVA_PATIENT_MERGE_DETAIL WHERE NEW_PRN = :prn))) 
        ) ORDER BY 8 DESC FETCH FIRST 5 ROWS ONLY) ORDER BY TO_DATE(DATE_TIME, 'DD Mon YYYY') ASC
    `
    list := make([]model.NovaPatientVitalSignsDetail, 0)
    err := db.SelectContext(s.ctx, &list, query, sql.Named("prn", prn), sql.Named("vitalSignsCode", vitalSignsCode))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVitalSignsDetailService) GetPatientVitalSignsHistoryBMI(prn string, vitalSignsCode string, conn *sqlx.DB) ([]model.NovaPatientVitalSignsDetail, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
          SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
          SELECT dbms_random.string('U', 30) as ID, d.REF_NO, d.CODE, d.DESCRIPTION, d.VALUE1, d.VALUE2, d.UNIT, to_char(vi.date_time, 'DD Mon YYYY') as DATE_TIME 
          FROM NOVA_PATIENT_VITALS_DET d, NOVA_PATIENT_VITALS vi, NOVA_VISIT v 
          WHERE d.CODE = :vitalSignsCode 
          AND d.REF_NO = vi.REF_NO 
          AND vi.ACCOUNT_NO = v.ACCOUNT_NO 
          AND d.REF_NO IN (SELECT REF_NO FROM NOVA_PATIENT_VITALS WHERE STATUS <> 'VOID' 
          AND (PRN = :prn OR PRN IN (SELECT OLD_PRN FROM NOVA_PATIENT_MERGE_DETAIL WHERE NEW_PRN = :prn))) 
		) ORDER BY 8 DESC FETCH FIRST 5 ROWS ONLY) ORDER BY TO_DATE(DATE_TIME, 'DD Mon YYYY') ASC
    `
    list := make([]model.NovaPatientVitalSignsDetail, 0)
    err := db.SelectContext(s.ctx, &list, query, sql.Named("prn", prn), sql.Named("vitalSignsCode", vitalSignsCode))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVitalSignsDetailService) GetPatientVitalSignsHistoryPulse(prn string, vitalSignsCode string, conn *sqlx.DB) ([]model.NovaPatientVitalSignsDetail, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
          SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
          SELECT dbms_random.string('U', 30) as ID, d.REF_NO, d.CODE, d.DESCRIPTION, d.VALUE1, d.VALUE2, d.UNIT, to_char(vi.date_time, 'DD Mon YYYY') as DATE_TIME 
          FROM NOVA_PATIENT_VITALS_DET d, NOVA_PATIENT_VITALS vi, NOVA_VISIT v 
          WHERE d.CODE = :vitalSignsCode 
          AND d.REF_NO = vi.REF_NO 
          AND vi.ACCOUNT_NO = v.ACCOUNT_NO 
          AND d.REF_NO IN (SELECT REF_NO FROM NOVA_PATIENT_VITALS WHERE STATUS <> 'VOID' 
          AND (PRN = :prn OR PRN IN (SELECT OLD_PRN FROM NOVA_PATIENT_MERGE_DETAIL WHERE NEW_PRN = :prn))) 
	    ) ORDER BY 8 DESC FETCH FIRST 5 ROWS ONLY) ORDER BY TO_DATE(DATE_TIME, 'DD Mon YYYY') ASC
    `
    list := make([]model.NovaPatientVitalSignsDetail, 0)
    err := db.SelectContext(s.ctx, &list, query, sql.Named("prn", prn), sql.Named("vitalSignsCode", vitalSignsCode))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVitalSignsDetailService) GetPatientVitalSignsHistoryBP(prn string, vitalSignsCode string, conn *sqlx.DB) ([]model.NovaPatientVitalSignsDetail, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, to_char(date_time, 'DD Mon YYYY') as DATE_TIME FROM ( 
          SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
          SELECT dbms_random.string('U', 30) as ID, d.REF_NO, d.CODE, d.DESCRIPTION, d.VALUE1, d.VALUE2, d.UNIT, vi.date_time as DATE_TIME 
          FROM NOVA_PATIENT_VITALS_DET d, NOVA_PATIENT_VITALS vi, NOVA_VISIT v 
          WHERE d.CODE = :vitalSignsCode 
          AND d.REF_NO = vi.REF_NO 
          AND vi.ACCOUNT_NO = v.ACCOUNT_NO 
          AND d.REF_NO IN (SELECT REF_NO FROM NOVA_PATIENT_VITALS WHERE STATUS <> 'VOID' 
          AND (PRN = :prn OR PRN IN (SELECT OLD_PRN FROM NOVA_PATIENT_MERGE_DETAIL WHERE NEW_PRN = :prn))) 
		) ORDER BY 8 DESC FETCH FIRST 5 ROWS ONLY) ORDER BY TO_DATE(DATE_TIME, 'DD Mon YYYY') ASC
    `
    list := make([]model.NovaPatientVitalSignsDetail, 0)
    err := db.SelectContext(s.ctx, &list, query, sql.Named("prn", prn), sql.Named("vitalSignsCode", vitalSignsCode))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaVitalSignsDetailService) GetPatientVitalSignsHistoryForDashboard(
    prn string, vitalCodeHeight string, vitalCodeWeight string, vitalCodeBP string, vitalCodeBMI string, 
    vitalCodePulse string, conn *sqlx.DB,
) ([]model.NovaPatientVitalSignsDetail, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT * FROM ( 
        WITH PATIENT_VISITS AS ( 
        SELECT PRN 
        FROM NOVA_VISIT 
        WHERE PRN = :prn 
        UNION 
        SELECT OLD_PRN 
        FROM NOVA_PATIENT_MERGE_DETAIL 
        WHERE NEW_PRN = :prn 
        ), 
        VALID_REF_NOS AS ( 
        SELECT REF_NO 
        FROM NOVA_PATIENT_VITALS 
        WHERE STATUS <> 'VOID' 
        AND PRN IN (SELECT PRN FROM PATIENT_VISITS) 
        ), 
        VITALS AS ( 
        SELECT d.ref_no, d.code, d.description, d.value1, d.value2, d.unit, vi.date_time 
        FROM NOVA_PATIENT_VITALS_DET d 
        JOIN NOVA_PATIENT_VITALS vi ON d.REF_NO = vi.ref_no 
        JOIN VALID_REF_NOS vr ON d.REF_NO = vr.REF_NO 
        WHERE LOWER(d.DESCRIPTION) IN (:vitalCodeHeight, :vitalCodeWeight) 
        OR d.CODE IN (:vitalCodeBP, :vitalCodeBMI, :vitalCodePulse) 
        ) 
        SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
        SELECT 
        dbms_random.string('U', 30) as ID, 
        REF_NO, 
        CODE, 
        DESCRIPTION, 
        CAST(ROUND(value1) AS NUMBER) as VALUE1, 
        CAST(value2 AS VARCHAR2(100)) as VALUE2, 
        CAST(unit AS VARCHAR2(50)) as UNIT, 
        TO_CHAR(date_time, 'DD Mon YYYY') as DATE_TIME 
        FROM VITALS 
        WHERE LOWER(DESCRIPTION) = :vitalCodeHeight 
        ORDER BY date_time DESC FETCH FIRST 5 ROWS ONLY 
        ) 
        UNION ALL 
        SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
        SELECT 
        dbms_random.string('U', 30) as ID, 
        REF_NO, 
        CODE, 
        DESCRIPTION, 
        CAST(value1 AS NUMBER) as VALUE1, 
        CAST(value2 AS VARCHAR2(100)) as VALUE2, 
        CAST(unit AS VARCHAR2(50)) as UNIT, 
        TO_CHAR(date_time, 'DD Mon YYYY') as DATE_TIME 
        FROM VITALS 
        WHERE LOWER(description) = :vitalCodeWeight 
        ORDER BY date_time DESC FETCH FIRST 5 ROWS ONLY 
        ) 
        UNION ALL 
        SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
        SELECT 
        dbms_random.string('U', 30) as ID, 
        REF_NO, 
        CODE, 
        DESCRIPTION, 
        CAST(value1 AS NUMBER) as VALUE1, 
        CAST(value2 AS VARCHAR2(100)) as VALUE2, 
        CAST(unit AS VARCHAR2(50)) as UNIT, 
        TO_CHAR(date_time, 'DD Mon YYYY') as DATE_TIME 
        FROM VITALS 
        WHERE CODE = :vitalCodeBP 
        ORDER BY date_time DESC FETCH FIRST 5 ROWS ONLY 
        ) 
        UNION ALL 
        SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
        SELECT 
        dbms_random.string('U', 30) as ID, 
        REF_NO, 
        CODE, 
        DESCRIPTION, 
        CAST(value1 AS NUMBER) as VALUE1, 
        CAST(value2 AS VARCHAR2(100)) as VALUE2, 
        CAST(unit AS VARCHAR2(50)) as UNIT, 
        TO_CHAR(date_time, 'DD Mon YYYY') as DATE_TIME 
        FROM VITALS 
        WHERE CODE = :vitalCodeBMI 
        ORDER BY date_time DESC FETCH FIRST 5 ROWS ONLY 
        ) 
        UNION ALL 
        SELECT ID, REF_NO, CODE, DESCRIPTION, VALUE1, VALUE2, UNIT, DATE_TIME FROM ( 
        SELECT 
        dbms_random.string('U', 30) as ID, 
        REF_NO, 
        CODE, 
        DESCRIPTION, 
        CAST(value1 AS NUMBER) as VALUE1, 
        CAST(value2 AS VARCHAR2(100)) as VALUE2, 
        CAST(unit AS VARCHAR2(50)) as UNIT, 
        TO_CHAR(date_time, 'DD Mon YYYY') as DATE_TIME 
        FROM VITALS 
        WHERE CODE = :vitalCodePulse 
        ORDER BY date_time DESC FETCH FIRST 5 ROWS ONLY 
        )) ORDER BY CODE, TO_DATE(date_time, 'DD Mon YYYY') ASC
    `
    list := make([]model.NovaPatientVitalSignsDetail, 0)
    err := db.SelectContext(s.ctx, &list, query, 
        sql.Named("prn", prn), 
        sql.Named("vitalCodeHeight", vitalCodeHeight), 
        sql.Named("vitalCodeWeight", vitalCodeWeight), 
        sql.Named("vitalCodeBP", vitalCodeBP), 
        sql.Named("vitalCodeBMI", vitalCodeBMI), 
        sql.Named("vitalCodePulse", vitalCodePulse),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
