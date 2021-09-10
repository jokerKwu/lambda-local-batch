package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func GetClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://ryan:W4JYJevdgfednj4c@ryan.ke0pv.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println(err)
	}
	defer cancel()
	return client, err
}
