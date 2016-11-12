package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestS3ListerNoConnection(t *testing.T) {
	l := newS3Lister("http://localhost:1234")
	_, err := l.List()
	if err == nil {
		t.Fatalf("Expecting an error but got none")
	}
}

func TestS3ListerNoXML(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "I'm not sending XML!")
	}))
	defer ts.Close()

	l := newS3Lister(ts.URL)
	_, err := l.List()
	if err == nil {
		t.Fatalf("Expecting an error but got none")
	}
}

func TestS3ListerValidXML(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w,
			`<?xml version="1.0" encoding="UTF-8"?>
				<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
				<Contents>
					<Key>only_one_key.jpg</Key>
				</Contents>
			</ListBucketResult>`)
	}))
	defer ts.Close()

	l := newS3Lister(ts.URL)
	urls, err := l.List()
	if err != nil {
		t.Fatalf("Expecting an error but got none")
	}

	k := urls[0]
	if k != "only_one_key.jpg" {
		t.Fatalf("Expecting key 'only_one_key.jpg' but got %q", k)
	}
}
