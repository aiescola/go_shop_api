package api

import (
	"shopify/api/products"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// Initializes the api
func Initialize(database *mongo.Database, router *mux.Router, logger *log.Logger) {
	//productDataSource := products.NewLocalProductDataSource()
	productDataSource := products.NewMongoProductDataSource(database)
	productService := products.NewService(productDataSource, logger)

	productService.AddRoutes(router)
}
