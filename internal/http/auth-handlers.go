package http

import (
	"net/http"
	"time"

	"github.com/charlieroth/reminders/internal/domain"
	"github.com/charlieroth/reminders/internal/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func Register(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req RegisterRequest
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := app.userService.CreateUser(gtx.Request.Context(), domain.CreateUserRequest{
			Email:        req.Email,
			PasswordHash: req.PasswordHash,
		})
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusCreated, nil)
	}
}

type LoginRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type LoginResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
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

		accessToken, accessClaims, err := jwt.CreateToken(app.config, user.ID, user.Email, false, app.config.AccessTokenDuration)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		refreshToken, refreshClaims, err := jwt.CreateToken(app.config, user.ID, user.Email, false, app.config.RefreshTokenDuration)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		session, err := app.authService.Login(gtx.Request.Context(), domain.CreateSessionRequest{
			ID:           uuid.MustParse(refreshClaims.RegisteredClaims.ID),
			Email:        user.Email,
			RefreshToken: refreshToken,
			IsRevoked:    false,
			ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
		})
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, LoginResponse{
			SessionID:             session.ID,
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		})
	}
}

func Logout(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		sessionID, err := uuid.Parse(gtx.Param("session_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
			return
		}

		err = app.authService.Logout(gtx.Request.Context(), domain.DeleteSessionRequest{
			ID: sessionID,
		})
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, nil)
	}
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func Refresh(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req RefreshRequest
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		refreshClaims, err := jwt.VerifyToken(app.config, req.RefreshToken)
		if err != nil {
			gtx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		session, err := app.authService.GetSession(gtx.Request.Context(), domain.GetSessionRequest{
			ID: uuid.MustParse(refreshClaims.RegisteredClaims.ID),
		})
		if err != nil {
			gtx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		if session.IsRevoked {
			gtx.JSON(http.StatusUnauthorized, gin.H{"error": "session revoked"})
			return
		}

		if session.Email != refreshClaims.Email {
			gtx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			return
		}

		accessToken, accessClaims, err := jwt.CreateToken(
			app.config,
			refreshClaims.ID,
			refreshClaims.Email,
			refreshClaims.IsAdmin,
			app.config.AccessTokenDuration,
		)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, RefreshResponse{
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
		})
	}
}

func RevokeSession(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		sessionID, err := uuid.Parse(gtx.Param("session_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
			return
		}

		err = app.authService.RevokeSession(gtx.Request.Context(), domain.RevokeSessionRequest{
			ID: sessionID,
		})
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, nil)
	}
}
