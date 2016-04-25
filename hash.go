package main

import (
	"os"
	"io"
	"strings"
	"crypto/md5"
	"encoding/hex"
)

func hashFile(filePath string) (string, error) {
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

func hashList(hashes []string) (string, error) {
	joinedString := strings.Join(hashes, ",")
	hash := md5.New()

	io.WriteString(hash, joinedString)

	return hex.EncodeToString(hash.Sum(nil)), nil
}