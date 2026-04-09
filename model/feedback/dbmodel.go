package feedback

import (
	"vesaliusm/utils"

	"github.com/guregu/null/v6"
)

type FeedbackAttachment struct {
    FeedbackAttachmentId null.Int64  `json:"feedback_attachment_id" db:"FEEDBACK_ATTACHMENT_ID" swaggertype:"integer"`
    FeedbackId           null.Int64  `json:"feedback_id" db:"FEEDBACK_ID" swaggertype:"integer"`
    AttachmentType       null.String `json:"attachmentType" db:"ATTACHMENT_TYPE" swaggertype:"string"`
    AttachmentContent    []byte      `json:"attachmentContent" db:"ATTACHMENT_CONTENT" swaggertype:"string"`
    AttachmentFilename   null.String `json:"attachmentFilename" db:"ATTACHMENT_FILENAME" swaggertype:"string"`
}

type Feedback struct {
    FeedbackId             null.Int64  `json:"feedback_id" db:"FEEDBACK_ID" swaggertype:"integer"`
    PatientPrn             null.String `json:"patientPrn" db:"PATIENT_PRN" swaggertype:"string"`
    AccountNo              null.String `json:"accountNo" db:"ACCOUNT_NO" swaggertype:"string"`
    VisitType              null.String `json:"visitType" db:"VISIT_TYPE" swaggertype:"string"`
    OverallRating          null.Int32  `json:"overallRating" db:"RATE_OVERALL_SATISFACTION" swaggertype:"integer"`
    HospitalServiceRating  null.Int32  `json:"hospitalServiceRating" db:"RATE_HOSPITAL_SERVICE" swaggertype:"integer"`
    StaffServiceRating     null.Int32  `json:"staffServiceRating" db:"RATE_STAFF_SERVICE" swaggertype:"integer"`
    ApptSchedulingRating   null.Int32  `json:"apptSchedulingRating" db:"RATE_APPT_SCHEDULING" swaggertype:"integer"`
    FoodBeveragesRating    null.Int32  `json:"foodBeveragesRating" db:"RATE_FOOD_BEVERAGES" swaggertype:"integer"`
    PaymentBillingRating   null.Int32  `json:"paymentBillingRating" db:"RATE_PAYMENT_BILLING" swaggertype:"integer"`
    RecommendUsRating      null.Int32  `json:"recommendUsRating" db:"RATE_RECOMMEND_US" swaggertype:"integer"`
    FeedbackDesc           null.String `json:"feedbackDesc" db:"FEEDBACK_DESC" swaggertype:"string"`
    SubmittedDateTime      null.String `json:"submittedDateTime" db:"DATE_SUBMIT" swaggertype:"string"`
    HasAttachmentV         null.String `json:"-" db:"HAS_ATTACHMENT" swaggertype:"string"`
    HasAttachment          bool        `json:"hasAttachment"`
    SubmittedDateTimeExcel null.String `json:"submittedDateTimeExcel" swaggertype:"string"`
}

func (o *Feedback) Set() {
    if !o.AccountNo.Valid {
        o.AccountNo = utils.NewNullString("-")
    }

    if !o.FeedbackDesc.Valid {
        o.FeedbackDesc = utils.NewNullString("-")
    }

    if o.HasAttachmentV.String == "Y" {
        o.HasAttachment = true
    } else {
        o.HasAttachment = false
    }
}
