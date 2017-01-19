package main

import (
	"path/filepath"
	"testing"
	"io"
	"bytes"
	"os"
)

func TestBuilder(t *testing.T) {
	path := filepath.Join("../../build", randomString(10))

	Builder().Build("assets/build", []Build{
		{
			Command: "cat data.txt > " + path,
		},
	})

	buf := bytes.NewBuffer(nil)
	f, err := os.Open(filepath.Join("build", filepath.Base(path)))

	if err != nil {
		t.Errorf("Error while reading file %s: %s", path, err)

		return
	}

	io.Copy(buf, f)
	f.Close()

	content := string(buf.Bytes())

	if content != "this is a test file" {
		t.Errorf("Expect content of 'this is a test file' in file %s, got %s", path, content)
	}
}
