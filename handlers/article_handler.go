package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/adityadeshlahre/multi-tenant-backend-app/model"
	"github.com/adityadeshlahre/multi-tenant-backend-app/pkg/middleware"
	"github.com/adityadeshlahre/multi-tenant-backend-app/repository"
)

type ArticleHandler struct {
	articleRepo repository.ArticleRepository
}

func NewArticleHandler(articleRepo repository.ArticleRepository) *ArticleHandler {
	return &ArticleHandler{
		articleRepo: articleRepo,
	}
}

type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Status  string `json:"status"`
}

type UpdateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	org, exists := c.Get("organization")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization not found in context"})
		return
	}

	orgModel := org.(*model.Organization)

	if req.Status == "" {
		req.Status = "draft"
	}

	article := &model.Article{
		Title:          req.Title,
		Content:        req.Content,
		Status:         req.Status,
		UserID:         userID.(uint),
		OrganizationID: orgModel.ID,
	}

	createdArticle, err := h.articleRepo.CreateArticle(c.Request.Context(), article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Article created successfully",
		"article": createdArticle,
	})
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	if !middleware.CanViewArticle(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this article"})
		return
	}

	article, exists := middleware.GetArticleFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Article not found in context"})
		return
	}

	permission, _ := middleware.GetUserPermissionFromContext(c)

	c.JSON(http.StatusOK, gin.H{
		"article":    article,
		"permission": permission,
	})
}

func (h *ArticleHandler) GetAllArticles(c *gin.Context) {
	org, exists := c.Get("organization")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization not found in context"})
		return
	}

	orgModel := org.(*model.Organization)
	articles, err := h.articleRepo.GetArticlesByOrganization(c.Request.Context(), orgModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
	})
}

func (h *ArticleHandler) GetPublishedArticles(c *gin.Context) {
	articles, err := h.articleRepo.GetPublishedArticles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch published articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
	})
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	if !middleware.CanEditArticle(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to edit this article"})
		return
	}

	article, exists := middleware.GetArticleFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Article not found in context"})
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Status != "" {
		article.Status = req.Status
	}

	updatedArticle, err := h.articleRepo.UpdateArticle(c.Request.Context(), article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article updated successfully",
		"article": updatedArticle,
	})
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	if !middleware.IsArticleOwner(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the article owner can delete this article"})
		return
	}

	article, exists := middleware.GetArticleFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Article not found in context"})
		return
	}

	err := h.articleRepo.DeleteArticle(c.Request.Context(), article.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article deleted successfully",
	})
}

func (h *ArticleHandler) GetMyArticles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	articles, err := h.articleRepo.GetArticlesByUserID(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
	})
}

func (h *ArticleHandler) CreateComment(c *gin.Context) {
	if !middleware.CanCommentOnArticle(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to comment on this article"})
		return
	}

	article, exists := middleware.GetArticleFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Article not found in context"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &model.Comment{
		Content:   req.Content,
		ArticleID: article.ID,
		AuthorID:  userID.(uint),
	}

	createdComment, err := h.articleRepo.CreateComment(c.Request.Context(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": createdComment,
	})
}

func (h *ArticleHandler) GetComments(c *gin.Context) {
	if !middleware.CanViewArticle(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view comments on this article"})
		return
	}

	article, exists := middleware.GetArticleFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Article not found in context"})
		return
	}

	comments, err := h.articleRepo.GetCommentsByArticleID(c.Request.Context(), article.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

func (h *ArticleHandler) UpdateComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &model.Comment{
		Content:  req.Content,
		AuthorID: userID.(uint),
	}
	comment.ID = uint(commentID)

	updatedComment, err := h.articleRepo.UpdateComment(c.Request.Context(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"comment": updatedComment,
	})
}

func (h *ArticleHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	if !middleware.IsArticleOwner(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the article owner can delete comments"})
		return
	}

	err = h.articleRepo.DeleteComment(c.Request.Context(), uint(commentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}
