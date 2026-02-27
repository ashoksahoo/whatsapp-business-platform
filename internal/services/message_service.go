package services

import (
	"fmt"
	"time"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/models"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/repositories"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/whatsapp"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/validator"
	"go.uber.org/zap"
)

// MessageService handles message business logic
type MessageService struct {
	messageRepo  *repositories.MessageRepository
	contactRepo  *repositories.ContactRepository
	waClient     *whatsapp.Client
	logger       *zap.Logger
}

// NewMessageService creates a new message service
func NewMessageService(
	messageRepo *repositories.MessageRepository,
	contactRepo *repositories.ContactRepository,
	waClient *whatsapp.Client,
	logger *zap.Logger,
) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		contactRepo: contactRepo,
		waClient:    waClient,
		logger:      logger,
	}
}

// SendTextMessage sends a text message
func (s *MessageService) SendTextMessage(phone, content string) (*models.Message, error) {
	// Validate phone number
	if err := validator.ValidatePhoneNumber(phone); err != nil {
		return nil, errors.NewInvalidPhoneNumberError(phone)
	}

	// Validate content
	if err := validator.ValidateNotEmpty(content, "content"); err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	// Get or create contact
	_, err := s.contactRepo.GetOrCreate(phone)
	if err != nil {
		s.logger.Error("Failed to get/create contact", zap.Error(err))
		return nil, errors.NewDatabaseError(err)
	}

	// Send message via WhatsApp
	resp, err := s.waClient.SendTextMessage(phone, content)
	if err != nil {
		s.logger.Error("Failed to send WhatsApp message", zap.Error(err))
		return nil, err
	}

	// Create message record
	message := &models.Message{
		WhatsAppMessageID: resp.Messages[0].ID,
		FromNumber:        resp.Contacts[0].Input,
		ToNumber:          phone,
		Direction:         "outbound",
		MessageType:       models.MessageTypeText,
		Content:           content,
		Status:            models.MessageStatusSent,
		Timestamp:         time.Now().UTC(),
	}

	if err := s.messageRepo.Create(message); err != nil {
		s.logger.Error("Failed to save message", zap.Error(err))
		return nil, errors.NewDatabaseError(err)
	}

	// Update contact last message time
	s.contactRepo.UpdateLastMessage(phone, message.Timestamp)
	s.contactRepo.IncrementMessageCount(phone, 1)

	s.logger.Info("Message sent successfully",
		zap.String("message_id", message.ID),
		zap.String("phone", phone),
	)

	return message, nil
}

// SendMediaMessage sends a media message
func (s *MessageService) SendMediaMessage(phone, mediaURL, caption, mediaType string) (*models.Message, error) {
	// Validate phone number
	if err := validator.ValidatePhoneNumber(phone); err != nil {
		return nil, errors.NewInvalidPhoneNumberError(phone)
	}

	// Validate media URL
	if err := validator.ValidateURL(mediaURL); err != nil {
		return nil, errors.NewBadRequest("invalid media URL: " + err.Error())
	}

	// Validate message type
	if err := validator.ValidateMessageType(mediaType); err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	// Get or create contact
	_, err := s.contactRepo.GetOrCreate(phone)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	// Send message via WhatsApp
	resp, err := s.waClient.SendMediaMessage(phone, mediaURL, caption, whatsapp.MediaType(mediaType))
	if err != nil {
		s.logger.Error("Failed to send media message", zap.Error(err))
		return nil, err
	}

	// Create message record
	message := &models.Message{
		WhatsAppMessageID: resp.Messages[0].ID,
		FromNumber:        resp.Contacts[0].Input,
		ToNumber:          phone,
		Direction:         "outbound",
		MessageType:       mediaType,
		Content:           caption,
		MediaURL:          mediaURL,
		Status:            models.MessageStatusSent,
		Timestamp:         time.Now().UTC(),
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	// Update contact
	s.contactRepo.UpdateLastMessage(phone, message.Timestamp)
	s.contactRepo.IncrementMessageCount(phone, 1)

	return message, nil
}

// SendTemplateMessage sends a template message
func (s *MessageService) SendTemplateMessage(phone, templateName, language string, params []string) (*models.Message, error) {
	// Validate inputs
	if err := validator.ValidatePhoneNumber(phone); err != nil {
		return nil, errors.NewInvalidPhoneNumberError(phone)
	}

	// Get or create contact
	_, err := s.contactRepo.GetOrCreate(phone)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	// Send template message
	resp, err := s.waClient.SendTemplateMessage(phone, templateName, language, params)
	if err != nil {
		return nil, err
	}

	// Create message record
	message := &models.Message{
		WhatsAppMessageID: resp.Messages[0].ID,
		FromNumber:        resp.Contacts[0].Input,
		ToNumber:          phone,
		Direction:         "outbound",
		MessageType:       models.MessageTypeTemplate,
		Content:           fmt.Sprintf("Template: %s", templateName),
		Status:            models.MessageStatusSent,
		Timestamp:         time.Now().UTC(),
		Metadata: models.JSONMap{
			"template_name": templateName,
			"language":      language,
			"parameters":    params,
		},
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	// Update contact
	s.contactRepo.UpdateLastMessage(phone, message.Timestamp)
	s.contactRepo.IncrementMessageCount(phone, 1)

	return message, nil
}

// GetMessage gets a message by ID
func (s *MessageService) GetMessage(messageID string) (*models.Message, error) {
	var message models.Message
	if err := s.messageRepo.FindByID(messageID, &message); err != nil {
		return nil, errors.NewNotFound("Message", messageID)
	}
	return &message, nil
}

// ListMessages lists messages with filters and pagination
func (s *MessageService) ListMessages(filters map[string]interface{}, pagination *utils.Pagination) ([]*models.Message, error) {
	return s.messageRepo.ListWithFilters(filters, pagination)
}

// SearchMessages searches messages by content
func (s *MessageService) SearchMessages(query string, filters map[string]interface{}, pagination *utils.Pagination) ([]*models.Message, error) {
	return s.messageRepo.Search(query, filters, pagination)
}

// ProcessIncomingMessage processes an incoming message from webhook
func (s *MessageService) ProcessIncomingMessage(event *whatsapp.MessageEvent) error {
	s.logger.Info("Processing incoming message",
		zap.String("from", event.From),
		zap.String("type", event.Type),
	)

	// Get or create contact
	contact, err := s.contactRepo.GetOrCreate(event.From)
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	// Update contact name if provided
	if event.ContactName != "" && contact.Name != event.ContactName {
		s.contactRepo.UpdateFields(contact.ID, contact, map[string]interface{}{
			"name": event.ContactName,
		})
	}

	// Create message record
	message := &models.Message{
		WhatsAppMessageID: event.MessageID,
		FromNumber:        event.From,
		ToNumber:          "", // Will be set from webhook metadata
		Direction:         "inbound",
		MessageType:       event.Type,
		Content:           event.Content,
		MediaURL:          event.MediaURL,
		MediaMimeType:     event.MimeType,
		Status:            "received",
		Timestamp:         event.Timestamp,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return errors.NewDatabaseError(err)
	}

	// Update contact
	s.contactRepo.UpdateLastMessage(event.From, event.Timestamp)
	s.contactRepo.IncrementMessageCount(event.From, 1)
	s.contactRepo.UpdateUnreadCount(event.From, 1)

	return nil
}

// UpdateMessageStatus updates the status of a message
func (s *MessageService) UpdateMessageStatus(whatsappMessageID, status string) error {
	return s.messageRepo.UpdateStatus(whatsappMessageID, status)
}
