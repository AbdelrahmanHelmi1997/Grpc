package dataBase

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://abdelrahmanhelmi:somini2A@cluster0-shard-00-00.x1ifh.mongodb.net:27017,cluster0-shard-00-01.x1ifh.mongodb.net:27017,cluster0-shard-00-02.x1ifh.mongodb.net:27017/?ssl=true&replicaSet=atlas-ibytxo-shard-0&authSource=admin&retryWrites=true&w=majority"
const dbName = "Shop"
const collName = "Items"
const dbName1 = "ShoppingCart"
const collName2 = "Items"
const dbName3 = "ShopUsers"
const collName3 = "Users"

var Collection *mongo.Collection
var ShoppingCartColliction *mongo.Collection
var UsersCollection *mongo.Collection

func DB() {

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal("mongoshit : ", err)
	}

	fmt.Println("Mongo Connection Success")

	Collection = client.Database(dbName).Collection(collName)
	ShoppingCartColliction = client.Database(dbName1).Collection(collName2)
	UsersCollection = client.Database(dbName3).Collection(collName3)
}
