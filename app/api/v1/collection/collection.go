package collection

import (
	createcollection "collection/app/api/v1/collection/create"
	deploycollection "collection/app/api/v1/collection/deploy"
	uploadcollection "collection/app/api/v1/collection/upload"
	"collection/app/middleware/pasetomiddleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ApplyRoutes applies the /v1.0 group and all child routes to given gin RouterGroup
func ApplyRoutes(dbClient *gorm.DB, r *gin.RouterGroup) {
	collectionGrp := r.Group("/collection")
	{
		pasetoMiddleWare := pasetomiddleware.PASETOMiddleWareService{
			Db: dbClient,
		}
		collectionGrp.Use(pasetoMiddleWare.PASETO)
		createcollection.ApplyRoutes(dbClient, collectionGrp)
		uploadcollection.ApplyRoutes(dbClient, collectionGrp)
		deploycollection.ApplyRoutes(dbClient, collectionGrp)
	}
}
