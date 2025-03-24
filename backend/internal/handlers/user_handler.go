package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"qckim-backend/config"
	"qckim-backend/internal/logger"
	"qckim-backend/internal/middleware"
	"qckim-backend/internal/models"
	"qckim-backend/internal/repository"
	"qckim-backend/internal/services"
	"qckim-backend/utils/cryptPass"
	"qckim-backend/utils/email"
	"qckim-backend/utils/jwt"
	"qckim-backend/utils/valid"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler - Обработчик запросов, связанных с пользователями

var cfg = config.GetConfig()

type UserHandler struct {
	UserService    *services.UserRepo
	TokenService   *services.TokenRepo
	QckLinkService *services.QckLinkRepo
}

// NewUserHandler - Создает новый обработчик для пользователя
func NewUserHandler(db *repository.QckRepo) *UserHandler {
	return &UserHandler{
		UserService:    services.NewUser(db.GetDB()),
		TokenService:   services.NewToken(db.GetDB()),
		QckLinkService: services.NewQckLink(db.GetDB()),
	}
}

// RegisterRoutes - Регистрирует маршруты для пользователя
func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/users/login", h.Login)
	r.POST("/api/users/register", h.Register)
	r.POST("/api/users/logout", h.Logout)
	r.GET("/api/users/activate", h.Activate)
	r.POST("/api/users/request-password-reset", h.RequestPasswordReset)
	r.POST("/api/users/reset-password", h.ResetPassword)
	r.POST("/api/users/refresh-token", h.RefreshToken)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	println("Рефреш токен", req.RefreshToken)

	userID, err := h.UserService.GetUserIDByRefreshToken(req.RefreshToken)
	if err != nil {
		logger.Error("h.UserService.GetUserIDByRefreshToken(req.RefreshToken)", userID, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	accessToken, err := jwt.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	newRefreshToken, err := jwt.GenerateToken32()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	err = h.TokenService.SaveToken(userID, c.ClientIP(), c.Request.UserAgent(), newRefreshToken)
	if err != nil {
		logger.Error("RefreshToken", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

func (h *UserHandler) RequestPasswordReset(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.UserService.CreateResetToken(req.Email)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		logger.Error("PasswordReset", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	email.SendEmailAsync(
		req.Email,
		"Подтверждение почты",
		fmt.Sprintf("Для изменения по ссылке %s/new-password?token=%s", cfg.BaseURL, token),
	)

}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	userID, expiresAt, err := h.TokenService.CheckResetToken(req.Token)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		logger.Error("CheckResetToken", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	if time.Now().After(expiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token expired"})
		return
	}

	hashedPassword, err := cryptPass.HashPassword(req.Password)
	if err != nil {
		logger.Error("hashedPassword", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	err = h.UserService.ChangePassword(userID, hashedPassword)

	if err != nil {
		logger.Error("h.UserService.ChangePassword", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update password"})
		return
	}
	// добавить удаление временного токена на сброс
	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// TODO проксировать через nginx?
// Activate - Активация учетной записи пользователя
func (h *UserHandler) Activate(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	urlInvalidToken := fmt.Sprintf("%s/activate?status=invalid_token", cfg.BaseURL)
	if token == "" {
		c.Redirect(http.StatusSeeOther, urlInvalidToken)
		return
	}

	userID64, err := h.UserService.GetUserIDByRefreshToken(token)
	if err != nil {
		logger.Error("Invalid activation token", "error", err)
		c.Redirect(http.StatusSeeOther, urlInvalidToken)
		return
	}

	urlActivationFailed := fmt.Sprintf("%s/activate?status=activation_failed", cfg.BaseURL)
	err = h.UserService.ActivateUser(userID64)
	if err != nil {
		logger.Error("Failed to activate user", "error", err)
		c.Redirect(http.StatusSeeOther, urlActivationFailed)
		return
	}

	urlActivationSucsess := fmt.Sprintf("%s/activate?status=success", cfg.BaseURL)
	c.Redirect(http.StatusSeeOther, urlActivationSucsess)
}

// Login - Логин и генерации JWT токена
func (h *UserHandler) Login(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) != 0 {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := jwt.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     "Token is valid, already logged in",
			"profile_url": "/profile",
		})
		return
	}

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON input"})
		return
	}

	validate := validator.New()

	if err := validate.Var(req.Password, "required,min=8,max=64,ascii"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid password"})
		return
	}

	userID, isActive, err := h.UserService.CheckUser(req.Login, req.Password)
	if err != nil {
		logger.Error("h.UserService.CheckUser", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверные данные"})
		return
	}

	if !isActive {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Подтвердите почту перейдя по ссылке из письма"})
		return
	}

	refreshToken, err := jwt.GenerateToken32()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate refresh token"})
		return
	}

	err = h.TokenService.SaveToken(userID, c.ClientIP(), c.Request.UserAgent(), refreshToken)
	if err != nil {
		logger.Error("h.TokenService.GetValidToken", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get valid token"})
		return
	}

	accessToken, err := jwt.GenerateJWT(userID)
	if err != nil {
		logger.Error("h.TokenService.GetValidToken", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get valid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Register - Сохранение пользователя в бд и возвращение JWT токена
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegistrarionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	validate := validator.New()

	if err := validate.Var(req.Email, "required,email"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email"})
		return
	}

	if err := validate.Var(req.Password, "required,min=8,max=64,ascii"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid password"})
		return
	}

	if err := validate.Var(req.QckLinkPassword, "required,min=8,max=64,ascii"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid password"})
		return
	}

	if !valid.ValidateQckLink(req.QckLink) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid qck link"})
		return
	}

	id, err := h.QckLinkService.GetQckLinkID(req.QckLink, req.QckLinkPassword)
	if err != nil {
		logger.Debug("Error qck_link_id", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid qck link or qck_link_password"})
		return
	}

	userAgent := c.Request.UserAgent()
	userIP := c.RemoteIP()

	userID, err := h.UserService.SaveUser(
		req.Login,
		req.Email,
		req.Password,
		id,
	)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateLogin):
			c.JSON(http.StatusConflict, gin.H{"message": "Логин уже существует"})
		case errors.Is(err, models.ErrDuplicateEmail):
			c.JSON(http.StatusConflict, gin.H{"message": "Email уже существует"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User creation failed"})
		}
		return
	}

	refreshToken, err := jwt.GenerateToken32()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate refresh token"})
		return
	}

	err = h.TokenService.SaveToken(userID, userIP, userAgent, refreshToken)
	if err != nil {
		logger.Error("Save jwt", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
	}

	subject := email.EncodeSubject("Подтверждение почты")
	email.SendEmailAsync(
		req.Email,
		subject,
		fmt.Sprintf("Для подтверждения почты перейдите по ссылке %s/api/users/activate?token=%s", cfg.BaseURL, refreshToken),
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered, check you email",
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	userID, exists := c.Get(string(middleware.UserKey))
	println(userID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userID64, err := strconv.ParseInt(userID.(string), 10, 64)
	if err != nil {
		logger.Error("Invalid user ID type", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
	}

	err = h.TokenService.InvalidateUserTokens(userID64)
	if err != nil {
		logger.Error("Error invalidating tokens", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}
