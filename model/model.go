package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"firstName"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	Token     string             `bson:"token"`
}
