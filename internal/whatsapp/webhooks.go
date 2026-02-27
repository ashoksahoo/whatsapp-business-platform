package whatsapp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
)

// VerifySignature verifies the webhook signature using HMAC-SHA256
func VerifySignature(body []byte, signature string, secret string) bool {
	return utils.VerifyHMAC(body, []byte(secret), signature)
}

// ParseWebhook parses the webhook payload
func ParseWebhook(body []byte) (*WebhookPayload, error) {
	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}
	return &payload, nil
}

// ParseMessageEvent extracts message events from webhook payload
func ParseMessageEvent(payload *WebhookPayload) ([]*MessageEvent, error) {
	var events []*MessageEvent

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			for _, msg := range change.Value.Messages {
				event, err := parseMessageValue(&msg, change.Value.Contacts)
				if err != nil {
					return nil, err
				}
				events = append(events, event)
			}
		}
	}

	return events, nil
}

// parseMessageValue converts a MessageValue to MessageEvent
func parseMessageValue(msg *MessageValue, contacts []ContactValue) (*MessageEvent, error) {
	timestamp, err := strconv.ParseInt(msg.Timestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	event := &MessageEvent{
		MessageID: msg.ID,
		From:      msg.From,
		Timestamp: time.Unix(timestamp, 0),
		Type:      msg.Type,
	}

	// Extract contact name
	for _, contact := range contacts {
		if contact.WaID == msg.From {
			event.ContactName = contact.Profile.Name
			break
		}
	}

	// Extract content based on message type
	switch msg.Type {
	case "text":
		if msg.Text != nil {
			event.Content = msg.Text.Body
		}

	case "image":
		if msg.Image != nil {
			event.MediaID = msg.Image.ID
			event.MimeType = msg.Image.MimeType
			event.Caption = msg.Image.Caption
		}

	case "document":
		if msg.Document != nil {
			event.MediaID = msg.Document.ID
			event.MimeType = msg.Document.MimeType
			event.Caption = msg.Document.Caption
			event.Filename = msg.Document.Filename
		}

	case "audio":
		if msg.Audio != nil {
			event.MediaID = msg.Audio.ID
			event.MimeType = msg.Audio.MimeType
		}

	case "video":
		if msg.Video != nil {
			event.MediaID = msg.Video.ID
			event.MimeType = msg.Video.MimeType
			event.Caption = msg.Video.Caption
		}
	}

	return event, nil
}

// ParseStatusEvent extracts status update events from webhook payload
func ParseStatusEvent(payload *WebhookPayload) ([]*StatusEvent, error) {
	var events []*StatusEvent

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			for _, status := range change.Value.Statuses {
				event, err := parseStatusValue(&status)
				if err != nil {
					return nil, err
				}
				events = append(events, event)
			}
		}
	}

	return events, nil
}

// parseStatusValue converts a StatusValue to StatusEvent
func parseStatusValue(status *StatusValue) (*StatusEvent, error) {
	timestamp, err := strconv.ParseInt(status.Timestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	event := &StatusEvent{
		MessageID:   status.ID,
		Status:      status.Status,
		Timestamp:   time.Unix(timestamp, 0),
		RecipientID: status.RecipientID,
	}

	// Extract error information if present
	if len(status.Errors) > 0 {
		event.ErrorCode = status.Errors[0].Code
		event.ErrorTitle = status.Errors[0].Title
		event.ErrorMsg = status.Errors[0].Message
	}

	return event, nil
}
