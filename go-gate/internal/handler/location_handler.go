package handler

import (
	"go-gate/internal/models"
	"go-gate/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	service *service.LocationService
}

func NewLocationHandler(service *service.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

func (h *LocationHandler) GetLocations(c *gin.Context) {
	var locations []models.Location

	locations, err := h.service.GetLocationList()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "조회 실패", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "조회 성공",
		"data":    locations,
	})
}
