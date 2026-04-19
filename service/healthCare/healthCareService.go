package healthCare

import (
    "context"
    "database/sql"
    "encoding/base64"
    "strconv"
    "strings"
    "vesaliusm/config"
    "vesaliusm/database"
    gm "vesaliusm/model"
    model "vesaliusm/model/healthCare"
    upck "vesaliusm/model/userPackage"
    "vesaliusm/service/adminUser"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/applicationUserFamily"
    "vesaliusm/service/branch"
    "vesaliusm/service/novaBill"
    "vesaliusm/service/novaDoctor"
    "vesaliusm/service/novaDoctorPatientAppt"
    "vesaliusm/service/novaHealthScreenRpt"
    "vesaliusm/service/novaHealthScreenRptDetail"
    "vesaliusm/service/novaInvestigation"
    "vesaliusm/service/novaInvestigationDetail"
    "vesaliusm/service/novaPatient"
    "vesaliusm/service/novaPatientAlert"
    "vesaliusm/service/novaPatientrx"
    "vesaliusm/service/novaReferralLetter"
    "vesaliusm/service/novaVisit"
    "vesaliusm/service/novaVisitSummary"
    "vesaliusm/service/novaVitalSigns"
    "vesaliusm/service/novaVitalSignsDetail"
    sqx "vesaliusm/sql"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/jmoiron/sqlx"
    "github.com/nleeper/goment"
)

var HealthCareSvc *HealthCareService = NewHealthCareService(database.GetDb(), database.GetCtx(), database.GetDbrs())

type HealthCareService struct {
    db   *sqlx.DB
    ctx  context.Context
    dbrs *sqlx.DB

    vitalCodeWeight                 string
    vitalCodeHeight                 string
    vitalCodeBMI                    string
    vitalCodeBP                     string
    vitalCodePulse                  string
    vitalCodeClinicaltemplateWeight string
    vitalCodeClinicaltemplateHeight string
    vitalCodeClinicaltemplateBMI    string
    vitalCodeClinicaltemplatePulse  string
    labInvestigationType            string
    labCodeHDL                      string
    labCodeLDL                      string
    labCodeGlucose                  string
    labCodeHemoglobin               string

    adminUserService                 *adminUser.AdminUserService
    applicationUserService           *applicationUser.ApplicationUserService
    applicationUserFamilyService     *applicationUserFamily.ApplicationUserFamilyService
    branchService                    *branch.BranchService
    novaDoctorService                *novaDoctor.NovaDoctorService
    novaDoctorPatientApptService     *novaDoctorPatientAppt.NovaDoctorPatientApptService
    novaVisitService                 *novaVisit.NovaVisitService
    novaPatientAlertService          *novaPatientAlert.NovaPatientAlertService
    novaPatientService               *novaPatient.NovaPatientService
    novaBillService                  *novaBill.NovaBillService
    novaVisitSummaryService          *novaVisitSummary.NovaVisitSummaryService
    novaInvestigationService         *novaInvestigation.NovaInvestigationService
    novaInvestigationDetailService   *novaInvestigationDetail.NovaInvestigationDetailService
    novaVitalSignsService            *novaVitalSigns.NovaVitalSignsService
    novaVitalSignsDetailService      *novaVitalSignsDetail.NovaVitalSignsDetailService
    novaReferralLetterService        *novaReferralLetter.NovaReferralLetterService
    novaHealthScreenRptService       *novaHealthScreenRpt.NovaHealthScreenRptService
    novaHealthScreenRptDetailService *novaHealthScreenRptDetail.NovaHealthScreenRptDetailService
    novaPatientrxService             *novaPatientrx.NovaPatientrxService
}

func NewHealthCareService(db *sqlx.DB, ctx context.Context, dbrs *sqlx.DB) *HealthCareService {
    return &HealthCareService{
        db:   db,
        ctx:  ctx,
        dbrs: dbrs,

        vitalCodeWeight:                 config.Config("vital.code.weight"),
        vitalCodeHeight:                 config.Config("vital.code.height"),
        vitalCodeBMI:                    config.Config("vital.code.bmi"),
        vitalCodeBP:                     config.Config("vital.code.bp"),
        vitalCodePulse:                  config.Config("vital.code.pulse"),
        vitalCodeClinicaltemplateWeight: config.Config("vital.code.weight.clinicaltemplate"),
        vitalCodeClinicaltemplateHeight: config.Config("vital.code.height.clinicaltemplate"),
        vitalCodeClinicaltemplateBMI:    config.Config("vital.code.bmi.clinicaltemplate"),
        vitalCodeClinicaltemplatePulse:  config.Config("vital.code.pulse.clinicaltemplate"),
        labInvestigationType:            config.Config("lab.investigation.type"),
        labCodeHDL:                      config.Config("lab.code.hdl"),
        labCodeLDL:                      config.Config("lab.code.ldl"),
        labCodeGlucose:                  config.Config("lab.code.glucose"),
        labCodeHemoglobin:               config.Config("lab.code.hemoglobin"),

        adminUserService:                 adminUser.NewAdminUserService(db, ctx),
        applicationUserService:           applicationUser.NewApplicationUserService(db, ctx),
        applicationUserFamilyService:     applicationUserFamily.NewApplicationUserFamilyService(db, ctx, dbrs),
        branchService:                    branch.NewBranchService(db, ctx),
        novaDoctorService:                novaDoctor.NewNovaDoctorService(db, ctx),
        novaDoctorPatientApptService:     novaDoctorPatientAppt.NewNovaDoctorPatientApptService(db, ctx),
        novaVisitService:                 novaVisit.NewNovaVisitService(db, ctx),
        novaPatientAlertService:          novaPatientAlert.NewNovaPatientAlertService(db, ctx),
        novaPatientService:               novaPatient.NewNovaPatientService(db, ctx),
        novaBillService:                  novaBill.NewNovaBillService(db, ctx),
        novaVisitSummaryService:          novaVisitSummary.NewNovaVisitSummaryService(db, ctx),
        novaInvestigationService:         novaInvestigation.NewNovaInvestigationService(db, ctx),
        novaInvestigationDetailService:   novaInvestigationDetail.NewNovaInvestigationDetailService(db, ctx),
        novaVitalSignsService:            novaVitalSigns.NewNovaVitalSignsService(db, ctx),
        novaVitalSignsDetailService:      novaVitalSignsDetail.NewNovaVitalSignsDetailService(db, ctx),
        novaReferralLetterService:        novaReferralLetter.NewNovaReferralLetterService(db, ctx),
        novaHealthScreenRptService:       novaHealthScreenRpt.NewNovaHealthScreenRptService(db, ctx),
        novaHealthScreenRptDetailService: novaHealthScreenRptDetail.NewNovaHealthScreenRptDetailService(db, ctx),
        novaPatientrxService:             novaPatientrx.NewNovaPatientrxService(db, ctx),
    }
}

func (s *HealthCareService) GetPatientVisit(prn string, pageId string) ([]model.NovaVisitDetails, error) {
    novaVisitDetails := make([]model.NovaVisitDetails, 0)
    visitHistoryType := pageId
    switch visitHistoryType {
    case "1":
        novaVisits, err := s.novaVisitService.GetSpecificPatientVisit(prn, s.dbrs)
        if err != nil {
            return nil, err
        }
        if len(novaVisits) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for i := range novaVisits {
            novaVisit := novaVisits[i]
            novaVisitSummaries, err := s.novaVisitSummaryService.FindByAccountNoAndCategoryNotAndCategoryNot(novaVisit.AccountNo.String, "Investigation", "Medication", s.dbrs)
            if err != nil {
                return nil, err
            }
            if len(novaVisitSummaries) > 0 {
                for range novaVisitSummaries {
                    novaVisitDetail := model.NovaVisitDetails{
                        NovaVisit:          &novaVisit,
                        NovaVisitSummaries: novaVisitSummaries,
                    }
                    novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
                }
            }
        }

    case "2":
        novaVisitsPrescription, err := s.novaVisitService.GetSpecificPatientVisit(prn, s.dbrs)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsPrescription) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for i := range novaVisitsPrescription {
            novaVisit := novaVisitsPrescription[i]
            novaVisitPatientRxList, err := s.novaPatientrxService.FindPatientRxByAccountNo(novaVisit.AccountNo.String, s.dbrs)
            if err != nil {
                return nil, err
            }
            if len(novaVisitPatientRxList) > 0 {
                for range novaVisitPatientRxList {
                    novaVisitDetail := model.NovaVisitDetails{
                        NovaVisit:              &novaVisit,
                        NovaVisitPatientRxList: novaVisitPatientRxList,
                    }
                    novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
                }
            }
        }

    case "3":
        novaVisitsInvestigation, err := s.novaVisitService.GetSpecificPatientVisit(prn, s.dbrs)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsInvestigation) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for i := range novaVisitsInvestigation {
            novaVisit := novaVisitsInvestigation[i]
            novaPatientInvestigationList, err := s.novaInvestigationService.FindNonPanelInvestiationByAccountNo(novaVisit.AccountNo.String, s.dbrs)
            if err != nil {
                return nil, err
            }
            uniquePanelList, err := s.novaInvestigationService.FindUniquePanelInvestigationByAccountNo(novaVisit.AccountNo.String, s.dbrs)
            if err != nil {
                return nil, err
            }
            if len(novaPatientInvestigationList) > 0 || len(uniquePanelList) > 0 {
                for range novaPatientInvestigationList {
                    novaVisitDetail := model.NovaVisitDetails{
                        NovaVisit: &novaVisit,
                    }
                    novaVisitInvestigationDetailList := make([]model.NovaVisitInvestigationDetail, 0)
                    if len(novaPatientInvestigationList) > 0 {
                        for _, novaPatientInvestigation := range novaPatientInvestigationList {
                            novaVisitInvestigationDetailDto := model.NovaVisitInvestigationDetail{}
                            novaPatientInvestigationDetailList, err := s.novaInvestigationDetailService.GetInvestigationRefNoAndCode(novaPatientInvestigation.InvestigationRefNo.String, novaPatientInvestigation.Code.String, s.dbrs)
                            if err != nil {
                                return nil, err
                            }
                            if len(novaPatientInvestigationDetailList) > 0 {
                                for _, novaPatientInvestigationDetail := range novaPatientInvestigationDetailList {
                                    novaVisitInvestigationDetailDto.InvestigationRefNo = novaPatientInvestigationDetail.InvestigationRefNo.String
                                    novaVisitInvestigationDetailDto.InvestigationType = novaPatientInvestigation.InvestigationType.String
                                    novaVisitInvestigationDetailDto.AccountNo = novaPatientInvestigation.AccountNo.String
                                    novaVisitInvestigationDetailDto.Code = novaPatientInvestigationDetail.Code.String
                                    novaVisitInvestigationDetailDto.Description = novaPatientInvestigation.Description.String
                                    novaVisitInvestigationDetailDto.PanelCode = novaPatientInvestigation.PanelCode.String
                                    novaVisitInvestigationDetailDto.PanelDescription = novaPatientInvestigation.PanelDescription.String
                                    if novaPatientInvestigationDetail.ResultValue.Valid && novaPatientInvestigationDetail.ResultValue.String != "null" {
                                        novaVisitInvestigationDetailDto.ResultValue = novaPatientInvestigationDetail.ResultValue.String
                                        novaVisitInvestigationDetailDto.ResultUnit = novaPatientInvestigationDetail.ResultUnit.String
                                    } else {
                                        novaVisitInvestigationDetailDto.ResultValue = novaPatientInvestigationDetail.ResultClob.String
                                        novaVisitInvestigationDetailDto.ResultUnit = ""
                                    }
                                    novaVisitInvestigationDetailDto.ReferenceRange = novaPatientInvestigationDetail.ReferenceRange.String
                                    novaVisitInvestigationDetailDto.RangeType = novaPatientInvestigationDetail.RangeType.String
                                    novaVisitInvestigationDetailDto.ResultClob = novaPatientInvestigationDetail.ResultClob.String
                                }
                            }
                            novaVisitInvestigationDetailList = append(novaVisitInvestigationDetailList, novaVisitInvestigationDetailDto)
                        }
                    }

                    if len(uniquePanelList) > 0 {
                        for i := range uniquePanelList {
                            uniquePanelCode := uniquePanelList[i]
                            novaPatientPanelInvestigationList, err := s.novaInvestigationService.FindPanelInvestigationByAccountNoPanelCode(novaVisit.AccountNo.String, uniquePanelCode, s.dbrs)
                            if err != nil {
                                return nil, err
                            }
                            if len(novaPatientPanelInvestigationList) > 0 {
                                novaVisitInvestigationDetailDto := model.NovaVisitInvestigationDetail{}
                                panelDetailList := make([]model.NovaVisitInvestigationPanelDetail, 0)
                                for j := range novaPatientPanelInvestigationList {
                                    novaPatientInvestigation := novaPatientPanelInvestigationList[j]
                                    novaPatientInvestigationDetailList, err := s.novaInvestigationDetailService.GetInvestigationRefNoAndCode(novaPatientInvestigation.InvestigationRefNo.String, novaPatientInvestigation.Code.String, s.db)
                                    if err != nil {
                                        return nil, err
                                    }
                                    if len(novaPatientInvestigationDetailList) > 0 {
                                        for k := range novaPatientInvestigationDetailList {
                                            novaPatientInvestigationDetail := novaPatientInvestigationDetailList[k]
                                            panelDetail := model.NovaVisitInvestigationPanelDetail{
                                                Code:        novaPatientInvestigationDetail.Code.String,
                                                Description: novaPatientInvestigation.Description.String,
                                            }
                                            if novaPatientInvestigationDetail.ResultValue.Valid && novaPatientInvestigationDetail.ResultValue.String != "null" {
                                                panelDetail.ResultValue = novaPatientInvestigationDetail.ResultValue.String
                                                panelDetail.ResultUnit = novaPatientInvestigationDetail.ResultUnit.String
                                            } else {
                                                panelDetail.ResultValue = novaPatientInvestigationDetail.ResultClob.String
                                                panelDetail.ResultUnit = ""
                                            }
                                            panelDetail.ReferenceRange = novaPatientInvestigationDetail.ReferenceRange.String
                                            panelDetail.RangeType = novaPatientInvestigationDetail.RangeType.String
                                            panelDetail.ResultClob = novaPatientInvestigationDetail.ResultClob.String
                                            panelDetailList = append(panelDetailList, panelDetail)
                                        }
                                    }
                                }
                                novaVisitInvestigationDetailDto.InvestigationRefNo = novaPatientPanelInvestigationList[0].InvestigationRefNo.String
                                novaVisitInvestigationDetailDto.InvestigationType = novaPatientPanelInvestigationList[0].InvestigationType.String
                                novaVisitInvestigationDetailDto.AccountNo = novaPatientPanelInvestigationList[0].AccountNo.String
                                novaVisitInvestigationDetailDto.PanelCode = novaPatientPanelInvestigationList[0].PanelCode.String
                                novaVisitInvestigationDetailDto.PanelDescription = novaPatientPanelInvestigationList[0].PanelDescription.String
                                novaVisitInvestigationDetailDto.PanelDetail = panelDetailList
                                novaVisitInvestigationDetailList = append(novaVisitInvestigationDetailList, novaVisitInvestigationDetailDto)
                            }
                        }
                    }
                    novaVisitDetail.NovaVisitInvestigationDetailList = novaVisitInvestigationDetailList
                    novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
                }
            }
        }

    case "4":
        novaVisitsBill, err := s.novaVisitService.GetSpecificPatientVisit(prn, s.dbrs)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsBill) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for i := range novaVisitsBill {
            novaVisit := novaVisitsBill[i]
            novaBills, err := s.novaBillService.GetNovaBillByPrnAndAccountNo(prn, novaVisit.AccountNo.String, s.dbrs)
            if err != nil {
                return nil, err
            }
            if len(novaBills) > 0 {
                novaVisitDetail := model.NovaVisitDetails{
                    NovaVisit: &novaVisit,
                    NovaBills: novaBills,
                }
                novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
            }
        }

    case "5":
        novaVisitsVital, err := s.novaVisitService.GetPatientVisitWithVitalSign(prn, s.dbrs)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsVital) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for i := range novaVisitsVital {
            novaVisit := novaVisitsVital[i]
            novaPatientVitalSignsList, err := s.novaVitalSignsService.FindPatientVitalSignsByAccountNo(novaVisit.AccountNo.String, s.dbrs)
            if err != nil {
                return nil, err
            }
            novaVisitVitalSignsDetailList := make([]model.NovaVisitVitalSignsDetail, 0)
            if len(novaPatientVitalSignsList) > 0 {
                for j := range novaPatientVitalSignsList {
                    patientVitalSigns := novaPatientVitalSignsList[j]
                    patientVitalSignsDetailList, err := s.novaVitalSignsDetailService.GetLocalPatientVitalSignsDetailsByRefNo(
                        patientVitalSigns.RefNo.String,
                        s.vitalCodeHeight,
                        s.vitalCodeWeight,
                        s.vitalCodeBP,
                        s.vitalCodeBMI,
                        s.vitalCodePulse,
                        s.dbrs,
                    )
                    if err != nil {
                        return nil, err
                    }
                    if len(patientVitalSignsDetailList) > 0 {
                        for k := range patientVitalSignsDetailList {
                            eachPatientVitalSignsDetail := patientVitalSignsDetailList[k]
                            visitVitalSignsDetail := model.NovaVisitVitalSignsDetail{
                                Code:   eachPatientVitalSignsDetail.Code.String,
                                Desc:   eachPatientVitalSignsDetail.Description.String,
                                Value1: eachPatientVitalSignsDetail.Value1.String,
                                Value2: eachPatientVitalSignsDetail.Value2.String,
                                Unit:   eachPatientVitalSignsDetail.Unit.String,
                            }
                            novaVisitVitalSignsDetailList = append(novaVisitVitalSignsDetailList, visitVitalSignsDetail)
                        }
                    }
                }
                novaVisitDetail := model.NovaVisitDetails{
                    NovaVisit:                     &novaVisit,
                    NovaVisitVitalSignsDetailList: novaVisitVitalSignsDetailList,
                }
                novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
            }
        }

    case "6":
        novaVisitsWithReferralLetterList, err := s.novaVisitService.GetPatientVisitWithReferralLetter(prn, s.dbrs)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsWithReferralLetterList) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for i := range novaVisitsWithReferralLetterList {
            novaVisit := novaVisitsWithReferralLetterList[i]
            novaReferralLetterList, err := s.novaReferralLetterService.FindAllReferralLetterByPrnAndAccountNo(novaVisit.PRN.String, novaVisit.AccountNo.String, s.dbrs)
            if err != nil {
                return nil, err
            }
            novaVisitReferralLetterList := make([]model.NovaVisitReferralLetter, 0)
            if len(novaReferralLetterList) > 0 {
                for j := range novaReferralLetterList {
                    novaExtReferralLetter := novaReferralLetterList[j]
                    novaVisitReferralLetterDetail := model.NovaVisitReferralLetter{
                        ReferralRefNo:            novaExtReferralLetter.ReferralRefNo.String,
                        PRN:                      novaExtReferralLetter.PRN.String,
                        AccountNo:                novaExtReferralLetter.AccountNo.String,
                        ReferralDateTime:         novaExtReferralLetter.ReferralDateTime.String,
                        ReferrerDoctor:           novaExtReferralLetter.ReferrerDoctorMCR.String,
                        ReferralDoctor:           novaExtReferralLetter.ReferralTo.String,
                        ReferralTitleDept:        novaExtReferralLetter.ReferralTitleDept.String,
                        ReferralAddressOrSubject: novaExtReferralLetter.ReferralAddressOrSubject.String,
                        ReferralLetter:           novaExtReferralLetter.ReferralLetter.String,
                        ReferralType:             novaExtReferralLetter.ReferralType.String,
                    }
                    novaVisitReferralLetterList = append(novaVisitReferralLetterList, novaVisitReferralLetterDetail)
                }
            }
            novaVisitDetail := model.NovaVisitDetails{
                NovaVisit:                   &novaVisit,
                NovaVisitReferralLetterList: novaVisitReferralLetterList,
            }
            novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
        }

    case "7":
        novaVisitsWithHealthScreeningReportList, err := s.novaVisitService.GetPatientVisitWithHealthScreeningReport(prn, s.dbrs)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsWithHealthScreeningReportList) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for i := range novaVisitsWithHealthScreeningReportList {
            novaVisit := novaVisitsWithHealthScreeningReportList[i]
            novaHealthScreeningRptList, err := s.novaHealthScreenRptService.FindHealthScreeningRptByPrnAndAccountNo(novaVisit.PRN.String, novaVisit.AccountNo.String, s.dbrs)
            if err != nil {
                return nil, err
            }
            novaHealthScreeningRptDtoList := make([]model.NovaHealthScreeningRpt, 0)
            if len(novaHealthScreeningRptList) > 0 {
                for j := range novaHealthScreeningRptList {
                    novaHealthScreeningRpt := novaHealthScreeningRptList[j]
                    novaHealthScreeningRptDetail := model.NovaHealthScreeningRpt{
                        HsrRefNo:   utils.NewNullString(novaHealthScreeningRpt.HsrRefNo.String),
                        AccountNo:  utils.NewNullString(novaHealthScreeningRpt.AccountNo.String),
                        ReportDate: utils.NewNullString(novaHealthScreeningRpt.ReportDate.String),
                        ReportUser: utils.NewNullString(novaHealthScreeningRpt.ReportUser.String),
                    }
                    novaHealthScreeningRptDtoList = append(novaHealthScreeningRptDtoList, novaHealthScreeningRptDetail)
                }
            }
            novaVisitDetail := model.NovaVisitDetails{
                NovaVisit:                  &novaVisit,
                NovaHealthScreeningRptList: novaHealthScreeningRptDtoList,
            }
            novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
        }
    }
    return novaVisitDetails, nil
}

func (s *HealthCareService) GetVitalSignsHistoryForDashboard(prn string) ([]model.NovaVitalSignsDashboard, error) {
    vitalSignsDashboardList := make([]model.NovaVitalSignsDashboard, 0)
    patientVitalSignsHistoryList, err := s.novaVitalSignsDetailService.GetPatientVitalSignsHistoryForDashboard(prn, s.vitalCodeHeight, s.vitalCodeWeight, s.vitalCodeBP, s.vitalCodeBMI, s.vitalCodePulse, nil)
    if err != nil {
        return nil, err
    }
    if len(patientVitalSignsHistoryList) < 1 {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    vitalSignsGroupHeight := model.NovaVitalSignsDashboard{}
    vitalSignsGroupWeight := model.NovaVitalSignsDashboard{}
    vitalSignsGroupBMI := model.NovaVitalSignsDashboard{}
    vitalSignsGroupBP := model.NovaVitalSignsDashboard{}
    vitalSignsGroupPulse := model.NovaVitalSignsDashboard{}

    vitalSignsGroupHeight.VitalSignCode = "HEIGHT"
    vitalSignsGroupWeight.VitalSignCode = "WEIGHT"
    vitalSignsGroupBMI.VitalSignCode = "BMI"
    vitalSignsGroupBP.VitalSignCode = "BP"
    vitalSignsGroupPulse.VitalSignCode = "PULSE RATE"

    patientVitalSignsDetailHeightList := make([]model.NovaPatientVitalSignsDetailDto, 0)
    patientVitalSignsDetailWeightList := make([]model.NovaPatientVitalSignsDetailDto, 0)
    patientVitalSignsDetailBMIList := make([]model.NovaPatientVitalSignsDetailDto, 0)
    patientVitalSignsDetailBPList := make([]model.NovaPatientVitalSignsDetailDto, 0)
    patientVitalSignsDetailPulseList := make([]model.NovaPatientVitalSignsDetailDto, 0)

    for i := range patientVitalSignsHistoryList {
        vitalSignsHistoryDetail := patientVitalSignsHistoryList[i]
        if strings.EqualFold(vitalSignsHistoryDetail.Description.String, s.vitalCodeHeight) ||
            strings.EqualFold(vitalSignsHistoryDetail.Code.String, s.vitalCodeClinicaltemplateHeight) {
            patientVitalSignsDetailHeight := model.NovaPatientVitalSignsDetail{
                Code:         utils.NewNullString(vitalSignsHistoryDetail.Code.String),
                Description:  utils.NewNullString(vitalSignsHistoryDetail.Description.String),
                RefNo:        utils.NewNullString(vitalSignsHistoryDetail.RefNo.String),
                RecordedDate: utils.NewNullString(vitalSignsHistoryDetail.RecordedDate.String),
                Unit:         utils.NewNullString(vitalSignsHistoryDetail.Unit.String),
                Value1:       utils.NewNullString(vitalSignsHistoryDetail.Value1.String),
                Value2:       utils.NewNullString(vitalSignsHistoryDetail.Value2.String),
            }

            novaPatientVitalSignsDetailHeightDto := model.NovaPatientVitalSignsDetailDto{
                NovaPatientVitalSignsDetail: &patientVitalSignsDetailHeight,
            }
            patientVitalSignsDetailHeightList = append(patientVitalSignsDetailHeightList, novaPatientVitalSignsDetailHeightDto)
        }

        if strings.EqualFold(vitalSignsHistoryDetail.Description.String, s.vitalCodeWeight) ||
            strings.EqualFold(vitalSignsHistoryDetail.Code.String, s.vitalCodeClinicaltemplateWeight) {
            patientVitalSignsDetailWeight := model.NovaPatientVitalSignsDetail{
                Code:         utils.NewNullString(vitalSignsHistoryDetail.Code.String),
                Description:  utils.NewNullString(vitalSignsHistoryDetail.Description.String),
                RefNo:        utils.NewNullString(vitalSignsHistoryDetail.RefNo.String),
                RecordedDate: utils.NewNullString(vitalSignsHistoryDetail.RecordedDate.String),
                Unit:         utils.NewNullString(vitalSignsHistoryDetail.Unit.String),
                Value1:       utils.NewNullString(vitalSignsHistoryDetail.Value1.String),
                Value2:       utils.NewNullString(vitalSignsHistoryDetail.Value2.String),
            }

            novaPatientVitalSignsDetailWeightDto := model.NovaPatientVitalSignsDetailDto{
                NovaPatientVitalSignsDetail: &patientVitalSignsDetailWeight,
            }
            patientVitalSignsDetailWeightList = append(patientVitalSignsDetailWeightList, novaPatientVitalSignsDetailWeightDto)
        }

        if strings.EqualFold(vitalSignsHistoryDetail.Code.String, s.vitalCodeBP) {
            patientVitalSignsDetailBP := model.NovaPatientVitalSignsDetail{
                Code:         utils.NewNullString(vitalSignsHistoryDetail.Code.String),
                Description:  utils.NewNullString(vitalSignsHistoryDetail.Description.String),
                RefNo:        utils.NewNullString(vitalSignsHistoryDetail.RefNo.String),
                RecordedDate: utils.NewNullString(vitalSignsHistoryDetail.RecordedDate.String),
                Unit:         utils.NewNullString(vitalSignsHistoryDetail.Unit.String),
                Value1:       utils.NewNullString(vitalSignsHistoryDetail.Value1.String),
                Value2:       utils.NewNullString(vitalSignsHistoryDetail.Value2.String),
            }

            novaPatientVitalSignsDetailBPDto := model.NovaPatientVitalSignsDetailDto{
                NovaPatientVitalSignsDetail: &patientVitalSignsDetailBP,
                Value1LowValue:              "120",
                Value1HighValue:             "130",
                Value2LowValue:              "80",
                Value2HighValue:             "90",
            }
            patientVitalSignsDetailBPList = append(patientVitalSignsDetailBPList, novaPatientVitalSignsDetailBPDto)
        }

        if strings.EqualFold(vitalSignsHistoryDetail.Code.String, s.vitalCodeBMI) ||
            strings.EqualFold(vitalSignsHistoryDetail.Code.String, s.vitalCodeClinicaltemplateBMI) {
            patientVitalSignsDetailBMI := model.NovaPatientVitalSignsDetail{
                Code:         utils.NewNullString(vitalSignsHistoryDetail.Code.String),
                Description:  utils.NewNullString(vitalSignsHistoryDetail.Description.String),
                RefNo:        utils.NewNullString(vitalSignsHistoryDetail.RefNo.String),
                RecordedDate: utils.NewNullString(vitalSignsHistoryDetail.RecordedDate.String),
                Unit:         utils.NewNullString(vitalSignsHistoryDetail.Unit.String),
                Value1:       utils.NewNullString(vitalSignsHistoryDetail.Value1.String),
                Value2:       utils.NewNullString(vitalSignsHistoryDetail.Value2.String),
            }

            novaPatientVitalSignsDetailBMIDto := model.NovaPatientVitalSignsDetailDto{
                NovaPatientVitalSignsDetail: &patientVitalSignsDetailBMI,
                Value1LowValue:              "17",
                Value1HighValue:             "24",
            }
            patientVitalSignsDetailBMIList = append(patientVitalSignsDetailBMIList, novaPatientVitalSignsDetailBMIDto)
        }

        if strings.EqualFold(vitalSignsHistoryDetail.Code.String, s.vitalCodePulse) ||
            strings.EqualFold(vitalSignsHistoryDetail.Code.String, s.vitalCodeClinicaltemplatePulse) {
            patientVitalSignsDetailPulse := model.NovaPatientVitalSignsDetail{
                Code:         utils.NewNullString(vitalSignsHistoryDetail.Code.String),
                Description:  utils.NewNullString(vitalSignsHistoryDetail.Description.String),
                RefNo:        utils.NewNullString(vitalSignsHistoryDetail.RefNo.String),
                RecordedDate: utils.NewNullString(vitalSignsHistoryDetail.RecordedDate.String),
                Unit:         utils.NewNullString(vitalSignsHistoryDetail.Unit.String),
                Value1:       utils.NewNullString(vitalSignsHistoryDetail.Value1.String),
                Value2:       utils.NewNullString(vitalSignsHistoryDetail.Value2.String),
            }

            novaPatientVitalSignsDetailPulseDto := model.NovaPatientVitalSignsDetailDto{
                NovaPatientVitalSignsDetail: &patientVitalSignsDetailPulse,
                Value1LowValue:              "60",
                Value1HighValue:             "100",
            }
            patientVitalSignsDetailPulseList = append(patientVitalSignsDetailPulseList, novaPatientVitalSignsDetailPulseDto)
        }
    }
    vitalSignsGroupHeight.VitalSignsData = patientVitalSignsDetailHeightList
    vitalSignsGroupWeight.VitalSignsData = patientVitalSignsDetailWeightList
    vitalSignsGroupBMI.VitalSignsData = patientVitalSignsDetailBMIList
    vitalSignsGroupBP.VitalSignsData = patientVitalSignsDetailBPList
    vitalSignsGroupPulse.VitalSignsData = patientVitalSignsDetailPulseList

    vitalSignsDashboardList = append(vitalSignsDashboardList, vitalSignsGroupBP)
    vitalSignsDashboardList = append(vitalSignsDashboardList, vitalSignsGroupBMI)
    vitalSignsDashboardList = append(vitalSignsDashboardList, vitalSignsGroupPulse)
    vitalSignsDashboardList = append(vitalSignsDashboardList, vitalSignsGroupHeight)
    vitalSignsDashboardList = append(vitalSignsDashboardList, vitalSignsGroupWeight)

    return vitalSignsDashboardList, nil
}

func (s *HealthCareService) GetVitalSignsHistory(prn string, visitDate string, vitalSignsCode string) ([]model.NovaPatientVitalSignsDetailDto, error) {
    patientVitalSignsHistoryList := make([]model.NovaPatientVitalSignsDetail, 0)
    var err error
    if strings.EqualFold(vitalSignsCode, "WEIGHT") {
        patientVitalSignsHistoryList, err = s.novaVitalSignsDetailService.GetPatientVitalSignsHistoryWeight(prn, s.vitalCodeWeight, nil)
        if err != nil {
            return nil, err
        }
    }
    if strings.EqualFold(vitalSignsCode, "HEIGHT") {
        patientVitalSignsHistoryList, err = s.novaVitalSignsDetailService.GetPatientVitalSignsHistoryHeight(prn, s.vitalCodeHeight, nil)
        if err != nil {
            return nil, err
        }
    }
    if strings.EqualFold(vitalSignsCode, "BP") {
        patientVitalSignsHistoryList, err = s.novaVitalSignsDetailService.GetPatientVitalSignsHistoryBP(prn, s.vitalCodeBP, nil)
        if err != nil {
            return nil, err
        }
    }
    if strings.EqualFold(vitalSignsCode, "BMI") {
        patientVitalSignsHistoryList, err = s.novaVitalSignsDetailService.GetPatientVitalSignsHistoryBMI(prn, s.vitalCodeBMI, nil)
        if err != nil {
            return nil, err
        }
    }
    if strings.EqualFold(vitalSignsCode, "PR") {
        patientVitalSignsHistoryList, err = s.novaVitalSignsDetailService.GetPatientVitalSignsHistoryPulse(prn, s.vitalCodePulse, nil)
        if err != nil {
            return nil, err
        }
    }

    if len(patientVitalSignsHistoryList) < 1 {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    novaPatientVitalSignsDetailDtoList := make([]model.NovaPatientVitalSignsDetailDto, 0)
    for _, vitalSignsHistoryDetail := range patientVitalSignsHistoryList {
        novaPatientVitalSignsDetailDto := model.NovaPatientVitalSignsDetailDto{}
        if strings.EqualFold(vitalSignsCode, "BMI") {
            novaPatientVitalSignsDetailDto.Value1LowValue = "17"
            novaPatientVitalSignsDetailDto.Value1HighValue = "24"
        } else if strings.EqualFold(vitalSignsCode, "BP") {
            novaPatientVitalSignsDetailDto.Value1LowValue = "120"
            novaPatientVitalSignsDetailDto.Value1HighValue = "130"
            novaPatientVitalSignsDetailDto.Value2LowValue = "80"
            novaPatientVitalSignsDetailDto.Value2HighValue = "90"
        } else if strings.EqualFold(vitalSignsCode, "PR") {
            novaPatientVitalSignsDetailDto.Value1LowValue = "60"
            novaPatientVitalSignsDetailDto.Value1HighValue = "100"
        }

        novaPatientVitalSignsDetailDto.NovaPatientVitalSignsDetail = &vitalSignsHistoryDetail
        novaPatientVitalSignsDetailDtoList = append(novaPatientVitalSignsDetailDtoList, novaPatientVitalSignsDetailDto)
    }

    return novaPatientVitalSignsDetailDtoList, nil
}

func (s *HealthCareService) GetLabHistoryForDahsboard(prn string) ([]model.NovaLabHistoryDashboard, error) {
    labHistoryDashboardList := make([]model.NovaLabHistoryDashboard, 0)
    patientLabHistoryList, err := s.novaInvestigationDetailService.GetLabHistoryTrendingForDashboard(
        prn, s.labInvestigationType, s.labCodeHDL,
        s.labCodeLDL, s.labCodeGlucose, s.labCodeHemoglobin, nil,
    )
    if err != nil {
        return nil, err
    }
    if len(patientLabHistoryList) < 1 {
        labHistoryHDL := model.NovaLabHistoryDashboard{}
        labHistoryLDL := model.NovaLabHistoryDashboard{}
        labHistoryGlucose := model.NovaLabHistoryDashboard{}
        labHistoryHemoglobin := model.NovaLabHistoryDashboard{}

        labHistoryHDL.LabCode = "HDL"
        labHistoryLDL.LabCode = "LDL"
        labHistoryGlucose.LabCode = "Glucose"
        labHistoryHemoglobin.LabCode = "Hemoglobin"

        patientLabHistoryHDLList := make([]model.NovaPatientInvestigationDetail, 0)
        patientLabHistoryLDLList := make([]model.NovaPatientInvestigationDetail, 0)
        patientLabHistoryGlucoseList := make([]model.NovaPatientInvestigationDetail, 0)
        patientLabHistoryHemoglobinList := make([]model.NovaPatientInvestigationDetail, 0)

        for _, labHistoryDetail := range patientLabHistoryList {
            if strings.EqualFold(labHistoryDetail.Code.String, s.labCodeHDL) {
                labDetailHDL := model.NovaPatientInvestigationDetail{
                    InvestigationRefNo: utils.NewNullString(labHistoryDetail.InvestigationRefNo.String),
                    Code:               utils.NewNullString(labHistoryDetail.Code.String),
                    ResultValue:        utils.NewNullString(labHistoryDetail.ResultValue.String),
                    ResultUnit:         utils.NewNullString(labHistoryDetail.ResultUnit.String),
                    RecordedDate:       utils.NewNullString(labHistoryDetail.RecordedDate.String),
                }
                patientLabHistoryHDLList = append(patientLabHistoryHDLList, labDetailHDL)
            }

            if strings.EqualFold(labHistoryDetail.Code.String, s.labCodeLDL) {
                labDetailLDL := model.NovaPatientInvestigationDetail{
                    InvestigationRefNo: utils.NewNullString(labHistoryDetail.InvestigationRefNo.String),
                    Code:               utils.NewNullString(labHistoryDetail.Code.String),
                    ResultValue:        utils.NewNullString(labHistoryDetail.ResultValue.String),
                    ResultUnit:         utils.NewNullString(labHistoryDetail.ResultUnit.String),
                    RecordedDate:       utils.NewNullString(labHistoryDetail.RecordedDate.String),
                }
                patientLabHistoryLDLList = append(patientLabHistoryLDLList, labDetailLDL)
            }

            if strings.EqualFold(labHistoryDetail.Code.String, s.labCodeGlucose) {
                labDetailGLU := model.NovaPatientInvestigationDetail{
                    InvestigationRefNo: utils.NewNullString(labHistoryDetail.InvestigationRefNo.String),
                    Code:               utils.NewNullString(labHistoryDetail.Code.String),
                    ResultValue:        utils.NewNullString(labHistoryDetail.ResultValue.String),
                    ResultUnit:         utils.NewNullString(labHistoryDetail.ResultUnit.String),
                    RecordedDate:       utils.NewNullString(labHistoryDetail.RecordedDate.String),
                }
                patientLabHistoryGlucoseList = append(patientLabHistoryGlucoseList, labDetailGLU)
            }

            if strings.EqualFold(labHistoryDetail.Code.String, s.labCodeHemoglobin) {
                labDetailHBA1 := model.NovaPatientInvestigationDetail{
                    InvestigationRefNo: utils.NewNullString(labHistoryDetail.InvestigationRefNo.String),
                    Code:               utils.NewNullString(labHistoryDetail.Code.String),
                    ResultValue:        utils.NewNullString(labHistoryDetail.ResultValue.String),
                    ResultUnit:         utils.NewNullString(labHistoryDetail.ResultUnit.String),
                    RecordedDate:       utils.NewNullString(labHistoryDetail.RecordedDate.String),
                }
                patientLabHistoryHemoglobinList = append(patientLabHistoryHemoglobinList, labDetailHBA1)
            }
        }

        labHistoryHDL.LabData = patientLabHistoryHDLList
        labHistoryLDL.LabData = patientLabHistoryLDLList
        labHistoryGlucose.LabData = patientLabHistoryGlucoseList
        labHistoryHemoglobin.LabData = patientLabHistoryHemoglobinList

        labHistoryDashboardList = append(labHistoryDashboardList, labHistoryHDL)
        labHistoryDashboardList = append(labHistoryDashboardList, labHistoryLDL)
        labHistoryDashboardList = append(labHistoryDashboardList, labHistoryGlucose)
        labHistoryDashboardList = append(labHistoryDashboardList, labHistoryHemoglobin)
    }

    return labHistoryDashboardList, nil
}

func (s *HealthCareService) GetPdfHealthScreeningReport(hsrRefNo string) ([]byte, error) {
    novaHealthScreeningRptDetailDtoList := make([]model.NovaHealthScreeningRptDetail, 0)
    novaHealthScreeningRptDetailList, err := s.novaHealthScreenRptDetailService.FindEachHealthScreeningRptByHSRRefNo(hsrRefNo, nil)
    if len(novaHealthScreeningRptDetailList) < 1 {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }

    for _, novaHealthScreeningDetailRpt := range novaHealthScreeningRptDetailList {
        novaHealthScreeningRptDetail := model.NovaHealthScreeningRptDetail{
            HsrRefNo:     novaHealthScreeningDetailRpt.HsrRefNo.String,
            HsrClobValue: novaHealthScreeningDetailRpt.HsrClobValue.String,
        }
        novaHealthScreeningRptDetailDtoList = append(novaHealthScreeningRptDetailDtoList, novaHealthScreeningRptDetail)
    }

    o := novaHealthScreeningRptDetailDtoList[0].HsrClobValue
    decoded, err := base64.StdEncoding.DecodeString(o)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return decoded, nil
}

func (s *HealthCareService) GetPatientAllergy(prn string) ([]model.NovaPatientAlert, error) {
    lx, err := s.novaPatientAlertService.FindPatientActiveAlertByPrn(prn, nil)
    if err != nil {
        return nil, err
    }
    if len(lx) < 1 {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }
    return lx, nil
}

func (s *HealthCareService) GetPatientFromReportSchemaByPrn(prn string) (*model.NovaPatient, error) {
    o, err := s.novaPatientService.FindByPrn(prn, nil)
    if o == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }
    if err != nil {
        return nil, err
    }
    return o, nil
}

func (s *HealthCareService) VesaliusGetPastAppointments(prn string) ([]gm.PatientPastAppointment, error) {
    lx := make([]gm.PatientPastAppointment, 0)
    lid := make([]string, 0)
    familyMembers := make([]gm.ApplicationUserFamily, 0)
    patient, err := s.applicationUserService.FindByPRN(prn, s.db)
    if err != nil {
        return nil, err
    }

    if patient == nil {
        return nil, fiber.NewError(fiber.StatusNotFound)
    }

    familyMembers, err = s.applicationUserFamilyService.FindAllByUserPrnAppt(patient.MasterPrn.String, true, true, s.db)
    if err != nil {
        return nil, err
    }

    for i := range familyMembers {
        f := familyMembers[i]
        q := `
            SELECT * FROM NOVA_PATIENT_APPOINTMENT 
             WHERE APPOINTMENT_DATE < SYSDATE 
             AND STATUS = 'REGISTER' 
             AND PRN = :prn 
            ORDER BY APPOINTMENT_DATE DESC
        `
        list := make([]model.PastAppointment, 0)
        err := s.db.SelectContext(s.ctx, &list, q, f.NokPrn.String)
        if err != nil {
            utils.LogError(err)
            return nil, err
        }

        var (
            sessionType      string = ""
            sessionStartTime string = ""
            sessionEndTime   string = ""
        )

        for j := range list {
            o := list[j]
            doc, err := s.novaDoctorService.FindDoctorByMcr(o.Mcr.String)
            if err != nil {
                return nil, err
            }

            if doc != nil {
                patientPastAppointment := gm.PatientPastAppointment{
                    DoctorId: doc.DoctorId.Int64,
                    Image:    doc.Image.String,
                    MCR:      doc.MCR.String,
                    Name:     doc.Name.String,
                }
                q := `
                    SELECT hp.PACKAGE_IMG, hp.PACKAGE_NAME, 
                    ppd.PACKAGE_PURCHASE_NO, ppd.EXPIRED_DATETIME
                    FROM NOVA_DOCTOR_PATIENT_APPT ndpa
                    JOIN PATIENT_PURCHASE_DETAILS ppd ON ndpa.PACKAGE_PURCHASE_NO = ppd.PACKAGE_PURCHASE_NO
                    LEFT JOIN HOSPITAL_PACKAGE hp ON ppd.PACKAGE_ID = hp.PACKAGE_ID 
                    WHERE ndpa.PATIENT_PRN = :prn AND ndpa.PACKAGE_PURCHASE_NO = :packagePurchaseNo
                `
                lp := make([]upck.UserPackage, 0)
                err := s.db.SelectContext(s.ctx, &lp, q, f.NokPrn.String, o.Reason.String)
                if err != nil {
                    utils.LogError(err)
                    return nil, err
                }

                var mobileAppt *upck.UserPackage
                if len(lp) > 0 {
                    mobileAppt = &lp[0]
                    patientPastAppointment.Image = ""
                    patientPastAppointment.Name = ""
                }

                g, _ := goment.New(o.AppointmentDate.String, "YYYY-MM-DD")
                o.AppointmentDate = utils.NewNullString(g.Format("DD-MMM-YYYY"))

                la := make([]gm.NovaDoctorAppointmentLists, 0)
                err = s.db.SelectContext(s.ctx, &la, sqx.GET_SINGLEDATE_DOCTOR_APPOINTMENTS,
                    sql.Named("doctorId", doc.DoctorId.Int64),
                    sql.Named("dt", o.AppointmentDate.String),
                )
                if err != nil {
                    utils.LogError(err)
                    return nil, err
                }

                var appt *gm.NovaDoctorAppointmentLists
                if len(la) > 0 {
                    appt = &la[0]
                }

                if appt.NormalStatus.String == "NOT AVAILABLE" {
                    q := `SELECT * FROM NOVA_DOCTOR_APPT_SLOT WHERE DOCTOR_ID = :doctorId`
                    ls := make([]gm.NovaDoctorApptSlot, 0)
                    err = s.db.SelectContext(s.ctx, &ls, q, doc.DoctorId.Int64)
                    if err != nil {
                        utils.LogError(err)
                        return nil, err
                    }

                    if len(ls) < 1 {
                        sessionType = ""
                        sessionStartTime = ""
                        sessionEndTime = ""
                    }

                    for k := range ls {
                        docAppt := ls[k]
                        if docAppt.SessionType.String == "MORNING" && appt.MorningStatus.String == "AVAILABLE" {
                            vesStartTime, _ := goment.New(o.AppointmentTime.String, "hh:mm")
                            docApptStartTime, _ := goment.New(docAppt.StartTime.String, "hh:mm")
                            docApptEndTime, _ := goment.New(docAppt.EndTime.String, "hh:mm")
                            isWithinMorningRange := !vesStartTime.IsBefore(docApptStartTime) && !vesStartTime.IsAfter(docApptEndTime)

                            if isWithinMorningRange {
                                sessionType = "Morning"
                                sessionStartTime = docAppt.StartTime.String
                                sessionEndTime = docAppt.EndTime.String
                            }
                        }

                        if docAppt.SessionType.String == "AFTERNOON" && appt.AfternoonStatus.String == "AVAILABLE" {
                            vesStartTime, _ := goment.New(o.AppointmentTime, "hh:mm")
                            docApptStartTime, _ := goment.New(docAppt.StartTime.String, "hh:mm")
                            docApptEndTime, _ := goment.New(docAppt.EndTime.String, "hh:mm")
                            isWithinAfternoonRange := !vesStartTime.IsBefore(docApptStartTime) && !vesStartTime.IsAfter(docApptEndTime)

                            if isWithinAfternoonRange {
                                sessionType = "Afternoon"
                                sessionStartTime = docAppt.StartTime.String
                                sessionEndTime = docAppt.EndTime.String
                            }
                        }
                    }
                }

                rel := "Self"
                if f.Relationship.String != "Self" {
                    rel = f.Fullname.String
                }

                packageName := ""
                packagePurchaseNo := ""
                packageImage := ""
                if mobileAppt != nil {
                    packageName = mobileAppt.PackageName.String
                    packagePurchaseNo = mobileAppt.PackagePurchaseNo.String
                    packageImage = mobileAppt.PackageImage.String
                }

                apptInfo := gm.VesaliusPastApptInfo{
                    AppointmentRefNo:  int(o.AppointmentRefNo.Int64),
                    HospitalCode:      o.HospitalCode.String,
                    PatientPRN:        o.PatientPrn.String,
                    PatientName:       o.PatientName.String,
                    PatientDOB:        o.PatientDob.String,
                    HomeContactNo:     o.HomeContactNo.String,
                    MobileContactNo:   o.MobileContactNo.String,
                    HomeAddress:       o.HomeAddress.String,
                    AppointmentDate:   o.AppointmentDate.String,
                    AppointmentTime:   o.AppointmentTime.String,
                    DurationMins:      int(o.DurationMins.Int64),
                    CaseType:          o.CaseType.String,
                    Status:            o.Status.String,
                    Reason:            o.Reason.String,
                    MCR:               o.Mcr.String,
                    DoctorName:        o.DoctorName.String,
                    RoomNo:            o.RoomNo.String,
                    ClinicName:        o.ClinicName.String,
                    AccountNo:         o.AccountNo.String,
                    AppointmentSource: o.AppointmentSource.String,
                    SourceOfReferral:  o.SourceOfReferral.String,
                    AppointmentType:   o.AppointmentType.String,
                    SyncDate:          o.SyncDate.String,
                    TransferFlg:       o.TransferFlg.String,
                    TransferDateTime:  o.TransferDateTime.String,
                    TransferSystem:    o.TransferSystem.String,

                    ApptPatientName:       rel,
                    ApptSessionType:       sessionType,
                    SessionStartTime:      sessionStartTime,
                    SessionEndTime:        sessionEndTime,
                    ApptPackageName:       packageName,
                    ApptPackagePurchaseNo: packagePurchaseNo,
                    ApptPackageImage:      packageImage,
                    ApptSlotType:          "Normal",
                }

                if sessionType == "Morning" || sessionType == "Afternoon" {
                    apptInfo.ApptSlotType = "Session"
                }

                patientPastAppointment.VesaliusPastApptInfo = apptInfo

                lid = append(lid, strconv.FormatInt(doc.DoctorId.Int64, 10))
                lx = append(lx, patientPastAppointment)
            }
        }

    }

    if len(lid) > 0 {
        doctorIds := strings.Join(lid, ",")
        m3, _ := s.novaDoctorService.FindAllNovaDoctorSpecialities(doctorIds, s.db)
        m4, _ := s.novaDoctorService.FindAllNovaDoctorClinicLocation(doctorIds, s.db)
        m5, _ := s.novaDoctorService.FindAllNovaDoctorContact(doctorIds, s.db)

        for i := range lx {
            o := lx[i]
            if list, ok := m3[o.DoctorId]; ok {
                lx[i].DoctorSpecialities = list
            }
            if list, ok := m4[o.DoctorId]; ok {
                lx[i].DoctorClinicLocation = list
            }
            if list, ok := m5[o.DoctorId]; ok {
                lx[i].DoctorContact = list
            }
        }
    }

    return lx, nil
}
