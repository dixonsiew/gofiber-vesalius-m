package maintenance

import (
    "database/sql"
    "fmt"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    "github.com/lnquy/cron"
    "github.com/nleeper/goment"
)

func (s *MaintenanceService) ListAllCronjobHistories(page string, limit string) (*model.PagedList, error) {
    total, err := s.CountAllCronjobHistories(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllCronjobHistories(pager.GetLowerBound(), pager.PageSize, s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *MaintenanceService) CountAllCronjobHistories(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT COUNT(*) AS COUNT
        FROM (
          SELECT CRONJOB_NAME, ROW_NUMBER() OVER (PARTITION BY CRONJOB_NAME ORDER BY START_DATE DESC) AS ROW_NUM
          FROM CRONJOB_HISTORY
          WHERE TRUNC(START_DATE) = TRUNC(SYSDATE)
        )
        WHERE ROW_NUM = 1
    `
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *MaintenanceService) FindAllCronjobHistories(offset int, limit int, conn *sqlx.DB) ([]model.CronjobHistory, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT 
          ch.CRONJOB_NAME, 
          cm.CRONJOB_DESC AS CRONJOB_NAME_DESC, 
          cm.CRONJOB_EXPRESSION, 
          ch.START_DATE, 
          ch.END_DATE, 
          ch.CRONJOB_STATUS,
          ch.REMARKS
        FROM (
          SELECT CRONJOB_NAME, START_DATE, END_DATE, CRONJOB_STATUS, REMARKS,
          ROW_NUMBER() OVER (PARTITION BY CRONJOB_NAME ORDER BY START_DATE DESC) AS ROW_NUM
          FROM CRONJOB_HISTORY
          WHERE TRUNC(START_DATE) = TRUNC(SYSDATE)
        ) ch
        JOIN CRONJOB_MASTER cm ON ch.CRONJOB_NAME = cm.CRONJOB_NAME
        WHERE ch.ROW_NUM = 1
        ORDER BY ch.START_DATE DESC
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    list := make([]model.CronjobHistory, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    descriptor, err := cron.NewDescriptor(
        cron.Use24HourTimeFormat(true),
        cron.Verbose(false),
    )
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        if list[i].CronjobExpression.Valid {
            desc, err := descriptor.ToDescription(
                list[i].CronjobExpression.String,
                cron.Locale_en,
            )
            if err != nil {
                utils.LogError(err)
                return nil, err
            }
            list[i].CronjobExpressionDesc = utils.NewNullString(desc)
            if strings.Contains(desc, "At") {
                list[i].CronjobExpressionDesc = utils.NewNullString(fmt.Sprintf("Every Day %s", desc))
            } else {
                list[i].CronjobExpressionDesc = utils.NewNullString(desc)
            }

            if list[i].StartDate.Valid {
                g, _ := goment.New(list[i].StartDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
                list[i].StartDate = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
            }
            if list[i].EndDate.Valid {
                g, _ := goment.New(list[i].EndDate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
                list[i].EndDate = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
            }
        }
    }
    return list, nil
}

func (s *MaintenanceService) SearchAllCronjobHistoriesByKeyword(cronjobName string, keyword string, keyword2 string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountCronjobHistoriesByKeyword(cronjobName, keyword, keyword2, s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindCronjobHistoriesByKeyword(cronjobName, keyword, keyword2, pager.GetLowerBound(), pager.PageSize, s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *MaintenanceService) CountCronjobHistoriesByKeyword(cronjobName string, keyword string, keyword2 string, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    conds, args := buildCronJobConditions(cronjobName, keyword, keyword2)
    base := `SELECT COUNT(*) AS COUNT FROM CRONJOB_HISTORY ch`
    query := base + whereClause(conds)

    var count int
    err := db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *MaintenanceService) FindCronjobHistoriesByKeyword(cronjobName string, keyword string, keyword2 string, offset int, limit int, conn *sqlx.DB) ([]model.CronjobHistory, error) {
    db := database.GetFromCon(conn, s.db)
    conditions, args := buildCronJobConditions(cronjobName, keyword, keyword2)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    base := `
        SELECT 
          ch.CRONJOB_NAME, 
          cm.CRONJOB_DESC AS CRONJOB_NAME_DESC, 
          cm.CRONJOB_EXPRESSION, 
          ch.START_DATE, 
          ch.END_DATE, 
          ch.CRONJOB_STATUS,
          ch.REMARKS
        FROM CRONJOB_HISTORY ch
        JOIN CRONJOB_MASTER cm ON ch.CRONJOB_NAME = cm.CRONJOB_NAME
    `
    query := base + whereClause(conditions) + ` ORDER BY ch.START_DATE DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    list := make([]model.CronjobHistory, 0)
    err := db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func buildCronJobConditions(cronjobName string, keyword string, keyword2 string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if cronjobName != "" {
        conds = append(conds, `ch.CRONJOB_NAME = :cronjobName`)
        args = append(args, sql.Named("cronjobName", cronjobName))
    }

    if keyword != "" {
        conds = append(conds, `TO_CHAR(ch.START_DATE, 'dd/mm/yyyy') = :keyword`)
        args = append(args, sql.Named("keyword", keyword))
    }

    if keyword2 != "" {
        conds = append(conds, `TO_CHAR(ch.END_DATE, 'dd/mm/yyyy') = :keyword2`)
        args = append(args, sql.Named("keyword2", keyword2))
    }

    return conds, args
}
