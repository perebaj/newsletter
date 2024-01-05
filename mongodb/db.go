// Package mongodb provides a MongoDB connection.
package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config is the configuration for the MongoDB connection.
type Config struct {
	URI string
}

// OpenDB connects to the MongoDB instance.
func OpenDB(ctx context.Context, cfg Config) (*mongo.Client, error) {
	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
		NilSliceAsEmpty:   true,
	}

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(cfg.URI).
		SetBSONOptions(bsonOpts))

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}