package main

import (
	"testing"
)

func TestSaveAndSearchImage(t *testing.T) {
	s, err := newBoltStorage()
	if err != nil {
		t.Fatalf("Can't create Bolt storage!")
	}
	defer s.Close()

	s.Save("test.jpg", map[string]string{"Make": "Apple"})
	photos := s.Search("test.jpg")

	if len(photos) == 0 {
		t.Fatalf("Expecting an Photo but got zero results!")
	}
	if photos[0].Tags["Make"] != "Apple" {
		t.Fatalf("Expecting tag with key 'Make' and value 'Apple', but was %q", photos[0].Tags["Make"])
	}
}

func TestQuery(t *testing.T) {
	s, err := newBoltStorage()
	if err != nil {
		t.Fatalf("Can't create Bolt storage!")
	}
	defer s.Close()

	photos := s.Search("test%")
	if len(photos) != 1 {
		t.Fatalf("Expecting one photo but got %d results!", len(photos))
	}
}

func TestQueryWithTwoPhotos(t *testing.T) {
	s, err := newBoltStorage()
	if err != nil {
		t.Fatalf("Can't create Bolt storage!")
	}
	defer s.Close()

	s.Save("another_test.jpg", map[string]string{"Make": "Apple"})
	photos := s.Search("%test%")
	if len(photos) != 2 {
		t.Fatalf("Expecting two photos but got zero results!")
	}
}

func TestQueryExpectingNoPhotos(t *testing.T) {
	s, err := newBoltStorage()
	if err != nil {
		t.Fatalf("Can't create Bolt storage!")
	}
	defer s.Close()

	photos := s.Search("%picture_me%")
	if len(photos) != 0 {
		t.Fatalf("Expecting zero results but got %d results!", len(photos))
	}
}
