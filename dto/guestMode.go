package dto

type GuestGetReturningPatientDto struct {
    IdentificationNumber string `json:"identificationNumber" validate:"required"`
}

type GuestMakeNewPatientDto struct {
    FullName             string `json:"fullName" validate:"required"`
    Dob                  string `json:"dob" validate:"required"`
    IdentificationNumber string `json:"identificationNumber" validate:"required"`
    Gender               string `json:"gender" validate:"required"`
    MaritalStatus        string `json:"maritalStatus"`
    Nationality          string `json:"nationality" validate:"required"`
    Country              string `json:"country"`
    Address1             string `json:"address1" validate:"required"`
    Address2             string `json:"address2"`
    TownCity             string `json:"townCity"`
    Postcode             string `json:"postcode"`
    State                string `json:"state"`
    ContactNumber        string `json:"contactNumber" validate:"required"`
    Email                string `json:"email"`
}

type GuestMakeAppointmentDto struct {
    IsReturningPatient        bool   `json:"isReturningPatient" validate:"required"`
    ReturningPatientDocNumber string `json:"returningPatientDocNumber"`
    NewPatientFullname        string `json:"newPatientFullname"`
    NewPatientDOB             string `json:"newPatientDOB"`
    NewPatientDocNumber       string `json:"newPatientDocNumber"`
    NewPatientGender          string `json:"newPatientGender"`
    NewPatientMaritalStatus   string `json:"newPatientMaritalStatus"`
    NewPatientNationality     string `json:"newPatientNationality"`
    NewPatientCountry         string `json:"newPatientCountry"`
    NewPatientAddress1        string `json:"newPatientAddress1"`
    NewPatientAddress2        string `json:"newPatientAddress2"`
    NewPatientTownCity        string `json:"newPatientTownCity"`
    NewPatientState           string `json:"newPatientState"`
    NewPatientContactNumber   string `json:"newPatientContactNumber"`
    NewPatientEmailAddress    string `json:"newPatientEmailAddress"`
}
