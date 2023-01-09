package event

import (
	"net/http"
)

// EvtReq defines information from http.Request to be logged.
type EvtReq struct {
	Method string
	Path   string
	Host   string
	Header http.Header
	Body   []byte
	URL    string
}

// EvtRes defines HTTP EvtRes information to be logged.
type EvtRes struct {
	Status int
	Header http.Header
	Body   []byte
}

// EvtMk defines core.EvtMk information to be logged.
type EvtMk struct {
	ID   string
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
