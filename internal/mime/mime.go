package mime

import "strings"

const (
	ContentTypeJSON           = "application/json"
	ContentTypeXML            = "application/xml"
	ContentTypeTextHTML       = "text/html"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

var ext = map[string]string{
	ContentTypeJSON:     ".json",
	ContentTypeXML:      ".xml",
	ContentTypeTextHTML: ".html",
}

func ExtensionFor(contenttype string) string {
	contenttype = strings.TrimSpace(strings.Split(contenttype, ",")[0])
	extension, ok := ext[contenttype]
	if !ok {
		return ".txt"
	}

	return extension
}
