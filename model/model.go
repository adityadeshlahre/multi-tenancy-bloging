package model

import (
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Name     string    `json:"name"`
	Users    []User    `gorm:"many2many:user_organization;"`
	Articles []Article `gorm:"foreignKey:OrganizationID"`
}

type User struct {
	gorm.Model
	Name         string         `json:"username"`
	Email        string         `json:"email"`
	Organization []Organization `gorm:"many2many:user_organization;"`
	Comments     []Comment      `gorm:"foreignKey:AuthorID"`
}

type Article struct {
	gorm.Model
	Title          string       `json:"title"`
	Content        string       `json:"content"`
	Status         string       `json:"status"`
	OrganizationID uint         `json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	UserID         uint         `json:"user_id"`
	User           User         `gorm:"foreignKey:UserID"`
	Comments       []Comment    `gorm:"foreignKey:ArticleID"`
}

type Comment struct {
	gorm.Model
	Message   string  `json:"message"`
	ArticleID uint    `json:"article_id"`
	Article   Article `gorm:"foreignKey:ArticleID"`
	AuthorID  uint    `json:"author_id"`
	Author    User    `gorm:"foreignKey:AuthorID"`
}
