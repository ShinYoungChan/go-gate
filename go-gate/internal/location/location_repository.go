package location

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	GetLocationByID(id uint) (*Location, error)
	GetAllLocations() ([]Location, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &locationRepository{db: db}
}

func (r *locationRepository) GetAllLocations() ([]Location, error) {
	var locations []Location
	err := r.db.Find(&locations).Error
	return locations, err
}

func (r *locationRepository) GetLocationByID(id uint) (*Location, error) {
	var location Location
	err := r.db.Where("id = ?", id).First(&location).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &location, nil
}
