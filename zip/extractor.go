package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type Extractor interface {
	Extract(zipFile string, destinationDir string, stripPrefix string) error
}

type extractorImpl struct {
}

func ExtractAndWriteFile(zipEntry *zip.File, normalizedDestination string) error {
	if zipEntry.FileInfo().IsDir() {
		return nil
	}

	zipReadCloser, err := zipEntry.Open()
	if err != nil {
		return err
	}
	defer zipReadCloser.Close()

	// Create the directory path if it doesn't exist
	dirPath := filepath.Dir(normalizedDestination)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("unable to create directory '%s': %v", dirPath, err)
	}

	// Open the file in write-only mode, create if it doesn't exist, truncate otherwise
	file, err := os.OpenFile(normalizedDestination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("unable to open file '%s': %v", normalizedDestination, err)
	}
	defer file.Close()

	log.Printf("extracting zipped file: %s -> %s", zipEntry.Name, normalizedDestination)

	_, err = io.Copy(file, zipReadCloser)
	if err != nil {
		return err
	}

	return nil
}

func (e *extractorImpl) Extract(zipFileName string, destinationDir string, stripPrefix string) error {

	r, err := zip.OpenReader(zipFileName)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {

		var normalizedFilePath string

		if runtime.GOOS == "windows" {
			// Convert to Windows path format
			normalizedFilePath = strings.ReplaceAll(f.Name, "/", "\\")
		} else {
			// Convert to Unix path format
			normalizedFilePath = strings.ReplaceAll(f.Name, "\\", "/")
		}

		normalizedFilePath, _ = strings.CutPrefix(normalizedFilePath, stripPrefix)

		destination := path.Join(destinationDir, normalizedFilePath)

		err := ExtractAndWriteFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func newExtractor() Extractor {
	return &extractorImpl{}
}
