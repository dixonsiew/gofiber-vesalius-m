package maintenance

import (
    "vesaliusm/model"
    "vesaliusm/utils"
)

func (s *MaintenanceService) GetAllStatisticAppointments() ([]model.StatisticAppointment, error) {
    query := `
        WITH past_12_months AS (
          SELECT
            ADD_MONTHS(TRUNC(SYSDATE, 'MM'), -LEVEL + 1) AS month_start
          FROM DUAL
          CONNECT BY LEVEL <= 12
        ),
        appointments_grouped AS (
          SELECT
            TRUNC(DATE_APPT, 'MM') AS month_start,
            APPT_STATUS,
            COUNT(*) AS total_appointments
          FROM NOVA_DOCTOR_PATIENT_APPT
          WHERE DATE_APPT >= ADD_MONTHS(TRUNC(SYSDATE, 'MM'), -11)
          GROUP BY TRUNC(DATE_APPT, 'MM'), APPT_STATUS
        ),
        pivoted AS (
          SELECT
            month_start,
            SUM(CASE WHEN APPT_STATUS = 'CHANGED' THEN total_appointments ELSE 0 END) AS TOTAL_CHANGED,
            SUM(CASE WHEN APPT_STATUS = 'CANCELLED' THEN total_appointments ELSE 0 END) AS TOTAL_CANCELLED,
            SUM(CASE WHEN APPT_STATUS = 'CONFIRMED' THEN total_appointments ELSE 0 END) AS TOTAL_CONFIRMED
          FROM appointments_grouped
          GROUP BY month_start
        )
        SELECT
          TO_CHAR(p.month_start, 'FMMonth') AS MONTH,
          TO_CHAR(p.month_start, 'YYYY') AS YEAR,
          NVL(a.TOTAL_CHANGED, 0) AS TOTAL_CHANGED,
          NVL(a.TOTAL_CANCELLED, 0) AS TOTAL_CANCELLED,
          NVL(a.TOTAL_CONFIRMED, 0) AS TOTAL_CONFIRMED
        FROM past_12_months p
        LEFT JOIN pivoted a ON p.month_start = a.month_start
        ORDER BY p.month_start
    `
    list := make([]model.StatisticAppointment, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) GetAllStatisticMobileRegistrations() ([]model.StatisticMobileRegistration, error) {
    query := `
        WITH past_12_months AS (
          SELECT
            ADD_MONTHS(TRUNC(SYSDATE, 'MM'), -LEVEL + 1) AS month_start
          FROM DUAL
          CONNECT BY LEVEL <= 12
        ),
        registrations_grouped AS (
          SELECT
            TRUNC(REGISTRATION_DATE_TIME, 'MM') AS month_start,
            COUNT(*) AS TOTAL_REGISTRATIONS
          FROM APPLICATION_USER
          WHERE REGISTRATION_DATE_TIME >= ADD_MONTHS(TRUNC(SYSDATE, 'MM'), -11)
          GROUP BY TRUNC(REGISTRATION_DATE_TIME, 'MM')
        )
        SELECT
          TO_CHAR(p.month_start, 'FMMonth') AS MONTH,
          TO_CHAR(p.month_start, 'YYYY') AS YEAR,
          NVL(r.total_registrations, 0) AS TOTAL_REGISTRATIONS
        FROM
          past_12_months p
        LEFT JOIN registrations_grouped r
          ON p.month_start = r.month_start
        ORDER BY
          p.month_start
    `
    list := make([]model.StatisticMobileRegistration, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) GetAllStatisticMobileFeedbacks() ([]model.StatisticMobileFeedback, error) {
    query := `
        WITH past_12_months AS (
          SELECT
            ADD_MONTHS(TRUNC(SYSDATE, 'MM'), -LEVEL + 1) AS month_start
          FROM DUAL
          CONNECT BY LEVEL <= 12
        ),
        feedback_grouped AS (
          SELECT
            TRUNC(DATE_SUBMIT, 'MM') AS month_start,
            COUNT(*) AS TOTAL_FEEDBACKS
          FROM PATIENT_FEEDBACK
          WHERE DATE_SUBMIT >= ADD_MONTHS(TRUNC(SYSDATE, 'MM'), -11)
          GROUP BY TRUNC(DATE_SUBMIT, 'MM')
        )
        SELECT
          TO_CHAR(p.month_start, 'FMMonth') AS MONTH,
          TO_CHAR(p.month_start, 'YYYY') AS YEAR,
          NVL(f.total_feedbacks, 0) AS TOTAL_FEEDBACKS
        FROM
          past_12_months p
        LEFT JOIN feedback_grouped f
          ON p.month_start = f.month_start
        ORDER BY
          p.month_start
    `
    list := make([]model.StatisticMobileFeedback, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) GetAllStatisticMobilePackages() ([]model.StatisticMobilePackage, error) {
    query := `
        WITH past_12_months AS (
          SELECT
            ADD_MONTHS(TRUNC(SYSDATE, 'MM'), -LEVEL + 1) AS month_start
          FROM DUAL
          CONNECT BY LEVEL <= 12
        )
        SELECT
          TO_CHAR(p12.month_start, 'Month') AS MONTH,
          TO_CHAR(p12.month_start, 'YYYY') AS YEAR,

          -- Count each event per month
          (SELECT COUNT(*) FROM PATIENT_PURCHASE_DETAILS p
          WHERE TRUNC(p.PURCHASED_DATETIME, 'MM') = p12.month_start) AS PURCHASED,

          (SELECT COUNT(*) FROM PATIENT_PURCHASE_DETAILS p
          WHERE TRUNC(p.REDEEMED_DATETIME, 'MM') = p12.month_start) AS REDEEMED,

          (SELECT COUNT(*) FROM PATIENT_PURCHASE_DETAILS p
          WHERE TRUNC(p.EXPIRED_DATETIME, 'MM') = p12.month_start) AS EXPIRED

        FROM past_12_months p12
        ORDER BY p12.month_start
    `
    list := make([]model.StatisticMobilePackage, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) GetAllStatisticMobileClubsKids() ([]model.StatisticMobileClubs, error) {
    query := `
        WITH past_12_months AS (
            SELECT
              ADD_MONTHS(TRUNC(SYSDATE, 'MM'), -LEVEL + 1) AS month_start
            FROM DUAL
            CONNECT BY LEVEL <= 12
          ),
          membership_grouped AS (
            SELECT
              TRUNC(KIDS_MEMBERSHIP_JOIN_DATE, 'MM') AS month_start,
              COUNT(*) AS TOTAL_JOINS
            FROM KIDS_CLUB_MEMBERSHIP
            WHERE KIDS_MEMBERSHIP_JOIN_DATE >= ADD_MONTHS(TRUNC(SYSDATE, 'MM'), -11)
            GROUP BY TRUNC(KIDS_MEMBERSHIP_JOIN_DATE, 'MM')
          )
          SELECT
            TO_CHAR(p.month_start, 'FMMonth') AS MONTH,
            TO_CHAR(p.month_start, 'YYYY') AS YEAR,
            NVL(m.total_joins, 0) AS TOTAL_JOINS
          FROM
            past_12_months p
          LEFT JOIN membership_grouped m
            ON p.month_start = m.month_start
          ORDER BY
            p.month_start
    `
    list := make([]model.StatisticMobileClubs, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *MaintenanceService) GetAllStatisticMobileClubsGoldenPearl() ([]model.StatisticMobileClubs, error) {
    query := `
        WITH past_12_months AS (
          SELECT
            ADD_MONTHS(TRUNC(SYSDATE, 'MM'), -LEVEL + 1) AS month_start
          FROM DUAL
          CONNECT BY LEVEL <= 12
        ),
        membership_grouped AS (
          SELECT
            TRUNC(GOLDEN_MEMBERSHIP_JOIN_DATE, 'MM') AS month_start,
            COUNT(*) AS total_joins
          FROM GOLDEN_CLUB_MEMBERSHIP
          WHERE GOLDEN_MEMBERSHIP_JOIN_DATE >= ADD_MONTHS(TRUNC(SYSDATE, 'MM'), -11)
          GROUP BY TRUNC(GOLDEN_MEMBERSHIP_JOIN_DATE, 'MM')
        )
        SELECT
          TO_CHAR(p.month_start, 'FMMonth') AS MONTH,
          TO_CHAR(p.month_start, 'YYYY') AS YEAR,
          NVL(m.total_joins, 0) AS TOTAL_JOINS
        FROM
          past_12_months p
        LEFT JOIN membership_grouped m
          ON p.month_start = m.month_start
        ORDER BY
          p.month_start
    `
    list := make([]model.StatisticMobileClubs, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}