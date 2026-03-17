package handler

import (
	"go-gate/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SignUpRequest 회원가입 요청 DTO
type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var req SignUpRequest

	// JSON 파일 읽기
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

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest

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
