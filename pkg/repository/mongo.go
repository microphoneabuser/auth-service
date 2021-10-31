package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	authCollection = "users"
	timeout        = 10 * time.Second
)

type MongoConfig struct {
	Uri      string
	Username string
	Password string
}

func NewMongoClient(config MongoConfig) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(config.Uri)
	if config.Username != "" && config.Password != "" {
		opts.SetAuth(options.Credential{
			Username: config.Username, Password: config.Password,
		})
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
