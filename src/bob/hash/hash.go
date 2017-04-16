package hash

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"sort"
	"strings"
)

func Value(key string, value string) (string, error) {
	h := md5.New()
	io.WriteString(h, key)
	io.WriteString(h, value)

	return hex.EncodeToString(h.Sum(nil)), nil
}

func File(filePath string) (string, error) {
	var result []byte

	file, err := os.Open(filePath)

	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(result)), nil
}

func List(hashes map[string]string) string {
	list := make([]string, 0, len(hashes))

	for _, value := range hashes {
		if len(value) > 0 {
			list = append(list, value)
		}
	}

	sort.Strings(list)
	joinedString := strings.Join(list, ",")
	hash := md5.New()

	io.WriteString(hash, joinedString)

	return hex.EncodeToString(hash.Sum(nil))
}
