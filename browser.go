package browser

import (
	"html/template"

	"github.com/spinsrv/browser/dom"
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

type Event interface{}

var Events = make(chan Event)

// call like:
// go browser.Dispatch(e)
// a hack to force a re-render is: go browser.Dispatch(nil)
func Dispatch(e Event) {
	Events <- e
}
