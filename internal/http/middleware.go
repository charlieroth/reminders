package http

import (
	"fmt"
	"net/http"

	"github.com/charlieroth/reminders/internal/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(app *App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userClaims, err := verifyClaimsFromAuthHeader(ctx, app)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("claims", userClaims)
		ctx.Next()
	}
}

func AdminMiddleware(app *App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userClaims, err := verifyClaimsFromAuthHeader(ctx, app)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !userClaims.IsAdmin {
			ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("user is not an admin"))
			return
		}

		ctx.Set("claims", userClaims)
		ctx.Next()
	}
}

func verifyClaimsFromAuthHeader(ctx *gin.Context, app *App) (*jwt.UserClaims, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return nil, fmt.Errorf("invalid authorization header")
	}

	token := authHeader[7:]
	userClaims, err := jwt.VerifyToken(app.config, token)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return nil, err
	}

	return userClaims, nil
}
