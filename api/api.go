package api

import (
	"shopify/api/products"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// Initializes the api
func Initialize(database *mongo.Database, router *mux.Router) {
	//productDataSource := products.MakeLocalProductDataSource()
	productDataSource := products.NewMongoProductDataSource(database)
	productService := products.NewService(productDataSource)

	productService.AddRoutes(router)
}
