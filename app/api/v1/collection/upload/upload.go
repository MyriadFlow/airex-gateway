package deploycollection

import (
	collectionservice "collection/app/api/v1/collection/service"
	"collection/app/middleware/pasetomiddleware"
	"collection/domain"
	"collection/logger"

	"github.com/TheLazarusNetwork/go-helpers/httpo"
	"github.com/gin-gonic/gin"
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
	g := r.Group("/upload")
	{
		g.GET("", collectionHandler.UploadCollection)
	}
}

func (u CollectionHandler) UploadCollection(c *gin.Context) {
	colId := c.Query("id")
	walletAddr := c.GetString(pasetomiddleware.GinContextWalletAddress)
	err := u.service.UploadCollection(walletAddr, colId)
	if err != nil {
		logger.Errorf("failed to upload collection: %s", err)
		httpo.NewErrorResponse(500, "failed to upload collection").Send(c, 500)
	}
}
