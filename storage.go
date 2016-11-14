package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"regexp"
	"strings"
)

var (
	// name of default photo bucket, master bucket containing all the photos
	PhotoBucket = []byte("PhotoBucket")
)

// Data struct describes an photo, with its name and its tags (EXIF) data
type Photo struct {
	Name string
	Tags map[string]string
}

// Storage implementation using the Bolt key, value store
type BoltStorage struct {
	db *bolt.DB
}

// Save an image to the storage engine, including the EXIF meta data (tags)
func (b *BoltStorage) Save(name string, tags map[string]string) error {
	b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(PhotoBucket)
		bkt, err := b.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		for key, value := range tags {
			err = bkt.Put([]byte(key), []byte(value))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

// Search for photos based on the given query, it is possible to use a '%' in
// the query for broader matching
func (b *BoltStorage) Search(query string) []*Photo {
	result := []*Photo{}

	b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(PhotoBucket)
		bkt.ForEach(func(k, v []byte) error {

			key := string(k)
			if compareKeys(key, query) {
				t := bkt.Bucket(k)
				tags := map[string]string{}
				t.ForEach(func(k, v []byte) error {
					tags[string(k)] = string(v)
					return nil
				})
				result = append(result, &Photo{Name: key, Tags: tags})
			}
			return nil
		})
		return nil
	})
	return result
}

// Close the database
func (b *BoltStorage) Close() {
	b.db.Close()
}

// Init the PhotoBucket bucket storage if it doesn't exists
func (b *BoltStorage) createPhotoBucket() {
	b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(PhotoBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

// Comparing keys, using a regex, replace SQL query like '%' by '.*'
func compareKeys(key, query string) bool {
	q := regexp.QuoteMeta(query)
	r, _ := regexp.Compile(strings.Replace("^"+q, "%", ".*", -1))
	return r.MatchString(key)
}

// Creates a new storage implementation based on the Bolt Key/Value store
func newBoltStorage() (Storage, error) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	s := &BoltStorage{db: db}
	s.createPhotoBucket()
	return s, nil
}

// Interface for saving and retreiving image data.
type Storage interface {
	Save(name string, tags map[string]string) error
	Search(query string) []*Photo
	Close()
}
