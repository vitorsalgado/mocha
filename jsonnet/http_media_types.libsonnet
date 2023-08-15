local _charsetUTF8 = 'charset:UTF-8';

{
  MIMEApplicationJSON: 'application/json',
  MIMEApplicationCharsetUTF8: self.MIMEApplicationJSON + '; ' + _charsetUTF8,
  MIMETextPlain: 'text/plain',
  MIMETextPlainCharsetUTF8: self.MIMETextPlain + '; ' + _charsetUTF8,
  MIMETextHTML: 'text/html',
  MIMETextHTMLCharsetUTF8: self.MIMETextHTML + '; ' + _charsetUTF8,
  MIMETextXML: 'text/xml',
  MIMETextXMLCharsetUTF8: self.MIMETextXML + '; ' + _charsetUTF8,
  MIMEFormURLEncoded: 'application/x-www-form-urlencoded',
  MIMEFormURLEncodedCharsetUTF8: self.MIMEFormURLEncoded + '; ' + _charsetUTF8,
  MIMEApplicationJavaScript: 'application/javascript',
  MIMEApplicationJavaScriptCharsetUTF8: self.MIMEApplicationJavaScript + '; ' + _charsetUTF8,
  MIMEApplicationXML: 'application/xml',
  MIMEApplicationXMLCharsetUTF8: self.MIMEApplicationXML + '; ' + _charsetUTF8,
  MIMEApplicationProtobuf: 'application/protobuf',
  MIMEMultipartForm: 'multipart/form-data',
  MIMEOctetStream: 'application/octet-stream',
}
