package browser

import (
	"html/template"
	"time"

	"github.com/nlandolfi/browser/dom"
)

type Browser interface {
	Document() dom.Document
	NewElement(template.HTML) dom.Element
}

type LocalStorage interface {
	Put(k, v string)
	Get(k string) string
	Del(k string)
}

// Events is the global events channel.
var Events = make(chan Event)

// Event is a helper type for events.
//
// Usually code type switches on these for handling.
type Event interface{}

// Use Dispatch to add an event to the queue. Usually you want to
// do so in in a non-blocking way, in which case, use the `go` keyword
//
//   go browser.Dispatch(e)
//
// A simple 'hack' to force a re-render is to dispatch a nil event
//
//   go browser.Dispatch(nil)
//
func Dispatch(e Event) {
	Events <- e
}

func DispatchAfter(e Event, d time.Duration) {
	time.Sleep(d)
	Dispatch(e)
}

func DispatchAt(e Event, t time.Time) {
	d := t.Sub(time.Now())
	if d < 0 {
		panic("can not dispatch in the past")
	}
	DispatchAfter(e, d)
}
