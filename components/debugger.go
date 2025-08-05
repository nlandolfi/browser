package components

import "github.com/nlandolfi/browser"

type Handler interface {
	Handle(browser.Event)
}

type State struct {
	Events []browser.Event

	Handler Handler
}

func (s *State) Handle(e browser.Event) {
	s.Events = append(s.Events, e)
}
