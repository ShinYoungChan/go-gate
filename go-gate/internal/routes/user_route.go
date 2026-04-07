package routes

import (
	"go-gate/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, h *handler.UserHandler) {
	r.POST("/signup", h.SignUp)
	r.POST("/login", h.Login)
	r.GET("/user/info/:id", h.GetUserInfo)
}
