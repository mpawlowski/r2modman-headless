package r2modman

import (
	"context"
	"encoding/json"
	"net/http"
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

var myClient = &http.Client{}

func GetPackagesMetadata(ctx context.Context, metadataUrl string) (packages map[string]*APIPackageResponse, err error) {

	req, err := http.NewRequestWithContext(ctx, "GET", metadataUrl, nil)
	if err != nil {
		return
	}

	r, err := myClient.Do(req)
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
