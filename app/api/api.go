// Package api provide support to create /api group
package api

import (
	apiv1 "collection/app/api/v1"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ApplyRoutes applies the /api group and v1 routes to given gin Engine
func ApplyRoutes(db *gorm.DB, r *gin.Engine) {
	api := r.Group("/api")
	{
		apiv1.ApplyRoutes(db, api)
	}
}
