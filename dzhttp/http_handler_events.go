package dzhttp

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jwalton/gchalk"
	"github.com/rs/zerolog"

	"github.com/vitorsalgado/mocha/v3/dzstd"
)

var (
	_ mockHTTPLifecycle = (*builtInMockHTTPLifecycle)(nil)
	_ mockHTTPLifecycle = (*builtInDescriptiveMockHTTPLifecycle)(nil)
)

type mockHTTPLifecycle interface {
	OnRequest(*RequestValues)
	OnMatch(*RequestValues, *Stub)
	OnNoMatch(*RequestValues, *dzstd.FindResult[*HTTPMock], *dzstd.Description)
	OnWarning(*RequestValues, error)
	OnError(*RequestValues, error)
}

// Standard HTTP Lifecycle

type builtInMockHTTPLifecycle struct {
	app *HTTPMockApp
}

func newBuiltInMockHTTPLifecycle(app *HTTPMockApp) *builtInMockHTTPLifecycle {
	return &builtInMockHTTPLifecycle{app}
}

func (h *builtInMockHTTPLifecycle) OnRequest(r *RequestValues) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	evt := h.app.logger.Info().Str("url", r.URL.String())
	req := zerolog.Dict()

	if h.app.config.LogVerbosity >= LogHeader {
		redactedHeader := redactHeader(r.RawRequest.Header, h.app.config.HeaderNamesToRedact)
		req.Any("header", redactedHeader)
	}

	if h.app.config.LogVerbosity >= LogBody && len(r.RawBody) > 0 {
		req.Bytes("body", r.RawBody)
	}

	evt.Dict("request", req)
	evt.Msgf("---> REQUEST RECEIVED %s %s", r.RawRequest.Method, r.URL.Path)
}

func (h *builtInMockHTTPLifecycle) OnMatch(r *RequestValues, s *Stub) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	evt := h.app.logger.Info().
		Str("url", r.URL.String()).
		Dur("elapsed", time.Since(r.StartedAt))

	dict := zerolog.Dict().Str("id", r.Mock.GetID())
	if r.Mock.GetName() != "" {
		dict.Str("name", r.Mock.GetName())
	}

	res := zerolog.Dict().Int("status", s.StatusCode)

	if h.app.config.LogVerbosity >= LogHeader {
		redactedHeader := redactHeader(s.Header, h.app.config.HeaderNamesToRedact)
		res.Any("header", redactedHeader)

		if len(s.Trailer) > 0 {
			redactedTrailer := redactHeader(s.Trailer, h.app.config.HeaderNamesToRedact)
			res.Any("trailer", redactedTrailer)
		}
	}

	bodyLen := int64(len(s.Body))
	if h.app.config.LogVerbosity >= LogBody && bodyLen > 0 &&
		(h.app.config.LogBodyMaxSize == 0 || bodyLen <= h.app.config.LogBodyMaxSize) {
		res.Bytes("body", s.Body)
		if s.Encoding != "" {
			res.Str("encoding", s.Encoding)
		}
	}

	evt.Dict("mock", dict)
	evt.Dict("response", res)

	evt.Msgf("<--- REQUEST MATCHED %s %s", r.RawRequest.Method, r.URL.Path)
}

func (h *builtInMockHTTPLifecycle) OnNoMatch(r *RequestValues, fr *dzstd.FindResult[*HTTPMock], d *dzstd.Description) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	evt := h.app.logger.Warn().Str("url", r.URL.String())
	if fr.ClosestMatch != nil {
		md := zerolog.Dict().Str("id", fr.ClosestMatch.GetID())
		if fr.ClosestMatch.GetName() != "" {
			md.Str("name", fr.ClosestMatch.GetName())
		}

		evt.Dict("nearest", md)
	}

	if d.Len() > 0 {
		evt.Strs("mismatches", d.Buf)
	}

	evt.Msgf("<--- REQUEST DID NOT MATCH %s %s", r.RawRequest.Method, r.URL.Path)
}

func (h *builtInMockHTTPLifecycle) OnWarning(r *RequestValues, err error) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	h.app.logger.Warn().
		Err(err).
		Str("url", r.URL.String()).
		Msgf("<--- WARNING %s %s", r.RawRequest.Method, r.URL.Path)
}

func (h *builtInMockHTTPLifecycle) OnError(r *RequestValues, err error) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	h.app.logger.Error().
		Err(err).
		Str("url", r.URL.String()).
		Msgf("<--- ERROR %s %s", r.RawRequest.Method, r.URL.Path)
}

// Descriptive Logger

type builtInDescriptiveMockHTTPLifecycle struct {
	app *HTTPMockApp
	cz  *gchalk.Builder
	out io.Writer
}

func newBuiltInDescriptiveMockHTTPLifecycle(app *HTTPMockApp, cz *gchalk.Builder, out io.Writer) *builtInDescriptiveMockHTTPLifecycle {
	return &builtInDescriptiveMockHTTPLifecycle{app, cz, out}
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnRequest(rv *RequestValues) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	buf := strings.Builder{}

	buf.WriteString(fmt.Sprintf("%s %s ---> %s %s\n%s %s",
		h.cz.BrightBlue(h.cz.Bold("REQUEST RECEIVED")),
		rv.StartedAt.Format(time.RFC3339),
		h.cz.Blue(rv.RawRequest.Method),
		h.cz.Blue(rv.URL.Path),
		rv.RawRequest.Method,
		rv.URL))

	if h.app.config.LogVerbosity >= LogHeader && len(rv.RawRequest.Header) > 0 {
		redactedHeader := redactHeader(rv.RawRequest.Header, h.app.config.HeaderNamesToRedact)

		buf.WriteString("\n")
		buf.WriteString(h.cz.Blue("Headers "))
		buf.WriteString(fmt.Sprintf("%s", redactedHeader))
	}

	if h.app.config.LogVerbosity >= LogBody && len(rv.RawBody) > 0 {
		buf.WriteString("\n")
		buf.WriteString(h.cz.Blue("Body: "))
		buf.WriteString(fmt.Sprintf("%v", string(rv.RawBody)))
	}

	buf.WriteString("\n\n")
	fmt.Fprint(h.out, buf.String())
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnMatch(rv *RequestValues, s *Stub) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	buf := strings.Builder{}

	fmt.Fprintf(&buf, "%s %s <--- %s %s\n%s ",
		h.cz.BrightGreen(h.cz.Bold("REQUEST MATCHED")),
		time.Now().Format(time.RFC3339),
		h.cz.Green(rv.RawRequest.Method),
		h.cz.Green(rv.URL.Path),
		rv.RawRequest.Method)
	buf.WriteString(rv.URL.String())
	buf.WriteString("\n")
	buf.WriteString(h.cz.Bold("Mock: "))
	buf.WriteString(rv.Mock.GetID() + " " + rv.Mock.GetName())
	buf.WriteString("\n")
	buf.WriteString(h.cz.Green("Took(ms): "))
	buf.WriteString(strconv.FormatInt(time.Since(rv.StartedAt).Milliseconds(), 10))
	buf.WriteString("\n")
	buf.WriteString(h.cz.Green("Status: "))
	buf.WriteString(strconv.FormatInt(int64(s.StatusCode), 10))

	if h.app.config.LogVerbosity >= LogHeader && len(s.Header) > 0 {
		buf.WriteString("\n")
		buf.WriteString(h.cz.Green("Headers: "))
		fmt.Fprintf(&buf, "%s", redactHeader(s.Header, h.app.config.HeaderNamesToRedact))

		if len(s.Trailer) > 0 {
			buf.WriteString("\n")
			buf.WriteString(h.cz.Green("Trailers: "))
			fmt.Fprintf(&buf, "%s", redactHeader(s.Trailer, h.app.config.HeaderNamesToRedact))
		}
	}

	bodyLen := int64(len(s.Body))
	if h.app.config.LogVerbosity >= LogBody && bodyLen > 0 &&
		(h.app.config.LogBodyMaxSize == 0 || bodyLen <= h.app.config.LogBodyMaxSize) {
		buf.WriteString("\n")
		buf.WriteString(h.cz.Green("Body: "))
		if s.Encoding == "" {
			buf.WriteString(string(s.Body))
		} else {
			buf.WriteString("<encoded body omitted>")
		}
	}

	buf.WriteString("\n")

	fmt.Fprint(h.out, buf.String())
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnNoMatch(
	reqValues *RequestValues,
	result *dzstd.FindResult[*HTTPMock],
	d *dzstd.Description,
) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	fmt.Fprintf(h.out, "%s %s <--- %s %s\n%s %s\n",
		h.cz.BrightYellow(h.cz.Bold("REQUEST DID NOT MATCH")),
		time.Now().Format(time.RFC3339),
		h.cz.Yellow(reqValues.RawRequest.Method),
		h.cz.Yellow(reqValues.URL.Path),
		reqValues.RawRequest.Method,
		reqValues.URL.String())

	if result.ClosestMatch != nil {
		fmt.Fprintf(h.out, "%s: %s %s\n",
			h.cz.Bold("Closest Match"), result.ClosestMatch.GetID(), result.ClosestMatch.GetName())
	}

	if d.Len() > 0 {
		if h.app.config.LogVerbosity <= LogHeader {
			fmt.Fprintf(h.out, "%s: %d", h.cz.Bold("Mismatches"), result.MismatchesCount)
		} else {
			fmt.Fprintf(h.out, "%s(%d):\n%s", h.cz.Bold("Mismatches"), result.MismatchesCount, d.String())
		}
	}
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnWarning(r *RequestValues, err error) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	fmt.Fprintf(h.out, "\n%s %s <--- %s %s\n%s %s\n%s: %v\n\n",
		h.cz.BrightRed(h.cz.Bold("WARNING")),
		time.Now().Format(time.RFC3339),
		h.cz.Red(r.RawRequest.Method),
		h.cz.Red(r.URL.Path),
		r.RawRequest.Method,
		r.URL.String(),
		h.cz.Red(h.cz.Bold("Error")),
		err,
	)
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnError(r *RequestValues, err error) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	fmt.Fprintf(h.out, "%s %s <--- %s %s\n%s %s\n%s: %v\n\n",
		h.cz.BrightRed(h.cz.Bold("ERROR")),
		time.Now().Format(time.RFC3339),
		h.cz.Red(r.RawRequest.Method),
		h.cz.Red(r.URL.Path),
		r.RawRequest.Method,
		r.URL.String(),
		h.cz.Red(h.cz.Bold("Error")),
		err,
	)
}

func redactHeader(h http.Header, toRedact map[string]struct{}) http.Header {
	if len(toRedact) == 0 {
		return h
	}

	redacted := h.Clone()
	for k := range h {
		if _, ok := toRedact[strings.ToLower(k)]; ok {
			redacted.Set(k, "<redacted>")
		}
	}

	return redacted
}
