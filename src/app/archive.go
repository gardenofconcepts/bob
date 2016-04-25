package app

import (
	"io"
	"archive/tar"
	"os"
	"path/filepath"
	"strings"
)

func Archive(path string, baseDir string, include []string, exclude []string) error {

	tarfile, err := os.Create(path)

	if err != nil {
		return err
	}

	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())

		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = strings.TrimPrefix(path, baseDir)
		}

		if err := tarball.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)

		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(tarball, file)

		return err
	})

	/*
	if match(include, path) && !match(exclude, path) {
		if err := addFile(tw, path); err != nil {
			log.Fatalln(err)
		}
	}
	*/

	return nil
}

func Extract(path string, directory string) {

}
