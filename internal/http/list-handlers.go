package http

import (
	"net/http"
	"time"

	"github.com/charlieroth/reminders/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateListRequestBody struct {
	Name string `json:"name"`
}

type CreateListResponseData struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCreateListResponseData(list domain.List) CreateListResponseData {
	return CreateListResponseData{
		ID:        list.ID,
		Name:      list.Name,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
}

func CreateList(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req CreateListRequestBody
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l, err := app.listService.CreateList(gtx, domain.NewCreateListRequest(req.Name))
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusCreated, NewCreateListResponseData(l))
	}
}

type UpdateListRequestBody struct {
	Name string `json:"name"`
}

func NewUpdateListRequest(name string) domain.UpdateListRequest {
	return domain.UpdateListRequest{Name: name}
}

type UpdateListResponseData struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUpdateListResponseData(list domain.List) UpdateListResponseData {
	return UpdateListResponseData{
		ID:        list.ID,
		Name:      list.Name,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
}

func UpdateList(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var req UpdateListRequestBody
		if err := gtx.ShouldBindJSON(&req); err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		listID, err := uuid.Parse(gtx.Param("list_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l, err := app.listService.UpdateList(gtx, listID, domain.UpdateListRequest{Name: req.Name})
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewUpdateListResponseData(l))
	}
}

type GetListResponseData struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewGetListResponseData(list domain.List) GetListResponseData {
	return GetListResponseData{
		ID:        list.ID,
		Name:      list.Name,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
}

func GetList(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		listID, err := uuid.Parse(gtx.Param("list_id"))
		if err != nil {
			gtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		l, err := app.listService.GetList(gtx, listID)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewGetListResponseData(l))
	}
}

func NewGetListsResponseData(lists []domain.List) []GetListResponseData {
	responseData := []GetListResponseData{}
	for _, l := range lists {
		responseData = append(responseData, NewGetListResponseData(l))
	}
	return responseData
}

func GetLists(app *App) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		lists, err := app.listService.GetLists(gtx)
		if err != nil {
			gtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gtx.JSON(http.StatusOK, NewGetListsResponseData(lists))
	}
}
