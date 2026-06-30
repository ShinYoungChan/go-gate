package entry

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, h *Handler) {
	r.POST("/api/v1/entry/token", h.GetEntryToken)
	r.POST("/api/v1/entry/verify", h.PostEntry)
}
