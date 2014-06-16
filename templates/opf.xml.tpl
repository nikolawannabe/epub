<?xml version="1.0"?>
<package version="3.0"
         xml:lang="en"
         xmlns="http://www.idpf.org/2007/opf"
         unique-identifier="pub-id">

    <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
        <dc:identifier
                id="pub-id">{{ .Metadata.ISBN13 }}</dc:identifier>

        <dc:language>{{ .Metadata.Language }} </dc:language>

        <dc:title>{{ .Metadata.Title }}</dc:title>

        <dc:creator id="creator">{{ .Metadata.Creator.Name }}</dc:creator>
        <meta refines="#creator"
              property="role"
              scheme="marc:relators">{{ .Metadata.Creator.Role }}</meta>

        <meta property="dcterms:modified">{{ .Metadata.Date}}</meta>

        <dc:publisher>{{ .Metadata.Publisher}}</dc:publisher>

        <dc:date>{{ .Metadata.Date}}</dc:date>

        <dc:identifier
                id="isbn13">urn:isbn:{{ .Metadata.ISBN13 }}</dc:identifier>
        <meta refines="#isbn13"
              property="identifier-type"
              scheme="onix:codelist5">15</meta>
    </metadata>

    <manifest>

        {{ range $manifestItem :=  .Manifest.ManifestItems }}
        <item id="{{ $manifestItem.Id}}"
              href="../{{ $manifestItem.Href }}"
              media-type="{{ $manifestItem.MediaType}}"/>
        {{ end }}
    </manifest>
    <spine>
        {{ range $manifestItem :=  .Manifest.ManifestItems }}
        <itemref idref="{{ $manifestItem.Id }}"/>
        {{ end }}
    </spine>
</package>