package routes

import (
	"go-gate/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupEntryRoutes(r *gin.Engine, h *handler.EntryHandler) {
	r.POST("/api/v1/users/entry", h.PostEntry)
}
