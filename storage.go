package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"log"
)

const S3_BUCKET = "cmsbuild"
const S3_REGION = "eu-central-1"

func Has(build BuildFile) bool {

	/*if _, err := os.Stat(build.Archive); err == nil {
		return true;
	}

	return false;*/

	fmt.Println("Checking storage for hash", build.Hash)

	svc := s3.New(session.New(&aws.Config{Region: aws.String(S3_REGION)}))

	params := &s3.HeadObjectInput{
		Bucket:	aws.String(S3_BUCKET),
		Key:    aws.String(build.Hash),
	}

	resp, err := svc.HeadObject(params)

	if err != nil {
		fmt.Println("Build not found")

		return false
	}

	fmt.Println("Build found", resp)

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

	fmt.Println("Upload archive", build.Archive, build.Hash)

	file, err := os.Open(build.Archive)

	if err != nil {
		log.Fatal("Failed to open file", err)
	}

	defer file.Close()

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(S3_REGION)}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: 	aws.String(S3_BUCKET),
		Key:    	aws.String(build.Hash),
		Body:   	file,
		/*Metadata:	map[string]*string{
			"Name":	aws.String(build.Name),
		},*/
	})

	if err != nil {
		log.Fatalln("Failed to upload", err)
	}

	log.Println("Successfully uploaded to", result.Location)
}
