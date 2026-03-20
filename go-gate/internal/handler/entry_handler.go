package handler

import (
	"go-gate/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EntryRequest struct {
	UserID     uint    `json:"user_id" binding:"required"`
	LocationID uint    `json:"location_id" binding:"required"`
	Lat        float64 `json:"lat" binding:"required"`
	Lon        float64 `json:"lon" binding:"required"`
}

type EntryHandler struct {
	service *service.EntryService
}

func NewEntryHandler(service *service.EntryService) *EntryHandler {
	return &EntryHandler{service: service}
}

func (h *EntryHandler) PostEntry(c *gin.Context) {
	var req EntryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userMembership, err := h.service.VerifyEntry(req.UserID, req.Lat, req.Lon, req.LocationID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "실패", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "입장 성공!",
		"remaining_count": userMembership.Count,
		"end_dt":          userMembership.EndDt.Format("2006-01-02"),
	})
}
