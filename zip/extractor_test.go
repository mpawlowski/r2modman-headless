package zip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWindowsDirAsFile(t *testing.T) {
	tempDir := t.TempDir()
	extractor := newExtractor()
	err := extractor.Extract("testdata/windows-dir-as-file.zip", tempDir, "")
	assert.Nil(t, err)
}
