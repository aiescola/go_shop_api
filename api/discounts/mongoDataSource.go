package discounts

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDiscountDataSource struct {
	collection *mongo.Collection
}

func (ds mongoDiscountDataSource) GetDiscounts() ([]Discount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := ds.collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}

	var discounts []Discount
	for cursor.Next(ctx) {
		var d Discount
		cursor.Decode(&d)
		discounts = append(discounts, d)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return discounts, nil
}

func (ds mongoDiscountDataSource) GetDiscount(code string) (*Discount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result := ds.collection.FindOne(ctx, bson.M{"code": code})

	if err := result.Err(); err != nil {
		return nil, err
	}
	discount := new(Discount)
	if err := result.Decode(discount); err != nil {
		return nil, err
	}

	return discount, nil
}

func (ds mongoDiscountDataSource) AddDiscount(discount Discount) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := ds.collection.InsertOne(ctx, discount)
	if err != nil {
		return err
	}

	return nil
}
