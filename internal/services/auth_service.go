package services

import (
	"time"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/models"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/repositories"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
)

// AuthService handles authentication business logic
type AuthService struct {
	apiKeyRepo *repositories.APIKeyRepository
}

// NewAuthService creates a new auth service
func NewAuthService(apiKeyRepo *repositories.APIKeyRepository) *AuthService {
	return &AuthService{
		apiKeyRepo: apiKeyRepo,
	}
}

// CreateAPIKey creates a new API key
func (s *AuthService) CreateAPIKey(name string, permissions []string, expiresAt *time.Time) (*models.APIKey, string, error) {
	// Generate API key
	rawKey, err := utils.GenerateAPIKey()
	if err != nil {
		return nil, "", errors.NewInternalError(err)
	}

	// Hash the key
	keyHash, err := utils.HashAPIKey(rawKey)
	if err != nil {
		return nil, "", errors.NewInternalError(err)
	}

	// Create API key model
	apiKey := &models.APIKey{
		Name:        name,
		KeyHash:     keyHash,
		KeyPrefix:   utils.GetKeyPrefix(rawKey),
		Permissions: permissions,
		ExpiresAt:   expiresAt,
	}

	if err := s.apiKeyRepo.Create(apiKey); err != nil {
		return nil, "", errors.NewDatabaseError(err)
	}

	// Return the API key with the raw key (only time it's visible)
	return apiKey, rawKey, nil
}

// ValidateAPIKey validates an API key and returns the key info
func (s *AuthService) ValidateAPIKey(rawKey string) (*models.APIKey, error) {
	// Get key prefix for optimization
	prefix := utils.GetKeyPrefix(rawKey)

	// Find keys with matching prefix
	apiKeys, err := s.apiKeyRepo.FindByKeyPrefix(prefix)
	if err != nil {
		return nil, errors.NewUnauthorized("Invalid API key")
	}

	// Check each key with matching prefix
	for _, apiKey := range apiKeys {
		if utils.CompareAPIKey(apiKey.KeyHash, rawKey) {
			// Check if expired
			if apiKey.IsExpired() {
				return nil, errors.NewAppError(errors.ErrAPIKeyExpired, "API key has expired", 401)
			}

			// Update last used
			s.apiKeyRepo.UpdateLastUsed(apiKey.ID)

			return apiKey, nil
		}
	}

	return nil, errors.NewUnauthorized("Invalid API key")
}

// RevokeAPIKey revokes an API key
func (s *AuthService) RevokeAPIKey(keyID string) error {
	var apiKey models.APIKey
	if err := s.apiKeyRepo.FindByID(keyID, &apiKey); err != nil {
		return errors.NewNotFound("API Key", keyID)
	}

	return s.apiKeyRepo.Delete(&apiKey)
}

// ListAPIKeys lists all API keys
func (s *AuthService) ListAPIKeys() ([]*models.APIKey, error) {
	return s.apiKeyRepo.ListAll()
}
