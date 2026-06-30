package database

import (
	"fmt"
	"go-gate/internal/entry"
	"go-gate/internal/location"
	"go-gate/internal/membership"
	"go-gate/internal/payment"
	"go-gate/internal/user"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	fmt.Println("Init DB Start!")

	host := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		host, dbUser, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB 연결 실패: ", err)
	}

	if err = db.AutoMigrate(
		&user.User{},
		&membership.MembershipItem{},
		&membership.UserMembership{},
		&location.Location{},
		&entry.AccessLog{},
		&payment.PaymentLog{},
	); err != nil {
		log.Fatal("DB 마이그레이션 실패: ", err)
	}

	fmt.Println("Init DB END......")
	return db
}
