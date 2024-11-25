package http

import (
	"net/http"
	"time"

	"github.com/charlieroth/reminders/internal/jwt"
	"github.com/charlieroth/reminders/internal/session"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req LoginRequest
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := app.userService.GetUserByEmail(gtx.Request.Context(), req.Email)
		if err != nil {
			gtx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if user.PasswordHash != req.PasswordHash {
			gtx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		expUtc := time.Now().Add(time.Minute * 60).Unix()
		token, err := jwt.GenerateToken(app.config, user.ID, expUtc)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		session, err := app.authService.Login(gtx.Request.Context(), session.CreateSessionRequest{
			UserID:    user.ID,
			Token:     token,
			UserAgent: gtx.Request.UserAgent(),
		})
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, LoginResponse{Token: session.Token})
	}
}

type LogoutRequest struct {
	Email string `json:"email"`
}

func Logout(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req LogoutRequest
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := app.authService.LogoutByEmail(gtx.Request.Context(), req.Email)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, nil)
	}
}

func Refresh(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		userID := gtx.GetString("user_id")
		if userID == "" {
			gtx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			gtx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}

		expUtc := time.Now().Add(time.Minute * 60).Unix()
		token, err := jwt.GenerateToken(app.config, userUUID, expUtc)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		expiresAt, err := jwt.ExpiresAt(app.config, token)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		session, err := app.authService.Refresh(gtx.Request.Context(), session.RefreshSessionRequest{
			UserID:    userUUID,
			Token:     token,
			ExpiresAt: expiresAt,
			UserAgent: gtx.Request.UserAgent(),
		})
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, session)
	}
}

func GetSessions(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		gtx.JSON(http.StatusOK, nil)
	}
}
