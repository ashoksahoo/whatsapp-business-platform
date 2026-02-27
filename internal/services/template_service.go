package services

import (
	"github.com/ashoksahoo/whatsapp-business-platform/internal/models"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/repositories"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
)

// TemplateService handles template business logic
type TemplateService struct {
	templateRepo *repositories.TemplateRepository
}

// NewTemplateService creates a new template service
func NewTemplateService(templateRepo *repositories.TemplateRepository) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
	}
}

// CreateTemplate creates a new template
func (s *TemplateService) CreateTemplate(template *models.Template) error {
	if err := template.Validate(); err != nil {
		return errors.NewBadRequest(err.Error())
	}

	return s.templateRepo.Create(template)
}

// GetTemplate gets a template by ID
func (s *TemplateService) GetTemplate(templateID string) (*models.Template, error) {
	var template models.Template
	if err := s.templateRepo.FindByID(templateID, &template); err != nil {
		return nil, errors.NewNotFound("Template", templateID)
	}
	return &template, nil
}

// GetTemplateByName gets a template by name and language
func (s *TemplateService) GetTemplateByName(name, language string) (*models.Template, error) {
	template, err := s.templateRepo.FindByName(name, language)
	if err != nil {
		return nil, errors.NewNotFound("Template", name)
	}
	return template, nil
}

// ListTemplates lists all templates
func (s *TemplateService) ListTemplates(pagination *utils.Pagination) ([]*models.Template, error) {
	return s.templateRepo.ListAll(pagination)
}

// UpdateTemplate updates a template
func (s *TemplateService) UpdateTemplate(templateID string, updates map[string]interface{}) (*models.Template, error) {
	var template models.Template
	if err := s.templateRepo.FindByID(templateID, &template); err != nil {
		return nil, errors.NewNotFound("Template", templateID)
	}

	if err := s.templateRepo.UpdateFields(templateID, &template, updates); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	// Fetch updated template
	if err := s.templateRepo.FindByID(templateID, &template); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	return &template, nil
}

// DeleteTemplate deletes a template
func (s *TemplateService) DeleteTemplate(templateID string) error {
	var template models.Template
	if err := s.templateRepo.FindByID(templateID, &template); err != nil {
		return errors.NewNotFound("Template", templateID)
	}

	return s.templateRepo.Delete(&template)
}
