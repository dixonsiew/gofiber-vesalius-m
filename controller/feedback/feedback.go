package feedback

import (
    "bytes"
    "fmt"
    "image"
    "image/png"
    "io"
    "strconv"
    "vesaliusm/dto"
    model "vesaliusm/model/feedback"
    "vesaliusm/service/feedback"
    "vesaliusm/service/mail"
    "vesaliusm/utils"

    "github.com/gofiber/fiber/v3"
    "github.com/nfnt/resize"
)

type FeedbackController struct {
    feedbackService *feedback.FeedbackService
    mailService     *mail.MailService
}

func NewFeedbackController() *FeedbackController {
    return &FeedbackController{
        feedbackService: feedback.FeedbackSvc,
        mailService:     mail.MailSvc,
    }
}

// GetAllFeedbacks
//
// @Tags Feedback
// @Produce json
// @Security BearerAuth
// @Param        _page             query       string  false  "_page"  default:"1"
// @Param        _limit            query       string  false  "_limit" default:"10"
// @Success 200 {array} feedback.FeedbackAttachment
// @Router /feedback/all [get]
func (cr *FeedbackController) GetAllFeedbacks(c fiber.Ctx) error {
    page := c.Query("_page", "1")
    limit := c.Query("_limit", strconv.Itoa(utils.PAGE_SIZE))
    m, err := cr.feedbackService.List(page, limit)
    if err != nil {
        return err
    }

    c.Set(utils.X_TOTAL_COUNT, strconv.Itoa(m.Total))
    c.Set(utils.X_TOTAL_PAGE, strconv.Itoa(m.TotalPages))
    return c.JSON(m.List)
}

// GetFeedbackById
//
// @Tags Feedback
// @Produce json
// @Security BearerAuth
// @Param        feedbackId        path        string  true  "feedbackId"
// @Success 200 {object} feedback.Feedback
// @Router /feedback/{feedbackId} [get]
func (cr *FeedbackController) GetFeedbackById(c fiber.Ctx) error {
    feedbackId := c.Params("feedbackId")
    ifeedbackId, _ := strconv.ParseInt(feedbackId, 10, 64)
    o, err := cr.feedbackService.FindByFeedbackId(ifeedbackId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// GetFeedbackAttachmentById
//
// @Tags Feedback
// @Produce json
// @Security BearerAuth
// @Param        feedbackId        path        string  true  "feedbackId"
// @Success 200 {object} feedback.FeedbackAttachment
// @Router /feedback/attachment/{feedbackId} [get]
func (cr *FeedbackController) GetFeedbackAttachmentById(c fiber.Ctx) error {
    feedbackId := c.Params("feedbackId")
    ifeedbackId, _ := strconv.ParseInt(feedbackId, 10, 64)
    o, err := cr.feedbackService.FindAttachmentByFeedbackId(ifeedbackId)
    if err != nil {
        return err
    }

    return c.JSON(o)
}

// DownloadAttachmentById
//
// @Tags Feedback
// @Produce json
// @Security BearerAuth
// @Param        attachmentId        path        string  true  "attachmentId"
// @Success 200 {file} binary
// @Router /feedback/attachment-download/{attachmentId} [get]
func (cr *FeedbackController) DownloadAttachmentById(c fiber.Ctx) error {
    attachmentId := c.Params("attachmentId")
    iattachmentId, _ := strconv.ParseInt(attachmentId, 10, 64)
    o, err := cr.feedbackService.GetAttachmentByAttachmentId(iattachmentId)
    if err != nil {
        return err
    }

    c.Set(fiber.HeaderCacheControl, "no-cache, no-store, must-revalidate")
    c.Set(fiber.HeaderPragma, "no-cache")
    c.Set(fiber.HeaderExpires, "0")
    c.Set(fiber.HeaderExpires, "0")
    c.Set("filename", o.AttachmentFilename.String)
    c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", o.AttachmentFilename.String))
    c.Set(fiber.HeaderContentType, o.AttachmentType.String)
    return c.Send(o.AttachmentContent)
}

// CreateFeedback
//
// @Tags Feedback
// @Accept multipart/form-data
// @Produce json
// @Param files formData file false "Files"
// @Param request body dto.FeedbackDto true "FeedbackDto"
// @Success 200
// @Router /feedback [post]
func (cr *FeedbackController) CreateFeedback(c fiber.Ctx) error {
    form, err := c.MultipartForm()
    if err != nil {
        return err
    }
    files := form.File["files"]
    hasAttachment := len(files) > 0

    data := new(dto.FeedbackDto)
    if err := utils.BindNValidate(c, data); err != nil {
        return err
    }

    o := model.Feedback{
        PatientPrn:            utils.NewNullString(data.PatientPrn),
        AccountNo:             utils.NewNullString(data.AccountNo),
        VisitType:             utils.NewNullString(data.VisitType),
        OverallRating:         utils.NewInt32(int32(data.OverallRating)),
        HospitalServiceRating: utils.NewInt32(int32(data.HospitalServiceRating)),
        StaffServiceRating:    utils.NewInt32(int32(data.StaffServiceRating)),
        ApptSchedulingRating:  utils.NewInt32(int32(data.ApptSchedulingRating)),
        FoodBeveragesRating:   utils.NewInt32(int32(data.FoodBeveragesRating)),
        PaymentBillingRating:  utils.NewInt32(int32(data.PaymentBillingRating)),
        RecommendUsRating:     utils.NewInt32(int32(data.RecommendUsRating)),
        FeedbackDesc:          utils.NewNullString(data.FeedbackDesc),
        HasAttachment:         hasAttachment,
    }

    la := make([]model.FeedbackAttachment, 0)

    if hasAttachment {
        for i, file := range files {
            fbFile := data.FeedbackFiles[i]
            mimeType := fbFile.MimeType
            filename := fbFile.Filename

            if mimeType != "application/pdf" && file.Size >= utils.FILE_SIZE_LIMIT {
                src, err := file.Open()
                if err != nil {
                    return err
                }
                defer src.Close()

                fileBytes, err := io.ReadAll(src)
                if err != nil {
                    return err
                }

                srcImg, _, err := image.Decode(bytes.NewReader(fileBytes))
                if err != nil {
                    return err
                }

                resizedImg := resize.Resize(3000, 4000, srcImg, resize.Lanczos3)
                var buf bytes.Buffer
                if err := png.Encode(&buf, resizedImg); err != nil {
                    return err
                }

                a := model.FeedbackAttachment{
                    AttachmentType:     utils.NewNullString(mimeType),
                    AttachmentContent:  buf.Bytes(),
                    AttachmentFilename: utils.NewNullString(filename),
                }
                la = append(la, a)
            }
        }
    }

    go func ()  {
        cr.mailService.SendPatientFeedbackSubmitted()
    }()
    err = cr.feedbackService.Save(o, la)
    if err != nil {
        return err
    }

    return c.JSON(fiber.Map{
        "message": "Feedback created",
    })
}
