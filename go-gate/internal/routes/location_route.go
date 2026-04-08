package routes

import (
	"go-gate/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupLocationRoutes(r *gin.Engine, h *handler.LocationHandler) {
	r.GET("/api/v1/locations", h.GetLocations)
}
