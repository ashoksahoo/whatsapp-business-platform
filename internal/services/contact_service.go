package services

import (
	"github.com/ashoksahoo/whatsapp-business-platform/internal/models"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/repositories"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
)

// ContactService handles contact business logic
type ContactService struct {
	contactRepo *repositories.ContactRepository
}

// NewContactService creates a new contact service
func NewContactService(contactRepo *repositories.ContactRepository) *ContactService {
	return &ContactService{
		contactRepo: contactRepo,
	}
}

// GetContact gets a contact by ID
func (s *ContactService) GetContact(contactID string) (*models.Contact, error) {
	var contact models.Contact
	if err := s.contactRepo.FindByID(contactID, &contact); err != nil {
		return nil, errors.NewNotFound("Contact", contactID)
	}
	return &contact, nil
}

// GetContactByPhone gets a contact by phone number
func (s *ContactService) GetContactByPhone(phone string) (*models.Contact, error) {
	contact, err := s.contactRepo.FindByPhone(phone)
	if err != nil {
		return nil, errors.NewNotFound("Contact", phone)
	}
	return contact, nil
}

// ListContacts lists all contacts with pagination and filters
func (s *ContactService) ListContacts(filters map[string]interface{}, pagination *utils.Pagination) ([]*models.Contact, error) {
	return s.contactRepo.ListWithFilters(filters, pagination)
}

// SearchContacts searches contacts by name or phone
func (s *ContactService) SearchContacts(query string, pagination *utils.Pagination) ([]*models.Contact, error) {
	return s.contactRepo.Search(query, pagination)
}

// UpdateContact updates contact information
func (s *ContactService) UpdateContact(contactID string, updates map[string]interface{}) (*models.Contact, error) {
	var contact models.Contact
	if err := s.contactRepo.FindByID(contactID, &contact); err != nil {
		return nil, errors.NewNotFound("Contact", contactID)
	}

	if err := s.contactRepo.UpdateFields(contactID, &contact, updates); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	// Fetch updated contact
	if err := s.contactRepo.FindByID(contactID, &contact); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	return &contact, nil
}

// GetOrCreateContact gets an existing contact or creates a new one
func (s *ContactService) GetOrCreateContact(phone string) (*models.Contact, error) {
	return s.contactRepo.GetOrCreate(phone)
}
