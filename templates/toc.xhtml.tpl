<nav epub:type="toc" id="toc">
    <h2>RSS Entries</h2>
    <ol>
        {{ range $manifestItem :=  .ManifestItems }}
        {{ if $manifestItem.Title }}
        <li>
            <a href="../{{ $manifestItem.Href }}">
                {{ $manifestItem.Title }}
            </a>
        </li>
        {{ end }}
        {{ end }}
    </ol>
</nav>