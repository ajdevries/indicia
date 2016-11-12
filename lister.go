package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

// Small struct for parsing the ListBuckerResult XML format that the S3 bucker returns
type ListBucketResult struct {
	Contents []struct {
		Key string
	}
}

// Implementation of Lister interface, specialized for (public) S3 buckets
type S3Lister struct {
	url string
}

// Implementation of the Lister interface, expects the ListBuckerResult XML format
// and parses the Key tag from the Content tags
// When something fails error is returned
func (s *S3Lister) List() (result []string, err error) {
	result = []string{}

	resp, err := http.Get(s.url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	v := ListBucketResult{}
	err = xml.Unmarshal(bytes, &v)
	if err != nil {
		return
	}

	for _, key := range v.Contents {
		result = append(result, strings.Trim(key.Key, " "))
	}
	return
}

func newS3Lister(url string) Lister {
	return &S3Lister{url: url}
}

// Interface that returns a list of strings containing image names that can be
// read using the S3Reader implementation
type Lister interface {
	List() ([]string, error)
}
