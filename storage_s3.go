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

type StorageS3Backend struct {
	StorageBackend

	instance *session.Session
	bucket   string
}

func StorageS3(region, bucket string) *StorageS3Backend {
	return &StorageS3Backend{
		instance: session.New(&aws.Config{Region: aws.String(region)}),
		bucket:   bucket,
	}
}

func (svc *StorageS3Backend) Has(build BuildFile) bool {
	log.WithFields(log.Fields{
		"bucket": svc.bucket,
		"hash":   build.Hash,
	}).Info("Checking storage for existing build")

	params := &s3.HeadObjectInput{
		Bucket: aws.String(svc.bucket),
		Key:    aws.String(build.Hash),
	}

	resp, err := s3.New(svc.instance).HeadObject(params)

	if err != nil {
		log.WithFields(log.Fields{
			"bucket": svc.bucket,
			"hash":   build.Hash,
		}).Warning("Build not found")

		return false
	}

	log.WithFields(log.Fields{
		"data": resp,
	}).Info("Found existing build")

	return true
}

func (svc *StorageS3Backend) Get(build BuildFile) {

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

	log.WithFields(log.Fields{
		"file":  file.Name(),
		"bytes": numBytes,
	}).Info("Downloaded file")
}

func (svc *StorageS3Backend) Put(build BuildFile) {

	log.WithFields(log.Fields{
		"file": build.Archive,
		"hash": build.Hash,
	}).Info("Uploading archive")

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

	log.WithFields(log.Fields{
		"location": result.Location,
	}).Info("Successfully uploaded to")
}
