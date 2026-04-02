package dto

type FeedbackFileDto struct {
    Filename string `json:"filename" validate:"required"`
    MimeType string `json:"mimeType" validate:"required"`
}

type FeedbackDto struct {
    PatientPrn            string            `json:"patientPrn" validate:"required"`
    AccountNo             string            `json:"accountNo" validate:"required"`
    VisitType             string            `json:"visitType" validate:"required"`
    OverallRating         int               `json:"overallRating" validate:"required,numeric"`
    HospitalServiceRating int               `json:"hospitalServiceRating" validate:"required,numeric"`
    StaffServiceRating    int               `json:"staffServiceRating" validate:"required,numeric"`
    ApptSchedulingRating  int               `json:"apptSchedulingRating" validate:"required,numeric"`
    FoodBeveragesRating   int               `json:"foodBeveragesRating" validate:"required,numeric"`
    PaymentBillingRating  int               `json:"paymentBillingRating" validate:"required,numeric"`
    RecommendUsRating     int               `json:"recommendUsRating" validate:"required,numeric"`
    FeedbackDesc          string            `json:"feedbackDesc" validate:"required"`
    FeedbackFiles         []FeedbackFileDto `json:"feedbackFiles"`
}
