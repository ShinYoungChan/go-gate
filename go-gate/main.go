package main

import (
	"fmt"
	"go-gate/internal/database"
	"go-gate/internal/handler"
	"go-gate/internal/repository"
	"go-gate/internal/routes"
	"go-gate/internal/service"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Main Start!")

	// 1. .env 파일 로드
	if err := godotenv.Load("go.env"); err != nil {
		log.Fatal(".env 파일을 찾을 수 없습니다.")
	}

	db := database.InitDB()

	// 의존성 주입
	// 1. DB
	locRepo := repository.NewLocationRepository(db)
	userRepo := repository.NewUserRepository(db)
	userMembershipRepo := repository.NewUserMembershipRepository(db)
	accessLogRepo := repository.NewAccessLogRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// 2. Service
	locService := service.NewLocationService(locRepo)
	membershipService := service.NewUserMembershipService(userMembershipRepo)
	accessLogService := service.NewAccessLogService(accessLogRepo)
	userService := service.NewUserService(userRepo, accessLogService, membershipService)
	entryService := service.NewEntryService(membershipService, accessLogRepo, locService)
	paymentService := service.NewPaymentService(paymentRepo, userMembershipRepo)

	// 3. Handler
	userHandler := handler.NewUserHandler(userService)
	userMembershipHandler := handler.NewUserMembershipHandler(membershipService)
	entryHandler := handler.NewEntryHandler(entryService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	locHandler := handler.NewLocationHandler(locService)

	r := gin.Default()
	r.Use(cors.Default())
	routes.SetupUserRoutes(r, userHandler)
	routes.SetupMembershipRoutes(r, userMembershipHandler)
	routes.SetupEntryRoutes(r, entryHandler)
	routes.SetupPaymentRoutes(r, paymentHandler)
	routes.SetupLocationRoutes(r, locHandler)

	r.Run(":8080")
}
