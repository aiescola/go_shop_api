package login

import "go.mongodb.org/mongo-driver/mongo"

type LoginDataSource interface {
	getPassword(user string) (*string, error)
	createUser(credentials Credentials) error
}

func NewMongoDataSource(database *mongo.Database) *mongoLoginDataSource {
	return &mongoLoginDataSource{
		collection: database.Collection("credentials"),
	}
}
