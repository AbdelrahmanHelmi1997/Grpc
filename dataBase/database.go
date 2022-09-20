package dataBase

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://abdelrahmanhelmi:somini2A@cluster0-shard-00-00.x1ifh.mongodb.net:27017,cluster0-shard-00-01.x1ifh.mongodb.net:27017,cluster0-shard-00-02.x1ifh.mongodb.net:27017/?ssl=true&replicaSet=atlas-ibytxo-shard-0&authSource=admin&retryWrites=true&w=majority"
const dbName = "Users"
const collName = "Users"

var UsersDB *mongo.Collection

func DB() {

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal("mongoshit : ", err)
	}

	fmt.Println("Mongo Connection Success")

	UsersDB = client.Database(dbName).Collection(collName)

}
