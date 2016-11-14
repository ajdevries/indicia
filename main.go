package main

import (
	"flag"
)

var (
	currentVersion  = "1.0"
	bucketURL       = flag.String("url", "/Photos", "File location of the photos")
	numberOfReaders = flag.Int("numberOfReaders", 16, "Number of concurrent readers")
)

func main() {
	flag.Parse()
	s, _ := newBoltStorage()
	defer s.Close()
	i := newIndicia(*bucketURL, newFileLister, newFileReader, s, *numberOfReaders)
	go func() {
		i.Start()
	}()
	StartServer(i)
}
