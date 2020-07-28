package main

import (
	"context"
	"net/http"
	"shopify/api"
	"shopify/util"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
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
	DEFAULT_SESSION_KEY        = "t0p-s3cr3t"
)

var logger *log.Logger
var cookieStore *sessions.CookieStore

func main() {
	port := util.GetEnv("PORT", DEFAULT_PORT)
	sessionKey := util.GetEnv("SESSION_KEY", DEFAULT_SESSION_KEY)

	logger = log.New()
	cookieStore = sessions.NewCookieStore([]byte(sessionKey))

	logger.Info("port: ", port)

	database := connectToDatabase()

	api := api.New(database, cookieStore, logger)
	router := api.CreateRouter()

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(router)

	server := http.Server{
		Addr:         ":" + port,
		Handler:      n,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Error(server.ListenAndServe())
}

func connectToDatabase() *mongo.Database {
	dbUri := util.GetEnv("BBDD_URI", DEFAULT_MONGO_URI)
	timeout := util.GetEnv("CONNECTION_TIMEOUT", DEFAULT_CONNECTION_TIMEOUT)
	dbName := util.GetEnv("BBDD_NAME", DEFAULT_BBDD_NAME)

	logger.Info("dbUri: ", dbUri)
	logger.Info("timeout: ", timeout)
	logger.Info("dbName: ", dbName)

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

	return client.Database(dbName)
}
