// Package mimetypes contains common mime types used internally by Mocha.
package mimetype

const (
	charsetUTF8 = "charset=UTF-8"
)

// Common mime types.
const (
	JSON                      = "application/json"
	JSONCharsetUTF8           = JSON + "; " + charsetUTF8
	TextPlain                 = "text/plain"
	TextPlainCharsetUTF8      = TextPlain + "; " + charsetUTF8
	TextHTML                  = "text/html"
	FormURLEncoded            = "application/x-www-form-urlencoded"
	FormURLEncodedCharsetUTF8 = FormURLEncoded + "; " + charsetUTF8
)
