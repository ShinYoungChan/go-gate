package models

type Location struct {
	ID        uint    `gorm:"primaryKey"`
	PlaceName string  `gorm:"size:20;not null;uniqueIndex"`
	Category  string  `gorm:"size:20;not null"` // 추가: 장소 종류 (예: 헬스장, 요가)
	Lat       float64 `gorm:"type:decimal(10,8);not null"`
	Lon       float64 `gorm:"type:decimal(11,8);not null"`
	Address   string  `gorm:"size:100;not null"`
	ImageURL  string  `gorm:"size:255"` // 추가: 장소 이미지 경로
}
