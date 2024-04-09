// Package r2modman contains utilities for interacting with r2modman/thunderstore
package r2modman

import "fmt"

const thunderstoreApiUrlFormat = "https://%s/live/repository/packages/%s-%d.%d.%d.zip"
const filenameFormat = "%s-%d.%d.%d.zip"

type ExportR2xModVersion struct {
	Major int
	Minor int
	Patch int
}

type ExportR2xMod struct {
	Name    string
	Version ExportR2xModVersion
	Enabled bool
}

func (e ExportR2xMod) Filename() string {
	return fmt.Sprintf(filenameFormat, e.Name, e.Version.Major, e.Version.Minor, e.Version.Patch)
}

func (e *ExportR2xMod) DownloadUrl(cdnHost string) string {
	return fmt.Sprintf(thunderstoreApiUrlFormat, cdnHost, e.Name, e.Version.Major, e.Version.Minor, e.Version.Patch)
}

func (e *ExportR2xMod) ThunderstoreKey() string {
	return e.Name
}

func (e *ExportR2xMod) ThunderstoreModVersion() string {
	return fmt.Sprintf("%s-%d.%d.%d", e.Name, e.Version.Major, e.Version.Minor, e.Version.Patch)
}

type ExportR2x struct {
	ProfileName string `yaml:"profileName"`
	Mods        []ExportR2xMod
}
