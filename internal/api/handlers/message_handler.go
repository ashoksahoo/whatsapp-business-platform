package handlers

import (
	"strconv"
	"time"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/services"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// MessageHandler handles message-related requests
type MessageHandler struct {
	messageService *services.MessageService
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(messageService *services.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// SendMessageRequest represents the request body for sending a message
type SendMessageRequest struct {
	Phone            string   `json:"phone" binding:"required"`
	Type             string   `json:"type" binding:"required"`
	Content          string   `json:"content"`
	MediaURL         string   `json:"media_url"`
	Caption          string   `json:"caption"`
	Filename         string   `json:"filename"`
	TemplateName     string   `json:"template_name"`
	TemplateLanguage string   `json:"template_language"`
	Parameters       []string `json:"parameters"`
}

// SendMessage handles POST /api/v1/messages
func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorJSON(c, errors.NewBadRequest("Invalid request body: "+err.Error()))
		return
	}

	var message interface{}
	var err error

	switch req.Type {
	case "text":
		message, err = h.messageService.SendTextMessage(req.Phone, req.Content)

	case "image", "video", "audio", "document":
		message, err = h.messageService.SendMediaMessage(req.Phone, req.MediaURL, req.Caption, req.Type)

	case "template":
		message, err = h.messageService.SendTemplateMessage(req.Phone, req.TemplateName, req.TemplateLanguage, req.Parameters)

	default:
		utils.ErrorJSON(c, errors.NewBadRequest("Invalid message type: "+req.Type))
		return
	}

	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorJSON(c, appErr)
		} else {
			utils.ErrorJSON(c, errors.NewInternalError(err))
		}
		return
	}

	utils.CreatedJSON(c, message)
}

// GetMessage handles GET /api/v1/messages/:id
func (h *MessageHandler) GetMessage(c *gin.Context) {
	messageID := c.Param("id")

	message, err := h.messageService.GetMessage(messageID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorJSON(c, appErr)
		} else {
			utils.ErrorJSON(c, errors.NewInternalError(err))
		}
		return
	}

	utils.SuccessJSON(c, 200, message)
}

// ListMessages handles GET /api/v1/messages
func (h *MessageHandler) ListMessages(c *gin.Context) {
	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	pagination := utils.NewPagination(limit, offset)

	// Build filters
	filters := make(map[string]interface{})

	if phone := c.Query("phone"); phone != "" {
		filters["phone"] = phone
	}
	if direction := c.Query("direction"); direction != "" {
		filters["direction"] = direction
	}
	if msgType := c.Query("type"); msgType != "" {
		filters["type"] = msgType
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			filters["start_date"] = t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			filters["end_date"] = t
		}
	}

	messages, err := h.messageService.ListMessages(filters, pagination)
	if err != nil {
		utils.ErrorJSON(c, errors.NewInternalError(err))
		return
	}

	utils.ListJSON(c, messages, pagination)
}

// SearchMessages handles GET /api/v1/messages/search
func (h *MessageHandler) SearchMessages(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorJSON(c, errors.NewBadRequest("Search query 'q' is required"))
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	pagination := utils.NewPagination(limit, offset)

	filters := make(map[string]interface{})
	if phone := c.Query("phone"); phone != "" {
		filters["phone"] = phone
	}

	messages, err := h.messageService.SearchMessages(query, filters, pagination)
	if err != nil {
		utils.ErrorJSON(c, errors.NewInternalError(err))
		return
	}

	utils.ListJSON(c, messages, pagination)
}
