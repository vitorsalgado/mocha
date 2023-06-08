// Package mhttpv implements utilities that might be useful to the users of this library.
package httpval

const _charsetUTF8 = "charset=UTF-8"

// Headers
const (
	HeaderAccept                          = "Accept"
	HeaderAcceptEncoding                  = "Accept-Encoding"
	HeaderContentType                     = "Content-Type"
	HeaderContentEncoding                 = "Content-Encoding"
	HeaderAllow                           = "Allow"
	HeaderAuthorization                   = "Authorization"
	HeaderContentDisposition              = "Content-Disposition"
	HeaderVary                            = "Vary"
	HeaderOrigin                          = "Origin"
	HeaderContentLength                   = "Content-length"
	HeaderConnection                      = "Connection"
	HeaderTrailer                         = "Trailer"
	HeaderLocation                        = "Location"
	HeaderCacheControl                    = "Cache-Control"
	HeaderCookie                          = "Cookie"
	HeaderSetCookie                       = "Set-Cookie"
	HeaderIfModifiedSince                 = "If-Modified-Since"
	HeaderLastModified                    = "Last-Modified"
	HeaderRetryAfter                      = "Retry-After"
	HeaderUpgrade                         = "Upgrade"
	HeaderWWWAuthenticate                 = "WWW-Authenticate"
	HeaderXForwardedFor                   = "X-Forwarded-For"
	HeaderXForwardedProto                 = "X-Forwarded-Proto"
	HeaderXForwardedProtocol              = "X-Forwarded-Protocol"
	HeaderXForwardedSsl                   = "X-Forwarded-Ssl"
	HeaderXUrlScheme                      = "X-Url-Scheme"
	HeaderXHTTPMethodOverride             = "X-HTTP-Method-Override"
	HeaderXRealIP                         = "X-Real-Ip"
	HeaderXRequestID                      = "X-Request-Id"
	HeaderXCorrelationID                  = "X-Correlation-Id"
	HeaderXRequestedWith                  = "X-Requested-With"
	HeaderServer                          = "Server"
	HeaderAccessControlRequestMethod      = "Access-Control-Request-Method"
	HeaderAccessControlAllowOrigin        = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods       = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders       = "Access-Control-Allow-Header"
	HeaderAccessControlExposeHeaders      = "Access-Control-Expose-Header"
	HeaderAccessControlMaxAge             = "Access-Control-max-Age"
	HeaderAccessControlAllowCredentials   = "Access-Control-Allow-Credentials"
	HeaderAccessControlRequestHeaders     = "Access-Control-Request-Header"
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderXCSRFToken                      = "X-CSRF-Token"
	HeaderReferrerPolicy                  = "Referrer-Policy"
)

// MIME Types
const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationCharsetUTF8           = MIMEApplicationJSON + "; " + _charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + _charsetUTF8
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + _charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + _charsetUTF8
	MIMEFormURLEncoded                   = "application/x-www-form-urlencoded"
	MIMEFormURLEncodedCharsetUTF8        = MIMEFormURLEncoded + "; " + _charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + _charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + _charsetUTF8
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)
