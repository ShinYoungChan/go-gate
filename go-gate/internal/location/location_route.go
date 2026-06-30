package location

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, h *Handler) {
	r.GET("/api/v1/locations", h.GetLocations)
}
