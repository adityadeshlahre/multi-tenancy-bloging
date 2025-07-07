package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/adityadeshlahre/multi-tenant-backend-app/model"
	"github.com/adityadeshlahre/multi-tenant-backend-app/pkg/middleware"
	"github.com/adityadeshlahre/multi-tenant-backend-app/repository"
)

type AuthHandler struct {
	userRepo repository.UserRepository
	orgRepo  repository.OrgRepository
}

func NewAuthHandler(userRepo repository.UserRepository, orgRepo repository.OrgRepository) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		orgRepo:  orgRepo,
	}
}

type RegisterRequest struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	OrganizationID uint   `json:"organization_id"`
	Role           string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = "member"
	}

	hashedPassword := middleware.HashPassword(req.Password)

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	createdUser, err := h.userRepo.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	if req.OrganizationID != 0 {
		org, err := h.orgRepo.GetOrganizationByID(c.Request.Context(), req.OrganizationID)
		if err == nil {
			// Use GORM's association mode to properly add the organization
			if err := h.userRepo.AddUserToOrganization(c.Request.Context(), createdUser.ID, org.ID); err != nil {
				log.Printf("Failed to add user to organization: %v", err)
			}
		}
	}

	token, err := middleware.GenerateJWT(createdUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    gin.H{"id": createdUser.ID, "name": createdUser.Name, "email": createdUser.Email, "role": createdUser.Role},
		"token":   token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !middleware.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := middleware.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role},
		"token":   token,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(*model.User)
	c.JSON(http.StatusOK, gin.H{
		"id":            userModel.ID,
		"name":          userModel.Name,
		"email":         userModel.Email,
		"role":          userModel.Role,
		"organizations": userModel.Organizations,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(*model.User)

	var updateData struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateData.Name != "" {
		userModel.Name = updateData.Name
	}

	updatedUser, err := h.userRepo.UpdateUser(c.Request.Context(), userModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    gin.H{"id": updatedUser.ID, "name": updatedUser.Name, "email": updatedUser.Email, "role": updatedUser.Role},
	})
}

func (h *AuthHandler) JoinOrganization(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	orgIDStr := c.Param("orgId")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	org, err := h.orgRepo.GetOrganizationByID(c.Request.Context(), uint(orgID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	userModel := user.(*model.User)

	err = h.userRepo.AddUserToOrganization(c.Request.Context(), userModel.ID, org.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Successfully joined organization",
		"organization": gin.H{"id": org.ID, "name": org.Name},
	})
}
