package repositories

import (
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"gorm.io/gorm"
)

// BaseRepository provides common CRUD operations
type BaseRepository struct {
	DB *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{DB: db}
}

// Create creates a new record
func (r *BaseRepository) Create(model interface{}) error {
	return r.DB.Create(model).Error
}

// FindByID finds a record by ID
func (r *BaseRepository) FindByID(id string, model interface{}) error {
	return r.DB.Where("id = ?", id).First(model).Error
}

// Update updates a record
func (r *BaseRepository) Update(model interface{}) error {
	return r.DB.Save(model).Error
}

// UpdateFields updates specific fields
func (r *BaseRepository) UpdateFields(id string, model interface{}, updates map[string]interface{}) error {
	return r.DB.Model(model).Where("id = ?", id).Updates(updates).Error
}

// Delete soft deletes a record
func (r *BaseRepository) Delete(model interface{}) error {
	return r.DB.Delete(model).Error
}

// HardDelete permanently deletes a record
func (r *BaseRepository) HardDelete(model interface{}) error {
	return r.DB.Unscoped().Delete(model).Error
}

// List returns paginated records
func (r *BaseRepository) List(model interface{}, pagination *utils.Pagination) error {
	query := r.DB.Model(model)

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return err
	}
	pagination.SetTotal(total)

	// Apply pagination
	return pagination.ApplyToQuery(query).Find(model).Error
}

// Exists checks if a record exists
func (r *BaseRepository) Exists(query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := r.DB.Model(&struct{}{}).Where(query, args...).Count(&count).Error
	return count > 0, err
}
