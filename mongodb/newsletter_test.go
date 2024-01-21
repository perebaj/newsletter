//go:build integration
// +build integration

package mongodb

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"
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

func TestNLStorageSavePage(t *testing.T) {
	ctx := context.Background()
	client, DBName := setup(ctx, t)

	database := client.Database(DBName)
	collection := database.Collection("pages")

	want := []Page{
		{IsMostRecent: true, URL: "https://www.google.com", Content: "HTML", HashMD5: md5.Sum([]byte("HTML")), ScrapeDatetime: time.Date(2023, time.August, 13, 15, 30, 0, 0, time.UTC)},
	}

	storage := NewNLStorage(client, DBName)
	err := storage.SavePage(ctx, want)
	if err != nil {
		t.Fatal("error saving page", err)
	}

	var got []Page
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		t.Fatal("error finding page", err)
	}

	if err := cursor.All(ctx, &got); err != nil {
		t.Fatal("error decoding page", err)
	}

	if len(got) == 1 {
		if !reflect.DeepEqual(got[0], want[0]) {
			t.Fatalf("got %v, want %v", got[0], want[0])
		}
	} else {
		t.Fatal("expected 1 page, got", len(got))
	}

	// if len(page) == 0 {
	// 	want[0].Updated = true
	// 	err := NLStorage.SavePage(ctx, want)
	// 	if err != nil {
	// 		t.Fatal("error saving page", err)
	// 	}
	// } else {
	// 	// Verifing if the new content is newest than the last one
	// 	lastScrapedPage, err := NLStorage.Page(ctx, want[0].URL)
	// 	if err != nil {
	// 		t.Fatal("error getting page", err)
	// 	}

	// 	if lastScrapedPage[0].HashMD5 != want[0].HashMD5 {
	// 		want[0].Updated = true
	// 		err = NLStorage.SavePage(ctx, want)
	// 		if err != nil {
	// 			t.Fatal("error saving page", err)
	// 		}
	// 	} else {
	// 		t.Fatal("page already scraped")
	// 		want[0].Updated = false
	// 		err = NLStorage.SavePage(ctx, want)
	// 		if err != nil {
	// 			t.Fatal("error saving page", err)
	// 		}
	// 	}

	t.Cleanup(teardown(ctx, client, DBName))
}

func TestNLStoragePage(t *testing.T) {
	ctx := context.Background()
	client, DBName := setup(ctx, t)

	want := []Page{
		{URL: "https://www.google.com", Content: "HTML", ScrapeDatetime: time.Date(2023, time.August, 13, 15, 30, 0, 0, time.UTC)},
		{URL: "https://www.google.com", Content: "HTML", ScrapeDatetime: time.Date(2023, time.August, 11, 15, 30, 0, 0, time.UTC)},
		{URL: "https://www.google.com", Content: "HTML", ScrapeDatetime: time.Date(2023, time.August, 12, 15, 30, 0, 0, time.UTC)},
	}

	storage := NewNLStorage(client, DBName)
	err := storage.SavePage(ctx, want)
	if err != nil {
		t.Fatal("error saving page", err)
	}

	got, err := storage.Page(ctx, "https://www.google.com")
	if err != nil {
		t.Fatal("error getting page", err)
	}

	if len(got) == 1 {
		if !reflect.DeepEqual(got[0], want[0]) {
			t.Fatalf("got %v, want %v", got[0], want[0])
		}
	} else {
		t.Fatal("expected 1 page, got", len(got))
	}

	t.Cleanup(teardown(ctx, client, DBName))
}

func TestNLStorageSaveEngineer(t *testing.T) {
	ctx := context.Background()
	client, DBName := setup(ctx, t)

	database := client.Database(DBName)
	collection := database.Collection("engineers")

	want := Engineer{
		Name: "John", URL: "https://www.1.com", Description: "John is a software engineer",
	}

	want2 := Engineer{
		Name: "John", URL: "https://www.2.com", Description: "John is a software engineer",
	}

	NLStorage := NewNLStorage(client, DBName)
	err := NLStorage.SaveEngineer(ctx, want)
	if err != nil {
		t.Fatal("error saving 1 engineer", err)
	}

	err = NLStorage.SaveEngineer(ctx, want2)
	if err != nil {
		t.Fatal("error saving 2 engineer", err)
	}

	var got []Engineer
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		t.Fatal("error finding engineer", err)
	}

	if err := cursor.All(ctx, &got); err != nil {
		t.Fatal("error decoding engineer", err)
	}

	if len(got) == 2 {
		if !reflect.DeepEqual(got, []Engineer{want, want2}) {
			t.Fatalf("got %v, want %v", got, []Engineer{want, want2})
		}
	} else {
		t.Fatal("expected 2 engineers, got", len(got))
	}

	t.Cleanup(teardown(ctx, client, DBName))
}

func TestNLStorageDistinctEngineerURLs(t *testing.T) {
	ctx := context.Background()
	client, DBName := setup(ctx, t)

	want := Engineer{
		Name: "John", URL: "https://www.1.com", Description: "John is a software engineer",
	}

	want2 := Engineer{
		Name: "John", URL: "https://www.2.com", Description: "John is a software engineer",
	}

	want3 := Engineer{
		Name: "John", URL: "https://www.2.com", Description: "John is a software engineer",
	}

	NLStorage := NewNLStorage(client, DBName)

	err := NLStorage.SaveEngineer(ctx, want)
	if err != nil {
		t.Fatal("error saving 1 engineer", err)
	}

	err = NLStorage.SaveEngineer(ctx, want2)
	if err != nil {
		t.Fatal("error saving 2 engineer", err)
	}

	err = NLStorage.SaveEngineer(ctx, want3)
	if err != nil {
		t.Fatal("error saving 3 engineer", err)
	}

	got, err := NLStorage.DistinctEngineerURLs(ctx)
	if err != nil {
		t.Fatal("error getting engineers", err)
	}

	if len(got) == 2 {
		if !reflect.DeepEqual(got, []interface{}{want.URL, want2.URL}) {
			t.Fatalf("got %v, want %v", got, []interface{}{want.URL, want2.URL})
		}
	} else {
		t.Fatal("expected 2 engineers, got", len(got))
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
	URI := os.Getenv("NL_MONGO_URI")
	client, err := OpenDB(ctx, Config{
		URI: URI,
	})
	if err != nil {
		panic(err)
	}

	DBName := t.Name() + fmt.Sprintf("%d", time.Now().UnixNano())
	return client, DBName
}
