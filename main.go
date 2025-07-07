package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adityadeshlahre/multi-tenant-backend-app/database"
	"github.com/adityadeshlahre/multi-tenant-backend-app/handlers"
	"github.com/adityadeshlahre/multi-tenant-backend-app/pkg/middleware"
	"github.com/adityadeshlahre/multi-tenant-backend-app/repository"
)

func main() {
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrgRepository(db)
	articleRepo := repository.NewArticleRepository(db)

	authHandler := handlers.NewAuthHandler(userRepo, orgRepo)
	orgHandler := handlers.NewOrganizationHandler(orgRepo)
	articleHandler := handlers.NewArticleHandler(articleRepo)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", middleware.AuthMiddleware(userRepo), authHandler.GetProfile)
			auth.PUT("/profile", middleware.AuthMiddleware(userRepo), authHandler.UpdateProfile)
			auth.POST("/join-org/:orgId", middleware.AuthMiddleware(userRepo), authHandler.JoinOrganization)
		}

		orgs := api.Group("/organizations")
		{
			orgs.POST("/", orgHandler.CreateOrganization)
			orgs.GET("/", orgHandler.GetAllOrganizations)

			orgRoutes := orgs.Group("/:orgId")
			orgRoutes.Use(middleware.AuthMiddleware(userRepo))
			orgRoutes.Use(middleware.OrganizationContext(orgRepo))
			{
				orgRoutes.GET("/", orgHandler.GetOrganization)
				orgRoutes.PUT("/", orgHandler.UpdateOrganization)
				orgRoutes.DELETE("/", orgHandler.DeleteOrganization)

				orgRoutes.POST("/articles", articleHandler.CreateArticle)
				orgRoutes.GET("/articles", articleHandler.GetAllArticles)

				articleRoutes := orgRoutes.Group("/articles/:id")
				articleRoutes.Use(middleware.ArticleContext(articleRepo))
				{
					articleRoutes.GET("/", articleHandler.GetArticle)
					articleRoutes.PUT("/", articleHandler.UpdateArticle)
					articleRoutes.DELETE("/", articleHandler.DeleteArticle)

					articleRoutes.POST("/comments", articleHandler.CreateComment)
					articleRoutes.GET("/comments", articleHandler.GetComments)
					articleRoutes.PUT("/comments/:commentId", articleHandler.UpdateComment)
					articleRoutes.DELETE("/comments/:commentId", articleHandler.DeleteComment)
				}
			}
		}

		articles := api.Group("/articles")
		{
			articles.GET("/published", articleHandler.GetPublishedArticles)
			articles.GET("/my", middleware.AuthMiddleware(userRepo), articleHandler.GetMyArticles)
		}
	}

	fmt.Println("🚀 Server starting on http://localhost:8080 ...")
	fmt.Println("📄 API endpoints available at http://localhost:8080/api/v1")

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
