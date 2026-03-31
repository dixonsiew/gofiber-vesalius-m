package feedback

import (
    "github.com/guregu/null/v6"
)

type FeedbackAttachment struct {
    FeedbackAttachmentId null.Int64  `json:"feedback_attachment_id" db:"FEEDBACK_ATTACHMENT_ID"`
    FeedbackID           null.Int64  `json:"feedback_id" db:"FEEDBACK_ID"`
    AttachmentType       null.String `json:"attachmentType" db:"ATTACHMENT_TYPE"`
    AttachmentContent    []byte      `json:"attachmentContent" db:"ATTACHMENT_CONTENT"`
    AttachmentFilename   null.String `json:"attachmentFilename" db:"ATTACHMENT_FILENAME"`
}

type Feedback struct {
    FeedbackId             null.Int64  `json:"feedback_id" db:"FEEDBACK_ID"`
    PatientPrn             null.String `json:"patientPrn" db:"PATIENT_PRN"`
    AccountNo              null.String `json:"accountNo" db:"ACCOUNT_NO"`
    VisitType              null.String `json:"visitType" db:"VISIT_TYPE"`
    OverallRating          null.Int32  `json:"overallRating" db:"RATE_OVERALL_SATISFACTION"`
    HospitalServiceRating  null.Int32  `json:"hospitalServiceRating" db:"RATE_HOSPITAL_SERVICE"`
    StaffServiceRating     null.Int32  `json:"staffServiceRating" db:"RATE_STAFF_SERVICE"`
    ApptSchedulingRating   null.Int32  `json:"apptSchedulingRating" db:"RATE_APPT_SCHEDULING"`
    FoodBeveragesRating    null.Int32  `json:"foodBeveragesRating" db:"RATE_FOOD_BEVERAGES"`
    PaymentBillingRating   null.Int32  `json:"paymentBillingRating" db:"RATE_PAYMENT_BILLING"`
    RecommendUsRating      null.Int32  `json:"recommendUsRating" db:"RATE_RECOMMEND_US" `
    FeedbackDesc           null.String `json:"feedbackDesc" db:"FEEDBACK_DESC" `
    SubmittedDateTime      null.String `json:"submittedDateTime" db:"DATE_SUBMIT"`
    HasAttachmentV         null.String `json:"-" db:"HAS_ATTACHMENT"`
    HasAttachment          bool        `json:"hasAttachment"`
    SubmittedDateTimeExcel null.String `json:"submittedDateTimeExcel"`
}
