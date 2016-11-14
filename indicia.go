package main

import (
	"fmt"
	"time"
)

// Factory method for creating a Lister
type NewLister func(base string) Lister
// Factory method for creating a Reader
type NewReader func(base string) Reader

type Options struct {
	// URL where the S3 bucket resists
	URL     string
	// Factory method for creating a Lister implementation
	Lister  NewLister
	// Factory method for creating a Reader implementation
	Reader  NewReader
	// Storage engine to save all the photos and their EXIF info (tags)
	Storage Storage
}

// Index photos based in the given configuration Options
func Index(options *Options) {
	start := time.Now()

	fmt.Printf("Getting photos from %s\n", options.URL)
	l := options.Lister(options.URL)
	list, _ := l.List()
	n := len(list)
	for i, photo := range list {
		fmt.Printf("\r%d/%d", i+1, n)
		r := options.Reader(options.URL)
		tags, _ := r.ReadAndParse(photo)
		options.Storage.Save(photo, tags)
	}
	elapsed := time.Since(start)
	fmt.Printf("\n\nIndexing took %s\n", elapsed)
}
