package main

import (
	"fmt"
	"go-gate/internal/database"
	"go-gate/internal/entry"
	"go-gate/internal/location"
	"go-gate/internal/membership"
	"go-gate/internal/payment"
	"go-gate/internal/user"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Main Start!")

	if err := godotenv.Load("go.env"); err != nil {
		log.Fatal(".env 파일을 찾을 수 없습니다.")
	}

	db := database.InitDB()

	// Repository
	locRepo := location.NewRepository(db)
	userRepo := user.NewUserRepository(db)
	membershipRepo := membership.NewRepository(db)
	accessLogRepo := entry.NewAccessLogRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)

	// Service
	locService := location.NewService(locRepo)
	membershipService := membership.NewService(membershipRepo)
	entryService := entry.NewService(membershipService, accessLogRepo, locService)
	userService := user.NewService(userRepo, accessLogRepo, membershipService)
	paymentService := payment.NewService(paymentRepo, membershipRepo)

	// Handler
	userHandler := user.NewHandler(userService)
	membershipHandler := membership.NewHandler(membershipService)
	entryHandler := entry.NewHandler(entryService)
	paymentHandler := payment.NewHandler(paymentService)
	locHandler := location.NewHandler(locService)

	r := gin.Default()
	r.Use(cors.Default())
	user.SetupRoutes(r, userHandler)
	membership.SetupRoutes(r, membershipHandler)
	entry.SetupRoutes(r, entryHandler)
	payment.SetupRoutes(r, paymentHandler)
	location.SetupRoutes(r, locHandler)

	r.Run(":8080")
}
