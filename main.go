package main

import (
	"context"
	"net/http"
	"shopify/api"
	"shopify/util"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DEFAULT_PORT               = "8080"
	DEFAULT_MONGO_URI          = "mongodb://localhost:27017"
	DEFAULT_BBDD_NAME          = "shopify"
	DEFAULT_CONNECTION_TIMEOUT = "30"
)

var logger *log.Logger

func main() {
	port := util.GetEnv("PORT", DEFAULT_PORT)
	dbUri := util.GetEnv("BBDD_URI", DEFAULT_MONGO_URI)
	dbName := util.GetEnv("BBDD_NAME", DEFAULT_BBDD_NAME)

	logger = log.New()

	logger.Info("port: ", port)
	logger.Info("dbUri: ", dbUri)
	logger.Info("dbName: ", dbName)

	database, err := createDatabase(dbUri, dbName)
	if err != nil {
		logger.Fatal(err)
	}

	router := mux.NewRouter()

	api.Initialize(database, router, logger)

	logger.Error(http.ListenAndServe(":"+port, router))
}

func createDatabase(dbUri string, dbName string) (*mongo.Database, error) {
	timeout := util.GetEnv("CONNECTION_TIMEOUT", DEFAULT_CONNECTION_TIMEOUT)
	logger.Info("CONNECTION_TIMEOUT: ", timeout)

	connectionTimeout, err := strconv.Atoi(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(connectionTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))

	if err != nil {
		logger.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatal(err)
	}

	names, _ := client.ListDatabaseNames(ctx, bson.M{}, options.ListDatabases())
	logger.Info("databases: ", names)

	return client.Database(dbName), nil
}
