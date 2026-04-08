package repository

import (
	"errors"
	"go-gate/internal/models"

	"gorm.io/gorm"
)

type LocationRepository interface {
	GetLocationByID(id uint) (*models.Location, error)
	GetAllLocations() ([]models.Location, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) GetAllLocations() ([]models.Location, error) {
	var locations []models.Location
	err := r.db.Find(&locations).Error
	return locations, err
}

func (r *locationRepository) GetLocationByID(id uint) (*models.Location, error) {
	var location models.Location

	err := r.db.Where("id = ?", id).First(&location).Error
	if err != nil {
		// 에러가 데이터를 찾지 못한 건지 확인
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// 그 외는 실제 에러(DB 연결 끊김)등
		return nil, err
	}
	return &location, nil
}
