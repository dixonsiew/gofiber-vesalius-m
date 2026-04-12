package hpackage

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "image"
    "image/png"
    "strconv"
    "strings"
    "vesaliusm/dto"
    "vesaliusm/middleware"
    model "vesaliusm/model/hpackage"
    "vesaliusm/service/exportExcel"
    "vesaliusm/service/hpackage"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/nfnt/resize"
    "github.com/nleeper/goment"
)

type PackageController struct {
    packageService     *hpackage.PackageService
    exportExcelService *exportExcel.ExportExcelService
}

func NewPackageController() *PackageController {
    return &PackageController{
        packageService:     hpackage.PackageSvc,
        exportExcelService: exportExcel.ExportExcelSvc,
    }
}

// ProcessResizeImage
//
// @Tags Package
// @Produce text/plain
// @Success 200
// @Router /package/process-resize-image [post]
func (cr *PackageController) ProcessResizeImage(c fiber.Ctx) error {
    lx, err := cr.packageService.FindAll(0, 500, nil)
    if err != nil {
        return err
    }

    for i := range lx {
        if lx[i].PackageImage.Valid {
            resized, err := resizeBase64Image(lx[i].PackageImage.String)
            if err != nil {
                return err
            }

            cr.packageService.ResizeAllPackageImage(resized, lx[i].PackageId.Int64)
        }
    }

    return c.SendString("done")
}

// GetAllPackages
//
// @Tags Package
// @Produce json
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} hpackage.Package
// @Router /package/all [get]
func (cr *PackageController) GetAllPackages(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.packageService.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetAllAppPackages
//
// @Tags Package
// @Produce json
// @Param        isHome            path        string  true   "isHome" default:"1"
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} hpackage.Package
// @Router /package/all/mobile/{isHome} [get]
func (cr *PackageController) GetAllAppPackages(c fiber.Ctx) error {
    isHome := c.Params("isHome")
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.packageService.ListApp(isHome == "1", page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// SearchAllPackages
//
// @Tags Package
// @Produce json
// @Security BearerAuth
// @Param         _page        query       string                false  "_page"  default:"1"
// @Param         _limit       query       string                false  "_limit" default:"10"
// @Param         keyword      body        dto.SearchKeyword3Dto false  "Search"
// @Success 200 {array} hpackage.Package
// @Router /package/all [post]
func (cr *PackageController) SearchAllPackages(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key2 := data.GetString("keyword2")
    key3 := data.GetString("keyword3")
    if key != "" {
        key = "%" + key + "%"
    }
    if key2 != "" {
        if _, err := goment.New(key2, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong start date format")
        }
    }
    if key3 != "" {
        if _, err := goment.New(key3, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong end date format")
        }
    }

    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.packageService.ListByKeyword(key, key2, key3, page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetPackageStatusById
//
// @Tags Package
// @Produce json
// @Security BearerAuth
// @Param         packageId    path        string                true  "packageId"
// @Success 200
// @Router /packageStatus/{packageId} [get]
func (cr *PackageController) GetPackageStatusById(c fiber.Ctx) error {
    packageId := c.Params("packageId")
    ipackageId, _ := strconv.ParseInt(packageId, 10, 64)
    o, err := cr.packageService.FindPackageStatusByPackageId(ipackageId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// GetPackageById
//
// @Tags Package
// @Produce json
// @Security BearerAuth
// @Param         packageId    path        string                true  "packageId"
// @Success 200
// @Router /package/{packageId} [get]
func (cr *PackageController) GetPackageById(c fiber.Ctx) error {
    packageId := c.Params("packageId")
    ipackageId, _ := strconv.ParseInt(packageId, 10, 64)
    o, err := cr.packageService.FindByPackageId(ipackageId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// CreatePackage
//
// @Tags Package
// @Produce json
// @Security BearerAuth
// @Param request body dto.PackageDto true "PackageDto"
// @Success 200
// @Router /package [post]
func (cr *PackageController) CreatePackage(c fiber.Ctx) error {
    adminId, _, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    data := new(dto.PackageDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    resizeValidation, err := resizeBase64Image(data.PackageImage)
    if err != nil {
        return err
    }

    o := model.Package{
        PackageCode:           utils.NewNullString(data.PackageCode),
        PackageName:           utils.NewNullString(data.PackageName),
        PackageDesc:           utils.NewNullString(data.PackageDesc),
        PackageImage:          utils.NewNullString(data.PackageImage),
        ResizePackageImage:    resizeValidation,
        PackageStartDateTime:  utils.NewNullString(data.PackageStartDateTime),
        PackageEndDateTime:    utils.NewNullString(data.PackageEndDateTime),
        PackageValidityType:   utils.NewNullString(data.PackageValidityType),
        PackageValidity:       utils.NewInt64(int64(data.PackageValidity)),
        PackageValidityDate:   utils.NewNullString(data.PackageValidityDate),
        PackageTnc:            utils.NewNullString(data.PackageTnc),
        PackagePrice:          utils.NewFloat(data.PackagePrice),
        PackageMaxPurchase:    utils.NewInt32(int32(data.PackageMaxPurchase)),
        PackageAssignedDoctor: utils.NewInt64(int64(data.PackageAssignedDoctor)),
        PackageAllowAppt:      utils.NewNullString(data.PackageAllowAppt),
        PackageDisplayOrder:   utils.NewInt32(int32(data.PackageDisplayOrder)),
        PackageExtLink:        utils.NewNullString(data.PackageExtLink),
    }

    err = cr.packageService.Save(o, adminId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Package created",
    })
}

// UpdatePackage
//
// @Tags Package
// @Produce json
// @Security BearerAuth
// @Param request body dto.PackageDto true "PackageDto"
// @Success 200
// @Router /package/{packageId} [put]
func (cr *PackageController) UpdatePackage(c fiber.Ctx) error {
    adminId, _, err := middleware.ValidateAdminToken(c)
    if err != nil {
        return err
    }

    data := new(dto.PackageDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    packageId := c.Params("packageId")
    ipackageId, _ := strconv.ParseInt(packageId, 10, 64)

    resizeValidation, err := resizeBase64Image(data.PackageImage)
    if err != nil {
        return err
    }

    o := model.Package{
        PackageId:             utils.NewInt64(ipackageId),
        PackageCode:           utils.NewNullString(data.PackageCode),
        PackageName:           utils.NewNullString(data.PackageName),
        PackageDesc:           utils.NewNullString(data.PackageDesc),
        PackageImage:          utils.NewNullString(data.PackageImage),
        ResizePackageImage:    resizeValidation,
        PackageStartDateTime:  utils.NewNullString(data.PackageStartDateTime),
        PackageEndDateTime:    utils.NewNullString(data.PackageEndDateTime),
        PackageValidityType:   utils.NewNullString(data.PackageValidityType),
        PackageValidity:       utils.NewInt64(int64(data.PackageValidity)),
        PackageValidityDate:   utils.NewNullString(data.PackageValidityDate),
        PackageTnc:            utils.NewNullString(data.PackageTnc),
        PackagePrice:          utils.NewFloat(data.PackagePrice),
        PackageMaxPurchase:    utils.NewInt32(int32(data.PackageMaxPurchase)),
        PackageAssignedDoctor: utils.NewInt64(int64(data.PackageAssignedDoctor)),
        PackageAllowAppt:      utils.NewNullString(data.PackageAllowAppt),
        PackageDisplayOrder:   utils.NewInt32(int32(data.PackageDisplayOrder)),
        PackageExtLink:        utils.NewNullString(data.PackageExtLink),
    }

    err = cr.packageService.Update(o, adminId)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Package updated",
    })
}

// GetAllExportHospitalPackage
//
// @Tags Package
// @Produce json
// @Security BearerAuth
// @Success 200 {array} hpackage.Package
// @Router /package/export/all [get]
func (cr *PackageController) GetAllExportHospitalPackage(c fiber.Ctx) error {
    lx, err := cr.exportExcelService.ExportAllHospitalPackage()
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

// GetSearchExportHospitalPackage
//
// @Tags Package
// @Produce json
// @Security BearerAuth
// @Param         keyword      body        dto.SearchKeyword3Dto false  "Search"
// @Success 200 {array} hpackage.Package
// @Router /package/export/search [post]
func (cr *PackageController) GetSearchExportHospitalPackage(c fiber.Ctx) error {
    var data utils.Map
    if err := c.Bind().Body(&data); err != nil {
        return err
    }

    key := data.GetString("keyword")
    key2 := data.GetString("keyword2")
    key3 := data.GetString("keyword3")
    if key != "" {
        key = "%" + key + "%"
    }
    if key2 != "" {
        if _, err := goment.New(key2, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong start date format")
        }
    }
    if key3 != "" {
        if _, err := goment.New(key3, "DD/MM/YYYY"); err != nil {
            return fiber.NewError(fiber.StatusBadRequest, "Wrong end date format")
        }
    }

    lx, err := cr.exportExcelService.ExportHospitalPackageByKeyword(key, key2, key3)
    if err != nil {
        return err
    }

    return c.JSON(lx)
}

func resizeBase64Image(base64s string) (string, error) {
    i := strings.Index(base64s, "base64,")

    if i < 0 {
        base64Data := base64s
        imgBytes, err := base64.StdEncoding.DecodeString(base64Data)
        if err != nil {
            return "", err
        }

        srcImg, _, err := image.Decode(bytes.NewReader(imgBytes))
        if err != nil {
            return "", err
        }

        if len(imgBytes) > 1000000 {
            resizedImg := resize.Resize(500, 0, srcImg, resize.Lanczos3)
            var buf bytes.Buffer
            if err := png.Encode(&buf, resizedImg); err != nil {
                return "", err
            }
            return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
        }
    } else {
        base64Data := base64s[i+7:]
        m := base64s[:i+7]
        imgBytes, err := base64.StdEncoding.DecodeString(base64Data)
        if err != nil {
            return "", err
        }

        srcImg, _, err := image.Decode(bytes.NewReader(imgBytes))
        if err != nil {
            return "", err
        }

        resizedImg := resize.Resize(500, 0, srcImg, resize.Lanczos3)
        var buf bytes.Buffer
        if err := png.Encode(&buf, resizedImg); err != nil {
            return "", err
        }
        s := fmt.Sprintf("%s%s", m, base64.StdEncoding.EncodeToString(buf.Bytes()))
        return s, nil
    }

    return base64s, nil
}
