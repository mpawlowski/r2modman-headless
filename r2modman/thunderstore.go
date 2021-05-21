package r2modman

import (
	"archive/zip"
	"strings"
)

type ModPackagingType int

const (
	ModPackagingTypeUnknown ModPackagingType = iota
	ModPackagingTypeRoot
	ModPackagingTypePlugin
)

func (d ModPackagingType) String() string {
	return [...]string{"ModPackagingTypeUnknown", "ModPackagingTypeRoot", "ModPackagingTypePlugin"}[d]
}

func (d ModPackagingType) Directory() string {
	return [...]string{"", "", "/BepInEx/plugins"}[d]
}

func DeterminePackagingType(modZipFilename string) (ModPackagingType, error) {

	// open zip file
	r, err := zip.OpenReader(modZipFilename)
	if err != nil {
		return ModPackagingTypeUnknown, err
	}
	defer r.Close()

	// loop through all the files in the archive to check if BepInEx/ exists
	for _, f := range r.File {
		if strings.Contains(f.Name, "BepInEx/") {
			return ModPackagingTypeRoot, nil
		}
	}

	// assume Plugin by default
	return ModPackagingTypePlugin, nil
}
