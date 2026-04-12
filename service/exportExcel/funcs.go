package exportExcel

import (
    "fmt"
    "encoding/json"
    "strings"
    "vesaliusm/model"
    "vesaliusm/model/clubs"
    hp "vesaliusm/model/hpackage"
    lg "vesaliusm/model/logistic"
    upck "vesaliusm/model/userPackage"
    "vesaliusm/utils"

    "github.com/nleeper/goment"
)

func (s *ExportExcelService) setMobileUser(o *model.ApplicationUser) {
    if o.RegisterDateExcel != "" {
        g, _ := goment.New(o.RegisterDateExcel, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.RegisterDateExcel = g.Format("DD/MM/YYYY HH:mm")
    }

    o.FullName = utils.NewNullString(strings.TrimSpace(fmt.Sprintf("%s %s %s", o.FirstName.String, o.MiddleName.String, o.LastName.String)))
}

func (s *ExportExcelService) setHospitalPackage(o *hp.Package) error {
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

func (s *ExportExcelService) setPurchaseHistory(o *upck.UserPackage) {
    if o.BillingContactCode.Valid && o.BillingContactNo.Valid {
        o.BillingFullContactExcel = fmt.Sprintf("%s%s", o.BillingContactCode.String, o.BillingContactNo.String) 
    }
    if o.PurchasedDateTime.Valid && o.PurchasedDateTime.String != "-" {
        g, _ := goment.New(o.PurchasedDateTime.String)
        o.PurchasedDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.PurchasedDateTimeExcel = o.PurchasedDateTime.String
    }
    if o.ExpiredDateTime.Valid && o.ExpiredDateTime.String != "-" {
        g, _ := goment.New(o.ExpiredDateTime.String)
        o.ExpiredDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.ExpiredDateTimeExcel = o.ExpiredDateTime.String
    }
    if o.OrderedDateTime.Valid && o.OrderedDateTime.String != "-" {
        g, _ := goment.New(o.OrderedDateTime.String)
        o.OrderedDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.OrderedDateTimeExcel = o.OrderedDateTime.String
    }
    if o.BookedDateTime.Valid && o.BookedDateTime.String != "-" {
        g, _ := goment.New(o.BookedDateTime.String)
        o.BookedDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.BookedDateTimeExcel = o.BookedDateTime.String
    }
    if o.RedeemedDateTime.Valid && o.RedeemedDateTime.String != "-" {
        g, _ := goment.New(o.RedeemedDateTime.String)
        o.RedeemedDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.RedeemedDateTimeExcel = o.RedeemedDateTime.String
    }
    if o.CancelledDateTime.Valid && o.CancelledDateTime.String != "-" {
        g, _ := goment.New(o.CancelledDateTime.String)
        o.CancelledDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.CancelledDateTimeExcel = o.CancelledDateTime.String
    }
    if o.PaymentTransDate.Valid && o.PaymentTransDate.String != "-" {
        g, _ := goment.New(o.PaymentTransDate.String)
        o.PaymentTransDateExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.PaymentTransDateExcel = o.PaymentTransDate.String
    }
}

func (s *ExportExcelService) setLogisticRequest(o *lg.LogisticRequest) {
    if o.RequestedPickupDate.Valid {
        g, _ := goment.New(o.RequestedPickupDate.String)
        o.RequestedPickupDateExcel = g.Format("DD/MM/YYYY")
    }
    if o.FlightArrivalDate.Valid {
        g, _ := goment.New(o.FlightArrivalDate.String)
        o.FlightArrivalDateExcel = g.Format("DD/MM/YYYY")
    }
    if o.DateCreate.Valid {
        g, _ := goment.New(o.DateCreate, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.DateCreateExcel = g.Format("DD/MM/YYYY HH:mm")
    }
    if o.UserDateUpdate.Valid && o.UserDateUpdate.String != "-" {
        g, _ := goment.New(o.UserDateUpdate, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.UserDateUpdateExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.UserDateUpdateExcel = o.UserDateUpdate.String
    }
    if o.AdminDateUpdate.Valid && o.AdminDateUpdate.String != "-" {
        g, _ := goment.New(o.AdminDateUpdate, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.AdminDateUpdateExcel = g.Format("DD/MM/YYYY HH:mm")
    } else {
        o.AdminDateUpdateExcel = o.AdminDateUpdate.String
    }
}

func (s *ExportExcelService) setMobileUserAuditLog(o *model.MobileUserAuditLog) {
	if o.DateCreate.Valid {
		g, _ := goment.New(o.DateCreate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
		o.DateCreateExcel = g.Format("DD/MM/YYYY HH:mm")
	}
}

func (s *ExportExcelService) setLittleKidsMembership(o *clubs.LittleExplorersKidsMembership) {
    if o.ActivityJoinDate.Valid {
		g, _ := goment.New(o.ActivityJoinDate.String)
		o.ActivityJoinDateExcel = g.Format("DD/MM/YYYY HH:mm")
	}
    if o.KidsMembershipJoinDate.Valid {
		g, _ := goment.New(o.KidsMembershipJoinDate.String)
		o.KidsMembershipJoinDateExcel = g.Format("DD/MM/YYYY HH:mm")
	}
    if o.KidsDob.Valid {
		g, _ := goment.New(o.KidsDob.String)
		o.KidsDobExcel = g.Format("DD/MM/YYYY")
	}
    if o.GuardianDob.Valid && o.GuardianDob.String != "-" {
		g, _ := goment.New(o.GuardianDob.String)
		o.GuardianDobExcel = g.Format("DD/MM/YYYY")
	} else {
        o.GuardianDobExcel = o.GuardianDob.String
    }
}

func (s *ExportExcelService) setGoldenPearlMembership(o *clubs.GoldenPearlMembership) {
    if o.ActivityJoinDate.Valid {
		g, _ := goment.New(o.ActivityJoinDate.String)
		o.ActivityJoinDateExcel = g.Format("DD/MM/YYYY HH:mm")
	}
    if o.GoldenMembershipJoinDate.Valid {
		g, _ := goment.New(o.GoldenMembershipJoinDate.String)
		o.GoldenMembershipJoinDateExcel = g.Format("DD/MM/YYYY HH:mm")
	}
    if o.GoldenDob.Valid {
		g, _ := goment.New(o.GoldenDob.String)
		o.GoldenDobExcel = g.Format("DD/MM/YYYY")
	}
    if o.NokDob.Valid && o.NokDob.String != "-" {
		g, _ := goment.New(o.NokDob.String)
		o.NokDobExcel = g.Format("DD/MM/YYYY")
	} else {
        o.NokDobExcel = o.NokDob.String
    }
}

func (s *ExportExcelService) setLittleKidsActivity(o *clubs.LittleExplorersKidsActivity) {
    if o.ActivityStartDateTime.Valid {
        g, _ := goment.New(o.ActivityStartDateTime, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.ActivityStartDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    }
    if o.ActivityEndDateTime.Valid && o.ActivityEndDateTime.String != "-" {
		g, _ := goment.New(o.ActivityEndDateTime.String, "YYYY-MM-DD[T]HH:mm:ssZ")
		o.ActivityEndDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
	} else {
        o.ActivityEndDateTimeExcel = "-"
    }
}

func (s *ExportExcelService) setGoldenPearlActivity(o *clubs.GoldenPearlActivity) {
    if o.ActivityStartDateTime.Valid {
        g, _ := goment.New(o.ActivityStartDateTime, "YYYY-MM-DD[T]HH:mm:ssZ")
        o.ActivityStartDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
    }
    if o.ActivityEndDateTime.Valid && o.ActivityEndDateTime.String != "-" {
		g, _ := goment.New(o.ActivityEndDateTime.String, "YYYY-MM-DD[T]HH:mm:ssZ")
		o.ActivityEndDateTimeExcel = g.Format("DD/MM/YYYY HH:mm")
	} else {
        o.ActivityEndDateTimeExcel = "-"
    }
}

func (s *ExportExcelService) setLittleKidsAttendees(o *clubs.LittleExplorersKidsMembership) {
    if o.ActivityJoinDate.Valid {
		g, _ := goment.New(o.ActivityJoinDate.String)
		o.ActivityJoinDateExcel = g.Format("DD/MM/YYYY HH:mm")
	}
    if o.KidsMembershipJoinDate.Valid {
		g, _ := goment.New(o.KidsMembershipJoinDate.String)
		o.KidsMembershipJoinDateExcel = g.Format("DD/MM/YYYY HH:mm")
	}
    if o.KidsDob.Valid {
		g, _ := goment.New(o.KidsDob.String)
		o.KidsDobExcel = g.Format("DD/MM/YYYY")
	}
}

func (s *ExportExcelService) setGoldenPearlAttendees(o *clubs.GoldenPearlMembership) {
    if o.ActivityJoinDate.Valid {
		g, _ := goment.New(o.ActivityJoinDate.String)
		o.ActivityJoinDateExcel = g.Format("DD/MM/YYYY HH:mm")
	}
    if o.GoldenMembershipJoinDate.Valid {
		g, _ := goment.New(o.GoldenMembershipJoinDate.String)
		o.GoldenMembershipJoinDateExcel = g.Format("DD/MM/YYYY HH:mm")
	}
    if o.GoldenDob.Valid {
		g, _ := goment.New(o.GoldenDob.String)
		o.GoldenDobExcel = g.Format("DD/MM/YYYY")
	}
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

func getPatientPurchaseDetailsCols(prefix string) string {
    lx := []string{
        "PATIENT_PURCHASE_ID",
        "PATIENT_PRN",
        "PATIENT_NAME",
        "PACKAGE_ID",
        "PACKAGE_PURCHASE_NO",
        "PACKAGE_STATUS",
        "ORDERED_DATETIME",
        "BOOKED_DATETIME",
        "REDEEMED_DATETIME",
        "CANCELLED_DATETIME",
        "PURCHASED_DATETIME",
        "EXPIRED_DATETIME",
    }
    ls := make([]string, 0)
    for _, v := range lx {
        ls = append(ls, prefix+v)
    }
    return strings.Join(ls, ", ")
}
