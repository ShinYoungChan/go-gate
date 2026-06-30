package entry

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type entryRequest struct {
	Token string  `json:"token" binding:"required"`
	Lat   float64 `json:"lat" binding:"required"`
	Lon   float64 `json:"lon" binding:"required"`
}

type tokenRequest struct {
	UserID     uint `json:"user_id" binding:"required"`
	LocationID uint `json:"location_id" binding:"required"`
}

type tokenResponse struct {
	Token   string `json:"token"`
	Expires int    `json:"expires"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) PostEntry(c *gin.Context) {
	var req entryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userMembership, err := h.service.VerifyEntry(req.Token, req.Lat, req.Lon)
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

func (h *Handler) GetEntryToken(c *gin.Context) {
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	token, err := h.service.GenerateEntryToken(req.UserID, req.LocationID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "토큰 생성 성공",
		"data":    tokenResponse{Token: token, Expires: 30},
	})
}
