package routes

import (
	"go-gate/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupMembershipRoutes(r *gin.Engine, h *handler.UserMembershipHandler) {
	r.GET("/membership/info/:id", h.GetUserMembershipInfo)
}
