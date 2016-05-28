package main

import (
	"archive/tar"
	"compress/gzip"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"errors"
)

type Archive struct {
	path string
}

func NewArchive(path string) *Archive {
	return &Archive{
		path: path,
	}
}

func (archive *Archive) Compress(baseDir string, includes []string, excludes []string) error {

	log.Info("Archive files", archive.path, baseDir)

	files :=0

	file, err := os.Create(archive.path)

	if err != nil {
		return err
	}

	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	filepath.Walk(baseDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !matchList(includes, filePath, baseDir) || matchList(excludes, filePath, baseDir) {
			return nil
		}

		log.Debug("Add file to archive", filePath)

		header, err := tar.FileInfoHeader(info, info.Name())

		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(filePath, baseDir))
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(filePath)

		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(tarWriter, file)

		files++

		return err
	})

	if files == 0 {
		log.Fatal("Archive contains no files")

		errors.New("Archive contains no files")
	}

	return nil
}

func (archive *Archive) Extract(dest string) error {
	log.Info("Extract file to directory", archive.path, dest)

	file, err := os.Open(archive.path)

	if err != nil {
		return err
	}

	defer file.Close()

	gzipReader, err := gzip.NewReader(file)

	if err != nil {
		return err
	}

	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if header.Name == "." {
			continue
		}

		err = extractTarArchiveFile(header, dest, tarReader)

		if err != nil {
			return err
		}
	}

	return nil
}

func extractTarArchiveFile(header *tar.Header, dest string, input io.Reader) error {
	filePath := filepath.Join(dest, header.Name)
	fileInfo := header.FileInfo()

	if fileInfo.IsDir() {
		err := os.MkdirAll(filePath, fileInfo.Mode())

		if err != nil {
			return err
		}
	} else {
		err := os.MkdirAll(filepath.Dir(filePath), 0755)

		if err != nil {
			return err
		}

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			return os.Symlink(header.Linkname, filePath)
		}

		fileCopy, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())

		if err != nil {
			return err
		}

		defer fileCopy.Close()

		_, err = io.Copy(fileCopy, input)

		if err != nil {
			return err
		}
	}

	return nil
}
