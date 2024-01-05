//go:build integration
// +build integration

package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNLStorageSaveNewsletter(t *testing.T) {
	ctx := context.Background()
	client, DBName := setup(ctx, t)

	database := client.Database(DBName)
	collection := database.Collection("newsletter")

	NLStorage := NewNLStorage(client, DBName)
	err := NLStorage.SaveNewsletter(ctx, Newsletter{
		UserEmail: "j@gmail.com",
		URLs:      []string{"https://www.google.com"},
	})

	if err != nil {
		t.Fatal("error saving newsletter", err)
	}

	var nls []Newsletter
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		t.Fatal("error finding newsletter", err)
	}

	if err := cursor.All(ctx, &nls); err != nil {
		t.Fatal("error decoding newsletter", err)
	}

	if len(nls) == 1 {
		reflect.DeepEqual(nls[0], Newsletter{
			UserEmail: "j@gmail.com",
			URLs:      []string{"https://www.google.com"},
		})
	} else {
		t.Fatal("expected 1 newsletter, got", len(nls))
	}

	t.Cleanup(teardown(ctx, client, DBName))
}

func TestNLStorageNewsletter(t *testing.T) {
	ctx := context.Background()
	client, DBName := setup(ctx, t)

	database := client.Database(DBName)
	collection := database.Collection("newsletter")

	_, err := collection.InsertOne(ctx, Newsletter{
		UserEmail: "j@gmail.com",
		URLs:      []string{"https://www.google.com"},
	})

	if err != nil {
		t.Fatal("error saving newsletter", err)
	}

	NLStorage := NewNLStorage(client, DBName)
	nls, err := NLStorage.Newsletter()
	if err != nil {
		t.Fatal("error getting newsletter", err)
	}

	if len(nls) == 1 {
		reflect.DeepEqual(nls[0], Newsletter{
			UserEmail: "j@gmail.com",
			URLs:      []string{"https://www.google.com"},
		})
	} else {
		t.Fatal("expected 1 newsletter, got", len(nls))
	}

	t.Cleanup(teardown(ctx, client, DBName))
}

func assert(t testing.TB, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func teardown(ctx context.Context, client *mongo.Client, DBName string) func() {
	return func() {
		if err := client.Database(DBName).Drop(ctx); err != nil {
			panic(err)
		}
	}
}

func setup(ctx context.Context, t testing.TB) (*mongo.Client, string) {
	// TODO: Receive the URI from the environment variable
	URI := "mongodb://root:root@mongodb:27017/"
	client, err := OpenDB(ctx, Config{
		URI: URI,
	})
	if err != nil {
		panic(err)
	}

	DBName := t.Name() + fmt.Sprintf("%d", time.Now().UnixNano())
	return client, DBName
}
