package service

import "go-gate/internal/repository"

type LocationService struct {
	repo *repository.LocationRepository
}

func (s *LocationService) CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // 지구 반지름 (미터 단위)

	// 각도를 라디안으로 변환하는 수식이 필요
	// 여기에 하버사인(Haversine) 공식을 구현

	return 0.0 // 임시 리턴
}
