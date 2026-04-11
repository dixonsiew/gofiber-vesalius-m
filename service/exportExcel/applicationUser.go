package exportExcel

import (
    "fmt"
    "strings"
    "vesaliusm/model"
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
