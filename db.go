package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDatabase(ctx context.Context, uri string) (*mongo.Client, error) {
	serverAPIOpts := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOpts)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	return client, err
}

func disconnectDatabase(ctx context.Context, client *mongo.Client) error {
	return client.Disconnect(ctx)
}
