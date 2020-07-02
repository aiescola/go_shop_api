package main

import (
	"context"
	"log"
	"net/http"
	"shopify/api"
	"shopify/util"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DEFAULT_PORT      = "8080"
	DEFAULT_MONGO_URI = "mongodb://localhost:27017"
	DEFAULT_BBDD_NAME = "shopify"
)

func main() {
	port := util.GetEnv("PORT", DEFAULT_PORT)
	dbUri := util.GetEnv("BBDD_URI", DEFAULT_MONGO_URI)
	dbName := util.GetEnv("BBDD_NAME", DEFAULT_BBDD_NAME)

	database, err := createDatabase(dbUri, dbName)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	api.Initialize(database, router)

	http.ListenAndServe(":"+port, router)
}

func createDatabase(dbUri string, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))

	if err != nil {
		log.Fatal(err)
	}

	return client.Database(dbName), nil
}
