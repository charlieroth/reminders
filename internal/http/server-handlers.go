package http

import (
	"net/http"
	"os"
	"runtime"

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

type ReadinessCheckResponse struct {
	Status string `json:"status"`
}

func ReadinessCheck(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		if err := app.databaseService.StatusCheck(gtx.Request.Context()); err != nil {
			gtx.JSON(http.StatusServiceUnavailable, ReadinessCheckResponse{
				Status: "unavailable",
			})
			return
		}

		gtx.JSON(http.StatusOK, ReadinessCheckResponse{
			Status: "ok",
		})
	}
}

type LivenessCheckResponse struct {
	Status     string `json:"status"`
	Host       string `json:"host"`
	GOMAXPROCS int    `json:"gomaxprocs"`
}

func LivenessCheck(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		host, err := os.Hostname()
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, LivenessCheckResponse{
			Status:     "up",
			Host:       host,
			GOMAXPROCS: runtime.GOMAXPROCS(0),
		})
	}
}
