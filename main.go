package main

import (
	"flag"
)

var (
	currentVersion  = "1.0"
	bucketURL       = flag.String("url", "http://s3.amazonaws.com/waldo-recruiting", "URL of the S3 bucket")
	numberOfReaders = flag.Int("numberOfReaders", 4, "Number of concurrent readers")
)

func main() {
	flag.Parse()
	s, _ := newBoltStorage()
	defer s.Close()
	i := newIndicia(*bucketURL, newS3Lister, newS3Reader, s, *numberOfReaders)
	go func() {
		i.Start()
	}()
	StartServer(i)
}
