package routes

import (
	"go-gate/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupMembershipRoutes(r *gin.Engine, h *handler.MembershipHandler) {
	r.GET("/membership/info/:user_id/:location_id", h.GetUserMembershipInfo)
	r.GET("/location/membership/:location_id", h.GetMembershipItems)
}
