package repositories

import (
	"time"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/models"
	"gorm.io/gorm"
)

// APIKeyRepository handles API key data access
type APIKeyRepository struct {
	*BaseRepository
}

// NewAPIKeyRepository creates a new API key repository
func NewAPIKeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindByKeyHash finds an API key by its hash
func (r *APIKeyRepository) FindByKeyHash(keyHash string) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := r.DB.Where("key_hash = ?", keyHash).First(&apiKey).Error
	return &apiKey, err
}

// FindByKeyPrefix finds API keys by prefix
func (r *APIKeyRepository) FindByKeyPrefix(prefix string) ([]*models.APIKey, error) {
	var apiKeys []*models.APIKey
	err := r.DB.Where("key_prefix = ?", prefix).Find(&apiKeys).Error
	return apiKeys, err
}

// UpdateLastUsed updates the last_used_at timestamp
func (r *APIKeyRepository) UpdateLastUsed(id string) error {
	return r.DB.Model(&models.APIKey{}).
		Where("id = ?", id).
		Update("last_used_at", time.Now().UTC()).Error
}

// FindValid finds all valid (not expired) API keys
func (r *APIKeyRepository) FindValid() ([]*models.APIKey, error) {
	var apiKeys []*models.APIKey
	err := r.DB.Where("expires_at IS NULL OR expires_at > ?", time.Now().UTC()).
		Find(&apiKeys).Error
	return apiKeys, err
}

// FindExpired finds all expired API keys
func (r *APIKeyRepository) FindExpired() ([]*models.APIKey, error) {
	var apiKeys []*models.APIKey
	err := r.DB.Where("expires_at IS NOT NULL AND expires_at <= ?", time.Now().UTC()).
		Find(&apiKeys).Error
	return apiKeys, err
}

// ListAll lists all API keys
func (r *APIKeyRepository) ListAll() ([]*models.APIKey, error) {
	var apiKeys []*models.APIKey
	err := r.DB.Order("created_at DESC").Find(&apiKeys).Error
	return apiKeys, err
}
