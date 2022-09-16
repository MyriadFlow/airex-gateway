package app

import (
	"collection/dto"
	"collection/service"

	// "encoding/json"
	"net/http"

	// "github.com/gorilla/mux"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CollectionHandler struct {
	service service.CollectionService
}

func (u CollectionHandler) CreateCollection(c *gin.Context) {
	var collection *dto.JsonFile
	id := uuid.New()
	var request dto.CollectionRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	request.Id = id.String()
	collection, appError := u.service.NewCollection(request)
	if appError != nil {
		c.JSON(appError.Code, appError.Message)
		return
	} else {
		response := dto.CollectionResponse{
			Id:     request.Id,
			Config: collection,
		}
		c.JSON(http.StatusOK, response)
	}

}

// func writeResponse(w http.ResponseWriter, code int, data interface{}) {
// 	w.Header().Add("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	if err := json.NewEncoder(w).Encode(data); err != nil {
// 		panic(err)
// 	}
// }
