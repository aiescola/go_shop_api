package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"shopify/api"
	"shopify/util"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	fmt.Println("port: ", port)
	fmt.Println("dbUri: ", dbUri)
	fmt.Println("dbName: ", dbName)

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

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	names, _ := client.ListDatabaseNames(ctx, bson.M{}, options.ListDatabases())
	fmt.Println("databases: ", names)

	return client.Database(dbName), nil
}
