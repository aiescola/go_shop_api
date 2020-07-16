package main

import (
	"context"
	"net/http"
	"shopify/api"
	"shopify/middleware"
	"shopify/util"
	"strconv"
	"time"

	"html/template"

	"github.com/gorilla/mux"
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
var templates *template.Template

func main() {
	port := util.GetEnv("PORT", DEFAULT_PORT)
	dbUri := util.GetEnv("BBDD_URI", DEFAULT_MONGO_URI)
	dbName := util.GetEnv("BBDD_NAME", DEFAULT_BBDD_NAME)
	sessionKey := util.GetEnv("SESSION_KEY", DEFAULT_SESSION_KEY)

	logger = log.New()
	cookieStore = sessions.NewCookieStore([]byte(sessionKey))
	templates = template.Must(template.ParseGlob("templates/*.html"))

	logger.Info("port: ", port)
	logger.Info("dbUri: ", dbUri)
	logger.Info("dbName: ", dbName)

	mongoClient := connectToDatabase(dbUri)

	router := createRouter(mongoClient.Database(dbName))

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(router)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      n,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Error(server.ListenAndServe())
}

func createRouter(database *mongo.Database) *mux.Router {
	router := mux.NewRouter()

	api := api.New(database, cookieStore, logger)

	var middleware = middleware.NewAuthMiddleware(cookieStore, logger)

	router.HandleFunc("/api/products", api.ProductController.GetProducts).Methods("GET")
	router.HandleFunc("/api/products/{Code}", api.ProductController.GetProduct).Methods("GET")
	router.HandleFunc("/api/products", middleware.Intercept(api.ProductController.AddProduct)).Methods("POST")

	router.HandleFunc("/register", api.LoginController.Register).Methods("POST")
	router.HandleFunc("/login", api.LoginController.Login).Methods("POST")

	//TODO: remove this, it's a test
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "login.html", nil)
	}).Methods("GET")

	return router
}

func connectToDatabase(dbUri string) *mongo.Client {
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

	return client
}
