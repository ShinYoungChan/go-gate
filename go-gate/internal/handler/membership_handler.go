package handler

import (
	"go-gate/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	service *service.MembershipService
}

func NewMembershipHandler(service *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{service: service}
}

func (h *MembershipHandler) GetUserMembershipInfo(c *gin.Context) {
	userIdStr := c.Param("user_id")
	locIdStr := c.Param("location_id")
	userId, _ := strconv.Atoi(userIdStr)
	locationId, _ := strconv.Atoi(locIdStr)

	userMembership, err := h.service.GetUserMembership(uint(userId), uint(locationId))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "내역 조회 중 오류 발생"})
		return
	} else if userMembership == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "조회 성공",
		"data":    userMembership,
	})
}

func (h *MembershipHandler) GetMembershipItems(c *gin.Context) {
	locIdStr := c.Param("location_id")
	locationId, err := strconv.Atoi(locIdStr)

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

	c.JSON(http.StatusOK, gin.H{
		"message": "조회 성공",
		"data":    membershipItems,
	})
}
