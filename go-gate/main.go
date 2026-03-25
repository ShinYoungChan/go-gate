package main

import (
	"fmt"
	"go-gate/internal/database"
	"go-gate/internal/handler"
	"go-gate/internal/repository"
	"go-gate/internal/routes"
	"go-gate/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Main Start!")

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
	userService := service.NewUserService(userRepo)
	entryService := service.NewEntryService(membershipService, accessLogRepo, locService)
	paymentService := service.NewPaymentService(paymentRepo, userMembershipRepo)

	// 3. Handler
	userHandler := handler.NewUserHandler(userService)
	entryHandler := handler.NewEntryHandler(entryService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	r := gin.Default()
	routes.SetupUserRoutes(r, userHandler)
	routes.SetupEntryRoutes(r, entryHandler)

	r.Run(":8080")
}
