package membership

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, h *Handler) {
	r.GET("/membership/info/:user_id/:location_id", h.GetUserMembershipInfo)
	r.GET("/location/membership/:location_id", h.GetMembershipItems)
}
