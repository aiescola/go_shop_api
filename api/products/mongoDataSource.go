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

func (dataSource *mongoDataSource) GetOne(code string) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result := dataSource.collection.FindOne(ctx, bson.M{"code": code})

	if err := result.Err(); err != nil {
		return nil, err
	}
	product := new(Product)
	if err := result.Decode(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (dataSource *mongoDataSource) GetProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := dataSource.collection.Find(ctx, bson.M{})
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

func (dataSource *mongoDataSource) AddProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := dataSource.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}

	return nil
}
