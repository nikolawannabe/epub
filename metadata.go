package epub

type Identifier struct {
	Value	string
	IdentifierType	string
}

type Creator struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type Metadata struct {
	Language        string  `json:"language"`
	Title           string  `json:"title"`
	Creator         Creator `json:"creator"`
	Publisher       string  `json:"publisher"`
	Date            string  `json:"date"`
	DateCopyrighted string  `json:"date_copyrighted"`
	ISBN13          string  `json:"isbn13"`
}

type ManifestItem struct {
	Id         string
	Href       string
	MediaType  string
	Properties []string // http://www.idpf.org/epub/30/spec/epub30-publications.html#sec-item-property-values
}

type Manifest struct {
	ManifestItems []ManifestItem
}

type OpfRootFile struct {
	FullPath   string
	MediaType  string
	Identifiers []Identifier
	Metadata   Metadata
	Manifest   Manifest
}

type Opf struct {
	RootFiles []OpfRootFile
}
