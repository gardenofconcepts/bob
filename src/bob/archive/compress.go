package archive

import (
	"os"
	"archive/tar"
	log "github.com/Sirupsen/logrus"
	"io"
	"compress/gzip"
	"path/filepath"
	"errors"
	"strings"
	"bob/path"
)

func (archive *Archive) Compress(baseDir string, includes []string, excludes []string) ([]string, error) {
	log.WithFields(log.Fields{
		"cwd":      baseDir,
		"file":     archive.path,
		"includes": includes,
		"excludes": excludes,
	}).Info("Packaging files")

	fileList := []string{}

	file, err := os.Create(archive.path)
	if err != nil {
		return fileList, err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	err = filepath.Walk(baseDir, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			log.Error("Error reading directory", err)

			return err
		}

		if !path.MatchList(includes, filePath, baseDir) || path.MatchList(excludes, filePath, baseDir) {
			log.WithField("file", filePath).Debug("Skipping file")

			return nil
		}

		if err := addNode(baseDir, filePath, fileInfo, tarWriter); err != nil {
			return err
		}

		fileList = append(fileList, filePath)

		return nil
	})

	if err != nil {
		log.WithFields(log.Fields{
			"cwd":   baseDir,
			"file":  archive.path,
			"error": err,
		}).Fatal("Error creating archive")
	}

	if len(fileList) == 0 {
		log.Fatal("Archive contains no files")

		errors.New("Archive contains no files")
	}

	return fileList, nil
}

func Compress(baseDir string, includes []string, excludes []string) ([]string, error) {
	log.WithFields(log.Fields{
		"cwd":      baseDir,
		"includes": includes,
		"excludes": excludes,
	}).Info("Packaging files")

	fileList := []string{}

	err := filepath.Walk(baseDir, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			log.Error("Error reading directory", err)

			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		if !path.MatchList(includes, filePath, baseDir) || path.MatchList(excludes, filePath, baseDir) {
			log.WithField("file", filePath).Debug("Skipping file")

			return nil
		}

		filePath, _ = filepath.Rel(baseDir, filePath)
		fileList = append(fileList, filePath)

		return nil
	})

	if err != nil {
		log.WithFields(log.Fields{
			"cwd": baseDir,
		}).Fatal("Error creating archive")
	}

	if len(fileList) == 0 {
		log.Fatal("Archive contains no files")

		errors.New("Archive contains no files")
	}

	return fileList, nil
}

func addNode(baseDir string, filePath string, fileInfo os.FileInfo, tarWriter *tar.Writer) error {
	log.WithFields(log.Fields{
		"cwd":  baseDir,
		"file": filePath,
	}).Debug("Adding file")

	header, err := tar.FileInfoHeader(fileInfo, filePath)

	if err != nil {
		log.Error("Cannot read header", err)

		return err
	}

	// To store relative paths
	header.Name = strings.TrimPrefix(filePath, baseDir+string(filepath.Separator))

	switch header.Typeflag {
	case tar.TypeDir:
		return addDirectory(filePath, fileInfo, tarWriter, header)
	case tar.TypeReg, tar.TypeRegA:
		return addFile(filePath, fileInfo, tarWriter, header)
	case tar.TypeSymlink:
		return addSymbolicLink(filePath, fileInfo, tarWriter, header)
	case tar.TypeLink:
		return addHardLink(filePath, fileInfo, tarWriter, header)
	default:
		log.Fatal("Error while adding; Unkown type: ", string(header.Typeflag))
	}

	return nil
}

func addHeader(tarWriter *tar.Writer, header *tar.Header) error {
	if err := tarWriter.WriteHeader(header); err != nil {
		log.WithFields(log.Fields{
			"header": header,
			"error":  err,
		}).Error("Cannot write header")

		return err
	}

	return nil
}

func addSymbolicLink(filePath string, fileInfo os.FileInfo, tarWriter *tar.Writer, header *tar.Header) error {
	symlinkTarget, _ := os.Readlink(filePath)

	header.Linkname = symlinkTarget

	return addHeader(tarWriter, header)
}

func addDirectory(filePath string, fileInfo os.FileInfo, tarWriter *tar.Writer, header *tar.Header) error {
	return nil
}

func addHardLink(filePath string, fileInfo os.FileInfo, tarWriter *tar.Writer, header *tar.Header) error {
	return addHeader(tarWriter, header)
}

func addFile(filePath string, fileInfo os.FileInfo, tarWriter *tar.Writer, header *tar.Header) error {
	if err := addHeader(tarWriter, header); err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  filePath,
			"error": err,
		}).Error("Cannot open file")

		return err
	}
	defer file.Close()

	if _, err = io.CopyN(tarWriter, file, fileInfo.Size()); err != nil {
		log.WithFields(log.Fields{
			"file":  filePath,
			"error": err,
		}).Error("Cannot copy file")

		return err
	}

	return nil
}
