package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

// Page is the struct that gather the scraped content of a website
type Page struct {
	URL            string    `bson:"url"`
	Content        string    `bson:"content"`
	ScrapeDatetime time.Time `bson:"scrape_date"`
	HashMD5        [16]byte  `bson:"hash_md5"`
	IsMostRecent   bool      `bson:"is_most_recent"`
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

// DistinctEngineerURLs returns all url sites of each distinct engineer
func (m *NLStorage) DistinctEngineerURLs(ctx context.Context) ([]interface{}, error) {
	database := m.client.Database(m.DBName)
	collection := database.Collection("engineers")

	resp, err := collection.Distinct(ctx, "url", bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error getting engineers: %v", err)
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

// SavePage saves the scraped content of a website
func (m *NLStorage) SavePage(ctx context.Context, pages []Page) error {
	database := m.client.Database(m.DBName)
	collection := database.Collection("pages")

	var docs []interface{}
	for _, site := range pages {
		docs = append(docs, site)
	}
	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

// PageIn returns the last scraped content of a given list of urls
func (m *NLStorage) PageIn(ctx context.Context, urls []string) ([]Page, error) {
	database := m.client.Database(m.DBName)
	collection := database.Collection("pages")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"url": bson.M{
					"$in": urls,
				},
			},
		},
		{
			"$sort": bson.M{
				"scrape_date": -1,
			},
		},
		{
			"$group": bson.M{
				"_id": "$url",
				"page": bson.M{
					"$first": "$$ROOT",
				},
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$page",
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error getting page: %v", err)
	}

	var page []Page
	if err = cursor.All(ctx, &page); err != nil {
		return page, fmt.Errorf("error decoding page: %v", err)
	}

	return page, nil
}

// Page returns the last scraped content of a given url
func (m *NLStorage) Page(ctx context.Context, url string) ([]Page, error) {
	var page []Page
	database := m.client.Database(m.DBName)
	collection := database.Collection("pages")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"url": url,
			},
		},
		{
			"$sort": bson.M{
				"scrape_date": -1,
			},
		},
		{
			"$limit": 1,
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return page, fmt.Errorf("error getting page: %v", err)
	}

	if err = cursor.All(ctx, &page); err != nil {
		return page, fmt.Errorf("error decoding page: %v", err)
	}

	return page, nil
}
