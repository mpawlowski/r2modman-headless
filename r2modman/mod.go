package r2modman

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ModUtil interface {
	Download(mod ExportR2xMod, metadata *APIPackageResponse) error
}

type modUtilImpl struct {
	config Config
}

func (m *modUtilImpl) verifyIntegrity(mod ExportR2xMod, metadata *APIPackageResponse) error {
	downloadedZipPath := m.config.WorkDirectory + "/" + mod.Filename()
	existingModFile, err := os.Stat(downloadedZipPath)
	if err != nil {
		return fmt.Errorf("mod package %s did not exist", downloadedZipPath)
	}
	for _, thunderstoreVersion := range metadata.Versions {
		if thunderstoreVersion.FullName == mod.ThunderstoreModVersion() {
			log.Printf("Verifying integrity: %s", downloadedZipPath)
			if thunderstoreVersion.FileSize == existingModFile.Size() {
				return nil
			}
			log.Printf("File Integrity Failed: %s, expected size: %d bytes, actual size: %d bytes", mod.Filename(), thunderstoreVersion.FileSize, existingModFile.Size())

		}
	}

	//package did not validate, delete it
	log.Printf("Removing invalid file: %s", downloadedZipPath)
	err = os.Remove(downloadedZipPath)
	if err != nil {
		return err
	}

	return fmt.Errorf("mod package %s did not validate", downloadedZipPath)
}

func (m *modUtilImpl) Download(mod ExportR2xMod, metadata *APIPackageResponse) error {
	downloadedZipPath := m.config.WorkDirectory + "/" + mod.Filename()

	client := http.Client{
		Timeout: m.config.ThunderstoreCDNTimeout,
	}

	err := m.verifyIntegrity(mod, metadata)
	if err == nil && !m.config.ThunderstoreForceDownload {
		return nil // file exists and validates
	}

	log.Printf("Downloading mod: %s", downloadedZipPath)
	resp, err := client.Get(mod.DownloadUrl(m.config.ThunderstoreCDN))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("unable to download, HTTP Error %d", resp.StatusCode)
	}

	// Create the file
	out, err := os.Create(downloadedZipPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return m.verifyIntegrity(mod, metadata)
}

func newModUtil(
	config Config,
) (ModUtil, error) {
	return &modUtilImpl{
		config: config,
	}, nil
}
