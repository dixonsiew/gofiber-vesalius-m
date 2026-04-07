package healthCare

import (
    "context"
    "encoding/base64"
    "vesaliusm/config"
    "vesaliusm/database"
    "vesaliusm/model/healthCare"
    "vesaliusm/service/adminUser"
    "vesaliusm/service/applicationUser"
    "vesaliusm/service/applicationUserFamily"
    "vesaliusm/service/branch"
    "vesaliusm/service/novaBill"
    "vesaliusm/service/novaDoctor"
    "vesaliusm/service/novaDoctorPatientAppt"
    "vesaliusm/service/novaHealthScreenRpt"
    "vesaliusm/service/novaHealthScreenrptDetail"
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
    "vesaliusm/utils"

    //"github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3"
    "github.com/jmoiron/sqlx"
)

var HealthCareSvc *HealthCareService = NewHealthCareService(database.GetDbrs(), database.GetCtx())

type HealthCareService struct {
    db  *sqlx.DB
    ctx context.Context

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

func NewHealthCareService(db *sqlx.DB, ctx context.Context) *HealthCareService {
    return &HealthCareService{
        db:                              db,
        ctx:                             ctx,

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
        applicationUserFamilyService:     applicationUserFamily.NewApplicationUserFamilyService(db, ctx),
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

func (s *HealthCareService) GetPatientVisit(prn string, pageId string) ([]healthCare.NovaVisitDetails, error) {
    novaVisitDetails := make([]healthCare.NovaVisitDetails, 0)
    visitHistoryType := pageId
    switch visitHistoryType {
    case "1":
        novaVisits, err := s.novaVisitService.GetSpecificPatientVisit(prn, s.db)
        if err != nil {
            return nil, err
        }
        if len(novaVisits) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for _, novaVisit := range novaVisits {
            novaVisitSummaries, err := s.novaVisitSummaryService.FindByAccountNoAndCategoryNotAndCategoryNot(novaVisit.AccountNo.String, "Investigation", "Medication", s.db)
            if err != nil {
                return nil, err
            }
            if len(novaVisitSummaries) > 0 {
                for _, _ = range novaVisitSummaries {
                    novaVisitDetail := healthCare.NovaVisitDetails{
                        NovaVisit: &novaVisit,
                        NovaVisitSummaries: novaVisitSummaries,
                    }
                    novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
                }
            }
        }
    
    case "2":
        novaVisitsPrescription, err := s.novaVisitService.GetSpecificPatientVisit(prn, s.db)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsPrescription) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for _, novaVisit := range novaVisitsPrescription {
            novaVisitPatientRxList, err := s.novaPatientrxService.FindPatientRxByAccountNo(novaVisit.AccountNo.String, s.db)
            if err != nil {
                return nil, err
            }
            if len(novaVisitPatientRxList) > 0 {
                for _, _ = range novaVisitPatientRxList {
                    novaVisitDetail := healthCare.NovaVisitDetails{
                        NovaVisit: &novaVisit,
                        NovaVisitPatientRxList: novaVisitPatientRxList,
                    }
                    novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
                }
            }
        }

    case "3":
        novaVisitsInvestigation, err := s.novaVisitService.GetSpecificPatientVisit(prn, s.db)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsInvestigation) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for _, novaVisit := range novaVisitsInvestigation {
            novaPatientInvestigationList, err := s.novaInvestigationService.FindNonPanelInvestiationByAccountNo(novaVisit.AccountNo.String, s.db)
            if err != nil {
                return nil, err
            }
            uniquePanelList, err := s.novaInvestigationService.FindUniquePanelInvestigationByAccountNo(novaVisit.AccountNo.String, s.db)
            if err != nil {
                return nil, err
            }
            if len(novaPatientInvestigationList) > 0 || len(uniquePanelList) > 0 {
                for _, _ = range novaPatientInvestigationList {
                    novaVisitDetail := healthCare.NovaVisitDetails{
                        NovaVisit: &novaVisit,
                    }
                    novaVisitInvestigationDetailList := make([]healthCare.NovaVisitInvestigationDetail, 0)
                    if len(novaPatientInvestigationList) > 0 {
                        for _, novaPatientInvestigation := range novaPatientInvestigationList {
                            novaVisitInvestigationDetailDto := healthCare.NovaVisitInvestigationDetail{}
                            novaPatientInvestigationDetailList, err := s.novaInvestigationDetailService.GetInvestigationRefNoAndCode(novaPatientInvestigation.InvestigationRefNo.String, novaPatientInvestigation.Code.String, s.db)
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
                        for _, uniquePanelCode := range uniquePanelList {
                            novaPatientPanelInvestigationList, err := s.novaInvestigationService.FindPanelInvestigationByAccountNoPanelCode(novaVisit.AccountNo.String, uniquePanelCode, s.db)
                            if err != nil {
                                return nil, err
                            }
                            if len(novaPatientPanelInvestigationList) > 0 {
                                novaVisitInvestigationDetailDto := healthCare.NovaVisitInvestigationDetail{}
                                panelDetailList := make([]healthCare.NovaVisitInvestigationPanelDetail, 0)
                                for _, novaPatientInvestigation := range novaPatientPanelInvestigationList {
                                    novaPatientInvestigationDetailList, err := s.novaInvestigationDetailService.GetInvestigationRefNoAndCode(novaPatientInvestigation.InvestigationRefNo.String, novaPatientInvestigation.Code.String, s.db)
                                    if err != nil {
                                        return nil, err
                                    }
                                    if len(novaPatientInvestigationDetailList) > 0 {
                                        for _, novaPatientInvestigationDetail := range novaPatientInvestigationDetailList {
                                            panelDetail := healthCare.NovaVisitInvestigationPanelDetail{
                                                Code: novaPatientInvestigationDetail.Code.String,
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
                                novaVisitInvestigationDetailDto.PanelCode = novaPatientPanelInvestigationList[0].PanelCode.String;
                                novaVisitInvestigationDetailDto.PanelDescription = novaPatientPanelInvestigationList[0].PanelDescription.String;
                                novaVisitInvestigationDetailDto.PanelDetail = panelDetailList;
                            }
                        }
                    }
                    novaVisitDetail.NovaVisitInvestigationDetailList = novaVisitInvestigationDetailList
                    novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
                }
            }
        }
    case "4":
        novaVisitsBill, err := s.novaVisitService.GetSpecificPatientVisit(prn, s.db)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsBill) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for _, novaVisit := range novaVisitsBill {
            novaBills, err := s.novaBillService.GetNovaBillByPrnAndAccountNo(prn, novaVisit.AccountNo.String, s.db)
            if err != nil {
                return nil, err
            }
            if len(novaBills) > 0 {
                novaVisitDetail := healthCare.NovaVisitDetails{
                    NovaVisit: &novaVisit,
                    NovaBills: novaBills,
                }
                novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
            }
        }
    case "5":
        novaVisitsVital, err := s.novaVisitService.GetPatientVisitWithVitalSign(prn, s.db)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsVital) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for _, novaVisit := range novaVisitsVital {
            novaPatientVitalSignsList, err := s.novaVitalSignsService.FindPatientVitalSignsByAccountNo(novaVisit.AccountNo.String, s.db)
            if err != nil {
                return nil, err
            }
            novaVisitVitalSignsDetailList := make([]healthCare.NovaVisitVitalSignsDetail, 0)
            if len(novaPatientVitalSignsList) > 0 {
                for _, patientVitalSigns := range novaPatientVitalSignsList {
                    patientVitalSignsDetailList, err := s.novaVitalSignsDetailService.GetLocalPatientVitalSignsDetailsByRefNo(
                        patientVitalSigns.RefNo.String,
                        s.vitalCodeHeight,
                        s.vitalCodeWeight,
                        s.vitalCodeBP,
                        s.vitalCodeBMI,
                        s.vitalCodePulse,
                        s.db,
                    )
                    if err != nil {
                        return nil, err
                    }
                    if len(patientVitalSignsDetailList) > 0 {
                        for _, eachPatientVitalSignsDetail := range patientVitalSignsDetailList {
                            visitVitalSignsDetail := healthCare.NovaVisitVitalSignsDetail{
                                Code: eachPatientVitalSignsDetail.Code.String,
                                Desc: eachPatientVitalSignsDetail.Description.String,
                                Value1: eachPatientVitalSignsDetail.Value1.String,
                                Value2: eachPatientVitalSignsDetail.Value2.String,
                                Unit: eachPatientVitalSignsDetail.Unit.String,
                            }
                            novaVisitVitalSignsDetailList = append(novaVisitVitalSignsDetailList, visitVitalSignsDetail)
                        }
                    }
                }
                novaVisitDetail := healthCare.NovaVisitDetails{
                    NovaVisit: &novaVisit,
                    NovaVisitVitalSignsDetailList: novaVisitVitalSignsDetailList,
                }
                novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
            }
        }
    case "6":
        novaVisitsWithReferralLetterList, err := s.novaVisitService.GetPatientVisitWithReferralLetter(prn, s.db)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsWithReferralLetterList) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for _, novaVisit := range novaVisitsWithReferralLetterList {
            novaReferralLetterList, err := s.novaReferralLetterService.FindAllReferralLetterByPrnAndAccountNo(novaVisit.PRN.String, novaVisit.AccountNo.String, s.db)
            if err != nil {
                return nil, err
            }
            novaVisitReferralLetterList := make([]healthCare.NovaVisitReferralLetter, 0)
            if len(novaReferralLetterList) > 0 {
                for _, novaExtReferralLetter := range novaReferralLetterList {
                    novaVisitReferralLetterDetail := healthCare.NovaVisitReferralLetter{
                        ReferralRefNo: novaExtReferralLetter.ReferralRefNo.String,
                        PRN: novaExtReferralLetter.PRN.String,
                        AccountNo: novaExtReferralLetter.AccountNo.String,
                        ReferralDateTime: novaExtReferralLetter.ReferralDateTime.String,
                        ReferrerDoctor: novaExtReferralLetter.ReferrerDoctorMCR.String,
                        ReferralDoctor: novaExtReferralLetter.ReferralTo.String,
                        ReferralTitleDept: novaExtReferralLetter.ReferralTitleDept.String,
                        ReferralAddressOrSubject: novaExtReferralLetter.ReferralAddressOrSubject.String,
                        ReferralLetter: novaExtReferralLetter.ReferralLetter.String,
                        ReferralType: novaExtReferralLetter.ReferralType.String,
                    }
                    novaVisitReferralLetterList = append(novaVisitReferralLetterList, novaVisitReferralLetterDetail)
                }
            }
            novaVisitDetail := healthCare.NovaVisitDetails{
                NovaVisit: &novaVisit,
                NovaVisitReferralLetterList: novaVisitReferralLetterList,
            }
            novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
        }
    case "7":
        novaVisitsWithHealthScreeningReportList, err := s.novaVisitService.GetPatientVisitWithHealthScreeningReport(prn, s.db)
        if err != nil {
            return nil, err
        }
        if len(novaVisitsWithHealthScreeningReportList) < 1 {
            return nil, fiber.NewError(fiber.StatusNoContent)
        }
        for _, novaVisit := range novaVisitsWithHealthScreeningReportList {
            novaHealthScreeningRptList, err := s.novaHealthScreenRptService.FindHealthScreeningRptByPrnAndAccountNo(novaVisit.PRN.String, novaVisit.AccountNo.String, s.db)
            if err != nil {
                return nil, err
            }
            novaHealthScreeningRptDtoList := make([]healthCare.NovaHealthScreeningRpt, 0)
            if len(novaHealthScreeningRptList) > 0 {
                for _, novaHealthScreeningRpt := range novaHealthScreeningRptList {
                    novaHealthScreeningRptDetail := healthCare.NovaHealthScreeningRpt{
                        HsrRefNo: utils.NewNullString(novaHealthScreeningRpt.HsrRefNo.String),
                        AccountNo: utils.NewNullString(novaHealthScreeningRpt.AccountNo.String),
                        ReportDate: utils.NewNullString(novaHealthScreeningRpt.ReportDate.String),
                        ReportUser: utils.NewNullString(novaHealthScreeningRpt.ReportUser.String),
                    }
                    novaHealthScreeningRptDtoList = append(novaHealthScreeningRptDtoList, novaHealthScreeningRptDetail)
                }
            }
            novaVisitDetail := healthCare.NovaVisitDetails{
                NovaVisit: &novaVisit,
                NovaHealthScreeningRptList: novaHealthScreeningRptDtoList,
            }
            novaVisitDetails = append(novaVisitDetails, novaVisitDetail)
        }
    }
    return novaVisitDetails, nil
}

func (s *HealthCareService) GetPdfHealthScreeningReport(hsrRefNo string) ([]byte, error) {
    novaHealthScreeningRptDetailDtoList := make([]healthCare.NovaHealthScreeningRptDetail, 0)
    novaHealthScreeningRptDetailList, err := s.novaHealthScreenRptDetailService.FindEachHealthScreeningRptByHSRRefNo(hsrRefNo, nil)
    if len(novaHealthScreeningRptDetailList) < 1 {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }
    
    for _, novaHealthScreeningDetailRpt := range novaHealthScreeningRptDetailList {
        novaHealthScreeningRptDetail := healthCare.NovaHealthScreeningRptDetail{
            HsrRefNo: novaHealthScreeningDetailRpt.HsrRefNo.String,
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

func (s *HealthCareService) GetPatientAllergy(prn string) ([]healthCare.NovaPatientAlert, error) {
    lx, err := s.novaPatientAlertService.FindPatientActiveAlertByPrn(prn, nil)
    if err != nil {
        return nil, err
    }
    if len(lx) < 0 {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }
    return lx, nil
}

func (s *HealthCareService) GetPatientFromReportSchemaByPrn(prn string) (*healthCare.NovaPatient, error) {
    o, err := s.novaPatientService.FindByPrn(prn, nil)
    if o == nil {
        return nil, fiber.NewError(fiber.StatusNoContent)
    }
    if err != nil {
        return nil, err
    }
    return o, nil
}

// func (s *HealthCareService)
