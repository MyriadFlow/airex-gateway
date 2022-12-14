package apiv1

import (
	"collection/app/api/v1/authenticate"
	"collection/app/api/v1/collection"
	"collection/app/api/v1/flowid"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ApplyRoutes applies the /v1.0 group and all child routes to given gin RouterGroup
func ApplyRoutes(dbClient *gorm.DB, r *gin.RouterGroup) {
	v1 := r.Group("/v1.0")
	{
		flowid.ApplyRoutes(dbClient, v1)
		authenticate.ApplyRoutes(dbClient, v1)

		collection.ApplyRoutes(dbClient, v1)
	}
}
