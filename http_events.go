package mocha

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
)

type mockHTTPLifecycle interface {
	OnRequest(*RequestValues)
	OnMatch(*RequestValues, *Stub)
	OnNoMatch(*RequestValues, *findResult)
	OnWarning(*RequestValues, error)
	OnError(*RequestValues, error)
}

// Standard HTTP Lifecycle

type builtInMockHTTPLifecycle struct {
	app *Mocha
}

func (h *builtInMockHTTPLifecycle) OnRequest(r *RequestValues) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	evt := h.app.logger.Info().Str("url", r.URL.String())
	req := zerolog.Dict()

	if h.app.config.LogVerbosity >= LogHeader {
		req.Any("header", r.RawRequest.Header)
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

	dict := zerolog.Dict().Str("id", r.Mock.ID)
	if r.Mock.Name != "" {
		dict.Str("name", r.Mock.Name)
	}

	res := zerolog.Dict().Int("status", s.StatusCode)

	if h.app.config.LogVerbosity >= LogHeader {
		res.Any("header", s.Header)
		if len(s.Trailer) > 0 {
			res.Any("trailer", s.Trailer)
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

func (h *builtInMockHTTPLifecycle) OnNoMatch(r *RequestValues, fr *findResult) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	evt := h.app.logger.Warn().Str("url", r.URL.String())
	if fr.ClosestMatch != nil {
		md := zerolog.Dict().Str("id", fr.ClosestMatch.ID)
		if fr.ClosestMatch.Name != "" {
			md.Str("name", fr.ClosestMatch.Name)
		}

		evt.Dict("closest_match", md)
	}

	if len(fr.MismatchDetails) > 0 {
		mismatches := zerolog.Arr()
		builder := strings.Builder{}

		for _, detail := range fr.MismatchDetails {
			builder.WriteString("[" + detail.Target.String())
			if detail.Key != "" {
				builder.WriteString("(" + detail.Key + ")")
			}
			builder.WriteString("] ")

			builder.WriteString(detail.MatchersName)

			if detail.Err == nil {
				builder.WriteString("(")

				if len(detail.Result.Ext) > 0 {
					builder.WriteString(strings.Join(detail.Result.Ext, ", "))
					builder.WriteString(") ")
					builder.WriteString(detail.Result.Message)
				} else {
					builder.WriteString(detail.Result.Message)
					builder.WriteString(") ")
				}
			} else {
				builder.WriteString(" ")
				builder.WriteString(detail.Err.Error())
			}

			mismatches.Str(builder.String())
			builder.Reset()
		}

		evt.Array("mismatches", mismatches)
	}

	evt.Msgf("<--- REQUEST WAS NOT MATCHED %s %s", r.RawRequest.Method, r.URL.Path)
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
	app *Mocha
	cz  *colorize.Colorize
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnRequest(rv *RequestValues) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s %s ---> %s %s\n%s %s",
		h.cz.BlueBright(h.cz.Bold("REQUEST RECEIVED")),
		rv.StartedAt.Format(time.RFC3339),
		h.cz.Blue(rv.RawRequest.Method),
		h.cz.Blue(rv.URL.Path),
		rv.RawRequest.Method,
		rv.URL))

	if h.app.config.LogVerbosity >= LogHeader && len(rv.RawRequest.Header) > 0 {
		builder.WriteString("\n")
		builder.WriteString(h.cz.Blue("Headers "))
		builder.WriteString(fmt.Sprintf("%s", rv.RawRequest.Header))
	}

	if h.app.config.LogVerbosity >= LogBody && len(rv.RawBody) > 0 {
		builder.WriteString("\n")
		builder.WriteString(h.cz.Blue("Body: "))
		builder.WriteString(fmt.Sprintf("%v", string(rv.RawBody)))
	}

	builder.WriteString("\n")

	fmt.Println(builder.String())
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnMatch(rv *RequestValues, s *Stub) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s %s <--- %s %s\n%s ",
		h.cz.GreenBright(h.cz.Bold("REQUEST MATCHED")),
		time.Now().Format(time.RFC3339),
		h.cz.Green(rv.RawRequest.Method),
		h.cz.Green(rv.URL.Path),
		rv.RawRequest.Method))
	builder.WriteString(rv.URL.String())
	builder.WriteString("\n")
	builder.WriteString(h.cz.Bold("Mock: "))
	builder.WriteString(rv.Mock.ID + " " + rv.Mock.Name)
	builder.WriteString("\n")
	builder.WriteString(h.cz.Green("Took(ms): "))
	builder.WriteString(strconv.FormatInt(time.Since(rv.StartedAt).Milliseconds(), 10))
	builder.WriteString("\n")
	builder.WriteString(h.cz.Green("Status: "))
	builder.WriteString(strconv.FormatInt(int64(s.StatusCode), 10))

	if h.app.config.LogVerbosity >= LogHeader && len(s.Header) > 0 {
		builder.WriteString("\n")
		builder.WriteString(h.cz.Green("Headers: "))
		builder.WriteString(fmt.Sprintf("%s", s.Header))
	}

	bodyLen := int64(len(s.Body))
	if h.app.config.LogVerbosity >= LogBody && bodyLen > 0 &&
		(h.app.config.LogBodyMaxSize == 0 || bodyLen <= h.app.config.LogBodyMaxSize) {
		builder.WriteString("\n")
		builder.WriteString(h.cz.Green("Body: "))
		if s.Encoding == "" {
			builder.WriteString(string(s.Body))
		} else {
			builder.WriteString("<encoded body omitted>")
		}
	}

	builder.WriteString("\n")

	fmt.Println(builder.String())
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnNoMatch(rv *RequestValues, fr *findResult) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s %s <--- %s %s\n%s %s\n",
		h.cz.YellowBright(h.cz.Bold("REQUEST WAS NOT MATCHED")),
		time.Now().Format(time.RFC3339),
		h.cz.Yellow(rv.RawRequest.Method),
		h.cz.Yellow(rv.URL.Path),
		rv.RawRequest.Method,
		rv.URL.String()))

	if fr.ClosestMatch != nil {
		builder.WriteString(fmt.Sprintf("%s: %s %s\n",
			h.cz.Bold("Closest Match"), fr.ClosestMatch.ID, fr.ClosestMatch.Name))
	}

	if len(fr.MismatchDetails) > 0 {
		if h.app.config.LogVerbosity <= LogHeader {
			builder.WriteString(fmt.Sprintf("%s: %d", h.cz.Bold("Mismatches"), len(fr.MismatchDetails)))
		} else {
			builder.WriteString(fmt.Sprintf("%s(%d):\n", h.cz.Bold("Mismatches"), len(fr.MismatchDetails)))
			for _, detail := range fr.MismatchDetails {
				builder.WriteString("[" + detail.Target.String())
				if detail.Key != "" {
					builder.WriteString("(" + detail.Key + ")")
				}
				builder.WriteString("] ")
				builder.WriteString(detail.MatchersName)
				builder.WriteString("(")

				if detail.Err != nil {
					builder.WriteString(detail.Err.Error())
					builder.WriteString(")\n")
					continue
				}

				if detail.Result == nil {
					builder.WriteString(")\n")
					continue
				}

				if len(detail.Result.Ext) == 0 {
					builder.WriteString(h.cz.Bold(detail.Result.Message))
					builder.WriteString(")")
				} else {
					builder.WriteString(h.cz.Bold(strings.Join(detail.Result.Ext, ", ")))
					builder.WriteString(") ")
					builder.WriteString(detail.Result.Message)
				}

				builder.WriteString("\n")
			}
		}
	}

	builder.WriteString("\n")

	fmt.Println(builder.String())
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnWarning(r *RequestValues, err error) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	fmt.Printf("\n%s %s <--- %s %s\n%s %s\n%s: %v\n\n",
		h.cz.RedBright(h.cz.Bold("WARNING")),
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

	fmt.Printf("%s %s <--- %s %s\n%s %s\n%s: %v\n\n",
		h.cz.RedBright(h.cz.Bold("ERROR")),
		time.Now().Format(time.RFC3339),
		h.cz.Red(r.RawRequest.Method),
		h.cz.Red(r.URL.Path),
		r.RawRequest.Method,
		r.URL.String(),
		h.cz.Red(h.cz.Bold("Error")),
		err,
	)
}
