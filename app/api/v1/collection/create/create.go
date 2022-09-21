package createcollection

import (
	collectionservice "collection/app/api/v1/collection/service"
	"collection/app/middleware/pasetomiddleware"
	"collection/domain"
	"collection/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	// "github.com/gorilla/mux"
)

type CollectionHandler struct {
	service collectionservice.CollectionService
}

// ApplyRoutes applies router to gin Router
func ApplyRoutes(dbClient *gorm.DB, r *gin.RouterGroup) {
	collectionRepo := domain.NewCollectionRepositoryDb(dbClient)
	userRepo := domain.NewUserRepositoryDb(dbClient)

	collectionHandler := CollectionHandler{collectionservice.NewCollectionService(collectionRepo, userRepo)}
	g := r.Group("/create")
	{
		g.POST("", collectionHandler.CreateCollection)
	}
}

func (u CollectionHandler) CreateCollection(c *gin.Context) {
	var collection *dto.JsonFile
	id := uuid.New()
	var request dto.CollectionRequest

	walletAddr := c.GetString(pasetomiddleware.GinContextWalletAddress)
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	request.Id = id.String()
	collection, appError := u.service.CreateCollection(walletAddr, request)
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
