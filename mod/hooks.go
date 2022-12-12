package mod

import (
	"fmt"
	"net/http"
)

// EvtReq defines information from http.Request to be logged.
type EvtReq struct {
	Method     string
	Path       string
	RequestURI string
	Host       string
	Header     http.Header
	Body       []byte
}

// EvtRes defines HTTP EvtRes information to be logged.
type EvtRes struct {
	Status int
	Header http.Header
	Body   []byte
}

// EvtMk defines core.EvtMk information to be logged.
type EvtMk struct {
	ID   int
	Name string
}

// EvtResult defines matching result to be logged.
type EvtResult struct {
	HasClosestMatch bool
	ClosestMatch    EvtMk
	Details         []EvtResultExt
}

// EvtResultExt defines matching result details to be logged.
type EvtResultExt struct {
	Name        string
	Target      string
	Description string
}

func (r *EvtReq) FullURL() string {
	return fmt.Sprintf("%s%s", r.Host, r.RequestURI)
}
