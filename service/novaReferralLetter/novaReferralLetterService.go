package novaReferralLetter

import (
    "context"
    "database/sql"
    "vesaliusm/database"
    model "vesaliusm/model/healthCare"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

var NovaReferralLetterSvc *NovaReferralLetterService = NewNovaReferralLetterService(database.GetDbrs(), database.GetCtx())

type NovaReferralLetterService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewNovaReferralLetterService(db *sqlx.DB, ctx context.Context) *NovaReferralLetterService {
    return &NovaReferralLetterService{db: db, ctx: ctx}
}

func (s *NovaReferralLetterService) FindExtReferralLetterByPrnAndAccountNo(prn string, accountNo string, conn *sqlx.DB) ([]model.NovaReferralLetter, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `
        SELECT REFERRAL_REF_NO, PRN, ACCOUNT_NO, REFERRAL_DATE_TIME, REFERRER_DOCTOR_MCR, 
        REFERRAL_TO, TITLE_DEPARTMENT, ADDRESS, REFERRAL_LETTER, 'EXTERNAL' AS REFFERAL_TYPE 
        FROM NOVA_EXT_REFERRAL_LETTER 
        WHERE PRN = :prn 
        AND ACCOUNT_NO = :accountNo 
        ORDER BY TO_DATE(REFERRAL_DATE_TIME, 'DD/MM/YYYY HH24:MI:SS') DESC
    `
    var list []model.NovaReferralLetter
    err := db.SelectContext(s.ctx, &list, query, prn, accountNo)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *NovaReferralLetterService) FindAllReferralLetterByPrnAndAccountNo(prn string, accountNo string, conn *sqlx.DB) ([]model.NovaReferralLetter, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    query := `
        SELECT REFERRAL_REF_NO, PRN, ACCOUNT_NO, REFERRAL_DATE_TIME, 
        (SELECT DOCTOR_NAME FROM NOVA_DOCTOR WHERE ROWNUM = 1 AND MCR_NO = REFERRER_DOCTOR_MCR) AS REFERRER_DOCTOR_MCR, 
        REFERRAL_TO, TITLE_DEPARTMENT, ADDRESS, REFERRAL_LETTER, 'EXTERNAL' AS REFFERAL_TYPE 
        FROM NOVA_EXT_REFERRAL_LETTER 
        WHERE PRN = :prn 
        AND ACCOUNT_NO = :accountNo 
        UNION ALL 
        SELECT REFERRAL_REF_NO, PRN, ACCOUNT_NO, REFERRAL_DATE_TIME, 
        (SELECT DOCTOR_NAME FROM NOVA_DOCTOR WHERE ROWNUM = 1 AND MCR_NO = REFERRER_DOCTOR_MCR) AS REFERRER_DOCTOR_MCR, 
        (SELECT DOCTOR_NAME FROM NOVA_DOCTOR WHERE ROWNUM = 1 AND MCR_NO = REFERRAL_DOCTOR_MCR) AS REFERRAL_DOCTOR_MCR, 
        SPECIALTY, SUBJECT, REFERRAL_LETTER, 'INTERNAL' AS REFFERAL_TYPE 
        FROM NOVA_INT_REFERRAL_LETTER 
        WHERE PRN = :prn 
        AND ACCOUNT_NO = :accountNo
    `
    var list []model.NovaReferralLetter
    err := db.SelectContext(s.ctx, &list, query, sql.Named("prn", prn), sql.Named("accountNo", accountNo))
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}
