package sql

const GET_ALL_DOCTOR_APPOINTMENTS = `
    SELECT
    A.CALENDAR_DATE,
    TRIM(A.APPT_DAY),
    A.normal_availability_status,
    A.morning_availability_status,
    A.afternoon_availability_status,
    A.night_availability_status,
    CASE
      WHEN A.normal_availability_status = 'AVAILABLE'
      OR A.morning_availability_status = 'AVAILABLE'
      OR A.afternoon_availability_status = 'AVAILABLE'
      OR A.night_availability_status = 'AVAILABLE' THEN 'AVAILABLE'
      WHEN A.normal_availability_status = 'FULL'
      OR A.morning_availability_status = 'FULL'
      OR A.afternoon_availability_status = 'FULL'
      OR A.night_availability_status = 'FULL' THEN 'FULL'
      ELSE 'NOT AVAILABLE'
    END AS daily_status
  FROM
    (
    SELECT
      F.CALENDAR_DATE,
      F.APPT_DAY,
      CASE
        WHEN normal_max_slots = 0
        AND normal_slots_taken = 0 THEN 'NOT AVAILABLE'
        ELSE 
              CASE
          WHEN normal_max_slots > normal_slots_taken THEN 'AVAILABLE'
          ELSE 'FULL'
        END
      END AS normal_availability_status,
      CASE
        WHEN morning_max_slots = 0
        AND morning_slots_taken = 0 THEN 'NOT AVAILABLE'
        ELSE 
              CASE
          WHEN morning_max_slots > morning_slots_taken THEN 'AVAILABLE'
          ELSE 'FULL'
        END
      END AS morning_availability_status,
      CASE
        WHEN afternoon_max_slots = 0
        AND afternoon_slots_taken = 0 THEN 'NOT AVAILABLE'
        ELSE 
              CASE
          WHEN afternoon_max_slots > afternoon_slots_taken THEN 'AVAILABLE'
          ELSE 'FULL'
        END
      END AS afternoon_availability_status,
      CASE
        WHEN night_max_slots = 0
        AND night_slots_taken = 0 THEN 'NOT AVAILABLE'
        ELSE 
              CASE
          WHEN night_max_slots > night_slots_taken THEN 'AVAILABLE'
          ELSE 'FULL'
        END
      END AS night_availability_status
    FROM
      (
      SELECT
        TO_CHAR(trunc(d.dt), 'DD/MM/YYYY') AS CALENDAR_DATE,
        TO_CHAR(d.dt, 'DAY') AS APPT_DAY,
        MAX(CASE WHEN t.SESSION_TYPE = 'NORMAL' THEN t.MAX_SLOTS ELSE 0 END) AS normal_max_slots,
        MAX(CASE WHEN t.SESSION_TYPE = 'NORMAL' THEN t.APPT_CNT ELSE 0 END) AS normal_slots_taken,
        MAX(CASE WHEN t.SESSION_TYPE = 'MORNING' THEN t.MAX_SLOTS ELSE 0 END) AS morning_max_slots,
        MAX(CASE WHEN t.SESSION_TYPE = 'MORNING' THEN t.APPT_CNT ELSE 0 END) AS morning_slots_taken,
        MAX(CASE WHEN t.SESSION_TYPE = 'AFTERNOON' THEN t.MAX_SLOTS ELSE 0 END) AS afternoon_max_slots,
        MAX(CASE WHEN t.SESSION_TYPE = 'AFTERNOON' THEN t.APPT_CNT ELSE 0 END) AS afternoon_slots_taken,
        MAX(CASE WHEN t.SESSION_TYPE = 'NIGHT' THEN t.MAX_SLOTS ELSE 0 END) AS night_max_slots,
        MAX(CASE WHEN t.SESSION_TYPE = 'NIGHT' THEN t.APPT_CNT ELSE 0 END) AS night_slots_taken
      FROM
        (
        SELECT
          TO_DATE(TO_CHAR(CALENDAR_DATE, 'DD-MON-YYYY'), 'DD-MON-YYYY') AS DATE_APPT,
          APPSLOT.SESSION_TYPE,
          APPSLOT.MAX_SLOTS,
          COUNT(PATAPPT.DATE_APPT) AS appt_cnt
        FROM
          (
          SELECT
            TRUNC(TO_DATE(:monthYear, 'MON-YYYY')-1) + LEVEL AS CALENDAR_DATE
          FROM
            DUAL
          CONNECT BY
            LEVEL <= (LAST_DAY(TO_DATE(:monthYear, 'MON-YYYY')) - TO_DATE(:monthYear, 'MON-YYYY')) + 1
      ) CALENDAR
        LEFT JOIN NOVA_DOCTOR_APPT_SLOT APPSLOT 
        ON
          TRIM(TO_CHAR(CALENDAR.CALENDAR_DATE, 'DAY')) = APPSLOT.DAY_OF_WEEK
          AND APPSLOT.DOCTOR_ID = :doctorId
        LEFT JOIN NOVA_DOCTOR_PATIENT_APPT PATAPPT
        ON
          APPSLOT.DOCTOR_ID = PATAPPT.DOCTOR_ID
          AND TRUNC(PATAPPT.DATE_APPT) = CALENDAR.CALENDAR_DATE
          AND TO_DATE(TO_CHAR(PATAPPT.DATE_APPT, 'DD-MM-YYYY HH24:MI:SS'), 'DD-MM-YYYY HH24:MI:SS')
              BETWEEN TO_DATE(TO_CHAR(PATAPPT.DATE_APPT, 'DD/MM/YYYY') || ' ' || APPSLOT.START_TIME, 'DD/MM/YYYY HH24:MI:SS') 
              AND TO_DATE(TO_CHAR(PATAPPT.DATE_APPT, 'DD/MM/YYYY') || ' ' || APPSLOT.END_TIME, 'DD/MM/YYYY HH24:MI:SS')
        WHERE
          CALENDAR.CALENDAR_DATE BETWEEN TO_DATE(:monthYear, 'MON-YYYY') AND LAST_DAY(TO_DATE(:monthYear, 'MON-YYYY')) + 1
        GROUP BY
          TO_DATE(TO_CHAR(CALENDAR.CALENDAR_DATE, 'DD-MON-YYYY'), 'DD-MON-YYYY'),
          APPSLOT.SESSION_TYPE,
          APPSLOT.MAX_SLOTS
        ORDER BY
          TO_DATE(TO_CHAR(CALENDAR.CALENDAR_DATE, 'DD-MON-YYYY'), 'DD-MON-YYYY')
  
      ) t
      RIGHT JOIN (
        SELECT
          TO_DATE(:monthYear, 'MON-YYYY') + LEVEL - 1 AS dt
        FROM
          dual
        CONNECT BY
          LEVEL <= to_number((SELECT LAST_DAY(TO_DATE(:monthYear, 'MON-YYYY')) - TO_DATE(:monthYear, 'MON-YYYY') + 1 AS days_left FROM dual))
      ) d 
      ON
        TRUNC(t.DATE_APPT) = d.dt
      GROUP BY
        d.dt
      ORDER BY
        d.dt
      ) F
  ) A
`

const GET_SINGLEDATE_DOCTOR_APPOINTMENTS = `
    SELECT
    A.CALENDAR_DATE,
    TRIM(A.APPT_DAY),
    A.normal_availability_status,
    A.morning_availability_status,
    A.afternoon_availability_status,
    A.night_availability_status,
    CASE
      WHEN A.normal_availability_status = 'AVAILABLE'
      OR A.morning_availability_status = 'AVAILABLE'
      OR A.afternoon_availability_status = 'AVAILABLE'
      OR A.night_availability_status = 'AVAILABLE' THEN 'AVAILABLE'
      WHEN A.normal_availability_status = 'FULL'
      OR A.morning_availability_status = 'FULL'
      OR A.afternoon_availability_status = 'FULL'
      OR A.night_availability_status = 'FULL' THEN 'FULL'
      ELSE 'NOT AVAILABLE'
    END AS daily_status
  FROM
    (
    SELECT
      F.CALENDAR_DATE,
      F.APPT_DAY,
      CASE
        WHEN normal_max_slots = 0
        AND normal_slots_taken = 0 THEN 'NOT AVAILABLE'
        ELSE 
              CASE
          WHEN normal_max_slots > normal_slots_taken THEN 'AVAILABLE'
          ELSE 'FULL'
        END
      END AS normal_availability_status,
      CASE
        WHEN morning_max_slots = 0
        AND morning_slots_taken = 0 THEN 'NOT AVAILABLE'
        ELSE 
              CASE
          WHEN morning_max_slots > morning_slots_taken THEN 'AVAILABLE'
          ELSE 'FULL'
        END
      END AS morning_availability_status,
      CASE
        WHEN afternoon_max_slots = 0
        AND afternoon_slots_taken = 0 THEN 'NOT AVAILABLE'
        ELSE 
              CASE
          WHEN afternoon_max_slots > afternoon_slots_taken THEN 'AVAILABLE'
          ELSE 'FULL'
        END
      END AS afternoon_availability_status,
      CASE
        WHEN night_max_slots = 0
        AND night_slots_taken = 0 THEN 'NOT AVAILABLE'
        ELSE 
              CASE
          WHEN night_max_slots > night_slots_taken THEN 'AVAILABLE'
          ELSE 'FULL'
        END
      END AS night_availability_status
    FROM
      (
      SELECT
        TO_CHAR(trunc(d.dt), 'DD/MM/YYYY') AS CALENDAR_DATE,
        TO_CHAR(d.dt, 'DAY') AS APPT_DAY,
        MAX(CASE WHEN t.SESSION_TYPE = 'NORMAL' THEN t.MAX_SLOTS ELSE 0 END) AS normal_max_slots,
        MAX(CASE WHEN t.SESSION_TYPE = 'NORMAL' THEN t.APPT_CNT ELSE 0 END) AS normal_slots_taken,
        MAX(CASE WHEN t.SESSION_TYPE = 'MORNING' THEN t.MAX_SLOTS ELSE 0 END) AS morning_max_slots,
        MAX(CASE WHEN t.SESSION_TYPE = 'MORNING' THEN t.APPT_CNT ELSE 0 END) AS morning_slots_taken,
        MAX(CASE WHEN t.SESSION_TYPE = 'AFTERNOON' THEN t.MAX_SLOTS ELSE 0 END) AS afternoon_max_slots,
        MAX(CASE WHEN t.SESSION_TYPE = 'AFTERNOON' THEN t.APPT_CNT ELSE 0 END) AS afternoon_slots_taken,
        MAX(CASE WHEN t.SESSION_TYPE = 'NIGHT' THEN t.MAX_SLOTS ELSE 0 END) AS night_max_slots,
        MAX(CASE WHEN t.SESSION_TYPE = 'NIGHT' THEN t.APPT_CNT ELSE 0 END) AS night_slots_taken
      FROM
        (
        SELECT
          TO_DATE(TO_CHAR(CALENDAR_DATE, 'DD-MON-YYYY'), 'DD-MON-YYYY') AS DATE_APPT,
          APPSLOT.SESSION_TYPE,
          APPSLOT.MAX_SLOTS,
          COUNT(PATAPPT.DATE_APPT) AS appt_cnt
        FROM
          (
          SELECT
            TRUNC(TO_DATE(:dt, 'DD-MON-YYYY')-1) + LEVEL AS CALENDAR_DATE
          FROM
            DUAL
          CONNECT BY
            LEVEL <= (TO_DATE(:dt, 'DD-MON-YYYY') - TO_DATE(:dt, 'DD-MON-YYYY')) + 1
      ) CALENDAR
        LEFT JOIN NOVA_DOCTOR_APPT_SLOT APPSLOT 
        ON
          TRIM(TO_CHAR(CALENDAR.CALENDAR_DATE, 'DAY')) = APPSLOT.DAY_OF_WEEK
          AND APPSLOT.DOCTOR_ID = :doctorId
        LEFT JOIN NOVA_DOCTOR_PATIENT_APPT PATAPPT
        ON
          APPSLOT.DOCTOR_ID = PATAPPT.DOCTOR_ID
          AND TRUNC(PATAPPT.DATE_APPT) = CALENDAR.CALENDAR_DATE
          AND TO_DATE(TO_CHAR(PATAPPT.DATE_APPT, 'DD-MM-YYYY HH24:MI:SS'), 'DD-MM-YYYY HH24:MI:SS')
              BETWEEN TO_DATE(TO_CHAR(PATAPPT.DATE_APPT, 'DD/MM/YYYY') || ' ' || APPSLOT.START_TIME, 'DD/MM/YYYY HH24:MI:SS') 
              AND TO_DATE(TO_CHAR(PATAPPT.DATE_APPT, 'DD/MM/YYYY') || ' ' || APPSLOT.END_TIME, 'DD/MM/YYYY HH24:MI:SS')
          AND PATAPPT.DOCTOR_ID = :doctorId
        WHERE
          CALENDAR.CALENDAR_DATE BETWEEN TO_DATE(:dt, 'DD-MON-YYYY') AND TO_DATE(:dt, 'DD-MON-YYYY') + 1
        GROUP BY
          TO_DATE(TO_CHAR(CALENDAR.CALENDAR_DATE, 'DD-MON-YYYY'), 'DD-MON-YYYY'),
          APPSLOT.SESSION_TYPE,
          APPSLOT.MAX_SLOTS
        ORDER BY
          TO_DATE(TO_CHAR(CALENDAR.CALENDAR_DATE, 'DD-MON-YYYY'), 'DD-MON-YYYY')
  
      ) t
      RIGHT JOIN (
        SELECT
          TO_DATE(:dt, 'DD-MON-YYYY') + LEVEL - 1 AS dt
        FROM
          dual
        CONNECT BY
          LEVEL <= to_number((SELECT TO_DATE(:dt, 'DD-MON-YYYY') - TO_DATE(:dt, 'DD-MON-YYYY') + 1 AS days_left FROM dual))
      ) d 
      ON
        TRUNC(t.DATE_APPT) = d.dt
      GROUP BY
        d.dt
      ORDER BY
        d.dt
      ) F
  ) A
`

const GET_PATIENT_NOK_FAMILY = `
    SELECT 
      N.PATIENT_PRN, 
      N.NOK_ID, 
      NVL2(N.PRN, 'Y', 'N') AS IS_PATIENT,
      INITCAP(N.PATIENT_NOK_NAME) AS NOK_FULLNAME,
      INITCAP(N.RELATION_DESCRIPTION) AS NOK_RELATIONSHIP,
      N.PRN AS NOK_PRN, 
      INITCAP(N.GENDER) AS NOK_GENDER,
      N.NOK_ID AS NOK_DOC_NUMBER,
      N.DOB AS NOK_DOB,
      INITCAP(N.NATIONALITY) AS NOK_NATIONALITY, 
      N.HOME_PHONE AS NOK_CONTACT, 
      N.HOME_ADDRESS AS NOK_ADDRESS,
      REF_NO,
      N.MARITAL_STATUS AS NOK_MARITAL,
      N.EMAIL AS NOK_EMAIL
  FROM NOVA_PATIENT_NOK N,
  (
      SELECT PATIENT_PRN, NOK_ID, MAX(REF_NO) AS LATEST_REF_NO
      FROM NOVA_PATIENT_NOK
      WHERE PATIENT_PRN = :prn
      GROUP BY PATIENT_PRN, NOK_ID
  ) R
  WHERE N.REF_NO = R.LATEST_REF_NO
`

const GET_PATIENT_NOK_NRICPASSPORT = `
    SELECT DOCUMENT_NUMBER FROM NOVA_PATIENT_DOCUMENT
    WHERE PRN = :prn AND DOCUMENT_TYPE = 'NRIC / Passport'
`

const GET_MOBILEAPP_SLOTS = `
    SELECT 
    lars.DAY_OF_WEEK,
    TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY'), 'YYYY-MM-DD') AS REQUESTED_PICKUP_DATE,
    lars.PICKUP_TIME,
    lars.MAX_SLOTS - COALESCE(SUM(
        CASE WHEN lar.VISIT_WITH_COMPANION = 'Y' 
        THEN 2 ELSE 1 END * lar.REQUESTED_SLOTS), 0) AS AVAILABLE_SLOTS
FROM 
    LOGISTIC_ARRANGEMENT_SLOT lars
LEFT JOIN (
    SELECT 
        REQUESTED_PICKUP_DAY AS REQUESTED_PICKUP_DAY,
        REQUESTED_PICKUP_TIME,
        VISIT_WITH_COMPANION,
        COUNT(*)
 AS REQUESTED_SLOTS
    FROM 
        LOGISTIC_ARRANGEMENT_REQUESTER
    WHERE 
        REQUESTED_PICKUP_DATE >= TO_DATE(:flightArrivalDate, 'DD/MM/YYYY')
        AND REQUESTED_PICKUP_DATE <= TO_DATE(:flightArrivalDate, 'DD/MM/YYYY')
        AND LOGISTIC_REQUEST_STATUS IN ('Requested', 'Confirmed')
    GROUP BY 
        REQUESTED_PICKUP_DAY, REQUESTED_PICKUP_TIME, VISIT_WITH_COMPANION
) lar 
ON TRIM(lars.DAY_OF_WEEK) = TRIM(lar.REQUESTED_PICKUP_DAY) 
AND TRIM(lars.PICKUP_TIME) = TRIM(lar.REQUESTED_PICKUP_TIME)
WHERE 
    TRIM(lars.DAY_OF_WEEK) = TRIM(TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY'), 'Day'))
    AND lars.PICKUP_TIME > :flightArrivalTime
GROUP BY 
    lars.DAY_OF_WEEK, 
    TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY'), 'YYYY-MM-DD'), 
    lars.PICKUP_TIME, 
    lars.MAX_SLOTS
HAVING 
    lars.MAX_SLOTS - COALESCE(SUM(
        CASE WHEN lar.VISIT_WITH_COMPANION = 'Y' 
        THEN 2 ELSE 1 END * lar.REQUESTED_SLOTS), 0) > 0

UNION

SELECT 
    lars.DAY_OF_WEEK,
    TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY') + INTERVAL '1' DAY, 'YYYY-MM-DD') AS REQUESTED_PICKUP_DATE,
    lars.PICKUP_TIME,
    lars.MAX_SLOTS - COALESCE(SUM(
        CASE WHEN lar.VISIT_WITH_COMPANION = 'Y' THEN 2 ELSE 1 
        END * lar.REQUESTED_SLOTS), 0) AS AVAILABLE_SLOTS
FROM 
    LOGISTIC_ARRANGEMENT_SLOT lars
LEFT JOIN (
    SELECT 
        REQUESTED_PICKUP_DAY AS REQUESTED_PICKUP_DAY,
        REQUESTED_PICKUP_TIME,
        VISIT_WITH_COMPANION,
        COUNT(*)
 AS REQUESTED_SLOTS
    FROM 
        LOGISTIC_ARRANGEMENT_REQUESTER
    WHERE 
        REQUESTED_PICKUP_DATE >= TO_DATE(:flightArrivalDate, 'DD/MM/YYYY') + INTERVAL '1' DAY
        AND REQUESTED_PICKUP_DATE <= TO_DATE(:flightArrivalDate, 'DD/MM/YYYY') + INTERVAL '1' DAY
        AND LOGISTIC_REQUEST_STATUS IN ('Requested', 'Confirmed')
    GROUP BY 
        REQUESTED_PICKUP_DAY, REQUESTED_PICKUP_TIME, VISIT_WITH_COMPANION
) lar 
ON TRIM(lars.DAY_OF_WEEK) = TRIM(lar.REQUESTED_PICKUP_DAY) 
AND TRIM(lars.PICKUP_TIME) = TRIM(lar.REQUESTED_PICKUP_TIME)
WHERE 
    TRIM(lars.DAY_OF_WEEK) = TRIM(TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY') + INTERVAL '1' DAY, 'Day'))
GROUP BY 
    lars.DAY_OF_WEEK, 
    TO_CHAR(TO_DATE(:flightArrivalDate, 'DD/MM/YYYY') + INTERVAL '1' DAY, 'YYYY-MM-DD'), 
    lars.PICKUP_TIME, 
    lars.MAX_SLOTS
HAVING 
    lars.MAX_SLOTS - COALESCE(SUM(
        CASE WHEN lar.VISIT_WITH_COMPANION = 'Y' 
        THEN 2 ELSE 1 END * lar.REQUESTED_SLOTS), 0) > 0
ORDER BY 
    REQUESTED_PICKUP_DATE, PICKUP_TIME
`
