package database

import (
	"fmt"
	"log"

	"github.com/adityadeshlahre/multi-tenant-backend-app/config"
	"github.com/adityadeshlahre/multi-tenant-backend-app/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&model.Organization{}, &model.User{}, &model.Article{}, &model.Comment{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
		return nil, err
	}

	return db, nil
}
