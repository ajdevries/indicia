package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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

type FileLister struct {
	S3Lister
	root string
}

func (f *FileLister) List() (result []string, err error) {
	result = []string{}
	count:=0
	filepath.Walk(f.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			fmt.Printf("\r%d", count)
			count+=1
			if r, err := regexp.MatchString("\\.jpg|\\.png", strings.ToLower(info.Name())); err == nil && r {
				result = append(result, path[len(f.root)+1:len(path)])
			}
		}
		return nil
	})
	return result, nil
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

func newFileLister(root string) Lister {
	return &FileLister{root: root}
}

// Interface that returns a list of strings containing image names that can be
// read using the S3Reader implementation
type Lister interface {
	List() ([]string, error)
}
