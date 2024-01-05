package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Newsletter is the struct that gather what websites to scrape for an user email
type Newsletter struct {
	UserEmail string   `bson:"user_email"`
	URLs      []string `bson:"urls"`
}

// Site is the struct that gather the scraped content of a website
type Site struct {
	UserEmail  string    `bson:"user_email"`
	URL        string    `bson:"url"`
	Content    string    `bson:"content"`
	ScrapeDate time.Time `bson:"scrape_date"`
}

// NLStorage joins the Mongo operations for the Newsletter collection
type NLStorage struct {
	client *mongo.Client
	DBName string
}

// NewNLStorage initializes a new NLStorage
func NewNLStorage(client *mongo.Client, DBName string) *NLStorage {
	return &NLStorage{
		client: client,
		DBName: DBName,
	}
}

// SaveNewsletter saves a newsletter in the database
func (m *NLStorage) SaveNewsletter(ctx context.Context, newsletter Newsletter) error {
	database := m.client.Database(m.DBName)
	collection := database.Collection("newsletter")
	_, err := collection.InsertOne(ctx, newsletter)
	if err != nil {
		return err
	}
	return nil
}

// Newsletter returns all the newsletters in the database
func (m *NLStorage) Newsletter() ([]Newsletter, error) {
	var newsletters []Newsletter
	database := m.client.Database(m.DBName)
	collection := database.Collection("newsletter")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.Background(), &newsletters); err != nil {
		return nil, err
	}

	return newsletters, nil
}
