package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
)

type L interface {
	Logf(format string, a ...any)
}

type InternalListener struct {
	l       L
	verbose bool
}

func NewInternalListener(l L, verbose bool) *InternalListener {
	return &InternalListener{l: l, verbose: verbose}
}

func (h *InternalListener) OnRequest(evt any) {
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

	if h.verbose {
		builder.WriteString("\n")

		if len(e.Request.Body) > 0 {
			builder.WriteString(colorize.Blue("Body: "))
			builder.WriteString(fmt.Sprintf("%v\n", string(e.Request.Body)))
		}
	}

	h.l.Logf(builder.String())
}

func (h *InternalListener) OnRequestMatched(evt any) {
	e := evt.(*OnRequestMatch)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("\n%s %s <--- %s %s\n%s %s\n",
		colorize.GreenBright(colorize.Bold("REQUEST DID MATCH")),
		time.Now().Format(time.RFC3339),
		colorize.Green(e.Request.Method),
		colorize.Green(e.Request.Path),
		e.Request.Method,
		e.Request.FullURL()))

	nm := e.Mock.Name
	if nm == "" {
		nm = "<unnamed>"
	}

	if h.verbose {
		builder.WriteString(fmt.Sprintf("\n%s %s %s\n%s %dms\n%s\n %s %d\n %s %v\n",
			colorize.Bold("Mock:"),
			e.Mock.ID,
			nm,
			colorize.Green("Took:"),
			e.Elapsed.Milliseconds(),
			colorize.Green("Response"),
			colorize.Green("Status:"),
			e.ResponseDefinition.Status,
			colorize.Green("Headers:"),
			e.ResponseDefinition.Header))

		if len(e.ResponseDefinition.Body) > 0 {
			builder.WriteString(
				fmt.Sprintf(" %s %s\n", colorize.Green("Body:"), string(e.ResponseDefinition.Body)))
		}
	}

	h.l.Logf(builder.String())
}

func (h *InternalListener) OnRequestNotMatched(evt any) {
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
		builder.WriteString(fmt.Sprintf("%s: %s %s\n\n",
			colorize.Bold("Closest Match"), e.Result.ClosestMatch.ID, e.Result.ClosestMatch.Name))
	}

	if h.verbose {
		builder.WriteString(fmt.Sprintf("%s:\n", colorize.Bold("Mismatches")))

		for _, detail := range e.Result.Details {
			builder.WriteString(detail.Description)
			builder.WriteString("\n")
		}
	}

	h.l.Logf(builder.String())
}

func (h *InternalListener) OnError(evt any) {
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
