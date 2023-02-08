// Package mimetypes contains common mime types used internally by Mocha.
package mimetype

const (
	JSON                             = "application/json"
	JSONCharsetUTF8                  = JSON + "; " + _charsetUTF8
	TextPlain                        = "text/plain"
	TextPlainCharsetUTF8             = TextPlain + "; " + _charsetUTF8
	TextHTML                         = "text/html"
	TextXML                          = "text/xml"
	TextXMLCharsetUTF8               = TextXML + "; " + _charsetUTF8
	FormURLEncoded                   = "application/x-www-form-urlencoded"
	FormURLEncodedCharsetUTF8        = FormURLEncoded + "; " + _charsetUTF8
	ApplicationJavaScript            = "application/javascript"
	ApplicationJavaScriptCharsetUTF8 = ApplicationJavaScript + "; " + _charsetUTF8
	ApplicationXML                   = "application/xml"
	ApplicationXMLCharsetUTF8        = ApplicationXML + "; " + _charsetUTF8
	ApplicationProtobuf              = "application/protobuf"
	ApplicationMsgpack               = "application/msgpack"
	TextHTMLCharsetUTF8              = TextHTML + "; " + _charsetUTF8
	MultipartForm                    = "multipart/form-data"
	OctetStream                      = "application/octet-stream"
)

const (
	_charsetUTF8 = "charset=UTF-8"
)
