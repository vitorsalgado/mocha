package mocha

import (
	"fmt"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
)

type internalEvents struct {
	l     TestingT
	level LogLevel
}

func newInternalEvents(l TestingT, level LogLevel) *internalEvents {
	return &internalEvents{l: l, level: level}
}

func (h *internalEvents) OnRequest(evt any) {
	e := evt.(*OnRequest)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("\n%s %s ---> %s %s\n%s %s\n\n%s: %v",
		colorize.BlueBright(colorize.Bold("REQUEST RECEIVED")),
		e.StartedAt.Format(time.RFC3339),
		colorize.Blue(e.Request.Method),
		colorize.Blue(e.Request.Path),
		e.Request.Method,
		e.Request.FullURL(),
		colorize.Blue("Headers"),
		e.Request.Header))

	if h.level == LogVerbose {
		var b any
		switch t := e.Request.Body.(type) {
		case []byte:
			b = string(t)
		default:
			b = t
		}

		builder.WriteString("\n")
		builder.WriteString(colorize.Blue("Body: "))
		builder.WriteString(fmt.Sprintf("%v\n", b))
	}

	h.l.Logf(builder.String())
}

func (h *internalEvents) OnRequestMatched(evt any) {
	e := evt.(*OnRequestMatch)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("\n%s %s <--- %s %s\n%s %s\n",
		colorize.GreenBright(colorize.Bold("REQUEST DID MATCH")),
		time.Now().Format(time.RFC3339),
		colorize.Green(e.Request.Method),
		colorize.Green(e.Request.Path),
		e.Request.Method,
		e.Request.FullURL()))

	if h.level == LogVerbose {
		builder.WriteString(fmt.Sprintf("\n%s%d %s\n\n%s: %dms\n%s:\n %s: %d\n %s: %v\n",
			colorize.Bold("Mock: "),
			e.Mock.ID,
			e.Mock.Name,
			colorize.Green("Took"),
			e.Elapsed.Milliseconds(),
			colorize.Green("Response Definition"),
			colorize.Green("Status"),
			e.ResponseDefinition.Status,
			colorize.Green("Headers"),
			e.ResponseDefinition.Header))
	}

	h.l.Logf(builder.String())
}

func (h *internalEvents) OnRequestNotMatched(evt any) {
	e := evt.(*OnRequestNotMatched)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("\n%s %s <--- %s %s\n%s %s\n\n",
		colorize.YellowBright(colorize.Bold("REQUEST DID NOT MATCH")),
		time.Now().Format(time.RFC3339),
		colorize.Yellow(e.Request.Method),
		colorize.Yellow(e.Request.Path),
		e.Request.Method,
		e.Request.FullURL()))

	if e.Result.HasClosestMatch {
		builder.WriteString(fmt.Sprintf("%s: %d %s\n\n",
			colorize.Bold("Closest Match"), e.Result.ClosestMatch.ID, e.Result.ClosestMatch.Name))
	}

	if h.level == LogVerbose {
		builder.WriteString(fmt.Sprintf("%s:\n", colorize.Bold("Mismatches")))

		for _, detail := range e.Result.Details {
			builder.WriteString(detail.Description)
			builder.WriteString("\n")
		}
	}

	h.l.Logf(builder.String())
}

func (h *internalEvents) OnError(evt any) {
	e := evt.(*OnError)

	h.l.Logf("\n%s %s <--- %s %s\n%s %s\n\n%s: %v",
		colorize.RedBright(colorize.Bold("ERROR")),
		time.Now().Format(time.RFC3339),
		colorize.Red(e.Request.Method),
		colorize.Red(e.Request.Path),
		e.Request.Method,
		e.Request.FullURL(),
		colorize.Red(colorize.Bold("Error: ")),
		e.Err,
	)
}
