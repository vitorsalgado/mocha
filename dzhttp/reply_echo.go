package dzhttp

import (
	"io"
	"net/http"
	"os"

	"github.com/vitorsalgado/mocha/v3/dzhttp/httpval"
)

type ReplyEcho struct {
	stdoutLog bool
}

func Echo() *ReplyEcho {
	return &ReplyEcho{}
}

func (e *ReplyEcho) Log() *ReplyEcho {
	e.stdoutLog = true
	return e
}

func (e *ReplyEcho) Build(w http.ResponseWriter, r *RequestValues) (*Stub, error) {
	w.Header().Add(httpval.HeaderContentType, httpval.MIMETextPlainCharsetUTF8)

	if e.stdoutLog {
		r.RawRequest.Write(io.MultiWriter(w, os.Stdout))
	} else {
		r.RawRequest.Write(w)
	}

	return nil, nil
}

func (e *ReplyEcho) Describe() any {
	return map[string]any{"echo": map[string]any{
		"stdout_log": e.stdoutLog,
	}}
}
