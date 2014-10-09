package epub

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// This test shows that there is a bug!
func TestFillTemplate(t *testing.T) {
	var epub EpubArchive
	var opf Opf
	var chapters []Chapter
	epub.Build("hello", opf, chapters)
	assert.True(t, false)
}

func TestGetMetadata(t *testing.T) {
	var epub EpubArchive
	m, err := epub.getMetadata("testdata/metadata.json")
	assert.NoError(t, err)
	assert.Equal(t, m.Publisher, "I am a publisher")
}
