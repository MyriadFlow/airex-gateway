package app

import (
	"fmt"
	"os"
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
	"github.com/gorilla/mux"
)

func Start() {
	envconfig.InitEnvVars()

	r := mux.NewRouter()

	dbClient := getDbClient()
	//wiring
	newRepositoryDb := domain.NewUserRepositoryDb(dbClient.Model(&dto.User{}))
	us := UserHandler{service.NewUserService(newRepositoryDb)}

	r.HandleFunc("/collections", us.CreateCollection).Methods("Post")
	ginApp := gin.Default()

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
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%s",
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

	logger.Info("Database is Connected")
	return db
}
