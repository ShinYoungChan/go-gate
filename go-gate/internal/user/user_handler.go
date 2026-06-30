package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type signUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	JoinDate string `json:"joindate"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SignUp(c *gin.Context) {
	var req signUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.SignUpUser(req.Name, req.Email, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "회원가입 실패", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "회원가입 성공!"})
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.AuthenticateUser(req.Email, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "로그인 실패", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "로그인 성공!"})
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	user, err := h.service.GetUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "내역 조회 중 오류 발생"})
		return
	} else if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "조회 성공",
		"data": userResponse{
			Name:     user.Name,
			Email:    user.Email,
			JoinDate: user.CreatedAt.String(),
		},
	})
}

func (h *Handler) GetUserSummary(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	summary, err := h.service.GetUserSummary(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "내역 조회 중 오류 발생"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "조회 성공", "data": summary})
}
