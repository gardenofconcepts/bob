package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"time"
)

type Session struct {
	instance *session.Session
	bucket   string
}

func Storage(region, bucket string) *Session {
	return &Session{
		instance: session.New(&aws.Config{Region: aws.String(region)}),
		bucket:   bucket,
	}
}

func (svc *Session) Has(build BuildFile) bool {

	log.Info("Checking storage for hash", build.Hash)

	params := &s3.HeadObjectInput{
		Bucket: aws.String(svc.bucket),
		Key:    aws.String(build.Hash),
	}

	resp, err := s3.New(svc.instance).HeadObject(params)

	if err != nil {
		log.Warning("Build not found")

		return false
	}

	log.Info("Build found", resp)

	return true
}

func (svc *Session) Get(build BuildFile) {

	file, err := os.Create(build.Archive)

	if err != nil {
		log.Fatal("Failed to create file", err)
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(svc.instance)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(svc.bucket),
			Key:    aws.String(build.Hash),
		})

	if err != nil {
		log.Fatal("Cannot download build", err)
	}

	log.Info("Downloaded file", file.Name(), numBytes, "bytes")
}

func (svc *Session) Put(build BuildFile) {

	log.Info("Upload archive", build.Archive, build.Hash)

	file, err := os.Open(build.Archive)

	if err != nil {
		log.Fatal("Failed to open file", err)
	}

	defer file.Close()

	uploader := s3manager.NewUploader(svc.instance)
	hostname, _ := os.Hostname()
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(svc.bucket),
		Key:    aws.String(build.Hash),
		Body:   file,
		Metadata: map[string]*string{
			"Name":      aws.String(build.Name),
			"Creator":   aws.String(hostname),
			"CreatedAt": aws.String(time.Now().Format(time.RFC850)),
		},
	})

	if err != nil {
		log.Fatal("Cannot upload build", err)
	}

	log.Println("Successfully uploaded to", result.Location)
}
