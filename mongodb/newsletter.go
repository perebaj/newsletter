package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Newsletter is the struct that gather what websites to scrape for an user email
type Newsletter struct {
	UserEmail string   `bson:"user_email"`
	URLs      []string `bson:"urls"`
}

// Engineer is the struct that gather the scraped content of an engineer
type Engineer struct {
	Name        string `bson:"name"`
	Description string `bson:"description"`
	URL         string `bson:"url"`
}

// Site is the struct that gather the scraped content of a website
type Site struct {
	UserEmail      string    `bson:"user_email"`
	URL            string    `bson:"url"`
	Content        string    `bson:"content"`
	ScrapeDatetime time.Time `bson:"scrape_date"`
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

// SaveEngineer saves an engineer in the database
func (m *NLStorage) SaveEngineer(ctx context.Context, e Engineer) error {
	database := m.client.Database(m.DBName)
	collection := database.Collection("engineers")
	_, err := collection.InsertOne(ctx, e)
	if err != nil {
		return err
	}
	return nil
}

// DistinctEngineerURL returns all url sites of each distinct engineer
func (m *NLStorage) DistinctEngineerURL(ctx context.Context) ([]interface{}, error) {
	database := m.client.Database(m.DBName)
	collection := database.Collection("engineers")

	resp, err := collection.Distinct(ctx, "url", bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error getting engineers: %w", err)
	}

	return resp, nil
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

// SaveSite saves a site in the database
func (m *NLStorage) SaveSite(ctx context.Context, sites []Site) error {
	database := m.client.Database(m.DBName)
	collection := database.Collection("sites")

	var docs []interface{}
	for _, site := range sites {
		docs = append(docs, site)
	}
	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

// Sites returns given an user email and a URL, the last scraped content of that URL
func (m *NLStorage) Sites(usrEmail, URL string) ([]Site, error) {
	database := m.client.Database(m.DBName)
	collection := database.Collection("sites")
	max := int64(2)

	filter := bson.M{"user_email": usrEmail, "url": URL}
	sort := bson.D{{Key: "scrape_date", Value: -1}}
	opts := options.Find().SetSort(sort)
	opts.Limit = &max

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	var sites []Site
	if err = cursor.All(context.Background(), &sites); err != nil {
		return nil, err
	}
	return sites, nil
}
