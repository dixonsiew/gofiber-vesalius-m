package exportExcel

import (
	model "vesaliusm/model/logistic"

	"github.com/nleeper/goment"
)

func (s *ExportExcelService) setLogisticRequest(o *model.LogisticRequest) {
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
