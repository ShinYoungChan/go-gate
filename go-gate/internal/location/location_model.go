package location

type Location struct {
	ID        uint    `gorm:"primaryKey"`
	PlaceName string  `gorm:"size:20;not null;uniqueIndex"`
	Category  string  `gorm:"size:20;not null"`
	Lat       float64 `gorm:"type:decimal(10,8);not null"`
	Lon       float64 `gorm:"type:decimal(11,8);not null"`
	Address   string  `gorm:"size:100;not null"`
	ImageURL  string  `gorm:"size:255"`
}
