package hpackage

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/model"
    "vesaliusm/model/hpackage"
    "vesaliusm/service/novaDoctor"
    "vesaliusm/service/patientPurchaseDetails"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
)

type PackageService struct {
    db                            *sqlx.DB
    ctx                           context.Context
    novaDoctorService             *novaDoctor.NovaDoctorService
    patientPurchaseDetailsService *patientPurchaseDetails.PatientPurchaseDetailsService
}

func NewPackageService(db *sqlx.DB, ctx context.Context) *PackageService {
    return &PackageService{
        db:                            db,
        ctx:                           ctx,
        novaDoctorService:             novaDoctor.NovaDoctorSvc,
        patientPurchaseDetailsService: patientPurchaseDetails.PatientPurchaseDetailsSvc,
    }
}

func (s *PackageService) ResizeAllPackageImage(image string, packageId int64) error {
    query := `UPDATE HOSPITAL_PACKAGE SET RESIZE_PACKAGE_IMG = :img WHERE PACKAGE_ID = :id`
    _, err := s.db.ExecContext(s.ctx, query, image, packageId)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return err
}

func (s *PackageService) ListByKeyword(keyword string, keyword2 string, keyword3 string, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountByKeyword(keyword, keyword2, keyword3, s.db)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindByKeyword(keyword, keyword2, keyword3, pager.GetLowerBound(), pager.PageSize, s.db)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *PackageService) CountByKeyword(keyword string, keyword2 string, keyword3 string, conn *sqlx.DB) (int, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    conditions, args := buildKeywordConditions(keyword, keyword2, keyword3)
    base := `SELECT COUNT(hp.PACKAGE_ID) AS COUNT FROM HOSPITAL_PACKAGE hp`
    query := base + whereClause(conditions)

    var count int
    err := db.GetContext(s.ctx, &count, query, args...)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *PackageService) FindByKeyword(keyword string, keyword2 string, keyword3 string, offset int, limit int, conn *sqlx.DB) ([]hpackage.Package, error) {
    db := conn
    if db == nil {
        db = s.db
    }
    conditions, args := buildKeywordConditions(keyword, keyword2, keyword3)
    args = append(args, sql.Named("offset", offset))
    args = append(args, sql.Named("limit", limit))

    base := `
        SELECT hp.*, (SELECT NVL(COUNT(ppd.PACKAGE_ID), 0) 
        FROM PATIENT_PURCHASE_DETAILS ppd 
        WHERE ppd.PACKAGE_ID = hp.PACKAGE_ID) AS PACKAGE_TOTAL_SOLD
        FROM HOSPITAL_PACKAGE hp
    `
    base = strings.Replace(base, "hp.*", getHospitalPackageCols("hp."), 1)
    query := base + whereClause(conditions) +
        ` ORDER BY hp.PACKAGE_NAME OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

    rows, err := db.QueryxContext(s.ctx, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    defer rows.Close()

    list := make([]hpackage.Package, 0)
    for rows.Next() {
        var o hpackage.Package
        if err = rows.StructScan(&o); err != nil {
            utils.LogError(err)
            return nil, err
        }
        if o.PackageAssignedDoctor.Valid {
            doctorName, err := s.novaDoctorService.FindDoctorNameByDoctorId(o.PackageAssignedDoctor.Int64, s.db)
            if err != nil {
                utils.LogError(err)
                return nil, err
            }
            o.PackageAssignedDoctorName = doctorName
        }
        list = append(list, o)
    }
    return list, nil
}

func buildKeywordConditions(keyword string, keyword2 string, keyword3 string) ([]string, []interface{}) {
    var (
        conds []string
        args  []interface{}
    )

    if keyword != "" {
        conds = append(conds, `(LOWER(hp.PACKAGE_CODE) LIKE :keyword OR LOWER(hp.PACKAGE_NAME) LIKE :keyword)`)
        args = append(args, sql.Named("keyword", strings.ToLower(keyword)))
    }
    if keyword2 != "" {
        conds = append(conds, `TRUNC(hp.PACKAGE_START_DATETIME) = TO_DATE(:keyword2, 'dd/mm/yyyy')`)
        args = append(args, sql.Named("keyword2", strings.ToLower(keyword2)))
    }
    if keyword3 != "" {
        conds = append(conds, `TRUNC(hp.PACKAGE_END_DATETIME) = TO_DATE(:keyword3, 'dd/mm/yyyy')`)
        args = append(args, sql.Named("keyword3", strings.ToLower(keyword3)))
    }
    return conds, args
}

func whereClause(conds []string) string {
    if len(conds) == 0 {
        return ""
    }
    return " WHERE " + strings.Join(conds, " AND ")
}

func getHospitalPackageCols(prefix string) string {
    lx := []string{
        "PACKAGE_NAME",
        "PACKAGE_IMG",
        "PACKAGE_VALIDITY",
        "PACKAGE_ALLOW_APPT",
    }
    ls := make([]string, 0)
    for _, v := range lx {
        ls = append(ls, prefix+v)
    }
    return strings.Join(ls, ", ")
}
