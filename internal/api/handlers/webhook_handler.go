package handlers

import (
	"io"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/services"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/whatsapp"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// WebhookHandler handles webhook-related requests
type WebhookHandler struct {
	messageService *services.MessageService
	verifyToken    string
	webhookSecret  string
	logger         *zap.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(
	messageService *services.MessageService,
	verifyToken string,
	webhookSecret string,
	logger *zap.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		messageService: messageService,
		verifyToken:    verifyToken,
		webhookSecret:  webhookSecret,
		logger:         logger,
	}
}

// VerifyWebhook handles GET /webhooks/whatsapp for webhook verification
func (h *WebhookHandler) VerifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == h.verifyToken {
		h.logger.Info("Webhook verified successfully")
		c.String(200, challenge)
		return
	}

	h.logger.Warn("Webhook verification failed",
		zap.String("mode", mode),
		zap.String("token", token),
	)
	utils.ErrorJSON(c, errors.NewUnauthorized("Webhook verification failed"))
}

// ReceiveWebhook handles POST /webhooks/whatsapp for receiving events
func (h *WebhookHandler) ReceiveWebhook(c *gin.Context) {
	// Read body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error("Failed to read webhook body", zap.Error(err))
		utils.ErrorJSON(c, errors.NewBadRequest("Failed to read request body"))
		return
	}

	// Verify signature
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature != "" && !whatsapp.VerifySignature(body, signature, h.webhookSecret) {
		h.logger.Warn("Webhook signature verification failed")
		utils.ErrorJSON(c, errors.NewUnauthorized("Invalid signature"))
		return
	}

	// Parse webhook payload
	payload, err := whatsapp.ParseWebhook(body)
	if err != nil {
		h.logger.Error("Failed to parse webhook", zap.Error(err))
		utils.ErrorJSON(c, errors.NewBadRequest("Invalid webhook payload"))
		return
	}

	// Process message events
	messageEvents, err := whatsapp.ParseMessageEvent(payload)
	if err != nil {
		h.logger.Error("Failed to parse message events", zap.Error(err))
	} else {
		for _, event := range messageEvents {
			if err := h.messageService.ProcessIncomingMessage(event); err != nil {
				h.logger.Error("Failed to process incoming message",
					zap.Error(err),
					zap.String("message_id", event.MessageID),
				)
			}
		}
	}

	// Process status events
	statusEvents, err := whatsapp.ParseStatusEvent(payload)
	if err != nil {
		h.logger.Error("Failed to parse status events", zap.Error(err))
	} else {
		for _, event := range statusEvents {
			if err := h.messageService.UpdateMessageStatus(event.MessageID, event.Status); err != nil {
				h.logger.Error("Failed to update message status",
					zap.Error(err),
					zap.String("message_id", event.MessageID),
				)
			}
		}
	}

	// Return success
	c.JSON(200, gin.H{"status": "received"})
}
