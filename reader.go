package main

import (
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"io"
	"net/http"
	"net/url"
	"path"
)

// S3Reader can read images from an HTTP server (S3 in this case) given a
// base URL.
type S3Reader struct {
	base string
	tags map[string]string
}

// Implementation of the Reader interface. When URL can't be found or image is invalid
// error is returned
func (s *S3Reader) ReadAndParse(name string) (tags map[string]string, err error) {
	u, err := url.Parse(s.base)
	if err != nil {
		return
	}

	u.Path = path.Join(u.Path, name)
	resp, err := http.Get(u.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	tags, err = s.parseExif(resp.Body)
	return
}

// Private method for parsing the EXIF fields based on a the io.Reader interface
func (f *S3Reader) parseExif(reader io.Reader) (map[string]string, error) {
	e, err := exif.Decode(reader)
	if err != nil {
		return nil, err
	}
	f.tags = make(map[string]string)
	e.Walk(f)
	return f.tags, nil
}

// Implementation of exif.Walker interface, is called for every EXIF field that is parsed
// and this fieldname and value (as string) is placed in the tags member
// of the S3Reader
func (f *S3Reader) Walk(name exif.FieldName, tag *tiff.Tag) error {
	s, _ := tag.StringVal()
	if s != "" {
		f.tags[string(name)] = s
	}
	return nil
}

// Creates a new S3 Reader that implements the Reader interface for parsing images
func newS3Reader(base string) Reader {
	return &S3Reader{base: base}
}

// Reader, reads and parses a image file based on the given name. Returns a map with the EXIF tags, for easiness
// tags are formatted as strings. When an image can't be found or parsed error is returned
type Reader interface {
	ReadAndParse(name string) (map[string]string, error)
}
