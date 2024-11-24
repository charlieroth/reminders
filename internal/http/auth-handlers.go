package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		gtx.JSON(http.StatusNotImplemented, nil)
	}
}

func Logout(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		gtx.JSON(http.StatusNotImplemented, nil)
	}
}

func Refresh(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		gtx.JSON(http.StatusNotImplemented, nil)
	}
}

func GetSessions(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		gtx.JSON(http.StatusNotImplemented, nil)
	}
}
