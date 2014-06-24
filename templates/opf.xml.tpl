<?xml version="1.0" encoding="utf-8" standalone="no"?>
<package xmlns="http://www.idpf.org/2007/opf" xmlns:dc="http://purl.org/dc/elements/1.1/"
         xmlns:dcterms="http://purl.org/dc/terms/" version="3.0" xml:lang="en"
         unique-identifier="pub_id_0">
    <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
        {{ range $i, $identifier := .Identifiers }}
        <dc:identifier
                id="pub_id_{{ $i }}">{{ .Value }}</dc:identifier>
            <meta refines="#pub_id_{{ $i }}"
                  property="identifier-type"
                  scheme="onix:codelist5">{{ $identifier.IdentifierType }}</meta>
        {{ end }}
        <dc:language>{{ .Metadata.Language }} </dc:language>
        <dc:title>{{ .Metadata.Title }}</dc:title>
        {{ if .Metadata.Creator.Name }}
        <dc:creator id="creator">{{ .Metadata.Creator.Name }}</dc:creator>
        <meta refines="#creator"
              property="role"
              scheme="marc:relators">{{ .Metadata.Creator.Role }}</meta>
        {{ end }}
        {{ if .Metadata.Date }}
        <meta property="dcterms:modified">{{ .Metadata.Date}}</meta>
        {{ end }}
        {{ if .Metadata.Publisher }}
        <dc:publisher>{{ .Metadata.Publisher}}</dc:publisher>
        {{ end }}
        {{ if .Metadata.Date }}
        <dc:date>{{ .Metadata.Date}}</dc:date>
        {{ end }}
    </metadata>
    <manifest>
        {{ range $manifestItem :=  .Manifest.ManifestItems }}
        <item id="{{ $manifestItem.Id}}"
              href="{{ $manifestItem.Href }}"
              media-type="{{ $manifestItem.MediaType}}"/>
        {{ end }}
    </manifest>
    <spine>
        {{ range $manifestItem :=  .Manifest.ManifestItems }}
        <itemref idref="{{ $manifestItem.Id }}"/>
        {{ end }}
    </spine>
</package>