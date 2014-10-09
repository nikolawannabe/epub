// Package epub provides a mechanism for generating a (nominally) valid ePub 3 file from a list
// of metadata and content.
package epub

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"text/template"
)

const (
	mimeTypeFileName = "mimetype"
	// mimetype is the required mime type of an epub file
	mimetype         = "application/epub+zip"

	// TODO: these files should be distributed as part of the binary, probably.
	// templateGlob is the location of the epub xml templates on your file system.
	templateGlob     = "/home/ctalbot/src/lugod/src/github.com/nikolawannabe/epub/templates/*.tpl"
)

// EpubArchive contains the necessary structs to generate an epub
type EpubArchive struct {
	zip.Writer
	Opf Opf
}

// getMetadata gets the book metadata from a JSON file.
func (w *EpubArchive) getMetadata(filePath string) (Metadata, error) {
	var metadata Metadata
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Unable to read metadata file: %v", err)
		return metadata, err
	}
	err = json.Unmarshal(bytes, &metadata)
	if err != nil {
		log.Printf("Unable to unmarshall metadat file: %v", err)
		return metadata, err
	}
	return metadata, nil
}

// MetaInfFile contains information about the meta inf files which must be included for an epub
// to be valid.
type MetaInfFile struct {
	Name    string
	Content []byte
}

// buildMetaInfFiles creates the xml files for the epub via the information specified in the opf.
func (w *EpubArchive) buildMetaInfFiles(title string, opf Opf) []MetaInfFile {
	var metaInfFiles []MetaInfFile
	tmpl := template.Must(template.ParseGlob(templateGlob))

	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "container.xml.tpl", opf)
	if err != nil {
		log.Fatal("template execution: %s", err) // TODO: this should return as err
	}
	containerXml := MetaInfFile{Name: "META-INF/container.xml",
		Content: buf.Bytes()}
	metaInfFiles = append(metaInfFiles, containerXml)

	buf = new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, "opf.xml.tpl", opf.RootFiles[0])
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	opfXml := MetaInfFile{Name: fmt.Sprintf("OEBPS/%s.opf", title),
		Content: buf.Bytes()}
	metaInfFiles = append(metaInfFiles, opfXml)

	return metaInfFiles
}

// Build generates an epub file if possible from the title, opf, and chapters
func (w *EpubArchive) Build(title string, opf Opf, chapters []Chapter) ([]byte, error) {
	buf := new(bytes.Buffer)
	w.Writer = *zip.NewWriter(buf)

	// The mimetype must always be first in the file, and must not be compressed.
	header := &zip.FileHeader{
		Name:   mimeTypeFileName,
		Method: zip.Store,
	}
	f, err := w.CreateHeader(header)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	_, err = f.Write([]byte(mimetype))
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	tocManifest := ManifestItem{
		Id:         "tocref",
		Href:       "chapters/toc.xhtml",
		MediaType:  "application/xhtml+xml",
		Properties: []string{"nav"},
	}

	opf.RootFiles[0].Manifest.ManifestItems = append(opf.RootFiles[0].Manifest.ManifestItems,
		tocManifest)

	metaFiles := w.buildMetaInfFiles(title, opf)
	w.addMetaInfFileToArchive(metaFiles)
	w.addChaptersToArchive(chapters)

	w.addChaptersToArchive([]Chapter{w.buildToc(opf)})

	err = w.Writer.Close()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return buf.Bytes(), nil

}

// buildToc generates the table of contents for the epub file.
func (w *EpubArchive) buildToc(opf Opf) Chapter {
	tmpl := template.Must(template.ParseGlob(templateGlob))
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "toc.xhtml.tpl", opf.RootFiles[0].Manifest)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	toc := Chapter{FileName: "chapters/toc.xhtml",
		Contents: string(buf.Bytes())}
	return toc
}

//TODO: chapter and meta-inf file can likely be combined into one info that takes bytes.
// addMetaInfFileToArchive adds the meta-inf files to the archive.
func (w *EpubArchive) addMetaInfFileToArchive(files []MetaInfFile) error {
	for _, file := range files {
		f, err := w.Writer.Create(file.Name)
		if err != nil {
			log.Printf("%v", err)
			return err
		}
		_, err = f.Write(file.Content)
		if err != nil {
			log.Printf("%v", err)
			return err
		}
	}
	return nil
}

// addChaptersToArchive appends the content for each chapter with the specified filename to the zip.
func (w *EpubArchive) addChaptersToArchive(chapters []Chapter) error {
	for _, chapter := range chapters {
		f, err := w.Writer.Create(filepath.Join("OEBPS", chapter.FileName))
		if err != nil {
			log.Printf("%v", err)
			return err
		}
		_, err = f.Write([]byte(chapter.Contents))
		if err != nil {
			log.Printf("%v", err)
			return err
		}
	}
	return nil
}
