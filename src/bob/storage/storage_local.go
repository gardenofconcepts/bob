package storage

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
)

type StorageLocalBackend struct {
	StorageBackend

	path string
}

func StorageLocal(path string) *StorageLocalBackend {
	return &StorageLocalBackend{
		path: path,
	}
}

func (svc *StorageLocalBackend) Has(build StorageRequest) bool {
	if len(build.Archive) == 0 {
		build.Archive = path.Join(svc.path, build.Hash+".tar.gz")
	}

	file := filepath.Join(svc.path, filepath.Base(build.Archive))

	log.WithFields(log.Fields{
		"file": file,
	}).Debug("Check for local file")

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	return true
}

func (svc *StorageLocalBackend) Get(build StorageRequest) {
	log.Debug("Expect there is a local file")
}

func (svc *StorageLocalBackend) Put(build StorageRequest) {
	log.Debug("Uploading file ;-D")
}
