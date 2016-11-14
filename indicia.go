package main

import (
	"fmt"
	"log"
	"time"
)

// Factory method for creating a Lister
type NewLister func(base string) Lister

// Factory method for creating a Reader
type NewReader func(base string) Reader

type Indicia struct {
	// URL where the S3 bucket resists
	URL string
	// Factory method for creating a Lister implementation
	Lister NewLister
	// Factory method for creating a Reader implementation
	Reader NewReader
	// Storage engine to save all the photos and their EXIF info (tags)
	Storage Storage
	// Number of photos that are found in the Lister
	count int
	// Number of photos that are indexed
	indexed int
	// Time it took to index all the photos
	elapsed time.Duration
}

// Index photos based in the given configuration Options
func newIndicia(URL string, Lister NewLister, Reader NewReader, Storage Storage) *Indicia {
	return &Indicia{URL: URL, Lister: Lister, Reader: Reader, Storage: Storage}
}

func (i *Indicia) Start() {
	start := time.Now()
	log.Printf("Getting photos from %s\n", i.URL)
	l := i.Lister(i.URL)
	list, _ := l.List()
	i.count = len(list)
	i.indexed = 0
	for _, photo := range list {
		fmt.Printf("\r%d/%d", i.indexed, i.count)
		r := i.Reader(i.URL)
		tags, _ := r.ReadAndParse(photo)
		i.Storage.Save(photo, tags)
		i.indexed += 1
		i.elapsed = time.Since(start)
	}
	fmt.Printf("\n\nIndexing took %s\n", i.elapsed)
}
