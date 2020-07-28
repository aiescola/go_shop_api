package products

import "go.mongodb.org/mongo-driver/mongo"

type ProductDataSource interface {
	GetProducts() ([]Product, error)
	GetProduct(code string) (*Product, error)
	AddProduct(product Product) error
}

func NewMongoDataSource(database *mongo.Database) mongoDataSource {
	return mongoDataSource{
		collection: database.Collection("products"),
	}
}
