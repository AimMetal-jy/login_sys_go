package handlers

import (
	"login_sys_go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserRepo *models.UserRepository
}

func NewAuthHandler(userRepo *models.UserRepository) *AuthHandler {
	return &AuthHandler{
		UserRepo: userRepo,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	User    *models.User `json:"user,omitempty"`
}

// Register 用户注册
func (ah *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// 创建用户
	user, err := ah.UserRepo.CreateUser(req.Username, req.Password)
	if err != nil {
		if err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, AuthResponse{
				Success: false,
				Message: "Username already exists",
			})
		} else {
			c.JSON(http.StatusInternalServerError, AuthResponse{
				Success: false,
				Message: "Failed to create user: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Success: true,
		Message: "User registered successfully",
		User:    user,
	})
}

// Login 用户登录
func (ah *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AuthResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// 验证用户
	user, err := ah.UserRepo.ValidateUser(req.Username, req.Password)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, AuthResponse{
				Success: false,
				Message: "Invalid username or password",
			})
		} else if err.Error() == "user account is not active" {
			c.JSON(http.StatusForbidden, AuthResponse{
				Success: false,
				Message: "User account is not active",
			})
		} else {
			c.JSON(http.StatusInternalServerError, AuthResponse{
				Success: false,
				Message: "Login failed: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Message: "Login successful",
		User:    user,
	})
}

// Health 健康检查
func (ah *AuthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "Login system is running",
	})
}