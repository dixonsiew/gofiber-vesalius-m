package feedback

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/godror/godror" // Oracle driver
)

// Feedback represents the PATIENT_FEEDBACK record.
type Feedback struct {
	FeedbackID              int       `db:"FEEDBACK_ID"`
	PatientPrn              string    `db:"PATIENT_PRN"`
	AccountNo               string    `db:"ACCOUNT_NO"`
	VisitType               string    `db:"VISIT_TYPE"`
	OverallRating           int       `db:"RATE_OVERALL_SATISFACTION"`
	HospitalServiceRating   int       `db:"RATE_HOSPITAL_SERVICE"`
	StaffServiceRating      int       `db:"RATE_STAFF_SERVICE"`
	ApptSchedulingRating    int       `db:"RATE_APPT_SCHEDULING"`
	FoodBeveragesRating     int       `db:"RATE_FOOD_BEVERAGES"`
	PaymentBillingRating    int       `db:"RATE_PAYMENT_BILLING"`
	RecommendUsRating       int       `db:"RATE_RECOMMEND_US"`
	FeedbackDesc            string    `db:"FEEDBACK_DESC"`
	HasAttachment           string    `db:"HAS_ATTACHMENT"` // "Y" or "N"
	DateSubmit              time.Time `db:"DATE_SUBMIT"`
}

// SubmittedDateTimeExcel returns the formatted date for Excel export.
func (f *Feedback) SubmittedDateTimeExcel() string {
	return f.DateSubmit.Format("02/01/2006 15:04")
}

// FeedbackAttachment represents the PATIENT_FEEDBACK_ATTACHMENT record.
type FeedbackAttachment struct {
	FeedbackAttachmentID int       `db:"FEEDBACK_ATTACHMENT_ID"`
	FeedbackID           int       `db:"FEEDBACK_ID"`
	AttachmentType       string    `db:"ATTACHMENT_TYPE"`
	AttachmentContent    []byte    `db:"ATTACHMENT_CONTENT"` // BLOB
	AttachmentFilename   string    `db:"ATTACHMENT_FILENAME"`
	DateCreated          time.Time `db:"DATE_CREATED"`
}

// PagedList holds paginated results.
type PagedList struct {
	Items      interface{} `json:"items"`
	TotalCount int         `json:"totalCount"`
	TotalPages int         `json:"totalPages"`
}

// Pager calculates pagination bounds.
type Pager struct {
	LowerBound int
	PageSize   int
	TotalPages int
}

// NewPager creates a new Pager from total, page, limit.
func NewPager(total, page, limit int) Pager {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	totalPages := (total + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}
	offset := (page - 1) * limit
	return Pager{
		LowerBound: offset,
		PageSize:   limit,
		TotalPages: totalPages,
	}
}

// FeedbackService provides feedback operations.
type FeedbackService struct {
	db *sqlx.DB
}

// NewFeedbackService creates a new service with the given DB connection.
func NewFeedbackService(db *sqlx.DB) *FeedbackService {
	return &FeedbackService{db: db}
}

// GetAttachmentByAttachmentID retrieves a single attachment by its ID.
func (s *FeedbackService) GetAttachmentByAttachmentID(ctx context.Context, attachmentID int) (*FeedbackAttachment, error) {
	var att FeedbackAttachment
	query := `SELECT ATTACHMENT_CONTENT, ATTACHMENT_FILENAME, ATTACHMENT_TYPE
	          FROM PATIENT_FEEDBACK_ATTACHMENT
	          WHERE FEEDBACK_ATTACHMENT_ID = :1`
	err := s.db.GetContext(ctx, &att, query, attachmentID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get attachment by id: %w", err)
	}
	return &att, nil
}

// FindByFeedbackID retrieves a feedback record by its ID.
func (s *FeedbackService) FindByFeedbackID(ctx context.Context, feedbackID int) (*Feedback, error) {
	var fb Feedback
	query := `SELECT * FROM PATIENT_FEEDBACK WHERE FEEDBACK_ID = :1`
	err := s.db.GetContext(ctx, &fb, query, feedbackID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find by feedback id: %w", err)
	}
	return &fb, nil
}

// FindAttachmentByFeedbackID retrieves all attachments for a given feedback ID.
func (s *FeedbackService) FindAttachmentByFeedbackID(ctx context.Context, feedbackID int) ([]FeedbackAttachment, error) {
	var attachments []FeedbackAttachment
	query := `SELECT * FROM PATIENT_FEEDBACK_ATTACHMENT WHERE FEEDBACK_ID = :1`
	err := s.db.SelectContext(ctx, &attachments, query, feedbackID)
	if err != nil {
		return nil, fmt.Errorf("find attachments by feedback id: %w", err)
	}
	return attachments, nil
}

// count returns the total number of feedback records.
func (s *FeedbackService) count(ctx context.Context) (int, error) {
	var total int
	query := `SELECT COUNT(FEEDBACK_ID) FROM PATIENT_FEEDBACK`
	err := s.db.GetContext(ctx, &total, query)
	if err != nil {
		return 0, fmt.Errorf("count feedback: %w", err)
	}
	return total, nil
}

// findAll retrieves a paginated slice of feedback records.
func (s *FeedbackService) findAll(ctx context.Context, offset, limit int) ([]Feedback, error) {
	var feedbacks []Feedback
	query := `SELECT * FROM PATIENT_FEEDBACK ORDER BY DATE_SUBMIT DESC
	          OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY`
	err := s.db.SelectContext(ctx, &feedbacks, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("find all paginated: %w", err)
	}
	return feedbacks, nil
}

// List returns a paged list of feedback records.
func (s *FeedbackService) List(ctx context.Context, page, limit int) (*PagedList, error) {
	total, err := s.count(ctx)
	if err != nil {
		return nil, err
	}
	pager := NewPager(total, page, limit)
	items, err := s.findAll(ctx, pager.LowerBound, pager.PageSize)
	if err != nil {
		return nil, err
	}
	return &PagedList{
		Items:      items,
		TotalCount: total,
		TotalPages: pager.TotalPages,
	}, nil
}

// FindAllExport returns all feedback records ordered by date submit descending,
// with the date already formatted for Excel.
func (s *FeedbackService) FindAllExport(ctx context.Context) ([]Feedback, error) {
	var feedbacks []Feedback
	query := `SELECT * FROM PATIENT_FEEDBACK ORDER BY DATE_SUBMIT DESC`
	err := s.db.SelectContext(ctx, &feedbacks, query)
	if err != nil {
		return nil, fmt.Errorf("find all export: %w", err)
	}
	return feedbacks, nil
}

// Save inserts a feedback record along with its attachments in a transaction.
func (s *FeedbackService) Save(ctx context.Context, fb *Feedback, attachments []FeedbackAttachment) error {
	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	hasAttachment := "N"
	if len(attachments) > 0 {
		hasAttachment = "Y"
	}

	// Insert feedback and retrieve the generated ID.
	var feedbackID int
	insertFeedback := `
		INSERT INTO PATIENT_FEEDBACK (
			PATIENT_PRN, ACCOUNT_NO, VISIT_TYPE,
			RATE_OVERALL_SATISFACTION, RATE_HOSPITAL_SERVICE, RATE_STAFF_SERVICE,
			RATE_APPT_SCHEDULING, RATE_FOOD_BEVERAGES, RATE_PAYMENT_BILLING,
			RATE_RECOMMEND_US, FEEDBACK_DESC, HAS_ATTACHMENT
		) VALUES (
			:patientPrn, :accountNo, :visitType,
			:overallRating, :hospitalServiceRating, :staffServiceRating,
			:apptSchedulingRating, :foodBeveragesRating, :paymentBillingRating,
			:recommendUsRating, :feedbackDesc, :hasAttachment
		) RETURNING FEEDBACK_ID INTO :feedbackID`

	params := map[string]interface{}{
		"patientPrn":            fb.PatientPrn,
		"accountNo":             fb.AccountNo,
		"visitType":             fb.VisitType,
		"overallRating":         fb.OverallRating,
		"hospitalServiceRating": fb.HospitalServiceRating,
		"staffServiceRating":    fb.StaffServiceRating,
		"apptSchedulingRating":  fb.ApptSchedulingRating,
		"foodBeveragesRating":   fb.FoodBeveragesRating,
		"paymentBillingRating":  fb.PaymentBillingRating,
		"recommendUsRating":     fb.RecommendUsRating,
		"feedbackDesc":          fb.FeedbackDesc,
		"hasAttachment":         hasAttachment,
		"feedbackID":            sql.Out{Dest: &feedbackID},
	}

	_, err = tx.NamedExecContext(ctx, insertFeedback, params)
	if err != nil {
		return fmt.Errorf("insert feedback: %w", err)
	}
	fb.FeedbackID = feedbackID

	// Insert attachments.
	for i := range attachments {
		att := attachments[i]
		att.FeedbackID = feedbackID
		insertAtt := `
			INSERT INTO PATIENT_FEEDBACK_ATTACHMENT (
				FEEDBACK_ID, ATTACHMENT_TYPE, ATTACHMENT_CONTENT, ATTACHMENT_FILENAME
			) VALUES (
				:feedbackID, :attachmentType, :attachmentContent, :attachmentFilename
			)`
		_, err = tx.NamedExecContext(ctx, insertAtt, att)
		if err != nil {
			return fmt.Errorf("insert attachment: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}