package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	log.WithFields(log.Fields{
		"cwd":      baseDir,
		"file":     archive.path,
		"includes": includes,
		"excludes": excludes,
	}).Info("Packaging files")

	files := 0

	file, err := os.Create(archive.path)
	if err != nil {
		return err
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

		if !matchList(includes, filePath, baseDir) || matchList(excludes, filePath, baseDir) {
			log.WithField("file", filePath).Debug("Skipping file")

			return nil
		}

		if err := addNode(baseDir, filePath, fileInfo, tarWriter); err != nil {
			return err
		}

		files++

		return nil
	})

	if err != nil {
		log.WithFields(log.Fields{
			"cwd":   baseDir,
			"file":  archive.path,
			"error": err,
		}).Fatal("Error creating archive")
	}

	if files == 0 {
		log.Fatal("Archive contains no files")

		errors.New("Archive contains no files")
	}

	return nil
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

func (archive *Archive) Extract(dest string) error {
	log.WithFields(log.Fields{
		"file": archive.path,
		"cwd":  dest,
	}).Info("Extracting archive")

	file, err := os.Open(archive.path)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  archive.path,
			"error": err,
		}).Error("Error while opening archive")

		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  archive.path,
			"error": err,
		}).Error("Error while reading archive")

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

		if err = extractTarArchiveFile(header, dest, tarReader); err != nil {
			return err
		}
	}

	return nil
}

func extractTarArchiveFile(header *tar.Header, dest string, input io.Reader) error {
	switch header.Typeflag {
	case tar.TypeDir:
		return writeNewDirectory(filepath.Join(dest, header.Name))
	case tar.TypeReg, tar.TypeRegA:
		return writeNewFile(filepath.Join(dest, header.Name), input, header.FileInfo().Mode().Perm())
	case tar.TypeSymlink:
		return writeNewSymbolicLink(filepath.Join(dest, header.Name), header.Linkname)
	case tar.TypeLink:
		return writeNewHardLink(filepath.Join(dest, header.Name), filepath.Join(dest, header.Linkname))
	default:
		log.WithFields(log.Fields{
			"file": header.Name,
			"type": header.Typeflag,
		}).Error("Unknown type flag")
	}

	return nil
}

func writeNewFile(filePath string, in io.Reader, fm os.FileMode) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  filePath,
			"error": err,
		}).Error("Making directory for file")

		return err
	}

	out, err := os.Create(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  filePath,
			"error": err,
		}).Error("Creating new file")

		return err
	}
	defer out.Close()

	err = out.Chmod(fm)
	if err != nil && runtime.GOOS != "windows" {
		log.WithFields(log.Fields{
			"file":  filePath,
			"error": err,
		}).Error("Changing file mode")

		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  filePath,
			"error": err,
		}).Error("Writing file", filePath, err)

		return err
	}

	return nil
}

func writeNewSymbolicLink(filePath string, target string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  filePath,
			"error": err,
		}).Error("Making directory for file", filePath, err)

		return err
	}

	err = os.Symlink(target, filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"file":   filePath,
			"target": target,
			"error":  err,
		}).Error("Making symbolic link")

		return err
	}

	return nil
}

func writeNewHardLink(filePath string, target string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  filePath,
			"error": err,
		}).Error("Making directory for file")

		return err
	}

	err = os.Link(target, filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  filePath,
			"error": err,
		}).Error("Making hard link", filePath, err)

		return err
	}

	return nil
}

func writeNewDirectory(dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  dirPath,
			"error": err,
		}).Error("Making directory")

		return err
	}

	return nil
}
