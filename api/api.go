package api

import (
	"shopify/api/login"
	"shopify/api/products"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Api struct {
	ProductController *products.ProductController
	LoginController   *login.LoginController
}

// Initializes the api
func New(database *mongo.Database, cookieStore *sessions.CookieStore, logger *log.Logger) *Api {
	productDataSource := products.NewMongoDataSource(database)
	productController := products.NewController(productDataSource, logger)

	loginDataSource := login.NewMongoDataSource(database)
	loginController := login.NewController(loginDataSource, cookieStore, logger)

	return &Api{
		productController,
		loginController,
	}
}
