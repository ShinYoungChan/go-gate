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
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()
	routes.SetupUserRoutes(r, userHandler)

	r.Run(":8080")
}
