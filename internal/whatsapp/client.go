package whatsapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// Config holds WhatsApp client configuration
type Config struct {
	APIToken      string
	PhoneNumberID string
	APIBaseURL    string
	APIVersion    string
	Logger        *zap.Logger
}

// Client represents a WhatsApp API client
type Client struct {
	httpClient    *resty.Client
	phoneNumberID string
	baseURL       string
	logger        *zap.Logger
}

// NewClient creates a new WhatsApp client
func NewClient(config Config) (*Client, error) {
	if config.APIToken == "" {
		return nil, fmt.Errorf("API token is required")
	}
	if config.PhoneNumberID == "" {
		return nil, fmt.Errorf("phone number ID is required")
	}

	baseURL := config.APIBaseURL
	if baseURL == "" {
		baseURL = "https://graph.facebook.com"
	}

	apiVersion := config.APIVersion
	if apiVersion == "" {
		apiVersion = "v18.0"
	}

	httpClient := resty.New()
	httpClient.SetBaseURL(fmt.Sprintf("%s/%s", baseURL, apiVersion))
	httpClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", config.APIToken))
	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetTimeout(30 * time.Second)
	httpClient.SetRetryCount(3)
	httpClient.SetRetryWaitTime(1 * time.Second)
	httpClient.SetRetryMaxWaitTime(5 * time.Second)

	return &Client{
		httpClient:    httpClient,
		phoneNumberID: config.PhoneNumberID,
		baseURL:       fmt.Sprintf("%s/%s/%s", baseURL, apiVersion, config.PhoneNumberID),
		logger:        config.Logger,
	}, nil
}

// SendTextMessage sends a text message
func (c *Client) SendTextMessage(to, text string) (*MessageResponse, error) {
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"body": text,
		},
	}

	return c.sendMessage(payload)
}

// SendMediaMessage sends a media message (image, document, audio, video)
func (c *Client) SendMediaMessage(to, mediaURL, caption string, mediaType MediaType) (*MessageResponse, error) {
	mediaObj := map[string]interface{}{
		"link": mediaURL,
	}

	if caption != "" && (mediaType == MediaTypeImage || mediaType == MediaTypeVideo || mediaType == MediaTypeDocument) {
		mediaObj["caption"] = caption
	}

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              string(mediaType),
		string(mediaType):   mediaObj,
	}

	return c.sendMessage(payload)
}

// SendTemplateMessage sends a template message
func (c *Client) SendTemplateMessage(to, templateName, language string, params []string) (*MessageResponse, error) {
	components := []map[string]interface{}{}

	if len(params) > 0 {
		parameters := make([]map[string]interface{}, len(params))
		for i, param := range params {
			parameters[i] = map[string]interface{}{
				"type": "text",
				"text": param,
			}
		}

		components = append(components, map[string]interface{}{
			"type":       "body",
			"parameters": parameters,
		})
	}

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "template",
		"template": map[string]interface{}{
			"name": templateName,
			"language": map[string]string{
				"code": language,
			},
			"components": components,
		},
	}

	return c.sendMessage(payload)
}

// sendMessage sends a message to WhatsApp API
func (c *Client) sendMessage(payload map[string]interface{}) (*MessageResponse, error) {
	endpoint := fmt.Sprintf("/%s/messages", c.phoneNumberID)

	c.logger.Debug("Sending message to WhatsApp",
		zap.String("endpoint", endpoint),
		zap.Any("payload", payload),
	)

	resp, err := c.httpClient.R().
		SetBody(payload).
		Post(endpoint)

	if err != nil {
		c.logger.Error("Failed to send message", zap.Error(err))
		return nil, errors.NewWhatsAppError(err)
	}

	if resp.IsError() {
		var errResp ErrorResponse
		if err := json.Unmarshal(resp.Body(), &errResp); err == nil {
			c.logger.Error("WhatsApp API error",
				zap.Int("code", errResp.Error.Code),
				zap.String("message", errResp.Error.Message),
				zap.String("type", errResp.Error.Type),
			)
			return nil, errors.NewWhatsAppError(fmt.Errorf("%s: %s", errResp.Error.Type, errResp.Error.Message))
		}

		c.logger.Error("WhatsApp API error",
			zap.Int("status", resp.StatusCode()),
			zap.String("body", string(resp.Body())),
		)
		return nil, errors.NewWhatsAppError(fmt.Errorf("WhatsApp API returned status %d", resp.StatusCode()))
	}

	var msgResp MessageResponse
	if err := json.Unmarshal(resp.Body(), &msgResp); err != nil {
		c.logger.Error("Failed to parse response", zap.Error(err))
		return nil, errors.NewInternalError(err)
	}

	c.logger.Info("Message sent successfully",
		zap.String("message_id", msgResp.Messages[0].ID),
	)

	return &msgResp, nil
}

// GetMessageStatus gets the delivery status of a message
func (c *Client) GetMessageStatus(messageID string) (*MessageStatus, error) {
	endpoint := fmt.Sprintf("/%s", messageID)

	resp, err := c.httpClient.R().Get(endpoint)
	if err != nil {
		return nil, errors.NewWhatsAppError(err)
	}

	if resp.IsError() {
		return nil, errors.NewWhatsAppError(fmt.Errorf("failed to get message status: %d", resp.StatusCode()))
	}

	var status MessageStatus
	if err := json.Unmarshal(resp.Body(), &status); err != nil {
		return nil, errors.NewInternalError(err)
	}

	return &status, nil
}

// SetTimeout sets the HTTP client timeout
func (c *Client) SetTimeout(duration time.Duration) {
	c.httpClient.SetTimeout(duration)
}

// SetRetryPolicy sets the retry policy for the HTTP client
func (c *Client) SetRetryPolicy(count int, waitTime time.Duration) {
	c.httpClient.SetRetryCount(count)
	c.httpClient.SetRetryWaitTime(waitTime)
}
