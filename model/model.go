package model

import (
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Name     string    `json:"name" gorm:"uniqueIndex"`
	Users    []User    `gorm:"many2many:user_organizations;"`
	Articles []Article `gorm:"foreignKey:OrganizationID"`
}

type User struct {
	gorm.Model
	Name          string         `json:"name"`
	Email         string         `json:"email" gorm:"uniqueIndex"`
	Password      string         `json:"-"`
	Role          string         `json:"role" gorm:"default:'member'"`
	Organizations []Organization `gorm:"many2many:user_organizations;"`
	Articles      []Article      `gorm:"foreignKey:UserID"`
	Comments      []Comment      `gorm:"foreignKey:AuthorID"`
}

type Article struct {
	gorm.Model
	Title          string       `json:"title"`
	Content        string       `json:"content"`
	Status         string       `json:"status" gorm:"default:'draft'"`
	OrganizationID uint         `json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	UserID         uint         `json:"user_id"`
	User           User         `gorm:"foreignKey:UserID"`
	Comments       []Comment    `gorm:"foreignKey:ArticleID"`
}

type Comment struct {
	gorm.Model
	Content   string  `json:"content"`
	ArticleID uint    `json:"article_id"`
	Article   Article `gorm:"foreignKey:ArticleID"`
	AuthorID  uint    `json:"author_id"`
	Author    User    `gorm:"foreignKey:AuthorID"`
}
