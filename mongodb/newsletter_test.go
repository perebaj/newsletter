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

	want := Newsletter{
		UserEmail: "j@gmail.com",
		URLs:      []string{"https://www.google.com"},
	}

	NLStorage := NewNLStorage(client, DBName)
	err := NLStorage.SaveNewsletter(ctx, want)

	if err != nil {
		t.Fatal("error saving newsletter", err)
	}

	var got []Newsletter
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		t.Fatal("error finding newsletter", err)
	}

	if err := cursor.All(ctx, &got); err != nil {
		t.Fatal("error decoding newsletter", err)
	}

	if len(got) == 1 {
		if !reflect.DeepEqual(got[0], want) {
			t.Fatalf("got %v, want %v", got[0], want)
		}
	} else {
		t.Fatal("expected 1 newsletter, got", len(got))
	}

	t.Cleanup(teardown(ctx, client, DBName))
}

func TestNLStorageNewsletter(t *testing.T) {
	ctx := context.Background()
	client, DBName := setup(ctx, t)

	database := client.Database(DBName)
	collection := database.Collection("newsletter")

	want := Newsletter{
		UserEmail: "j@gmail.com",
		URLs:      []string{"https://www.google.com"},
	}
	_, err := collection.InsertOne(ctx, want)

	if err != nil {
		t.Fatal("error saving newsletter", err)
	}

	NLStorage := NewNLStorage(client, DBName)
	got, err := NLStorage.Newsletter()
	if err != nil {
		t.Fatal("error getting newsletter", err)
	}

	if len(got) == 1 {
		if !reflect.DeepEqual(got[0], want) {
			t.Fatalf("got %v, want %v", got[0], want)
		}
	} else {
		t.Fatal("expected 1 newsletter, got", len(got))
	}

	t.Cleanup(teardown(ctx, client, DBName))
}

func TestNLStorageSaveSite(t *testing.T) {
	ctx := context.Background()
	client, DBName := setup(ctx, t)

	database := client.Database(DBName)
	collection := database.Collection("sites")

	want := []Site{
		{UserEmail: "j@gmail.com", URL: "https://www.google.com", Content: "HTML", ScrapeDatetime: time.Date(2023, time.August, 14, 15, 30, 0, 0, time.UTC)},
		{UserEmail: "j@gmail.com", URL: "https://www.google.com", Content: "HTML", ScrapeDatetime: time.Date(2023, time.August, 14, 15, 30, 0, 0, time.UTC)},
		{UserEmail: "jj@gmail.com", URL: "https://www.jj.com", Content: "HTML", ScrapeDatetime: time.Date(2023, time.August, 14, 15, 30, 0, 0, time.UTC)},
	}

	NLStorage := NewNLStorage(client, DBName)
	err := NLStorage.SaveSite(ctx, want)

	if err != nil {
		t.Fatal("error saving site", err)
	}

	var got []Site
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		t.Fatal("error finding site", err)
	}

	if err := cursor.All(ctx, &got); err != nil {
		t.Fatal("error decoding site", err)
	}

	if len(got) == 3 {
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	} else {
		t.Fatal("expected 2 sites, got", len(got))
	}

	t.Cleanup(teardown(ctx, client, DBName))
}

func TestNLStorageSites(t *testing.T) {
	ctx := context.Background()
	client, DBName := setup(ctx, t)

	want := []Site{
		{UserEmail: "j@gmail.com", URL: "https://www.google.com", Content: "HTML", ScrapeDatetime: time.Date(2023, time.August, 13, 15, 30, 0, 0, time.UTC)},
		{UserEmail: "j@gmail.com", URL: "https://www.google.com", Content: "HTML", ScrapeDatetime: time.Date(2023, time.August, 12, 15, 30, 0, 0, time.UTC)},
		{UserEmail: "j@gmail.com", URL: "https://www.google.com", Content: "HTML", ScrapeDatetime: time.Date(2023, time.August, 11, 15, 30, 0, 0, time.UTC)},
	}

	NLStorage := NewNLStorage(client, DBName)
	err := NLStorage.SaveSite(ctx, want)
	if err != nil {
		t.Fatal("error saving site", err)
	}

	got, err := NLStorage.Sites("j@gmail.com", "https://www.google.com")
	if err != nil {
		t.Fatal("error getting site", err)
	}

	if len(got) == 2 {
		assert(t, got[0].UserEmail, want[0].UserEmail)
		assert(t, got[0].URL, want[0].URL)
		assert(t, got[0].Content, want[0].Content)
		assert(t, got[0].ScrapeDatetime, want[0].ScrapeDatetime)

		assert(t, got[1].UserEmail, want[1].UserEmail)
		assert(t, got[1].URL, want[1].URL)
		assert(t, got[1].Content, want[1].Content)
		assert(t, got[1].ScrapeDatetime, want[1].ScrapeDatetime)
	} else {
		t.Fatal("expected 2 sites, got", len(got))
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
