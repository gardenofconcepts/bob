package storage

import (
	log "github.com/Sirupsen/logrus"
	"reflect"
)

type StorageBackend interface {
	Has(build StorageRequest) bool
	Get(build StorageRequest)
	Put(build StorageRequest)
}

type StorageRequest struct {
	Name	string
	Archive string
	Hash    string
}

type StorageBag struct {
	backend []StorageBackend
}

func Storage() StorageBag {
	return StorageBag{
		backend: []StorageBackend{},
	}
}

func (svc *StorageBag) Register(backend StorageBackend) {
	svc.backend = append(svc.backend, backend)

	log.WithField("type", reflect.TypeOf(backend)).Debug("Register backend")
}

func (svc *StorageBag) Has(build StorageRequest) bool {
	for _, backend := range svc.backend {
		if backend.Has(build) {
			return true
		}
	}

	return false
}

func (svc *StorageBag) Get(build StorageRequest) {
	for _, backend := range svc.backend {
		backend.Get(build)
	}
}

func (svc *StorageBag) Put(build StorageRequest) {
	for _, backend := range svc.backend {
		backend.Put(build)
	}
}
