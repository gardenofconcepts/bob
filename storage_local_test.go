package main

import (
	"os"
	"testing"
)

func TestStorageLocal(t *testing.T) {
	storage := StorageLocal(os.TempDir())

	result := storage.Has(BuildFile{
		Hash: "123",
	})

	if result == true {
		t.Errorf("Expect %s, got %s", "false", true)
	}
}
