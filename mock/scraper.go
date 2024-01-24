package mock

import (
	"context"
	"crypto/md5"

	"github.com/perebaj/newsletter/mongodb"
)

type StorageMockImpl struct{}

const FakeURL = "http://fakeurl.test"

func NewStorageMock() StorageMockImpl                                        { return StorageMockImpl{} }
func (s StorageMockImpl) SavePage(_ context.Context, _ []mongodb.Page) error { return nil }
func (s StorageMockImpl) DistinctEngineerURLs(_ context.Context) ([]interface{}, error) {
	return []interface{}{FakeURL}, nil
}
func (s StorageMockImpl) Page(_ context.Context, _ string) ([]mongodb.Page, error) {
	return []mongodb.Page{}, nil
}
func (s StorageMockImpl) Newsletter() ([]mongodb.Newsletter, error) {
	return []mongodb.Newsletter{{URLs: []string{FakeURL}}}, nil
}
func (s StorageMockImpl) PageIn(_ context.Context, _ []string) ([]mongodb.Page, error) {
	return []mongodb.Page{
		{IsMostRecent: true, URL: FakeURL, Content: "Hello, World!", HashMD5: md5.Sum([]byte("Hello, World!"))},
		{IsMostRecent: true, URL: FakeURL, Content: "Hello, World! 2", HashMD5: md5.Sum([]byte("Hello, World! 2"))},
	}, nil
}
