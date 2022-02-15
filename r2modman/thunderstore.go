package r2modman

import (
	"archive/zip"
	"log"
	"strings"
)

type ModPackagingType int

const (
	// ModPackagingTypeUnknown is the default packaging type. This means the packaging is unrecognized
	ModPackagingTypeUnknown ModPackagingType = iota
	ModPackagingTypeRoot
	ModPackagingTypePlugin
	ModPackagingTypeRootDLL
)

func (d ModPackagingType) String() string {
	return [...]string{"ModPackagingTypeUnknown", "ModPackagingTypeRoot", "ModPackagingTypePlugin", "ModPackagingTypeRootDLL"}[d]
}

func (d ModPackagingType) Directory() string {
	return [...]string{"", "", "BepInEx", "BepInEx/plugins"}[d]
}

func DeterminePackagingType(modZipFilename string) (modType ModPackagingType, prefixToStrip string, err error) {

	// open zip file
	r, err := zip.OpenReader(modZipFilename)
	if err != nil {
		return ModPackagingTypeUnknown, "", err
	}
	defer r.Close()

	containsDLLinRoot := false
	containsBepInExFolder := false
	prefixBeforeBepInExFolder := ""

	// loop through all the files in the archive to check determine what kind of packaging the mod has
	for _, f := range r.File {

		if !strings.ContainsAny(f.Name, "/") && strings.HasSuffix(f.Name, ".dll") {
			containsDLLinRoot = true
		}

		if strings.Contains(f.Name, "BepInEx/") {
			containsBepInExFolder = true
			index := strings.LastIndex(f.Name, "BepInEx/")
			prefixBeforeBepInExFolder = f.Name[0:index]
		}

	}

	log.Printf("containsDLL: %v, containsBepInExFolder: %v, %s", containsDLLinRoot, containsBepInExFolder, modZipFilename)

	if containsBepInExFolder {

		return ModPackagingTypeRoot, prefixBeforeBepInExFolder, nil
	}

	if containsDLLinRoot {
		return ModPackagingTypeRootDLL, "", nil
	}

	// assume Plugin by default
	return ModPackagingTypePlugin, "", nil
}
