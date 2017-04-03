package main

import (
	log "github.com/Sirupsen/logrus"
	"reflect"
)

type StorageBackend interface {
	Has(build BuildFile) bool
	Get(build BuildFile)
	Put(build BuildFile)
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

func (svc *StorageBag) Has(build BuildFile) bool {
	for _, backend := range svc.backend {
		if backend.Has(build) {
			return true
		}
	}

	return false
}

func (svc *StorageBag) Get(build BuildFile) {
	for _, backend := range svc.backend {
		backend.Get(build)
	}
}

func (svc *StorageBag) Put(build BuildFile) {
	for _, backend := range svc.backend {
		backend.Put(build)
	}
}
