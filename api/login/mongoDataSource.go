package login

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoLoginDataSource struct {
	collection *mongo.Collection
}

func (ds *mongoLoginDataSource) getPassword(user string) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result := ds.collection.FindOne(ctx, bson.M{"user": user})

	if err := result.Err(); err != nil {
		return nil, err
	}

	credentials := new(Credentials)
	if err := result.Decode(credentials); err != nil {
		return nil, err
	}

	return &credentials.Password, nil

}

func (ds *mongoLoginDataSource) createUser(credentials Credentials) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := ds.collection.InsertOne(ctx, credentials)
	if err != nil {
		return err
	}

	return nil
}
