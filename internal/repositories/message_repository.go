package repositories

import (
	"time"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/models"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"gorm.io/gorm"
)

// MessageRepository handles message data access
type MessageRepository struct {
	*BaseRepository
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindByPhone finds messages by phone number with pagination
func (r *MessageRepository) FindByPhone(phone string, pagination *utils.Pagination) ([]*models.Message, error) {
	var messages []*models.Message

	query := r.DB.Where("from_number = ? OR to_number = ?", phone, phone).
		Order("timestamp DESC")

	// Get total count
	var total int64
	if err := query.Model(&models.Message{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.SetTotal(total)

	// Apply pagination and fetch
	err := pagination.ApplyToQuery(query).Find(&messages).Error
	return messages, err
}

// FindByDateRange finds messages within a date range
func (r *MessageRepository) FindByDateRange(start, end time.Time, pagination *utils.Pagination) ([]*models.Message, error) {
	var messages []*models.Message

	query := r.DB.Where("timestamp >= ? AND timestamp <= ?", start, end).
		Order("timestamp DESC")

	// Get total count
	var total int64
	if err := query.Model(&models.Message{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.SetTotal(total)

	// Apply pagination
	err := pagination.ApplyToQuery(query).Find(&messages).Error
	return messages, err
}

// FindByStatus finds messages by status
func (r *MessageRepository) FindByStatus(status string, pagination *utils.Pagination) ([]*models.Message, error) {
	var messages []*models.Message

	query := r.DB.Where("status = ?", status).
		Order("created_at DESC")

	// Get total count
	var total int64
	if err := query.Model(&models.Message{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.SetTotal(total)

	// Apply pagination
	err := pagination.ApplyToQuery(query).Find(&messages).Error
	return messages, err
}

// FindByWhatsAppMessageID finds a message by WhatsApp message ID
func (r *MessageRepository) FindByWhatsAppMessageID(waMessageID string) (*models.Message, error) {
	var message models.Message
	err := r.DB.Where("whatsapp_message_id = ?", waMessageID).First(&message).Error
	return &message, err
}

// Search performs full-text search on message content
func (r *MessageRepository) Search(query string, filters map[string]interface{}, pagination *utils.Pagination) ([]*models.Message, error) {
	var messages []*models.Message

	dbQuery := r.DB.Where("content ILIKE ?", "%"+query+"%")

	// Apply filters
	if phone, ok := filters["phone"].(string); ok && phone != "" {
		dbQuery = dbQuery.Where("from_number = ? OR to_number = ?", phone, phone)
	}
	if direction, ok := filters["direction"].(string); ok && direction != "" {
		dbQuery = dbQuery.Where("direction = ?", direction)
	}
	if msgType, ok := filters["type"].(string); ok && msgType != "" {
		dbQuery = dbQuery.Where("message_type = ?", msgType)
	}

	dbQuery = dbQuery.Order("timestamp DESC")

	// Get total count
	var total int64
	if err := dbQuery.Model(&models.Message{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.SetTotal(total)

	// Apply pagination
	err := pagination.ApplyToQuery(dbQuery).Find(&messages).Error
	return messages, err
}

// CountByPhone counts messages for a phone number
func (r *MessageRepository) CountByPhone(phone string) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Message{}).
		Where("from_number = ? OR to_number = ?", phone, phone).
		Count(&count).Error
	return count, err
}

// ListWithFilters lists messages with various filters
func (r *MessageRepository) ListWithFilters(filters map[string]interface{}, pagination *utils.Pagination) ([]*models.Message, error) {
	var messages []*models.Message

	query := r.DB.Model(&models.Message{})

	// Apply filters
	if phone, ok := filters["phone"].(string); ok && phone != "" {
		query = query.Where("from_number = ? OR to_number = ?", phone, phone)
	}
	if direction, ok := filters["direction"].(string); ok && direction != "" {
		query = query.Where("direction = ?", direction)
	}
	if msgType, ok := filters["type"].(string); ok && msgType != "" {
		query = query.Where("message_type = ?", msgType)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate, ok := filters["start_date"].(time.Time); ok && !startDate.IsZero() {
		query = query.Where("timestamp >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(time.Time); ok && !endDate.IsZero() {
		query = query.Where("timestamp <= ?", endDate)
	}

	query = query.Order("timestamp DESC")

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.SetTotal(total)

	// Apply pagination
	err := pagination.ApplyToQuery(query).Find(&messages).Error
	return messages, err
}

// UpdateStatus updates the status of a message
func (r *MessageRepository) UpdateStatus(whatsappMessageID, status string) error {
	return r.DB.Model(&models.Message{}).
		Where("whatsapp_message_id = ?", whatsappMessageID).
		Update("status", status).Error
}
