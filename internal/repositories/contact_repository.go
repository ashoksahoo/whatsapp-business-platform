package repositories

import (
	"time"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/models"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ContactRepository handles contact data access
type ContactRepository struct {
	*BaseRepository
}

// NewContactRepository creates a new contact repository
func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindByPhone finds a contact by phone number
func (r *ContactRepository) FindByPhone(phone string) (*models.Contact, error) {
	var contact models.Contact
	err := r.DB.Where("phone_number = ?", phone).First(&contact).Error
	return &contact, err
}

// GetOrCreate gets an existing contact or creates a new one
func (r *ContactRepository) GetOrCreate(phone string) (*models.Contact, error) {
	var contact models.Contact

	// Use upsert to handle race conditions
	result := r.DB.Where("phone_number = ?", phone).
		Attrs(models.Contact{PhoneNumber: phone}).
		FirstOrCreate(&contact)

	return &contact, result.Error
}

// Search searches contacts by name or phone
func (r *ContactRepository) Search(query string, pagination *utils.Pagination) ([]*models.Contact, error) {
	var contacts []*models.Contact

	dbQuery := r.DB.Where("name ILIKE ? OR phone_number ILIKE ?", "%"+query+"%", "%"+query+"%").
		Order("last_message_at DESC NULLS LAST")

	// Get total count
	var total int64
	if err := dbQuery.Model(&models.Contact{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.SetTotal(total)

	// Apply pagination
	err := pagination.ApplyToQuery(dbQuery).Find(&contacts).Error
	return contacts, err
}

// UpdateLastMessage updates the last message timestamp for a contact
func (r *ContactRepository) UpdateLastMessage(phone string, timestamp time.Time) error {
	return r.DB.Model(&models.Contact{}).
		Where("phone_number = ?", phone).
		Updates(map[string]interface{}{
			"last_message_at": timestamp,
			"updated_at":      time.Now().UTC(),
		}).Error
}

// IncrementMessageCount increments the message count for a contact
func (r *ContactRepository) IncrementMessageCount(phone string, delta int) error {
	return r.DB.Model(&models.Contact{}).
		Where("phone_number = ?", phone).
		UpdateColumn("message_count", gorm.Expr("message_count + ?", delta)).Error
}

// UpdateUnreadCount updates the unread count for a contact
func (r *ContactRepository) UpdateUnreadCount(phone string, delta int) error {
	return r.DB.Model(&models.Contact{}).
		Where("phone_number = ?", phone).
		UpdateColumn("unread_count", gorm.Expr("GREATEST(unread_count + ?, 0)", delta)).Error
}

// ResetUnreadCount resets the unread count to zero
func (r *ContactRepository) ResetUnreadCount(phone string) error {
	return r.DB.Model(&models.Contact{}).
		Where("phone_number = ?", phone).
		Update("unread_count", 0).Error
}

// FindActive finds contacts with recent activity
func (r *ContactRepository) FindActive(limit int, pagination *utils.Pagination) ([]*models.Contact, error) {
	var contacts []*models.Contact

	query := r.DB.Where("last_message_at IS NOT NULL").
		Order("last_message_at DESC").
		Limit(limit)

	err := query.Find(&contacts).Error
	return contacts, err
}

// ListWithFilters lists contacts with filters
func (r *ContactRepository) ListWithFilters(filters map[string]interface{}, pagination *utils.Pagination) ([]*models.Contact, error) {
	var contacts []*models.Contact

	query := r.DB.Model(&models.Contact{})

	// Apply sorting
	sortField := "last_message_at"
	sortOrder := "DESC"
	if sf, ok := filters["sort"].(string); ok && sf != "" {
		sortField = sf
	}
	if so, ok := filters["order"].(string); ok && so != "" {
		sortOrder = so
	}

	query = query.Order(sortField + " " + sortOrder + " NULLS LAST")

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.SetTotal(total)

	// Apply pagination
	err := pagination.ApplyToQuery(query).Find(&contacts).Error
	return contacts, err
}

// UpsertContact creates or updates a contact
func (r *ContactRepository) UpsertContact(contact *models.Contact) error {
	return r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "phone_number"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "profile_url", "updated_at"}),
	}).Create(contact).Error
}
