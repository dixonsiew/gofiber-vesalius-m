package exportExcel

import (
	"vesaliusm/model"

	"github.com/nleeper/goment"
)

func (s *ExportExcelService) setMobileUserAuditLog(o *model.MobileUserAuditLog) {
	if o.DateCreate.Valid {
		g, _ := goment.New(o.DateCreate.String, "YYYY-MM-DD[T]HH:mm:ssZ")
		o.DateCreateExcel = g.Format("DD/MM/YYYY HH:mm")
	}
}
