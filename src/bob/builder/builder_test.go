package builder

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
	"bob/util"
)

func TestBuilder(t *testing.T) {
	path := filepath.Join("../../../../build", util.RandomString(10))

	Builder().Build("test-fixtures", []Build{
		{
			Command: "cat data.txt > " + path,
		},
	})

	buf := bytes.NewBuffer(nil)
	f, err := os.Open(filepath.Join("../../../build", filepath.Base(path)))

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
