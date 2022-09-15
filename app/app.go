package app

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"log"

	_ "github.com/lib/pq"

	"collection/domain"
	"collection/logger"
	"collection/service"
	"github.com/jmoiron/sqlx"

	
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func Start(){
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	r :=mux.NewRouter()

	dbClient := getDbClient()
	//wiring
	newRepositoryDb := domain.NewUserRepositoryDb(dbClient)
	us := UserHandler{service.NewUserService(newRepositoryDb)}
	
	r.HandleFunc("/collections",us.CreateCollection).Methods("Post")

	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), r))
	
}

func getDbClient() *sqlx.DB{
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbAddr := os.Getenv("DB_ADDR")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	db:= os.Getenv("DB")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?sslmode=disable", dbUser, dbPasswd, dbAddr, dbPort, dbName)

	client, err := sqlx.Open(db, dataSource)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 10)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	logger.Info("Database is Connected")
	return client
}


