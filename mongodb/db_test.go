//go:build integration
// +build integration

package mongodb

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func setup(ctx context.Context, t testing.TB) (*mongo.Collection, func()) {
	client, err := Connect(ctx, Config{
		URI: "mongodb://root:root@mongodb:27017",
	})

	if err != nil {
		panic(err)
	}

	randCollection := t.Name()

	collection := client.Database("test").Collection(randCollection)

	teardown := func() {
		err = collection.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}

	return collection, teardown
}

func TestConnect(t *testing.T) {
	collection, teardown := setup(context.Background(), t)

	defer teardown()

	resp, err := collection.InsertOne(context.Background(), map[string]string{"name": "pi", "value": "3.14159"})

	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.InsertedID)

	result := collection.FindOne(context.Background(), map[string]string{"name": "pi", "value": "3.14159"})
	t.Log(result.Raw())

}
