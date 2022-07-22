package events

import (
	"fmt"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/internal/colorize"
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

func (h *InternalEvents) OnRequestMatch(e OnRequestMatch) {
	h.l.Logf("\n%s %s <--- %s %s\n%s %s\n\n%s%d - %s\n\n%s: %dms\n%s:\n %s: %d\n %s: %v\n",
		colorize.GreenBright(colorize.Bold("REQUEST DID MATCH")),
		time.Now().Format(time.RFC3339),
		colorize.Green(e.Request.Method),
		colorize.Green(e.Request.Path),
		e.Request.Method,
		fullURL(e.Request.Host, e.Request.RequestURI),
		colorize.Bold("AddMocks: "),
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
		builder.WriteString("Closest Match:\n")
		builder.WriteString(
			fmt.Sprintf("id: %d\nname: %s\n\n", e.Result.ClosestMatch.ID, e.Result.ClosestMatch.Name))
	}

	builder.WriteString("Mismatches:\n")

	for _, detail := range e.Result.Details {
		builder.WriteString(fmt.Sprintf("%s\n%s\n%s",
			fmt.Sprintf("Matcher \"%s\"", colorize.Bold(detail.Name)),
			"Target: "+detail.Target,
			detail.Description))
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
