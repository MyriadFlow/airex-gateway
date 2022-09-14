package app

import (
	"collection/dto"
	"collection/service"
	"encoding/json"
	"net/http"

	// "github.com/gorilla/mux"
	"github.com/google/uuid"
)

type UserHandler struct {
	service service.UserService
}

func (u UserHandler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var collection *dto.JsonFile
	id := uuid.New()
	var request dto.CollectionRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	}
	request.User_id = id.String()
	collection, appError := u.service.NewCollection(request)
	if appError != nil {
		writeResponse(w, appError.Code, appError.Message)
	} else {
			response := dto.CollectionResponse{
			Id:     request.User_id,
			Config: collection,
		}
		writeResponse(w,http.StatusOK,response)
	}

}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
