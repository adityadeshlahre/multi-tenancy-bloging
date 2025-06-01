package database

import (
	"github.com/adityadeshlahre/multi-tenant-backend-app/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectDatabase() (*gorm.DB, error) {
	dns := "host=localhost user=postgres password=yourpassword dbname=multitenantapp port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&model.Organization{}, &model.User{}, &model.Article{}, &model.Comment{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
