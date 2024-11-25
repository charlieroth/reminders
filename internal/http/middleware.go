package http

import (
	"net/http"

	"github.com/charlieroth/reminders/internal/jwt"
	"github.com/charlieroth/reminders/internal/session"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(app *App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := authHeader[7:]
		if err := jwt.VerifyToken(app.config, token); err != nil {
			app.authService.Logout(ctx, session.InvalidateSessionRequest{Token: token})

			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := jwt.UserID(app.config, token)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user_id", userID.String())
		ctx.Set("token", token)
		ctx.Next()
	}
}
