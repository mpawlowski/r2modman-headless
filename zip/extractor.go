package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Extractor interface {
	Extract(zipFile string, destinationDir string, stripPrefix string) error
}

type extractorImpl struct {
}

func (e *extractorImpl) Extract(zipFileName string, destinationDir string, stripPrefix string) error {

	r, err := zip.OpenReader(zipFileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(zipFile *zip.File, stripPrefix string) error {
		zipReadCloser, err := zipFile.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := zipReadCloser.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(destinationDir, zipFile.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(destinationDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if stripPrefix != "" {
			stripped := strings.Replace(zipFile.Name, stripPrefix, "", 1)
			path = filepath.Join(destinationDir, stripped)
			// log.Println("wat1", stripPrefix, path)
		}

		if zipFile.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {

			err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
			if err != nil {
				return err
			}

			// log.Println("wat2", stripPrefix, path, f.Name())

			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			// log.Printf("extracting zipped file: %s -> %s", zipFile.Name, f.Name())
			_, err = io.Copy(f, zipReadCloser)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {

		// stripPrefix := getPrefixToStrip(f)
		// if stripPrefix != "" {
		// 	log.Printf("stripping prefix '%s' from %s\n", stripPrefix, f.Name)
		// }

		err := extractAndWriteFile(f, stripPrefix)
		if err != nil {
			return err
		}
	}

	return nil
}

func newExtractor() Extractor {
	return &extractorImpl{}
}
