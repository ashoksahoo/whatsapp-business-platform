package repositories

import (
	"github.com/ashoksahoo/whatsapp-business-platform/internal/models"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"gorm.io/gorm"
)

// TemplateRepository handles template data access
type TemplateRepository struct {
	*BaseRepository
}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindByName finds a template by name and language
func (r *TemplateRepository) FindByName(name, language string) (*models.Template, error) {
	var template models.Template
	err := r.DB.Where("name = ? AND language = ?", name, language).First(&template).Error
	return &template, err
}

// FindByCategory finds templates by category
func (r *TemplateRepository) FindByCategory(category string, pagination *utils.Pagination) ([]*models.Template, error) {
	var templates []*models.Template

	query := r.DB.Where("category = ?", category).Order("created_at DESC")

	var total int64
	if err := query.Model(&models.Template{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.SetTotal(total)

	err := pagination.ApplyToQuery(query).Find(&templates).Error
	return templates, err
}

// FindByStatus finds templates by status
func (r *TemplateRepository) FindByStatus(status string, pagination *utils.Pagination) ([]*models.Template, error) {
	var templates []*models.Template

	query := r.DB.Where("status = ?", status).Order("created_at DESC")

	var total int64
	if err := query.Model(&models.Template{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.SetTotal(total)

	err := pagination.ApplyToQuery(query).Find(&templates).Error
	return templates, err
}

// FindApproved finds all approved templates
func (r *TemplateRepository) FindApproved(pagination *utils.Pagination) ([]*models.Template, error) {
	return r.FindByStatus(models.TemplateStatusApproved, pagination)
}

// ListAll lists all templates with pagination
func (r *TemplateRepository) ListAll(pagination *utils.Pagination) ([]*models.Template, error) {
	var templates []*models.Template

	query := r.DB.Model(&models.Template{}).Order("created_at DESC")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.SetTotal(total)

	err := pagination.ApplyToQuery(query).Find(&templates).Error
	return templates, err
}
