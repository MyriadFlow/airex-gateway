package app

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"collection/config/envconfig"
	"collection/domain"
	"collection/dto"
	"collection/logger"
	"collection/service"

	"gorm.io/driver/postgres"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/gorilla/mux"
)

func Start() {

	envconfig.InitEnvVars()

	// r := mux.NewRouter()

	dbClient := getDbClient()

	//wiring
	newRepositoryDb := domain.NewCollectionRepositoryDb(dbClient.Model(&dto.User{}))

	us := CollectionHandler{service.NewCollectionService(newRepositoryDb)}

	// r.HandleFunc("/collections", us.CreateCollection).Methods("Post")
	ginApp := gin.Default()
	ginApp.POST("/collection", us.CreateCollection)

	corsM := cors.New(cors.Config{AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowOrigins:     envconfig.EnvVars.ALLOWED_ORIGIN})
	ginApp.Use(corsM)
	// api.ApplyRoutes(ginApp)
	port := envconfig.EnvVars.APP_PORT
	err := ginApp.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalf("failed to serve app on port %s: %s", port, err)
	}

}

func getDbClient() *gorm.DB {
	dbUser := envconfig.EnvVars.DB_USERNAME
	dbPasswd := envconfig.EnvVars.DB_PASSWORD
	dbHost := envconfig.EnvVars.DB_HOST
	dbPort := envconfig.EnvVars.DB_PORT
	dbName := envconfig.EnvVars.DB_NAME
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%d",
		dbHost, dbUser, dbPasswd, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("failed to connect to database: %s", err)
	}

	// Get underlying sql database to ping it
	sqlDb, err := db.DB()
	if err != nil {
		logger.Fatalf("failed to ping database: %s", err)
	}

	// If ping fails then log error and exit
	if err = sqlDb.Ping(); err != nil {
		logger.Fatalf("failed to ping database: %s", err)
	}

	err = db.AutoMigrate(&domain.User{}, &domain.Collection{}, &dto.FlowId{})
	if err != nil {
		logger.Fatalf("Auto migration failed: %s", err)
	}

	err = db.SetupJoinTable(&domain.Collection{}, "Sellers", &domain.CollectionSeller{})
	if err != nil {
		log.Fatalf("Auto migration failed: %s", err)
	}
	err = db.AutoMigrate(&domain.CollectionSeller{})
	if err != nil {
		log.Fatalf("Auto migration failed: %s", err)
	}
	logger.Info("Database is Connected")
	return db
}
