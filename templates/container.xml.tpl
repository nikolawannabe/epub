<?xml version="1.0"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
    <rootfiles>
        {{ range .RootFiles }}
        <rootfile full-path="{{ .FullPath}}"
                  media-type="{{ .MediaType }}" />
        {{ end }}
    </rootfiles>
</container>