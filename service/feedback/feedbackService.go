package feedback

import (
	"context"
	"database/sql"
	"strings"
	"vesaliusm/database"
	"vesaliusm/model"
	"vesaliusm/model/feedback"
	"vesaliusm/utils"

	"github.com/jmoiron/sqlx"
	"github.com/nleeper/goment"
	go_ora "github.com/sijms/go-ora/v2"
)

var FeedbackSvc *FeedbackService = NewFeedbackService(database.GetDb(), database.GetCtx())

type FeedbackService struct {
    db  *sqlx.DB
    ctx context.Context
}

func NewFeedbackService(db *sqlx.DB, ctx context.Context) *FeedbackService {
    return &FeedbackService{
        db:  db,
        ctx: ctx,
    }
}

func (s *FeedbackService) GetAttachmentByAttachmentId(attachmentId int64) (*feedback.FeedbackAttachment, error) {
    query := `
        SELECT ATTACHMENT_CONTENT, ATTACHMENT_FILENAME, ATTACHMENT_TYPE
        FROM PATIENT_FEEDBACK_ATTACHMENT 
        WHERE FEEDBACK_ATTACHMENT_ID = :id
    `
    var o feedback.FeedbackAttachment
    err := s.db.GetContext(s.ctx, &o, query, attachmentId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    return &o, err
}

func (s *FeedbackService) FindByFeedbackId(feedbackId int64) (*feedback.Feedback, error) {
    query := `SELECT * FROM PATIENT_FEEDBACK WHERE FEEDBACK_ID = :id`
    query = strings.Replace(query, "*", utils.GetDbCols(feedback.Feedback{}, ""), 1)
    var o feedback.Feedback
    err := s.db.GetContext(s.ctx, &o, query, feedbackId)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        utils.LogError(err)
        return nil, err
    }
    o.Set()
    return &o, err
}

func (s *FeedbackService) FindAttachmentByFeedbackId(feedbackId int64) ([]feedback.FeedbackAttachment, error) {
    query := `SELECT * FROM PATIENT_FEEDBACK_ATTACHMENT WHERE FEEDBACK_ID = :id`
    query = strings.Replace(query, "*", utils.GetDbCols(feedback.FeedbackAttachment{}, ""), 1)
    list := make([]feedback.FeedbackAttachment, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return list, nil
}

func (s *FeedbackService) List(page string, limit string) (*model.PagedList, error) {
    total, err := s.Count(s.db)
    if err != nil {
        return nil, err
    }
    pager := model.GetPager(total, page, limit)
    list, err := s.FindAll(pager.GetLowerBound(), pager.PageSize, s.db)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    return &model.PagedList{
        List:       list,
        Total:      total,
        TotalPages: pager.GetTotalPages(),
    }, nil
}

func (s *FeedbackService) Count(conn *sqlx.DB) (int, error) {
    db := database.GetFromCon(conn, s.db)
    var count int
    query := `SELECT COUNT(FEEDBACK_ID) AS COUNT FROM PATIENT_FEEDBACK`
    err := db.GetContext(s.ctx, &count, query)
    if err != nil {
        utils.LogError(err)
        return 0, err
    }
    return count, nil
}

func (s *FeedbackService) FindAllExport() ([]feedback.Feedback, error) {
    query := `SELECT * FROM PATIENT_FEEDBACK ORDER BY DATE_SUBMIT DESC`
    query = strings.Replace(query, "*", utils.GetDbCols(feedback.Feedback{}, ""), 1)
    list := make([]feedback.Feedback, 0)
    err := s.db.SelectContext(s.ctx, &list, query)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
        if list[i].SubmittedDateTimeExcel.Valid {
            g, _ := goment.New(list[i].SubmittedDateTimeExcel, "YYYY-MM-DD[T]HH:mm:ssZ")
            list[i].SubmittedDateTimeExcel = utils.NewNullString(g.Format("DD/MM/YYYY HH:mm"))
        }
    }
    return list, nil
}

func (s *FeedbackService) FindAll(offset int, limit int, conn *sqlx.DB) ([]feedback.Feedback, error) {
    db := database.GetFromCon(conn, s.db)
    query := `SELECT * FROM PATIENT_FEEDBACK ORDER BY DATE_SUBMIT DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
    query = strings.Replace(query, "*", utils.GetDbCols(feedback.Feedback{}, ""), 1)
    list := make([]feedback.Feedback, 0)
    err := db.SelectContext(s.ctx, &list, query, offset, limit)
    if err != nil {
        utils.LogError(err)
        return nil, err
    }
    for i := range list {
        list[i].Set()
    }
    return list, nil
}

func (s *FeedbackService) Save(o feedback.Feedback, la []feedback.FeedbackAttachment) error {
    tx, err := s.db.BeginTxx(s.ctx, nil)
    if err != nil {
        utils.LogError(err)
        return err
    }
    defer func() {
        if err != nil {
            utils.LogError(err)
            tx.Rollback()
        }
    }()

    query := `
        INSERT INTO PATIENT_FEEDBACK
        (
          PATIENT_PRN, ACCOUNT_NO, VISIT_TYPE, RATE_OVERALL_SATISFACTION, RATE_HOSPITAL_SERVICE,
          RATE_STAFF_SERVICE, RATE_APPT_SCHEDULING, RATE_FOOD_BEVERAGES, RATE_PAYMENT_BILLING,
          RATE_RECOMMEND_US, FEEDBACK_DESC, HAS_ATTACHMENT
        ) VALUES (
          :patientPrn, :accountNo, :visitType, :overallRating, :hospitalServiceRating,
          :staffServiceRating, :apptSchedulingRating, :foodBeveragesRating, :paymentBillingRating,
          :recommendUsRating, :feedbackDesc, :hasAttachment
        ) RETURNING FEEDBACK_ID INTO :feedback_id
    `
    hasAttachment := "N"
    if o.HasAttachment {
        hasAttachment = "Y"
    }
    var feedbackId go_ora.Number
    _, err = tx.ExecContext(s.ctx, query,
        sql.Named("patientPrn", o.PatientPrn.String),
        sql.Named("accountNo", o.AccountNo.String),
        sql.Named("visitType", o.VisitType.String),
        sql.Named("overallRating", o.OverallRating.Int32),
        sql.Named("hospitalServiceRating", o.HospitalServiceRating.Int32),
        sql.Named("staffServiceRating", o.StaffServiceRating.Int32),
        sql.Named("apptSchedulingRating", o.ApptSchedulingRating.Int32),
        sql.Named("foodBeveragesRating", o.FoodBeveragesRating.Int32),
        sql.Named("paymentBillingRating", o.PaymentBillingRating.Int32),
        sql.Named("recommendUsRating", o.RecommendUsRating.Int32),
        sql.Named("feedbackDesc", o.FeedbackDesc.String),
        sql.Named("hasAttachment", hasAttachment),
        go_ora.Out{Dest: &feedbackId},
    )
    if err != nil {
        utils.LogError(err)
        return err
    }
    
    ifeedbackId, _ := feedbackId.Int64()
    o.FeedbackId.Int64 = ifeedbackId
    
    for _, a := range la {
        q := `
           INSERT INTO PATIENT_FEEDBACK_ATTACHMENT
           (
            FEEDBACK_ID, ATTACHMENT_TYPE, ATTACHMENT_CONTENT, ATTACHMENT_FILENAME
           ) VALUES (
            :feedbackId, :attachmentType, :attachmentContent, :attachmentFilename
           )
        `
        _, err = tx.ExecContext(s.ctx, q,
            sql.Named("feedbackId", o.FeedbackId.Int64),
            sql.Named("attachmentType", a.AttachmentType.String),
            sql.Named("attachmentContent", a.AttachmentContent),
            sql.Named("attachmentFilename", a.AttachmentFilename.String),
           )
        if err != nil {
            utils.LogError(err)
            return err
        }
    }

    return tx.Commit()
}
