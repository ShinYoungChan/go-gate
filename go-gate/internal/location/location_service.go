package location

import (
	"errors"
	"math"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetLocationList() ([]Location, error) {
	return s.repo.GetAllLocations()
}

func (s *Service) CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000

	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)
	rLat1 := lat1 * (math.Pi / 180)
	rLat2 := lat2 * (math.Pi / 180)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(rLat1)*math.Cos(rLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (s *Service) GetLocation(id uint) (*Location, error) {
	location, err := s.repo.GetLocationByID(id)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, errors.New("장소를 찾을 수 없습니다.")
	}
	return location, nil
}
