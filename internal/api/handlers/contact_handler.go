package handlers

import (
	"strconv"

	"github.com/ashoksahoo/whatsapp-business-platform/internal/services"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/ashoksahoo/whatsapp-business-platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ContactHandler handles contact-related requests
type ContactHandler struct {
	contactService *services.ContactService
}

// NewContactHandler creates a new contact handler
func NewContactHandler(contactService *services.ContactService) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
	}
}

// GetContact handles GET /api/v1/contacts/:id
func (h *ContactHandler) GetContact(c *gin.Context) {
	contactID := c.Param("id")

	contact, err := h.contactService.GetContact(contactID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorJSON(c, appErr)
		} else {
			utils.ErrorJSON(c, errors.NewInternalError(err))
		}
		return
	}

	utils.SuccessJSON(c, 200, contact)
}

// ListContacts handles GET /api/v1/contacts
func (h *ContactHandler) ListContacts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	pagination := utils.NewPagination(limit, offset)

	filters := make(map[string]interface{})
	if sort := c.Query("sort"); sort != "" {
		filters["sort"] = sort
	}
	if order := c.Query("order"); order != "" {
		filters["order"] = order
	}

	contacts, err := h.contactService.ListContacts(filters, pagination)
	if err != nil {
		utils.ErrorJSON(c, errors.NewInternalError(err))
		return
	}

	utils.ListJSON(c, contacts, pagination)
}

// UpdateContact handles PATCH /api/v1/contacts/:id
func (h *ContactHandler) UpdateContact(c *gin.Context) {
	contactID := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorJSON(c, errors.NewBadRequest("Invalid request body: "+err.Error()))
		return
	}

	contact, err := h.contactService.UpdateContact(contactID, updates)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorJSON(c, appErr)
		} else {
			utils.ErrorJSON(c, errors.NewInternalError(err))
		}
		return
	}

	utils.SuccessJSON(c, 200, contact)
}

// SearchContacts handles GET /api/v1/contacts/search
func (h *ContactHandler) SearchContacts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorJSON(c, errors.NewBadRequest("Search query 'q' is required"))
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	pagination := utils.NewPagination(limit, offset)

	contacts, err := h.contactService.SearchContacts(query, pagination)
	if err != nil {
		utils.ErrorJSON(c, errors.NewInternalError(err))
		return
	}

	utils.ListJSON(c, contacts, pagination)
}
