package handler

import (
	"go-gate/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserMembershipHandler struct {
	service *service.UserMembershipService
}

func NewUserMembershipHandler(service *service.UserMembershipService) *UserMembershipHandler {
	return &UserMembershipHandler{service: service}
}

func (h *UserMembershipHandler) GetUserMembershipInfo(c *gin.Context) {
	userIdStr := c.Param("id")
	userID, _ := strconv.Atoi(userIdStr)

	userMembership, err := h.service.GetUserMembership(uint(userID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "내역 조회 중 오류 발생"})
	} else if userMembership == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "조회 성공",
		"data":    userMembership,
	})
}
