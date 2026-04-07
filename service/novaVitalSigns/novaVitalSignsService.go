package novaVitalSigns

import (
    "context"
    "strings"
    "vesaliusm/database"
    model "vesaliusm/model/healthCare"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var NovaVitalSignsSvc *NovaVitalSignsService = NewNovaVitalSignsService(database.GetDbrs(), database.GetCtx())

type NovaVitalSignsService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaVitalSignsService(db *sqlx.DB, ctx context.Context) *NovaVitalSignsService {
    return &NovaVitalSignsService{db: db, ctx: ctx}
}

func (s *NovaVitalSignsService) FindPatientVitalSignsByAccountNo(accountNo string, conn *sqlx.DB) ([]model.NovaPatientVitalSigns, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `
        SELECT * FROM (  
          SELECT t.TEMPLATE_REF_NO AS REF_NO, t.PRN, t.ACCOUNT_NO, t.DATE_TIME   
          FROM NOVA_PATIENT_TEMPLATE t, NOVA_VISIT v  
          WHERE t.ACCOUNT_NO = v.ACCOUNT_NO  
          AND v.ACCOUNT_NO = :accountNo    
          AND ((v.VISIT_TYPE IN ('DAY-SURGERY', 'INPATIENT') AND v.VISIT_STATUS <> 'CANCEL' AND v.ADMISSION_STATUS = 'DISCHARGED')    
          OR (v.VISIT_TYPE IN ('OUTPATIENT', 'EXTERNAL', 'ACCIDENT-EMERGENCY') AND v.VISIT_STATUS <> 'CANCEL'))   
          AND t.TEMPLATE_REF_NO IN (SELECT TEMPLATE_REF_NO FROM NOVA_PATIENT_TEMPLATE_DETAIL  
          WHERE TEMPLATE_REF_NO IN (SELECT TEMPLATE_REF_NO FROM NOVA_PATIENT_TEMPLATE WHERE TEMPLATE_NAME = 'PHYSICAL EXAMINATION')  
          AND SECTION_NAME = 'GENERAL CONDITION AND VITAL SIGNS')  
          UNION ALL  
          SELECT p.REF_NO, p.PRN, p.ACCOUNT_NO, p.DATE_TIME    
          FROM NOVA_PATIENT_VITALS p, NOVA_VISIT v    
          WHERE p.ACCOUNT_NO = v.ACCOUNT_NO   
          AND v.ACCOUNT_NO = :accountNo  
          AND ((v.VISIT_TYPE IN ('DAY-SURGERY', 'INPATIENT') AND v.VISIT_STATUS <> 'CANCEL' AND v.ADMISSION_STATUS = 'DISCHARGED')    
          OR (v.VISIT_TYPE IN ('OUTPATIENT', 'EXTERNAL', 'ACCIDENT-EMERGENCY') AND v.VISIT_STATUS <> 'CANCEL'))   
          AND p.STATUS <> 'VOID'  
            )  
            ORDER BY DATE_TIME DESC  
            FETCH FIRST 1 ROW ONLY
    `
    query = strings.Replace(query, "*", utils.GetDbCols(model.NovaPatientVitalSigns{}, ""), 1)
    list := make([]model.NovaPatientVitalSigns, 0)
    err := db.SelectContext(s.ctx, &list, query, accountNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
