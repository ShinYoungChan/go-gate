package location

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetLocations(c *gin.Context) {
	locations, err := h.service.GetLocationList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "조회 실패", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "조회 성공", "data": locations})
}
