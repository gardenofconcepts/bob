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
	tarfile, err := os.Create(archive.path)

	if err != nil {
		return err
	}

	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)

	defer tarball.Close()

	filepath.Walk(baseDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
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
			header.Name = strings.TrimPrefix(filePath, baseDir)
		}

		if err := tarball.WriteHeader(header); err != nil {
			return err
		}

		file, err := os.Open(filePath)

		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(tarball, file)

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

	fd, err := os.Open(archive.path)
	if err != nil {
		return err
	}
	defer fd.Close()

	gReader, err := gzip.NewReader(fd)
	if err != nil {
		return err
	}
	defer gReader.Close()

	tarReader := tar.NewReader(gReader)

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if hdr.Name == "." {
			continue
		}

		err = extractTarArchiveFile(hdr, dest, tarReader)
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
