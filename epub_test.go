package epub

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFillTemplate(t *testing.T) {
	var epub EpubArchive
	var opf Opf
	epub.fillTemplates(opf)
	assert.True(t, false)
}

func TestGetMetadata(t *testing.T) {
	var epub EpubArchive
	m, err := epub.getMetadata("testdata/metadata.json")
	assert.NoError(t, err)
	assert.Equal(t, m.Publisher, "I am a publisher")
}
