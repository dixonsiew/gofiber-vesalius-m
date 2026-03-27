package dto

import (
    "vesaliusm/model"
    "vesaliusm/utils"
)

type NovaDoctorAppointmentDto struct {
    DisplaySequence int    `json:"displaySequence" validate:"required,numeric"`
    ApptDayOfWeek   string `json:"apptDayOfWeek" validate:"required,max=20"`
    ApptSlotType    string `json:"apptSlotType" validate:"required"`
    ApptSessionType string `json:"apptSessionType" validate:"required,max=10"`
    ApptStartTime   string `json:"apptStartTime" validate:"required"`
    ApptEndTime     string `json:"apptEndTime" validate:"required"`
    ApptMaxSlots    int    `json:"apptMaxSlots" validate:"required,numeric"`
}

func (o NovaDoctorAppointmentDto) ToDbModel() model.NovaDoctorAppointment {
    return model.NovaDoctorAppointment{
        DisplaySequence: utils.NewInt32(int32(o.DisplaySequence)),
        ApptDayOfWeek:   utils.NewNullString(o.ApptDayOfWeek),
        ApptSlotType:    utils.NewNullString(o.ApptSlotType),
        ApptSessionType: utils.NewNullString(o.ApptSessionType),
        ApptStartTime:   utils.NewNullString(o.ApptStartTime),
        ApptEndTime:     utils.NewNullString(o.ApptEndTime),
        ApptMaxSlots:    utils.NewInt32(int32(o.ApptMaxSlots)),
    }
}

type NovaDoctorClinicHoursDto struct {
    DisplaySequence   int    `json:"displaySequence" validate:"required,numeric"`
    DayOfTheWeek      string `json:"dayOfTheWeek" validate:"required,max=20"`
    DayStartTime      string `json:"dayStartTime" validate:"required,max=20"`
    DayEndTime        string `json:"dayEndTime" validate:"required,max=20"`
    ByAppointmentOnly bool   `json:"byAppointmentOnly" validate:"required"`
}

func (o NovaDoctorClinicHoursDto) ToDbModel() model.NovaDoctorClinicHours {
    byAppointmentOnly := 0
    if o.ByAppointmentOnly {
        byAppointmentOnly = 1
    }
    return  model.NovaDoctorClinicHours{
        DisplaySequence:    utils.NewInt32(int32(o.DisplaySequence)),
        DayOfTheWeek:       utils.NewNullString(o.DayOfTheWeek),
        DayStartTime:       utils.NewNullString(o.DayStartTime),
        DayEndTime:         utils.NewNullString(o.DayEndTime),
        ByAppointmentOnlyV: utils.NewInt32(int32(byAppointmentOnly)),
        ByAppointmentOnly:  o.ByAppointmentOnly,
    }
}

type NovaDoctorClinicLocationDto struct {
    Location string `json:"location" validate:"required,max=200"`
    Building string `json:"building" validate:"required,max=200"`
}

func (o NovaDoctorClinicLocationDto) ToDbModel() model.NovaDoctorClinicLocation {
    return model.NovaDoctorClinicLocation{
        Location: utils.NewNullString(o.Location),
        Building: utils.NewNullString(o.Building),
    }
}

type NovaDoctorContactDto struct {
    DisplaySequence int    `json:"displaySequence" validate:"required,numeric"`
    ContactType     string `json:"contactType" validate:"required,max=255"`
    ContactValue    string `json:"contactValue" validate:"required,max=255"`
}

func (o NovaDoctorContactDto) ToDbModel() model.NovaDoctorContact {
    return model.NovaDoctorContact{
        DisplaySequence: utils.NewInt32(int32(o.DisplaySequence)),
        ContactType:     utils.NewNullString(o.ContactType),
        ContactValue:    utils.NewNullString(o.ContactValue),
    }
}

type NovaDoctorQualificationsDto struct {
    DisplaySequence int    `json:"displaySequence" validate:"required,numeric"`
    Qualification   string `json:"qualification" validate:"required,max=200"`
}

func (o NovaDoctorQualificationsDto) ToDbModel() model.NovaDoctorQualifications {
    return model.NovaDoctorQualifications{
        DisplaySequence: utils.NewInt32(int32(o.DisplaySequence)),
        Qualification:   utils.NewNullString(o.Qualification),
    }
}

type NovaDoctorSpecialitiesDto struct {
    DisplaySequence int    `json:"displaySequence" validate:"required,numeric"`
    Specialities    string `json:"specialities" validate:"required,max=200"`
    Subspecialty    string `json:"subspecialty" validate:"required,max=200"`
}

func (o NovaDoctorSpecialitiesDto) ToDbModel() model.NovaDoctorSpecialities {
    return model.NovaDoctorSpecialities{
        DisplaySequence: utils.NewInt32(int32(o.DisplaySequence)),
        Specialities:    utils.NewNullString(o.Specialities),
        Subspecialty:    utils.NewNullString(o.Subspecialty),
    }
}

type NovaDoctorSpecialtyDto struct {
    SpecialtyId      int  `json:"specialtyId" validate:"required,numeric"`
    PrimarySpecialty bool `json:"primarySpecialty" validate:"required"`
}

func (o NovaDoctorSpecialtyDto) ToDbModel() model.NovaDoctorSpecialty {
    primarySpecialty := 0
    if o.PrimarySpecialty {
        primarySpecialty = 1
    }
    return model.NovaDoctorSpecialty{
        SpecialtyId:      utils.NewInt64(int64(o.SpecialtyId)),
        PrimarySpecialtyV: utils.NewInt32(int32(primarySpecialty)),
        PrimarySpecialty: o.PrimarySpecialty,
    }
}

type NovaDoctorSpokenLanguageDto struct {
    DisplaySequence int    `json:"displaySequence" validate:"required,numeric"`
    SpokenLanguage  string `json:"spokenLanguage" validate:"required,max=200"`
}

func (o NovaDoctorSpokenLanguageDto) ToDbModel() model.NovaDoctorSpokenLanguage {
    return model.NovaDoctorSpokenLanguage{
        DisplaySequence: utils.NewInt32(int32(o.DisplaySequence)),
        SpokenLanguage:  utils.NewNullString(o.SpokenLanguage),
    }
}

type NovaDoctorDto struct {
    Mcr                  string                        `json:"mcr" validate:"required,max=100"`
    Name                 string                        `json:"name" validate:"required,max=200"`
    Gender               string                        `json:"gender" validate:"required,max=20"`
    ConsultantType       string                        `json:"consultantType" validate:"required,max=20"`
    Nationality          string                        `json:"nationality" validate:"required,max=200"`
    DisplaySequence      int                           `json:"displaySequence" validate:"numeric"`
    AllowAppointment     string                        `json:"allowAppointment" validate:"required,max=1"`
    Qualifications       string                        `json:"qualifications" validate:"max=200"`
    RegistrationNum      string                        `json:"registrationNum" validate:"max=200"`
    Subspecialty         string                        `json:"subspecialty" validate:"max=200"`
    IsForPackage         string                        `json:"isForPackage" validate:"required,max=1"`
    Image                string                        `json:"image"`
    DoctorSpokenLanguage []NovaDoctorSpokenLanguageDto `json:"doctorSpokenLanguage"`
    DoctorQualifications []NovaDoctorQualificationsDto `json:"doctorQualifications"`
    DoctorSpecialities   []NovaDoctorSpecialitiesDto   `json:"doctorSpecialities"`
    DoctorClinicLocation []NovaDoctorClinicLocationDto `json:"doctorClinicLocation"`
    DoctorClinicHours    []NovaDoctorClinicHoursDto    `json:"doctorClinicHours"`
    DoctorAppointment    []NovaDoctorAppointmentDto    `json:"doctorAppointment"`
    DoctorContact        []NovaDoctorContactDto        `json:"doctorContact"`
    DoctorSpecialty      []NovaDoctorSpecialtyDto      `json:"doctorSpecialty"`
}
