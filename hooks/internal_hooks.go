package hooks

import (
	"fmt"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
)

// Logger defines internal events logger contract.
type Logger interface {
	Logf(string, ...any)
}

// InternalEvents implements default event handlers that logs event information.
type InternalEvents struct {
	l Logger
}

// NewInternalEvents creates an internal event handlers.
func NewInternalEvents(l Logger) *InternalEvents {
	return &InternalEvents{l: l}
}

func (h *InternalEvents) OnRequest(e OnRequest) {
	h.l.Logf("\n%s %s ---> %s %s\n%s %s\n\n%s:\n %s: %v\n",
		colorize.BlueBright(colorize.Bold("REQUEST RECEIVED")),
		e.StartedAt.Format(time.RFC3339),
		colorize.Blue(e.Request.Method),
		colorize.Blue(e.Request.Path),
		e.Request.Method,
		fullURL(e.Request.Host, e.Request.RequestURI),
		colorize.Blue("Request"),
		colorize.Blue("Headers"),
		e.Request.Header,
	)
}

func (h *InternalEvents) OnRequestMatched(e OnRequestMatch) {
	h.l.Logf("\n%s %s <--- %s %s\n%s %s\n\n%s%d %s\n\n%s: %dms\n%s:\n %s: %d\n %s: %v\n",
		colorize.GreenBright(colorize.Bold("REQUEST DID MATCH")),
		time.Now().Format(time.RFC3339),
		colorize.Green(e.Request.Method),
		colorize.Green(e.Request.Path),
		e.Request.Method,
		fullURL(e.Request.Host, e.Request.RequestURI),
		colorize.Bold("Mock: "),
		e.Mock.ID,
		e.Mock.Name,
		colorize.Green("Took"),
		e.Elapsed.Milliseconds(),
		colorize.Green("Response Definition"),
		colorize.Green("Status"),
		e.ResponseDefinition.Status,
		colorize.Green("Headers"),
		e.ResponseDefinition.Header,
	)
}

func (h *InternalEvents) OnRequestNotMatched(e OnRequestNotMatched) {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("\n%s %s <--- %s %s\n%s %s\n\n",
		colorize.YellowBright(colorize.Bold("REQUEST DID NOT MATCH")),
		time.Now().Format(time.RFC3339),
		colorize.Yellow(e.Request.Method),
		colorize.Yellow(e.Request.Path),
		e.Request.Method,
		fullURL(e.Request.Host, e.Request.RequestURI)))

	if e.Result.HasClosestMatch {
		builder.WriteString(fmt.Sprintf("%s: %d %s\n\n",
			colorize.Bold("Closest Match"), e.Result.ClosestMatch.ID, e.Result.ClosestMatch.Name))
	}

	builder.WriteString(fmt.Sprintf("%s:\n", colorize.Bold("Mismatches")))

	for _, detail := range e.Result.Details {
		builder.WriteString(fmt.Sprintf("%s, reason=%s, applied-to=%s\n",
			colorize.Bold(detail.Name), detail.Description, detail.Target))
	}

	h.l.Logf(builder.String())
}

func (h *InternalEvents) OnError(e OnError) {
	h.l.Logf("\n%s %s <--- %s %s\n%s %s\n\n%s: %v",
		colorize.RedBright(colorize.Bold("REQUEST DID NOT MATCH")),
		time.Now().Format(time.RFC3339),
		colorize.Red(e.Request.Method),
		colorize.Red(e.Request.Path),
		e.Request.Method,
		fullURL(e.Request.Host, e.Request.RequestURI),
		colorize.Red(colorize.Bold("Error: ")),
		e.Err,
	)
}
