package products

import "go.mongodb.org/mongo-driver/mongo"

type ProductDataSource interface {
	GetProducts() ([]Product, error)
	GetOne(code string) (*Product, error)
	AddProduct(product Product) error
}

func NewLocalProductDataSource() *localDataSource {
	return &localDataSource{
		[]Product{{"MUG", "Mug", 5.4}, {"TSHIRT", "T-Shirt", 12.5}, {"PEN", "Pen", 3.2}},
	}
}

func NewMongoProductDataSource(database *mongo.Database) *mongoDataSource {
	return &mongoDataSource{
		collection: database.Collection("products"),
	}
}
