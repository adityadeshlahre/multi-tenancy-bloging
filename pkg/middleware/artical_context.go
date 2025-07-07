package middleware

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/adityadeshlahre/multi-tenant-backend-app/model"
	"github.com/adityadeshlahre/multi-tenant-backend-app/repository"
)

const (
	ArticleKey        = "article"
	UserPermissionKey = "userPermission"
)

const (
	PermissionNone    = "none"
	PermissionView    = "view"
	PermissionComment = "comment"
	PermissionEdit    = "edit"
	PermissionOwner   = "owner"
)

func ArticleContext(articleRepo repository.ArticleRepository, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		articleIDParam := c.Param("articleId")
		if articleIDParam == "" {
			articleIDParam = c.Param("id")
		}

		if articleIDParam == "" {
			log.Println("Article ID not found in URL parameters")
			c.AbortWithStatusJSON(400, gin.H{"error": "Article ID is required"})
			return
		}

		articleID, err := strconv.ParseUint(articleIDParam, 10, 32)
		if err != nil {
			log.Printf("Invalid article ID format: %v", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "Invalid article ID format"})
			return
		}

		article, err := articleRepo.GetArticleByID(c.Request.Context(), uint(articleID))
		if err != nil {
			log.Printf("Error fetching article: %v", err)
			c.AbortWithStatusJSON(404, gin.H{"error": "Article not found"})
			return
		}

		organization, exists := c.Get("organization")
		if !exists {
			log.Println("Organization not found in context")
			c.AbortWithStatusJSON(500, gin.H{"error": "Organization context not found"})
			return
		}

		org := organization.(*model.Organization)

		if article.OrganizationID != org.ID {
			log.Printf("Article %d does not belong to organization %d", article.ID, org.ID)
			c.AbortWithStatusJSON(403, gin.H{"error": "Access denied: Article not in your organization"})
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			log.Println("User ID not found in context")
			c.AbortWithStatusJSON(401, gin.H{"error": "User authentication required"})
			return
		}

		userRole, exists := c.Get("userRole")
		if !exists {
			log.Println("User role not found in context")
			c.AbortWithStatusJSON(401, gin.H{"error": "User authentication required"})
			return
		}

		permission := determineUserPermissionWithRoles(article, userID.(uint), userRole.(string))

		c.Set(ArticleKey, article)
		c.Set(UserPermissionKey, permission)

		log.Printf("User %d has %s permission for article %d", userID.(uint), permission, article.ID)

		c.Next()
	}
}

func determineUserPermissionWithRoles(article *model.Article, userID uint, userRole string) string {
	if article.UserID == userID {
		return PermissionOwner
	}

	switch userRole {
	case "admin":
		return PermissionEdit
	case "editor":
		if article.Status == "published" {
			return PermissionEdit
		}
		return PermissionComment
	case "member":
		if article.Status == "published" {
			return PermissionComment
		}
		return PermissionView
	case "viewer":
		if article.Status == "published" {
			return PermissionView
		}
		return PermissionNone
	default:
		return PermissionView
	}
}

func CanViewArticle(c *gin.Context) bool {
	permission, exists := c.Get(UserPermissionKey)
	if !exists {
		return false
	}

	perm := permission.(string)
	return perm == PermissionView || perm == PermissionComment ||
		perm == PermissionEdit || perm == PermissionOwner
}

func CanCommentOnArticle(c *gin.Context) bool {
	permission, exists := c.Get(UserPermissionKey)
	if !exists {
		return false
	}

	perm := permission.(string)
	return perm == PermissionComment || perm == PermissionEdit || perm == PermissionOwner
}

func CanEditArticle(c *gin.Context) bool {
	permission, exists := c.Get(UserPermissionKey)
	if !exists {
		return false
	}

	perm := permission.(string)
	return perm == PermissionEdit || perm == PermissionOwner
}

func IsArticleOwner(c *gin.Context) bool {
	permission, exists := c.Get(UserPermissionKey)
	if !exists {
		return false
	}

	return permission.(string) == PermissionOwner
}

func GetArticleFromContext(c *gin.Context) (*model.Article, bool) {
	article, exists := c.Get(ArticleKey)
	if !exists {
		return nil, false
	}
	return article.(*model.Article), true
}

func GetUserPermissionFromContext(c *gin.Context) (string, bool) {
	permission, exists := c.Get(UserPermissionKey)
	if !exists {
		return "", false
	}
	return permission.(string), true
}
