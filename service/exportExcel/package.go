package exportExcel

import (
    "encoding/json"
    model "vesaliusm/model/hpackage"
    "vesaliusm/utils"

    "github.com/nleeper/goment"
)

func (s *ExportExcelService) setHospitalPackage(o *model.Package) error {
    if o.PackageStartDateTime.Valid {
        g, _ := goment.New(o.PackageStartDateTime.String)
        o.PackageStartDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    }
    if o.PackageEndDateTime.Valid && o.PackageEndDateTime.String != "-" {
        g, _ := goment.New(o.PackageEndDateTime.String)
        o.PackageEndDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.PackageEndDateTimeExcel = "-"
    }
    if o.PackageAssignedDoctor.Valid {
        doctorName, err := s.novaDoctorService.FindDoctorNameByDoctorId(o.PackageAssignedDoctor.Int64, s.db)
        if err != nil {
            return err
        } else {
            o.PackageAssignedDoctorName = doctorName
        }
    }
    if o.PackageValidityType.String == "Date" && o.PackageValidityDate.Valid {
        g, _ := goment.New(o.PackageValidityDate.String)
        o.PackageValidityDateExcel = g.Format("DD/MM/YYYY")
    } else {
        o.PackageValidityDateExcel = "-"
    }
    if !o.PackageValidity.Valid {
        o.PackageValidityExcel = "-"
    } else {
        jsonBytes, err := json.Marshal(o.PackageValidity)
        if err != nil {
            return err
        } else {
            o.PackageValidityExcel = string(jsonBytes)
        }
    }
    if !o.PackageTnc.Valid {
        o.PackageTnc = utils.NewNullString("-")
    }
    if !o.PackageExtLink.Valid {
        o.PackageExtLink = utils.NewNullString("-")
    }

    return nil
}
