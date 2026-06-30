package user

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, h *Handler) {
	r.POST("/signup", h.SignUp)
	r.POST("/login", h.Login)
	r.GET("/user/info/:id", h.GetUserInfo)
	r.GET("user/mypage/:id", h.GetUserSummary)
}
