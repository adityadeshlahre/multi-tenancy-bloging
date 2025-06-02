package repository

import (
	"context"
	"github.com/adityadeshlahre/multi-tenant-backend-app/model"
	"gorm.io/gorm"
	"log"
)

type articleRepository struct {
	db *gorm.DB
}

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article *model.Article) (*model.Article, error)
	GetArticleByID(ctx context.Context, id uint) (*model.Article, error)
	GetAllArticles(ctx context.Context) ([]model.Article, error)
	UpdateArticle(ctx context.Context, article *model.Article) (*model.Article, error)
	DeleteArticle(ctx context.Context, id uint) error
	GetArticlesByUserID(ctx context.Context, userID uint) ([]model.Article, error)
	CreateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	GetCommentsByArticleID(ctx context.Context, articleID uint) ([]model.Comment, error)
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) CreateArticle(ctx context.Context, article *model.Article) (*model.Article, error) {
	if err := r.db.WithContext(ctx).Create(article).Error; err != nil {
		log.Printf("Error creating article: %v", err)
		return nil, err
	}
	return article, nil
}

func (r *articleRepository) GetArticleByID(ctx context.Context, id uint) (*model.Article, error) {
	var article model.Article
	if err := r.db.WithContext(ctx).First(&article, id).Error; err != nil {
		log.Printf("Error fetching article by ID %d: %v", id, err)
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) GetAllArticles(ctx context.Context) ([]model.Article, error) {
	var articles []model.Article
	if err := r.db.WithContext(ctx).Find(&articles).Error; err != nil {
		log.Printf("Error fetching all articles: %v", err)
		return nil, err
	}
	return articles, nil
}

func (r *articleRepository) UpdateArticle(ctx context.Context, article *model.Article) (*model.Article, error) {
	if err := r.db.WithContext(ctx).Save(article).Error; err != nil {
		log.Printf("Error updating article ID %d: %v", article.ID, err)
		return nil, err
	}
	return article, nil
}

func (r *articleRepository) DeleteArticle(ctx context.Context, id uint) error {

	if err := r.db.WithContext(ctx).Delete(&model.Article{}, id).Error; err != nil {
		log.Printf("Error deleting article ID %d: %v", id, err)
	}
	return nil
}

func (r *articleRepository) GetArticlesByUserID(ctx context.Context, userID uint) ([]model.Article, error) {
	var articles []model.Article
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&articles).Error; err != nil {
		log.Printf("Error fetching articles by user ID %d: %v", userID, err)
		return nil, err
	}
	return articles, nil
}

func (r *articleRepository) CreateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
		log.Printf("Error creating comment: %v", err)
		return nil, err
	}
	return comment, nil
}

func (r *articleRepository) GetCommentsByArticleID(ctx context.Context, articleID uint) ([]model.Comment, error) {

	var comments []model.Comment
	if err := r.db.WithContext(ctx).Where("article_id = ?", articleID).Find(&comments).Error; err != nil {
		log.Printf("Error fetching comments by article ID %d: %v", articleID, err)
		return nil, err
	}
	return comments, nil
}


