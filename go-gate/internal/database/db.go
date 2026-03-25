package database

import (
	"fmt"
	"go-gate/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	fmt.Println("Init DB Start!")

	// 1. 환경 변수 값 읽기
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// 2. postgreSQL 연결
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("DB 연결 실패: ", err)
	}

	if err = db.AutoMigrate(&models.User{}, &models.MembershipItem{}, &models.UserMembership{}, &models.Location{}, &models.AccessLog{}, &models.PaymentLog{}); err != nil {
		log.Fatal("DB 마이그레이션 실패: ", err)
	}

	fmt.Println("Init DB END......")
	return db
}
