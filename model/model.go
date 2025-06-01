package model

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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

func ConnectDatabase() (*gorm.DB, error) {
	dns := "host=localhost user=postgres password=yourpassword dbname=multitenantapp port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&Organization{}, &User{}, &Article{}, &Comment{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
