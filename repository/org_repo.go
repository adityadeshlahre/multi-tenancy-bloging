package repository

import (
	"context"
	"log"

	"github.com/adityadeshlahre/multi-tenant-backend-app/model"
	"gorm.io/gorm"
)

type orgRepository struct {
	db *gorm.DB
}

type OrgRepository interface {
	CreateOrganization(ctx context.Context, org *model.Organization) (*model.Organization, error)
	GetOrganizationByID(ctx context.Context, id uint) (*model.Organization, error)
	GetOrganizationByName(ctx context.Context, name string) (*model.Organization, error)
	UpdateOrganization(ctx context.Context, org *model.Organization) (*model.Organization, error)
	DeleteOrganization(ctx context.Context, id uint) error
	GetAllOrganizations(ctx context.Context) ([]model.Organization, error)
}

func NewOrgRepository(db *gorm.DB) OrgRepository {
	return &orgRepository{db: db}
}

func (r *orgRepository) CreateOrganization(ctx context.Context, org *model.Organization) (*model.Organization, error) {
	if err := r.db.WithContext(ctx).Create(org).Error; err != nil {
		log.Printf("Error creating organization: %v", err)
		return nil, err
	}
	return org, nil
}

func (r *orgRepository) GetOrganizationByID(ctx context.Context, id uint) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.WithContext(ctx).First(&org, id).Error; err != nil {
		log.Printf("Error fetching organization by ID %d: %v", id, err)
		return nil, err
	}
	return &org, nil
}

func (r *orgRepository) GetOrganizationByName(ctx context.Context, name string) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&org).Error; err != nil {
		log.Printf("Error fetching organization by name %s: %v", name, err)
		return nil, err
	}
	return &org, nil
}

func (r *orgRepository) UpdateOrganization(ctx context.Context, org *model.Organization) (*model.Organization, error) {
	if err := r.db.WithContext(ctx).Save(org).Error; err != nil {
		log.Printf("Error updating organization ID %d: %v", org.ID, err)
		return nil, err
	}
	return org, nil
}

func (r *orgRepository) DeleteOrganization(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Organization{}, id).Error; err != nil {
		log.Printf("Error deleting organization ID %d: %v", id, err)
	}
	return nil
}

func (r *orgRepository) GetAllOrganizations(ctx context.Context) ([]model.Organization, error) {
	var orgs []model.Organization
	if err := r.db.WithContext(ctx).Find(&orgs).Error; err != nil {
		log.Printf("Error fetching all organizations: %v", err)
		return nil, err
	}
	return orgs, nil
}
