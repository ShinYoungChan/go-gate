package database

import (
	"fmt"
	"go-gate/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	fmt.Println("Init DB Start!")

	// 1. postgreSQL 연결
	dsn := "..."

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("DB 연결 실패: ", err)
	}

	if err = db.AutoMigrate(&models.User{}, &models.MembershipItem{}, &models.UserMembership{}, &models.Location{}, &models.AccessLog{}); err != nil {
		log.Fatal("DB 마이그레이션 실패: ", err)
	}

	return db
}
