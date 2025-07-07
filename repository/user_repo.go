package repository

import (
	"context"
	"log"

	"github.com/adityadeshlahre/multi-tenant-backend-app/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id uint) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUsersByOrganization(ctx context.Context, orgID uint) ([]model.User, error)
	AddUserToOrganization(ctx context.Context, userID, orgID uint) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Preload("Organizations").First(&user, id).Error; err != nil {
		log.Printf("Error fetching user by ID %d: %v", id, err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Preload("Organizations").Where("email = ?", email).First(&user).Error; err != nil {
		log.Printf("Error fetching user by email %s: %v", email, err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		log.Printf("Error updating user ID %d: %v", user.ID, err)
		return nil, err
	}
	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		log.Printf("Error deleting user ID %d: %v", id, err)
		return err
	}
	return nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).Preload("Organizations").Find(&users).Error; err != nil {
		log.Printf("Error fetching all users: %v", err)
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUsersByOrganization(ctx context.Context, orgID uint) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).Joins("JOIN user_organizations ON users.id = user_organizations.user_id").
		Where("user_organizations.organization_id = ?", orgID).Find(&users).Error; err != nil {
		log.Printf("Error fetching users by organization ID %d: %v", orgID, err)
		return nil, err
	}
	return users, nil
}

func (r *userRepository) AddUserToOrganization(ctx context.Context, userID, orgID uint) error {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	var org model.Organization
	if err := r.db.WithContext(ctx).First(&org, orgID).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&user).Association("Organizations").Append(&org)
}
