# Indicia

Small [go](https://golang.org) program to extract [EXIF](https://en.wikipedia.org/wiki/Exif) data from images. Build for [Waldo Photos](http://waldo.photos/).

# Structure

### Reader
Interface for reading and parsing EXIF tags from an image file. Implementation depends on what you need. For now
there is a S3Reader that can read images from a relative URL, given the base URL. I.e.
```go
reader := newS3Reader("http://s3.amazonaws.com/waldo-recruiting")
tags, err := reader.ReadAndParse("0003b8d6-d2d8-4436-a398-eab8d696f0f9.68cccdd4-e431-457d-8812-99ab561bf867.jpg")
```

# How to build
