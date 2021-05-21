package r2modman

import (
	"archive/zip"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ExportParser interface {
	Parse(file string) (*ExportR2x, error)
}

type exportParserImpl struct {
}

func (p *exportParserImpl) Parse(file string) (*ExportR2x, error) {

	reader, err := zip.OpenReader(file)
	if err != nil {
		return nil, err
	}

	var exportMetadata *ExportR2x

	for _, v := range reader.File {

		fileread, err := v.Open()
		if err != nil {
			return nil, err
		}
		defer fileread.Close()

		// parse metadata file
		if v.Name == "export.r2x" {
			content, err := ioutil.ReadAll(fileread)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(content, &exportMetadata)
			if err != nil {
				return nil, err
			}
		}
	}

	if exportMetadata == nil {
		return nil, fmt.Errorf("%s does not contain an export.r2x file", file)
	}

	return exportMetadata, nil
}

func newExportParser() (ExportParser, error) {
	return &exportParserImpl{}, nil
}
