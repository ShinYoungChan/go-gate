package routes

import (
	"go-gate/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupMembershipRoutes(r *gin.Engine, h *handler.UserMembershipHandler) {
	r.GET("/membership/info/:user_id/:location_id", h.GetUserMembershipInfo)
}
