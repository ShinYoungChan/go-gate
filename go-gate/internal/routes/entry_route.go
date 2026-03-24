package routes

import (
	"go-gate/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupEntryRoutes(r *gin.Engine, h *handler.EntryHandler) {
	r.GET("/api/v1/entry/token", h.GetEntryToken)
	r.POST("/api/v1/entry/verify", h.PostEntry)
}
