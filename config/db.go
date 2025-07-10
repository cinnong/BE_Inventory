package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://dinaoktafiani04:dina2004@cluster0.e6hyv3q.mongodb.net/InventarisKantor?retryWrites=true&w=majority")

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("InventarisKantor")
	fmt.Println("MongoDB Atlas connected!")
}
