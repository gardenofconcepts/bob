package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"log"
	"io"
	"compress/gzip"
)

const S3_BUCKET = "cms"
const S3_REGION = "us-west-2"

func Has(build BuildFile) bool {

	fmt.Println("Checking storage for hash", build.Hash)

	svc := s3.New(session.New(&aws.Config{Region: aws.String(S3_REGION)}))

	params := &s3.HeadObjectInput{
		Bucket:	aws.String(S3_BUCKET),
		Key:    aws.String(build.Hash),
	}

	resp, err := svc.HeadObject(params)

	if err != nil {
		return false
	}

	fmt.Println(resp)

	return true;
}

func Get(build BuildFile) {

	file, err := os.Create(build.Archive)

	if err != nil {
		log.Fatal("Failed to create file", err)
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String(S3_REGION)}))
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(S3_BUCKET),
			Key:    aws.String(build.Hash),
		})

	if err != nil {
		fmt.Println("Failed to download file", err)
		return
	}

	fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")
}

func Put(build BuildFile) {

	fmt.Println("Upload archive", build.Archive)

	file, err := os.Open(build.Archive)

	if err != nil {
		log.Fatal("Failed to open file", err)
	}

	// Not required, but you could zip the file before uploading it
	// using io.Pipe read/writer to stream gzip'd file contents.
	reader, writer := io.Pipe()

	go func() {
		gw := gzip.NewWriter(writer)
		io.Copy(gw, file)

		file.Close()
		gw.Close()
		writer.Close()
	}()

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(S3_REGION)}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   	reader,
		Bucket: 	aws.String(S3_BUCKET),
		Key:    	aws.String(build.Hash),
		Metadata:	map[string]*string{
			"Name":	aws.String(build.Name),
		},
	})

	if err != nil {
		log.Fatalln("Failed to upload", err)
	}

	log.Println("Successfully uploaded to", result.Location)
}
