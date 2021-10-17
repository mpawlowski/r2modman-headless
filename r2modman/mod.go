package r2modman

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type ModUtil interface {
	Download(mod ExportR2xMod) error
}

type modUtilImpl struct {
	config Config
}

func (m *modUtilImpl) Download(mod ExportR2xMod) error {

	client := http.Client{
		Timeout: m.config.ThunderstoreCDNTimeout,
	}

	// Get the data
	resp, err := client.Get(mod.DownloadUrl())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fmt.Sprintf("%s/%s", m.config.WorkDirectory, mod.Filename()))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func newModUtil(
	config Config,
) (ModUtil, error) {
	return &modUtilImpl{
		config: config,
	}, nil
}
