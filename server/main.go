package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aiescola/go_shop_api/api"
	"github.com/aiescola/go_shop_api/util"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DEFAULT_PORT        = "8080"
	DEFAULT_MONGO_URI   = "mongodb://mongo:27017"
	DEFAULT_BBDD_NAME   = "go_shop_db"
	DEFAULT_SESSION_KEY = "t0p-s3cr3t"
	DEFAULT_REDIS_ADDR  = "redis:6379"
)

var logger *log.Logger
var cookieStore *sessions.CookieStore

func main() {
	port := util.GetEnv("PORT", DEFAULT_PORT)
	sessionKey := util.GetEnv("SESSION_KEY", DEFAULT_SESSION_KEY)

	logger = log.New()
	cookieStore = sessions.NewCookieStore([]byte(sessionKey))

	logger.Info("port: ", port)

	client := createRedisClient()
	fmt.Println("client: ", client)
	database := connectToDatabase()

	api := api.New(database, client, cookieStore, logger)
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
	dbName := util.GetEnv("BBDD_NAME", DEFAULT_BBDD_NAME)

	logger.Info("dbUri: ", dbUri)
	logger.Info("dbName: ", dbName)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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

func createRedisClient() *redis.Client {
	addr := util.GetEnv("REDIS_ADDR", DEFAULT_REDIS_ADDR)
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()

	fmt.Println(pong, err)

	if err != nil {
		logger.Fatal(err)
	}

	return client
}
