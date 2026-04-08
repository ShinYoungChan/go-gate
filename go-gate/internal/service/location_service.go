package service

import (
	"errors"
	"go-gate/internal/models"
	"go-gate/internal/repository"
	"math"
)

type LocationService struct {
	repo repository.LocationRepository
}

func NewLocationService(repo repository.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

func (s *LocationService) GetLocationList() ([]models.Location, error) {
	// 여기서 나중에 위도/경도 기반 거리 계산 등을 넣을 수 있습니다.
	return s.repo.GetAllLocations()
}

func (s *LocationService) CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // 지구 반지름 (미터 단위)

	// 1. 위도, 경도를 라디안(Radian) 단위로 변환
	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)

	rLat1 := lat1 * (math.Pi / 180)
	rLat2 := lat2 * (math.Pi / 180)

	// 2. 하버사인 공식 적용
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(rLat1)*math.Cos(rLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// 3. 최종 거리 산출(미터 단위)
	return earthRadius * c
}

func (s *LocationService) GetLocation(id uint) (*models.Location, error) {
	location, err := s.repo.GetLocationByID(id)
	if err != nil {
		return nil, err
	}

	if location == nil {
		return nil, errors.New("장소를 찾을 수 없습니다.")
	}
	return location, nil
}
