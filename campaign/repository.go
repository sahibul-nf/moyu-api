package campaign

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindByUserID(userID int) ([]Campaign, error)
	FindAll() ([]Campaign, error)
	Save(campaign Campaign) (Campaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Preload("CampaignImages", "campaign_images.isPrimary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) FindByUserID(userID int) ([]Campaign, error) {
	var campaign []Campaign

	err := r.db.Where("id = ?", userID).Preload("CampaignImages", "campaign_images.isPrimary = 1").Find(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
