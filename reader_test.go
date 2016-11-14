package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseNoConnection(t *testing.T) {
	reader := newS3Reader("http://localhost:9000")
	_, err := reader.ReadAndParse("test_image.png")
	if err == nil {
		t.Fatalf("Expecting an error got none")
	}
}

func TestParseNoImage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Sending something else")
	}))
	defer ts.Close()
	reader := newS3Reader(ts.URL)
	_, err := reader.ReadAndParse("test_image.png")
	if err == nil {
		t.Fatalf("Expecting an error got none")
	}
}

func TestParseImage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test_image.png")
	}))
	defer ts.Close()
	reader := newS3Reader(ts.URL)
	tags, err := reader.ReadAndParse("test_image.png")

	if err != nil {
		t.Fatalf("Got an error: %q", err)
	}
	EqualsTag(t, tags, "Model", "NIKON D750")
	EqualsTag(t, tags, "Make", "NIKON CORPORATION")
}

func TestParseImageFromFile(t *testing.T) {
	reader := newFileReader(".")
	tags, err := reader.ReadAndParse("test_image.png")

	if err != nil {
		t.Fatalf("Got an error: %q", err)
	}
	EqualsTag(t, tags, "Model", "NIKON D750")
	EqualsTag(t, tags, "Make", "NIKON CORPORATION")
}

func EqualsTag(t *testing.T, tags map[string]string, name, value string) {
	if tags[name] != value {
		t.Fatalf("Expecting tag with key '%v' and value '%v', but was %q", name, value, tags[name])
	}
}
