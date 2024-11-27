package http

import (
	"net/http"
	"time"

	"github.com/charlieroth/reminders/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateUserRequestBody struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type CreateUserResponseData struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCreateUserResponseData(user domain.User) CreateUserResponseData {
	return CreateUserResponseData{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func CreateUser(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req CreateUserRequestBody
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, err := app.userService.CreateUser(gtx, domain.CreateUserRequest{Email: req.Email, PasswordHash: req.PasswordHash})
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusCreated, NewCreateUserResponseData(u))
	}
}

type GetUserResponseData struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetUserResponseData(user domain.User) GetUserResponseData {
	return GetUserResponseData{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func GetUser(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		userId, err := uuid.Parse(gtx.Param("user_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, err := app.userService.GetUser(gtx, userId)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewGetUserResponseData(u))
	}
}

func NewGetUsersResponseData(users []domain.User) []GetUserResponseData {
	responseData := []GetUserResponseData{}
	for _, u := range users {
		responseData = append(responseData, NewGetUserResponseData(u))
	}
	return responseData
}

func GetUsers(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		users, err := app.userService.GetUsers(gtx)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewGetUsersResponseData(users))
	}
}

type UpdateUserRequestBody struct {
	Email        *string `json:"email"`
	PasswordHash *string `json:"password_hash"`
}

func UpdateUser(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req UpdateUserRequestBody
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userId, err := uuid.Parse(gtx.Param("user_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, err := app.userService.UpdateUser(gtx, userId, domain.UpdateUserRequest{Email: req.Email, PasswordHash: req.PasswordHash})
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewGetUserResponseData(u))
	}
}
