package hpackage

import (
    "context"
    "database/sql"
    "strings"
    "vesaliusm/database"
    "vesaliusm/model"
    "vesaliusm/model/hpackage"
    "vesaliusm/service/novaDoctor"
    "vesaliusm/service/patientPurchaseDetails"
    "vesaliusm/utils"

    "github.com/jmoiron/sqlx"
    "github.com/gofiber/fiber/v3"
)

var PackageSvc *PackageService = NewPackageService(database.GetDb(), database.GetCtx())

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
    db := database.GetFromCon(conn, s.db)
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
    db := database.GetFromCon(conn, s.db)
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

    list := make([]hpackage.Package, 0)
    err := db.SelectContext(s.ctx, &list, query, args...)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    for i := range list {
        if list[i].PackageAssignedDoctor.Valid {
            doctorName, err := s.novaDoctorService.FindDoctorNameByDoctorId(list[i].PackageAssignedDoctor.Int64, s.db)
            if err != nil {
                utils.LogError(err)
                return nil, err
            }
            list[i].PackageAssignedDoctorName = doctorName
        }
    }
    return list, nil
}

func (s *PackageService) Save(o hpackage.Package, adminId int64) error {
    query := `
        INSERT INTO HOSPITAL_PACKAGE
        (PACKAGE_CODE, PACKAGE_NAME, PACKAGE_DESC, PACKAGE_IMG, RESIZE_PACKAGE_IMG,
        PACKAGE_START_DATETIME, PACKAGE_END_DATETIME, PACKAGE_VALIDITY_TYPE, PACKAGE_VALIDITY, 
        PACKAGE_VALIDITY_DATETIME, PACKAGE_TNC, PACKAGE_PRICE, PACKAGE_MAX_PURCHASE, 
        PACKAGE_ASSIGNED_DOCTOR, PACKAGE_ALLOW_APPT, PACKAGE_DISPLAY_ORDER, PACKAGE_EXT_LINK, USER_CREATE)
        VALUES (:packageCode, :packageName, :packageDesc, :packageImage, :resizePackageImage,
        TO_DATE(:packageStartDateTime, 'DD/MM/YYYY hh24:mi'), TO_DATE(:packageEndDateTime, 'DD/MM/YYYY hh24:mi'), :packageValidityType, :packageValidity, 
        TO_DATE(:packageValidityDate, 'DD/MM/YYYY'), :packageTnc, :packagePrice, :packageMaxPurchase,
        :packageAssignedDoctor, :packageAllowAppt, :packageDisplayOrder, :packageExtLink, :adminId)
    `
    
    _, err := s.db.ExecContext(s.ctx, query,
        sql.Named("packageCode", o.PackageCode.String),
        sql.Named("packageName", o.PackageName.String),
        sql.Named("packageDesc", o.PackageDesc.String),
        sql.Named("packageImage", o.PackageImage.String),
        sql.Named("resizePackageImage", o.ResizePackageImage),
        sql.Named("packageStartDateTime", o.PackageStartDateTime.String),
        sql.Named("packageEndDateTime", o.PackageEndDateTime.String),
        sql.Named("packageValidityType", o.PackageValidityType.String),
        sql.Named("packageValidity", o.PackageValidity.Int64),
        sql.Named("packageValidityDate", o.PackageValidityDate.String),
        sql.Named("packageTnc", o.PackageTnc.String),
        sql.Named("packagePrice", o.PackagePrice.Float64),
        sql.Named("packageMaxPurchase", o.PackageMaxPurchase.Int32),
        sql.Named("packageAssignedDoctor", o.PackageAssignedDoctor.Int64),
        sql.Named("packageAllowAppt", o.PackageAllowAppt.String),
        sql.Named("packageDisplayOrder", o.PackageDisplayOrder.Int32),
        sql.Named("packageExtLink", o.PackageExtLink.String),
        sql.Named("adminId", adminId),
    )
    
    if err != nil {
        utils.LogError(err)
        return err
    }
    
    return nil
}

func (s *PackageService) Update(o hpackage.Package, adminId int64) error {
    query := `
        UPDATE HOSPITAL_PACKAGE SET
          PACKAGE_CODE = :packageCode,
          PACKAGE_NAME = :packageName,
          PACKAGE_DESC = :packageDesc,
          PACKAGE_IMG = :packageImage,
          RESIZE_PACKAGE_IMG = :resizePackageImage,
          PACKAGE_START_DATETIME = TO_DATE(:packageStartDateTime, 'DD/MM/YYYY hh24:mi'),
          PACKAGE_END_DATETIME = TO_DATE(:packageEndDateTime, 'DD/MM/YYYY hh24:mi'),
          PACKAGE_VALIDITY_TYPE = :packageValidityType,
          PACKAGE_VALIDITY = :packageValidity,
          PACKAGE_VALIDITY_DATETIME = TO_DATE(:packageValidityDate, 'DD/MM/YYYY'),
          PACKAGE_TNC = :packageTnc,
          PACKAGE_PRICE = :packagePrice,
          PACKAGE_MAX_PURCHASE = :packageMaxPurchase,
          PACKAGE_ASSIGNED_DOCTOR = :packageAssignedDoctor,
          PACKAGE_ALLOW_APPT = :packageAllowAppt,
          PACKAGE_DISPLAY_ORDER = :packageDisplayOrder,
          PACKAGE_EXT_LINK = :packageExtLink,
          USER_UPDATE = :adminId,
          DATE_UPDATE = CURRENT_TIMESTAMP
        WHERE PACKAGE_ID = :packageId`
    
    _, err := s.db.Exec(query,
        sql.Named("packageCode", o.PackageCode.String),
        sql.Named("packageName", o.PackageName.String),
        sql.Named("packageDesc", o.PackageDesc.String),
        sql.Named("packageImage", o.PackageImage.String),
        sql.Named("resizePackageImage", o.ResizePackageImage),
        sql.Named("packageStartDateTime", o.PackageStartDateTime.String),
        sql.Named("packageEndDateTime", o.PackageEndDateTime.String),
        sql.Named("packageValidityType", o.PackageValidityType.String),
        sql.Named("packageValidity", o.PackageValidity.Int64),
        sql.Named("packageValidityDate", o.PackageValidityDate.String),
        sql.Named("packageTnc", o.PackageTnc.String),
        sql.Named("packagePrice", o.PackagePrice.Float64),
        sql.Named("packageMaxPurchase", o.PackageMaxPurchase.Int32),
        sql.Named("packageAssignedDoctor", o.PackageAssignedDoctor.Int64),
        sql.Named("packageAllowAppt", o.PackageAllowAppt.String),
        sql.Named("packageDisplayOrder", o.PackageDisplayOrder.Int32),
        sql.Named("packageExtLink", o.PackageExtLink.String),
        sql.Named("adminId", adminId),
        sql.Named("packageId", o.PackageId.Int64),
    )
    
    if err != nil {
        utils.LogError(err)
        return err
    }
    
    return nil
}

func (s *PackageService) FindAllApp(offset int, limit int, isHome bool, conn *sqlx.DB) ([]hpackage.Package, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT PACKAGE_ID, PACKAGE_CODE, PACKAGE_NAME, PACKAGE_DESC, RESIZE_PACKAGE_IMG AS PACKAGE_IMG, 
        PACKAGE_START_DATETIME, PACKAGE_END_DATETIME, PACKAGE_VALIDITY, PACKAGE_TNC, PACKAGE_PRICE, 
        PACKAGE_MAX_PURCHASE, PACKAGE_ASSIGNED_DOCTOR, PACKAGE_ALLOW_APPT, PACKAGE_DISPLAY_ORDER, PACKAGE_EXT_LINK, 
        USER_CREATE, DATE_CREATE, USER_UPDATE, DATE_UPDATE, PACKAGE_VALIDITY_TYPE, 
        PACKAGE_VALIDITY_DATETIME
        FROM HOSPITAL_PACKAGE
        WHERE (
        PACKAGE_START_DATETIME <= CURRENT_TIMESTAMP AND (
            PACKAGE_END_DATETIME >= CURRENT_TIMESTAMP OR
            PACKAGE_END_DATETIME IS NULL
        )
        )
        ORDER BY PACKAGE_START_DATETIME
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    if isHome {
        query = `
            SELECT PACKAGE_ID, PACKAGE_CODE, PACKAGE_NAME, PACKAGE_DESC, RESIZE_PACKAGE_IMG AS PACKAGE_IMG, 
            PACKAGE_START_DATETIME, PACKAGE_END_DATETIME, PACKAGE_VALIDITY, PACKAGE_TNC, PACKAGE_PRICE, 
            PACKAGE_MAX_PURCHASE, PACKAGE_ASSIGNED_DOCTOR, PACKAGE_ALLOW_APPT, PACKAGE_DISPLAY_ORDER, PACKAGE_EXT_LINK, 
            USER_CREATE, DATE_CREATE, USER_UPDATE, DATE_UPDATE, PACKAGE_VALIDITY_TYPE, 
            PACKAGE_VALIDITY_DATETIME
            FROM HOSPITAL_PACKAGE
            WHERE (
            PACKAGE_START_DATETIME <= CURRENT_TIMESTAMP AND (
            PACKAGE_END_DATETIME >= CURRENT_TIMESTAMP OR
            PACKAGE_END_DATETIME IS NULL
            )
        ) AND PACKAGE_DISPLAY_ORDER = 1
        ORDER BY PACKAGE_START_DATETIME
        FETCH FIRST 5 ROWS ONLY
        `
    }
    list := make([]hpackage.Package, 0)
    var err error
    if isHome {
        err = db.SelectContext(s.ctx, &list, query)
    } else {
        err = db.SelectContext(s.ctx, &list, query, offset, limit)
    }

    if err != nil {
        utils.LogError(err)
        return nil, err
    }

    return list, nil
}

func (s *PackageService) FindAll(offset int, limit int, conn *sqlx.DB) ([]hpackage.Package, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT hp.*, (SELECT NVL(COUNT(ppd.PACKAGE_ID), 0) 
        FROM PATIENT_PURCHASE_DETAILS ppd 
        WHERE ppd.PACKAGE_ID = hp.PACKAGE_ID) AS PACKAGE_TOTAL_SOLD
        FROM HOSPITAL_PACKAGE hp
        ORDER BY hp.PACKAGE_NAME 
        OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
    `
    query = strings.Replace(query, "hp.*", getHospitalPackageCols("hp."), 1)
    list := make([]hpackage.Package, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        if list[i].PackageAssignedDoctor.Valid {
            doctorName, err := s.novaDoctorService.FindDoctorNameByDoctorId(list[i].PackageAssignedDoctor.Int64, s.db)
            if err != nil {
                utils.LogError(err)
                return nil, err
            }
            list[i].PackageAssignedDoctorName = doctorName
        }
    }
    return list, nil
}

func (s *PackageService) ListApp(isHome bool, page string, limit string) (*model.PagedList, error) {
    total, err := s.CountApp(isHome, s.db)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAllApp(pager.GetLowerBound(), pager.PageSize, isHome, s.db)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *PackageService) CountApp(isHome bool, conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `
        SELECT COUNT(PACKAGE_ID) AS COUNT FROM HOSPITAL_PACKAGE
         WHERE (
          PACKAGE_START_DATETIME <= CURRENT_TIMESTAMP AND (
            PACKAGE_END_DATETIME >= CURRENT_TIMESTAMP OR
            PACKAGE_END_DATETIME IS NULL
          )
        )
    `
    if isHome {
        query = `
            SELECT COUNT(PACKAGE_ID) AS COUNT FROM HOSPITAL_PACKAGE
             WHERE (
              PACKAGE_START_DATETIME <= CURRENT_TIMESTAMP AND (
              PACKAGE_END_DATETIME >= CURRENT_TIMESTAMP OR
              PACKAGE_END_DATETIME IS NULL
             )
            ) AND PACKAGE_DISPLAY_ORDER = 1
        `
    }
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *PackageService) List(page string, limit string) (*model.PagedList, error) {
    total, err := s.Count(s.db)
    if err != nil {
        return nil, err
    }

    pager := model.GetPager(total, page, limit)
    list, err := s.FindAll(pager.GetLowerBound(), pager.PageSize, s.db)
    if err != nil {
        return nil, err
    }

    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *PackageService) Count(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT COUNT(PACKAGE_ID) AS COUNT FROM HOSPITAL_PACKAGE`
    var count int
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *PackageService) FindNameByPackageId(packageId int64) (*hpackage.Package, error) {
    query := `SELECT PACKAGE_NAME FROM HOSPITAL_PACKAGE WHERE PACKAGE_ID = :id`
    var o hpackage.Package
    err := s.db.GetContext(s.ctx, &o, query, packageId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
}

func (s *PackageService) FindPackageStatusByPackageId(packageId int64) (fiber.Map, error) {
    packageStatus := fiber.Map{}
    expiry, err := s.patientPurchaseDetailsService.GetPackageExpiryStatus(s.db, packageId)
    if err != nil {
        return packageStatus, err
    }

    if expiry == "Expired" {
        packageStatus["expired"] = 1
    } else {
        packageStatus["expired"] = 0
    }
    
    soldout, err := s.patientPurchaseDetailsService.GetPackageSoldoutStatus(packageId)
    if err != nil {
        return packageStatus, err
    }

    if soldout == "Sold Out" {
        packageStatus["soldout"] = 1
    } else {
        packageStatus["soldout"] = 0
    }

    exceedPurchase, err := s.patientPurchaseDetailsService.GetPackageExceedPurchaseStatus(s.db, packageId, 999)
    if err != nil {
        return packageStatus, err
    }

    if exceedPurchase != nil {
        availableToPurchase := exceedPurchase.RecommendedQuantity.Int32
        packageStatus["availableToPurchase"] = availableToPurchase
    }
    
    return packageStatus, nil
}

func (s *PackageService) FindByPackageId(packageId int64) (*hpackage.Package, error) {
    query := `
        SELECT hp.*, (SELECT NVL(COUNT(ppd.PACKAGE_ID), 0) 
        FROM PATIENT_PURCHASE_DETAILS ppd 
        WHERE ppd.PACKAGE_ID = hp.PACKAGE_ID) AS PACKAGE_TOTAL_SOLD
        FROM HOSPITAL_PACKAGE hp
        WHERE hp.PACKAGE_ID = :id
    `
    query = strings.Replace(query, "hp.*", getHospitalPackageCols("hp."), 1)
    var o hpackage.Package
    err := s.db.GetContext(s.ctx, &o, query, packageId)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &o, nil
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
