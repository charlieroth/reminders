package http

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ServeOpenAPISpec(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		specBytes, err := os.ReadFile("docs/reminders.openapi.json")
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read OpenAPI specification"})
			return
		}

		gtx.Header("Content-Type", "application/json")
		gtx.String(http.StatusOK, string(specBytes))
	}
}

func ReadinessCheck(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		gtx.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func LivenessCheck(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		gtx.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
