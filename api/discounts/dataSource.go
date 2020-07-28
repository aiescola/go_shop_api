package discounts

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type DiscountDataSource interface {
	GetDiscounts() ([]Discount, error)
	GetDiscount(code string) (*Discount, error)
	AddDiscount(discount Discount) error
}

func NewMongoDataSource(database *mongo.Database) mongoDiscountDataSource {
	return mongoDiscountDataSource{
		collection: database.Collection("discounts"),
	}
}
