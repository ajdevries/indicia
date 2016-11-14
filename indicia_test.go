package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	mux.HandleFunc("/", index)
	mux.HandleFunc("/key_1.jpg", image)
	mux.HandleFunc("/key_2.jpg", image)
	mux.HandleFunc("/key_3.jpg", image)
	mux.HandleFunc("/key_4.jpg", image)
	defer ts.Close()

	s, _ := newBoltStorage()
	defer s.Close()

	i := newIndicia(ts.URL, newS3Lister, newS3Reader, s)
	i.Start()

	photos := s.Search("key%")
	if len(photos) != 5 {
		t.Fatalf("Expecting 5 photos but got %d results!", len(photos))
	}

	EqualsTag(t, photos[0].Tags, "Model", "NIKON D750")
	EqualsTag(t, photos[1].Tags, "Model", "NIKON D750")
	EqualsTag(t, photos[2].Tags, "Model", "NIKON D750")
	EqualsTag(t, photos[3].Tags, "Model", "NIKON D750")
	if (len(photos[4].Tags)) > 0 {
		t.Fatalf("Expecting zero tags for the last photo, but got", len(photos[4].Tags))
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,
		`<?xml version="1.0" encoding="UTF-8"?>
			<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
			<Contents>
				<Key>key_1.jpg</Key>
			</Contents>
			<Contents>
				<Key>key_2.jpg</Key>
			</Contents>
			<Contents>
				<Key>key_3.jpg</Key>
			</Contents>
			<Contents>
				<Key>key_4.jpg</Key>
			</Contents>
			<Contents>
				<Key>key_5.jpg</Key>
			</Contents>
		</ListBucketResult>`)
}
func image(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "test_image.png")
}
