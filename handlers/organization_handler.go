package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/adityadeshlahre/multi-tenant-backend-app/model"
	"github.com/adityadeshlahre/multi-tenant-backend-app/repository"
)

type OrganizationHandler struct {
	orgRepo repository.OrgRepository
}

func NewOrganizationHandler(orgRepo repository.OrgRepository) *OrganizationHandler {
	return &OrganizationHandler{
		orgRepo: orgRepo,
	}
}

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org := &model.Organization{
		Name: req.Name,
	}

	createdOrg, err := h.orgRepo.CreateOrganization(c.Request.Context(), org)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Organization created successfully",
		"organization": gin.H{"id": createdOrg.ID, "name": createdOrg.Name},
	})
}

func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	org, exists := c.Get("organization")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization not found in context"})
		return
	}

	orgModel := org.(*model.Organization)
	c.JSON(http.StatusOK, gin.H{
		"id":   orgModel.ID,
		"name": orgModel.Name,
	})
}

func (h *OrganizationHandler) GetAllOrganizations(c *gin.Context) {
	organizations, err := h.orgRepo.GetAllOrganizations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"organizations": organizations,
	})
}

func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	org, exists := c.Get("organization")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization not found in context"})
		return
	}

	orgModel := org.(*model.Organization)

	var updateData struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateData.Name != "" {
		orgModel.Name = updateData.Name
	}

	updatedOrg, err := h.orgRepo.UpdateOrganization(c.Request.Context(), orgModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Organization updated successfully",
		"organization": gin.H{"id": updatedOrg.ID, "name": updatedOrg.Name},
	})
}

func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	err = h.orgRepo.DeleteOrganization(c.Request.Context(), uint(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization deleted successfully",
	})
}
