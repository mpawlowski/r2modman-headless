package r2modman

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportParser(t *testing.T) {
	parser, err := newExportParser()
	assert.Nil(t, err)

	metadata, err := parser.Parse("testdata/Realm_Of_Atflac_2021_05_15.r2z")
	assert.Nil(t, err)

	assert.Equal(t, 10, len(metadata.Mods))
	assert.Equal(t, "Realm_Of_Atflac_2021_05_15", metadata.ProfileName)
}
