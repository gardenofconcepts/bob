package archive

import (
	"os"
	log "github.com/Sirupsen/logrus"
	"path/filepath"
	"compress/gzip"
	"archive/tar"
	"io"
	"runtime"
)

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

	if _, err := os.Stat(filePath); err == nil {
		os.RemoveAll(filePath)
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
