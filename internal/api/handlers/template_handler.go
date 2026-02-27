package handlers

import (
	"strconv"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/models"
	"github.com/ashoksahoo/whatsapp-business-platform/internal/services"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// TemplateHandler handles template-related requests
type TemplateHandler struct {
	templateService *services.TemplateService
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler(templateService *services.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
	}
}

// CreateTemplate handles POST /api/v1/templates
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var template models.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		utils.ErrorJSON(c, errors.NewBadRequest("Invalid request body: "+err.Error()))
		return
	}

	if err := h.templateService.CreateTemplate(&template); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorJSON(c, appErr)
		} else {
			utils.ErrorJSON(c, errors.NewInternalError(err))
		}
		return
	}

	utils.CreatedJSON(c, template)
}

// GetTemplate handles GET /api/v1/templates/:id
func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	templateID := c.Param("id")

	template, err := h.templateService.GetTemplate(templateID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorJSON(c, appErr)
		} else {
			utils.ErrorJSON(c, errors.NewInternalError(err))
		}
		return
	}

	utils.SuccessJSON(c, 200, template)
}

// ListTemplates handles GET /api/v1/templates
func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	pagination := utils.NewPagination(limit, offset)

	templates, err := h.templateService.ListTemplates(pagination)
	if err != nil {
		utils.ErrorJSON(c, errors.NewInternalError(err))
		return
	}

	utils.ListJSON(c, templates, pagination)
}

// UpdateTemplate handles PATCH /api/v1/templates/:id
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	templateID := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorJSON(c, errors.NewBadRequest("Invalid request body: "+err.Error()))
		return
	}

	template, err := h.templateService.UpdateTemplate(templateID, updates)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorJSON(c, appErr)
		} else {
			utils.ErrorJSON(c, errors.NewInternalError(err))
		}
		return
	}

	utils.SuccessJSON(c, 200, template)
}

// DeleteTemplate handles DELETE /api/v1/templates/:id
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("id")

	if err := h.templateService.DeleteTemplate(templateID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorJSON(c, appErr)
		} else {
			utils.ErrorJSON(c, errors.NewInternalError(err))
		}
		return
	}

	utils.NoContentJSON(c)
}
