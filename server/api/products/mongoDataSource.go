package products

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDataSource struct {
	collection *mongo.Collection
}

func (ds mongoDataSource) GetProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := ds.collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}

	var products []Product
	for cursor.Next(ctx) {
		var p Product
		cursor.Decode(&p)
		products = append(products, p)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (ds mongoDataSource) GetProduct(code string) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result := ds.collection.FindOne(ctx, bson.M{"code": code})

	if err := result.Err(); err != nil {
		return nil, err
	}
	product := new(Product)
	if err := result.Decode(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (ds mongoDataSource) AddProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := ds.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}

	return nil
}
