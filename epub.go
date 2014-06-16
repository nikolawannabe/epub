package epub

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"text/template"
	"io/ioutil"
	"log"
	"fmt"
)

const (
	mimeTypeFileName = "mimetype"
	mimetype         = "application/epub+zip"
	templateGlob     = "/home/ctalbot/src/lugod/src/github.com/nikolawannabe/epub/templates/*.tpl"
)

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

type MetaInfFile struct {
	Name string
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

func (w *EpubArchive) Build(title string, opf Opf, chapters []Chapter) ([]byte, error) {
	buf := new(bytes.Buffer)
	w.Writer = *zip.NewWriter(buf)

	// The mimetype must always be first in the archive.
	f, err := w.Writer.Create(mimeTypeFileName)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	_, err = f.Write([]byte(mimetype))
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	metaFiles := w.buildMetaInfFiles(title, opf)
	w.addMetaInfFileToArchive(metaFiles)
	w.addChaptersToArchive(chapters)

	err = w.Writer.Close()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return buf.Bytes(), nil

}

//TODO: chapter and meta-inf file can likely be combined into one info that takes bytes.
// addMetaInfFileToArchive adds the meta-inf files to the archive.
func (w *EpubArchive) addMetaInfFileToArchive(files []MetaInfFile) error {
	for _, file := range files {
		f, err := w.Writer.Create(file.Name)
		if err != nil {
			log.Printf("%v", err)
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
		f, err := w.Writer.Create(chapter.FileName)
		if err != nil {
			log.Printf("%v", err)
		}
		_, err = f.Write([]byte(chapter.Contents))
		if err != nil {
			log.Printf("%v", err)
			return err
		}
	}
	return nil
}
