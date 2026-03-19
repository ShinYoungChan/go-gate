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

	// 2. Service
	locService := service.NewLocationService(locRepo)
	userService := service.NewUserService(userRepo, userMembershipRepo, accessLogRepo, locService)

	// 3. Handler
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()
	routes.SetupUserRoutes(r, userHandler)

	r.Run(":8080")
}
