package membership

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUserMembershipInfo(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("user_id"))
	locationId, _ := strconv.Atoi(c.Param("location_id"))

	userMembership, err := h.service.GetUserMembership(uint(userId), uint(locationId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "내역 조회 중 오류 발생"})
		return
	} else if userMembership == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "조회 성공", "data": userMembership})
}

func (h *Handler) GetMembershipItems(c *gin.Context) {
	locationId, err := strconv.Atoi(c.Param("location_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 지점 ID입니다."})
		return
	}

	membershipItems, err := h.service.GetAvailableMemberships(uint(locationId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "내역 조회 중 오류 발생"})
		return
	}

	log.Println(membershipItems)
	c.JSON(http.StatusOK, gin.H{"message": "조회 성공", "data": membershipItems})
}
