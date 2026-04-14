package feedback

import (
    "vesaliusm/middleware"

    "github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    feedbackController := NewFeedbackController()
    feedbackController.registerRoutes(router)
}

func (c *FeedbackController) registerRoutes(router fiber.Router) {
    api := router.Group("/feedback")
    api.Post("/", c.CreateFeedback)
    
    api.Use(middleware.JWTProtected, middleware.ValidateUser)
    api.Get("/export", c.GetAllExportFeedbacks)
    api.Get("/all", c.GetAllFeedbacks)
    api.Get("/:feedbackId", c.GetFeedbackById)
    api.Get("/attachment/:feedbackId", c.GetFeedbackAttachmentById)
    api.Get("/attachment-download/:attachmentId", c.DownloadAttachmentById)
}
