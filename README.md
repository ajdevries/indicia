# Indicia

Small [go](https://golang.org) program to extract [EXIF](https://en.wikipedia.org/wiki/Exif) data from images. Build for [Waldo Photos](http://waldo.photos/).

# Structure

### Reader
Interface for reading and parsing EXIF tags from an image file. Implementation depends on what you need. For now
there is a S3Reader that can read images from a relative URL, given the base URL.

Usage:

```go
reader := newS3Reader("http://s3.amazonaws.com/waldo-recruiting")
tags, err := reader.ReadAndParse("0003b8d6-d2d8-4436-a398-eab8d696f0f9.68cccdd4-e431-457d-8812-99ab561bf867.jpg")
```

### Lister
Interface for listing images from a base location. Could be an URL or a file path. For now there is the
S3Lister implementation expects the ListBucketResult XML format and parses the Key tag from the Content tags.

Usage:

```go
lister := newS3Lister("http://s3.amazonaws.com/waldo-recruiting")
urls, err := lister.List()
```

### Storage
Interface for storing images and there EXIF data. First implementation is the `BoltStorage` implementation. It
is possible to Save and Search for images. With the Search method a SQL like syntax is possible, i.e. `%test%` returns all
photos that contain the keyword test in the file name.

Usage:

```go
s, err := newBoltStorage()
if err != nil {
  t.Fatalf("Can't create Bolt storage!")
}
defer s.Close()

s.Save("test.jpg", map[string]string{"Make": "Apple"})
photos := s.Search("test%") // returns a slice containing one pointer to a Photo struct
```

### Indicia
The component that glues every thing together. First it lists all the photos from a `Lister` component (S3 bucket). Then it reads and parses the EXIF tags using a `Reader`,
and then is stores the photo and their tags in to a `Storage` engine, i.e. `BoltStorage`.
# How to build
