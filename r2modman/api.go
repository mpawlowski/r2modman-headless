package r2modman

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type APIPackageResponse struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Versions []struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		FileSize int64  `json:"file_size"`
	} `json:"versions"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetPackagesMetadata(ctx context.Context, metadataUrl string) (packages map[string]*APIPackageResponse, err error) {

	r, err := myClient.Get(metadataUrl)
	if err != nil {
		return
	}
	defer r.Body.Close()

	var packageMetadataList []*APIPackageResponse
	err = json.NewDecoder(r.Body).Decode(&packageMetadataList)
	if err != nil {
		return
	}

	packages = map[string]*APIPackageResponse{}
	for _, p := range packageMetadataList {
		packages[p.FullName] = p
	}
	return
}
